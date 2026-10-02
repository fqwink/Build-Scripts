package archive

import "testing"

func TestOwnerFileContract(t *testing.T) {
	if Owner() != "archive" {
		t.Fatalf("unexpected owner: %s", Owner())
	}
	if !executeOwnerFileContract(newOwnerFileContract()) {
		t.Fatalf("archive owner file contract failed")
	}
}
