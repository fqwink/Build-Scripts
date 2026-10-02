package statefile

import "testing"

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "statefile" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("statefile owner file contract failed")
	}
}
