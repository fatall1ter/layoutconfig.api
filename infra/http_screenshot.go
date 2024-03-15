package infra

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ScreenshotsResponse http wrapper with metadata
type ScreenshotsResponse struct {
	Data domain.Screenshots `json:"data"`
	Metadata
}

// nolint:lll
// apiGetScreenshots docs
// @Summary Get screenshots from the storage
// @Description get slice of screenshots with layout_id, store_id, zone_id, device_id, status, offset, limit parameters
// @Description fields - comma separated values of field names, can be layout_id,store_id,device_id,url,notes... all of them described in the model
// @Description screenshot_status can be new - default processed - blocked for changes, to_delete - marked for delete, archived - marked for all time store
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss then the server's local timezone is applied
// @Produce  json
// @Tags screenshots
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param zone_id query string false "default=*"
// @Param device_id query string false "default=*"
// @Param screenshot_status query string false "default=*"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Param fields query string false "layout_id,store_id,device_id..default=all"
// @Success 200 {object} infra.ScreenshotsResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/screenshots [get]
func (s *Server) apiGetScreenshots(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	layoutID := c.QueryParam("layout_id")
	storeID := c.QueryParam("store_id")
	zoneID := c.QueryParam("zone_id")
	deviceID := c.QueryParam("device_id")
	status := c.QueryParam("screenshot_status")
	from, to, err := s.getFromToParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	// define devices IDs by zone
	deviceIDs := []string{deviceID}
	if zoneID != "" && (deviceID == "" || deviceID == "*") {
		defDeviceIDs, er := s.getDevicesByZoneID(c.Request().Context(), layoutID, storeID, zoneID)
		if er != nil {
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
		}
		deviceIDs = defDeviceIDs
	}
	// request screens with devices ids
	screens, count, err := s.dmRepo.FindScreens(layoutID, storeID, status, deviceIDs, from, to, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	if len(screens) == 0 {
		s.log.With(zap.String("layoutID", layoutID), zap.String("storeID", storeID), zap.String("deviceID", deviceID)).
			Warnf("from=%s, to=%s, offset=%d, limit=%d, nothing not found", from, to, offset, limit)
		return c.JSON(http.StatusNotFound,
			ErrNotFound(fmt.Errorf("from=%s, to=%s for layout_id=%s, store_id=%s and device_id=%s does't have screens",
				from, to, layoutID, storeID, deviceID)))
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = []string{"device_id", "layout_id", "store_id", "screenshot_time", "url"}
	}
	screens.SetZeroValue(fields)
	response := ScreenshotsResponse{
		Data: screens,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(screens)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiGetScreenshotsAtTime docs
// @Summary Get all screenshots from the storage on the specified time
// @Description get slice of screenshots with layout_id, store_id, zone_id, device_id, at time momen parameters
// @Description fields - comma separated values of field names, can be layout_id,store_id,device_id,url,notes... all of them described in the model
// @Description at can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss then the server's local timezone is applied
// @Produce  json
// @Tags screenshots
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param zone_id query string false "default=*"
// @Param device_id query string false "default=*"
// @Param at query string false "ISO8601 datetime, default=currentMoment-1min accuracy 15s"
// @Param fields query string false "layout_id,store_id,device_id..default=all"
// @Success 200 {object} domain.Screenshots
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/screenshots/attime [get]
func (s *Server) apiGetScreenshotsAtTime(c echo.Context) error {
	layoutID := c.QueryParam("layout_id")
	storeID := c.QueryParam("store_id")
	zoneID := c.QueryParam("zone_id")
	deviceID := c.QueryParam("device_id")
	sAt := c.QueryParam("at")
	if sAt == "" {
		sAt = time.Now().Truncate(time.Second * 15).Add(-time.Minute).Format(time.RFC3339)
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
	// define devices IDs by zone
	deviceIDs := []string{deviceID}
	if zoneID != "" && (deviceID == "" || deviceID == "*") {
		defDeviceIDs, err := s.getDevicesByZoneID(c.Request().Context(), layoutID, storeID, zoneID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
		}
		deviceIDs = defDeviceIDs
	}

	// request screens with devices ids
	screens, err := s.dmRepo.FindScreensAtTime(layoutID, storeID, deviceIDs, at)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	if len(screens) == 0 {
		s.log.With(zap.String("layoutID", layoutID), zap.String("storeID", storeID), zap.Strings("deviceIDs", deviceIDs)).Warnf("at=%s, nothing not found", at)
		return c.JSON(http.StatusNotFound, ErrNotFound(fmt.Errorf("at=%s for layout_id=%s, store_id=%s and device_ids=%v does't have screens", at, layoutID, storeID, deviceIDs)))
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = []string{"device_id", "layout_id", "store_id", "screenshot_time", "url"}
	}
	screens.SetZeroValue(fields)
	return c.JSON(http.StatusOK, screens)
}

func (s *Server) getDevicesByZoneID(ctx context.Context, layoutID, storeID, zoneID string) ([]string, error) {
	dt := time.Now().Format("2006-01-02")
	var lOffset, lLimit int64 = 0, 999999
	repo, ok := s.repoM.RepoByID(layoutID)
	if !ok {
		s.log.Errorf("not found repo for layout_id=%s", layoutID)
		return nil, errLayoutRepoNotFound
	}
	devices, _, err := repo.FindChainDevices(ctx, "ru", dt, layoutID, storeID, "", "", "", "", "", lOffset, lLimit)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("devices not found")
	}
	var shortSensorFieldSet []string = []string{"sensor_id", "device_id", "external_id", "kind"}
	sensors, _, err := repo.FindChainSensors(ctx, "ru", dt, "*", "*", "*", "*", lOffset, lLimit)
	if err == nil {
		sensors.SetZeroValue(shortSensorFieldSet)
		devices.IncludeSensors(sensors)
	}
	binds, _, err := repo.FindBindChainSensorZone(ctx, zoneID, "*", lOffset, lLimit)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, 4)
	for _, bind := range binds {
		d := devices.FindDeviceBySensorID(bind.SensorID)
		if d != nil {
			res = append(res, d.DeviceID)
		}
	}
	return res, nil
}

// apiUpdScreenshotsStatus docs
// @Summary Updates screenshots status
// @Description Updates screenshots status
// @Description example upd screenshots statuses:
// @Description [{
// @Description 	"device_id": "1234567",
// @Description 	"screenshot_time": "2020-09-01T00:15:15+03:00",
// @Description 	"screenshot_status": "to_delete"
// @Description },
// @Description {
// @Description 	"device_id": "1234568",
// @Description 	"screenshot_time": "2020-09-01T00:15:15+03:00",
// @Description 	"screenshot_status": "to_delete"
// @Description }]
// @Accept  json
// @Produce  json
// @Tags screenshots
// @Param params body domain.ParamsScreenUpd true "screenshots update status properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/screenshots [put]
func (s *Server) apiUpdScreenshotsStatus(c echo.Context) error {
	params := domain.ParamsScreenUpd{}
	if err := c.Bind(&params); err != nil {
		s.log.Errorf("apiUpdScreenshotsStatus, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	cu, err := s.dmRepo.UpdStatusManyScreens(params)
	if err != nil {
		if strings.Contains(err.Error(), "incorrect status value, allowed only") {
			return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
		}
		s.log.Errorf("UpdStatusManyScreens failed %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	return c.JSON(http.StatusOK, OkStatus(`updated `+strconv.FormatInt(cu, 10)))
}
