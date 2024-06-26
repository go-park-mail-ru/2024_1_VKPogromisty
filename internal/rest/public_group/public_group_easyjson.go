// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package rest

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

func easyjsonEc7c1227DecodeSocioInternalRestPublicGroup(in *jlexer.Lexer, out *DeletePublicGroupAdminInput) {
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
		case "userId":
			out.UserID = uint64(in.Uint64())
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
func easyjsonEc7c1227EncodeSocioInternalRestPublicGroup(out *jwriter.Writer, in DeletePublicGroupAdminInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"userId\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v DeletePublicGroupAdminInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEc7c1227EncodeSocioInternalRestPublicGroup(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DeletePublicGroupAdminInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEc7c1227EncodeSocioInternalRestPublicGroup(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *DeletePublicGroupAdminInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEc7c1227DecodeSocioInternalRestPublicGroup(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DeletePublicGroupAdminInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEc7c1227DecodeSocioInternalRestPublicGroup(l, v)
}
func easyjsonEc7c1227DecodeSocioInternalRestPublicGroup1(in *jlexer.Lexer, out *CreatePublicGroupAdminInput) {
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
		case "userId":
			out.UserID = uint64(in.Uint64())
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
func easyjsonEc7c1227EncodeSocioInternalRestPublicGroup1(out *jwriter.Writer, in CreatePublicGroupAdminInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"userId\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreatePublicGroupAdminInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEc7c1227EncodeSocioInternalRestPublicGroup1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreatePublicGroupAdminInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEc7c1227EncodeSocioInternalRestPublicGroup1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreatePublicGroupAdminInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEc7c1227DecodeSocioInternalRestPublicGroup1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreatePublicGroupAdminInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEc7c1227DecodeSocioInternalRestPublicGroup1(l, v)
}
func easyjsonEc7c1227DecodeSocioInternalRestPublicGroup2(in *jlexer.Lexer, out *CheckIfUserIsAdminRes) {
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
		case "isAdmin":
			out.IsAdmin = bool(in.Bool())
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
func easyjsonEc7c1227EncodeSocioInternalRestPublicGroup2(out *jwriter.Writer, in CheckIfUserIsAdminRes) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"isAdmin\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.IsAdmin))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CheckIfUserIsAdminRes) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEc7c1227EncodeSocioInternalRestPublicGroup2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CheckIfUserIsAdminRes) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEc7c1227EncodeSocioInternalRestPublicGroup2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CheckIfUserIsAdminRes) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEc7c1227DecodeSocioInternalRestPublicGroup2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CheckIfUserIsAdminRes) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEc7c1227DecodeSocioInternalRestPublicGroup2(l, v)
}
