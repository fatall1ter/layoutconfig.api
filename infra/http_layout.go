package infra

import (
	"errors"
	"net/http"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

var (
	errLayoutRepoNotFound error = errors.New("doesn't have layout repo")
)

// LayoutResponse http wrapper with metadata
type LayoutResponse struct {
	Data domain.Layouts `json:"data"`
	Metadata
}

// apiLayouts docs
// @Summary Get all layouts
// @Description get slice of layouts/projects configuration with location, date, offset, limit parameters
// @Produce  json
// @Produce  xml
// @Tags common
// @Param loc query string false "default=ru"
// @Param date query string false "YYYY-MM-DD, default=today"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Success 200 {object} infra.LayoutResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/layouts [get]
func (s *Server) apiLayouts(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	dt := c.QueryParam("date")
	if _, err := time.Parse("2006-01-02", dt); err != nil {
		dt = time.Now().Format("2006-01-02")
	}
	var count int64
	layouts := make(domain.Layouts, 0, limit)
	perm := s.perm.FromRequest(c.Request())
	for _, repo := range s.repoM.Repos() {
		_layouts, _count, err := repo.FindLayouts(c.Request().Context(), loc, dt, offset, limit)
		if err != nil {
			s.log.Errorf("repo.FindLayouts error %s", err)
			return s.responserMIME(c, http.StatusInternalServerError, ErrServerInternal(err))
		}
		currLen := int64(len(_layouts))
	LLOOP:
		for _, _layout := range _layouts { // TODO: check limit offset working!!!!
			if !perm.CheckLayout(_layout.ID, "read") {
				continue LLOOP
			}
			if currLen < limit {
				layouts = append(layouts, _layout)
				currLen++
			}
		}
		count += _count
	}
	//
	response := LayoutResponse{
		Data: layouts,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(layouts)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	s.log.Debugf("headers: %+v", c.Request().Header)
	return s.responserMIME(c, http.StatusOK, response)
}

// apiLayoutByID docs
// @Summary Get specified layout
// @Description get layout with location, date and layout_id parameters
// @Produce  json
// @Tags common
// @Param layout_id path string true "digit/uuid format"
// @Param loc query string false "default=ru"
// @Param date query string false "YYYY-MM-DD, default=today"
// @Success 200 {object} domain.Layout
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/layouts/{layout_id} [get]
func (s *Server) apiLayoutByID(c echo.Context) error {
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	dt := c.QueryParam("date")
	if _, err = time.Parse("2006-01-02", dt); err != nil {
		dt = time.Now().Format("2006-01-02")
	}
	layout, err := repo.FindLayoutByID(c.Request().Context(), loc, layoutID, dt)
	if err != nil {
		s.log.Errorf("repo.FindLayoutByID failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, layout)
}
