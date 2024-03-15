package infra

import (
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// RentersResponse http wrapper with metadata
type RentersResponse struct {
	Data domain.Renters `json:"data"`
	Metadata
}

// nolint:lll
// apiRenters docs
// @Summary Get all renters in the mall schema
// @Description get slice of renters with loc, date, layout_id, category_id, price_segment_id, contract, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field renter_id,title,layout_id,price_segment_id...
// @Produce  json
// @Tags malls/renters
// @Param layout_id query string true "uuid format, default=*"
// @Param categor_id query string false "default=*"
// @Param price_segment_id query string false "default=*"
// @Param contract query string false "default=*"
// @Param loc query string false "location, default=ru"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "renter_id,title,layout_id... default=all"
// @Success 200 {object} infra.RentersResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/renters [get]
func (s *Server) apiRenters(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	category := c.QueryParam("categor_id")
	priceSegment := c.QueryParam("price_segment_id")
	contract := c.QueryParam("contract")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm FilteredRenters read
	renters, count, err := repo.FindRenters(c.Request().Context(), loc, dt, layoutID, category, priceSegment, contract, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindRenters error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if renters == nil {
		response := RentersResponse{
			Data: renters,
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
	renters.SetZeroValue(fields)

	response := RentersResponse{
		Data: renters,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(renters)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiRenterByID docs
// @Summary Get specified renter in the mall schema
// @Description get renter with loc, date, renter_id, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field renter_id,title,layout_id,price_segment_id...
// @Description include - comma separated list of entities, embedded in current, for renters it can be only zones
// @Produce  json
// @Tags malls/renters
// @Param layout_id query string true "uuid format, default=*"
// @Param renter_id path string true "uuid/digits format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "renter_id,title,layout_id... default=all"
// @Param include query string false "zones default=none"
// @Success 200 {object} domain.Renter
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/renters/{renter_id} [get]
func (s *Server) apiRenterByID(c echo.Context) error {
	id := c.Param("renter_id")
	if id == "" {
		s.log.Errorf("bad request apiRenterByID, %v", errEmptyID)
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
	renter, err := repo.FindRenterByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindRenterByID error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if renter == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have records", id, loc, dt)
		return c.JSON(http.StatusOK, renter)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	renter.SetZeroValue(fields)
	include := c.QueryParam("include")
	if include == "zones" {
		zones, err := repo.FindMallZonesByRenter(c.Request().Context(), loc, dt, id)
		if err == nil {
			var shortZoneFieldSet []string = []string{"zone_id", "kind", "title", "options"}
			zones.SetZeroValue(shortZoneFieldSet)
			renter.Zones = zones
		}
	}
	return c.JSON(http.StatusOK, renter)
}
