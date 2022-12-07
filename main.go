package main

import (
	"flag"
	"os"
	"os/exec"
)

var showNotify = false //nolint:gochecknoglobals // dn

func main() {
	checkDependencies()
	captureType, path, needEdit := parseFlag()

	DoSave(
		MakeCapture(captureType),
		path,
		needEdit,
	)

	SendNotify(NotifyLow, "Screenshot saved")

}

// wrapErr wraps error and send notify if needed
//
//nolint:unparam // dn
func wrapErr(err error, exit bool) {
	if err != nil {
		SendNotify(NotifyCritical, err.Error())
		if exit {
			os.Exit(1)
		}
	}
}

// DoSave saves captured image to path\clipboard
func DoSave(c Captured, path string, edit bool) {
	if edit {
		reCaptured, err := c.Edit()
		wrapErr(err, true)

		c = reCaptured
	}

	if path == "" {
		wrapErr(c.SaveClipboard(), true)
		return
	}

	wrapErr(c.SaveFile(path), true)
}

// MakeCapture is interface for capture
func MakeCapture(captureType Capture) Captured {
	var (
		captured Captured
		err      error
	)
	switch captureType {
	case Area:
		region, selectErr := SelectRegion()
		wrapErr(selectErr, true)

		captured, err = CaptureArea(region)
		wrapErr(err, true)

	case All:
		captured, err = CaptureAll()
		wrapErr(err, true)

	default:
		wrapErr(ErrUnknownCaptureType, true)
	}

	return captured
}

func parseFlag() (Capture, string, bool) {
	flag.BoolVar(&showNotify, "notify", false, "Show notification with notify-send")

	captureType := flag.String("capture", Area.String(), "Capture area or all")
	path := flag.String("path", "", "Path to save file")
	edit := flag.Bool("edit", false, "Edit image with swappy")

	flag.Parse()

	capture, captErr := ParseCapture(*captureType)
	wrapErr(captErr, true)

	return capture, *path, *edit
}

func checkDependencies() {
	var needMessage string
	for _, dependency := range []string{"slurp", "grim", "wl-copy", "swappy"} {
		if _, err := exec.LookPath(dependency); err != nil {
			needMessage += dependency + " "
		}
	}

	if needMessage != "" {
		SendNotify(NotifyCritical, "Missing dependencies: "+needMessage)
	}

}
