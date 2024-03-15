package domain

// need generate...
type chainFields []string

var (
	allChainFields chainFields = chainFields{"layout_id", "kind", "title", "languages", "crm_key", "brands", "currency", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "read_only", "stores"}
)

func (fs chainFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs chainFields) invertList(list []string) []string {
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

func (ch *Chain) setZeroByTagName(name string) {
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
	case "brands":
		ch.Brands = ""
	case "currency":
		ch.Currency = ""
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
	case "stores":
		ch.Stores = nil
	}
}

func (ch *Chain) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs Chains) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

func (chs Chains) IncludeStores(stores ChainStores) {
	storesInChain := make(map[string]ChainStores)
	for _, store := range stores {
		if instores, ok := storesInChain[store.LayoutID]; !ok {
			sts := make(ChainStores, 1)
			sts[0] = store
			storesInChain[store.LayoutID] = sts
		} else {
			instores = append(instores, store)
			storesInChain[store.LayoutID] = instores
		}
	}
	for i, chain := range chs {
		chs[i].Stores = storesInChain[chain.LayoutID]
	}
}

// ChainStore
// need generate...
type chainStoreFields []string

var (
	allChainStoreFields chainStoreFields = chainStoreFields{"store_id", "layout_id", "kind", "title", "crm_key", "brands", "statistics", "location_id", "area", "currency", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "entrances", "zones", "devices"}
)

