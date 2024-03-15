package infra

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

const (
	userMock string = "layoutConfig_api_get_username_from_somewhere"
)

var (
	errEmptyID error = errors.New("empty id not allowed")
)

// ChainsResponse http wrapper with metadata
type ChainsResponse struct {
	Data domain.Chains `json:"data"`
	Metadata
}

// NewChain is parameter for new Layout of chain kind
type NewChain struct {
	Kind      string     `json:"kind,omitempty"`
	Title     string     `json:"title,omitempty"`
	Languages string     `json:"languages,omitempty"`
	CRMKey    string     `json:"crm_key,omitempty"`
	Brands    string     `json:"brands,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	Options   string     `json:"options,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
	Creator   string     `json:"creator,omitempty"`
	ReadOnly  bool       `json:"read_only,omitempty"`
}

func (nc *NewChain) makeChain() domain.Chain {
	return domain.Chain{
		Kind:      nc.Kind,
		Title:     nc.Title,
		Languages: nc.Languages,
		CRMKey:    nc.CRMKey,
		Brands:    nc.Brands,
		Currency:  nc.Currency,
		Options:   nc.Options,
		Notes:     nc.Notes,
		ValidFrom: nc.ValidFrom,
		ValidTo:   nc.ValidTo,
		Creator:   nc.Creator,
		ReadOnly:  nc.ReadOnly,
	}
}

