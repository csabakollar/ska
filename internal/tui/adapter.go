package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"regexp"
	"strings"
)

func NewModelFromInteractiveForm(iForm InteractiveForm, header string) *Model {
	m := &Model{
		header: header,
		inputs: make([]textinput.Model, len(iForm.Inputs)),
	}

	for i := range iForm.Inputs {
		t := textinput.New()

		// Prompt
		t.Prompt = fmt.Sprintf("%s: ", iForm.Inputs[i].Label)
		// Placeholder
		t.Placeholder = iForm.Inputs[i].Help
		t.PlaceholderStyle = helpStyle
		// Validation
		t.Validate = validator(iForm.Inputs[i].RegExp)
		if iForm.Inputs[i].MaxLength > 0 {
			t.CharLimit = iForm.Inputs[i].MaxLength
		}
		// Default
		if iForm.Inputs[i].Default != "" {
			t.SetValue(iForm.Inputs[i].Default)
		}
		// First Item
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = noStyle
		}
		m.inputs[i] = t
	}
	return m
}

func validator(regexpString string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("value cannot be empty")
		}
		if strings.TrimSpace(regexpString) != "" {
			re := regexp.MustCompile(regexpString)
			if !re.MatchString(s) {
				return fmt.Errorf("invalid value. It should match %s", regexpString)
			}
		}
		return nil
	}
}