func (fs chainStoreFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs chainStoreFields) invertList(list []string) []string {
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

func (ch *ChainStore) setZeroByTagName(name string) {
	switch name {
	case "store_id":
		ch.StoreID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "kind":
		ch.Kind = ""
	case "title":
		ch.Title = ""
	case "crm_key":
		ch.CRMKey = ""
	case "brands":
		ch.Brands = ""
	case "location_id":
		ch.LocationID = ""
	case "area":
		ch.Area = 0.0
	case "statistics":
		ch.Statistics = ""
	case "currency":
		ch.Currency = ""
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
	case "entrances":
		ch.Entrances = nil
	case "zones":
		ch.Zones = nil
	case "devices":
		ch.Devices = nil
	}
}

func (ch *ChainStore) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainStoreFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (ch *ChainStore) IncludeDevices(ins ChainDevices) {
	mapIN := make(map[string]ChainDevices)
	for _, inItem := range ins {
		if in, ok := mapIN[inItem.StoreID]; !ok {
			sts := make(ChainDevices, 1)
			sts[0] = inItem
			mapIN[inItem.StoreID] = sts
		} else {
			in = append(in, inItem)
			mapIN[inItem.StoreID] = in
		}
	}
	ch.Devices = mapIN[ch.StoreID]
}

func (chs ChainStores) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainStoreFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

func (chs ChainStores) IncludeEntrances(ins ChainEntrances) {
	mapIN := make(map[string]ChainEntrances)
	for _, inItem := range ins {
		if in, ok := mapIN[inItem.StoreID]; !ok {
			sts := make(ChainEntrances, 1)
			sts[0] = inItem
			mapIN[inItem.StoreID] = sts
		} else {
			in = append(in, inItem)
			mapIN[inItem.StoreID] = in
		}
	}
	for i, item := range chs {
		chs[i].Entrances = mapIN[item.StoreID]
	}
}

func (chs ChainStores) IncludeZones(ins ChainZones) { // TODO: make tests!! and check USE POINTER!!!!
	mapIN := make(map[string]ChainZones)
	for _, inItem := range ins {
		if in, ok := mapIN[inItem.StoreID]; !ok {
			sts := make(ChainZones, 1)
			sts[0] = inItem
			mapIN[inItem.StoreID] = sts
		} else {
			in = append(in, inItem)
			mapIN[inItem.StoreID] = in
		}
	}
	for i, item := range chs {
		chs[i].Zones = mapIN[item.StoreID]
	}
}

func (chs ChainStores) IncludeDevices(ins ChainDevices) {
	mapIN := make(map[string]ChainDevices)
	for _, inItem := range ins {
		if in, ok := mapIN[inItem.StoreID]; !ok {
			sts := make(ChainDevices, 1)
			sts[0] = inItem
			mapIN[inItem.StoreID] = sts
		} else {
			in = append(in, inItem)
			mapIN[inItem.StoreID] = in
		}
	}
	for i, item := range chs {
		chs[i].Devices = mapIN[item.StoreID]
	}
}

// ChainEntrance
// need generate...
type chainEntranceFields []string

var (
	allChainEntranceFields chainEntranceFields = chainEntranceFields{"entrance_id", "layout_id", "store_id", "kind", "title", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "sensors"}
)

func (fs chainEntranceFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs chainEntranceFields) invertList(list []string) []string {
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

func (ch *ChainEntrance) setZeroByTagName(name string) {
	switch name {
	case "entrance_id":
		ch.EntranceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "store_id":
		ch.StoreID = ""
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

func (ch *ChainEntrance) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainEntranceFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs ChainEntrances) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainEntranceFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

// ChainZone
// need generate...
type chainZoneFields []string

var (
	allChainZoneFields chainZoneFields = chainZoneFields{"zone_id", "parent_id", "layout_id", "store_id", "kind", "title", "area", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "entrances", "sensors"}
)

func (fs chainZoneFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs chainZoneFields) invertList(list []string) []string {
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

func (ch *ChainZone) setZeroByTagName(name string) {
	switch name {
	case "zone_id":
		ch.ZoneID = ""
	case "parent_id":
		ch.ParentID = nil
	case "layout_id":
		ch.LayoutID = ""
	case "store_id":
		ch.StoreID = ""
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

func (ch *ChainZone) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainZoneFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs ChainZones) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainZoneFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

// ChainDevice
// need generate...
type chainDeviceFields []string

var (
	allChainDeviceFields chainDeviceFields = chainDeviceFields{"device_id", "layout_id", "store_id", "master_id",
		"kind", "title", "is_active", "ip", "port", "sn", "mode", "dcmode", "login", "password", "options", "notes",
		"valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at", "sensors", "delay"}
)

func (fs chainDeviceFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs chainDeviceFields) invertList(list []string) []string {
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

func (ch *ChainDevice) setZeroByTagName(name string) {
	switch name {
	case "device_id":
		ch.DeviceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "store_id":
		ch.StoreID = ""
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

func (ch *ChainDevice) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainDeviceFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (ch *ChainDevice) FillDelay(delays DeviceDelayDatas) {
	if len(delays) == 0 {
		return
	}
	for _, d := range delays {
		if d.DeviceID == ch.DeviceID {
			ch.Delay = &DelayPoint{Value: d.DelayMinutes}
		}
	}
}

func (chs *ChainDevices) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainDeviceFields.invertList(fields)
	for i, c := range *chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		(*chs)[i] = c
	}
}

func (chs *ChainDevices) IncludeSensors(ins ChainSensors) {
	mapIN := make(map[string]ChainSensors)
	for _, cs := range ins {
		if in, ok := mapIN[cs.DeviceID]; !ok {
			sts := make(ChainSensors, 1)
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

func (chs *ChainDevices) FindDeviceBySensorID(sid string) *ChainDevice {
	for _, dev := range *chs {
		for _, s := range dev.Sensors {
			if s.SensorID == sid {
				return &dev
			}
		}
	}
	return nil
}

func (chs *ChainDevices) SliceIDs() []string {
	result := make([]string, len(*chs))
	for i, d := range *chs {
		result[i] = d.DeviceID
	}
	return result
}

func (chs *ChainDevices) FillDelays(delays DeviceDelayDatas) {
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
type chainSensorFields []string

var (
	allChainSensorFields chainSensorFields = chainSensorFields{"sensor_id", "device_id", "layout_id", "store_id", "external_id",
		"kind", "title", "options", "notes", "valid_from", "valid_to", "creator", "created_at", "modifier", "modified_at"}
)

func (fs chainSensorFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs chainSensorFields) invertList(list []string) []string {
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

func (ch *ChainSensor) setZeroByTagName(name string) {
	switch name {
	case "sensor_id":
		ch.SensorID = ""
	case "device_id":
		ch.DeviceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "store_id":
		ch.StoreID = ""
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

func (ch *ChainSensor) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allChainSensorFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs ChainSensors) SetZeroValue(fields []string) {
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

// SCREENSHOTS

type screenshotFields []string

var (
	allScreenshotFields screenshotFields = screenshotFields{"device_id", "layout_id", "store_id", "screenshot_time", "screenshot_status", "url", "layout_info", "notes", "creator", "created_at"}
)

func (fs screenshotFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs screenshotFields) invertList(list []string) []string {
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

func (ch *Screenshot) setZeroByTagName(name string) {
	switch name {
	case "device_id":
		ch.DeviceID = ""
	case "layout_id":
		ch.LayoutID = ""
	case "layout_info":
		ch.LayoutInfo = nil
	case "store_id":
		ch.StoreID = ""
	case "screenshot_status":
		ch.ScreenshotStatus = ""
	case "url":
		ch.URL = ""
	case "notes":
		ch.Notes = ""
	case "creator":
		ch.Creator = ""
	}
}

func (ch *Screenshot) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allScreenshotFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs Screenshots) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allScreenshotFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}

// EVENTS

type eventFields []string

var (
	allEventFields eventFields = eventFields{"id", "key", "event_time", "kind", "message", "severity", "layout_id", "store_id", "source", "creator", "created_at"}
)

func (fs eventFields) allFields() []string {
	res := make([]string, len(fs))
	for i, f := range fs {
		res[i] = f
	}
	return res
}

func (fs eventFields) invertList(list []string) []string {
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

func (ch *Event) setZeroByTagName(name string) {
	switch name {
	case "key":
		ch.Key = ""
	case "kind":
		ch.Kind = ""
	case "message":
		ch.Message = ""
	case "severity":
		ch.Severity = ""
	case "layout_id":
		ch.LayoutID = ""
	case "store_id":
		ch.StoreID = ""
	case "source":
		ch.Source = nil
	case "creator":
		ch.Creator = ""
	case "created_at":
		ch.CreatedAt = nil
	}
}

func (ch *Event) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allEventFields.invertList(fields)
	for _, tag := range defList {
		ch.setZeroByTagName(tag)
	}
}

func (chs Events) SetZeroValue(fields []string) {
	if len(fields) == 0 {
		return
	}
	defList := allEventFields.invertList(fields)
	for i, c := range chs {
		for _, tag := range defList {
			c.setZeroByTagName(tag)
		}
		chs[i] = c
	}
}
