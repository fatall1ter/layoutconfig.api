package infra

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// nolint:lll
// apiZoneDataInsideDay docs
// @Summary Get data inside at daily intervals by cumulative total
// @Description get data on the number of people inside the zone at intervals by cumulative total, can be filtered by zone_id and day parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags data/inside
// @Param layout_id path string true "uuid format, default=*"
// @Param zone_id query string false "default=*"
// @Param day query string false "YYYY-MM-DD default=today"
// @Success 200 {object} domain.DatasInside
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/inside/days [get]
func (s *Server) apiZoneDataInsideDay(c echo.Context) error {
	pZoneID := c.QueryParam("zone_id") // TODO: add zoneIDs param
	var zoneID *string
	if pZoneID != "" {
		zoneID = &pZoneID
	}
	sDay := c.QueryParam("day")
	day := time.Now()
	if sDay != "" {
		pday, err := time.Parse("2006-01-02", sDay)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				ErrInvalidRequest(fmt.Errorf("wrong day [%v] parameter %v", sDay, err)))
		}
		day = pday
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add ZoneFiltered read
	data, err := repo.FindZoneDataInsideDay(c.Request().Context(), zoneID, &day)
	if err != nil {
		s.log.Errorf("repo.FindZoneDataInsideDay error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// nolint:lll
// apiZoneDataInsideRange docs
// @Summary Get data inside by range interval of the day
// @Description get data on the number of people inside the zone by range intervals,  can be filtered by the zone_id parameter and the day parameter.
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description from, to and day must be at same day
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss then the server's local timezone is applied
// @Produce  json
// @Tags data/inside
// @Param layout_id path string true "uuid format, default=*"
// @Param zone_id query string false "default=*"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Param day query string false "YYYY-MM-DD default=today"
// @Success 200 {object} domain.DatasInside
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/inside/days/range [get]
func (s *Server) apiZoneDataInsideRange(c echo.Context) error {
	from, to, err := s.getFromToParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	threshold := time.Hour * 24
	if to.Sub(from) > threshold {
		return c.JSON(http.StatusBadRequest,
			ErrInvalidRequest(fmt.Errorf("from [%s]-to [%s] diff more than [%v]", from, to, threshold)))
	}
	pZoneID := c.QueryParam("zone_id") // TODO: add zoneIDs param
	var zoneID *string = nil
	if pZoneID != "" {
		zoneID = &pZoneID
	}
	sDay := c.QueryParam("day")
	day := time.Now()
	if sDay != "" {
		pday, er := time.Parse("2006-01-02", sDay)
		if er != nil {
			s.log.Warnf("parse day %s error, %s", sDay, err)
			return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("wrong day [%v] parameter %v", sDay, err)))
		}
		day = pday
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add ZoneFiltered read
	data, err := repo.FindZoneDataInsideRange(c.Request().Context(), from, to, zoneID, &day)
	if err != nil {
		s.log.Errorf("repo.FindZoneDataInsideRange error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}

// apiZoneDataInside docs
// @Summary Get data inside at now
// @Description get data on the number of people inside the zone at now, can be filtered by zone_id parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags data/inside
// @Param layout_id path string true "uuid format, default=*"
// @Param zone_id query string false "default=*"
// @Success 200 {object} domain.DatasInside
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/data/inside [get]
func (s *Server) apiZoneDataInside(c echo.Context) error {
	pZoneID := c.QueryParam("zone_id") // TODO: add zoneIDs list param
	var zoneID *string
	if pZoneID != "" {
		zoneID = &pZoneID
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add ZoneFiltered read
	data, err := repo.FindZoneDataInsideNow(c.Request().Context(), zoneID) // TODO: add zoneIDs list param
	if err != nil {
		s.log.Errorf("repo.FindZoneDataInsideNow error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, data)
}
