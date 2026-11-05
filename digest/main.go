package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

func hashFile(f *os.File, algo string) ([]byte, error) {
	var h hash.Hash

	switch algo {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha384":
		// https://www.illumos.org/man/1/digest
		// the manual says sha385, but uhh i'm just gonna say this is a typo
		h = sha512.New384()
	case "sha512":
		h = sha512.New()
	// these are algorithms not supported in the normal program,
	// but still i felt like adding
	case "sha3-256":
		h = sha3.New256()
	case "sha3-384":
		h = sha3.New384()
	case "sha3-512":
		h = sha3.New512()
	case "blake2b":
		h, _ = blake2b.New512(nil)
	case "blake2s":
		h, _ = blake2s.New256(nil)
	default:
		return nil, errors.New("unsupported algorithm")
	}

	_, err := io.Copy(h, f)

	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func main() {
	var algo string
	var list bool
	var verbose bool

	flag.StringVar(&algo, "a", "blake2s", "digest algorithm to be used, see -l for a list.")
	flag.BoolVar(&list, "l", false, "display list of algorithms available")
	flag.BoolVar(&verbose, "v", false, "verbose output, includes algorithm and file name in the output")
	flag.Parse()

	// not the best way to do this but what ever, it works :p
	if list {
		fmt.Println(`md5
sha1
sha256
sha384
sha512
sha3-256
sha3-384
sha3-512
blake2b
blake2s`)

		os.Exit(0)
	}

	files := flag.Args()

	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}

		digest, err := hashFile(f, algo)
		encoded := hex.EncodeToString(digest)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			fmt.Printf("%s (%s) = %s\n", algo, fname, encoded)
		} else if len(files) > 1 {
			fmt.Printf("(%s) = %s\n", fname, encoded)
		} else {
			fmt.Println(encoded)
		}

		f.Close()
	}

	if len(files) == 0 {
		digest, err := hashFile(os.Stdin, algo)
		encoded := hex.EncodeToString(digest)

		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			fmt.Printf("%s %s\n", algo, encoded)
		} else {
			fmt.Println(encoded)
		}
	}
}
