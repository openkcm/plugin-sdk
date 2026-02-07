package catalog

import (
	"reflect"
	"testing"
)

func TestBuiltInRegistry(t *testing.T) {
	t.Parallel()

	t.Run("new registry is empty", func(t *testing.T) {
		t.Parallel()

		reg := DefaultBuiltInPluginRegistry()

		if got := len(reg.Get()); got != 0 {
			t.Fatalf("expected empty registry, got %d plugins", got)
		}
	})

	t.Run("register and get preserves order", func(t *testing.T) {
		t.Parallel()

		reg := DefaultBuiltInPluginRegistry()

		p1 := BuiltInPlugin{Name: "plugin-1"}
		p2 := BuiltInPlugin{Name: "plugin-2"}

		reg.Register(p1)
		reg.Register(p2)

		got := reg.Get()
		want := []BuiltInPlugin{p1, p2}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected plugins\nwant: %#v\ngot:  %#v", want, got)
		}
	})

	t.Run("get returns a copy", func(t *testing.T) {
		t.Parallel()

		reg := DefaultBuiltInPluginRegistry()

		p1 := BuiltInPlugin{Name: "plugin-1"}
		reg.Register(p1)

		plugins := reg.Get()
		plugins[0] = BuiltInPlugin{Name: "mutated"}

		got := reg.Get()
		if got[0].Name != p1.Name {
			t.Fatalf(
				"registry was mutated via Get(); want %+v, got %+v",
				p1,
				got[0],
			)
		}
	})
}
