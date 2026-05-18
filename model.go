package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application states
type state int

const (
	stateWelcome state = iota
	stateInstructions
	stateFirstEntry
	stateTransition
	stateSecondEntry
	stateDone
)

// ── Claude Code inspired palette ──
// Warm amber/orange primary, soft whites, muted grays, no purple.
var (
	cAmber    = lipgloss.Color("#E8975B")
	cOrange   = lipgloss.Color("#D4723C")
	cPeach    = lipgloss.Color("#F5C396")
	cRed      = lipgloss.Color("#E05252")
	cGreen    = lipgloss.Color("#73C991")
	cDimGreen = lipgloss.Color("#3D6B4F")
	cWhite    = lipgloss.Color("#E8E0D5")
	cGray     = lipgloss.Color("#6B6560")
	cDimGray  = lipgloss.Color("#3D3936")
	cBg       = lipgloss.Color("#1A1816")
	cBgLight  = lipgloss.Color("#252220")
)

// ── Styles ──
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cAmber)

	slotFilled = lipgloss.NewStyle().
			Bold(true).
			Foreground(cAmber)

	slotActiveType = lipgloss.NewStyle().
			Bold(true).
			Foreground(cGreen)

	slotActiveDel = lipgloss.NewStyle().
			Bold(true).
			Foreground(cRed)

	slotEmpty = lipgloss.NewStyle().
			Foreground(cDimGray)

	bigDigitStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cOrange)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(cWhite)

	typeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cGreen)

	deleteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cRed)

	progressStyle = lipgloss.NewStyle().
			Foreground(cGray)

	doneStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cGreen)

	helpStyle = lipgloss.NewStyle().
			Foreground(cDimGray).
			Italic(true)

	accentBar = lipgloss.NewStyle().
			Foreground(cAmber)

	filledDot = lipgloss.NewStyle().
			Foreground(cAmber)

	emptyDot = lipgloss.NewStyle().
			Foreground(cDimGray)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cDimGray).
			Padding(1, 3)

	phaseTag = lipgloss.NewStyle().
			Background(cOrange).
			Foreground(cBg).
			Bold(true).
			Padding(0, 1)

	doneTag = lipgloss.NewStyle().
			Background(cGreen).
			Foreground(cBg).
			Bold(true).
			Padding(0, 1)
)

type model struct {
	state       state
	passcode    [4]int
	firstSeq    []Step
	secondSeq   []Step
	currentStep int
}

