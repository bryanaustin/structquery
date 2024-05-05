package structquery

import (
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-multierror"
	"reflect"
	"testing"
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
	focus := map[string][]string{"node": {"one", "two"}}
	ms, errs := Match(focus, []string{"node", "1"})
	checkErrs(t, errs)
	expectedchild := reflect.ValueOf(focus["node"][1])
	expected := []MatchStack{{Child: &expectedchild}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestBasicMapMapMatch(t *testing.T) {
	t.Parallel()
	focus := map[string]map[string]string{
		"alpha": {"one": "111", "two": "222"},
		"bravo": {"nine": "999", "five": "555"},
	}
	ms, errs := Match(focus, []string{"alpha", "two"})
	checkErrs(t, errs)
	expectedchild := reflect.ValueOf(focus["alpha"]["two"])
	expected := []MatchStack{{Child: &expectedchild}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestBasicStructStructMatch(t *testing.T) {
	t.Parallel()
	focus := Alpha{A: Bravo{456, 123}, B: Bravo{9999, 20}}
	ms, errs := Match(focus, []string{"*", "X"})
	checkErrs(t, errs)
	expectedchild1 := reflect.ValueOf(focus.A.X)
	expectedchild2 := reflect.ValueOf(focus.B.X)
	expected := []MatchStack{
		{Child: &expectedchild1},
		{Child: &expectedchild2}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestBasicSliceSliceMatch(t *testing.T) {
	t.Parallel()
	focus := [][]int{{9, 8, 7}, {5}}
	ms, errs := Match(focus, []string{"0", "2"})
	checkErrs(t, errs)
	expectedchild := reflect.ValueOf(focus[0][2])
	expected := []MatchStack{{Child: &expectedchild}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMultiSliceSliceMatch(t *testing.T) {
	t.Parallel()
	focus := [][]int{{22, 33, 44, 555, 666}, {99, 44, 1}, {11, 55, 77, 88}}
	ms, errs := Match(focus, []string{"1", "*"})
	checkErrs(t, errs)
	expectedchild1 := reflect.ValueOf(focus[1][0])
	expectedchild2 := reflect.ValueOf(focus[1][1])
	expectedchild3 := reflect.ValueOf(focus[1][2])
	expected := []MatchStack{
		{Child: &expectedchild1},
		{Child: &expectedchild2},
		{Child: &expectedchild3}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMultiSliceSliceMatchTwo(t *testing.T) {
	t.Parallel()
	focus := [][]int{{6}, {99, 44, 1}}
	ms, errs := Match(focus, []string{"*", "0"})
	checkErrs(t, errs)
	expectedchild1 := reflect.ValueOf(focus[0][0])
	expectedchild2 := reflect.ValueOf(focus[1][0])
	expected := []MatchStack{
		{Child: &expectedchild1},
		{Child: &expectedchild2}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMultiArrayArrayMatch(t *testing.T) {
	t.Parallel()
	focus := [2][3]int{{6}, {99, 44, 1}}
	ms, errs := Match(focus, []string{"*", "1"})
	checkErrs(t, errs)
	expectedchild1 := reflect.ValueOf(focus[0][1])
	expectedchild2 := reflect.ValueOf(focus[1][1])
	expected := []MatchStack{
		{Child: &expectedchild1},
		{Child: &expectedchild2}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

func TestMapNilChild(t *testing.T) {
	t.Skip() // Can't compare maps
	t.Parallel()
	focus := map[string]int{"alpha": 1, "barvo": 2}
	ms, errs := Match(focus, []string{"zulu"})
	checkErrs(t, errs)
	expectedparent := reflect.ValueOf(focus)
	expected := []MatchStack{{Parent: &expectedparent}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

type PS3 struct {
	Number int
}

type PS2 struct {
	Childer *PS3
}

type PS1 struct {
	Child *PS2
}

func TestPointerStruct(t *testing.T) {
	t.Parallel()
	focus := PS1{Child:&PS2{Childer:&PS3{Number:7}}}
	ms, errs := Match(focus, []string{"Child", "Childer", "Number"})
	checkErrs(t, errs)
	// I need to document how parent firld works
	//expectedparent := reflect.ValueOf(focus.Child)
	expectedchild := reflect.ValueOf(focus.Child.Childer.Number)
	//expected := []MatchStack{{Parent: &expectedparent, Child:&expectedchild}}
	expected := []MatchStack{{Child:&expectedchild}}
	if msg := cmp.Diff(ms, expected); msg != "" {
		t.Error(msg)
	}
}

//TODO: test nil pointer error

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
