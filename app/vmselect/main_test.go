package vmselect

import (
	"fmt"
	"testing"
)

func TestLastWrappedErr(t *testing.T) {
	f := func(err, exp error) {
		t.Helper()
		got := lastWrappedErr(err)
		if got == nil && exp == nil {
			return
		}
		if exp.Error() != got.Error() {
			t.Fatalf("expected: \n%s\n got: \n%s", exp, got)
		}
	}

	f(nil, nil)
	f(fmt.Errorf("foo"), fmt.Errorf("foo"))

	err1 := fmt.Errorf("baz")
	f(fmt.Errorf("bar: %w", err1), err1)

	err2 := fmt.Errorf("bar: %w", err1)
	f(fmt.Errorf("foo: %w", err2), err1)
}
