package bootstrap

import (
	"context"
	"errors"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	goplugin "github.com/hashicorp/go-plugin"

	pluginerrors "github.com/openkcm/plugin-sdk/api/plugin-errors"
	pluginoption "github.com/openkcm/plugin-sdk/api/plugin-option"
	initv1 "github.com/openkcm/plugin-sdk/internal/proto/service/init/v1"
)

type mockValidator struct {
	err error
}

func (v *mockValidator) Validate(_ proto.Message, _ ...protovalidate.ValidationOption) error {
	return v.err
}

func cancelledTestConfig() *goplugin.ServeTestConfig {
	// Create a context that is already cancelled to prevent the server from actually starting during tests.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return &goplugin.ServeTestConfig{Context: ctx}
}

func TestServe_NoPluginServer(t *testing.T) {
	err := Serve()
	assert.ErrorIs(t, err, pluginerrors.ErrServerRequired)
}

func TestServe_WithValidationEnabled(t *testing.T) {
	mock := &pluginMock{typ: "test"}

	err := Serve(
		pluginoption.WithPluginServer(mock),
		pluginoption.EnableInputValidation(),
		pluginoption.WithTestConfig(cancelledTestConfig()),
	)
	assert.NoError(t, err)

	err = Serve(
		pluginoption.WithPluginServer(mock),
		pluginoption.EnableOutputValidation(),
		pluginoption.WithTestConfig(cancelledTestConfig()),
	)
	assert.NoError(t, err)
}

func TestServe_WithValidationAndExistingServerOptions(t *testing.T) {
	mock := &pluginMock{typ: "test"}
	err := Serve(
		pluginoption.WithPluginServer(mock),
		pluginoption.SetServerOption(grpc.MaxRecvMsgSize(1024)),
		pluginoption.EnableInputValidation(),
		pluginoption.WithTestConfig(cancelledTestConfig()),
	)
	assert.NoError(t, err)
}

func TestHCServer_GRPCServer(t *testing.T) {
	mock := &pluginMock{typ: "test"}
	p := newHCPlugin(hclog.Default(), mock, nil)
	srv := grpc.NewServer()
	defer srv.Stop()

	err := p.GRPCServer(nil, srv)
	assert.NoError(t, err)
}

func TestCustomGRPCServer(t *testing.T) {
	factory := customGRPCServer([]grpc.ServerOption{})
	srv := factory([]grpc.ServerOption{})
	assert.NotNil(t, srv)
	srv.Stop()
}

func TestNewHCPlugin(t *testing.T) {
	mock := &pluginMock{typ: "test"}
	p := newHCPlugin(hclog.Default(), mock, nil)
	assert.NotNil(t, p)
	assert.Len(t, p.servers, 1)
}

func TestHCServer_GRPCClient(t *testing.T) {
	p := &hcServer{}
	result, err := p.GRPCClient(context.Background(), nil, nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestHCDialer_DialHost_CachedConn(t *testing.T) {
	mock := &mockClientConn{}
	d := &hcDialer{conn: mock}

	conn, err := d.DialHost(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, mock, conn)
}

type mockClientConn struct {
	grpc.ClientConnInterface
}

func TestValidationUnaryInterceptor_SkipsNonProtoRequest(t *testing.T) {
	v := &mockValidator{err: errors.New("should not be called")}
	interceptor := ValidationUnaryInterceptor(v, true, false)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	// Non-proto value: validation is skipped, handler is called.
	resp, err := interceptor(context.Background(), "not-a-proto", nil, handler)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestValidationUnaryInterceptor_HandlerError(t *testing.T) {
	v := &mockValidator{}
	interceptor := ValidationUnaryInterceptor(v, false, false)

	handlerErr := errors.New("handler failed")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, handlerErr
	}

	_, err := interceptor(context.Background(), "req", nil, handler)
	assert.ErrorIs(t, err, handlerErr)
}

func TestValidationUnaryInterceptor_NoValidation(t *testing.T) {
	v := &mockValidator{}
	interceptor := ValidationUnaryInterceptor(v, false, false)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	resp, err := interceptor(context.Background(), "req", nil, handler)
	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestValidationUnaryInterceptor_RequestValidationFails(t *testing.T) {
	v := &mockValidator{err: errors.New("bad request")}
	interceptor := ValidationUnaryInterceptor(v, true, false)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	_, err := interceptor(context.Background(), &initv1.InitRequest{}, nil, handler)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok, "expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestValidationUnaryInterceptor_ResponseValidationFails(t *testing.T) {
	v := &mockValidator{err: errors.New("bad response")}
	interceptor := ValidationUnaryInterceptor(v, false, true)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &initv1.InitRequest{}, nil
	}

	_, err := interceptor(context.Background(), "req", nil, handler)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok, "expected gRPC status error")
	assert.Equal(t, codes.Internal, st.Code())
}
