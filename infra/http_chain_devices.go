package infra

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

var (
	errInternal      error = errors.New("internal error")
	errEmptyLayoutID error = errors.New("empty layout_id not allowed")
)

// NewChainDevice properties of new device,
// controller, navigator, etc...
type NewChainDevice struct {
	LayoutID string  `json:"layout_id,omitempty"`
	StoreID  string  `json:"store_id,omitempty"`
	MasterID *string `json:"master_id,omitempty"`
	Kind     string  `json:"kind,omitempty"`
	Title    string  `json:"title,omitempty"`
	IsActive bool    `json:"is_active,omitempty"`
	IP       string  `json:"ip,omitempty"`
	Port     string  `json:"port,omitempty"`
	SN       string  `json:"sn,omitempty"`
	// Mode of people-counting: single/master/slave
	Mode string `json:"mode,omitempty"`
	// DCMode data collector mode: active - device transmit data to server; passive - server request data from device
	DCMode    string     `json:"dcmode,omitempty"`
	Login     string     `json:"login,omitempty"`
	Password  string     `json:"password,omitempty"`
	Options   string     `json:"options,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
	Creator   string     `json:"creator,omitempty"`
}

func (cd *NewChainDevice) makeDomainChainDevice() domain.ChainDevice {
	return domain.ChainDevice{
		LayoutID:  cd.LayoutID,
		StoreID:   cd.StoreID,
		MasterID:  cd.MasterID,
		Kind:      cd.Kind,
		Title:     cd.Title,
		IsActive:  cd.IsActive,
		IP:        cd.IP,
		Port:      cd.Port,
		SN:        cd.SN,
		Mode:      cd.Mode,
		DCMode:    cd.DCMode,
		Login:     cd.Login,
		Password:  cd.Password,
		Options:   cd.Options,
		Notes:     cd.Notes,
		ValidFrom: cd.ValidFrom,
		ValidTo:   cd.ValidTo,
		Creator:   cd.Creator,
	}
}

// nolint:lll
// apiCreateRetailDevice docs
// @Summary Create new device in the retail schema
// @Description create new device in the retail schema
// @Description example new device, layout_id and store_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"dcmode": "active",
// @Description 	"ip": "98.99.100.101",
// @Description 	"is_active": true,
// @Description 	"kind": "device.3dv",
// @Description 	"layout_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description 	"login": "admin",
// @Description 	"mode": "single",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"options": "{\\"localIP\\":\\"192.168.0.10\\", \\"localPort\\": 80}",
// @Description 	"password": "passW0rd",
// @Description 	"port": "8080",
// @Description 	"sn": "00:00:00:00:11:22",
// @Description 	"store_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description 	"title": "{\\"ru\\":\\"Наименование Устройства на русском языке\\",\\"en\\":\\"Device name in English\\"}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/devices
// @Param device body infra.NewChainDevice true "device properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/devices [post]
func (s *Server) apiCreateChainDevice(c echo.Context) error {
	device := &NewChainDevice{}
	if err := c.Bind(device); err != nil {
		s.log.Errorf("apiCreateChainDevice, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if device.LayoutID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyLayoutID))
	}
	repo, ok := s.repoM.RepoByID(device.LayoutID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errLayoutRepoNotFound))
	}
	id, err := repo.AddChainDevice(c.Request().Context(), device.makeDomainChainDevice())
	if err != nil {
		s.log.Errorf("repo.AddChainDevice error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	href := `{"href":"/v2/chains/devices/` + id + `"}`
	return c.JSON(http.StatusCreated, CreatedStatus(href))
}

// ChainDevicesResponse http wrapper with metadata
type ChainDevicesResponse struct {
	Data domain.ChainDevices `json:"data"`
	Metadata
}

// nolint:lll
// apiChainDevices docs
// @Summary Get all chains devices in the retail schema
// @Description get slice of devices with loc, date, layout_id, store_id, kind, is_active, dc_mode, mode, sn, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be device_id,layout_id,store_id,master_id,kind,title... all of them described at the model
// @Description include - comma separated list of entities, embedded in current, for devices it can be sensors,delay
// @Produce json
// @Tags chains/devices
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
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
// @Success 200 {object} infra.ChainDevicesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/devices [get]
func (s *Server) apiChainDevices(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	storeID := c.QueryParam("store_id")
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
	devices, count, err := repo.FindChainDevices(c.Request().Context(), loc, dt, layoutID, storeID, kind, isActive, dcMode, sn, mode,
		offset, limit) // TODO: add list storeIDs param
	if err != nil {
		s.log.Errorf("repo.FindChainDevices error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if devices == nil {
		response := ChainDevicesResponse{
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
			sensors, _, err := repo.FindChainSensors(c.Request().Context(), loc, dt, "*", "*", "*", "*", lOffset, lLimit) // TODO: add store list ids param
			s.log.Debugf("got sensors len()=%d", len(sensors))
			if err == nil {
				sensors.SetZeroValue(shortSensorFieldSet)
				devices.IncludeSensors(sensors)
			}
		case "delay":
			delays, err := repo.FindChainDeviceDelays(c.Request().Context(), layoutID, nil, devices.SliceIDs()) // TODO: add store list ids param
			if err == nil {
				devices.FillDelays(delays)
			}
		}
	}
	s.log.Debugf("got devices %+v", devices)
	response := ChainDevicesResponse{
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

// apiChainDeviceByID docs
// @Summary Get specified device in the retail schema
// @Description get device with loc, date, device_id, fields parameters from retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be device_id,layout_id...
// @Description all of them described at the model
// @Description include - comma separated list of entities, embedded in current, for devices it can be sensors
// @Produce  json
// @Tags chains/devices
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
// @Router /v2/chains/devices/{device_id} [get]
func (s *Server) apiChainDeviceByID(c echo.Context) error { // TODO: add store_id params?
	id := c.Param("device_id")
	if id == "" {
		s.log.Errorf("bad request apiChainDeviceByID, %v", errEmptyID)
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
	device, err := repo.FindChainDeviceByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindChainDeviceByID error, %s", err)
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
			sensors, _, err := repo.FindChainSensors(c.Request().Context(), loc, dt, "*", "*", device.DeviceID, "*", lOffset, lLimit)
			if err == nil {
				sensors.SetZeroValue(shortSensorFieldSet)
				device.Sensors = sensors
			}
		case "delay":
			delays, err := repo.FindChainDeviceDelays(c.Request().Context(), layoutID, nil, []string{device.DeviceID})
			if err == nil {
				device.FillDelay(delays)
			}
		}
	}
	return c.JSON(http.StatusOK, device)
}

// apiDeleteChainDevice docs
// @Summary Delete specified device in the retail schema
// @Description delete specified device in the retail schema by device_id parameter
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/devices
// @Param layout_id query string false "default=*"
// @Param device_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/devices/{device_id} [delete]
func (s *Server) apiDeleteChainDevice(c echo.Context) error {
	id := c.Param("device_id")
	if id == "" {
		s.log.Errorf("bad request apiDeleteChainDevice, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainDevice(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.DelChainDevice error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d device(s)", cnt)))
}

// UpdChainDevice properties for upd device,
// controller, navigator, etc...
type UpdChainDevice struct {
	DeviceID string  `json:"device_id,omitempty"`
	MasterID *string `json:"master_id,omitempty"`
	Kind     string  `json:"kind,omitempty"`
	Title    string  `json:"title,omitempty"`
	IsActive bool    `json:"is_active,omitempty"`
	IP       string  `json:"ip,omitempty"`
	Port     string  `json:"port,omitempty"`
	SN       string  `json:"sn,omitempty"`
	// Mode of people-counting: single/master/slave
	Mode string `json:"mode,omitempty"`
	// DCMode data collector mode: active - device transmit data to server; passive - server request data from device
	DCMode    string     `json:"dcmode,omitempty"`
	Login     string     `json:"login,omitempty"`
	Password  string     `json:"password,omitempty"`
	Options   string     `json:"options,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	ValidFrom *time.Time `json:"valid_from,omitempty"`
	NoHistory bool       `json:"no_history"`
}

