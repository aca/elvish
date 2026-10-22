//go:build unix

package term

import (
	"os"
	"testing"

	"github.com/creack/pty"
	"golang.org/x/sys/unix"
	"src.elv.sh/pkg/sys/eunix"
)

func TestSetupForTUIOnce(t *testing.T) {
	_, tty := setupPTY(t)

	setupForTUIOnce(tty, tty)
	// TODO: Test whether the interesting flags in the termios were indeed set.
}

func TestSetupForTUI(t *testing.T) {
	_, tty := setupPTY(t)

	_, err := setupForTUI(tty, tty)
	if err != nil {
		t.Errorf("setupForTUI returns an error")
	}
	// TODO: Test whether the interesting flags in the termios were indeed set.
	// termios, err := eunix.TermiosForFd(int(tty.Fd()))
}

func setupPTY(t *testing.T) (ptySide, ttySide *os.File) {
	t.Helper()
	ptySide, ttySide, err := pty.Open()
	if err != nil {
		t.Skip("cannot open pty")
	}
	t.Cleanup(func() {
		ptySide.Close()
		ttySide.Close()
	})
	return ptySide, ttySide
}

func TestResetSane(t *testing.T) {
	_, tty := setupPTY(t)
	fd := int(tty.Fd())

	broken, err := eunix.TermiosForFd(fd)
	if err != nil {
		t.Fatal(err)
	}
	broken.SetICanon(false)
	broken.SetEcho(false)
	broken.SetIXON(false)
	broken.Oflag &^= unix.OPOST
	broken.Cc[unix.VINTR] = 0
	broken.ApplyToFd(fd)

	ResetSane(tty)

	got, _ := eunix.TermiosForFd(fd)
	if got.Lflag&unix.ICANON == 0 || got.Lflag&unix.ECHO == 0 || got.Oflag&unix.OPOST == 0 {
		t.Errorf("icanon/echo/opost not restored")
	}
	if got.Cc[unix.VINTR] != 3 {
		t.Errorf("VINTR = %d, want 3", got.Cc[unix.VINTR])
	}
	if got.Iflag&unix.IXON != 0 {
		t.Errorf("IXON was changed")
	}
}
