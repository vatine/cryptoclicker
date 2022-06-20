package mixer

import (
	"math/bits"
)

var k = []uint32{
	0x428a2f98,
	0x71374491,
	0xb5c0fbcf,
	0xe9b5dba5,
	0x3956c25b,
	0x59f111f1,
	0x923f82a4,
	0xab1c5ed5,
	0xd807aa98,
	0x12835b01,
	0x243185be,
	0x550c7dc3,
	0x72be5d74,
	0x80deb1fe,
	0x9bdc06a7,
	0xc19bf174,
	0xe49b69c1,
	0xefbe4786,
	0x0fc19dc6,
	0x240ca1cc,
	0x2de92c6f,
	0x4a7484aa,
	0x5cb0a9dc,
	0x76f988da,
	0x983e5152,
	0xa831c66d,
	0xb00327c8,
	0xbf597fc7,
	0xc6e00bf3,
	0xd5a79147,
	0x06ca6351,
	0x14292967,
	0x27b70a85,
	0x2e1b2138,
	0x4d2c6dfc,
	0x53380d13,
	0x650a7354,
	0x766a0abb,
	0x81c2c92e,
	0x92722c85,
	0xa2bfe8a1,
	0xa81a664b,
	0xc24b8b70,
	0xc76c51a3,
	0xd192e819,
	0xd6990624,
	0xf40e3585,
	0x106aa070,
	0x19a4c116,
	0x1e376c08,
	0x2748774c,
	0x34b0bcb5,
	0x391c0cb3,
	0x4ed8aa4a,
	0x5b9cca4f,
	0x682e6ff3,
	0x748f82ee,
	0x78a5636f,
	0x84c87814,
	0x8cc70208,
	0x90befffa,
	0xa4506ceb,
	0xbef9a3f7,
	0xc67178f2,
}

var hash = [8]uint32{
	0x6a09e667,
	0xbb67ae85,
	0x3c6ef372,
	0xa54ff53a,
	0x510e527f,
	0x9b05688c,
	0x1f83d9ab,
	0x5be0cd19,
}

// The Mixer package is essentially an attempt at implementing SHA-256
// from scratch. However, no promises are made to implementation
// quality, speed, or correctness.

// This is basically our hash function, but as we're not intending to
// do any padding, but want a "first and final" block to operate on,
// we will just stuff this as-is into the MD rounds, after re-chunking
// to 32-bit values.
func Mix(in [64]byte) [32]byte {
	var internal [64]uint32

	for ix := 0; ix < 16; ix++ {
		j := ix * 4
		internal[ix] = uint32(in[j])<<24 | uint32(in[j+1])<<16 | uint32(in[j+2])<<8 | uint32(in[j+3])
	}

	// Extension
	for ix := 16; ix < 64; ix++ {
		v1 := internal[ix-2]
		t1 := bits.RotateLeft32(v1, -17) ^ bits.RotateLeft32(v1, -19) ^ (v1 >> 10)
		v2 := internal[ix-15]
		t2 := bits.RotateLeft32(v2, -7) ^ bits.RotateLeft32(v2, -18) ^ (v2 >> 3)
		internal[ix] = t1 + internal[ix-7] + t2 + internal[ix-16]
	}

	// Compression
	a, b, c, d, e, f, g, h := hash[0], hash[1], hash[2], hash[3], hash[4], hash[5], hash[6], hash[7]

	for ix := 0; ix < 64; ix++ {
		t1 := h + ((bits.RotateLeft32(e, -6)) ^ (bits.RotateLeft32(e, -11)) ^ (bits.RotateLeft32(e, -25))) + ((e & f) ^ (^e & g)) + k[ix] + internal[ix]

		t2 := ((bits.RotateLeft32(a, -2)) ^ (bits.RotateLeft32(a, -13)) ^ (bits.RotateLeft32(a, -22))) + ((a & b) ^ (a & c) ^ (b & c))

		h = g
		g = f
		f = e
		e = d + t1
		d = c
		c = b
		b = a
		a = t1 + t2

	}

	var tmp [8]uint32
	tmp[0] = hash[0] + a
	tmp[1] = hash[1] + b
	tmp[2] = hash[2] + c
	tmp[3] = hash[3] + d
	tmp[4] = hash[4] + e
	tmp[5] = hash[5] + f
	tmp[6] = hash[6] + g
	tmp[7] = hash[7] + h

	var out [32]uint8

	for ix, v := range tmp {
		outIx := ix * 4
		out[outIx+0] = uint8((v & 0xff000000) >> 24)
		out[outIx+1] = uint8((v & 0x00ff0000) >> 16)
		out[outIx+2] = uint8((v & 0x0000ff00) >> 8)
		out[outIx+3] = uint8((v & 0x000000ff))
	}

	return out
}
