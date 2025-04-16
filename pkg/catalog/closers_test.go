package catalog

import (
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestCloserGroupClose(t *testing.T) {
	// create test cases
	tests := []struct {
		name      string
		closers   closerGroup
		wantError bool
	}{
		{
			name: "zero values",
		}, {
			name:    "one successful closer",
			closers: closerGroup{closerFunc(func() {})},
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := tc.closers.Close()

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

func TestGroupCloserFuncs(t *testing.T) {
	// create test cases
	tests := []struct {
		name    string
		closers []func()
	}{
		{
			name: "zero values",
		}, {
			name: "two closers",
			closers: []func(){
				func() {},
				func() {},
			},
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := groupCloserFuncs(tc.closers...)

			// Assert
			if len(got) != len(tc.closers) {
				t.Errorf("expected length: %d, got: %d", len(tc.closers), len(got))
			}
		})
	}
}

func TestCloserFuncClose(t *testing.T) {
	// create test cases
	tests := []struct {
		name      string
		fn        closerFunc
		wantError bool
	}{
		{
			name: "one successful closer",
			fn:   closerFunc(func() {}),
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := tc.fn.Close()

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

func TestGracefulStopWithTimeout(t *testing.T) {
	// create test cases
	tests := []struct {
		name    string
		timeout time.Duration
		want    bool
	}{
		{
			name:    "long timeout",
			timeout: 5 * time.Second,
			want:    true,
		}, {
			name:    "short timeout",
			timeout: time.Nanosecond,
			want:    false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := gracefulStopWithTimeout(grpc.NewServer(), tc.timeout)

			// Assert
			if got != tc.want {
				t.Errorf("expected value: %v, got: %v", tc.want, got)
			}
		})
	}
}
