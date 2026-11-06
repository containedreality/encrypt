package shared

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func GetPass(prompt string) []byte {
	fmt.Print(prompt)
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal()
	}
	fmt.Println()

	return pwd
}