// nolint:lll
// apiCreateChain docs
// @Summary Create new chain in the retail schema
// @Description create new chain in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new chain:
// @Description {
// @Description     "brands": "{\\"ru\\":[\\"Брэнд1\\",\\"Брэнд2\\"], \\"en\\":[\\"Brand1\\",\\"Brand2\\"]}",
// @Description     "creator": "username of creator",
// @Description     "crm_key": "code 1C",
// @Description     "currency": "rub/usd/eur...",
// @Description     "kind": "chain",
// @Description     "languages": "[\\"ru\\",\\"en\\"]",
// @Description     "notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description     "options": "{}",
// @Description     "read_only": false,
// @Description     "title": "{\\"ru\\":\\"Наименование сети на русском языке\\",\\"en\\":\\"Chain name in English\\"}",
// @Description     "valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Description valid_to можно указать когда заранее известна дата изменения состояния сети и ее атрибутов
// @Accept  json
// @Produce  json
// @Tags chains
// @Param layout_id query string false "default=*"
// @Param chain body infra.NewChain true "chain properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains [post]
func (s *Server) apiCreateChain(c echo.Context) error {
	newChain := &NewChain{}
	if err := c.Bind(newChain); err != nil {
		s.log.Errorf("apiCreateChain, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	// TODO: add check permissions
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckLayout create
	id, err := repo.AddChain(c.Request().Context(), newChain.makeChain())
	if err != nil {
		s.log.Errorf("AddChain(chain.title=%s) error %v", newChain.Title, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	href := `{"href":"/v2/chains/` + id + `"}`
	return c.JSON(http.StatusCreated, CreatedStatus(href))
}

// nolint:lll
// apiChains docs
// @Summary Get all chains in the retail schema
// @Description get slice of chains with loc (location), date, offset, limit, fields, include parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be layout_id,kind,title,languages... all of them described at the model
// @Description include - comma separated list of entities, embedded in current, for chain it can be stores
// @Produce  json
// @Produce  xml
// @Tags chains
// @Param layout_id query string false "default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD, default=today"
// @Param fields query string false "layout_id,title...default=all"
// @Param include query string false "stores default=none"
// @Success 200 {object} infra.ChainsResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains [get]
func (s *Server) apiChains(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}

	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	repo, layout, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		// return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
		return s.responserMIME(c, http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	perm := s.perm.FromRequest(c.Request())
	if !perm.CheckLayout(layout, "read") {
		// return c.JSON(http.StatusForbidden, ErrForbidden(nil))
		return s.responserMIME(c, http.StatusForbidden, ErrForbidden(nil))
	}
	// TODO: add perm CheckLayout read how layouts!!! for every chain
	chains, count, err := repo.FindChains(c.Request().Context(), loc, dt, layout, offset, limit)
	if err != nil {
		s.log.Errorf("FindChains error %v", err)
		return s.responserMIME(c, http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	chains.SetZeroValue(fields)

	include := c.QueryParam("include")
	if include == "stores" {
		var lOffset, lLimit int64 = 0, 999999
		var shortStoreFieldSet []string = []string{"store_id", "layout_id", "title", "crm_key"}
		stores, _, err := repo.FindChainStores(c.Request().Context(), loc, dt, layout, "*", lOffset, lLimit, "*")
		if err == nil {
			stores.SetZeroValue(shortStoreFieldSet)
			chains.IncludeStores(stores)
		}
	}

	response := ChainsResponse{
		Data: chains,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(chains)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}

	return s.responserMIME(c, http.StatusOK, response)
}

// apiChainByID docs
// @Summary Get specified chain in the retail schema
// @Description get chain with loc(ation), date, layout_id, fields, include parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be layout_id,kind,title,languages...
// @Description include - comma separated list of entities, embedded in current, for chain it can be stores
// @Produce  json
// @Tags chains
// @Param layout_id path string true "uuid format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "layout_id,title...default=all"
// @Param include query string false "stores default=none"
// @Success 200 {object} domain.Chain
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/{layout_id} [get]
func (s *Server) apiChainByID(c echo.Context) error {
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	repo, id, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	if id == "" {
		s.log.Errorf("bad request apiChainByID %s", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	chain, err := repo.FindChainByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("FindChainByID(%s, %s, %s) error %v", loc, id, dt, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if chain == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have chain", id, loc, dt)
		return c.JSON(http.StatusOK, chain)
	}
	chain.SetZeroValue(fields)
	include := c.QueryParam("include")
	if include == "stores" {
		var lOffset, lLimit int64 = 0, 999999
		var shortStoreFieldSet []string = []string{"store_id", "layout_id", "title", "crm_key"}
		stores, _, err := repo.FindChainStores(c.Request().Context(), loc, dt, chain.LayoutID, "*", lOffset, lLimit, "*")
		if err == nil {
			stores.SetZeroValue(shortStoreFieldSet)
			chain.Stores = stores
		} else {
			s.log.Errorf("FindChainStores(%s, %s, %s, %s) error %s",
				loc, dt, chain.LayoutID, "*, 0, 999999", err)
		}
	}
	return c.JSON(http.StatusOK, chain)
}

// apiDeleteChain docs
// @Summary Delete specified chain in the retail schema
// @Description delete chain by layout_id parameter
// @Produce  json
// @Tags chains
// @Param layout_id path string true "uuid format, default=*"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/{layout_id} [delete]
func (s *Server) apiDeleteChain(c echo.Context) error {
	repo, id, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	if id == "" {
		s.log.Errorf("bad request apiDeleteChain, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	// TODO: add perm CheckLayout delete
	cnt, err := repo.DelChain(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.DelChain(%s) error %s", id, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d chain(s)", cnt)))
}

// UpdChain is parameter for new Layout of chain kind
type UpdChain struct {
	LayoutID  string     `json:"layout_id"`
	Title     string     `json:"title,omitempty"`
	Languages string     `json:"languages,omitempty"`
	CRMKey    string     `json:"crm_key,omitempty"`
	Brands    string     `json:"brands,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	Options   string     `json:"options,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
	Creator   string     `json:"creator,omitempty"`
	NoHistory bool       `json:"no_history"`
}

func (nc *UpdChain) makeChain() domain.Chain {
	return domain.Chain{
		LayoutID:  nc.LayoutID,
		Title:     nc.Title,
		Languages: nc.Languages,
		CRMKey:    nc.CRMKey,
		Brands:    nc.Brands,
		Currency:  nc.Currency,
		Options:   nc.Options,
		Notes:     nc.Notes,
		ValidFrom: nc.ValidFrom,
		ValidTo:   nc.ValidTo,
		Creator:   nc.Creator,
	}
}

// apiUpdChain docs
// @Summary Update exists chain in the retail schema
// @Description updates a chain with upd_chain parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description no_history can be true/false; if true = without save history, if false = save history of changes
// @Description valid_from set if need save history for specified date
// @Description example upd_chain:
// @Description {
// @Description     "layout_id": "970195a2-e722-4222-94b9-c5266d37b1b8",
// @Description     "crm_key": "code 1C",
// @Description     "notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description     "title": "{\\"ru\\":\\"Новое наименование сети на русском языке\\",\\"en\\":\\"New chain name in English\\"}",
// @Description     "valid_from": "2020-06-19T00:00:00+03:00",
// @Description     "no_history":false
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains
// @Param layout_id path string true "uuid format, default=*"
// @Param chain body infra.UpdChain true "chain properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains [put]
func (s *Server) apiUpdChain(c echo.Context) error {
	updChain := &UpdChain{}
	if err := c.Bind(updChain); err != nil {
		s.log.Errorf("apiUpdChain, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if updChain.LayoutID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty layout_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckLayout update
	count, err := repo.UpdChain(c.Request().Context(), updChain.makeChain(), updChain.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdChain(chain.title=%s) error %s", updChain.Title, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("update %d chain(s)", count)))
}
