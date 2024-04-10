package contextlogger_test

import (
	"context"
	"errors"
	"socio/pkg/contextlogger"
	"socio/pkg/requestcontext"
	"testing"

	"go.uber.org/zap"
)

func TestLogInfo(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		return
	}

	sugar := logger.Sugar()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "error",
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.WithValue(
					context.Background(),
					requestcontext.LoggerKey,
					sugar.With("tya", "zhelo"),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contextlogger.LogInfo(tt.args.ctx)
		})
	}
}

func TestLogErr(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		return
	}

	sugar := logger.Sugar()

	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				err: nil,
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.WithValue(
					context.Background(),
					requestcontext.LoggerKey,
					sugar.With("tya", "zhelo"),
				),
				err: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contextlogger.LogErr(tt.args.ctx, tt.args.err)
		})
	}
}

func TestLogSQL(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		return
	}

	sugar := logger.Sugar()

	type args struct {
		ctx   context.Context
		query string
		args  []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "error",
			args: args{
				ctx:   context.Background(),
				query: "",
				args:  nil,
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.WithValue(
					context.Background(),
					requestcontext.LoggerKey,
					sugar.With("tya", "zhelo"),
				),
				query: "SELECT * FROM users",
				args:  []interface{}{"arg1", "arg2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contextlogger.LogSQL(tt.args.ctx, tt.args.query, tt.args.args...)
		})
	}
}

func TestLogRedisAction(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		return
	}

	sugar := logger.Sugar()

	type args struct {
		ctx    context.Context
		action string
		key    interface{}
		value  interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "error",
			args: args{
				ctx:    context.Background(),
				action: "",
				key:    nil,
				value:  nil,
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.WithValue(
					context.Background(),
					requestcontext.LoggerKey,
					sugar.With("tya", "zhelo"),
				),
				action: "GET",
				key:    "key",
				value:  "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contextlogger.LogRedisAction(tt.args.ctx, tt.args.action, tt.args.key, tt.args.value)
		})
	}
}
