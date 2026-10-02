package archive

import "testing"

func TestPhase12OwnerContract(t *testing.T) {
	if Owner() != "archive" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executePhase12Model(newPhase12Model()) {
		t.Fatalf("archive phase12 model contract failed")
	}
}
