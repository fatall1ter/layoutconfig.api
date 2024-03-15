package infra

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	"git.countmax.ru/countmax/layoutconfig.api/internal/permission"

	"github.com/labstack/echo/v4"
)

// NewChainStoreChainStore is retail object abstraction of store/shop/renter
type NewChainStore struct {
	LayoutID   string     `json:"layout_id,omitempty"`
	Kind       string     `json:"kind,omitempty"`
	Title      string     `json:"title,omitempty"`
	CRMKey     string     `json:"crm_key,omitempty"`
	Brands     string     `json:"brands,omitempty"`
	LocationID string     `json:"location_id,omitempty"`
	Area       float64    `json:"area,omitempty"`
	Currency   string     `json:"currency,omitempty"`
	Options    string     `json:"options,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
}

func (ncs *NewChainStore) makeDomainChainStore() domain.ChainStore {
	return domain.ChainStore{
		LayoutID:   ncs.LayoutID,
		Kind:       ncs.Kind,
		Title:      ncs.Title,
		CRMKey:     ncs.CRMKey,
		Brands:     ncs.Brands,
		LocationID: ncs.LocationID,
		Area:       ncs.Area,
		Currency:   ncs.Currency,
		Options:    ncs.Options,
		Notes:      ncs.Notes,
		ValidFrom:  ncs.ValidFrom,
		ValidTo:    ncs.ValidTo,
		Creator:    ncs.Creator,
	}
}

// nolint:lll
// apiCreateChainStore docs
// @Summary Create new store in the retail schema
// @Description create new store in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example, use layout_id for exists layout
// @Description {
// @Description    "area": 109.94,
// @Description    "brands": "{\\"ru\\":[\\"Брэнд1\\"], \\"en\\":[\\"Brand1\\"]}",
// @Description	   "creator": "username",
// @Description	   "crm_key": "code1s for store",
// @Description	   "currency": "rub",
// @Description	   "kind": "store",
// @Description	   "layout_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description	   "notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description	   "options": "{\\"tz\\":3}",
// @Description	   "title": "{\\"ru\\":\\"Наименование Магазина на русском языке\\",\\"en\\":\\"Store name in English\\"}",
// @Description	   "valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/stores
// @Param layout_id query string false "uuid format, default=*"
// @Param store body infra.NewChainStore true "store properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/stores [post]
func (s *Server) apiCreateChainStore(c echo.Context) error {
	store := &NewChainStore{}
	if err := c.Bind(store); err != nil {
		s.log.Errorf("apiCreateChainStore, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	id, err := repo.AddChainStore(c.Request().Context(), store.makeDomainChainStore())
	if err != nil {
		s.log.Errorf("repo.AddChainStore error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	href := `{"href":"/v2/chains/stores/` + id + `"}`
	return c.JSON(http.StatusCreated, CreatedStatus(href))
}

// ChainStoresResponse http wrapper with metadata
type ChainStoresResponse struct {
	Data domain.ChainStores `json:"data"`
	Metadata
}

// apiChainStores docs
// @Summary Get all stores in the retail schema
// @Description get slice of stores with loc, date, layout_id, crm_key, offset, limit, fields, include parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description store_ids - comma separated list of the store's ids, default *
// @Description fields - comma separated values of field names, can be store_id,layout_id,kind,title...
// @Description all of them described at the model
// @Description include - comma separated list of entities, embedded in current,
// @Description for store it can be entrances,zones,devices
// @Produce  json
// @Tags chains/stores
// @Param layout_id query string false "default=*"
// @Param store_ids query string false "default=*"
// @Param crm_key query string false "default=*"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "store_id,title... default=all"
// @Param include query string false "entrances,zones,devices default=none"
// @Success 200 {object} infra.ChainStoresResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/stores [get]
func (s *Server) apiChainStores(c echo.Context) error {
	inputIDs := c.QueryParam("store_ids")
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	crmKey := c.QueryParam("crm_key")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	filteredList, err := s.perm.FilteredStores(c.Request(), layoutID, inputIDs, acl.ActionRead)
	if err != nil {
		s.log.Errorf("perm.FilteredStores error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errors.New("permission process failed")))
	}
	if filteredList == "" {
		response := ChainStoresResponse{
			Data: domain.ChainStores{},
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
	stores, count, err := repo.FindChainStores(c.Request().Context(), loc, dt, layoutID, crmKey, offset, limit, filteredList)
	if err != nil {
		s.log.Errorf("repo.FindChainStores error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if stores == nil {
		response := ChainStoresResponse{
			Data: stores,
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
	stores.SetZeroValue(fields)
	includes := strings.Split(c.QueryParam("include"), ",")
	for _, include := range includes {
		var lOffset, lLimit int64 = 0, 999999
		switch include {
		case "entrances":
			var shortEntranceFieldSet []string = []string{"entrance_id", "store_id", "kind", "title"}
			entrances, _, err := repo.FindChainEntrances(c.Request().Context(), loc, dt, layoutID, "*", "*", "*", lOffset, lLimit)
			if err == nil {
				entrances.SetZeroValue(shortEntranceFieldSet)
				stores.IncludeEntrances(entrances)
			}
		case "zones":
			var shortZoneFieldSet []string = []string{"zone_id", "store_id", "kind", "title"}
			zones, _, err := repo.FindChainZones(c.Request().Context(), loc, dt, layoutID, "*", "*", "*", "*", "*", lOffset, lLimit)
			if err == nil {
				zones.SetZeroValue(shortZoneFieldSet)
				stores.IncludeZones(zones)
			}
		case "devices":
			var shortDeviceFieldSet []string = []string{"device_id", "store_id", "kind", "ip", "port"}
			devices, _, err := repo.FindChainDevices(c.Request().Context(), loc, dt, layoutID, "*", "*", "*", "*", "*", "*", lOffset, lLimit)
			if err == nil {
				devices.SetZeroValue(shortDeviceFieldSet)
				stores.IncludeDevices(devices)
			}
		}

	}
	//
	response := ChainStoresResponse{
		Data: stores,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(stores)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiChainStoreByID docs
// @Summary Get specified store in the retail schema
// @Description get store with loc(ation), date, store_id, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names,
// @Description can be store_id,layout_id,kind,title... all of them described at the model
// @Description include - comma separated list of entities,
// @Description embedded in current, for store it can be entrances,zones,devices
// @Produce  json
// @Tags chains/stores
// @Param layout_id query string false "default=*"
// @Param store_id path string true "uuid format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "store_id,title... default=all"
// @Param include query string false "entrances,zones,devices default=none"
// @Success 200 {object} domain.ChainStore
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/stores/{store_id} [get]
func (s *Server) apiChainStoreByID(c echo.Context) error {
	id := c.Param("store_id")
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	store, err := repo.FindChainStoreByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindChainStoreByID error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if store == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have store", id, loc, dt)
		return c.JSON(http.StatusOK, store)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	store.SetZeroValue(fields)
	includes := strings.Split(c.QueryParam("include"), ",")
	for _, include := range includes {
		var lOffset, lLimit int64 = 0, 999999
		switch include {
		case "entrances":
			var shortEntranceFieldSet []string = []string{"entrance_id", "store_id", "kind", "title"}
			entrances, _, err := repo.FindChainEntrances(c.Request().Context(), loc, dt, layoutID, store.StoreID, "*", "", lOffset, lLimit)
			if err == nil {
				entrances.SetZeroValue(shortEntranceFieldSet)
				store.Entrances = entrances
			}
		case "zones":
			var shortZoneFieldSet []string = []string{"zone_id", "store_id", "kind", "title"}
			zones, _, err := repo.FindChainZones(c.Request().Context(), loc, dt, layoutID, store.StoreID, "*", "*", "*", "*", lOffset, lLimit)
			if err == nil {
				zones.SetZeroValue(shortZoneFieldSet)
				store.Zones = zones
			}
		case "devices":
			var shortDeviceFieldSet []string = []string{"device_id", "store_id", "kind", "ip", "port"}
			devices, _, err := repo.FindChainDevices(c.Request().Context(), loc, dt, layoutID, store.StoreID, "*", "*", "*", "*", "*", lOffset, lLimit)
			if err == nil {
				devices.SetZeroValue(shortDeviceFieldSet)
				store.Devices = devices
			}
		}
	}
	return c.JSON(http.StatusOK, store)
}

// apiDeleteStore docs
// @Summary Delete specified store in the retail schema
// @Description delete store by store_id parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/stores
// @Param layout_id query string false "default=*"
// @Param store_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/stores/{store_id} [delete]
func (s *Server) apiDeleteChainStore(c echo.Context) error {
	id := c.Param("store_id")
	if id == "" {
		s.log.Errorf("bad request apiDeleteStore, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainStore(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.DelChainStore error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d row(s)", cnt)))
}

// UpdChainStore is retail object abstraction of store/shop/renter
type UpdChainStore struct {
	StoreID    string     `json:"store_id,omitempty"`
	LayoutID   string     `json:"layout_id,omitempty"`
	Title      string     `json:"title,omitempty"`
	CRMKey     string     `json:"crm_key,omitempty"`
	Brands     string     `json:"brands,omitempty"`
	LocationID string     `json:"location_id,omitempty"`
	Area       float64    `json:"area,omitempty"`
	Currency   string     `json:"currency,omitempty"`
	Options    string     `json:"options,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	NoHistory  bool       `json:"no_history"`
}

