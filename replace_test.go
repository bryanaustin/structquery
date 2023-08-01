package structquery

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestReplaceMapSingle(t *testing.T) {
	t.Parallel()
	focus := []string{"one", "two", "ten", "six"}

	nu, err := Replace(focus, "2", "net")
	if em := cmp.Diff(err, nil); em != "" {
		t.Fatal(em)
	}
	if em := cmp.Diff(nu, "ten"); em != "" {
		t.Error(em)
	}
	if em := cmp.Diff(len(focus), 4); em != "" {
		t.Fatal(em)
	}
	expected := []string{"one", "two", "net", "six"}
	if em := cmp.Diff(focus, expected); em != "" {
		t.Error(em)
	}
}
