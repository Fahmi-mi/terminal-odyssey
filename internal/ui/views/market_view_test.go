package views

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestFormatCommodityName(t *testing.T) {
	cases := map[string]string{
		"lumber":        "Kayu Bangunan",
		"stone":         "Batu Tambang",
		"rations":       "Ransum Kering",
		"iron_ore":      "Bijih Besi Murni",
		"herbal_salve":  "Salep Herbal",
		"spices":        "Rempah Lembah",
		"caravan_silk":  "Sutra Musafir",
		"ancient_relic": "Relik Kuno",
	}

	for id, expected := range cases {
		got := formatCommodityName(id)
		if got != expected {
			t.Errorf("expected %s -> %s, got %s", id, expected, got)
		}
	}
}

func TestRenderMarketView_Tab0_TableAlignment(t *testing.T) {
	eng := engine.NewGame("Petualang", "Oakhaven")

	rendered := RenderMarketView(eng, 0, 0, 0, 80)
	lines := strings.Split(rendered, "\n")

	var itemLines []string
	var headerLine string

	for _, l := range lines {
		clean := lipgloss.NewStyle().Render(l)
		if strings.Contains(clean, "Komoditas") && strings.Contains(clean, "Tren") {
			headerLine = clean
		}
		if strings.Contains(clean, "Kayu Bangunan") || strings.Contains(clean, "Batu Tambang") || strings.Contains(clean, "Relik Kuno") {
			itemLines = append(itemLines, clean)
		}
	}

	if headerLine == "" {
		t.Fatalf("expected to find header line in local market tab")
	}
	if len(itemLines) < 3 {
		t.Fatalf("expected at least 3 commodity lines, got %d", len(itemLines))
	}

	headerPipes := findVisualPipeIndices(headerLine)
	if len(headerPipes) < 5 {
		t.Fatalf("expected 5 column separators '|' in header, got %d", len(headerPipes))
	}

	for i, it := range itemLines {
		rowPipes := findVisualPipeIndices(it)
		if len(rowPipes) != len(headerPipes) {
			t.Errorf("item line %d has %d pipes, expected %d", i, len(rowPipes), len(headerPipes))
			continue
		}
		for pIdx := range headerPipes {
			if rowPipes[pIdx] != headerPipes[pIdx] {
				t.Errorf("item line %d pipe %d position mismatch: got col %d, expected col %d. Line: %q",
					i, pIdx, rowPipes[pIdx], headerPipes[pIdx], it)
			}
		}
	}
}

func TestRenderMarketView_Tab1_TableAlignment(t *testing.T) {
	eng := engine.NewGame("Petualang", "Oakhaven")
	eng.Village.Buildings[settlement.BuildingCaravanPost] = 2

	rendered := RenderMarketView(eng, 0, 1, 0, 80)
	lines := strings.Split(rendered, "\n")

	var routeLines []string
	var headerLine string

	for _, l := range lines {
		clean := lipgloss.NewStyle().Render(l)
		if strings.Contains(clean, "Rute Tujuan") && strings.Contains(clean, "Modal & Kargo") {
			headerLine = clean
		}
		if strings.Contains(clean, "Riverfall") || strings.Contains(clean, "Ironpeak") || strings.Contains(clean, "Dunhallow") {
			routeLines = append(routeLines, clean)
		}
	}

	if headerLine == "" {
		t.Fatalf("expected to find header line in caravan tab")
	}
	if len(routeLines) < 3 {
		t.Fatalf("expected at least 3 caravan routes, got %d", len(routeLines))
	}

	// Extract pipe positions in header
	headerPipes := findVisualPipeIndices(headerLine)
	if len(headerPipes) < 4 {
		t.Fatalf("expected 4 column separators '|' in header, got %d", len(headerPipes))
	}

	for i, r := range routeLines {
		rowPipes := findVisualPipeIndices(r)
		if len(rowPipes) != len(headerPipes) {
			t.Errorf("route line %d has %d pipes, expected %d", i, len(rowPipes), len(headerPipes))
			continue
		}
		for pIdx := range headerPipes {
			if rowPipes[pIdx] != headerPipes[pIdx] {
				t.Errorf("route line %d pipe %d position mismatch: got col %d, expected col %d. Line: %q",
					i, pIdx, rowPipes[pIdx], headerPipes[pIdx], r)
			}
		}
	}
}

func TestRenderMarketView_Tab1_DetailColons(t *testing.T) {
	eng := engine.NewGame("Petualang", "Oakhaven")
	eng.Village.Buildings[settlement.BuildingCaravanPost] = 2

	rendered := RenderMarketView(eng, 0, 1, 0, 80)
	lines := strings.Split(rendered, "\n")

	labels := []string{
		"Modal & Kargo",
		"Waktu Tempuh",
		"Estimasi Hasil",
		"Risiko Bandit",
		"Status Armada",
	}

	colonPositions := make(map[string]int)

	for _, l := range lines {
		for _, label := range labels {
			if strings.Contains(l, label) && strings.Contains(l, ":") {
				colons := findVisualColonIndices(l)
				if len(colons) > 0 {
					colonPositions[label] = colons[0]
				}
			}
		}
	}

	for _, label := range labels {
		if _, ok := colonPositions[label]; !ok {
			t.Errorf("missing detail label: %s", label)
		}
	}

	// Verify all colons are at the same visual column position
	var refPos int
	for _, label := range labels {
		pos := colonPositions[label]
		if refPos == 0 {
			refPos = pos
		} else if pos != refPos {
			t.Errorf("colon for %s at visual col %d, expected visual col %d", label, pos, refPos)
		}
	}
}

func findVisualPipeIndices(line string) []int {
	var indices []int
	var currentVisualCol int
	for _, r := range line {
		if r == '|' {
			indices = append(indices, currentVisualCol)
		}
		currentVisualCol += lipgloss.Width(string(r))
	}
	return indices
}

func findVisualColonIndices(line string) []int {
	var indices []int
	var currentVisualCol int
	for _, r := range line {
		if r == ':' {
			indices = append(indices, currentVisualCol)
		}
		currentVisualCol += lipgloss.Width(string(r))
	}
	return indices
}