func (ncs *UpdChainStore) makeDomainChainStore() domain.ChainStore {
	return domain.ChainStore{
		StoreID:    ncs.StoreID,
		LayoutID:   ncs.LayoutID,
		Title:      ncs.Title,
		CRMKey:     ncs.CRMKey,
		Brands:     ncs.Brands,
		LocationID: ncs.LocationID,
		Area:       ncs.Area,
		Currency:   ncs.Currency,
		Options:    ncs.Options,
		Notes:      ncs.Notes,
		ValidFrom:  ncs.ValidFrom,
	}
}

// nolint:lll
// apiUpdChainStore docs
// @Summary Update store in the retail schema
// @Description update a store with upd_store parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description no_history can be true/false; if true = without save history, if false = save history of changes
// @Description valid_from set if need save history for specified date
// @Description example, store_id must to exists
// @Description {
// @Description    "store_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description    "area": 109.94,
// @Description	   "crm_key": "code1s for store",
// @Description	   "options": "{\\"tz\\":3}",
// @Description	   "title": "{\\"ru\\":\\"Наименование Магазина на русском языке\\",\\"en\\":\\"Store name in English\\"}",
// @Description	   "valid_from": "2020-04-01T00:00:00+03:00",
// @Description	   "no_history": true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/stores
// @Param layout_id query string false "default=*"
// @Param store body infra.UpdChainStore true "upd_store properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/stores [put]
func (s *Server) apiUpdChainStore(c echo.Context) error {
	store := &UpdChainStore{}
	if err := c.Bind(store); err != nil {
		s.log.Errorf("apiUpdChainStore, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if store.StoreID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty store_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	count, err := repo.UpdChainStore(c.Request().Context(), store.makeDomainChainStore(),
		store.NoHistory, c.Request().Header.Get(permission.XUserID)) // TODO: add method extract user id
	if err != nil {
		s.log.Errorf("repo.DelChainStore error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("update %d row(s)", count)))
}
