package worker

import (
	"testing"
)

func TestCreatePool(t *testing.T) {

	p := NewDispatcher(4)

	if &p == nil {
		t.Error("Error while creating Pool")
	}
}
