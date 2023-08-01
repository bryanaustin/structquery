package structquery

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestRemoveMapSingle(t *testing.T) {
	t.Skip()
	t.Parallel()
	focus := []string{"one", "two", "ten", "six"}

	err := Remove(focus, "2")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(focus), 3); em != "" {
		t.Fatal(em)
	}
	expected := []string{"one", "two", "six"}
	if em := cmp.Diff(focus, expected); em != "" {
		t.Error(em)
	}
}

func TestRemoveMapMulti(t *testing.T) {
	t.Parallel()
	focus := map[string]*map[string]int{
		"first": {"one": 1, "two": 2},
		"last":  {"ten": 10, "six": 6},
	}

	err := Remove(focus, "last.*")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(*focus["last"]), 0); em != "" {
		t.Error(em)
	}
	if em := cmp.Diff(len(*focus["first"]), 2); em != "" {
		t.Error(em)
	}
}

func TestRemoveStructSingle(t *testing.T) {
	focus := map[string]*Testmany{
		"zeebra":   {Primary: "uuu", Secondary: "yty"},
		"elephant": {Primary: "(O)", Secondary: "^u^"},
	}

	err := Remove(&focus, "zeebra.Primary")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(focus), 2); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(focus["zeebra"].Primary, ""); em != "" {
		t.Error(em)
	}
	if em := cmp.Diff(focus["zeebra"].Secondary, "yty"); em != "" {
		t.Error(em)
	}
}

func TestRemoveStructMulti(t *testing.T) {
	focus := map[string]*Testmany{
		"zeebra":   {Primary: "uuu", Secondary: "yty"},
		"elephant": {Primary: "(O)", Secondary: "^u^"},
	}

	err := Remove(&focus, "elephant.*")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(focus), 2); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(focus["elephant"].Primary, ""); em != "" {
		t.Error(em)
	}
	if em := cmp.Diff(focus["elephant"].Secondary, ""); em != "" {
		t.Error(em)
	}
}
