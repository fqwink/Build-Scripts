package commitstatus

import "testing"

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "commitstatus" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("commitstatus owner file contract failed")
	}
}
