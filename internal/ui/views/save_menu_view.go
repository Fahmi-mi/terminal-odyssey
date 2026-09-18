package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/save"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderSaveMenuView displays slot selection for saving progress
func RenderSaveMenuView(e *engine.Engine, selectedIdx int, slots []save.SaveSlotInfo, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	header := styles.TitleStyle.Render("SIMPAN PROGRES PERMAINAN")
	sub := styles.SubtitleStyle.Render(fmt.Sprintf("Status Saat Ini: Hari ke-%d | Desa %s | %s", e.DayCounter, e.Village.Name, e.Player.Name))
	divider := strings.Repeat("─", contentWidth)

	var slotRows []string
	slotRows = append(slotRows, styles.SubtitleStyle.Render("PILIH SLOT TUJUAN:"))
	slotRows = append(slotRows, "")

	// Only show manual save slots 1, 2, 3
	manualSlots := make([]save.SaveSlotInfo, 0)
	for _, s := range slots {
		if s.SlotID != save.SlotAutosave {
			manualSlots = append(manualSlots, s)
		}
	}

	for i, s := range manualSlots {
		var slotTitle string
		switch s.SlotID {
		case save.Slot1:
			slotTitle = "Slot 1"
		case save.Slot2:
			slotTitle = "Slot 2"
		case save.Slot3:
			slotTitle = "Slot 3"
		default:
			slotTitle = s.SlotID
		}

		var detailText string
		if s.Exists {
			detailText = fmt.Sprintf("Hari ke-%d | %s (%s)", s.DayCounter, s.VillageName, s.SaveTime)
		} else {
			detailText = "[ Kosong - Siap Digunakan ]"
		}

		var row string
		if i == selectedIdx {
			row = styles.ItemHighlight.Render(fmt.Sprintf("▶ %-10s : %s", slotTitle, detailText))
		} else {
			row = styles.ItemNormal.Render(fmt.Sprintf("  %-10s : %s", slotTitle, detailText))
		}
		slotRows = append(slotRows, row)
	}

	slotContent := strings.Join(slotRows, "\n")

	navHelp := styles.KeyBadge.Render("Enter") + " Simpan ke Slot  " +
		styles.KeyBadge.Render("▲/▼") + " Pilih Slot  " +
		styles.KeyBadge.Render("Esc / B") + " Kembali ke Desa"

	var alertBanner string
	if e.StatusAlert != "" {
		alertBanner = styles.AlertSuccess.Render("[INFO] " + e.StatusAlert)
	}

	elements := []string{
		header,
		sub,
		divider,
		"",
		slotContent,
		"",
		divider,
		navHelp,
	}

	if alertBanner != "" {
		elements = append(elements, "", alertBanner)
	}

	boxContent := lipgloss.JoinVertical(lipgloss.Left, elements...)
	return styles.ActiveBox.Width(boxWidth).Render(boxContent)
}
