//go:build unix

// Copyright 2015 go-termios Author. All Rights Reserved.
// https://github.com/go-termios/termios
// Author: John Lenton <chipaca@github.com>

package eunix

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// Termios represents terminal attributes.
type Termios unix.Termios

// TermiosForFd returns a pointer to a Termios structure if the file
// descriptor is open on a terminal device.
func TermiosForFd(fd int) (*Termios, error) {
	term, err := unix.IoctlGetTermios(fd, getAttrIOCTL)
	return (*Termios)(term), err
}

// ApplyToFd applies term to the given file descriptor.
func (term *Termios) ApplyToFd(fd int) error {
	return unix.IoctlSetTermios(fd, setAttrNowIOCTL, (*unix.Termios)(unsafe.Pointer(term)))
}

// Copy returns a copy of term.
func (term *Termios) Copy() *Termios {
	v := *term
	return &v
}

// SetVTime sets the timeout in deciseconds for noncanonical read.
func (term *Termios) SetVTime(v uint8) {
	term.Cc[unix.VTIME] = v
}

// SetVMin sets the minimal number of characters for noncanonical read.
func (term *Termios) SetVMin(v uint8) {
	term.Cc[unix.VMIN] = v
}

// SetICanon sets the canonical flag.
func (term *Termios) SetICanon(v bool) {
	setFlag(&term.Lflag, unix.ICANON, v)
}

// SetIExten sets the iexten flag.
func (term *Termios) SetIExten(v bool) {
	setFlag(&term.Lflag, unix.IEXTEN, v)
}

// SetEcho sets the echo flag.
func (term *Termios) SetEcho(v bool) {
	setFlag(&term.Lflag, unix.ECHO, v)
}

// SetICRNL sets the CRNL iflag bit.
func (term *Termios) SetICRNL(v bool) {
	setFlag(&term.Iflag, unix.ICRNL, v)
}

// SetIXON sets the IXON iflag bit.
func (term *Termios) SetIXON(v bool) {
	setFlag(&term.Iflag, unix.IXON, v)
}

func setFlag(flag *termiosFlag, mask termiosFlag, v bool) {
	if v {
		*flag |= mask
	} else {
		*flag &= ^mask
	}
}

// SetSane resets term to a sane state for line-oriented programs, similar to
// "stty sane": cooked mode with echo, signal generation, CR/NL translation and
// default control characters. IXON, IUTF8, the character size and baud rates
// are left unchanged.
func (term *Termios) SetSane() {
	term.Iflag |= unix.BRKINT | unix.ICRNL | unix.IMAXBEL
	term.Iflag &^= unix.IGNBRK | unix.INLCR | unix.IGNCR | unix.IXOFF | unix.IXANY
	term.Oflag |= unix.OPOST | unix.ONLCR
	term.Oflag &^= unix.OCRNL | unix.ONOCR | unix.ONLRET
	term.Cflag |= unix.CREAD
	term.Lflag |= unix.ISIG | unix.ICANON | unix.IEXTEN |
		unix.ECHO | unix.ECHOE | unix.ECHOK | unix.ECHOCTL | unix.ECHOKE
	term.Lflag &^= unix.ECHONL | unix.NOFLSH | unix.TOSTOP | unix.ECHOPRT | unix.FLUSHO

	cc := &term.Cc
	cc[unix.VINTR] = 'C' & 0x1f
	cc[unix.VQUIT] = '\\' & 0x1f
	cc[unix.VERASE] = 0x7f
	cc[unix.VKILL] = 'U' & 0x1f
	cc[unix.VEOF] = 'D' & 0x1f
	cc[unix.VEOL] = vdisable
	cc[unix.VEOL2] = vdisable
	cc[unix.VSTART] = 'Q' & 0x1f
	cc[unix.VSTOP] = 'S' & 0x1f
	cc[unix.VSUSP] = 'Z' & 0x1f
	cc[unix.VREPRINT] = 'R' & 0x1f
	cc[unix.VWERASE] = 'W' & 0x1f
	cc[unix.VLNEXT] = 'V' & 0x1f
	cc[unix.VDISCARD] = 'O' & 0x1f
}
