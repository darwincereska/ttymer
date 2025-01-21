package main

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"time"
)

type DurationTimer struct {
	seconds int
	minutes int
	hours   int
	days    int
}

type TimeTimer struct {
	time int
	pm   bool
}

func convertDurationTime(days int, hours int, minutes int, seconds int) int {
	a := DurationTimer{
		days:    days,
		hours:   hours,
		minutes: minutes,
		seconds: seconds,
	}
	total := 0
	total += a.days * 24 * 60 * 60
	total += a.hours * 60 * 60
	total += a.minutes * 60
	total += a.seconds
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

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%ds ", seconds)
	}
	return result
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
	fmt.Println()

	// Send notification
	exec.Command("notify-send", "Timer Complete", fmt.Sprintf("Your %s timer has ended", initialTime)).Run()
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

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	// Duration mode
	if *time_flag == 0 && (*seconds + *minutes + *hours + *days) == 0 {
		err := errors.New("Please specify a duration using -s, -m, -h or -d flags")
		fmt.Println(err) 
		return
	}

	if *time_flag != 0 {
		if *pm_flag {
			countdown(convertClockTime(*time_flag, true)) // Converts clock time into seconds and then calls the countdown
		} else {
			countdown(convertClockTime(*time_flag, false)) // Converts clock time into seconds and then calls the countdown
		}
	} else {
		countdown(convertDurationTime(*days, *hours, *minutes, *seconds))
	}
}
