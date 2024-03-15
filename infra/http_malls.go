package infra

import (
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// MallsResponse http wrapper with metadata
type MallsResponse struct {
	Data domain.Malls `json:"data"`
	Metadata
}

// apiMalls docs
// @Summary Get all malls in the mall schema
// @Description get slice of malls with loc (location), date, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be layout_id,kind,title,languages...
// @Produce  json
// @Tags malls
// @Param layout_id query string false "default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD, default=today"
// @Param fields query string false "layout_id,title...default=all"
// @Success 200 {object} infra.MallsResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls [get]
func (s *Server) apiMalls(c echo.Context) error {
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
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckLayout read as layouts...
	malls, count, err := repo.FindMalls(c.Request().Context(), loc, dt, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindMalls error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	malls.SetZeroValue(fields)
	response := MallsResponse{
		Data: malls,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(malls)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiMallByID docs
// @Summary Get specified mall in the mall schema
// @Description get mall with loc(ation), date, layout_id, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be layout_id,kind,title,languages...
// @Produce  json
// @Tags malls
// @Param layout_id path string true "uuid format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "layout_id,title...default=all"
// @Success 200 {object} domain.Mall
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/{layout_id} [get]
func (s *Server) apiMallByID(c echo.Context) error {
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
		s.log.Errorf("bad request apiMallByID, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	// TODO: add CheckLayout read
	mall, err := repo.FindMallByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindMallByID(%s, %s, %s) error %v", loc, id, dt, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if mall == nil {
		s.log.Warnf("for id=%s, loc=%s and date=%s doesn't have mall", id, loc, dt)
		return c.JSON(http.StatusOK, mall)
	}
	mall.SetZeroValue(fields)
	return c.JSON(http.StatusOK, mall)
}
