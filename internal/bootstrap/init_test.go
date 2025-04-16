package bootstrap

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	initv1 "github.com/openkcm/plugin-sdk/internal/proto/service/init/v1"
)

type bootstrapServerMockOK struct {
	initv1.UnimplementedBootstrapServer
}

func (s *bootstrapServerMockOK) Init(ctx context.Context, req *initv1.InitRequest) (*initv1.InitResponse, error) {
	return &initv1.InitResponse{PluginServiceNames: []string{}}, nil
}

func (s *bootstrapServerMockOK) Deinit(ctx context.Context, req *initv1.DeinitRequest) (*initv1.DeinitResponse, error) {
	return &initv1.DeinitResponse{}, nil
}

type bootstrapServerMockUnimplemented struct {
	initv1.UnimplementedBootstrapServer
}

func (s *bootstrapServerMockUnimplemented) Init(ctx context.Context, req *initv1.InitRequest) (*initv1.InitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method not implemented")
}

func (s *bootstrapServerMockUnimplemented) Deinit(ctx context.Context, req *initv1.DeinitRequest) (*initv1.DeinitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method not implemented")
}

type bootstrapServerMockFailing struct {
	initv1.UnimplementedBootstrapServer
}

func (s *bootstrapServerMockFailing) Init(ctx context.Context, req *initv1.InitRequest) (*initv1.InitResponse, error) {
	return nil, errors.New("failed")
}

func (s *bootstrapServerMockFailing) Deinit(ctx context.Context, req *initv1.DeinitRequest) (*initv1.DeinitResponse, error) {
	return nil, errors.New("failed")
}

func TestInit(t *testing.T) {
	// Arrange
	const bufSize = 1024 * 1024

	// create test cases
	tests := []struct {
		name      string
		srv       initv1.BootstrapServer
		wantError bool
	}{
		{
			name: "OK mock",
			srv:  &bootstrapServerMockOK{},
		}, {
			name: "Unimplemented mock",
			srv:  &bootstrapServerMockUnimplemented{},
		}, {
			name:      "Failing mock",
			srv:       &bootstrapServerMockFailing{},
			wantError: true,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			lis := bufconn.Listen(bufSize)
			s := grpc.NewServer()
			initv1.RegisterBootstrapServer(s, tc.srv)
			go func() {
				if err := s.Serve(lis); err != nil {
					log.Fatalf("Server exited with error: %v", err)
				}
			}()
			ctx := context.Background()
			conn, err := grpc.NewClient("passthrough://bufnet",
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			// Act
			_, err = Init(ctx, conn, []string{})

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

func TestDeinit(t *testing.T) {
	// Arrange
	const bufSize = 1024 * 1024

	// create test cases
	tests := []struct {
		name      string
		srv       initv1.BootstrapServer
		wantError bool
	}{
		{
			name: "OK mock",
			srv:  &bootstrapServerMockOK{},
		}, {
			name: "Unimplemented mock",
			srv:  &bootstrapServerMockUnimplemented{},
		}, {
			name:      "Failing mock",
			srv:       &bootstrapServerMockFailing{},
			wantError: true,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			lis := bufconn.Listen(bufSize)
			s := grpc.NewServer()
			initv1.RegisterBootstrapServer(s, tc.srv)
			go func() {
				if err := s.Serve(lis); err != nil {
					log.Fatalf("Server exited with error: %v", err)
				}
			}()
			ctx := context.Background()
			conn, err := grpc.NewClient("passthrough://bufnet",
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf("Failed to dial bufnet: %v", err)
			}
			defer conn.Close()

			// Act
			err = Deinit(ctx, conn)

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
