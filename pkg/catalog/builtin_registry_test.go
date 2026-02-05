package catalog

import (
	"reflect"
	"testing"
)

func TestNewBuildInRegistry_Empty(t *testing.T) {
	reg := DefaultBuiltInPluginRegistry()

	plugins := reg.Get()
	if len(plugins) != 0 {
		t.Fatalf("expected empty registry, got %d plugins", len(plugins))
	}
}

func TestBuildInRegistry_RegisterAndGet(t *testing.T) {
	reg := DefaultBuiltInPluginRegistry()

	p1 := BuiltIn{Name: "plugin-1"}
	p2 := BuiltIn{Name: "plugin-2"}

	reg.Register(p1)
	reg.Register(p2)

	got := reg.Get()
	want := []BuiltIn{p1, p2}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected plugins.\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestBuildInRegistry_GetReturnsCopy(t *testing.T) {
	reg := DefaultBuiltInPluginRegistry()

	p1 := BuiltIn{Name: "plugin-1"}
	reg.Register(p1)

	plugins := reg.Get()
	plugins[0] = BuiltIn{Name: "mutated"}

	// Fetch again — internal state must be unchanged
	got := reg.Get()
	if got[0].Name != p1.Name {
		t.Fatalf("registry was mutated via Get(); expected %v, got %v", p1, got[0])
	}
}
