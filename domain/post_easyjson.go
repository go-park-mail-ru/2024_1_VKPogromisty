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

func easyjson5a72dc82DecodeSocioDomain(in *jlexer.Lexer, out *PostWithAuthorAndGroup) {
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
		case "post":
			if in.IsNull() {
				in.Skip()
				out.Post = nil
			} else {
				if out.Post == nil {
					out.Post = new(Post)
				}
				(*out.Post).UnmarshalEasyJSON(in)
			}
		case "author":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(User)
				}
				easyjson5a72dc82DecodeSocioDomain1(in, out.Author)
			}
		case "group":
			if in.IsNull() {
				in.Skip()
				out.Group = nil
			} else {
				if out.Group == nil {
					out.Group = new(PublicGroup)
				}
				easyjson5a72dc82DecodeSocioDomain2(in, out.Group)
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
func easyjson5a72dc82EncodeSocioDomain(out *jwriter.Writer, in PostWithAuthorAndGroup) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix[1:])
		if in.Post == nil {
			out.RawString("null")
		} else {
			(*in.Post).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		if in.Author == nil {
			out.RawString("null")
		} else {
			easyjson5a72dc82EncodeSocioDomain1(out, *in.Author)
		}
	}
	{
		const prefix string = ",\"group\":"
		out.RawString(prefix)
		if in.Group == nil {
			out.RawString("null")
		} else {
			easyjson5a72dc82EncodeSocioDomain2(out, *in.Group)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostWithAuthorAndGroup) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeSocioDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostWithAuthorAndGroup) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeSocioDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostWithAuthorAndGroup) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeSocioDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostWithAuthorAndGroup) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeSocioDomain(l, v)
}
func easyjson5a72dc82DecodeSocioDomain2(in *jlexer.Lexer, out *PublicGroup) {
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
func easyjson5a72dc82EncodeSocioDomain2(out *jwriter.Writer, in PublicGroup) {
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
func easyjson5a72dc82DecodeSocioDomain1(in *jlexer.Lexer, out *User) {
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
func easyjson5a72dc82EncodeSocioDomain1(out *jwriter.Writer, in User) {
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
func easyjson5a72dc82DecodeSocioDomain3(in *jlexer.Lexer, out *PostWithAuthor) {
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
		case "post":
			if in.IsNull() {
				in.Skip()
				out.Post = nil
			} else {
				if out.Post == nil {
					out.Post = new(Post)
				}
				(*out.Post).UnmarshalEasyJSON(in)
			}
		case "author":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(User)
				}
				easyjson5a72dc82DecodeSocioDomain1(in, out.Author)
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
func easyjson5a72dc82EncodeSocioDomain3(out *jwriter.Writer, in PostWithAuthor) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix[1:])
		if in.Post == nil {
			out.RawString("null")
		} else {
			(*in.Post).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		if in.Author == nil {
			out.RawString("null")
		} else {
			easyjson5a72dc82EncodeSocioDomain1(out, *in.Author)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostWithAuthor) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeSocioDomain3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostWithAuthor) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeSocioDomain3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostWithAuthor) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeSocioDomain3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostWithAuthor) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeSocioDomain3(l, v)
}
func easyjson5a72dc82DecodeSocioDomain4(in *jlexer.Lexer, out *PostLike) {
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
		case "likeId":
			out.ID = uint(in.Uint())
		case "postId":
			out.PostID = uint(in.Uint())
		case "userId":
			out.UserID = uint(in.Uint())
		case "createdAt":
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
func easyjson5a72dc82EncodeSocioDomain4(out *jwriter.Writer, in PostLike) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"likeId\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"postId\":"
		out.RawString(prefix)
		out.Uint(uint(in.PostID))
	}
	{
		const prefix string = ",\"userId\":"
		out.RawString(prefix)
		out.Uint(uint(in.UserID))
	}
	if true {
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostLike) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeSocioDomain4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostLike) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeSocioDomain4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostLike) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeSocioDomain4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostLike) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeSocioDomain4(l, v)
}
func easyjson5a72dc82DecodeSocioDomain5(in *jlexer.Lexer, out *Post) {
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
		case "postId":
			out.ID = uint(in.Uint())
		case "authorId":
			out.AuthorID = uint(in.Uint())
		case "groupId":
			out.GroupID = uint(in.Uint())
		case "content":
			out.Content = string(in.String())
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
		case "likedBy":
			if in.IsNull() {
				in.Skip()
				out.LikedByIDs = nil
			} else {
				in.Delim('[')
				if out.LikedByIDs == nil {
					if !in.IsDelim(']') {
						out.LikedByIDs = make([]uint64, 0, 8)
					} else {
						out.LikedByIDs = []uint64{}
					}
				} else {
					out.LikedByIDs = (out.LikedByIDs)[:0]
				}
				for !in.IsDelim(']') {
					var v2 uint64
					v2 = uint64(in.Uint64())
					out.LikedByIDs = append(out.LikedByIDs, v2)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson5a72dc82EncodeSocioDomain5(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"postId\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"authorId\":"
		out.RawString(prefix)
		out.Uint(uint(in.AuthorID))
	}
	if in.GroupID != 0 {
		const prefix string = ",\"groupId\":"
		out.RawString(prefix)
		out.Uint(uint(in.GroupID))
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		out.String(string(in.Content))
	}
	{
		const prefix string = ",\"attachments\":"
		out.RawString(prefix)
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Attachments {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"likedBy\":"
		out.RawString(prefix)
		if in.LikedByIDs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.LikedByIDs {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v6))
			}
			out.RawByte(']')
		}
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
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeSocioDomain5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeSocioDomain5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeSocioDomain5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeSocioDomain5(l, v)
}
func easyjson5a72dc82DecodeSocioDomain6(in *jlexer.Lexer, out *GroupPost) {
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
		case "postId":
			out.PostID = uint(in.Uint())
		case "groupId":
			out.GroupID = uint(in.Uint())
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
func easyjson5a72dc82EncodeSocioDomain6(out *jwriter.Writer, in GroupPost) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.ID))
	}
	{
		const prefix string = ",\"postId\":"
		out.RawString(prefix)
		out.Uint(uint(in.PostID))
	}
	{
		const prefix string = ",\"groupId\":"
		out.RawString(prefix)
		out.Uint(uint(in.GroupID))
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
func (v GroupPost) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeSocioDomain6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GroupPost) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeSocioDomain6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GroupPost) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeSocioDomain6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GroupPost) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeSocioDomain6(l, v)
}