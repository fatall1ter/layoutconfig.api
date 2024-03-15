package infra

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// NewChainZone is area of some physical territory, for creating
type NewChainZone struct {
	ParentID  *string    `json:"parent_id,omitempty"`
	LayoutID  string     `json:"layout_id,omitempty"`
	StoreID   string     `json:"store_id,omitempty"`
	Kind      string     `json:"kind,omitempty"`
	Title     string     `json:"title,omitempty"`
	Options   string     `json:"options,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
	Creator   string     `json:"creator,omitempty"`
}

func (ncz *NewChainZone) makeDomainChainZone() domain.ChainZone {
	return domain.ChainZone{
		ParentID:  ncz.ParentID,
		LayoutID:  ncz.LayoutID,
		StoreID:   ncz.StoreID,
		Kind:      ncz.Kind,
		Title:     ncz.Title,
		Options:   ncz.Options,
		Notes:     ncz.Notes,
		ValidFrom: ncz.ValidFrom,
		ValidTo:   ncz.ValidTo,
		Creator:   ncz.Creator,
	}
}

// nolint:lll
// apiCreateChainZone docs
// @Summary Create new zone in the retail schema
// @Description creates new zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new zone, layout_id and store_id must exists in the database:
// @Description {
// @Description	   "creator": "username",
// @Description	   "kind": "zone",
// @Description	   "layout_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description	   "notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description	   "options": "{\\"is_online\\":1,\\"borders\\":[{\\"title\\":\\"lowlevel\\",\\"low\\":0,\\"high\\":37,\\"color\\":\\"#008000\\"},{\\"title\\":\\"middle\\",\\"low\\":37,\\"high\\":52,\\"color\\":\\"#FFFF80\\"},{\\"title\\":\\"high\\",\\"low\\":52,\\"high\\":70,\\"color\\":\\"#FF0000\\"}]}",
// @Description	   "store_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description    "title": "{\\"ru\\":\\"Наименование Зоны на русском языке\\",\\"en\\":\\"Zone name in English\\"}",
// @Description    "valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string true "default=*"
// @Param zone body infra.NewChainZone true "zone properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones [post]
func (s *Server) apiCreateChainZone(c echo.Context) error {
	zone := &NewChainZone{}
	if err := c.Bind(zone); err != nil {
		s.log.Errorf("apiCreateChainZone, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	id, err := repo.AddChainZone(c.Request().Context(), zone.makeDomainChainZone())
	if err != nil {
		s.log.Errorf("repo.AddChainZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	href := `{"href":"/v2/chains/zones/` + id + `"}`
	return c.JSON(http.StatusCreated, CreatedStatus(href))
}

// ChainZonesResponse http wrapper with metadata
type ChainZonesResponse struct {
	Data domain.ChainZones `json:"data"`
	Metadata
}

// nolint:lll
// apiChainZones docs
// @Summary Get all zones in the retail schema
// @Description get slice of zones with loc, date, layout_id, store_id, is_online, is_active, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field zone_id,parent_id,layout_id,store_id,kind,title... all of them described at the model
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param kind query string false "default=*"
// @Param parent_id query string false "default=*"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param is_online query string false "default=*"
// @Param is_active query string false "default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "zone_id,kind,title... default=all"
// @Success 200 {object} infra.ChainZonesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones [get]
func (s *Server) apiChainZones(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	storeID := c.QueryParam("store_id") // TODO: add listIDs param
	kind := c.QueryParam("kind")
	parentID := c.QueryParam("parent_id")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	isOnline := c.QueryParam("is_online")
	isActive := c.QueryParam("is_active")
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: FilteredStores create
	zones, count, err := repo.FindChainZones(c.Request().Context(),
		loc, dt, layoutID, storeID, kind, parentID, isOnline, isActive, offset, limit) // TODO: add store list param
	if err != nil {
		s.log.Errorf("repo.FindChainZones error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if zones == nil {
		response := ChainZonesResponse{
			Data: zones,
			Metadata: Metadata{
				ResultSet: ResultSet{
					Count:  0,
					Offset: offset,
					Limit:  limit,
					Total:  count,
				},
			},
		}
		return c.JSON(http.StatusOK, response)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	zones.SetZeroValue(fields)
	response := ChainZonesResponse{
		Data: zones,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(zones)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// nolint:lll
// apiChainZoneByID docs
// @Summary Get specified zone in the retail schema
// @Description get zone with loc, date, zone_id, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field zone_id,parent_id,layout_id,store_id,kind,title... all of them described at the model
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string false "default=*"
// @Param zone_id path string true "uuid format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "zone_id,kind,title... default=all"
// @Success 200 {object} domain.ChainZone
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones/{zone_id} [get]
func (s *Server) apiChainZoneByID(c echo.Context) error {
	id := c.Param("zone_id")
	if id == "" {
		s.log.Errorf("bad request apiChainZoneByID, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add ZoneCheck read
	zone, err := repo.FindChainZoneByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindChainZoneByID error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if zone == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have zone", id, loc, dt)
		return c.JSON(http.StatusOK, zone)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	zone.SetZeroValue(fields)
	return c.JSON(http.StatusOK, zone)
}

// apiDeleteZone docs
// @Summary Delete specified zone in the retail schema
// @Description delete zone by zone_id parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string false "default=*"
// @Param zone_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones/{zone_id} [delete]
func (s *Server) apiDeleteChainZone(c echo.Context) error {
	id := c.Param("zone_id")
	if id == "" {
		s.log.Errorf("bad request apiDeleteChainZone, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add CheckLayout delete
	cnt, err := repo.DelChainZone(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.DelChainZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d zone(s)", cnt)))
}

// UpdChainZone is area of some physical territory, for creating
type UpdChainZone struct {
	ZoneID    string     `json:"zone_id,omitempty"`
	ParentID  *string    `json:"parent_id,omitempty"`
	Kind      string     `json:"kind,omitempty"`
	Title     string     `json:"title,omitempty"`
	Options   string     `json:"options,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	NoHistory bool       `json:"no_history"`
}

func (ncz *UpdChainZone) makeDomainChainZone() domain.ChainZone {
	return domain.ChainZone{
		ZoneID:    ncz.ZoneID,
		ParentID:  ncz.ParentID,
		Kind:      ncz.Kind,
		Title:     ncz.Title,
		Options:   ncz.Options,
		Notes:     ncz.Notes,
		ValidFrom: ncz.ValidFrom,
	}
}

// nolint:lll
// apiCreateChainZone docs
// @Summary Update zone in the retail schema
// @Description update exists zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description no_history can be true/false; if true = without save history, if false = save history of changes, default = false
// @Description valid_from set if need save history for specified date
// @Description example upd zone, zone_id must to be pass in parameters:
// @Description {
// @Description 	"zone_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description 	"kind": "zone",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"options": "{}",
// @Description 	"title": "{\\"ru\\":\\"Наименование Зоны на русском языке\\",\\"en\\":\\"Zone name in English\\"}",
// @Description 	"valid_from": "2020-06-09T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string false "default=*"
// @Param zone body infra.UpdChainZone true "zone properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones [put]
func (s *Server) apiUpdChainZone(c echo.Context) error {
	zone := &UpdChainZone{}
	if err := c.Bind(zone); err != nil {
		s.log.Errorf("apiUpdChainZone, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if zone.ZoneID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty zone_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add CheckLayout update
	count, err := repo.UpdChainZone(c.Request().Context(), zone.makeDomainChainZone(), zone.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdChainZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("update %d row(s)", count)))
}
