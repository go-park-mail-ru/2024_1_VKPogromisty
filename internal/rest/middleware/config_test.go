package middleware

import (
	"net/http"
	"testing"
)

func TestCheckOrigin(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid origin",
			args: args{
				r: &http.Request{
					Header: map[string][]string{
						"Origin": {"http://localhost"},
					},
				},
			},
			want: true,
		},
		{
			name: "invalid origin",
			args: args{
				r: &http.Request{
					Header: map[string][]string{
						"Origin": {"http://invalid"},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckOrigin(tt.args.r); got != tt.want {
				t.Errorf("CheckOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}
