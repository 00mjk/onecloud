package nbd

import (
	"testing"
)

func TestGetNBDManager(t *testing.T) {
	Init()
	nbdman := GetNBDManager()
	t.Logf("Acquire nbd: %s, %v", nbdman.AcquireNbddev(), nbdman.nbdDevs)
}
