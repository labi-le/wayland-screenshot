package main

import (
	"os"
	"os/exec"
)

type NotifyLevel int

func (l *NotifyLevel) String() string {
	switch *l {
	case NotifyLow:
		return "low"
	case NotifyNormal:
		return "normal"
	case NotifyCritical:
		return "critical"
	}
	return ""
}

const (
	NotifyLow NotifyLevel = iota
	NotifyNormal
	NotifyCritical
)

// SendNotify sends a notification
func SendNotify(level NotifyLevel, message string) {
	if showNotify {
		err := exec.Command("notify-send", "-u", level.String(), "wl-screenshot", message).Run() //nolint:gosec // dn
		if err != nil {
			_, _ = os.Stderr.WriteString(err.Error())
		}

	}
	_, _ = os.Stderr.WriteString(message)
}

// SendNotifyAndExit same with exit
func SendNotifyAndExit(level NotifyLevel, message string, exitCode int) {
	SendNotify(level, message)
	os.Exit(exitCode)
}
