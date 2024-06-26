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

func easyjson6bf45bf4DecodeSocioDomain(in *jlexer.Lexer, out *Dialog) {
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
		case "user1":
			if in.IsNull() {
				in.Skip()
				out.User1 = nil
			} else {
				if out.User1 == nil {
					out.User1 = new(User)
				}
				easyjson6bf45bf4DecodeSocioDomain1(in, out.User1)
			}
		case "user2":
			if in.IsNull() {
				in.Skip()
				out.User2 = nil
			} else {
				if out.User2 == nil {
					out.User2 = new(User)
				}
				easyjson6bf45bf4DecodeSocioDomain1(in, out.User2)
			}
		case "lastMessage":
			if in.IsNull() {
				in.Skip()
				out.LastMessage = nil
			} else {
				if out.LastMessage == nil {
					out.LastMessage = new(PersonalMessage)
				}
				easyjson6bf45bf4DecodeSocioDomain2(in, out.LastMessage)
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
func easyjson6bf45bf4EncodeSocioDomain(out *jwriter.Writer, in Dialog) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user1\":"
		out.RawString(prefix[1:])
		if in.User1 == nil {
			out.RawString("null")
		} else {
			easyjson6bf45bf4EncodeSocioDomain1(out, *in.User1)
		}
	}
	{
		const prefix string = ",\"user2\":"
		out.RawString(prefix)
		if in.User2 == nil {
			out.RawString("null")
		} else {
			easyjson6bf45bf4EncodeSocioDomain1(out, *in.User2)
		}
	}
	{
		const prefix string = ",\"lastMessage\":"
		out.RawString(prefix)
		if in.LastMessage == nil {
			out.RawString("null")
		} else {
			easyjson6bf45bf4EncodeSocioDomain2(out, *in.LastMessage)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Dialog) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6bf45bf4EncodeSocioDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Dialog) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6bf45bf4EncodeSocioDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Dialog) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6bf45bf4DecodeSocioDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Dialog) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6bf45bf4DecodeSocioDomain(l, v)
}
func easyjson6bf45bf4DecodeSocioDomain2(in *jlexer.Lexer, out *PersonalMessage) {
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
		case "senderId":
			out.SenderID = uint(in.Uint())
		case "receiverId":
			out.ReceiverID = uint(in.Uint())
		case "content":
			out.Content = string(in.String())
		case "sticker":
			if in.IsNull() {
				in.Skip()
				out.Sticker = nil
			} else {
				if out.Sticker == nil {
					out.Sticker = new(Sticker)
				}
				easyjson6bf45bf4DecodeSocioDomain3(in, out.Sticker)
			}
		case "createdAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "updatedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
		case "attachments":
			if in.IsNull() {
				in.Skip()
				out.Attachments = nil
			} else {
				in.Delim('[')
				if out.Attachments == nil {
					if !in.IsDelim(']') {
						out.Attachments = make([]string, 0, 4)
					} else {
						out.Attachments = []string{}
					}
				} else {
					out.Attachments = (out.Attachments)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Attachments = append(out.Attachments, v1)
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
func easyjson6bf45bf4EncodeSocioDomain2(out *jwriter.Writer, in PersonalMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"senderId\":"
		out.RawString(prefix)
		out.Uint(uint(in.SenderID))
	}
	{
		const prefix string = ",\"receiverId\":"
		out.RawString(prefix)
		out.Uint(uint(in.ReceiverID))
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		out.String(string(in.Content))
	}
	if in.Sticker != nil {
		const prefix string = ",\"sticker\":"
		out.RawString(prefix)
		easyjson6bf45bf4EncodeSocioDomain3(out, *in.Sticker)
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
	{
		const prefix string = ",\"attachments\":"
		out.RawString(prefix)
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Attachments {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjson6bf45bf4DecodeSocioDomain3(in *jlexer.Lexer, out *Sticker) {
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
		case "authorId":
			out.AuthorID = uint(in.Uint())
		case "fileName":
			out.FileName = string(in.String())
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
func easyjson6bf45bf4EncodeSocioDomain3(out *jwriter.Writer, in Sticker) {
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
		const prefix string = ",\"authorId\":"
		out.RawString(prefix)
		out.Uint(uint(in.AuthorID))
	}
	{
		const prefix string = ",\"fileName\":"
		out.RawString(prefix)
		out.String(string(in.FileName))
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
func easyjson6bf45bf4DecodeSocioDomain1(in *jlexer.Lexer, out *User) {
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
			out.ID = uint(in.Uint())
		case "firstName":
			out.FirstName = string(in.String())
		case "lastName":
			out.LastName = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "dateOfBirth":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateOfBirth).UnmarshalJSON(data))
			}
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
func easyjson6bf45bf4EncodeSocioDomain1(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"userId\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"firstName\":"
		out.RawString(prefix)
		out.String(string(in.FirstName))
	}
	{
		const prefix string = ",\"lastName\":"
		out.RawString(prefix)
		out.String(string(in.LastName))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	if true {
		const prefix string = ",\"dateOfBirth\":"
		out.RawString(prefix)
		out.Raw((in.DateOfBirth).MarshalJSON())
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
