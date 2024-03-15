package domain

// need generate...
type mallFields []string

var (
	allMallFields mallFields = mallFields{"layout_id", "kind", "title", "languages", "crm_key", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "read_only", "stores"}
)

func (fs mallFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs mallFields) invertList(list []string) []string {
	if len(list) == 0 {
		return fs.allFields()
	}
	res := make([]string, 0, len(fs)-len(list))
	for _, field := range fs {
		exists := false
		for _, item := range list {
			if item == field {
				exists = true
			}
		}
		if !exists {
			res = append(res, field)
		}
	}
	if len(res) == 0 {
		return fs.allFields()
	}
	return res
}

func (ch *Mall) setZeroByTagName(name string) {
	switch name {
	case "layout_id":
		ch.LayoutID = ""
	case "kind":
		ch.Kind = ""
	case "title":
		ch.Title = ""
	case "languages":
		ch.Languages = ""
	case "crm_key":
		ch.CRMKey = ""
	case "options":
		ch.Options = ""
	case "notes":
		ch.Notes = ""
	case "valid_from":
		ch.ValidFrom = nil
	case "valid_to":
		ch.ValidTo = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	case "modifier":
		ch.Modifier = nil
	case "modified_at":
		ch.ModifiedAt = nil
	case "read_only":
		ch.ReadOnly = false
	}
}

func (ch *Mall) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs Malls) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

// MallZone
// need generate...
type mallZoneFields []string

var (
	allMallZoneFields mallZoneFields = mallZoneFields{"zone_id", "parent_id", "layout_id", "kind", "title", "area", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "entrances", "sensors"}
)

func (fs mallZoneFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs mallZoneFields) invertList(list []string) []string {
	if len(list) == 0 {
		return fs.allFields()
	}
	res := make([]string, 0, len(fs)-len(list))
	for _, field := range fs {
		exists := false
		for _, item := range list {
			if item == field {
				exists = true
			}
		}
		if !exists {
			res = append(res, field)
		}
	}
	if len(res) == 0 {
		return fs.allFields()
	}
	return res
}

func (ch *MallZone) setZeroByTagName(name string) {
	switch name {
	case "zone_id":
		ch.ZoneID = ""
	case "parent_id":
		ch.ParentID = nil
	case "layout_id":
		ch.LayoutID = ""
	case "kind":
		ch.Kind = ""
	case "title":
		ch.Title = ""
	case "area":
		ch.Area = 0
	case "options":
		ch.Options = ""
	case "notes":
		ch.Notes = ""
	case "valid_from":
		ch.ValidFrom = nil
	case "valid_to":
		ch.ValidTo = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	case "modifier":
		ch.Modifier = nil
	case "modified_at":
		ch.ModifiedAt = nil
	case "sensors":
		ch.Sensors = nil
	case "entrances":
		ch.Entrances = nil
	}
}

func (ch *MallZone) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallZoneFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs MallZones) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallZoneFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

//----------------------------------------------------------------
// Renters
//----------------------------------------------------------------

type renterFields []string

var (
	allRenterFields renterFields = renterFields{"renter_id", "title", "layout_id", "category_id", "price_segment_id", "time_open", "time_close", "contract", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "zones"}
)

func (fs renterFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs renterFields) invertList(list []string) []string {
	if len(list) == 0 {
		return fs.allFields()
	}
	res := make([]string, 0, len(fs)-len(list))
	for _, field := range fs {
		exists := false
		for _, item := range list {
			if item == field {
				exists = true
			}
		}
		if !exists {
			res = append(res, field)
		}
	}
	if len(res) == 0 {
		return fs.allFields()
	}
	return res
}

func (ch *Renter) setZeroByTagName(name string) {
	switch name {
	case "renter_id":
		ch.RenterID = ""
	case "title":
		ch.Title = ""
	case "layout_id":
		ch.LayoutID = ""
	case "category_id":
		ch.CategoryID = ""
	case "preice_segment_id":
		ch.PriceSegmentID = ""
	case "time_open":
		ch.TimeOpen = nil
	case "time_close":
		ch.TimeClose = nil
	case "contract":
		ch.Contract = ""
	case "options":
		ch.Options = ""
	case "notes":
		ch.Notes = ""
	case "valid_from":
		ch.ValidFrom = nil
	case "valid_to":
		ch.ValidTo = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	case "modifier":
		ch.Modifier = nil
	case "modified_at":
		ch.ModifiedAt = nil
	case "zones":
		ch.Zones = nil
	}
}

func (ch *Renter) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allRenterFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs Renters) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allRenterFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

//----------------------------------------------------------------
// MallEntrances
//----------------------------------------------------------------

type mallEntranceFields []string

var (
	allMallEntranceFields mallEntranceFields = mallEntranceFields{"entrance_id", "layout_id", "floor_id", "kind", "title", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "sensors"}
)

func (fs mallEntranceFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs mallEntranceFields) invertList(list []string) []string {
	if len(list) == 0 {
		return fs.allFields()
	}
	res := make([]string, 0, len(fs)-len(list))
	for _, field := range fs {
		exists := false
		for _, item := range list {
			if item == field {
				exists = true
			}
		}
		if !exists {
			res = append(res, field)
		}
	}
	if len(res) == 0 {
		return fs.allFields()
	}
	return res
}

