package infra

import (
	"fmt"
	"net/http"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// ChainZonesStatesResponse http wrapper with metadata
type ChainZonesStatesResponse struct {
	Data domain.ZoneStates `json:"data"`
	Metadata
}

// apiChainZonesStates docs
// @Summary Get zones (service_channel) states
// @Description get the zones states changes for zones with kind=service_channel, returns states online/offline
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param zone_id query string false "default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Success 200 {object} infra.ChainZonesStatesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones/states [get]
func (s *Server) apiChainZonesStates(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	storeID := c.QueryParam("store_id")
	zoneID := c.QueryParam("zone_id")
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Warnf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	zoneStates, count, err := repo.FindChainZonesStates(c.Request().Context(), layoutID, storeID, zoneID, from, to, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindChainZonesStates error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	response := ChainZonesStatesResponse{
		Data: zoneStates,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(zoneStates)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiChainZonesLastStates docs
// @Summary Get zones (service_channel) last states
// @Description get the zones last states for zones with kind=service_channel, returns states online/offline
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param zone_id query string false "default=*"
// @Success 200 {object} domain.ZoneStates
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones/states/last [get]
func (s *Server) apiChainZonesLastStates(c echo.Context) error {
	storeID := c.QueryParam("store_id")
	zoneID := c.QueryParam("zone_id")
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm LayoutCheck read only. X5 specified not 4 all layouts
	zoneStates, err := repo.FindChainZonesStatesLast(c.Request().Context(), layoutID, storeID, zoneID)
	if err != nil {
		s.log.Errorf("repo.FindChainZonesStatesLast error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, zoneStates)
}

// apiChainZonesLastStates docs
// @Summary Get zones (service_channel) last states
// @Description get the zones last states for zones with kind=service_channel, returns states online/offline
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/zones
// @Param layout_id query string true "default=*"
// @Param zone_id query string true "uuid format"
// @Param at query string false "ISO8601 datetime, default=current moment"
// @Success 200 {object} domain.ZoneState
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/zones/states/attime [get]
func (s *Server) apiChainZoneStateAtTime(c echo.Context) error {
	zoneID := c.QueryParam("zone_id")
	sAt := c.QueryParam("at")
	if sAt == "" {
		sAt = time.Now().Format(time.RFC3339)
	}
	_, toffset := time.Now().Zone()
	at, err := time.Parse(time.RFC3339, sAt)
	if err != nil {
		preErr := err
		toNaive, er := time.Parse(naiveTimeFormat, sAt)
		if er != nil {
			return c.JSON(http.StatusBadRequest,
				ErrInvalidRequest(fmt.Errorf("wrong at [%s] parameter %v/%v", sAt, preErr, er)))
		}
		at = toNaive.Local().Add(-time.Duration(toffset) * time.Second)
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm LayoutCheck read only. X5 specified not 4 all layouts
	zoneState, err := repo.FindChainZoneStateAtTime(c.Request().Context(), layoutID, zoneID, at)
	if err != nil {
		s.log.Errorf("repo.FindChainZoneStateAtTime error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, zoneState)
}
