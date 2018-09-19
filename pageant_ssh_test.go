// +build sshtest

package pageant

import (
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// This test requires all of the following to work:
//    - build tag sshtest is active
//    - environment variable PAGEANT_TEST_SSH_ADDR is set to a valid SSH
//      server address (host:port)
//    - environment variable PAGEANT_TEST_SSH_USER is set to a user name
//      that the SSH server recognizes
//    - Pageant is running on the local machine
//    - Pageant has a key that is authorized for the user on the server
func TestSSHConnect(t *testing.T) {
	pageantConn, err := NewConn()
	if err != nil {
		t.Fatalf("error on NewConn: %s", err)
	}
	defer pageantConn.Close()
	sshAgent := agent.NewClient(pageantConn)
	signers, err := sshAgent.Signers()
	if err != nil {
		t.Fatalf("cannot obtain signers from SSH agent: %s", err)
	}
	sshUser := os.Getenv("PAGEANT_TEST_SSH_USER")
	config := ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signers...)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            sshUser,
	}
	sshAddr := os.Getenv("PAGEANT_TEST_SSH_ADDR")
	sshConn, err := ssh.Dial("tcp", sshAddr, &config)
	if err != nil {
		t.Fatalf("failed to connect to %s@%s due to error: %s", sshUser, sshAddr, err)
	}
	sshConn.Close()
}
