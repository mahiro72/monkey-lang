package helper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertEqual(t *testing.T, want, got any) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}
