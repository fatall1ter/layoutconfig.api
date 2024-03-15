package domain

import "time"

// Screenshot image, extracted or getted from device and putted to S3 storage,
// link to that image in the storage contains in the URL parameter
type Screenshot struct {
	DeviceID         string      `json:"device_id"`                   // text not null, -- идентификатор устройства, с которого получен скриншот
	LayoutID         string      `json:"layout_id,omitempty"`         // text null, -- идентификатор схемы размещения, которой принадлежит устройство с которого получен скриншот
	StoreID          string      `json:"store_id,omitempty"`          // text not null, -- идентификатор магазина, в котором находилось устройство с которого получен скриншот
	ScreenshotTime   time.Time   `json:"screenshot_time"`             // stamptz not null, -- дата-время создания скриншота
	ScreenshotStatus string      `json:"screenshot_status,omitempty"` // статус скриншота: new-по умолчанию, processed-блокировка от изменений другими пользователями (не более 1 минуты), to_delete-помечен для удаления, archived-запись для вечного хранения
	URL              string      `json:"url,omitempty"`               // text not null, -- адрес местонахождения скриншота
	URLALiases       []string    `json:"url_aliases,omitempty"`       // api layer only, alt links to screenshots
	LayoutInfo       *LayoutInfo `json:"layout_info,omitempty"`       // jsonb not null default '{}'::jsonb, -- более детальная проектная информация устройства
	Notes            string      `json:"notes,omitempty"`             // text null, -- заметки
	Creator          string      `json:"creator,omitempty"`           // varchar(128) not null default current_user, -- пользователь, создавший запись
	CreatedAt        time.Time   `json:"created_at,omitempty"`        // timestamptz not null default current_timestamp, -- дата-время создания записи
}

// LayoutInfo data about layout in k/v form
type LayoutInfo struct {
	LayoutID string  `json:"layout_id"`
	Params   []Param `json:"params"`
}

// GetParamByName helper to extract specified parameter from LayoutInfo
func (li *LayoutInfo) GetParamByName(name string) string {
	for _, p := range li.Params {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

// SetParamByName set value of param in the LayoutInfo
func (li *LayoutInfo) SetParamByName(name, value string) {
	if par := li.GetParamByName(name); par == "" {
		li.Params = append(li.Params, Param{name, value})
	}
}

//easyjson:json
type Screenshots []Screenshot

// ParamScreenUpd properties for update state of screenshot
type ParamScreenUpd struct {
	DeviceID         string    `json:"device_id"`
	ScreenshotTime   time.Time `json:"screenshot_time"`
	ScreenshotStatus string    `json:"screenshot_status"`
}

type ParamsScreenUpd []ParamScreenUpd

type IScreenRepo interface {
	FindScreens(layoutID, storeID, status string, deviceID []string, from, to time.Time, limit, offset int64) (Screenshots, int64, error)
	FindScreensAtTime(layoutID, storeID string, deviceID []string, t time.Time) (Screenshots, error)
	UpdStatusManyScreens(ParamsScreenUpd) (int64, error)
}