func initialModel() model {
	passcode := generatePasscode()
	return model{
		state:    stateWelcome,
		passcode: passcode,
		firstSeq: generateSequence(passcode),
		secondSeq: generateSequence(passcode),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == stateDone || m.state == stateWelcome {
				return m, tea.Quit
			}
		case "enter", " ":
			switch m.state {
			case stateWelcome:
				m.state = stateInstructions
			case stateInstructions:
				m.state = stateFirstEntry
				m.currentStep = 0
			case stateFirstEntry:
				if m.currentStep < len(m.firstSeq)-1 {
					m.currentStep++
				} else {
					m.state = stateTransition
				}
			case stateTransition:
				m.state = stateSecondEntry
				m.currentStep = 0
			case stateSecondEntry:
				if m.currentStep < len(m.secondSeq)-1 {
					m.currentStep++
				} else {
					m.state = stateDone
				}
			case stateDone:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var content string
	switch m.state {
	case stateWelcome:
		content = m.viewWelcome()
	case stateInstructions:
		content = m.viewInstructions()
	case stateFirstEntry:
		content = m.viewEntry(m.firstSeq, "FIRST ENTRY", 1)
	case stateTransition:
		content = m.viewTransition()
	case stateSecondEntry:
		content = m.viewEntry(m.secondSeq, "CONFIRM", 2)
	case stateDone:
		content = m.viewDone()
	default:
		content = ""
	}
	return "\n" + content + "\n"
}

func ruler() string {
	return accentBar.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func thinRuler() string {
	return lipgloss.NewStyle().Foreground(cDimGray).Render("┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄")
}

func (m model) viewWelcome() string {
	logo := bigDigitStyle.Render(
		"   _____ _______ _\n" +
			"  / ____|__   __| |\n" +
			" | (___    | |  | |     ___   ___| | _____ _ __\n" +
			"  \\___ \\   | |  | |    / _ \\ / __| |/ / _ \\ '__|\n" +
			"  ____) |  | |  | |___| (_) | (__|   <  __/ |\n" +
			" |_____/   |_|  |______\\___/ \\___|_|\\_\\___|_|")

	desc := subtitleStyle.Render(
		"  Generates a random Screen Time passcode and walks\n" +
			"  you through entering it digit by digit — with fake\n" +
			"  digits mixed in so you can't memorize it.")

	warn := lipgloss.NewStyle().Foreground(cOrange).Render(
		"  ⚠  The passcode is never shown in full.\n" +
			"  ⚠  It is not stored anywhere. This is a one-way trip.")

	help := helpStyle.Render("  [enter] start  ·  [ctrl+c] quit")

	return logo + "\n\n" + ruler() + "\n\n" + desc + "\n\n" + warn + "\n\n" + thinRuler() + "\n" + help
}

func (m model) viewInstructions() string {
	tag := phaseTag.Render(" SETUP ")
	title := titleStyle.Render("  Before we begin")

	steps := subtitleStyle.Render(
		"  ❶  Open  System Settings  on your Mac\n" +
			"  ❷  Navigate to  Screen Time\n" +
			"  ❸  Click  Lock Screen Time Settings\n" +
			"  ❹  Have the passcode input field visible")

	note := lipgloss.NewStyle().Foreground(cPeach).Render(
		"  You'll type digits and delete them as instructed.\n" +
			"  Don't try to peek — just follow the steps.")

	help := helpStyle.Render("  [enter] I'm ready  ·  [ctrl+c] abort")

	return ruler() + "\n" + tag + "  " + title + "\n" + ruler() + "\n\n" +
		steps + "\n\n" + thinRuler() + "\n\n" + note + "\n\n" + help
}

func (m model) viewEntry(seq []Step, label string, round int) string {
	step := seq[m.currentStep]

	tag := phaseTag.Render(fmt.Sprintf(" %s %d/%d ", label, m.currentStep+1, len(seq)))

	// Compute field count BEFORE this step
	countBefore := 0
	for i := 0; i < m.currentStep; i++ {
		if seq[i].Action == "type" {
			countBefore++
		} else if countBefore > 0 {
			countBefore--
		}
	}

	// Determine each slot's character and style (ASCII only for alignment)
	chars := [4]string{"_", "_", "_", "_"}
	charStyles := [4]lipgloss.Style{slotEmpty, slotEmpty, slotEmpty, slotEmpty}
	underStyles := [4]lipgloss.Style{slotEmpty, slotEmpty, slotEmpty, slotEmpty}

	for i := 0; i < countBefore; i++ {
		chars[i] = "*"
		charStyles[i] = slotFilled
	}

	var actionLine string
	if step.Action == "type" {
		chars[countBefore] = fmt.Sprintf("%d", step.Digit)
		charStyles[countBefore] = slotActiveType
		underStyles[countBefore] = slotActiveType
		actionLine = typeStyle.Render(fmt.Sprintf("  ▸ TYPE: %d", step.Digit))
	} else {
		if countBefore > 0 {
			chars[countBefore-1] = "X"
			charStyles[countBefore-1] = slotActiveDel
			underStyles[countBefore-1] = slotActiveDel
		}
		actionLine = deleteStyle.Render("  ✕ DELETE (backspace)")
	}

	// Build field line: fixed positions, 1-char slots, 5 spaces between
	// Positions: col 8, 14, 20, 26
	field := "        "
	for i := 0; i < 4; i++ {
		field += charStyles[i].Render(chars[i])
		if i < 3 {
			field += "     "
		}
	}

	// Build underline row: "───" centered under each 1-char slot
	// Positions: col 7-9, 13-15, 19-21, 25-27
	under := "       "
	for i := 0; i < 4; i++ {
		under += underStyles[i].Render("───")
		if i < 3 {
			under += "   "
		}
	}

	// Progress bar
	total := len(seq)
	done := m.currentStep + 1
	barWidth := 36
	filled := (done * barWidth) / total
	bar := accentBar.Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(cDimGray).Render(strings.Repeat("░", barWidth-filled))
	progressLine := progressStyle.Render(fmt.Sprintf("  %s  %d/%d", bar, done, total))

	help := helpStyle.Render("  [enter/space] next step")

	return ruler() + "\n" + tag + "\n" + ruler() + "\n\n" +
		field + "\n" + under + "\n\n" +
		actionLine + "\n\n" +
		progressLine + "\n\n" + help
}

func (m model) viewTransition() string {
	tag := doneTag.Render(" ROUND 1 COMPLETE ")

	check := doneStyle.Render(
		"  ✓  First entry done.\n")

	body := subtitleStyle.Render(
		"  macOS will now ask you to confirm the passcode.\n" +
			"  Same digits, different distraction pattern.")

	warn := lipgloss.NewStyle().Foreground(cOrange).Render(
		"  ⚠  Don't try to remember what you typed.")

	help := helpStyle.Render("  [enter] continue to confirmation")

	return ruler() + "\n" + tag + "\n" + ruler() + "\n\n" +
		check + "\n" + body + "\n\n" + warn + "\n\n" + help
}

func (m model) viewDone() string {
	tag := doneTag.Render(" LOCKED ")

	check := doneStyle.Render(
		"  ✓  Screen Time passcode set successfully.\n")

	body := subtitleStyle.Render(
		"  You should have no idea what the passcode is.\n" +
			"  That's the whole point.\n\n" +
			"  To reset: Apple ID  ▸  Screen Time  ▸  Change Passcode\n" +
			"  Or a full device reset.")

	help := helpStyle.Render("  [q/enter] exit")

	skull := bigDigitStyle.Render(
		"     ░░░░░░░\n" +
			"   ░░       ░░\n" +
			"  ░  ■   ■   ░\n" +
			"  ░     ▼     ░\n" +
			"  ░  ▀▀▀▀▀   ░\n" +
			"   ░░       ░░\n" +
			"     ░░░░░░░")

	return ruler() + "\n" + tag + "\n" + ruler() + "\n\n" +
		skull + "\n\n" + check + "\n" + body + "\n\n" + help
}
