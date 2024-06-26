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

func easyjsonEc7c1227DecodeSocioDomain(in *jlexer.Lexer, out *PublicGroupSubscription) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "publicGroupId":
			out.PublicGroupID = uint(in.Uint())
		case "subscriberId":
			out.SubscriberID = uint(in.Uint())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
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
func easyjsonEc7c1227EncodeSocioDomain(out *jwriter.Writer, in PublicGroupSubscription) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"publicGroupId\":"
		out.RawString(prefix)
		out.Uint(uint(in.PublicGroupID))
	}
	{
		const prefix string = ",\"subscriberId\":"
		out.RawString(prefix)
		out.Uint(uint(in.SubscriberID))
	}
	if true {
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	if true {
		const prefix string = ",\"updatedAt\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PublicGroupSubscription) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEc7c1227EncodeSocioDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PublicGroupSubscription) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEc7c1227EncodeSocioDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PublicGroupSubscription) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEc7c1227DecodeSocioDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PublicGroupSubscription) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEc7c1227DecodeSocioDomain(l, v)
}
func easyjsonEc7c1227DecodeSocioDomain1(in *jlexer.Lexer, out *PublicGroupAdmin) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "publicGroupId":
			out.PublicGroupID = uint(in.Uint())
		case "adminId":
			out.UserID = uint(in.Uint())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
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
func easyjsonEc7c1227EncodeSocioDomain1(out *jwriter.Writer, in PublicGroupAdmin) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"publicGroupId\":"
		out.RawString(prefix)
		out.Uint(uint(in.PublicGroupID))
	}
	{
		const prefix string = ",\"adminId\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	if true {
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	if true {
		const prefix string = ",\"updatedAt\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PublicGroupAdmin) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEc7c1227EncodeSocioDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PublicGroupAdmin) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEc7c1227EncodeSocioDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PublicGroupAdmin) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEc7c1227DecodeSocioDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PublicGroupAdmin) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEc7c1227DecodeSocioDomain1(l, v)
}
func easyjsonEc7c1227DecodeSocioDomain2(in *jlexer.Lexer, out *PublicGroup) {
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
		case "id":
			out.ID = uint(in.Uint())
		case "name":
			out.Name = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "subscribersCount":
			out.SubscribersCount = uint(in.Uint())
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
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
func easyjsonEc7c1227EncodeSocioDomain2(out *jwriter.Writer, in PublicGroup) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"subscribersCount\":"
		out.RawString(prefix)
		out.Uint(uint(in.SubscribersCount))
	}
	if true {
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	if true {
		const prefix string = ",\"updatedAt\":"
		out.RawString(prefix)
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PublicGroup) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEc7c1227EncodeSocioDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PublicGroup) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEc7c1227EncodeSocioDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PublicGroup) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEc7c1227DecodeSocioDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PublicGroup) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEc7c1227DecodeSocioDomain2(l, v)
}
