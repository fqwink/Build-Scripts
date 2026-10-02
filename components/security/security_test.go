package security

import "testing"

func TestPhase12OwnerContract(t *testing.T) {
	if Owner() != "security" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executePhase12Model(newPhase12Model()) {
		t.Fatalf("security phase12 model contract failed")
	}
}
