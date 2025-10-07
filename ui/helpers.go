package ui

import (
	"fmt"
	"strings"
	"time"
)

func (ui *UI) clear() {
	ui.print("\033[2J\033[H")
}

func (ui *UI) print(text string) {
	ui.term.Write([]byte(text))
}

func (ui *UI) println(text string) {
	ui.term.Write([]byte(text + "\n"))
}

func (ui *UI) readLine(prompt string) string {
	if prompt != "" {
		ui.print(prompt)
	}
	line, _ := ui.term.ReadLine()
	return line
}

func (ui *UI) printLine() {
	ui.println(strings.Repeat("─", 60))
}

func (ui *UI) printHeader(text string) {
	ui.println("")
	ui.println(fmt.Sprintf("╔%s╗", strings.Repeat("═", len(text)+2)))
	ui.println(fmt.Sprintf("║ %s ║", text))
	ui.println(fmt.Sprintf("╚%s╝", strings.Repeat("═", len(text)+2)))
}

func (ui *UI) printError(msg string) {
	ui.println(fmt.Sprintf("\033[31m✗ %s\033[0m", msg))
}

func (ui *UI) printSuccess(msg string) {
	ui.println(fmt.Sprintf("\033[32m✓ %s\033[0m", msg))
}

func (ui *UI) formatTime(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 02, 2006")
	}
}
