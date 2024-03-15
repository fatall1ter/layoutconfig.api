// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package infra

import (
	json "encoding/json"
	domain "git.countmax.ru/countmax/layoutconfig.api/domain"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra(in *jlexer.Lexer, out *UpdChainDevice) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "device_id":
			out.DeviceID = string(in.String())
		case "master_id":
			if in.IsNull() {
				in.Skip()
				out.MasterID = nil
			} else {
				if out.MasterID == nil {
					out.MasterID = new(string)
				}
				*out.MasterID = string(in.String())
			}
		case "kind":
			out.Kind = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "is_active":
			out.IsActive = bool(in.Bool())
		case "ip":
			out.IP = string(in.String())
		case "port":
			out.Port = string(in.String())
		case "sn":
			out.SN = string(in.String())
		case "mode":
			out.Mode = string(in.String())
		case "dcmode":
			out.DCMode = string(in.String())
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "options":
			out.Options = string(in.String())
		case "notes":
			out.Notes = string(in.String())
		case "valid_from":
			if in.IsNull() {
				in.Skip()
				out.ValidFrom = nil
			} else {
				if out.ValidFrom == nil {
					out.ValidFrom = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ValidFrom).UnmarshalJSON(data))
				}
			}
		case "no_history":
			out.NoHistory = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra(out *jwriter.Writer, in UpdChainDevice) {
	out.RawByte('{')
	first := true
	_ = first
	if in.DeviceID != "" {
		const prefix string = ",\"device_id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.DeviceID))
	}
	if in.MasterID != nil {
		const prefix string = ",\"master_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.MasterID))
	}
	if in.Kind != "" {
		const prefix string = ",\"kind\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Kind))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.IsActive {
		const prefix string = ",\"is_active\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.IsActive))
	}
	if in.IP != "" {
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.IP))
	}
	if in.Port != "" {
		const prefix string = ",\"port\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Port))
	}
	if in.SN != "" {
		const prefix string = ",\"sn\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SN))
	}
	if in.Mode != "" {
		const prefix string = ",\"mode\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Mode))
	}
	if in.DCMode != "" {
		const prefix string = ",\"dcmode\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.DCMode))
	}
	if in.Login != "" {
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	if in.Options != "" {
		const prefix string = ",\"options\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Options))
	}
	if in.Notes != "" {
		const prefix string = ",\"notes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Notes))
	}
	if in.ValidFrom != nil {
		const prefix string = ",\"valid_from\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.ValidFrom).MarshalJSON())
	}
	{
		const prefix string = ",\"no_history\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.NoHistory))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpdChainDevice) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpdChainDevice) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpdChainDevice) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpdChainDevice) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra(l, v)
}
func easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra1(in *jlexer.Lexer, out *NewChainDevice) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "layout_id":
			out.LayoutID = string(in.String())
		case "store_id":
			out.StoreID = string(in.String())
		case "master_id":
			if in.IsNull() {
				in.Skip()
				out.MasterID = nil
			} else {
				if out.MasterID == nil {
					out.MasterID = new(string)
				}
				*out.MasterID = string(in.String())
			}
		case "kind":
			out.Kind = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "is_active":
			out.IsActive = bool(in.Bool())
		case "ip":
			out.IP = string(in.String())
		case "port":
			out.Port = string(in.String())
		case "sn":
			out.SN = string(in.String())
		case "mode":
			out.Mode = string(in.String())
		case "dcmode":
			out.DCMode = string(in.String())
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "options":
			out.Options = string(in.String())
		case "notes":
			out.Notes = string(in.String())
		case "valid_from":
			if in.IsNull() {
				in.Skip()
				out.ValidFrom = nil
			} else {
				if out.ValidFrom == nil {
					out.ValidFrom = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ValidFrom).UnmarshalJSON(data))
				}
			}
		case "valid_to":
			if in.IsNull() {
				in.Skip()
				out.ValidTo = nil
			} else {
				if out.ValidTo == nil {
					out.ValidTo = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ValidTo).UnmarshalJSON(data))
				}
			}
		case "creator":
			out.Creator = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra1(out *jwriter.Writer, in NewChainDevice) {
	out.RawByte('{')
	first := true
	_ = first
	if in.LayoutID != "" {
		const prefix string = ",\"layout_id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.LayoutID))
	}
	if in.StoreID != "" {
		const prefix string = ",\"store_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.StoreID))
	}
	if in.MasterID != nil {
		const prefix string = ",\"master_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.MasterID))
	}
	if in.Kind != "" {
		const prefix string = ",\"kind\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Kind))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.IsActive {
		const prefix string = ",\"is_active\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.IsActive))
	}
	if in.IP != "" {
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.IP))
	}
	if in.Port != "" {
		const prefix string = ",\"port\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Port))
	}
	if in.SN != "" {
		const prefix string = ",\"sn\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SN))
	}
	if in.Mode != "" {
		const prefix string = ",\"mode\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Mode))
	}
	if in.DCMode != "" {
		const prefix string = ",\"dcmode\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.DCMode))
	}
	if in.Login != "" {
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	if in.Options != "" {
		const prefix string = ",\"options\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Options))
	}
	if in.Notes != "" {
		const prefix string = ",\"notes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Notes))
	}
	if in.ValidFrom != nil {
		const prefix string = ",\"valid_from\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.ValidFrom).MarshalJSON())
	}
	if in.ValidTo != nil {
		const prefix string = ",\"valid_to\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.ValidTo).MarshalJSON())
	}
	if in.Creator != "" {
		const prefix string = ",\"creator\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Creator))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NewChainDevice) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NewChainDevice) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NewChainDevice) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NewChainDevice) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra1(l, v)
}
func easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra2(in *jlexer.Lexer, out *ChainDevicesResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				in.Delim('[')
				if out.Data == nil {
					if !in.IsDelim(']') {
						out.Data = make(domain.ChainDevices, 0, 0)
					} else {
						out.Data = domain.ChainDevices{}
					}
				} else {
					out.Data = (out.Data)[:0]
				}
				for !in.IsDelim(']') {
					var v1 domain.ChainDevice
					(v1).UnmarshalEasyJSON(in)
					out.Data = append(out.Data, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "result_set":
			easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra3(in, &out.ResultSet)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra2(out *jwriter.Writer, in ChainDevicesResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		if in.Data == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Data {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"result_set\":"
		out.RawString(prefix)
		easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra3(out, in.ResultSet)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChainDevicesResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChainDevicesResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChainDevicesResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChainDevicesResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra2(l, v)
}
func easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra3(in *jlexer.Lexer, out *ResultSet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "count":
			out.Count = int64(in.Int64())
		case "offset":
			out.Offset = int64(in.Int64())
		case "limit":
			out.Limit = int64(in.Int64())
		case "total":
			out.Total = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra3(out *jwriter.Writer, in ResultSet) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Count))
	}
	{
		const prefix string = ",\"offset\":"
		out.RawString(prefix)
		out.Int64(int64(in.Offset))
	}
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix)
		out.Int64(int64(in.Limit))
	}
	{
		const prefix string = ",\"total\":"
		out.RawString(prefix)
		out.Int64(int64(in.Total))
	}
	out.RawByte('}')
}
func easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra4(in *jlexer.Lexer, out *ChainDeviceTracksResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			(out.Data).UnmarshalEasyJSON(in)
		case "result_set":
			easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra3(in, &out.ResultSet)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra4(out *jwriter.Writer, in ChainDeviceTracksResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		(in.Data).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"result_set\":"
		out.RawString(prefix)
		easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra3(out, in.ResultSet)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChainDeviceTracksResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChainDeviceTracksResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC3ba8e86EncodeGitCountmaxRuCountmaxLayoutconfigApiInfra4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChainDeviceTracksResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChainDeviceTracksResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC3ba8e86DecodeGitCountmaxRuCountmaxLayoutconfigApiInfra4(l, v)
}