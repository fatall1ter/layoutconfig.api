// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domain

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain(in *jlexer.Lexer, out *Screenshots) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Screenshots, 0, 0)
			} else {
				*out = Screenshots{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Screenshot
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain(out *jwriter.Writer, in Screenshots) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Screenshots) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Screenshots) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Screenshots) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Screenshots) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain(l, v)
}
func easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(in *jlexer.Lexer, out *Screenshot) {
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
		case "layout_id":
			out.LayoutID = string(in.String())
		case "store_id":
			out.StoreID = string(in.String())
		case "screenshot_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ScreenshotTime).UnmarshalJSON(data))
			}
		case "screenshot_status":
			out.ScreenshotStatus = string(in.String())
		case "url":
			out.URL = string(in.String())
		case "url_aliases":
			if in.IsNull() {
				in.Skip()
				out.URLALiases = nil
			} else {
				in.Delim('[')
				if out.URLALiases == nil {
					if !in.IsDelim(']') {
						out.URLALiases = make([]string, 0, 4)
					} else {
						out.URLALiases = []string{}
					}
				} else {
					out.URLALiases = (out.URLALiases)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.URLALiases = append(out.URLALiases, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "layout_info":
			if in.IsNull() {
				in.Skip()
				out.LayoutInfo = nil
			} else {
				if out.LayoutInfo == nil {
					out.LayoutInfo = new(LayoutInfo)
				}
				(*out.LayoutInfo).UnmarshalEasyJSON(in)
			}
		case "notes":
			out.Notes = string(in.String())
		case "creator":
			out.Creator = string(in.String())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
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
func easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(out *jwriter.Writer, in Screenshot) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"device_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.DeviceID))
	}
	if in.LayoutID != "" {
		const prefix string = ",\"layout_id\":"
		out.RawString(prefix)
		out.String(string(in.LayoutID))
	}
	if in.StoreID != "" {
		const prefix string = ",\"store_id\":"
		out.RawString(prefix)
		out.String(string(in.StoreID))
	}
	{
		const prefix string = ",\"screenshot_time\":"
		out.RawString(prefix)
		out.Raw((in.ScreenshotTime).MarshalJSON())
	}
	if in.ScreenshotStatus != "" {
		const prefix string = ",\"screenshot_status\":"
		out.RawString(prefix)
		out.String(string(in.ScreenshotStatus))
	}
	if in.URL != "" {
		const prefix string = ",\"url\":"
		out.RawString(prefix)
		out.String(string(in.URL))
	}
	if len(in.URLALiases) != 0 {
		const prefix string = ",\"url_aliases\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v5, v6 := range in.URLALiases {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	if in.LayoutInfo != nil {
		const prefix string = ",\"layout_info\":"
		out.RawString(prefix)
		(*in.LayoutInfo).MarshalEasyJSON(out)
	}
	if in.Notes != "" {
		const prefix string = ",\"notes\":"
		out.RawString(prefix)
		out.String(string(in.Notes))
	}
	if in.Creator != "" {
		const prefix string = ",\"creator\":"
		out.RawString(prefix)
		out.String(string(in.Creator))
	}
	if true {
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Screenshot) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Screenshot) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Screenshot) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Screenshot) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(l, v)
}
func easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain2(in *jlexer.Lexer, out *ParamScreenUpd) {
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
		case "screenshot_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ScreenshotTime).UnmarshalJSON(data))
			}
		case "screenshot_status":
			out.ScreenshotStatus = string(in.String())
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
func easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain2(out *jwriter.Writer, in ParamScreenUpd) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"device_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.DeviceID))
	}
	{
		const prefix string = ",\"screenshot_time\":"
		out.RawString(prefix)
		out.Raw((in.ScreenshotTime).MarshalJSON())
	}
	{
		const prefix string = ",\"screenshot_status\":"
		out.RawString(prefix)
		out.String(string(in.ScreenshotStatus))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ParamScreenUpd) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ParamScreenUpd) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ParamScreenUpd) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ParamScreenUpd) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain2(l, v)
}
func easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain3(in *jlexer.Lexer, out *LayoutInfo) {
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
		case "params":
			if in.IsNull() {
				in.Skip()
				out.Params = nil
			} else {
				in.Delim('[')
				if out.Params == nil {
					if !in.IsDelim(']') {
						out.Params = make([]Param, 0, 2)
					} else {
						out.Params = []Param{}
					}
				} else {
					out.Params = (out.Params)[:0]
				}
				for !in.IsDelim(']') {
					var v7 Param
					(v7).UnmarshalEasyJSON(in)
					out.Params = append(out.Params, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain3(out *jwriter.Writer, in LayoutInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"layout_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.LayoutID))
	}
	{
		const prefix string = ",\"params\":"
		out.RawString(prefix)
		if in.Params == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Params {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LayoutInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LayoutInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB361eb5cEncodeGitCountmaxRuCountmaxLayoutconfigApiDomain3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LayoutInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LayoutInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB361eb5cDecodeGitCountmaxRuCountmaxLayoutconfigApiDomain3(l, v)
}
