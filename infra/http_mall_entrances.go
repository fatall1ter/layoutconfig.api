package infra

import (
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// MallEntrancesResponse http wrapper with metadata
type MallEntrancesResponse struct {
	Data domain.MallEntrances `json:"data"`
	Metadata
}

// apiMallEntrances docs
// @Summary Get all entrances in the mall schema
// @Description get slice of entrances with loc, date, layout_id, floor_id, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be entrance_id,layout_id,floor_id,kind,title...
// @Produce  json
// @Tags malls/entrances
// @Param layout_id query string false "uuid/digit format, default=*"
// @Param floor_id query string false "default=*"
// @Param kind query string false "default=entrance"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param entrance_ids query string false "comma separated list ids, default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "entrance_id,floor_id,title... default=all"
// @Success 200 {object} infra.MallEntrancesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/entrances [get]
func (s *Server) apiMallEntrances(c echo.Context) error {
	//
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	floorID := c.QueryParam("floor_id")
	kind := c.QueryParam("kind")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	inputEnterIDs := c.QueryParam("entrance_ids")
	entrances, count, err := repo.FindMallEntrances(c.Request().Context(),
		loc, dt, layoutID, floorID, kind, inputEnterIDs, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindMallEntrances error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if entrances == nil {
		s.log.With(zap.String("layoutID", layoutID), zap.String("storeID", floorID)).
			Infof("kind=%s, loc=%s, dt=%s, offset=%d, limit=%d", kind, loc, dt, offset, limit)
		response := MallEntrancesResponse{
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

	response := MallEntrancesResponse{
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

// apiMallEntranceByID docs
// @Summary Get specified entrance in the retail schema
// @Description get specified entrance with specified, loc, date, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be entrance_id,layout_id,floor_id,kind,title...
// @Produce  json
// @Tags malls/entrances
// @Param layout_id query string false "uuid/digit format, default=*"
// @Param entrance_id path string true "uuid format"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "entrance_id,floor_id,title... default=all"
// @Success 200 {object} domain.MallEntrance
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/entrances/{entrance_id} [get]
func (s *Server) apiMallEntranceByID(c echo.Context) error {
	id := c.Param("entrance_id")
	if id == "" {
		s.log.Errorf("bad request apiMallEntranceByID, %v", errEmptyID)
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
		s.log.Errorf("repo.FindChainEntranceByID error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if entrance == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have entrance", id, loc, dt)
		return c.JSON(http.StatusOK, entrance)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	entrance.SetZeroValue(fields)
	return c.JSON(http.StatusOK, entrance)
}
