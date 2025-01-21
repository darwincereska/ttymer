package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)

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
	return
}

func main() {
	help := flag.Bool("help", false, "show help menu")
	days := flag.Int("d", 0, "days")
	hours := flag.Int("h", 0, "hours")
	minutes := flag.Int("m", 0, "minutes")
	seconds := flag.Int("s", 0, "seconds")
	time_flag := flag.Int("t", 0, "time mode (format: 1-12 or 100-1159)")
	pm_flag := flag.Bool("p", false, "am/pm")
	ui_flag := flag.Bool("ui", false, "Toggles a tui")
	flag.Parse()

	if *help {
		flag.Usage()
	}
	type Settings struct {
		days    int
		hours   int
		minutes int
		seconds int
		time    int
		pm      bool
		ui      bool
	}
	settings := Settings{
		days:    *days,
		hours:   *hours,
		minutes: *minutes,
		seconds: *seconds,
		time:    *time_flag,
		pm:      *pm_flag,
		ui:      *ui_flag,
	}
	if settings.ui {
		if settings.time != 0 {
			if settings.pm {
				countdown(convertClockTime(settings.time, true))
			} else {
				countdown(convertClockTime(settings.time, false))
			}
		} else {
			countdown(convertDurationTime(settings.days, settings.hours, settings.minutes, settings.seconds))
		}
	} else {
		if settings.time != 0 {
			if settings.pm {
				countdown(convertClockTime(settings.time, true))
			} else {
				countdown(convertClockTime(settings.time, false))
			}
		} else {
			countdown(convertDurationTime(settings.days, settings.hours, settings.minutes, settings.seconds))
		}
	}

}
