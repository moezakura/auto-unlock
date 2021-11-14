package soundmeter

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
)

type SoundMeter struct {
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewSoundMeter() *SoundMeter {
	return &SoundMeter{}
}

func (s *SoundMeter) Exec() error {
	cmd := exec.Command("/usr/local/bin/soundmeter", "--segment", "0.2")
	var err error

	s.stdout, _ = cmd.StdoutPipe()
	s.stderr, _ = cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		o, _ := io.ReadAll(s.stdout)
		e, _ := io.ReadAll(s.stderr)

		return xerrors.Errorf("failed to execute command:\n\tstdout: %s\n\tstderr: %s", string(o), string(e))
	}
	return nil
}

func (s *SoundMeter) GetLine() string {
	if s.stdout == nil {
		return ""
	}

	scanner := bufio.NewScanner(s.stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		return strings.Trim(m, "\r\n\t ")
	}

	return ""
}
