package zros_go

import "testing"

func TestNewZmqPubStub(t *testing.T) {
	_, addr := NewZmqPubStub("test")
	if addr == "" {
		t.Error("NewZmqPubStub failed")
	} else {
		t.Logf("bind address on %s", addr)
	}
}