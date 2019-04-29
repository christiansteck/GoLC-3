package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildInstr(instr string) uint16 {
	if len(instr) != 12 {
		panic("Wrong length for rest parameter")
	}
	instr64, _ := strconv.ParseUint(instr, 2, 16)
	return uint16(instr64)
}

func TestHandleBR(t *testing.T) {
	tests := []struct {
		description string
		condReg     uint16
		pcReg       uint16
		instr       uint16

		expectedPCReg uint16
	}{
		{
			description: "Set PCReg correctly for CondPOS",
			condReg:     CondPOS,
			pcReg:       PCStart,
			instr:       buildInstr("001" + "000000001"),

			expectedPCReg: 0x3001,
		},
		{
			description: "Set PCReg correctly for CondZRO",
			condReg:     CondZRO,
			pcReg:       PCStart,
			instr:       buildInstr("010" + "000000001"),

			expectedPCReg: 0x3001,
		},
		{
			description: "Set PCReg correctly for CondNEG",
			condReg:     CondNEG,
			pcReg:       PCStart,
			instr:       buildInstr("100" + "000000001"),

			expectedPCReg: 0x3001,
		},
		{
			description: "Does nothing for wrong CondReg",
			condReg:     CondNEG,
			pcReg:       PCStart,
			instr:       buildInstr("001" + "000000001"),

			expectedPCReg: 0x3000,
		},
		{
			description: "Handle sign extension correctly",
			condReg:     CondPOS,
			pcReg:       PCStart,
			instr:       buildInstr("001" + "111111111"),

			expectedPCReg: 0x2FFF,
		},
	}

	for _, testData := range tests {
		a := ALU{
			CondReg: testData.condReg,
			PCReg:   testData.pcReg,
		}

		a.handleBR(testData.instr)

		if a.PCReg != testData.expectedPCReg {
			t.Errorf("PCReg should be equal for '%s'\n"+
				"Expected %b\n"+
				"Got      %b",
				testData.description,
				testData.expectedPCReg,
				a.PCReg,
			)
		}
	}
}

