package infra

import (
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// MallDevicesResponse http wrapper with metadata
type MallDevicesResponse struct {
	Data domain.MallDevices `json:"data"`
	Metadata
}

// nolint:lll
// apiMallDevices docs
// @Summary Get all mall devices in the mall schema
// @Description get slice of devices with loc, date, layout_id, kind, is_active, dc_mode, mode, sn, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be device_id,layout_id,master_id,kind,title... all of them described at the model
// @Description include - comma separated list of entities, embedded in current, for devices it can be sensors,delay
// @Produce json
// @Tags malls/devices
// @Param layout_id query string false "default=*"
// @Param kind query string false "default=*"
// @Param is_active query string false "default=*"
// @Param dc_mode query string false "default=*"
// @Param mode query string false "default=*"
// @Param sn query string false "default=*"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "device_id,kind,title... default=all"
// @Param include query string false "sensors,delay default=none"
// @Success 200 {object} infra.MallDevicesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/devices [get]
func (s *Server) apiMallDevices(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	kind := c.QueryParam("kind")
	isActive := c.QueryParam("is_active")
	dcMode := c.QueryParam("dc_mode")
	mode := c.QueryParam("mode")
	sn := c.QueryParam("sn")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add list allowed read stores
	devices, count, err := repo.FindMallDevices(c.Request().Context(), loc, dt, layoutID, kind, isActive, dcMode, sn, mode,
		offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindMallDevices error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if devices == nil {
		response := MallDevicesResponse{
			Data: devices,
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
	devices.SetZeroValue(fields)
	includes := strings.Split(c.QueryParam("include"), ",")
	s.log.Debugf("include=%s", includes)
	for _, include := range includes {
		switch include {
		case "sensors":
			var lOffset, lLimit int64 = 0, 999999
			var shortSensorFieldSet []string = []string{"sensor_id", "device_id", "external_id", "kind"}
			sensors, _, err := repo.FindMallSensors(c.Request().Context(), loc, dt, "*", "*", "*", lOffset, lLimit)
			s.log.Debugf("got sensors len()=%d", len(sensors))
			if err == nil {
				sensors.SetZeroValue(shortSensorFieldSet)
				devices.IncludeSensors(sensors)
			}
		case "delay":
			delays, err := repo.FindMallDeviceDelays(c.Request().Context(), layoutID, devices.SliceIDs())
			if err == nil {
				devices.FillDelays(delays)
			}
		}
	}
	s.log.Debugf("got devices %+v", devices)
	response := MallDevicesResponse{
		Data: devices,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(devices)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiMallDeviceByID docs
// @Summary Get specified device in the mall schema
// @Description get device with loc, date, device_id, fields parameters from mall schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be device_id,layout_id...
// @Description all of them described at the model
// @Description include - comma separated list of entities, embedded in current, for devices it can be sensors
// @Produce  json
// @Tags malls/devices
// @Param layout_id query string false "default=*"
// @Param device_id path string true "uuid format"
// @Param loc query string false "location, default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "device_id,kind,title... default=all"
// @Param include query string false "sensors default=none"
// @Success 200 {object} domain.ChainDevice
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/malls/devices/{device_id} [get]
func (s *Server) apiMallDeviceByID(c echo.Context) error {
	id := c.Param("device_id")
	if id == "" {
		s.log.Errorf("bad request apiMallDeviceByID, %v", errEmptyID)
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
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add check by device
	device, err := repo.FindMallDeviceByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindMallDeviceByID error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if device == nil {
		s.log.Warnf("for id=%s, loc=%s and date=%s doesn't have device", id, loc, dt)
		return c.JSON(http.StatusOK, device)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	device.SetZeroValue(fields)
	includes := strings.Split(c.QueryParam("include"), ",")
	for _, include := range includes {
		switch include {
		case "sensors":
			var lOffset, lLimit int64 = 0, 999999
			var shortSensorFieldSet []string = []string{"sensor_id", "external_id", "kind"}
			sensors, _, err := repo.FindMallSensors(c.Request().Context(), loc, dt, "*", device.DeviceID, "*", lOffset, lLimit)
			if err == nil {
				sensors.SetZeroValue(shortSensorFieldSet)
				device.Sensors = sensors
			}
		case "delay":
			delays, err := repo.FindMallDeviceDelays(c.Request().Context(), layoutID, []string{device.DeviceID})
			if err == nil {
				device.FillDelay(delays)
			}
		}
	}
	return c.JSON(http.StatusOK, device)
}
