package infra

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// NewChainSensor is new sensor entity
type NewChainSensor struct {
	DeviceID   string     `json:"device_id,omitempty"`
	LayoutID   string     `json:"layout_id,omitempty"`
	StoreID    string     `json:"store_id,omitempty"`
	ExternalID string     `json:"external_id,omitempty"`
	Kind       string     `json:"kind,omitempty"`
	Title      string     `json:"title,omitempty"`
	Options    string     `json:"options,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
}

func (cs *NewChainSensor) makeDomainChainSensor() domain.ChainSensor {
	return domain.ChainSensor{
		DeviceID:   cs.DeviceID,
		LayoutID:   cs.LayoutID,
		StoreID:    cs.StoreID,
		ExternalID: cs.ExternalID,
		Kind:       cs.Kind,
		Title:      cs.Title,
		Options:    cs.Options,
		Notes:      cs.Notes,
		ValidFrom:  cs.ValidFrom,
		ValidTo:    cs.ValidTo,
		Creator:    cs.Creator,
	}
}

// nolint:lll
// apiCreateChainSensor docs
// @Summary Create new sensor in the retail schema
// @Description create new sensor in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new sensor, layout_id, device_id and store_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"device_id": "47cbe28b-5370-4721-986e-468a2a0c9b87",
// @Description 	"external_id": "Rule-1",
// @Description 	"kind": "sensor.people_count",
// @Description 	"layout_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"store_id": "29587d9a-05d3-4d2c-a974-f2c11fcb30fa",
// @Description 	"title": "{\\"ru\\":\\"Наименование Сенсора на русском языке\\",\\"en\\":\\"Sensor name in English\\"}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"options": "{}"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "default=*"
// @Param sensor body infra.NewChainSensor true "sensor properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors [post]
func (s *Server) apiCreateChainSensor(c echo.Context) error {
	sensor := &NewChainSensor{}
	if err := c.Bind(sensor); err != nil {
		s.log.Errorf("apiCreateChainSensor, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckLayout create
	id, err := repo.AddChainSensor(c.Request().Context(), sensor.makeDomainChainSensor())
	if err != nil {
		s.log.Errorf("repo.AddChainSensor error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	href := `{"href":"/v2/chains/sensors/` + id + `"}`
	return c.JSON(http.StatusCreated, CreatedStatus(href))
}

// ChainSensorsResponse http wrapper with metadata
type ChainSensorsResponse struct {
	Data domain.ChainSensors `json:"data"`
	Metadata
}

// nolint:lll
// apiChainSensors docs
// @Summary Get all sensors in the retail schema
// @Description get slice of sensors with loc, date, layout_id, store_id, device_id, kind, offset, limit, fields parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be sensor_id,device_id,layout_id,store_id,external_id,kind,title... all of them described at the model
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param store_id query string false "uuid format, default=*"
// @Param device_id query string false "uuid format, default=*"
// @Param kind query string false "default=*"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "sensor_id,external_id,kind... default=all"
// @Success 200 {object} infra.ChainSensorsResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors [get]
func (s *Server) apiChainSensors(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	kind := c.QueryParam("kind")
	storeID := c.QueryParam("store_id")
	deviceID := c.QueryParam("device_id")
	dt := c.QueryParam("date")
	if dt == "" {
		dt = time.Now().Format("2006-01-02")
	}
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckLayout read
	sensors, count, err := repo.FindChainSensors(c.Request().Context(), loc, dt, layoutID, storeID, deviceID, kind, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindChainSensors error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if sensors == nil {
		response := ChainSensorsResponse{
			Data: sensors,
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
	sensors.SetZeroValue(fields)

	response := ChainSensorsResponse{
		Data: sensors,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(sensors)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// nolint:lll
// apiChainSensorByID docs
// @Summary Get specified sensor in the retail schema
// @Description get sensor with loc, date, sensor_id parameters from retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description fields - comma separated values of field names, can be sensor_id,device_id,layout_id,store_id,external_id,kind,title... all of them described at the model
// @Produce json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param sensor_id path string true "uuid format"
// @Param loc query string false "default=ru"
// @Param date query string false "ISO8601 YYYY-MM-DD date"
// @Param fields query string false "sensor_id,external_id,kind... default=all"
// @Success 200 {object} domain.ChainSensor
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/{sensor_id} [get]
func (s *Server) apiChainSensorByID(c echo.Context) error {
	id := c.Param("sensor_id")
	if id == "" {
		s.log.Errorf("bad request apiChainSensorByID, %v", errEmptyID)
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
	sensor, err := repo.FindChainSensorByID(c.Request().Context(), loc, id, dt)
	if err != nil {
		s.log.Errorf("repo.FindChainSensors error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	if sensor == nil {
		s.log.Infof("for id=%s, loc=%s and date=%s doesn't have sensor", id, loc, dt)
		return c.JSON(http.StatusOK, sensor)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = nil
	}
	sensor.SetZeroValue(fields)
	return c.JSON(http.StatusOK, sensor)
}

// apiDeleteChainSensor docs
// @Summary Delete specified sensor in the retail schema
// @Description delete specified sensor by sensor_id parameter in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param sensor_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/{sensor_id} [delete]
func (s *Server) apiDeleteChainSensor(c echo.Context) error {
	id := c.Param("sensor_id")
	if id == "" {
		s.log.Errorf("bad request apiDeleteChainSensor, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainSensor(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.DelChainSensor error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d sensor(s)", cnt)))
}

// apiCreateBindChainSensorEntrance docs
// @Summary Create new binding sensor to entrance in the retail schema
// @Description create new binding sensor to entrance in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new binding sensor, entrance_id and sensor_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"direction": "forward",
// @Description 	"entrance_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"k_in": 1.0,
// @Description 	"k_out": 1.0,
// @Description 	"kind_entrance": "entrance",
// @Description 	"options": "{}",
// @Description 	"sensor_id": "9070be58-b4bf-467a-be64-6077edbd9867",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param params body domain.BindingChainSensorEntrance true "sensor bind entrance properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/entrance [post]
func (s *Server) apiCreateBindChainSensorEntrance(c echo.Context) error {
	bind := domain.BindingChainSensorEntrance{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiCreateBindChainSensorEntrance, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.BindChainSensorEntrance(c.Request().Context(), bind)
	if err != nil {
		s.log.Errorf("repo.BindChainSensorEntrance error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusCreated, CreatedStatus("binded"))
}

type UpdBindingChainSensorEntrance struct {
	SensorID      string    `json:"sensor_id,omitempty"`
	OldEntranceID string    `json:"old_entrance_id,omitempty"`
	NewEntranceID string    `json:"new_entrance_id,omitempty"`
	Direction     string    `json:"direction,omitempty"`
	KIn           float64   `json:"k_in,omitempty"`
	KOut          float64   `json:"k_out,omitempty"`
	Options       string    `json:"options,omitempty"`
	ValidFrom     time.Time `json:"valid_from,omitempty"`
	Modifier      string    `json:"modifier,omitempty"`
	NoHistory     bool      `json:"no_history"`
}

func (bse *UpdBindingChainSensorEntrance) makeDomainBindingChainSensorEntrance() domain.BindingChainSensorEntrance {
	return domain.BindingChainSensorEntrance{
		SensorID:   bse.SensorID,
		EntranceID: bse.OldEntranceID,
		Direction:  bse.Direction,
		KIn:        bse.KIn,
		KOut:       bse.KOut,
		Options:    bse.Options,
		ValidFrom:  &bse.ValidFrom,
	}
}

// apiUpdBindChainSensorEntrance docs
// @Summary Update binding sensor to entrance in the retail schema
// @Description update exists binding sensor to new entrance in retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example upd binding old_entrance_id is exists binding entrance, new_entrance_id is target entrance
// @Description valid_from set if need save history for specified date
// @Description and sensor_id must to be pass in parameters:
// @Description {
// @Description 	"sensor_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"old_entrance_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"new_entrance_id": "3595c9ad-8116-408e-a007-ec31d48f9667",
// @Description 	"direction": "forward",
// @Description 	"k_in": 1.0,
// @Description 	"k_out": 1.0,
// @Description 	"options": "{}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param params body infra.UpdBindingChainSensorEntrance true "sensor upd bind entrance properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/entrance [put]
func (s *Server) apiUpdBindChainSensorEntrance(c echo.Context) error {
	bind := UpdBindingChainSensorEntrance{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiUpdBindChainSensorEntrance, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if bind.SensorID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty sensor_is not allowed")))
	}
	if bind.OldEntranceID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty old_entrance_id not allowed")))
	}
	if bind.NewEntranceID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty new_entrance_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.UpdBindChainSensorEntrance(c.Request().Context(), bind.makeDomainBindingChainSensorEntrance(),
		bind.NewEntranceID, bind.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdBindChainSensorEntrance error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus("bind updated"))
}

// apiDeleteChainBindSensorEntrance docs
// @Summary Delete binding sensor to entrance in the retail schema
// @Description delete binding sensor to entrance in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param sensor_id path string true "uuid format"
// @Param entrance_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/entrance [delete]
func (s *Server) apiDeleteChainBindSensorEntrance(c echo.Context) error {
	sensor_id := c.QueryParam("sensor_id")
	entrance_id := c.QueryParam("entrance_id")
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainBindSensorEntrance(c.Request().Context(), sensor_id, entrance_id)
	if err != nil {
		s.log.Errorf("repo.DelChainBindSensorEntrance error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d row(s)", cnt)))
}

// apiCreateChainBindSensorZone docs
// @Summary Create new binding sensor to zone in the retail schema
// @Description create new binding sensor to zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example new binding sensor, zone_id and sensor_id must exists in the database:
// @Description {
// @Description 	"creator": "username",
// @Description 	"kind_zone": "zone",
// @Description 	"options": "{}",
// @Description 	"sensor_id": "9070be58-b4bf-467a-be64-6077edbd9867",
// @Description 	"zone_id": "2b076f28-c5d1-4fba-8b3b-2f58cba07c07",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00"
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param params body domain.BindingChainSensorZone true "sensor bind zone properties"
// @Success 201 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/zone [post]
func (s *Server) apiCreateChainBindSensorZone(c echo.Context) error {
	bind := domain.BindingChainSensorZone{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiCreateChainBindSensorZone, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.BindChainSensorZone(c.Request().Context(), bind)
	if err != nil {
		s.log.Errorf("repo.BindChainSensorZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusCreated, CreatedStatus("binded"))
}

// ChainBindSensorZoneResponse http wrapper with metadata
type ChainBindSensorZoneResponse struct {
	Data domain.BindingsChainSensorZone `json:"data"`
	Metadata
}

// apiGetChainBindSensorZone docs
// @Summary Get binding sensors to zones in the retail schema
// @Description get binding sensors to zones in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Produce json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param zone_id query string false "uuid/digits format"
// @Param sensor_id query string false "uuid/digits format"
// @Success 200 {object} infra.ChainBindSensorZoneResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/zone [get]
func (s *Server) apiGetChainBindSensorZone(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	sid := c.QueryParam("sensor_id")
	zid := c.QueryParam("zone_id")
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	binds, count, err := repo.FindBindChainSensorZone(c.Request().Context(), zid, sid, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindBindChainSensorZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	response := ChainBindSensorZoneResponse{
		Data: binds,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(binds)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

type UpdBindingChainSensorZone struct {
	SensorID  string    `json:"sensor_id,omitempty"`
	OldZoneID string    `json:"old_zone_id,omitempty"`
	NewZoneID string    `json:"new_zone_id,omitempty"`
	Options   string    `json:"options,omitempty"`
	ValidFrom time.Time `json:"valid_from,omitempty"`
	NoHistory bool      `json:"no_history"`
}

func (bsz *UpdBindingChainSensorZone) makeDomainBindingChainSensorZone() domain.BindingChainSensorZone {
	return domain.BindingChainSensorZone{
		SensorID:  bsz.SensorID,
		ZoneID:    bsz.OldZoneID,
		Options:   bsz.Options,
		ValidFrom: &bsz.ValidFrom,
	}
}

// apiUpdBindChainSensorZone docs
// @Summary Update binding sensor to zone in the retail schema
// @Description update exists binding sensor to new zone in retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description example upd binding old_zone_id is exists binding zone, new_zone_id is target zone
// @Description valid_from set if need save history for specified date
// @Description and sensor_id must to be pass in parameters:
// @Description {
// @Description 	"sensor_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"old_zone_id": "3595c9ad-8116-408e-a007-ec31d48f9669",
// @Description 	"new_zone_id": "3595c9ad-8116-408e-a007-ec31d48f9667",
// @Description 	"options": "{}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param params body infra.UpdBindingChainSensorZone true "sensor upd bind entrance properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/zone [put]
func (s *Server) apiUpdBindChainSensorZone(c echo.Context) error {
	bind := UpdBindingChainSensorZone{}
	if err := c.Bind(&bind); err != nil {
		s.log.Errorf("apiUpdBindChainSensorZone, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if bind.SensorID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty sensor_is not allowed")))
	}
	if bind.OldZoneID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty old_zone_id not allowed")))
	}
	if bind.NewZoneID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty new_zone_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	err = repo.UpdBindChainSensorZone(c.Request().Context(), bind.makeDomainBindingChainSensorZone(), bind.NewZoneID, bind.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdBindChainSensorZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus("bind updated"))
}

// apiDeleteChainBindSensorZone docs
// @Summary Delete binding sensor to zone in the retail schema
// @Description delete binding sensor to zone in the retail schema
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param sensor_id path string true "uuid format"
// @Param zone_id path string true "uuid format"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors/bindto/zone [delete]
func (s *Server) apiDeleteChainBindSensorZone(c echo.Context) error {
	sensor_id := c.QueryParam("sensor_id")
	zone_id := c.QueryParam("zone_id")
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	cnt, err := repo.DelChainBindSensorZone(c.Request().Context(), sensor_id, zone_id)
	if err != nil {
		s.log.Errorf("repo.DelChainBindSensorZone error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("deleted %d row(s)", cnt)))
}

// UpdChainSensor is upd sensor entity
type UpdChainSensor struct {
	SensorID   string     `json:"sensor_id,omitempty"`
	ExternalID string     `json:"external_id,omitempty"`
	Kind       string     `json:"kind,omitempty"`
	Title      string     `json:"title,omitempty"`
	Options    string     `json:"options,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	NoHistory  bool       `json:"no_history"`
}

