package infra

import (
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// MallZonesResponse http wrapper with metadata
type MallZonesResponse struct {
	Data domain.MallZones `json:"data"`
	Metadata
}

// apiChainZones docs
// @Summary Get all zones in the mall schema
// @Description get slice of zones with loc, date, layout_id, is_online, is_active, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field zone_id,parent_id,layout_id,kind,title...
// @Produce  json
// @Tags malls/zones
// @Param layout_id query string false "default=*"
// @Param kind query string false "default=*"
// @Param parent_id query string false "default=*"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param is_online query string false "default=*"
// @Param is_active query string false "default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "zone_id,kind,title... default=all"
// @Success 200 {object} infra.MallZonesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/zones [get]
func (s *Server) apiMallZones(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	kind := c.QueryParam("kind")
	parentID := c.QueryParam("parent_id")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	isOnline := c.QueryParam("is_online")
	isActive := c.QueryParam("is_active")
	//
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: get perm FilteredMallZones
	zones, count, err := repo.FindMallZones(c.Request().Context(), loc, dt, layoutID, kind, parentID, isOnline, isActive, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindMallZones error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if zones == nil {
		s.log.With(zap.String("layoutID", layoutID)).
			Infof("kind=%s, parentID=%s, loc=%s, dt=%s, offset=%d, limit=%d", kind, parentID, loc, dt, offset, limit)
		response := MallZonesResponse{
			Data: zones,
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
	zones.SetZeroValue(fields)
	response := MallZonesResponse{
		Data: zones,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(zones)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiMallZoneByID docs
// @Summary Get specified zone in the mall schema
// @Description get zone with loc, date, zone_id, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field zone_id,parent_id,layout_id,kind,title...
// @Description include - comma separated list of entities, embedded in current, for zone it can be only entrances
// @Produce  json
// @Tags malls/zones
// @Param layout_id query string false "default=*"
// @Param zone_id path string true "uuid/digits format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "zone_id,kind,title... default=all"
// @Param include query string false "entrances default=none"
// @Success 200 {object} domain.MallZone
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/zones/{zone_id} [get]
func (s *Server) apiMallZoneByID(c echo.Context) error {
	id := c.Param("zone_id")
	if id == "" {
		s.log.Errorf("bad request apiMallZoneByID, %v", errEmptyID)
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
	// TODO: add perm CheckMallZone read
	zone, err := repo.FindMallZoneByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindMallZoneByID error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if zone == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have zone", id, loc, dt)
		return c.JSON(http.StatusOK, zone)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	zone.SetZeroValue(fields)
	include := c.QueryParam("include")
	if include == "entrances" {
		enters, err := repo.FindCrossesZoneEnter(c.Request().Context(), id, "")
		if err == nil {
			// enters.SetZeroValue(shortSensorFieldSet)
			zone.Entrances = enters
		}
	}
	return c.JSON(http.StatusOK, zone)
}
