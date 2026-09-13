package views

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
)

func TestRenderViewsBorderAlignment(t *testing.T) {
	eng := engine.NewGame("Sang Petualang", "Oakhaven")

	testCases := []struct {
		name     string
		rendered string
	}{
		{"TownView", RenderTownView(eng, 80)},
		{"WorkerView", RenderWorkerView(eng, 0, 80)},
		{"BuildView", RenderBuildView(eng, 0, 80)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.rendered, "\n")
			for i, line := range lines {
				if len(strings.TrimSpace(line)) == 0 {
					continue
				}
				w := lipgloss.Width(line)
				if w != 78 {
					t.Errorf("line %d has width %d (expected 78): %q", i, w, line)
				}
			}
		})
	}
}
