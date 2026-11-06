package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rc4"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/containedreality/crdpa"
	"github.com/containedreality/encrypt/shared"
	"github.com/containedreality/nekoget/progresstracker"
	"github.com/containedreality/speck"
	"golang.org/x/crypto/chacha20"
)

func encrypt(f *os.File, key []byte, salt []byte, algo string, o *os.File, pt *progresstracker.ProgressTracker) error {
	streamcipher := false
	iv := make([]byte, 16)
	rand.Read(iv)

	if salt == nil {
		salt = make([]byte, 16)
	}

	var s cipher.Stream
	var b cipher.Block
	var err error

	switch algo {
	case "aes":
		b, err = aes.NewCipher(key)
	case "speck":
		b, err = speck.New128(key)
	case "crdpa":
		fmt.Fprintf(os.Stderr, "probably shouldn't be using %s unless you know what you're doing\n", algo)
		b, err = crdpa.New(key)
	case "arcfour":
		streamcipher = true
		s, err = rc4.NewCipher(key)
	case "chacha20":
		streamcipher = true
		s, err = chacha20.NewUnauthenticatedCipher(key, iv[:12])
	default:
		return errors.New("unsupported cipher")
	}

	if !streamcipher {
		s = cipher.NewCTR(b, iv)
	}

	if err != nil {
		return err
	}

	var sw cipher.StreamWriter

	if pt.Total > 0 {
		pt.Underlying = o
		pt.Total += 32

		pt.Write(salt)
		pt.Write(iv)

		sw = cipher.StreamWriter{S: s, W: pt}

	} else {
		o.Write(salt)
		o.Write(iv)

		sw = cipher.StreamWriter{S: s, W: o}
	}

	_, err = io.Copy(sw, f)

	return err
}

func main() {
	var algorithm string
	var input string
	var output string
	var keyfile string
	var err error
	pt := &progresstracker.ProgressTracker{Total: 0}

	f := os.Stdin
	o := os.Stdout

	var verbose bool
	var list bool

	flag.StringVar(&algorithm, "a", "aes", "encryption algorithm")
	flag.StringVar(&keyfile, "k", "", "keyfile")
	flag.StringVar(&input, "i", "", "input file")
	flag.StringVar(&output, "o", "", "output file")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&list, "l", false, "list algorithms")

	flag.Parse()

	if input != "" {
		f, err = os.Open(input)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		if verbose {
			stats, err := os.Stat(input)

			if err != nil {
				log.Fatal(err)
			}

			pt.Total = stats.Size()
		}
	}

	if output != "" {
		o, err = os.Create(output)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
	}

	if list {
		fmt.Println(`cipher: key sizes (bits)
aes (default): 128, 192, 256
speck: 128, 256
crdpa: 256
arcfour: 8 -> 2048
chacha20: 256`)

		os.Exit(0)
	}

	var salt []byte
	var key []byte

	if keyfile == "" {
		p1 := shared.GetPass("Password: ")
		p2 := shared.GetPass("Confirm Password: ")

		if bytes.Equal(p1, p2) {
			key, salt = shared.DeriveKeysPassword(string(p1), nil)
		} else {
			log.Fatal("passwords don't match")
			os.Exit(1)
		}
	} else {
		kf, err := os.Open(keyfile)
		if err != nil {
			log.Fatal(err)
		}

		key, err = io.ReadAll(kf)
		if err != nil {
			log.Fatal(err)
		}

		salt = nil

		kf.Close()
	}

	err = encrypt(f, key, salt, algorithm, o, pt)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		fmt.Println()
	}
}
