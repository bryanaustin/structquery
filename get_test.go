package structquery

import (
	"testing"
	"github.com/google/go-cmp/cmp"
)

func TestGetSingle(t *testing.T) {
	t.Parallel()

	focus := struct {
		A string
		B struct {
			C string
			D int
		}
	}{
		A: "root-ish",
		B: struct{
			C string
			D int
		}{C: "nested", D:8998},
	}

	seelist, err := Get(focus, "B.C")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(seelist), 1); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(seelist[0], "nested"); em != "" {
		t.Error(em)
	}

	deelist, err := Get(focus, "B.D")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(deelist), 1); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(deelist[0], 8998); em != "" {
		t.Error(em)
	}
}

func TestGetMapSingle(t *testing.T) {
	t.Parallel()
	focus := map[string]string{
		"one": "1",
		"two": "2",
		"six": "6",
	}

	sixlist, err := Get(focus, "six")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(sixlist), 1); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(sixlist[0], "6"); em != "" {
		t.Error(em)
	}
}

type Testmany struct {
	Primary string
	Secondary string
}

func TestGetManyStruct(t *testing.T) {
	focus := map[string]Testmany {
		"zeebra":{Primary:"uuu",Secondary:"yty"},
		"elephant":{Primary:"(O)",Secondary:"^u^"},
	}

	zeelist, err := Get(focus, "zeebra.*")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(zeelist), 2); em != "" {
		t.Fatal(em)
	}
	expected := []interface{}{"uuu","yty"}
	if em := cmp.Diff(zeelist, expected); em != "" {
		t.Error(em)
	}
}

func TestGetNumberAddress(t *testing.T) {
	t.Parallel()
	focus := []string{"one", "two", "ten", "six"}

	snglist, err := Get(focus, "2")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(snglist), 1); em != "" {
		t.Fatal(em)
	}
	expected := "ten"
	if em := cmp.Diff(snglist[0], expected); em != "" {
		t.Error(em)
	}
}

func TestGetRoot(t *testing.T) {
	t.Parallel()
	focus := 0x44
	
	v, err := Get(focus, "")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(len(v), 1); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(v[0], 0x44); em != "" {
		t.Error(em)
	}
}