func (cd *UpdChainDevice) makeDomainChainDevice() domain.ChainDevice {
	return domain.ChainDevice{
		DeviceID:  cd.DeviceID,
		MasterID:  cd.MasterID,
		Kind:      cd.Kind,
		Title:     cd.Title,
		IsActive:  cd.IsActive,
		IP:        cd.IP,
		Port:      cd.Port,
		SN:        cd.SN,
		Mode:      cd.Mode,
		DCMode:    cd.DCMode,
		Login:     cd.Login,
		Password:  cd.Password,
		Options:   cd.Options,
		Notes:     cd.Notes,
		ValidFrom: cd.ValidFrom,
	}
}

// nolint:lll
// apiUpdRetailDevice docs
// @Summary Update device in the retail schema
// @Description update exists device in the retail schema
// @Description valid_from set if need save history for specified date
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example upd device, device_id must to be pass in parameters:
// @Description {
// @Description 	"device_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description 	"dcmode": "active",
// @Description 	"ip": "98.99.100.101",
// @Description 	"is_active": true,
// @Description 	"kind": "device.3dv",
// @Description 	"login": "admin",
// @Description 	"mode": "single",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"options": "{\\"localIP\\":\\"192.168.0.10\\", \\"localPort\\": 80}",
// @Description 	"password": "passW0rd",
// @Description 	"port": "8080",
// @Description 	"sn": "00:00:00:00:11:22",
// @Description 	"title": "{\\"ru\\":\\"Наименование Устройства на русском языке\\",\\"en\\":\\"Device name in English\\"}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/devices
// @Param layout_id query string false "default=*"
// @Param device body infra.UpdChainDevice true "device properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/devices [put]
func (s *Server) apiUpdChainDevice(c echo.Context) error {
	device := &UpdChainDevice{}
	if err := c.Bind(device); err != nil {
		s.log.Errorf("apiUpdRetailDevice, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if device.DeviceID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty device_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	count, err := repo.UpdChainDevice(c.Request().Context(), device.makeDomainChainDevice(), device.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdChainDevice error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("updated %d device(s)", count)))
}

// ChainDeviceTracksResponse http wrapper with metadata
type ChainDeviceTracksResponse struct {
	Data domain.Tracks `json:"data"`
	Metadata
}

// apiChainDeviceTracks docs
// @Summary Get device customer tracks
// @Description get device customer tracks in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/devices
// @Param layout_id query string false "digit/uuid format"
// @Param store_id query string false "digit/uuid format"
// @Param device_id query string false "digit/uuid format"
// @Param from query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=start today"
// @Param to query string false "ISO8601 YYYY-MM-DD HH:mm:SS timestamp default=current time"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Success 200 {object} infra.ChainDeviceTracksResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/devices/tracks [get]
func (s *Server) apiChainDeviceTracks(c echo.Context) error { // TODO: remove me
	store := c.QueryParam("store_id")
	device := c.QueryParam("device_id")
	offset, limit := s.getPageParams(c)
	from, to, err := s.getFromToParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	tracks, count, err := repo.FindChainDeviceTracks(c.Request().Context(), layoutID, store, device, from, to, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindChainDeviceTracks failed: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	response := ChainDeviceTracksResponse{
		Data: tracks,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(tracks)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiChainDeviceTracksAt docs
// @Summary Get device customer tracks at time moment
// @Description get device customer tracks at time moment with specified accuracy in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/devices
// @Param layout_id query string false "digit/uuid format"
// @Param store_id query string false "digit/uuid format"
// @Param device_id query string false "digit/uuid format"
// @Param accuracy query string false "1s,5s,15s golang time.Duration format, default=1s"
// @Param at query string false "ISO8601 timestamp default=current moment"
// @Success 200 {object} domain.Tracks
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/devices/tracks/attime [get]
func (s *Server) apiChainDeviceTracksAt(c echo.Context) error { // TODO: remove me
	store := c.QueryParam("store_id")
	device := c.QueryParam("device_id")
	sAcc := c.QueryParam("accuracy")
	acc := time.Second
	if sAcc != "" {
		tAcc, err := time.ParseDuration(sAcc)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("wrong accuracy [%v] parameter %v", sAcc, err)))
		}
		acc = tAcc
	}
	sAt := c.QueryParam("at")
	_, toffset := time.Now().Zone()
	if sAt == "" {
		sAt = time.Now().Format("2006-01-02T00:00:00Z07:00")
	}
	at, err := time.Parse(time.RFC3339, sAt)
	if err != nil {
		fromNaive, er := time.Parse(naiveTimeFormat, sAt)
		if er != nil {
			return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("wrong at [%v] parameter %v", sAt, er)))
		}
		at = fromNaive.Local().Add(-time.Duration(toffset) * time.Second)
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	tracks, err := repo.FindChainDeviceTracksAt(c.Request().Context(), layoutID, store, device, at, acc)
	if err != nil {
		s.log.Errorf("repo.FindChainDeviceTracksAt failed, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, tracks)
}
