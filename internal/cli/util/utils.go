package util

import (
	"os/exec"
	"runtime"
	"time"
)

func FormatLastSeenAt(lastSeenAt *time.Time) string {
	if lastSeenAt == nil {
		return "Never"
	}

	return lastSeenAt.Local().Format("2006-01-02 15:04")
}

func OpenBrowser(url string) error {
	var command *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		command = exec.Command("xdg-open", url)
	}

	return command.Run()
}
