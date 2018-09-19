package pageant

import (
	"testing"

	"golang.org/x/crypto/ssh/agent"
)

// Pageant must be running for this test to work.
func TestNewConn(t *testing.T) {
	conn, err := NewConn()
	if err != nil {
		t.Fatalf("error on NewConn: %s", err)
	} else if conn == nil {
		t.Fatalf("NewConn returned nil")
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("error on Conn.Close: %s", err)
	}
}

// Pageant must be running and have at least 1 key loaded for this test to work.
func TestSSHAgentList(t *testing.T) {
	conn, err := NewConn()
	if err != nil {
		t.Fatalf("error on NewConn: %s", err)
	}
	defer conn.Close()
	sshAgent := agent.NewClient(conn)
	keys, err := sshAgent.List()
	if err != nil {
		t.Fatalf("error on agent.List: %s", err)
	}
	if len(keys) == 0 {
		t.Fatalf("no keys listed by Pagent")
	}
	for i, key := range keys {
		t.Logf("key %d: %s", i, key.Comment)
	}
}
