package version

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	// Test Set function to update build information
	Set("1.0", "2023-09-01", "abcdef")

	// Check if the values are correctly updated
	if BuildVersion() != "1.0" {
		t.Errorf("Expected BuildVersion() to return '1.0', but got '%s'", BuildVersion())
	}

	if BuildDate() != "2023-09-01" {
		t.Errorf("Expected BuildDate() to return '2023-09-01', but got '%s'", BuildDate())
	}

	if BuildCommit() != "abcdef" {
		t.Errorf("Expected BuildCommit() to return 'abcdef', but got '%s'", BuildCommit())
	}
}

func TestPrintConsole(t *testing.T) {
	// Capture the output of PrintConsole
	output := captureOutput(func() {
		PrintConsole()
	})

	// Verify that the printed output matches the expected format
	expectedOutput := "Build version: 1.0\nBuild date: 2023-09-01\nBuild commit: abcdef\n"

	if output != expectedOutput {
		t.Errorf("Printed output does not match the expected output.\nExpected:\n%s\nActual:\n%s", expectedOutput, output)
	}
}

// captureOutput captures the output of a function and returns it as a string
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}
