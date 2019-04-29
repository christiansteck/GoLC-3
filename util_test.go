package main

import (
	"strconv"
	"testing"
)

func TestSubBits(t *testing.T) {
	/*
	 *  Pos:   15 14 13 12 11 10  9  8  7  6  5  4  3  2  1  0
	 *  x:      1  1  1  1  1  0  1  0  1  1  0  0  1  1  1  0
	 */

	x64, _ := strconv.ParseUint("1111101011001110", 2, 16)
	x := uint16(x64)

	tests := []struct {
		hi       uint
		lo       uint
		expected string
	}{
		{
			hi:       15,
			lo:       12,
			expected: "1111",
		},
		{
			hi:       9,
			lo:       2,
			expected: "10110011",
		},
		{
			hi:       11,
			lo:       5,
			expected: "1010110",
		},
		{
			hi:       7,
			lo:       7,
			expected: "1",
		},
	}

	for _, testData := range tests {
		res := subBits(x, testData.hi, testData.lo)
		expected64, _ := strconv.ParseUint(testData.expected, 2, 16)
		expected := uint16(expected64)

		if res != expected {
			t.Errorf("Expected %b but got %b", expected, res)
		}
	}
}

func TestSignExtend(t *testing.T) {
	tests := []struct {
		x        string
		expected string
	}{
		{
			x:        "1",
			expected: "1111111111111111",
		},
		{
			x:        "01",
			expected: "0000000000000001",
		},
		{
			x:        "1010",
			expected: "1111111111111010",
		},
		{
			x:        "0100011",
			expected: "0000000000100011",
		},
		{
			x:        "110100011",
			expected: "1111111110100011",
		},
		{
			x:        "0100011",
			expected: "0000000000100011",
		},
		{
			x:        "110100011",
			expected: "1111111110100011",
		},
	}

	for _, testData := range tests {
		x64, _ := strconv.ParseUint(testData.x, 2, 16)
		expected64, _ := strconv.ParseUint(testData.expected, 2, 16)

		x, expected := uint16(x64), uint16(expected64)
		numBits := uint(len(testData.x))

		res := signExtend(x, numBits)
		if res != expected {
			t.Errorf("Expected %b but got %b", expected, res)
		}
	}
}
