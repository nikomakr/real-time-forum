package utils_test

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

func init() { // As any init function is executed before any tests are run, even the main function.
	if os.Getenv("TRIGGER_CRASH") == "1" {
		log.Fatalf("simulated UUID failure — this kills the entire process")
	}
}

func TestNewUUID_FatalfCrashesEntireProcess(t *testing.T) {
	// exec.Command(...) this tells Go to prepare to launch a separate, isolated background process running a copy of this exact test binary.
	// The os.Args[0] points to the compiled test binary file that Go creates when you type go test. You see where it goes hah?
	cmd := exec.Command(os.Args[0], "-test.run=TestNewUUID_FatalfCrashesEntireProcess")
	cmd.Env = append(os.Environ(), "TRIGGER_CRASH=1") // The os.Environ()clones your current system environment variables
	err := cmd.Run()

	if err == nil {
		t.Fatal("expected the process to crash but it ran cleanly")
	}

	t.Logf("PROVEN: log.Fatalf killed the entire process — %v", err)
}
