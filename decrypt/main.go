package main

import (
	"crypto/aes"
	"crypto/cipher"
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

func decrypt(f *os.File, key []byte, algo string, o *os.File, pt *progresstracker.ProgressTracker) error {
	streamcipher := false

	salt := make([]byte, 16)
	iv := make([]byte, 16)

	f.Read(salt)
	f.Read(iv)

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

	sr := cipher.StreamReader{S: s, R: f}
	var writer io.Writer

	if pt.Total > 0 {
		pt.Underlying = o
		writer = pt
	} else {
		writer = o
	}

	_, err = io.Copy(writer, sr)

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
	flag.BoolVar(&verbose, "v", false, "unimplemented")
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

			pt.Total = stats.Size() - 32
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

	var key []byte

	if keyfile == "" {
		p1 := shared.GetPass("Password: ")

		salt := make([]byte, 16)
		f.Read(salt)

		key, _ = shared.DeriveKeysPassword(string(p1), salt)

		f.Seek(0, io.SeekStart)
	} else {
		kf, err := os.Open(keyfile)
		if err != nil {
			log.Fatal(err)
		}

		key, err = io.ReadAll(kf)
		if err != nil {
			log.Fatal(err)
		}

		kf.Close()
	}

	err = decrypt(f, key, algorithm, o, pt)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		fmt.Println()
	}
}
