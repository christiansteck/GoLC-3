package main

import (
	"os"
	"os/exec"
)

func disableInputBuffering() {
    // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func processInput(a *ALU) {
    var b []byte = make([]byte, 1)
    for {
        os.Stdin.Read(b)
		a.Memory[KBSR] = 0x8000
		a.Memory[KBDR] = uint16(b[0])

		// send to KBSRChan in non-blocking way
		select {
			case a.KBSRChan <- struct{}{}:
			default:
		}
    }
}

