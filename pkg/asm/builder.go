package asm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Build(macro string, debug bool) (err error) {

	// Step 1: Run `common` commands
	fmt.Println("Running common target...")
	cmd := fmt.Sprintf("as -o %s_test.o %s_test.asm && gcc -o test %s_test.o", macro, macro, macro)
	if err = runCommand(cmd, debug); err != nil {
		return fmt.Errorf("%s: %w", cmd, err)
	}

	// Step 2: Run `tests` (just running the ./test binary and cleaning up)
	fmt.Println("Running tests target...")

	if err = runCommand("./test", debug); err != nil {
		fmt.Println("Error running tests:", err)
		return fmt.Errorf("test %s: %w", macro, err)
	}

	if err = runCommand("rm -f *.o test", debug); err != nil {
		return fmt.Errorf("clean %s: %w", macro, err)
	}

	return nil
}

// Executes a shell command with optional quiet mode based on $V
func runCommand(command string, quiet bool) error {
	if quiet {
		command = strings.ReplaceAll(command, "$Q", "@")
	} else {
		command = strings.ReplaceAll(command, "$Q", "")
		fmt.Println("Running:", command)
	}

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
