package bootstrap

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"

	"github.com/openkcm/plugin-sdk/api"
	initv1 "github.com/openkcm/plugin-sdk/internal/proto/service/init/v1"
)

type needsHostServiceMock struct {
	fail bool
}

func (hsm *needsHostServiceMock) BrokerHostServices(sb api.ServiceBroker) error {
	if hsm.fail {
		return errors.New("failed")
	}
	return nil
}

func (hsm *needsHostServiceMock) SetLogger(logger hclog.Logger) {}

func (hsm *needsHostServiceMock) Close() error {
	if hsm.fail {
		return errors.New("failed")
	}
	return nil
}

func TestInitServiceInit(t *testing.T) {
	// Arrange
	mock := &needsHostServiceMock{}

	// create test cases
	tests := []struct {
		name      string
		svc       *initService
		wantError bool
	}{
		{
			name: "dialer fails",
			svc: &initService{
				logger: hclog.Default(),
				dialer: &hostDialerMock{fail: true},
				impls:  []any{&needsHostServiceMock{}},
			},
			wantError: true,
		}, {
			name: "broker host service fails",
			svc: &initService{
				logger: hclog.Default(),
				dialer: &hostDialerMock{},
				impls:  []any{&needsHostServiceMock{fail: true}},
			},
			wantError: true,
		}, {
			name: "success",
			svc: &initService{
				logger: hclog.Default(),
				dialer: &hostDialerMock{},
				impls:  []any{mock, mock},
			},
			wantError: false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			_, err := tc.svc.Init(context.Background(), &initv1.InitRequest{})

			// Assert
			if tc.wantError && err != nil { // expected error and got it
				return
			} else if tc.wantError && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err == nil {
				return
			}
		})
	}
}

func TestIinitServiceDeinit(t *testing.T) {
	// Arrange
	mock := &needsHostServiceMock{}

	// create test cases
	tests := []struct {
		name      string
		svc       *initService
		wantError bool
	}{
		{
			name: "close fails",
			svc: &initService{
				logger: hclog.Default(),
				impls:  []any{&needsHostServiceMock{fail: true}},
			},
			wantError: false,
		}, {
			name: "success",
			svc: &initService{
				logger: hclog.Default(),
				impls:  []any{mock, mock},
			},
			wantError: false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			_, err := tc.svc.Deinit(context.Background(), &initv1.DeinitRequest{})

			// Assert
			if tc.wantError && err != nil { // expected error and got it
				return
			} else if tc.wantError && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err == nil {
				return
			}
		})
	}
}

type hostDialerMock struct {
	fail bool
}

func (hdm hostDialerMock) DialHost(ctx context.Context) (grpc.ClientConnInterface, error) {
	if hdm.fail {
		return nil, errors.New("failed")
	}
	return nil, nil
}

type serviceServerMock struct {
	name string
}

func (ssm *serviceServerMock) GRPCServiceName() string {
	return ssm.name
}

func (ssm *serviceServerMock) RegisterServer(s *grpc.Server) any {
	return nil
}

func TestRegister(t *testing.T) {
	// create test cases
	tests := []struct {
		name string
		svcs []api.ServiceServer
	}{
		{
			name: "zero values",
		}, {
			name: "with services",
			svcs: []api.ServiceServer{
				&serviceServerMock{},
				&serviceServerMock{},
			},
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			grpcSrv := grpc.NewServer()

			// Act
			Register(grpcSrv, tc.svcs, hclog.Default(), hostDialerMock{})
		})
	}
}

type serviceClientMock struct{}

func (sc *serviceClientMock) GRPCServiceName() string {
	return "foo"
}
func (sc *serviceClientMock) InitClient(conn grpc.ClientConnInterface) any {
	return nil
}

func TestBrokerClient(t *testing.T) {
	// Arrange
	scm := &serviceClientMock{}

	// create test cases
	tests := []struct {
		name             string
		hostServiceNames []string
		want             bool
	}{
		{
			name: "zero values",
		}, {
			name:             "no match",
			hostServiceNames: []string{"bar"},
		}, {
			name:             "match",
			hostServiceNames: []string{"foo"},
			want:             true,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			sb := serviceBroker{
				hostServiceNames: tc.hostServiceNames,
			}

			// Act
			got := sb.BrokerClient(scm)

			// Assert
			if got != tc.want {
				t.Errorf("BrokerClient() = %v, want %v", got, tc.want)
			}
		})
	}
}
