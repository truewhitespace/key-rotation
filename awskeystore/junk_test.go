package awskeystore

import "testing"

func TestStringSlice_Contains(t *testing.T) {
	slice := StringSlice{"test", "feeling", "in", "the", "way"}
	if !slice.Contains("feeling") {
		t.Errorf("expected to contain 'feelings' but got false")
	}

	if slice.Contains("sometimes") {
		t.Errorf("claims 'sometimes' exists within the slice but has no such element")
	}
}
