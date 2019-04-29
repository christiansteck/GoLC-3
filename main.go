package main

import (
//	"fmt"
//	"os"
)

// Condition flags for conditional register
const (
	CondPOS uint16 = 1 << iota
	CondZRO
	CondNEG
)

// Starting position of pc register
const PCStart = 0x3000

// Instructions
const (
	OpBR   uint16 = iota // branch
	OpADD                // add
	OpLD                 // load
	OpST                 // store
	OpJSR                // jump register
	OpAND                // bitwise and
	OpLDR                // load register
	OpSTR                // store register
	OpRTI                // unused
	OpNOT                // bitwise not
	OpLDI                // load indirect
	OpSTI                // store indirect
	OpJMP                // jump
	OpRES                // reserved (unused)
	OpLEA                // load effective address
	OpTRAP               // execute trap
)

// Traps
const (
	TrapGETC  = 0x20 // get character from keyboard, not echoed onto the terminal
	TrapOUT   = 0x21 // output a character
	TrapPUTS  = 0x22 // output a word string
	TrapIN    = 0x23 // get character from keyboard, echoed onto the terminal
	TrapPUTSP = 0x24 // output a byte string
	TrapHALT  = 0x25 // halt the program
)

type ALU struct {
	Reg     [8]uint16 // registers
	CondReg uint16    // conditional register
	PCReg   uint16    // program counter register

	Memory [65536]uint16 // memory
}

func (a *ALU) EmulateInstruction() {
	instr := a.Memory[a.PCReg]
	a.PCReg++

	switch op := subBits(instr, 15, 12); op {
	case OpBR:
		a.handleBR(instr)
	case OpADD:
		a.handleADD(instr)
	case OpLD:
		a.handleLD(instr)
	case OpST:
		a.handleST(instr)
	case OpJSR:
		a.handleJSR(instr)
	case OpAND:
		a.handleAND(instr)
	case OpLDR:
		a.handleLDR(instr)
	case OpSTR:
		a.handleSTR(instr)
	case OpRTI:
		panic("OpRTI is unused")
	case OpNOT:
		a.handleNOT(instr)
	case OpLDI:
		a.handleLDI(instr)
	case OpSTI:
		a.handleSTI(instr)
	case OpJMP:
		a.handleJMP(instr)
	case OpRES:
		panic("OpRES is unused")
	case OpLEA:
		a.handleLEA(instr)
	case OpTRAP:
		a.handleTRAP(instr)
	}
}

func (a *ALU) handleTRAP(instr uint16) {
	switch trapVector := subBits(instr, 7, 0); trapVector {
	//TODO
	}
}

func (a *ALU) SetCC(r uint16) {
	if a.Reg[r] == 0 {
		a.CondReg = CondZRO
	} else if subBits(a.Reg[r], 15, 15) == 1 { // Right-most bit is 1 for negative numbers
		a.CondReg = CondNEG
	} else {
		a.CondReg = CondPOS
	}
}

func (a *ALU) handleBR(instr uint16) {
	pcOffset := subBits(instr, 8, 0)
	flag := subBits(instr, 11, 9)

	if flag == a.CondReg {
		a.PCReg += signExtend(pcOffset, 9)
	}
}

func (a *ALU) handleADD(instr uint16) {
	dr := subBits(instr, 11, 9)
	sr1 := subBits(instr, 8, 6)

	var s uint16
	if subBits(instr, 5, 5) == 0 {
		s = a.Reg[subBits(instr, 2, 0)]
	} else {
		imm := subBits(instr, 4, 0)
		s = signExtend(imm, 5)
	}

	a.Reg[dr] = a.Reg[sr1] + s
	a.SetCC(dr)
}

func (a *ALU) handleLD(instr uint16) {
	dr := subBits(instr, 11, 9)
	pcOffset := subBits(instr, 8, 0)

	a.Reg[dr] = a.Memory[a.PCReg+signExtend(pcOffset, 9)]
	a.SetCC(dr)
}

func (a *ALU) handleAND(instr uint16) {
	dr := subBits(instr, 11, 9)
	sr1 := subBits(instr, 8, 6)

	var s uint16
	if subBits(instr, 5, 5) == 0 {
		s = a.Reg[subBits(instr, 2, 0)]
	} else {
		imm := subBits(instr, 4, 0)
		s = signExtend(imm, 5)
	}

	a.Reg[dr] = a.Reg[sr1] & s
	a.SetCC(dr)
}

func (a *ALU) handleJSR(instr uint16) {
	a.Reg[7] = a.PCReg

	if subBits(instr, 11, 11) == 0 {
		baseR := subBits(instr, 8, 6)
		a.PCReg = a.Reg[baseR]
	} else {
		a.PCReg += signExtend(subBits(instr, 10, 0), 11)
	}
}

func (a *ALU) handleJMP(instr uint16) {
	baseR := subBits(instr, 8, 6)
	a.PCReg = a.Reg[baseR]
}

func (a *ALU) handleLDI(instr uint16) {
	dr := subBits(instr, 11, 9)
	pcOffset := subBits(instr, 8, 0)

	a.Reg[dr] = a.Memory[a.Memory[a.PCReg+signExtend(pcOffset, 9)]]
	a.SetCC(dr)
}

func (a *ALU) handleLDR(instr uint16) {
	dr := subBits(instr, 11, 9)
	baseR := subBits(instr, 8, 6)
	offset := subBits(instr, 5, 0)

	a.Reg[dr] = a.Memory[a.Reg[baseR]+signExtend(offset, 6)]
	a.SetCC(dr)
}

func (a *ALU) handleLEA(instr uint16) {
	dr := subBits(instr, 11, 9)
	pcOffset := subBits(instr, 8, 0)

	a.Reg[dr] = a.PCReg + signExtend(pcOffset, 9)
	a.SetCC(dr)
}

func (a *ALU) handleNOT(instr uint16) {
	dr := subBits(instr, 11, 9)
	sr := subBits(instr, 8, 6)

	a.Reg[dr] = ^a.Reg[sr]
	a.SetCC(dr)
}

func (a *ALU) handleST(instr uint16) {
	sr := subBits(instr, 11, 9)
	pcOffset := subBits(instr, 8, 0)

	a.Memory[a.PCReg+signExtend(pcOffset, 9)] = a.Reg[sr]
}

func (a *ALU) handleSTI(instr uint16) {
	sr := subBits(instr, 11, 9)
	pcOffset := subBits(instr, 8, 0)

	a.Memory[a.Memory[a.PCReg+signExtend(pcOffset, 9)]] = a.Reg[sr]
}

func (a *ALU) handleSTR(instr uint16) {
	sr := subBits(instr, 11, 9)
	baseR := subBits(instr, 8, 6)
	offset := subBits(instr, 5, 0)

	a.Memory[a.Reg[baseR]+signExtend(offset, 6)] = a.Reg[sr]
}

func main() {}
