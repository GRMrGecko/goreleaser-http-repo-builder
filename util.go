package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Helper for CLI to ask for confirmation.
func askForConfirmation(message string) bool {
	// Read stdanrd input for each new line.
	scanner := bufio.NewScanner(os.Stdin)

	// Loop the question until answered.
	for {
		fmt.Printf("%s [y/n]: ", message)

		// Get next line.
		scanner.Scan()
		resp := strings.ToLower(strings.TrimSpace(scanner.Text()))

		// Check if yes or no.
		switch resp {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Invalid answer.")
		}
	}
}

// Helper for copying files.
func copyFile(srcFile, dstFile string) (err error) {
	// Open the source file.
	f, err := os.Open(srcFile)
	if err != nil {
		return
	}
	defer f.Close()

	// Open the destination file.
	d, err := os.Create(dstFile)
	if err != nil {
		return
	}
	defer d.Close()

	// Copy the data to the new file.
	_, err = io.Copy(d, f)
	if err != nil {
		return
	}

	// Ensure new file is fully written.
	err = d.Sync()
	return
}
