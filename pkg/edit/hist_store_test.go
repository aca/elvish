package edit

import (
	"reflect"
	"testing"

	"src.elv.sh/pkg/store"
	"src.elv.sh/pkg/store/storedefs"
)

func TestHistStore_SuggestCandidates(t *testing.T) {
	st := store.MustTempStore(t)
	for _, c := range []string{"a", "b", "a", "c"} {
		st.AddCmd(c)
	}
	hs, err := newHistStore(st)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := hs.SuggestCandidates(), []string{"b", "a", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	hs.AddCmd(storedefs.Cmd{Text: "b"})
	if got, want := hs.SuggestCandidates(), []string{"b", "a", "c", "b"}; !reflect.DeepEqual(got, want) {
		t.Errorf("after AddCmd: got %q, want %q", got, want)
	}
	if got, want := defaultHistoryProvider(hs, "b"), ""; got != want {
		t.Errorf("exact match: got %q, want %q", got, want)
	}

	st.AddCmd("bx")
	hs.FastForward()
	if got, want := defaultHistoryProvider(hs, "b"), "x"; got != want {
		t.Errorf("after FastForward: got %q, want %q", got, want)
	}
}
