package infra

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// NewChainEntrance is enter/exit entity and its properties for creation
type NewChainEntrance struct {
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

func (nce *NewChainEntrance) makeDomainChainEntrance() domain.ChainEntrance {
	return domain.ChainEntrance{
		LayoutID:  nce.LayoutID,
		StoreID:   nce.StoreID,
		Kind:      nce.Kind,
		Title:     nce.Title,
		Options:   nce.Options,
		Notes:     nce.Notes,
		ValidFrom: nce.ValidFrom,
		ValidTo:   nce.ValidTo,
		Creator:   nce.Creator,
	}
}

// nolint:lll
// apiCreateChainEntrance docs
// @Summary Create new entrance in the retail schema
// @Description creates new entrance
// @Description example new entrance, layout_id and store_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"kind": "entrance",
// @Description 	"layout_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"options": "{}",
// @Description 	"store_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description 	"title": "{\\"ru\\":\\"Наименование Входа на русском языке\\",\\"en\\":\\"Entrance name in English\\"}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param store body infra.NewChainEntrance true "entrance properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances [post]
func (s *Server) apiCreateChainEntrance(c echo.Context) error {
	enter := &NewChainEntrance{}
	if err := c.Bind(enter); err != nil {
		s.log.Errorf("apiCreateChainEntrance, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if enter.LayoutID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyLayoutID))
	}
	repo, ok := s.repoM.RepoByID(enter.LayoutID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errLayoutRepoNotFound))
	}
	id, err := repo.AddChainEntrance(c.Request().Context(), enter.makeDomainChainEntrance())
	if err != nil {
		s.log.Errorf("repo.AddChainEntrance failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	href := `{"href":"/v2/chains/entrances/` + id + `"}`
	return c.JSON(http.StatusCreated, CreatedStatus(href))
}

// ChainEntrancesResponse http wrapper with metadata
type ChainEntrancesResponse struct {
	Data domain.ChainEntrances `json:"data"`
	Metadata
}

// apiEntrances docs
// @Summary Get all entrances in the retail schema
// @Description get slice of entrances with loc, date, layout_id, store_id, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description entrance_ids comma separated list entrance ids
// @Description fields - comma separated values of field names, can be entrance_id,layout_id,store_id,kind,title
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param entrance_ids query string false "comma separated list ids, default=*"
// @Param kind query string false "default=entrance"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "entrance_id,store_id,title... default=all"
// @Success 200 {object} infra.ChainEntrancesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances [get]
func (s *Server) apiChainEntrances(c echo.Context) error {
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	storeID := c.QueryParam("store_id") // TODO: add list stores ids param
	kind := c.QueryParam("kind")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	inputEnterIDs := c.QueryParam("entrance_ids")
	filteredEnterIDsList, err := s.perm.FilteredEnters(c.Request(), layoutID, inputEnterIDs, acl.ActionRead)
	if err != nil {
		s.log.Errorf("perm.FilteredEnters error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errors.New("permission process failed")))
	}
	if filteredEnterIDsList == "" {
		response := ChainEntrancesResponse{
			Data: domain.ChainEntrances{},
			Metadata: Metadata{
				ResultSet: ResultSet{
					Count:  0,
					Offset: offset,
					Limit:  limit,
					Total:  0,
				},
			},
		}
		return c.JSON(http.StatusOK, response)
	}
	entrances, count, err := repo.FindChainEntrances(c.Request().Context(),
		loc, dt, layoutID, storeID, kind, filteredEnterIDsList, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindChainEntrances failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if entrances == nil {
		s.log.With(
			zap.String("layoutID", layoutID),
			zap.String("storeID", storeID)).
			Errorf("kind=%s, loc=%s, dt=%s, offset=%d, limit=%d", kind, loc, dt, offset, limit)
		response := ChainEntrancesResponse{
			Data: entrances,
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
	entrances.SetZeroValue(fields)

	response := ChainEntrancesResponse{
		Data: entrances,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(entrances)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiChainEntranceByID docs
// @Summary Get specified entrance in the retail schema
// @Description get specified entrance with specified, loc, date, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be entrance_id,layout_id,store_id,kind,title
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param entrance_id path string true "uuid format"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "entrance_id,store_id,title... default=all"
// @Success 200 {object} domain.ChainEntrance
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/{entrance_id} [get]
func (s *Server) apiChainEntranceByID(c echo.Context) error {
	id := c.Param("entrance_id")
	if id == "" {
		s.log.Errorf("bad request apiChainEntranceByID, %v", errEmptyID)
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
	entrance, err := repo.FindChainEntranceByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindChainEntranceByID failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if entrance == nil {
		s.log.Warnf("for id=%s, loc=%s and date=%s doesn't have entrance", id, loc, dt)
		return c.JSON(http.StatusOK, entrance)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	entrance.SetZeroValue(fields)
	return c.JSON(http.StatusOK, entrance)
}

// apiDeleteChainEntrance docs
// @Summary Delete specified entrance in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description delete entrance by entrance_id parameter
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param entrance_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/{entrance_id} [delete]
func (s *Server) apiDeleteChainEntrance(c echo.Context) error {
	id := c.Param("entrance_id")
	if id == "" {
		s.log.Errorf("bad request apiDeleteChainEntrance, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainEntrance(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.DelChainEntrance failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d entrance(s)", cnt)))
}

// apiCreateChainBindEntranceZone docs
// @Summary Create new binding entrance to zone in the retail schema
// @Description create new binding entrance to zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new binding entrance, entrance_id and zone_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"direction": "forward",
// @Description 	"entrance_id": "6fce2865-4b81-45ae-bdb6-9130a365b2b5",
// @Description 	"kind_zone": "zone",
// @Description 	"options": "{}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"zone_id": "9ea85f05-f02f-4b70-a464-ab85273471b7"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param params body domain.BindingChainEntranceZone true "entrance bind to zone parameters"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/bindto/zone [post]
func (s *Server) apiCreateChainBindEntranceZone(c echo.Context) error {
	bind := domain.BindingChainEntranceZone{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiCreateChainBindEntranceZone, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.BindChainEntranceZone(c.Request().Context(), bind)
	if err != nil {
		s.log.Errorf("repo.BindChainEntranceZone failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusCreated, CreatedStatus("binded"))
}

type UpdBindingChainEntranceZone struct {
	EntranceID string    `json:"entrance_id,omitempty"`
	OldZoneID  string    `json:"old_zone_id,omitempty"`
	NewZoneID  string    `json:"new_zone_id,omitempty"`
	Direction  string    `json:"direction,omitempty"`
	Options    string    `json:"options,omitempty"`
	ValidFrom  time.Time `json:"valid_from,omitempty"`
	NoHistory  bool      `json:"no_history"`
}

func (bez *UpdBindingChainEntranceZone) makeDomainBindingChainEntranceZone() domain.BindingChainEntranceZone {
	return domain.BindingChainEntranceZone{
		EntranceID: bez.EntranceID,
		ZoneID:     bez.OldZoneID,
		Direction:  bez.Direction,
		Options:    bez.Options,
		ValidFrom:  &bez.ValidFrom,
	}
}

// apiUpdChainBindEntranceZone docs
// @Summary Update binding entrance to zone in the retail schema
// @Description update binding entrance to zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description valid_from set if need save history for specified date
// @Description example upd binding old_zone_id is exists binding zone, new_zone_id is target zone
// @Description and entrance_id must to be pass in parameters:
// @Description {
// @Description 	"entrance_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"old_zone_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"new_zone_id": "3595c9ad-8116-408e-a007-ec31d48f9667",
// @Description 	"direction": "forward",
// @Description 	"options": "{}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param params body infra.UpdBindingChainEntranceZone true "entrance bind to zone parameters"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/bindto/zone [put]
func (s *Server) apiUpdChainBindEntranceZone(c echo.Context) error {
	bind := UpdBindingChainEntranceZone{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiUpdChainBindEntranceZone, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if bind.EntranceID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty entrance_id not allowed")))
	}
	if bind.OldZoneID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty old_zone_id not allowed")))
	}
	if bind.NewZoneID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty new_zone_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.UpdBindChainEntranceZone(c.Request().Context(), bind.makeDomainBindingChainEntranceZone(),
		bind.NewZoneID, bind.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdBindChainEntranceZone failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus("bind updated"))
}

// apiDeleteChainBindEntranceZone docs
// @Summary Delete binding entrance to zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description delete binding entrance to zone in the retail schema
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param entrance_id path string true "uuid format"
// @Param zone_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/bindto/zone [delete]
func (s *Server) apiDeleteChainBindEntranceZone(c echo.Context) error {
	entrance_id := c.QueryParam("entrance_id")
	zone_id := c.QueryParam("zone_id")
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainBindEntranceZone(c.Request().Context(), entrance_id, zone_id)
	if err != nil {
		s.log.Errorf("repo.DelChainBindEntranceZone failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d bind(s)", cnt)))
}

// apiCreateChainBindEntranceStore docs
// @Summary Create new binding entrance to store in the retail schema
// @Description create new binding entrance to store in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new binding entrance entrance_id and store_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"direction": "forward",
// @Description 	"entrance_id": "6fce2865-4b81-45ae-bdb6-9130a365b2b5",
// @Description 	"kind_store": "store",
// @Description 	"store_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"options": "{}"
// @Description   }
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param params body domain.BindingChainEntranceStore true "sensor bind entrance properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/bindto/store [post]
func (s *Server) apiCreateChainBindEntranceStore(c echo.Context) error {
	bind := domain.BindingChainEntranceStore{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiCreateChainBindEntranceStore, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.BindChainEntranceStore(c.Request().Context(), bind)
	if err != nil {
		s.log.Errorf("repo.BindChainEntranceStore, error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusCreated, CreatedStatus("binded"))
}

type UpdBindingChainEntranceStore struct {
	EntranceID string    `json:"entrance_id,omitempty"`
	OldStoreID string    `json:"old_store_id,omitempty"`
	NewStoreID string    `json:"new_store_id,omitempty"`
	Direction  string    `json:"direction,omitempty"`
	Options    string    `json:"options,omitempty"`
	ValidFrom  time.Time `json:"valid_from,omitempty"`
	NoHistory  bool      `json:"no_history"`
}

func (bez *UpdBindingChainEntranceStore) makeDomainBindingChainEntranceStore() domain.BindingChainEntranceStore {
	return domain.BindingChainEntranceStore{
		EntranceID: bez.EntranceID,
		StoreID:    bez.OldStoreID,
		Direction:  bez.Direction,
		Options:    bez.Options,
		ValidFrom:  &bez.ValidFrom,
	}
}

// apiUpdChainBindEntranceStore docs
// @Summary Update binding entrance to store in the retail schema
// @Description update binding entrance to store in the retail schema
// @Description valid_from set if need save history for specified date
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example upd binding old_store_id is exists binding store, new_store_id is target store
// @Description and entrance_id must to be pass in parameters:
// @Description {
// @Description 	"entrance_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"old_store_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"new_store_id": "3595c9ad-8116-408e-a007-ec31d48f9667",
// @Description 	"direction": "forward",
// @Description 	"options": "{}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param params body infra.UpdBindingChainEntranceStore true "entrance bind to store parameters"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/bindto/store [put]
func (s *Server) apiUpdChainBindEntranceStore(c echo.Context) error {
	bind := UpdBindingChainEntranceStore{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiUpdChainBindEntranceStore, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if bind.EntranceID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty entrance_id not allowed")))
	}
	if bind.OldStoreID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty old_store_id not allowed")))
	}
	if bind.NewStoreID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty new_store_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.UpdBindChainEntranceStore(c.Request().Context(), bind.makeDomainBindingChainEntranceStore(),
		bind.NewStoreID, bind.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdBindChainEntranceStore, error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus("bind updated"))
}

// apiDeleteChainBindEntranceStore docs
// @Summary Delete binding entrance to store in the retail schema
// @Description delete binding entrance to store in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param entrance_id path string true "uuid format"
// @Param store_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances/bindto/store [delete]
func (s *Server) apiDeleteChainBindEntranceStore(c echo.Context) error {
	entrance_id := c.QueryParam("entrance_id")
	store_id := c.QueryParam("store_id")
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainBindEntranceStore(c.Request().Context(), entrance_id, store_id)
	if err != nil {
		s.log.Errorf("repo.DelChainBindEntranceStore, error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d bind(s)", cnt)))
}

// UpdChainEntrance is enter/exit entity and its properties for update entrance
type UpdChainEntrance struct {
	EntranceID string     `json:"entrance_id,omitempty"`
	Kind       string     `json:"kind,omitempty"`
	Title      string     `json:"title,omitempty"`
	Options    string     `json:"options,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	NoHistory  bool       `json:"no_history"`
}

func (nce *UpdChainEntrance) makeDomainChainEntrance() domain.ChainEntrance {
	return domain.ChainEntrance{
		EntranceID: nce.EntranceID,
		Kind:       nce.Kind,
		Title:      nce.Title,
		Options:    nce.Options,
		Notes:      nce.Notes,
		ValidFrom:  nce.ValidFrom,
	}
}

// nolint:lll
// apiUpdChainEntrance docs
// @Summary Update entrance
// @Description update exists entrance in the retail schema
// @Description no_history can be true/false; if true = without save history, if false = save history of changes, default = false
// @Description valid_from set if need save history for specified date
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example upd entrance, entrance_id must to be pass in parameters:
// @Description {
// @Description 	"entrance_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description 	"kind": "entrance",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"options": "{}",
// @Description 	"title": "{\\"ru\\":\\"Наименование Входа на русском языке\\",\\"en\\":\\"Entrance name in English\\"}",
// @Description 	"valid_from": "2020-06-09T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/entrances
// @Param layout_id query string false "default=*"
// @Param upd_entrance body infra.UpdChainEntrance true "upd_entrance properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/entrances [put]
func (s *Server) apiUpdChainEntrance(c echo.Context) error {
	enter := &UpdChainEntrance{}
	if err := c.Bind(enter); err != nil {
		s.log.Errorf("apiUpdChainEntrance, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if enter.EntranceID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty entrance_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	count, err := repo.UpdChainEntrance(c.Request().Context(), enter.makeDomainChainEntrance(), enter.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdChainEntrance, error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("update %d entrance(s)", count)))
}
