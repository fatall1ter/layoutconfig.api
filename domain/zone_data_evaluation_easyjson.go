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

func easyjson252a9d82DecodeGitCountmaxRuCountmaxLayoutconfigApiDomain(in *jlexer.Lexer, out *ZoneDataEvaluations) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ZoneDataEvaluations, 0, 0)
			} else {
				*out = ZoneDataEvaluations{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 ZoneDataEvaluation
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
func easyjson252a9d82EncodeGitCountmaxRuCountmaxLayoutconfigApiDomain(out *jwriter.Writer, in ZoneDataEvaluations) {
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
func (v ZoneDataEvaluations) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson252a9d82EncodeGitCountmaxRuCountmaxLayoutconfigApiDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ZoneDataEvaluations) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson252a9d82EncodeGitCountmaxRuCountmaxLayoutconfigApiDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ZoneDataEvaluations) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson252a9d82DecodeGitCountmaxRuCountmaxLayoutconfigApiDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ZoneDataEvaluations) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson252a9d82DecodeGitCountmaxRuCountmaxLayoutconfigApiDomain(l, v)
}
func easyjson252a9d82DecodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(in *jlexer.Lexer, out *ZoneDataEvaluation) {
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
		case "service_channel_block_id":
			out.ServiceChannelBlockID = string(in.String())
		case "record_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.RecordTime).UnmarshalJSON(data))
			}
		case "is_full":
			out.IsFull = bool(in.Bool())
		case "comment":
			out.Comment = string(in.String())
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
func easyjson252a9d82EncodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(out *jwriter.Writer, in ZoneDataEvaluation) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"layout_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.LayoutID))
	}
	{
		const prefix string = ",\"store_id\":"
		out.RawString(prefix)
		out.String(string(in.StoreID))
	}
	{
		const prefix string = ",\"service_channel_block_id\":"
		out.RawString(prefix)
		out.String(string(in.ServiceChannelBlockID))
	}
	{
		const prefix string = ",\"record_time\":"
		out.RawString(prefix)
		out.Raw((in.RecordTime).MarshalJSON())
	}
	{
		const prefix string = ",\"is_full\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsFull))
	}
	{
		const prefix string = ",\"comment\":"
		out.RawString(prefix)
		out.String(string(in.Comment))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ZoneDataEvaluation) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson252a9d82EncodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ZoneDataEvaluation) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson252a9d82EncodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ZoneDataEvaluation) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson252a9d82DecodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ZoneDataEvaluation) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson252a9d82DecodeGitCountmaxRuCountmaxLayoutconfigApiDomain1(l, v)
}
