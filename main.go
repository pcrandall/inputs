package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	te "github.com/muesli/termenv"
)

const focusedTextColor = "205"

var (
	color               = te.ColorProfile().Color
	focusedPrompt       = te.String("> ").Foreground(color("205")).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + te.String("Submit").Foreground(color("205")).String() + " ]"
	blurredSubmitButton = "[ " + te.String("Submit").Foreground(color("240")).String() + " ]"

	focusedAddButton = "[ " + te.String("Add").Foreground(color("205")).String() + " ]"
	blurredAddButton = "[ " + te.String("Add").Foreground(color("240")).String() + " ]"

	keyword       textinput.Model
	location      textinput.Model
	add           textinput.Model
	packageInputs = [][]textinput.Model{}
)

type model struct {
	index         int
	keywordInput  textinput.Model
	locationInput textinput.Model
	addInput      textinput.Model
	submitButton  string
	submit        bool
}

func main() {
	// s := &search{}
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Inputs here %+v", INPUTS)
	for _, val := range packageInputs {
		for i, v := range val {
			fmt.Println(i, v)
			fmt.Printf("Index: %d, Value: %+v\n\n", i, v.Value())
			// fmt.Println(i, v)
		}
	}
}

func initialModel() model {
	keyword = textinput.NewModel()
	keyword.Placeholder = "Keyword eg: Web Developer"
	keyword.Focus()
	keyword.Prompt = focusedPrompt
	keyword.TextColor = focusedTextColor
	keyword.CharLimit = 32

	location = textinput.NewModel()
	location.Placeholder = "Location eg: Boulder CO, Salt Lake City UT"
	location.Prompt = blurredPrompt
	location.CharLimit = 64

	add = textinput.NewModel()
	add.Placeholder = "[Y/n]"
	add.Prompt = blurredPrompt
	add.CharLimit = 8

	return model{0, keyword, location, add, blurredSubmitButton, false}

}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if m.submit == true {
		m.submit = false
		return m.UpdateAdd(msg)
	}
	return m.UpdateQuery(msg)
}

func (m model) UpdateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
		case "esc":
			return m, tea.Quit

		case "Y", "y", "enter":
			return m, nil
		}
	}

	// Handle character input and blinks
	m, cmd = updateAdd(msg, m)

	return m, cmd
}

func updateAdd(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.addInput, cmd = m.addInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) UpdateQuery(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
		case "esc":
			return m, tea.Quit

		// Cycle between input
		case "tab", "shift+tab", "enter", "up", "down":

			input := []textinput.Model{
				m.keywordInput,
				m.locationInput,
			}

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.index == len(input) {
				packageInputs = append(packageInputs, input)
				m, cmd = updateInputs("", m)
				initialModel().Init()
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.index--
			} else {
				m.index++
			}

			if m.index > len(input) {
				m.index = 0
			} else if m.index < 0 {
				m.index = len(input)
			}

			for i := 0; i <= len(input)-1; i++ {
				if i == m.index {
					// Set focused state
					input[i].Focus()
					input[i].Prompt = focusedPrompt
					input[i].TextColor = focusedTextColor
					continue
				}
				// Remove focused state
				input[i].Blur()
				input[i].Prompt = blurredPrompt
				input[i].TextColor = ""
			}

			m.keywordInput = input[0]
			m.locationInput = input[1]

			if m.index == len(input) {
				m.submitButton = focusedSubmitButton
			} else {
				m.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	}

	// Handle character input and blinks
	m, cmd = updateInputs(msg, m)

	return m, cmd
}

// Pass messages and models through to text input components. Only text inputs
// with Focus() set will respond, so it's safe to simply update all of them
// here without any further logic.
func updateInputs(msg tea.Msg, m model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.keywordInput, cmd = m.keywordInput.Update(msg)
	cmds = append(cmds, cmd)

	m.locationInput, cmd = m.locationInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string

	if m.submit == true {
		s = m.additionalView()
	} else {
		s = m.queryView()
	}

	return indent.String("\n"+s+"\n\n", 2)
}

func (m model) queryView() string {
	s := "\n"
	inputs := []string{
		m.keywordInput.View(),
		m.locationInput.View(),
	}
	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}
	s += "\n\n" + m.submitButton + "\n"
	return s
}

func (m model) additionalView() string {

	s := "Add another search item? [Y/n]"
	inputs := []string{
		m.addInput.View(),
	}
	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}
	s += "\n\n" + m.submitButton + "\n"
	return s
}
