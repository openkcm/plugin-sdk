package bootstrap

import (
	"testing"

	"google.golang.org/grpc"
)

type pluginMock struct {
	typ string
}

func (m *pluginMock) Type() string {
	return m.typ
}
func (m *pluginMock) GRPCServiceName() string {
	return "mock"
}
func (m *pluginMock) RegisterServer(server *grpc.Server) any {
	return nil
}
func (m *pluginMock) InitClient(conn grpc.ClientConnInterface) any {
	return nil
}

func TestServerHandshakeConfig(t *testing.T) {
	// Arrange
	mock := &pluginMock{typ: "test"}

	// Act
	got := ServerHandshakeConfig(mock)

	// Assert
	if got.ProtocolVersion != 1 {
		t.Errorf("Expected ProtocolVersion to be 1, but got %d", got.ProtocolVersion)
	}
	if got.MagicCookieKey != "test" {
		t.Errorf("Expected MagicCookieKey to be 'test', but got %s", got.MagicCookieKey)
	}
	if got.MagicCookieValue != "test" {
		t.Errorf("Expected MagicCookieValue to be 'test', but got %s", got.MagicCookieValue)
	}
}

func TestClientHandshakeConfig(t *testing.T) {
	// Arrange
	mock := &pluginMock{typ: "test"}

	// Act
	got := ClientHandshakeConfig(mock)

	// Assert
	if got.ProtocolVersion != 1 {
		t.Errorf("Expected ProtocolVersion to be 1, but got %d", got.ProtocolVersion)
	}
	if got.MagicCookieKey != "test" {
		t.Errorf("Expected MagicCookieKey to be 'test', but got %s", got.MagicCookieKey)
	}
	if got.MagicCookieValue != "test" {
		t.Errorf("Expected MagicCookieValue to be 'test', but got %s", got.MagicCookieValue)
	}
}
