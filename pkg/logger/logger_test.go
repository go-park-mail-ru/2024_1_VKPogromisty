package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"socio/pkg/requestcontext"

	"github.com/google/uuid"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

type mockUnaryHandler struct {
	resp interface{}
	err  error
}

func (m *mockUnaryHandler) handle(ctx context.Context, req interface{}) (interface{}, error) {
	return m.resp, m.err
}

func TestUnaryLoggerInterceptor(t *testing.T) {
	logger, err := NewZapLogger([]string{"stderr"})
	if err != nil {
		t.Errorf("failed to create logger: %v", err)
	}

	l := NewLogger(logger)

	tests := []struct {
		name         string
		ctx          context.Context
		req          interface{}
		info         *grpc.UnaryServerInfo
		handler      *mockUnaryHandler
		expectedResp interface{}
		expectedErr  error
	}{
		{
			name:         "Test OK",
			ctx:          context.WithValue(context.Background(), requestcontext.RequestIDKey, uuid.New().String()),
			req:          "test request",
			info:         &grpc.UnaryServerInfo{FullMethod: "/test.method"},
			handler:      &mockUnaryHandler{resp: "test response", err: nil},
			expectedResp: "test response",
			expectedErr:  nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := l.UnaryLoggerInterceptor(tt.ctx, tt.req, tt.info, tt.handler.handle)

			if resp != tt.expectedResp {
				t.Errorf("handler returned unexpected response: got %v want %v",
					resp, tt.expectedResp)
			}

			if err != tt.expectedErr {
				t.Errorf("handler returned unexpected error: got %v want %v",
					err, tt.expectedErr)
			}
		})
	}
}

type mockHandler struct {
	resp string
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(m.resp))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TestLoggerMiddleware(t *testing.T) {
	logger := zaptest.NewLogger(t)
	l := NewLogger(logger)

	tests := []struct {
		name         string
		req          *http.Request
		handler      http.Handler
		expectedResp string
	}{
		{
			name:         "Test OK",
			req:          httptest.NewRequest("GET", "http://example.com/foo", nil),
			handler:      &mockHandler{resp: "OK"},
			expectedResp: "OK",
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			l.LoggerMiddleware(tt.handler).ServeHTTP(rr, tt.req)

			if rr.Body.String() != tt.expectedResp {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedResp)
			}

			// Here you can also check the logged values by inspecting the `logger` object
		})
	}
}
