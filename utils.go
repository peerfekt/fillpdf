/*
 *  FillPDF - Fill PDF forms
 *  Copyright DesertBit
 *  Author: Roland Singer, Alexander Félix
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package fillpdf

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getAbs(path string) (string, error) {
	// Get the absolute paths.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get the absolute path: %v", err)
	}

	// Check if the form file exists.
	e, err := exists(path)
	if err != nil {
		return "", fmt.Errorf("failed to check if file exists: %v", err)
	} else if !e {
		return "", fmt.Errorf("file does not exists: '%s'", path)
	}

	return absPath, nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// copyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// runCommandInPath runs a command and waits for it to exit.
// The working directory is also set.
// The stderr error message is returned on error.
func runCommandInPath(dir, name string, args ...string) error {
	// Create the command.
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Dir = dir

	// Start the command and wait for it to exit.
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	return nil
}

func runCommandWithOutput(dir, name string, args ...string) ([]byte, error) {
	// Create the command.
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Dir = dir
	// Start the command and wait for it to exit.
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	return stdout.Bytes(), nil
}

//create Random ID
func GetID(prefix string) (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%x", prefix, b), nil
}
