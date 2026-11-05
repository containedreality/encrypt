# encrypt

this is a rewrite of [encrypt(1)](https://www.illumos.org/man/1/encrypt), [decrypt(1)](https://www.illumos.org/man/1/decrypt) and [digest(1)](https://www.illumos.org/man/1/digest) from illumos.

## Differences

### digest

There are the new algorithms:
* blake2s
* blake2b
* sha3-256
* sha3-384
* sha3-512

digest defaults to blake2s, rather than erroring out because an algorithm wasn't choosen.