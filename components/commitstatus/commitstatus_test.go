package commitstatus

import "testing"

func TestPhase12OwnerContract(t *testing.T) {
	if Owner() != "commitstatus" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executePhase12Model(newPhase12Model()) {
		t.Fatalf("commitstatus phase12 model contract failed")
	}
}
