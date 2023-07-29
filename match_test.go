package structquery

import (
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-multierror"
	"testing"
	"reflect"
)

type Bravo struct {
	Z int
	X int
}

type Alpha struct {
	A Bravo
	B Bravo
}

func TestBasicMapSliceMatch(t *testing.T) {
	t.Parallel()
	focus := map[string][]string{"node":{"one", "two"}}
	ms, errs := Match(focus, "node.1")
	checkErrs(t, errs)
	expected := []reflect.Value{reflect.ValueOf(focus["node"][1])}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestBasicMapMapMatch(t *testing.T) {
	t.Parallel()
	focus := map[string]map[string]string{
		"alpha":{"one":"111","two":"222"},
		"bravo":{"nine":"999","five":"555"},
	}
	ms, errs := Match(focus, "alpha.two")
	checkErrs(t, errs)
	expected := []reflect.Value{reflect.ValueOf(focus["alpha"]["two"])}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestBasicStructStructMatch(t *testing.T) {
	t.Parallel()
	focus := Alpha{A:Bravo{456,123},B:Bravo{9999,20}}
	ms, errs := Match(focus, "*.X")
	checkErrs(t, errs)
	expected := []reflect.Value{
		reflect.ValueOf(focus.A.X),
		reflect.ValueOf(focus.B.X)}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestBasicSliceSliceMatch(t *testing.T) {
	t.Parallel()
	focus := [][]int{{9,8,7},{5}}
	ms, errs := Match(focus, "0.2")
	checkErrs(t, errs)
	expected := []reflect.Value{reflect.ValueOf(focus[0][2])}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMultiSliceSliceMatch(t *testing.T) {
	t.Parallel()
	focus := [][]int{{22,33,44,555,666},{99,44,1},{11,55,77,88}}
	ms, errs := Match(focus, "1.*")
	checkErrs(t, errs)
	expected := []reflect.Value{
		reflect.ValueOf(focus[1][0]),
		reflect.ValueOf(focus[1][1]),
		reflect.ValueOf(focus[1][2])}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMultiSliceSliceMatchTwo(t *testing.T) {
	t.Parallel()
	focus := [][]int{{6},{99,44,1}}
	ms, errs := Match(focus, "*.0")
	checkErrs(t, errs)
	expected := []reflect.Value{
		reflect.ValueOf(focus[0][0]),
		reflect.ValueOf(focus[1][0])}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMultiArrayArrayMatch(t *testing.T) {
	t.Parallel()
	focus := [2][3]int{{6},{99,44,1}}
	ms, errs := Match(focus, "*.1")
	checkErrs(t, errs)
	expected := []reflect.Value{
		reflect.ValueOf(focus[0][1]),
		reflect.ValueOf(focus[1][1])}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func checkErrs(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if merr, ok := err.(*multierror.Error); ok {
		for _, err := range merr.Errors {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
	if t.Failed() {
		t.FailNow()
	}
}
