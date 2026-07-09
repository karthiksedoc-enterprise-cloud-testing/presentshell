package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSlideNavigation(t *testing.T) {
	m := newModel("examples/demo.md")
	if m.err != nil {
		t.Fatalf("Failed to load: %v", m.err)
	}

	// Simulate WindowSizeMsg to initialize (skip terminal creation for test)
	m.ready = true
	m.layout.Width = 120
	m.layout.Height = 40

	if m.currentSlide != 0 {
		t.Fatalf("Expected slide 0, got %d", m.currentSlide)
	}

	// Test next slide with "right"
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = result.(model)
	if m.currentSlide != 1 {
		t.Fatalf("After 'right': expected slide 1, got %d", m.currentSlide)
	}

	// Test next slide with "n"
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = result.(model)
	if m.currentSlide != 2 {
		t.Fatalf("After 'n': expected slide 2, got %d", m.currentSlide)
	}

	// Test prev slide with "left"
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = result.(model)
	if m.currentSlide != 1 {
		t.Fatalf("After 'left': expected slide 1, got %d", m.currentSlide)
	}

	// Test prev slide with "p"
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = result.(model)
	if m.currentSlide != 0 {
		t.Fatalf("After 'p': expected slide 0, got %d", m.currentSlide)
	}

	// Test boundary - can't go below 0
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = result.(model)
	if m.currentSlide != 0 {
		t.Fatalf("After 'left' at 0: expected slide 0, got %d", m.currentSlide)
	}

	t.Logf("✅ All navigation tests passed! Total slides: %d", m.totalSlides)
}

func TestTabFocusSwitch(t *testing.T) {
	m := newModel("examples/demo.md")
	m.ready = true
	m.layout.Width = 120
	m.layout.Height = 40

	if m.layout.FocusedPane != 0 { // SlidePane
		t.Fatal("Expected initial focus on slides")
	}

	// Tab to terminal
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = result.(model)
	if m.layout.FocusedPane != 1 { // TerminalPane
		t.Fatal("Expected focus on terminal after Tab")
	}

	// Tab back to slides
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = result.(model)
	if m.layout.FocusedPane != 0 { // SlidePane
		t.Fatal("Expected focus back on slides after Tab")
	}

	t.Log("✅ Tab focus switching works!")
}
