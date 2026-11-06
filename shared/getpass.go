package shared

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/term"
)

func GetPass(prompt string) ([]byte, error) {
	var tty *os.File
	var err error

	if runtime.GOOS == "windows" {
		tty, err = os.OpenFile("CONIN$", os.O_RDWR, 0)
	} else {
		tty, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	}

	if err != nil {
		tty = os.Stdin
	} else {
		defer tty.Close()
	}

	fmt.Fprint(os.Stderr, prompt)

	bytes, err := term.ReadPassword(int(tty.Fd()))
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stderr)

	return bytes, nil
}
