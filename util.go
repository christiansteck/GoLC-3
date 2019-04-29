package main

// signExtend extends a uint16 x of numBits relevant bits with 0s if positive and 1s if negative
func signExtend(x uint16, numBits uint) uint16 {
	if ((x >> (numBits - 1)) & 1) > 0 {
		x |= (0xFFFF << numBits)
	}
	return x
}

// subBits extracts the bits between hi and lo inclusively
func subBits(x uint16, hi, lo uint) uint16 {
	return (x & (0xFFFF >> (15 - hi))) >> lo
}
