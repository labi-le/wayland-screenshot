package main

import (
	"bytes"
	"errors"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const FormatFileDefault = "2006-01-02_15-04-05"

type Capture int

var ErrUnknownCaptureType = errors.New("unknown capture type")

const (
	Area Capture = iota
	All
)

// ParseCapture parses a string to a Capture type
func ParseCapture(s string) (Capture, error) {
	switch s {
	case "area":
		return Area, nil
	case "all":
		return All, nil
	}

	return 0, ErrUnknownCaptureType
}

func (c Capture) String() string {
	switch c {
	case Area:
		return "area"
	case All:
		return "all"
	}

	return ""
}

// CaptureArea captures a region
func CaptureArea(r Region) (Captured, error) {
	output, err := exec.Command("grim", "-g", r.String(), "-").Output() //nolint:gosec // dn
	return Captured{output}, err
}

// CaptureAll captures the whole screen
func CaptureAll() (Captured, error) {
	output, err := exec.Command("grim", "-").Output()
	return Captured{output}, err
}

// Captured is a captured image
type Captured struct {
	Raw []byte
}

// SaveFile saves the captured image to a file
func (c *Captured) SaveFile(path string) error {
	return saveFileToDir(c.Raw, path)
}

// SaveClipboard saves the captured image to the clipboard
func (c *Captured) SaveClipboard() error {
	cmd := exec.Command("wl-copy")
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewReader(c.Raw)
	return cmd.Run()
}

// Edit edits the captured image with swappy
func (c *Captured) Edit() (Captured, error) {
	cmd := exec.Command("swappy", "-f", "-", "-o", "-")
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewReader(c.Raw)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return Captured{}, err
	}

	return Captured{out.Bytes()}, nil
}

// generateNameDefaultFormat generates a name for a file with the default format
func generateNameDefaultFormat(raw []byte) (string, error) {
	ext, err := getExt(raw)
	if err != nil {
		return "", err
	}

	name := time.Now().Format(FormatFileDefault) + ext

	return name, nil
}

// getExt returns the extension of a mime type
func getExt(raw []byte) (string, error) {
	byType, mimeTypeErr := mime.ExtensionsByType(http.DetectContentType(raw))
	if mimeTypeErr != nil {
		return "", mimeTypeErr
	}

	if len(byType) == 0 {
		return "", errors.New("unknown file type")
	}

	return byType[0], nil
}

// saveFileToDir saves a file to a directory
func saveFileToDir(raw []byte, path string) error {
	if err := ensureDir(path); err != nil {
		return err
	}

	p, err := isPath(path)
	if err != nil {
		return err
	}

	if p {
		name, genErr := generateNameDefaultFormat(raw)
		if genErr != nil {
			return err
		}

		path = path + "/" + name
	}

	file, createErr := os.Create(path)
	if createErr != nil {
		return createErr
	}

	_, err = file.Write(raw)

	return err
}

// ensureDir create recursively a directory
func ensureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		return os.MkdirAll(dirName, os.ModePerm)
	}

	return nil
}

// isPath checks if a path is a directory or a file
func isPath(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		// if no such file or directory, is ok
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return stat.IsDir(), nil
}
