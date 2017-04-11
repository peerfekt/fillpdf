package fillpdf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Multistamp stamps one PDF ontop of another, returns a reader to bytes generated.
func Multistamp(stampontoPDFFile, stampPDFFile string) (io.Reader, error) {
	var err error
	stampontoPDFFile, err = getAbs(stampontoPDFFile)
	if err != nil {
		return nil, nil
	}

	stampPDFFile, err = getAbs(stampPDFFile)
	if err != nil {
		return nil, nil
	}

	// Check if the pdftk utility exists.
	_, err = exec.LookPath("pdftk")
	if err != nil {
		return nil, fmt.Errorf("pdftk utility is not installed!")
	}

	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "fillpdf-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %v", err)
	}

	// Remove the temporary directory on defer again.
	defer func() {
		errD := os.RemoveAll(tmpDir)
		// Log the error only.
		if errD != nil {
			log.Printf("fillpdf: failed to remove temporary directory '%s' again: %v", tmpDir, errD)
		}
	}()

	// Create the temporary output file path.
	outputFile := filepath.Clean(tmpDir + "/output.pdf")

	// Create the pdftk command line arguments.
	args := []string{
		stampontoPDFFile,
		"multistamp", stampPDFFile,
		"output", outputFile,
	}

	// Run the pdftk utility.
	err = runCommandInPath(tmpDir, "pdftk", args...)
	if err != nil {
		return nil, fmt.Errorf("pdftk error: %v", err)
	}

	fb, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open generated file")
	}

	return bytes.NewReader(fb), nil
}
