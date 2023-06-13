package dna_processor

import "testing"

func TestRopeStorage(t *testing.T) {
	var storage DnaStorage
	storage = NewDnaStorageRope("ICFP")

	if storage.GetChar() != 'I' {
		t.Errorf("expected I")
	}

	if storage.GetChar() != 'C' {
		t.Errorf("expected C")
	}

	if storage.GetChar() != 'F' {
		t.Errorf("expected C")
	}

	storage.UndoGet()
	storage.UndoGet()

	if storage.GetChar() != 'C' {
		t.Errorf("expected C")
	}

	if storage.GetChar() != 'F' {
		t.Errorf("expected C")
	}

	if storage.GetChar() != 'P' {
		t.Errorf("expected C")
	}
}