func (cs *UpdChainSensor) makeDomainChainSensor() domain.ChainSensor {
	return domain.ChainSensor{
		SensorID:   cs.SensorID,
		ExternalID: cs.ExternalID,
		Kind:       cs.Kind,
		Title:      cs.Title,
		Options:    cs.Options,
		Notes:      cs.Notes,
		ValidFrom:  cs.ValidFrom,
	}
}

// nolint:lll
// apiUpdChainSensor docs
// @Summary Update sensor in the retail schema
// @Description update exists sensor
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description valid_from set if need save history for specified date
// @Description example upd sensor, sensor_id must to be pass in parameters:
// @Description {
// @Description 	"sensor_id": "8056fa1e-b63e-4d37-b014-744c4246621b",
// @Description 	"external_id":"Rule-2",
// @Description 	"kind": "sensor.people_count",
// @Description 	"notes": "{\\"ru\\":\\"описание/комментарии на русском языке\\",\\"en\\":\\"Comments/notes in English\\"}",
// @Description 	"options": "{\\"localIP\\":\\"192.168.0.10\\", \\"localPort\\": 80}",
// @Description 	"title": "{\\"ru\\":\\"Наименование сенсора на русском языке\\",\\"en\\":\\"Sensor name in English\\"}",
// @Description 	"valid_from": "2020-04-01T00:00:00+03:00",
// @Description 	"no_history":true
// @Description }
// @Accept  json
// @Produce  json
// @Tags chains/sensors
// @Param layout_id query string false "uuid format, default=*"
// @Param sensor body infra.UpdChainSensor true "sensor properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/sensors [put]
func (s *Server) apiUpdChainSensor(c echo.Context) error {
	sensor := &UpdChainSensor{}
	if err := c.Bind(sensor); err != nil {
		s.log.Errorf("apiUpdChainSensor, bad request error %v", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	if sensor.SensorID == "" {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(fmt.Errorf("empty sensor_id not allowed")))
	}
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	count, err := repo.UpdChainSensor(c.Request().Context(), sensor.makeDomainChainSensor(), sensor.NoHistory, userMock)
	if err != nil {
		s.log.Errorf("repo.UpdChainSensor error %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, OkStatus(fmt.Sprintf("update %d sensor(s)", count)))
}
