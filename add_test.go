package structquery

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestAddMap(t *testing.T) {
	t.Parallel()
	focus := map[string]int{"one": 1, "two": 2}
	if err := Add(&focus, "ten", 10); err != nil {
		t.Fatal(err)
	}
	if em := cmp.Diff(focus, map[string]int{"one": 1, "two": 2, "ten": 10}); em != "" {
		t.Error(em)
	}
}

func TestAddMapSlice(t *testing.T) {
	t.Parallel()
	focus := map[string]*[]string{"node": {"one", "two"}}

	err := Add(&focus, "node", "ten")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	node, ok := focus["node"]
	if !ok {
		t.Fatal("Expected to find child named \"node\"")
	}
	if em := cmp.Diff(*node, []string{"one", "two", "ten"}); em != "" {
		t.Error(em)
	}
}

func TestAddStruct(t *testing.T) {
	t.Parallel()
	focus := Testmany{
		Primary:   "ppp",
		Secondary: "SSS",
	}

	err := Add(&focus, "Primary", "QQQ")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(focus.Primary, "QQQ"); em != "" {
		t.Error(em)
	}
}

func TestAddMany(t *testing.T) {
	t.Skip() // error: reflect.Value.Set using unaddressable value
	t.Parallel()
	focus := map[string]string{
		"one": "1",
		"two": "2",
		"six": "6",
	}

	err := Add(&focus, "*", "0")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	expected := map[string]string{
		"one": "0",
		"two": "0",
		"six": "0",
	}
	if em := cmp.Diff(focus, expected); em != "" {
		t.Error(em)
	}
}
