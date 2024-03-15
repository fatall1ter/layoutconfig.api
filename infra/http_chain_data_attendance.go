package infra

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	"git.countmax.ru/countmax/layoutconfig.api/repos"
	"github.com/labstack/echo/v4"
)

const (
	defaultGroupBy  string = "interval"
	naiveTimeFormat string = "2006-01-02 15:04:05"
	boolFalse       string = "false"
	boolTrue        string = "true"
)

var (
	errInvalidDataSource error = errors.New("invalid datasource")
)

// apiChainStoresDataAttendance docs
// @Summary Get data attendance for stores
// @Description get data enters and exits in/out the stores by group_by intervals,
// @Description can be filtered by list store_id as comma separated list (12345,12344567,8488...)
// @Description and daterange parameters from and to.
// @Description intervals can be: interval, 1m, hour, day, week, month, quarter, year
// @Description use_rawdata for request data by sensors, can be: true/false/1/0; default false
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags data/attendance
// @Param layout_id query string false "default=*"
// @Param store_ids query string false "comma separated list store_id"
// @Param group_by query string false "default=interval"
// @Param use_rawdata query string false "default=false"
// @Param from query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=start today"
// @Param to query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=current time"
// @Success 200 {object} domain.StoresAttendance
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/attendance/stores [get]
func (s *Server) apiChainStoresDataAttendance(c echo.Context) error {
	storeIDs, err := s.joinParam(c, "store_id", "store_ids")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("%s", err)))
	}
	groupBy := c.QueryParam("group_by")
	if groupBy == "" {
		groupBy = defaultGroupBy
	}
	useRawData := c.QueryParam("use_rawdata")
	if useRawData == "" {
		useRawData = boolFalse
	}
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Errorf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	filteredList, err := s.perm.FilteredStores(c.Request(), layoutID, strings.Join(storeIDs, ","), acl.ActionRead)
	if err != nil {
		s.log.Errorf("perm.FilteredStores error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errors.New("permission process failed")))
	}
	if filteredList == "" {
		return c.JSON(http.StatusOK, domain.StoresAttendance{})
	}
	data, err := repo.FindChainStoresDataAttendance(c.Request().Context(),
		from, to, groupBy, layoutID, useRawData, filteredList)
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindChainStoresDataAttendance error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiChainZonesDataAttendance docs
// @Summary Get data attendance for zones of the stores
// @Description get data enters and exits in/out the zones of the stores by group_by intervals,
// @Description can be filtered by list zone_id as comma separated list (12345,12344567,8488...)
// @Description and daterange parameters from and to.
// @Description intervals can be: interval, hour, day, week, month, quarter, year
// @Description use_rawdata for request data by sensors, can be: true/false/1/0; default false
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags data/attendance
// @Param layout_id query string false "default=*"
// @Param zone_ids query string false "default=*"
// @Param group_by query string false "default=interval"
// @Param use_rawdata query string false "default=false"
// @Param from query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=start today"
// @Param to query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=current time"
// @Success 200 {object} domain.ZonesAttendance
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/attendance/stores/zones [get]
func (s *Server) apiChainZonesDataAttendance(c echo.Context) error {
	zoneIDs, err := s.joinParam(c, "zone_id", "zone_ids")
	if err != nil {
		s.log.Warnf("joinParam error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	groupBy := c.QueryParam("group_by")
	if groupBy == "" {
		groupBy = defaultGroupBy
	}
	useRawData := c.QueryParam("use_rawdata")
	if useRawData == "" {
		useRawData = boolFalse
	}
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Warnf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	//
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add get filtered list zones
	data, err := repo.FindChainZonesDataAttendance(c.Request().Context(), from, to, groupBy, layoutID, useRawData, strings.Join(zoneIDs, ","))
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindChainZonesDataAttendance error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiChainEntrancesDataAttendance docs
// @Summary Get data attendance for entrances of the stores
// @Description get data enters and exits in/out the entrances of the stores by group_by intervals,
// @Description entrance_ids - can be filtered by list entrance_id as comma separated list (12345,12344567,8488...)
// @Description and daterange parameters from and to.
// @Description intervals can be: interval, hour, day, week, month, quarter, year
// @Description use_rawdata for request data by sensors, can be: true/false/1/0; default false
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags data/attendance
// @Param layout_id query string false "default=*"
// @Param entrance_ids query string false "default=*"
// @Param group_by query string false "default=interval"
// @Param use_rawdata query string false "default=false"
// @Param from query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=start today"
// @Param to query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=current time"
// @Success 200 {object} domain.EntrancesAttendance
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/attendance/stores/entrances [get]
func (s *Server) apiChainEntrancesDataAttendance(c echo.Context) error {
	enterIDs, err := s.joinParam(c, "entrance_id", "entrance_ids")
	if err != nil {
		s.log.Warnf("joinParam error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	groupBy := c.QueryParam("group_by")
	if groupBy == "" {
		groupBy = defaultGroupBy
	}
	useRawData := c.QueryParam("use_rawdata")
	if useRawData == "" {
		useRawData = boolFalse
	}
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
	//
	filteredEnterIDsList, err := s.perm.FilteredEnters(c.Request(), layoutID, strings.Join(enterIDs, ","), acl.ActionRead)
	if err != nil {
		s.log.Errorf("perm.FilteredEnters error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errors.New("permission process failed")))
	}
	if filteredEnterIDsList == "" {
		s.log.Infof("for input %v and permissions %s, got empty filtered enter list",
			enterIDs, s.perm.FromRequest(c.Request()))
		return c.JSON(http.StatusOK, "[]")
	}
	data, err := repo.FindChainEntrancesDataAttendance(c.Request().Context(), from, to, groupBy, layoutID, useRawData, filteredEnterIDsList)
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindChainEntrancesDataAttendance error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}
