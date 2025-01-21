package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)

type DisplayMode int

const (
	StandardMode DisplayMode = iota
	UIMode
)

type TimerType int

const (
	DurationTimer TimerType = iota
	ClockTimer
)

// These settings hold all the configuration for the timer
type Settings struct {
	DisplayMode DisplayMode
	TimerType   TimerType

	// Duration Mode Settings
	Duration struct {
		Days    int
		Hours   int
		Minutes int
		Seconds int
	}

	// Clock mode Settings
	Clock struct {
		Time int
		PM   bool
	}
}

// This NewSettings function creates and returns settings with default values
func NewSettings() *Settings {
	return &Settings{
		DisplayMode: StandardMode,
		TimerType:   DurationTimer,
	}
}

// this parseFlags func parses the command line flags into the settings struct
func ParseFlags() *Settings {
	settings := NewSettings()

	// Display mode flag
	uiFlag := flag.Bool("ui", false, "Use the TUI interface")

	// Duration Timer Flags
	flag.IntVar(&settings.Duration.Days, "d", 0, "days for timer")
	flag.IntVar(&settings.Duration.Hours, "h", 0, "hours for timer")
	flag.IntVar(&settings.Duration.Minutes, "m", 0, "minutes for timer")
	flag.IntVar(&settings.Duration.Seconds, "s", 0, "seconds for timer")

	// Clock Timer Flags
	flag.IntVar(&settings.Clock.Time, "t", 0, "time mode (format: 1-12 or 100-1159)")
	flag.BoolVar(&settings.Clock.PM, "p", false, "PM flag for clock mode")

	// Help flag
	helpFlag := flag.Bool("help", false, "Shows help menu")

	// Parse the flags
	flag.Parse()

	// Set Display mode
	if *uiFlag {
		settings.DisplayMode = UIMode
	}
	if *helpFlag {
		flag.Usage()
	}

	// Set timer type
	if settings.Clock.Time != 0 {
		settings.TimerType = ClockTimer
	}

	return settings

}

func convertDurationTime(days int, hours int, minutes int, seconds int) int {
	var total int
	total += days * 24 * 60 * 60
	total += hours * 60 * 60
	total += minutes * 60
	total += seconds
	return total
}

func convertClockTime(clockTime int, pm bool) int {
	var hours, minutes int

	// Handle 4 digit time (e.g. 1159) vs 2 digit time (e.g. 11)
	if clockTime >= 100 {
		hours = clockTime / 100
		minutes = clockTime % 100
	} else {
		hours = clockTime
		minutes = 0
	}

	// Convert PM times
	if pm && hours != 12 {
		hours += 12
	}
	// Convert 12 AM to 0 hours
	if !pm && hours == 12 {
		hours = 0
	}

	// Calculate total seconds until time
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, now.Location())

	// If target time is earlier today, move to tomorrow
	if target.Before(now) {
		target = target.Add(24 * time.Hour)
	}

	return int(target.Sub(now).Seconds())
}

// GetSeconds calculates total seconds based on settings
func (s *Settings) GetSeconds() int {
	switch s.TimerType {
	case DurationTimer:
		return convertDurationTime(
			s.Duration.Days,
			s.Duration.Hours,
			s.Duration.Minutes,
			s.Duration.Seconds,
		)
	case ClockTimer:
		return convertClockTime(
			s.Clock.Time,
			s.Clock.PM,
		)
	default:
		return 0
	}
}

func formatTimeString(timeInSeconds int) string {
	days := timeInSeconds / 86400
	hours := (timeInSeconds % 86400) / 3600
	minutes := (timeInSeconds % 3600) / 60
	seconds := timeInSeconds % 60
	var result string

	if days > 0 {
		result += fmt.Sprintf("%dd", days)
		if hours > 0 || minutes > 0 || seconds > 0 {
			result += " "
		}
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
		if minutes > 0 || seconds > 0 {
			result += " "
		}
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm", minutes)
		if seconds > 0 {
			result += " "
		}
	}
	if seconds > 0 {
		result += fmt.Sprintf("%ds", seconds)
	}
	return result
}

func notification(title string, text string) {
	exec.Command("notify-send", title, text).Run()
}

func countdown(time int) {
	if time > 0 {
		initialTime := formatTimeString(time)
		for time > 0 {
			fmt.Print("\r\033[K") // Clear the line before printing new time
			fmt.Print(formatTimeString(time))

			time--
			// Sleep for 1 second
			exec.Command("sleep", "1").Run()
		}

		fmt.Print("\r\033[K")
		fmt.Println("Your", initialTime, "timer has ended") // Prints Timer after its finished

		// Send notification
		notification("Timer Ended", fmt.Sprintf("Your %s timer has ended", initialTime))
	}
	return
}

func main() {
	settings := ParseFlags()

	switch settings.DisplayMode {
	case StandardMode:
		countdown(settings.GetSeconds())
	case UIMode:
		fmt.Println("UI mode is not implemented yet.")
	}
}
