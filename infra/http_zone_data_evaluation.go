package infra

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// nolint:lll
// apiChainStoresDataAttendance docs
// @Summary Get zone data evaluation for serviceChannel blocks
// @Description get zone data evaluation for serviceChannel blocks by intervals
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description can be filtered by layout_id, store_id, parent_zone_id and daterange parameters from and to.
// @Description service_channel_block_id is serviceChannel block identifier
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss then the server's local timezone is applied
// @Produce  json
// @Tags data/queue
// @Param layout_id path string true "uuid format, default=*"
// @Param store_id query string false "default=*"
// @Param service_channel_block_id query string false "default=*"
// @Param is_full query string false "default=*"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Success 200 {object} domain.ZoneDataEvaluations
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/queue/evaluations [get]
func (s *Server) apiZoneDataEvaluation(c echo.Context) error {
	pStoreID := c.QueryParam("store_id")
	pSCB := c.QueryParam("service_channel_block_id")
	pIsFull := c.QueryParam("is_full")
	from, to, err := s.getFromToParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, pLayoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: CheckLayout read
	data, err := repo.FindZoneDataEvaluation(c.Request().Context(), from, to, pLayoutID, pStoreID, pSCB, pIsFull)
	if err != nil {
		s.log.Errorf("repo.FindZoneDataEvaluation error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}
