package main

import (
	"context"
	"testing"
)

func TestTest(t *testing.T) {
	tp := TestPlugin{}
	res, err := tp.Test(context.Background(), nil)
	if err != nil {
		t.Errorf("Test() error = %v, want nil", err)
	}
	if res.Response != "test" {
		t.Errorf("Test() = %v, want %v", res.Response, "test")
	}
}

func TestConfigure(t *testing.T) {
	tp := TestPlugin{}
	res, err := tp.Configure(context.Background(), nil)
	if err != nil {
		t.Errorf("Configure() error = %v, want nil", err)
	}
	if res == nil {
		t.Errorf("Configure() = nil, want non-nil")
	}
}
