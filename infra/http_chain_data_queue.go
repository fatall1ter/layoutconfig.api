package infra

import (
	"fmt"
	"net/http"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/repos"
	"github.com/labstack/echo/v4"
)

const (
	defaultWindow int = 10
)

// apiChainStoresDataQueue docs
// @Summary Get data queue length for stores
// @Description get data on the number of people inside the queue in the stores by group intervals,
// @Description group can be interval means rawdata, 1m means group by minutes
// @Description agg_func applies to raw data for calculate rawdata, can be max, min, avg
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description can be filtered by store_id and daterange parameters from and to.
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00
// @Description or naive YYYY-MM-DD HH:mm:ss then the server's local timezone is applied
// @Produce  json
// @Tags data/queue
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param group_by query string false "default=interval"
// @Param agg_func query string false "default=max"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Success 200 {object} domain.StoresDataQueue
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/queue/length/stores [get]
func (s *Server) apiChainStoresDataQueue(c echo.Context) error {
	pStoreID := c.QueryParam("store_id") // TODO: add list of the ids stores as params
	groupBy := c.QueryParam("group_by")
	if groupBy == "" {
		groupBy = "interval"
	}
	groupFunc := c.QueryParam("agg_func")
	if groupFunc == "" {
		groupFunc = "max"
	}
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Warnf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add filtered list of the stores
	data, err := repo.FindChainStoresDataQueue(c.Request().Context(), from, to, &pStoreID, groupBy, groupFunc, defaultWindow)
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindChainStoresDataQueue error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiChainStoresDataQueue docs
// @Summary Get data queue length for stores in current moment
// @Description get data on the number of people inside the queue in the stores in current moment,
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description can be filtered by store_id
// @Produce  json
// @Tags data/queue
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Success 200 {object} domain.StoresDataQueue
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/queue/length/stores/live [get]
func (s *Server) apiChainStoresDataQueueNow(c echo.Context) error {
	pStoreID := c.QueryParam("store_id") // TODO: add list of the ids stores as params
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add filtered list of the stores
	data, err := repo.FindChainStoresDataQueueNow(c.Request().Context(), &pStoreID)
	if err != nil {
		s.log.Errorf("repo.FindChainStoresDataQueueNow error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiChainZonesDataQueue docs
// @Summary Get data queue length for zones of stores by intervals
// @Description get data on the number of people inside the queue in the zones of stores by intervals,
// @Description can be filtered by zone_id, store_id and daterange parameters from and to
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss
// @Description then the server's local timezone is applied
// @Produce  json
// @Tags data/queue
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param zone_id query string false "default=*"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Success 200 {object} domain.ZonesDataQueue
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/queue/length/zones [get]
func (s *Server) apiChainZonesDataQueue(c echo.Context) error {
	zoneID := c.QueryParam("zone_id") // TODO: add list of the ids zones as params
	storeID := c.QueryParam("store_id")
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Warnf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: get filtered list of the zones/detections
	data, err := repo.FindChainZonesDataQueue(c.Request().Context(), from, to, storeID, zoneID) // TODO: add param list zones
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindChainZonesDataQueue error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiChainZonesDataQueueNow docs
// @Summary Get data queue length for zones of stores in current moment
// @Description get data on the number of people inside the queue in the zones of stores in current moment,
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description can be filtered by zone_id
// @Produce  json
// @Tags data/queue
// @Param layout_id query string false "default=*"
// @Param zone_id query string false "default=*"
// @Success 200 {object} domain.ZonesDataQueue
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/queue/length/zones/live [get]
func (s *Server) apiChainZonesDataQueueNow(c echo.Context) error {
	zoneID := c.QueryParam("zone_id") // TODO: add list of the ids zones as params
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: get filtered list of the zones/detections
	data, err := repo.FindChainZonesDataQueueNow(c.Request().Context(), &zoneID) // TODO: add param list zones
	if err != nil {
		s.log.Errorf("repo.FindChainZonesDataQueueNow error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiChainRecommendationsQueue docs
// @Summary Get data recommendations
// @Description get data recommendations on the need to open service channels
// @Description to prevent queuing by time range and timerange
// @Description can be filtered by store_id, split_interval can fill intervals between original points
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss
// @Description then the server's local timezone is applied to is the upper bound,
// @Description but the recommendations have points in the future, usually up to 10 minutes ahead
// @Produce  json
// @Tags data/queue
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param split_interval query string false "default=60s"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Success 200 {object} domain.PredictionsQueue
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/queue/recommendations [get]
func (s *Server) apiChainRecommendationsQueue(c echo.Context) error {
	storeID := c.QueryParam("store_id")
	split := c.QueryParam("split_interval")
	var splitInterval time.Duration = 60 * time.Second
	if split != "" {
		splitParse, err := time.ParseDuration(split)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				ErrInvalidRequest(fmt.Errorf("wrong split_interval=%s parameter, need 1s/5s/15s/30s..., %v",
					split, err)))
		}
		splitInterval = splitParse
	}
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Warnf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckStore
	data, err := repo.FindChainPredictionQueue(c.Request().Context(), from, to, &storeID, splitInterval)
	if err != nil {
		s.log.Errorf("repo.FindChainPredictionQueue error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}
