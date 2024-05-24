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

func easyjson6fbf8f0cDecodeSocioInternalRestSubscriptions(in *jlexer.Lexer, out *SubscriptionInput) {
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
		case "subscribedTo":
			out.SubscribedToID = uint(in.Uint())
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
func easyjson6fbf8f0cEncodeSocioInternalRestSubscriptions(out *jwriter.Writer, in SubscriptionInput) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"subscribedTo\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.SubscribedToID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SubscriptionInput) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6fbf8f0cEncodeSocioInternalRestSubscriptions(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SubscriptionInput) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6fbf8f0cEncodeSocioInternalRestSubscriptions(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SubscriptionInput) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6fbf8f0cDecodeSocioInternalRestSubscriptions(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SubscriptionInput) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6fbf8f0cDecodeSocioInternalRestSubscriptions(l, v)
}
