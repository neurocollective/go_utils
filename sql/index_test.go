package sql

import (
	"testing"
)

func TestSQLArgSequence(t *testing.T) {

	argSeq := new(SQLArgSequence)

	id := argSeq.Next()

	if id != 1 {
		t.Fatalf("expected 1 but got: %d", id)
	}

	id = argSeq.Next()

	if id != 2 {
		t.Fatalf("expected 2 but got: %d", id)
	}

	stringId := argSeq.NextString()

	if stringId != "$3" {
		t.Fatalf("expected '$3' but got: %s", stringId)
	}

	stringId = argSeq.NextString()

	if stringId != "$4" {
		t.Fatalf("expected '$4' but got: %s", stringId)
	}

}
