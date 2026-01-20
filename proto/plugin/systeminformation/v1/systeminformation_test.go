package systeminformationv1

import "testing"

func TestGetRequestTypeString(t *testing.T) {
	req := &GetRequest{
		Id:   "id-1",
		Type: "SYSTEM",
	}

	if got := req.GetType(); got != "SYSTEM" {
		t.Fatalf("expected type SYSTEM, got %q", got)
	}
}
