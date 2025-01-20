package main

import (
	"flag"
	"fmt"
	"os/exec"
	"errors"
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

func NewDurationTimer(days int, hours int, minutes int, seconds int) DurationTimer {
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
	countdown(total)
	return a
}

func NewTimeTimer(clockTime int, pm bool) TimeTimer {
	a := TimeTimer{
		time: clockTime,
		pm:   pm,
	}

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

	secondsUntil := int(target.Sub(now).Seconds())
	countdown(secondsUntil)
	return a
}

func countdown(time int) {
	initialTime := time // Store initial time for notification
	for time > 0 {
		days := time / 86400
		hours := (time % 86400) / 3600
		minutes := (time % 3600) / 60
		seconds := time % 60

		fmt.Print("\r\033[K") // Clear the line before printing new time
		if days > 0 {
			fmt.Printf("%d"+"d ", days)
		}
		if hours > 0 {
			fmt.Printf("%d"+"h ", hours)
		}
		if minutes > 0 {
			fmt.Printf("%d"+"m ", minutes)
		}
		if seconds > 0 {
			fmt.Printf("%d"+"s ", seconds)
		}

		time--
		// Sleep for 1 second
		exec.Command("sleep", "1").Run()
	}
	fmt.Println()

	// Calculate original duration components for notification
	days := initialTime / 86400
	hours := (initialTime % 86400) / 3600
	minutes := (initialTime % 3600) / 60
	seconds := initialTime % 60

	// Build timer duration string
	durationStr := ""
	if days > 0 {
		durationStr += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		durationStr += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		durationStr += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 {
		durationStr += fmt.Sprintf("%ds", seconds)
	}

	// Send notification
	exec.Command("notify-send", "Timer Complete", fmt.Sprintf("Your %s timer has ended", durationStr)).Run()
	return
}

func main() {
	days := flag.Int("d", 0, "days")
	hours := flag.Int("h", 0, "hours")
	minutes := flag.Int("m", 0, "minutes")
	seconds := flag.Int("s", 0, "seconds")
	time_flag := flag.Int("t", 0, "time mode (format: 1-12 or 100-1159)")
	pm_flag := flag.Bool("p", false, "am/pm")

	flag.Parse()
	// Duration mode
	if *time_flag == 0 && *seconds == 0 && *minutes == 0 && *hours == 0 && *days == 0 {
		err := errors.New("at least seconds (-s) is required in duration mode")
		fmt.Println(err)
		return
	}

	if *time_flag != 0 {
		if *pm_flag == true {
			NewTimeTimer(*time_flag, true)
		} else {
			NewTimeTimer(*time_flag, false)
		}
	} else {
		NewDurationTimer(*days, *hours, *minutes, *seconds)
	}
}