package pdu

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	expect "github.com/google/goexpect"
	"github.com/ziutek/telnet"
)

type PDU struct {
	ex      expect.Expecter
	timeout time.Duration
}

func Dial(network, addr string, timeout time.Duration) (*PDU, error) {
	conn, err := telnet.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	err = conn.SetEcho(false)
	if err != nil {
		return nil, err
	}

	resCh := make(chan error)
	ex, _, err := expect.SpawnGeneric(&expect.GenOptions{
		In:  conn,
		Out: conn,
		Wait: func() error {
			return <-resCh
		},
		Close: func() error {
			close(resCh)
			return conn.Close()
		},
		Check: func() bool { return true },
	}, timeout)
	if err != nil {
		return nil, err
	}

	return &PDU{
		ex:      ex,
		timeout: timeout,
	}, nil
}

func (p *PDU) Authenticate(username, password string) error {
	res, err := p.ex.ExpectBatch([]expect.Batcher{
		&expect.BExp{R: `User Name\s*: `},
		&expect.BSnd{S: username + "\r\n"},
		&expect.BExp{R: `Password\s*: `},
		&expect.BSnd{S: password + " -c" + "\r\n"},
		&expect.BExp{R: `APC>\s*`},
	}, p.timeout)
	if err != nil {
		log.Fatalf("ExpectBatch failed: %v , res: %#v", err, res)
	}

	return nil
}

func (p *PDU) Run(args ...string) ([]string, error) {
	if len(args) < 0 {
		return nil, errors.New("invalid")
	}

	cmd := ""
	for i := 0; i < len(args); i++ {
		if i == 0 {
			cmd += args[i]
		} else {
			cmd += fmt.Sprintf(" \"%s\"", args[i])
		}
	}

	res, err := p.ex.ExpectBatch([]expect.Batcher{
		&expect.BSnd{S: cmd + "\r\n"},
		&expect.BExp{R: `\r\n\s+OK\r\n(([^\r\n]+\r\n)*)APC>\s*`},
	}, p.timeout)
	if err != nil {
		log.Fatalf("ExpectBatch failed: %v , res: %#v", err, res)
	}

	lines := res[len(res)-1].Match[1]
	var out []string
	for _, line := range strings.Split(lines, "\r\n") {
		if len(line) > 0 {
			out = append(out, strings.TrimSpace(line))
		}
	}
	return out, nil
}

func (p *PDU) Whoami() (string, error) {
	lines, err := p.Run("whoami")
	if err != nil {
		return "", err
	}

	return lines[0], nil
}

func (p *PDU) Logout() error {
	res, err := p.ex.ExpectBatch([]expect.Batcher{
		&expect.BSnd{S: "logout\r\n"},
		&expect.BExp{R: `\r\nBye.\r\n\r\nConnection Closed - Bye\r\n`},
	}, p.timeout)
	if err != nil {
		log.Fatalf("ExpectBatch failed: %v , res: %#v", err, res)
	}

	return p.ex.Close()
}
