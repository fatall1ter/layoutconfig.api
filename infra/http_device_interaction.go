package infra

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.countmax.ru/countmax/devcount"
	"git.countmax.ru/countmax/devcount/vvtk"
	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

const (
	deviceTimeout time.Duration = 60 * time.Second
)

// apiDeviceInfo docs
// @Summary Get specified device info
// @Description get device info: serialNumber, sensors with parameter device_id=id
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags common
// @Param layout_id path string true "uuid format, default=*"
// @Param device_id path string true "uuid format"
// @Success 200 {object} domain.DevConfig
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/devices/{device_id}/info [get]
func (s *Server) apiDeviceInfo(c echo.Context) error {
	id := c.Param("device_id")
	if id == "" {
		s.log.Errorf("bad request apiRetailDeviceByID, %v", errEmptyID)
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
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	device, err := repo.FindChainDeviceByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("FindChainDeviceByID(%s, %s, %s) error %v", loc, id, dt, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	if device == nil {
		s.log.Warnf("for id=%s, loc=%s and date=%s doesn't have device", id, "", "")
		return c.JSON(http.StatusNotFound, ErrNotFound(fmt.Errorf("for id=%s doesn't have device", id)))
	}

	dc, err := newDevCount(device.Kind, device.IP, device.Port, device.Login, device.Password, deviceTimeout)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	cfg, err := dc.GetConfig(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	response := domain.DevConfig{
		SN:      cfg.SN,
		Sensors: make(domain.DevSensors, len(cfg.Sensors)),
	}
	for i, ds := range cfg.Sensors {
		response.Sensors[i].ID = ds.ID
		response.Sensors[i].Kind = string(ds.Kind)
	}
	return c.JSON(http.StatusOK, response)
}

// apiNewDeviceInfo docs
// @Summary Get some device info
// @Description get device info about serialNumber, sensors
// @Description parameters:
// @Description kind: device.3dh, device.3dv, device.3d, device.bdv, device.3dx...
// @Description timeout in seconds. Attention 3dv little big latency, recommended 60 sec
// @Produce  json
// @Tags common
// @Param kind query string false "device kind, default=device.3dv"
// @Param ip query string false "ip/fqdn access to device, default=192.168.0.1"
// @Param port query string false "tcp port access to device, default=80"
// @Param login query string false "login, default=admin"
// @Param password query string false "password, default=passw0rd"
// @Param timeout query string false "timeout in seconds, default=30"
// @Success 200 {object} domain.DevConfig
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/devices/info [get]
func (s *Server) apiNewDeviceInfo(c echo.Context) error {
	kind := c.QueryParam("kind")
	if kind == "" {
		kind = "device.3dv"
	}
	ip := c.QueryParam("ip")
	if ip == "" {
		ip = "192.168.0.1"
	}
	port := c.QueryParam("port")
	if port == "" {
		port = "80"
	}
	login := c.QueryParam("login")
	if login == "" {
		login = "admin"
	}
	pass := c.QueryParam("password")
	if pass == "" {
		pass = "passw0rd"
	}
	timeoutStr := c.QueryParam("timeout")
	if timeoutStr == "" {
		timeoutStr = "30"
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 30
	}

	dc, err := newDevCount(kind, ip, port, login, pass, time.Duration(timeout)*time.Second)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	cfg, err := dc.GetConfig(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	response := domain.DevConfig{
		SN:      cfg.SN,
		Sensors: make(domain.DevSensors, len(cfg.Sensors)),
	}
	for i, ds := range cfg.Sensors {
		response.Sensors[i].ID = ds.ID
		response.Sensors[i].Kind = string(ds.Kind)
	}
	return c.JSON(http.StatusOK, response)
}

func newDevCount(kind, ip, port, login, pass string, timeout time.Duration) (devcount.DevCount, error) {
	switch kind {
	case "device.3dv":
		return vvtk.NewSC81XX(ip, port, login, pass, timeout), nil
	//case "device.3dx":
	//	return xovis.NewXovis(ip, port, login, pass, timeout), nil
	//case "device.cm_15", "device.cm_18":
	default:
		return nil, fmt.Errorf("device kind=%s not supported, allow only %v", kind, []string{"device.3d", "device.3dx", "device.cm_15", "device.cm_18"})
	}
}
