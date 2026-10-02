package statefile

import "testing"

func TestPhase12OwnerContract(t *testing.T) {
	if Owner() != "statefile" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executePhase12Model(newPhase12Model()) {
		t.Fatalf("statefile phase12 model contract failed")
	}
}
