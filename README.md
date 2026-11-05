# encrypt

this is a unfaithful rewrite of [encrypt(1)](https://www.illumos.org/man/1/encrypt), [decrypt(1)](https://www.illumos.org/man/1/decrypt) and [digest(1)](https://www.illumos.org/man/1/digest) from illumos.

it's supposed to be similar to the original tools, but you should expect major incompatibilities and major differences.

## Differences

### digest

There are the new algorithms:
* [blake2s](https://pkg.go.dev/golang.org/x/crypto@v0.58.0/blake2s)
* [blake2b](https://pkg.go.dev/golang.org/x/crypto@v0.58.0/blake2b)
* [sha3-256](https://pkg.go.dev/crypto/sha3)
* [sha3-384](https://pkg.go.dev/crypto/sha3)
* [sha3-512](https://pkg.go.dev/crypto/sha3)

digest defaults to blake2s, rather than erroring out because an algorithm wasn't choosen.
