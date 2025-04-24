package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	tm               *TuringMachine
	state            string // "input", "mode", "running", "finished"
	sequenceInput    textinput.Model
	modeSelected     bool
	stepByStep       bool
	stepCount        int
	maxSteps         int
	spinner          spinner.Model
	transitionsTable table.Model
	width, height    int
	err              error
}

// Init initializes the model
func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter sequence or leave blank for default"
	ti.Focus()
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		state:         "input",
		sequenceInput: ti,
		modeSelected:  false,
		stepByStep:    false,
		stepCount:     0,
		maxSteps:      100,
		spinner:       s,
		width:         80,
		height:        24,
	}
}

// Define some styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Italic(true)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B7"))

	symbolStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF8700")).
			Bold(true)

	stateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	headStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3333")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3333")).
			Bold(true)

	tableStyle = table.Styles{
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true),
		Cell: lipgloss.NewStyle().
			Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C3C3C")).
			Bold(true),
	}
)

// Messages for Bubbletea
type initTMMsg struct{ tm *TuringMachine }
type stepMsg struct{}
type errMsg struct{ err error }

// Initialize the Turing Machine
func initTM(sequence string) tea.Cmd {
	return func() tea.Msg {
		if sequence == "" {
			sequence = "01010101001101001010010011010001000100010110001001000101011000100010010010110001010010010111111100"
		}
		return initTMMsg{tm: NewTuringMachine(sequence)}
	}
}

// Command for stepping the Turing Machine
func stepTM() tea.Cmd {
	return func() tea.Msg {
		return stepMsg{}
	}
}

// Command for automatic stepping with a delay
func autoStep() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return stepMsg{}
	})
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case "input":
			switch msg.String() {
			case "enter":
				m.state = "mode"
				return m, initTM(m.sequenceInput.Value())
			case "ctrl+c", "q":
				return m, tea.Quit
			default:
				m.sequenceInput, cmd = m.sequenceInput.Update(msg)
				return m, cmd
			}

		case "mode":
			switch msg.String() {
			case "1":
				m.stepByStep = false
				m.modeSelected = true
				m.state = "running"
				return m, autoStep()
			case "2":
				m.stepByStep = true
				m.modeSelected = true
				m.state = "running"
				return m, stepTM()
			case "ctrl+c", "q":
				return m, tea.Quit
			}

		case "running":
			if m.stepByStep && msg.String() == "enter" {
				return m, stepTM()
			}
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}

		case "finished":
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.tm != nil {
			m.transitionsTable.SetWidth(m.width - 4)
			m.transitionsTable.SetHeight(10)
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case initTMMsg:
		m.tm = msg.tm

		// Setup transitions table
		columns := []table.Column{
			{Title: "#", Width: 4},
			{Title: "From", Width: 8},
			{Title: "Read", Width: 8},
			{Title: "To", Width: 8},
			{Title: "Write", Width: 8},
			{Title: "Move", Width: 8},
		}

		rows := make([]table.Row, 0)
		for key, transition := range m.tm.transitions {
			parts := strings.Split(key, ",")
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", len(rows)+1),
				parts[0],
				parts[1],
				transition[0],
				transition[1],
				transition[2],
			})
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithHeight(10),
			table.WithWidth(m.width-4),
			table.WithStyles(tableStyle),
		)

		m.transitionsTable = t
		return m, nil

	case stepMsg:
		if m.tm != nil && !m.tm.isFinished && m.maxSteps > 0 {
			m.stepCount++
			m.maxSteps--
			m.tm.Step()

			if m.tm.isFinished || m.maxSteps == 0 {
				m.state = "finished"
				return m, nil
			}

			if !m.stepByStep {
				return m, autoStep()
			}
		}

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	switch m.state {
	case "input":
		return inputView(m)
	case "mode":
		return modeView(m)
	case "running":
		return runningView(m)
	case "finished":
		return finishedView(m)
	default:
		return "Something went wrong!"
	}
}

func inputView(m Model) string {
	title := titleStyle.Render("Universal Turing Machine Simulator")

	info := infoStyle.Render(`
A Universal Turing Machine can simulate any other Turing Machine.
Enter a binary sequence representing a Turing Machine program followed by initial tape contents.
Format: <program>111<tape>
`)

	input := fmt.Sprintf(
		"%s\n%s\n\n%s\n%s",
		"Enter sequence (or leave blank for default example):",
		m.sequenceInput.View(),
		"Press Enter to continue",
		"Press q to quit",
	)

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, info, input)
}

func modeView(m Model) string {
	title := titleStyle.Render("Universal Turing Machine Simulator")

	if m.tm == nil {
		return fmt.Sprintf("%s\n\nInitializing Turing Machine... %s", title, m.spinner.View())
	}

	// Display transitions table
	transitionsView := fmt.Sprintf("\nTransitions Table:\n%s\n", m.transitionsTable.View())

	// Mode selection
	modePrompt := `
Select execution mode:
1) Automatic - Run steps automatically
2) Step-by-step - Press Enter to proceed to next step

Press 1 or 2 to select mode...
Press q to quit
`

	return fmt.Sprintf("%s\n%s\n%s", title, transitionsView, modePrompt)
}

func runningView(m Model) string {
	title := titleStyle.Render("Universal Turing Machine Simulator")

	if m.tm == nil {
		return fmt.Sprintf("%s\n\nInitializing...", title)
	}

	// Execution information
	info := fmt.Sprintf("\nStep: %s\nState: %s\nSymbol: %s",
		highlightStyle.Render(fmt.Sprintf("%d", m.stepCount)),
		stateStyle.Render(m.tm.currentState),
		symbolStyle.Render(m.tm.tape.Read()),
	)

	// Tape visualization
	tapeView := renderTape(m.tm.tape)

	// Controls
	controls := fmt.Sprintf("\nMax steps remaining: %d", m.maxSteps)
	if m.stepByStep {
		controls += "\n\nPress Enter for next step..."
	} else {
		controls += fmt.Sprintf("\n\n%s Running automatically...", m.spinner.View())
	}
	controls += "\nPress q to quit"

	return fmt.Sprintf("%s%s\n\n%s\n%s", title, info, tapeView, controls)
}

func finishedView(m Model) string {
	title := titleStyle.Render("Universal Turing Machine Simulator")

	result := "\nResult: "
	if m.tm.IsAccepting() {
		result += successStyle.Render("✓ Accepted")
	} else {
		result += errorStyle.Render("✗ Rejected")
	}

	info := fmt.Sprintf("\nCompleted in %d steps\nFinal state: %s",
		m.stepCount,
		stateStyle.Render(m.tm.currentState),
	)

	tapeView := renderTape(m.tm.tape)

	return fmt.Sprintf("%s%s%s\n\n%s\n\nPress q to quit", title, result, info, tapeView)
}

func renderTape(tape *Tape) string {
	// Find display range (show roughly 40 chars centered on head)
	start := tape.headPosition - 20
	if start < 0 {
		start = 0
	}
	end := start + 40
	if end > len(tape.content) {
		end = len(tape.content)
	}

	// Create head position indicator
	headPos := strings.Repeat(" ", tape.headPosition-start) + headStyle.Render(" ▼")

	// Create tape view with highlighted current position
	tapeView := ""
	for i := start; i < end; i++ {
		symb := "_"
		if i < len(tape.content) {
			symb = tape.content[i]
		}

		if i == tape.headPosition {
			tapeView += highlightStyle.Render("[" + symb + "]")
		} else {
			tapeView += symb
		}
	}

	return fmt.Sprintf("%s\n%s", headPos, tapeView)
}