func TestHandleADD(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0x7890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly adds in immediate mode, result is positive",
			instr:       buildInstr("101" + "011" + "1" + "00111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x345D, 0x6789, 0x7890},
			expectedCondReg: CondPOS,
		},
		{
			description: "Correctly adds in immediate mode, result is zero",
			instr:       buildInstr("101" + "000" + "1" + "11111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x0, 0x6789, 0x7890},
			expectedCondReg: CondZRO,
		},
		{
			description: "Correctly adds in immediate mode, result is negative",
			instr:       buildInstr("101" + "000" + "1" + "11110"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0xFFFF, 0x6789, 0x7890},
			expectedCondReg: CondNEG,
		},
		{
			description: "Correctly adds in register mode",
			instr:       buildInstr("101" + "001" + "0" + "00" + "010"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x3579, 0x6789, 0x7890},
			expectedCondReg: CondPOS,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg: initialReg,
		}

		a.handleADD(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleAND(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly ands in immediate mode, result is positive",
			instr:       buildInstr("101" + "011" + "1" + "00111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x0006, 0x6789, 0xF890},
			expectedCondReg: CondPOS,
		},
		{
			description: "Correctly ands in immediate mode, result is zero",
			instr:       buildInstr("101" + "000" + "1" + "00010"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x0, 0x6789, 0xF890},
			expectedCondReg: CondZRO,
		},
		{
			description: "Correctly ands in immediate mode, result is negative",
			instr:       buildInstr("101" + "111" + "1" + "11111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0xF890, 0x6789, 0xF890},
			expectedCondReg: CondNEG,
		},
		{
			description: "Correctly adds in register mode",
			instr:       buildInstr("101" + "001" + "0" + "00" + "010"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x0204, 0x6789, 0xF890},
			expectedCondReg: CondPOS,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg: initialReg,
		}

		a.handleAND(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleJMP(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedPCReg uint16
	}{
		{
			description: "Correctly sets PCReg",
			instr:       buildInstr("000" + "011" + "000000"),

			expectedPCReg: 0x3456,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg:   initialReg,
			PCReg: PCStart,
		}

		a.handleJMP(testData.instr)

		assert.Equal(testData.expectedPCReg, a.PCReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleJSR(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg   [8]uint16
		expectedPCReg uint16
	}{
		{
			description: "Correctly sets PCReg in immediate mode",
			instr:       buildInstr("1" + "010" + "1100" + "0000"),

			expectedReg:   [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, PCStart},
			expectedPCReg: 0x32C0,
		},
		{
			description: "Correctly sets PCReg in register mode",
			instr:       buildInstr("0" + "00" + "110" + "000000"),

			expectedReg:   [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, PCStart},
			expectedPCReg: 0x6789,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg:   initialReg,
			PCReg: PCStart,
		}

		a.handleJSR(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedPCReg, a.PCReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleLD(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly loads memory",
			instr:       buildInstr("101" + "011000000"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x7777, 0x6789, 0xF890},
			expectedCondReg: CondPOS,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg:   initialReg,
			PCReg: PCStart,
		}
		a.Memory[0x30C0] = 0x7777

		a.handleLD(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleLDI(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly loads memory",
			instr:       buildInstr("101" + "011000000"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x9000, 0x6789, 0xF890},
			expectedCondReg: CondNEG,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg:   initialReg,
			PCReg: PCStart,
		}
		a.Memory[0x30C0] = 0x7777
		a.Memory[0x7777] = 0x9000

		a.handleLDI(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleLDR(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly loads memory",
			instr:       buildInstr("101" + "011" + "000111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x7777, 0x6789, 0xF890},
			expectedCondReg: CondPOS,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg:   initialReg,
			PCReg: PCStart,
		}
		a.Memory[0x345D] = 0x7777

		a.handleLDR(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleLEA(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly sets PCReg",
			instr:       buildInstr("101" + "011000111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x30C7, 0x6789, 0xF890},
			expectedCondReg: CondPOS,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg:   initialReg,
			PCReg: PCStart,
		}

		a.handleLEA(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleNOT(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	tests := []struct {
		description string
		instr       uint16

		expectedReg     [8]uint16
		expectedCondReg uint16
	}{
		{
			description: "Correctly stores bitwise complement",
			instr:       buildInstr("101" + "011" + "111111"),

			expectedReg:     [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0xCBA9, 0x6789, 0xF890},
			expectedCondReg: CondNEG,
		},
	}

	for _, testData := range tests {
		a := ALU{
			Reg: initialReg,
		}

		a.handleNOT(testData.instr)

		assert.Equal(testData.expectedReg, a.Reg, "Should be equal for %s", testData.description)
		assert.Equal(testData.expectedCondReg, a.CondReg, "Should be equal for %s", testData.description)
	}
}

func TestHandleST(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	a := ALU{
		Reg:   initialReg,
		PCReg: PCStart,
	}

	a.handleST(buildInstr("101" + "000000101"))

	assert.Equal(initialReg, a.Reg, "Should be equal for 'HandleST'")
	assert.Equal(a.Reg[5], a.Memory[0x3005], "Should be equal for 'HandleST'")
}

func TestHandleSTI(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	a := ALU{
		Reg:   initialReg,
		PCReg: PCStart,
	}
	a.Memory[0x3005] = 0x66FF

	a.handleSTI(buildInstr("101" + "000000101"))

	assert.Equal(initialReg, a.Reg, "Should be equal for 'HandleSTI'")
	assert.Equal(a.Reg[5], a.Memory[0x66FF], "Should be equal for 'HandleSTI'")
}

func TestHandleSTR(t *testing.T) {
	assert := assert.New(t)
	initialReg := [8]uint16{0x0001, 0x1234, 0x2345, 0x3456, 0x4567, 0x5678, 0x6789, 0xF890}

	a := ALU{
		Reg:   initialReg,
		PCReg: PCStart,
	}

	a.handleSTR(buildInstr("101" + "000" + "000101"))

	assert.Equal(initialReg, a.Reg, "Should be equal for 'HandleSTR'")
	assert.Equal(a.Reg[5], a.Memory[0x0006], "Should be equal for 'HandleSTR'")
}
