package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExamples runs all example files and verifies their output
func TestExamples(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected []string // Lines that must appear in output
		skip     bool     // Skip if file requires user input or external resources
		skipMsg  string
	}{
		{
			name: "hello",
			file: "hello.pseudo",
			expected: []string{
				"Hello, World!",
				"Welcome to Cambridge Pseudocode",
			},
		},
		{
			name: "variables",
			file: "variables.pseudo",
			expected: []string{
				"Hello Alice!",
				"Age: 17",
				"Height: 1.65 meters",
				"Is a student: TRUE",
				"PI = 3.14159",
			},
		},
		{
			name: "selection",
			file: "selection.pseudo",
			expected: []string{
				"Score: 75 Grade: C",
				"Wednesday",
			},
		},
		{
			name: "loops",
			file: "loops.pseudo",
			expected: []string{
				"FOR loop:",
				"i = 1",
				"i = 5",
				"FOR loop with STEP 2:",
				"0",
				"10",
				"WHILE loop (sum until > 50):",
				"Count: 10 Sum: 55",
				"REPEAT loop (countdown):",
				"5",
				"1",
				"Liftoff!",
			},
		},
		{
			name: "arrays",
			file: "arrays.pseudo",
			expected: []string{
				"Numbers array:",
				"Numbers[1] = 10",
				"Numbers[5] = 50",
				"Sum: 150",
				"Average: 30",
				"Names array:",
				"1: Alice",
				"3: Charlie",
				"3x3 Matrix (multiplication table):",
				"1 2 3",
				"3 6 9",
			},
		},
		{
			name: "strings",
			file: "strings.pseudo",
			expected: []string{
				"Original: Hello, World!",
				"Length: 13",
				"LEFT(Text, 5): Hello",
				"RIGHT(Text, 6): World!",
				"MID(Text, 8, 5): World",
				"UCASE: HELLO, WORLD!",
				"LCASE: hello, world!",
				"Character: A",
				"ASC('A'): 65",
				"CHR(66): B",
				"Full Name: John Smith",
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			},
		},
		{
			name: "functions",
			file: "functions.pseudo",
			expected: []string{
				"Hello, World!",
				"Factorial of 5: 120",
				"Factorial of 7: 5040",
				"Is 7 prime? TRUE",
				"Is 10 prime? FALSE",
				"Is 13 prime? TRUE",
				"Max of 10 and 25: 25",
				"Max of 100 and 50: 100",
			},
		},
		{
			name: "records",
			file: "records.pseudo",
			expected: []string{
				"Student 1:",
				"Name: Alice Johnson",
				"Age: 17",
				"Grade: A",
				"Student 2:",
				"Name: Bob Smith",
				"Age: 16",
				"Grade: B",
			},
		},
		{
			name: "oop",
			file: "oop.pseudo",
			expected: []string{
				"My dog's name: Buddy",
				"My dog's age: 3",
				"My dog's breed: Golden Retriever",
				"Buddy barks: Woof!",
				"My cat's name: Whiskers",
				"My cat's age: 5",
				"Whiskers meows: Meow!",
			},
		},
		{
			name:    "fileio",
			file:    "fileio.pseudo",
			skip:    true,
			skipMsg: "File I/O test requires filesystem access and creates files",
		},
	}

	// Find the examples directory
	examplesDir := findExamplesDir(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipMsg)
			}

			filePath := filepath.Join(examplesDir, tt.file)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read example file %s: %v", tt.file, err)
			}

			output, err := runProgram(string(content))
			if err != nil {
				t.Fatalf("program execution failed: %v", err)
			}

			for _, expectedLine := range tt.expected {
				if !strings.Contains(output, expectedLine) {
					t.Errorf("expected output to contain %q\nGot output:\n%s", expectedLine, output)
				}
			}
		})
	}
}

// findExamplesDir locates the examples directory relative to the test
func findExamplesDir(t *testing.T) string {
	// Try common relative paths
	paths := []string{
		"../../examples",
		"../examples",
		"examples",
	}

	for _, p := range paths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p
		}
	}

	// Try to find from working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}

	// Walk up looking for examples directory
	dir := wd
	for i := 0; i < 5; i++ {
		examplesPath := filepath.Join(dir, "examples")
		if info, err := os.Stat(examplesPath); err == nil && info.IsDir() {
			return examplesPath
		}
		dir = filepath.Dir(dir)
	}

	t.Fatal("could not find examples directory")
	return ""
}

// TestExamplesFileIO tests file I/O example separately with cleanup
func TestExamplesFileIO(t *testing.T) {
	examplesDir := findExamplesDir(t)
	filePath := filepath.Join(examplesDir, "fileio.pseudo")

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read example file: %v", err)
	}

	// Clean up any existing test file before running
	os.Remove("test_output.txt")
	defer os.Remove("test_output.txt")

	output, err := runProgram(string(content))
	if err != nil {
		t.Fatalf("program execution failed: %v", err)
	}

	expected := []string{
		"Writing to file...",
		"File written successfully!",
		"Reading from file...",
		"File read successfully!",
		"Appending to file...",
		"Append complete!",
	}

	for _, expectedLine := range expected {
		if !strings.Contains(output, expectedLine) {
			t.Errorf("expected output to contain %q\nGot output:\n%s", expectedLine, output)
		}
	}
}
