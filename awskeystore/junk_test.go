package awskeystore

import "testing"

func (l StringSlice) assertContains(t *testing.T, expected string) {
	if !l.Contains(expected) {
		t.Helper()
		t.Errorf("Expected %+q to contain %q but did not", l, expected)
	}
}

func (l StringSlice) assertMissing(t *testing.T, expectedMissing string) {
	if l.Contains(expectedMissing) {
		t.Helper()
		t.Errorf("Expected %+q not contain %q but did", l, expectedMissing)
	}
}

func TestStringSlice_Contains(t *testing.T) {
	slice := StringSlice{"test", "feeling", "in", "the", "way"}

	slice.assertContains(t, "feeling")
	slice.assertMissing(t, "sometimes")
}
