package integer

import "testing"

func TestAdd(t *testing.T) {
	assertResult := func(t *testing.T, got, expected int) {
		t.Helper()
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	}


	sum := Add(2, 2)
	expected := 4
	assertResult(t, sum, expected)
}