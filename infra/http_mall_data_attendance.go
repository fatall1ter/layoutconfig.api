package infra

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.countmax.ru/countmax/layoutconfig.api/repos"
	"github.com/labstack/echo/v4"
)

// apiMallZonesDataAttendance docs
// @Summary Get data attendance for zones of the mall
// @Description get data enters and exits in/out the zones of the mall by group_by intervals,
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description can be filtered by list zone_id as comma separated list (12345,12344567,8488...)
// @Description and daterange parameters from and to.
// @Description intervals can be: interval, hour, day, week, month, quarter, year
// @Description use_rawdata for request data by sensors, can be: true/false/1/0; default false
// @Produce  json
// @Tags data/attendance
// @Param zone_ids query string false "default=*"
// @Param layout_id query string false "default=*"
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
// @Router /v2/data/attendance/malls/zones [get]
func (s *Server) apiMallZonesDataAttendance(c echo.Context) error {
	zoneIDs, err := s.joinParam(c, "zone_id", "zone_ids")
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
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("wrong parameter %v", err)))
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add FilteredZone read
	data, err := repo.FindMallZonesDataAttendance(c.Request().Context(), from, to, groupBy, layoutID, useRawData, strings.Join(zoneIDs, ","))
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindMallZonesDataAttendance error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiMallEntrancesDataAttendance docs
// @Summary Get data attendance for entrances of the mall
// @Description get data enters and exits in/out the entrances of the mall by group_by intervals,
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description can be filtered by list entrance_id as comma separated list (12345,12344567,8488...)
// @Description and daterange parameters from and to.
// @Description intervals can be: interval, hour, day, week, month, quarter, year
// @Description use_rawdata for request data by sensors, can be: true/false/1/0; default false
// @Produce  json
// @Tags data/attendance
// @Param entrance_ids query string false "default=*"
// @Param layout_id query string false "default=*"
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
// @Router /v2/data/attendance/malls/entrances [get]
func (s *Server) apiMallEntrancesDataAttendance(c echo.Context) error {
	enterIDs, err := s.joinParam(c, "entrance_id", "entrance_ids")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("%s", err)))
	}
	groupBy := c.QueryParam("group_by")
	if groupBy == "" {
		groupBy = defaultGroupBy
	}
	useRawData := c.QueryParam("use_rawdata")
	if useRawData == "" {
		useRawData = "false"
	}
	from, to, err := s.getFromToParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("wrong parameter %v", err)))
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add FilteredEnter read
	data, err := repo.FindMallEntrancesDataAttendance(c.Request().Context(), from, to, groupBy, layoutID, useRawData, strings.Join(enterIDs, ","))
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindMallEntrancesDataAttendance error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	return c.JSON(http.StatusOK, data)
}

// apiRenterDataAttendance docs
// @Summary Get data attendance for renters of the mall
// @Description get data enters and exits in/out the entrances of the mall's renters by group_by intervals,
// @Description can be filtered by list renter_ids as comma separated list (12345,12344567,8488...) and layout_id
// @Description and daterange parameters from and to.
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description intervals can be: interval, hour, day, week, month, quarter, year
// @Description use_rawdata for request data by sensors, can be: true/false/1/0; default false
// @Produce  json
// @Tags data/attendance
// @Param renter_ids query string false "default=*"
// @Param layout_id query string false "default=*"
// @Param group_by query string false "default=interval"
// @Param use_rawdata query string false "default=false"
// @Param from query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=start today"
// @Param to query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=current time"
// @Success 200 {object} domain.RentersAttendance
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 413 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/attendance/malls/renters [get]
func (s *Server) apiRenterDataAttendance(c echo.Context) error {
	renterIDs, err := s.joinParam(c, "renter_id", "renter_ids")
	if err != nil {
		s.log.Errorf("joinParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errors.New("wrong renter_id/renter_ids param")))
	}
	groupBy := c.QueryParam("group_by")
	if groupBy == "" {
		groupBy = defaultGroupBy
	}
	useRawData := c.QueryParam("use_rawdata")
	if useRawData == "" {
		useRawData = "false"
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
	// TODO: add perm ListFilteredRenters read
	data, err := repo.FindRenterDataAttendance(c.Request().Context(), from, to, groupBy, layoutID, useRawData, strings.Join(renterIDs, ","))
	if err != nil {
		if err == repos.ErrNotAllowedDataRange {
			s.log.Errorf("repo.FindRenterDataAttendance too large error, %s", err)
			return c.JSON(http.StatusRequestEntityTooLarge, ErrPayloadTooLarge(err))
		}
		s.log.Errorf("repo.FindRenterDataAttendance error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}