func (ch *MallEntrance) setZeroByTagName(name string) {
	switch name {
	case "entrance_id":
		ch.EntranceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "floor_id":
		ch.FloorID = ""
	case "kind":
		ch.Kind = ""
	case "title":
		ch.Title = ""
	case "options":
		ch.Options = ""
	case "notes":
		ch.Notes = ""
	case "valid_from":
		ch.ValidFrom = nil
	case "valid_to":
		ch.ValidTo = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	case "modifier":
		ch.Modifier = nil
	case "modified_at":
		ch.ModifiedAt = nil
	case "sensors":
		ch.Sensors = nil
	}
}

func (ch *MallEntrance) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallEntranceFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs MallEntrances) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallEntranceFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

// MallDevice

type mallDeviceFields []string

var (
	allMallDeviceFields mallDeviceFields = mallDeviceFields{"device_id", "layout_id", "floor_id", "master_id",
		"kind", "title", "is_active", "ip", "port", "sn", "mode", "dcmode", "login", "password", "options", "notes",
		"valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "sensors", "delay"}
)

func (fs mallDeviceFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs mallDeviceFields) invertList(list []string) []string {
	if len(list) == 0 {
		return fs.allFields()
	}
	res := make([]string, 0, len(fs)-len(list))
	for _, field := range fs {
		exists := false
		for _, item := range list {
			if item == field {
				exists = true
			}
		}
		if !exists {
			res = append(res, field)
		}
	}
	if len(res) == 0 {
		return fs.allFields()
	}
	return res
}

func (ch *MallDevice) setZeroByTagName(name string) {
	switch name {
	case "device_id":
		ch.DeviceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "floor_id":
		ch.FloorID = ""
	case "master_id":
		ch.MasterID = nil
	case "kind":
		ch.Kind = ""
	case "title":
		ch.Title = ""
	case "is_active":
		ch.IsActive = false
	case "ip":
		ch.IP = ""
	case "port":
		ch.Port = ""
	case "sn":
		ch.SN = ""
	case "mode":
		ch.Mode = ""
	case "dcmode":
		ch.DCMode = ""
	case "login":
		ch.Login = ""
	case "password":
		ch.Password = ""
	case "options":
		ch.Options = ""
	case "notes":
		ch.Notes = ""
	case "valid_from":
		ch.ValidFrom = nil
	case "valid_to":
		ch.ValidTo = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	case "modifier":
		ch.Modifier = nil
	case "modified_at":
		ch.ModifiedAt = nil
	case "sensors":
		ch.Sensors = nil
	case "delay":
		ch.Delay = nil
	}
}

func (ch *MallDevice) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallDeviceFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (ch *MallDevice) FillDelay(delays DeviceDelayDatas) {
	if len(delays) == 0 {
		return
	}
	for _, d := range delays {
		if d.DeviceID == ch.DeviceID {
			ch.Delay = &DelayPoint{Value: d.DelayMinutes}
		}
	}
}

func (chs *MallDevices) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallDeviceFields.invertList(fields)
	for i, c := range *chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		(*chs)[i] = c
	}
}

func (chs *MallDevices) IncludeSensors(ins MallSensors) {
	mapIN := make(map[string]MallSensors)
	for _, cs := range ins {
		if in, ok := mapIN[cs.DeviceID]; !ok {
			sts := make(MallSensors, 1)
			sts[0] = cs
			mapIN[cs.DeviceID] = sts
		} else {
			in = append(in, cs)
			mapIN[cs.DeviceID] = in
		}
	}
	for i, item := range *chs {
		(*chs)[i].Sensors = mapIN[item.DeviceID]
	}
}

func (chs *MallDevices) FindDeviceBySensorID(sid string) *MallDevice {
	for _, dev := range *chs {
		for _, s := range dev.Sensors {
			if s.SensorID == sid {
				return &dev
			}
		}
	}
	return nil
}

func (chs *MallDevices) SliceIDs() []string {
	result := make([]string, len(*chs))
	for i, d := range *chs {
		result[i] = d.DeviceID
	}
	return result
}

func (chs *MallDevices) FillDelays(delays DeviceDelayDatas) {
LOOP:
	for _, d := range delays {
		for i, dev := range *chs {
			if dev.DeviceID == d.DeviceID {
				dev.Delay = &DelayPoint{Value: d.DelayMinutes}
				(*chs)[i] = dev
				continue LOOP
			}
		}
	}
}

// ChainSensor
// need generate...
type mallSensorFields []string

var (
	allMallSensorFields mallSensorFields = mallSensorFields{"sensor_id", "device_id", "layout_id", "external_id",
		"kind", "title", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at"}
)

func (fs mallSensorFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs mallSensorFields) invertList(list []string) []string {
	if len(list) == 0 {
		return fs.allFields()
	}
	res := make([]string, 0, len(fs)-len(list))
	for _, field := range fs {
		exists := false
		for _, item := range list {
			if item == field {
				exists = true
			}
		}
		if !exists {
			res = append(res, field)
		}
	}
	if len(res) == 0 {
		return fs.allFields()
	}
	return res
}

func (ch *MallSensor) setZeroByTagName(name string) {
	switch name {
	case "sensor_id":
		ch.SensorID = ""
	case "device_id":
		ch.DeviceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "external_id":
		ch.ExternalID = ""
	case "kind":
		ch.Kind = ""
	case "title":
		ch.Title = ""
	case "options":
		ch.Options = ""
	case "notes":
		ch.Notes = ""
	case "valid_from":
		ch.ValidFrom = nil
	case "valid_to":
		ch.ValidTo = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	case "modifier":
		ch.Modifier = nil
	case "modified_at":
		ch.ModifiedAt = nil
	}
}

func (ch *MallSensor) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allMallSensorFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs MallSensors) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainSensorFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}
