# Go Pageant client

This repository contains a library for Go that provides a native
[PuTTY][putty] Pageant SSH agent implementation compatible with the
[golang.org/x/crypto/ssh/agent][go-ssh-agent] package.

This package, rather unsuprisingly, only works with Windows.
See below for alternatives on Unix/Linux platforms. 

[putty]: https://www.chiark.greenend.org.uk/~sgtatham/
[go-ssh-agent]: https://godoc.org/golang.org/x/crypto/ssh/agent

## Usage

```golang
import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"github.com/kbolino/pageant"
)

func main() {
	agentConn, err := pageant.NewConn()
	if err != nil {
		// failed to connect to Pageant
	}
	defer agentConn.Close()
	sshAgent := agent.NewClient(agentConn)
	signers, err := sshAgent.Signers()
	if err != nil {
		// failed to get signers from Pageant
	}
	config := ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signers...)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            "somebody",
	}
	sshConn, err := ssh.Dial("tcp", "someserver:22", &config)
	if err != nil {
		// failed to connect to SSH
	}
	defer sshConn.Close()
	// now connected to SSH with public key auth from Pageant
	// ...
}
```

## Unix/Linux Alternatives

The `ssh-agent` command implements the same [SSH agent protocol][ssh-agent]
as Pageant, but over a `Unix domain socket` instead of shared memory.
The path to this socket is exposed through the environment variable
`SSH_AUTH_SOCK`.

Replace the connection to Pageant with one to the socket:
```golang
	// instead of this:
	agentConn, err := pageant.NewConn()
	// do this:
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
```

## OpenSSH for Windows Alternatives

The `ssh-add`, `ssh` commands of `OpenSSH for Windows` implements the same [SSH agent protocol][ssh-agent]
as Unix/Linux, but over a `Named Pipe` instead of `Unix domain socket`.<br>
The path to this pipe is exposed through the environment variable `SSH_AUTH_SOCK` like `\\.\pipe\somepath`<br>
The `ssh-agent` daemon of `OpenSSH for Windows` used `Named Pipe` `\\.\pipe\openssh-ssh-agent`<br>
The `sshd` daemon of `OpenSSH for Windows` used `Unix domain socket` like `/tmp/somepath` for some versions of Windows it works

Replace the connection to Pageant with one to the socket:
```golang
	// instead of this:
	agentConn, err := pageant.NewConn()
	// do this:
	agentConn, err := winio.DialPipe(os.Getenv("SSH_AUTH_SOCK"), nil)
```


[ssh-agent]: https://tools.ietf.org/html/draft-miller-ssh-agent-02

## Testing

The standard tests require Pageant to be running and to have at least 1
key loaded.
To test connecting to an SSH server, set the `sshtest` build flag and
see the comments in `pageant_ssh_test.go` for how to set up the test. 