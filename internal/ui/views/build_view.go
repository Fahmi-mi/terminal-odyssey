package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderBuildView renders the building and infrastructure upgrade screen
func RenderBuildView(e *engine.Engine, selectedIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2
	v := e.Village
	buildings := settlement.AllBuildingsList()

	title := styles.TitleStyle.Render("PEMBANGUNAN & PENINGKATAN FASILITAS DESA")
	resourcesSummary := fmt.Sprintf("Kayu: %s   •   Batu: %s   •   Emas: %s",
		styles.ResourceWood.Render(fmt.Sprintf("%d", v.Lumber)),
		styles.ResourceStone.Render(fmt.Sprintf("%d", v.Stone)),
		styles.ResourceGold.Render(fmt.Sprintf("%d Gold", v.Treasury)),
	)
	divider := strings.Repeat("─", contentWidth)

	headerBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			divider,
			resourcesSummary,
		),
	)

	// Buildings List with spacing between items
	var lines []string
	for i, bName := range buildings {
		info := v.GetBuildingInfo(bName)
		cursor := "  "
		isSelected := i == selectedIdx
		if isSelected {
			cursor = "▶ "
		}

		statusStr := fmt.Sprintf("Lv.%d/%d", info.Level, info.MaxLevel)
		if info.Level >= info.MaxLevel {
			statusStr = "MAKS"
		}

		costStr := fmt.Sprintf("%d Kayu, %d Batu, %d Emas", info.WoodCost, info.StoneCost, info.GoldCost)
		if info.Level >= info.MaxLevel {
			costStr = "Level maksimum tercapai"
		}

		line1 := fmt.Sprintf("%s[%d] %-20s %-8s  Biaya: %s", cursor, i+1, bName, statusStr, costStr)
		line2 := fmt.Sprintf("      ↳ %s", info.Description)

		if isSelected {
			lines = append(lines, styles.ItemHighlight.Render(line1))
			lines = append(lines, styles.ItemHighlight.Render(line2))
		} else {
			lines = append(lines, styles.ItemNormal.Render(line1))
			lines = append(lines, styles.ItemNormal.Render(line2))
		}

		if i < len(buildings)-1 {
			lines = append(lines, "")
		}
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lines...,
		),
	)

	// Controls guide
	controls := fmt.Sprintf("%s/%s Pilih  |  %s Tingkatkan  |  %s Kembali",
		styles.KeyBadge.Render("↑/↓"),
		styles.KeyBadge.Render("1-7"),
		styles.KeyBadge.Render("ENTER"),
		styles.KeyBadge.Render("ESC"),
	)

	bottomBox := styles.BaseBox.Width(boxWidth).Render(controls)

	var alertBox string
	if e.StatusAlert != "" {
		alertBox = styles.AlertWarning.Render("[INFO] " + e.StatusAlert)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBox,
		contentBox,
		bottomBox,
		alertBox,
	)
}
