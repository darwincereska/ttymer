package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mbndr/figlet4go"
	"golang.org/x/term"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TUI CODE START ___________________________________
type model struct {
	timeRemaining int
	initialTime   string
	quitting      bool
}

// Init implements bubbletea.Model
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		tea.EnterAltScreen,
	)
}

// Tick command
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Messages
type tickMsg struct{}

// Update implements bubbletea.Model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		if m.timeRemaining > 0 {
			m.timeRemaining--
			return m, tick()
		}
		if m.timeRemaining == 0 {
			notification("Timer Ended", fmt.Sprintf("Your %s timer has ended", m.initialTime))
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// View implements bubbletea.Model
func (m model) View() string {
	if m.quitting {
		return "Timer ended!\n"
	}

	// Get terminal dimensions
	width, height := 80, 24 // default values
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width, height = w, h
	}

	// Create styles
	containerStyle := lipgloss.NewStyle().
		Width(width - 2).   // Subtract 2 for borders
		Height(height - 2). // Subtract 2 for top and bottom borders
		Align(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FAFAFA")).
		Padding(1)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF75B5")).
		Bold(true).
		Width(width - 4). // Subtract border and padding
		Align(lipgloss.Center)

	quitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(width - 4). // Subtract border and padding
		Align(lipgloss.Center)

	timerStyle := lipgloss.NewStyle().
		Width(width - 4). // Subtract border and padding
		Align(lipgloss.Center)

	// Create figlet
	ascii := figlet4go.NewAsciiRender()
	renderStr := formatTimeString(m.timeRemaining)

	// Set some options
	options := figlet4go.NewRenderOptions()
	options.FontName = "standard"

	ascii.LoadFont("standard")
	rendered, _ := ascii.RenderOpts(renderStr, options)

	// Create the content
	title := titleStyle.Render("TTYMER")
	timerText := timerStyle.Render(rendered)
	quitText := quitStyle.Render("Press Ctrl+C to quit")

	// Calculate content height
	contentHeight := lipgloss.Height(title) + lipgloss.Height(timerText) + lipgloss.Height(quitText) + 4 // +4 for spacing

	// Calculate padding for vertical centering
	paddingTop := (height - contentHeight - 4) / 2 // -4 for borders and padding
	if paddingTop < 0 {
		paddingTop = 0
	}
	paddingBottom := height - contentHeight - paddingTop - 4 // -4 for borders and padding
	if paddingBottom < 0 {
		paddingBottom = 0
	}

	// Combine content with padding
	content := strings.Repeat("\n", paddingTop) +
		title + "\n\n" +
		timerText + "\n\n" +
		quitText +
		strings.Repeat("\n", paddingBottom)

	// Return the final view
	return containerStyle.Render(content)
}

// TUI CODE END ____________________________

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
		initialSeconds := settings.GetSeconds()
		p := tea.NewProgram(
			model{
				timeRemaining: initialSeconds,
				initialTime:   formatTimeString(initialSeconds),
			},
			tea.WithAltScreen(),       // use the full screen
			tea.WithMouseCellMotion(), // enable mouse support
		)

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}
}
