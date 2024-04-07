package sanitizer

import (
	"reflect"
	"socio/domain"
	"testing"

	"github.com/microcosm-cc/bluemonday"
)

func TestNewSanitizer(t *testing.T) {
	type args struct {
		sanitizer *bluemonday.Policy
	}
	tests := []struct {
		name string
		args args
		want *Sanitizer
	}{
		{
			name: "New sanitizer",
			args: args{
				sanitizer: bluemonday.UGCPolicy(),
			},
			want: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSanitizer(tt.args.sanitizer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSanitizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizer_SanitizePostWithAuthor(t *testing.T) {
	type args struct {
		post *domain.PostWithAuthor
	}
	tests := []struct {
		name string
		s    *Sanitizer
		args args
	}{
		{
			name: "Sanitize post with author",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				post: &domain.PostWithAuthor{
					Post: &domain.Post{
						Content: "<script>alert('Hello')</script>",
					},
					Author: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
				},
			},
		},
		{
			name: "nil postWithAuthor",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				post: nil,
			},
		},
		{
			name: "nil post",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				post: &domain.PostWithAuthor{
					Post: nil,
					Author: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SanitizePostWithAuthor(tt.args.post)
		})
	}
}

func TestSanitizer_SanitizeDialog(t *testing.T) {
	type args struct {
		dialog *domain.Dialog
	}
	tests := []struct {
		name string
		s    *Sanitizer
		args args
	}{
		{
			name: "Sanitize dialog",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				dialog: &domain.Dialog{
					User1: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
					User2: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
					LastMessage: &domain.PersonalMessage{
						Content: "<script>alert('Hello')</script>",
					},
				},
			},
		},
		{
			name: "nil dialog",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				dialog: nil,
			},
		},
		{
			name: "nil user",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				dialog: &domain.Dialog{
					User1: nil,
					User2: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
					LastMessage: &domain.PersonalMessage{
						Content: "<script>alert('Hello')</script>",
					},
				},
			},
		},
		{
			name: "nil message",
			s: &Sanitizer{
				sanitizer: bluemonday.UGCPolicy(),
			},
			args: args{
				dialog: &domain.Dialog{
					User1: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
					User2: &domain.User{
						FirstName: "<script>alert('Hello')</script>",
						LastName:  "<script>alert('Hello')</script>",
						Email:     "<script>alert('Hello')</script>",
						Avatar:    "<script>alert('Hello')</script>",
					},
					LastMessage: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SanitizeDialog(tt.args.dialog)
		})
	}
}
