package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/save"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// TitleMenuMode differentiates between main title screen and load slot selector
type TitleMenuMode int

const (
	TitleModeMain TitleMenuMode = iota
	TitleModeSelectSlot
)

// RenderTitleView renders the start screen and save slot loader
func RenderTitleView(mode TitleMenuMode, selectedIdx int, slots []save.SaveSlotInfo, alert string, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	// ASCII banner
	bannerLines := []string{
		" _____ _____ ____  __  __ ___ _   _    _    _     ",
		"|_   _| ____|  _ \\|  \\/  |_ _| \\ | |  / \\  | |    ",
		"  | | |  _| | |_) | |\\/| || ||  \\| | / _ \\ | |    ",
		"  | | | |___|  _ <| |  | || || |\\  |/ ___ \\| |___ ",
		"  |_| |_____|_| \\_\\_|  |_|___|_| \\_/_/   \\_\\_____|",
		"  ___  ______   ______ ____ _____   __            ",
		" / _ \\|  _ \\ \\ / / ___/ ___| ____\\ \\ / /            ",
		"| | | | | | \\ V /\\___ \\___ \\  _|  \\ V /             ",
		"| |_| | |_| || |  ___) |__) | |___ | |              ",
		" \\___/|____/ |_| |____/____/_____| |_|              ",
	}

	styledBannerLines := make([]string, len(bannerLines))
	for i, bl := range bannerLines {
		pad := (contentWidth - lipgloss.Width(bl)) / 2
		if pad < 0 {
			pad = 0
		}
		styledBannerLines[i] = strings.Repeat(" ", pad) + styles.TitleStyle.Render(bl)
	}
	banner := strings.Join(styledBannerLines, "\n")

	subtitle := styles.SubtitleStyle.Render("SIMULASI PEMUKIMAN & EKSPLORASI KATAKOMBE ROGUELIKE")
	subPad := (contentWidth - lipgloss.Width(subtitle)) / 2
	if subPad < 0 {
		subPad = 0
	}
	subHeader := strings.Repeat(" ", subPad) + subtitle

	divider := strings.Repeat("─", contentWidth)

	var menuContent string
	var helpText string
	if mode == TitleModeMain {
		options := []string{
			"Mulai Petualangan Baru",
			"Muat Permainan Tersimpan",
			"Keluar ke Terminal",
		}

		var menuRows []string
		for i, opt := range options {
			var cursor, label string
			if i == selectedIdx {
				cursor = styles.DefenseStyle.Render("▶ ")
				label = styles.ItemHighlight.Render(fmt.Sprintf("[%d] %s", i+1, opt))
			} else {
				cursor = "  "
				label = styles.ItemNormal.Render(fmt.Sprintf("[%d] %s", i+1, opt))
			}
			rowPad := (contentWidth - lipgloss.Width(cursor+label)) / 2
			if rowPad < 0 {
				rowPad = 0
			}
			menuRows = append(menuRows, strings.Repeat(" ", rowPad)+cursor+label)
		}
		menuContent = strings.Join(menuRows, "\n\n")
		helpText = styles.ResourceLabel.Render("Gunakan tombol [▲/▼] atau [1-3] untuk memilih, [Enter] untuk konfirmasi")
	} else {
		// Slot Selection Mode
		var slotRows []string
		slotRows = append(slotRows, styles.SubtitleStyle.Render("PILIH SLOT SIMPANAN UNTUK DIMUAT:"))
		slotRows = append(slotRows, "")

		for i, s := range slots {
			var slotBadge string
			switch s.SlotID {
			case save.Slot1:
				slotBadge = "[1] Slot 1"
			case save.Slot2:
				slotBadge = "[2] Slot 2"
			case save.Slot3:
				slotBadge = "[3] Slot 3"
			case save.SlotAutosave:
				slotBadge = "[4] Autosave"
			default:
				slotBadge = fmt.Sprintf("[%d] %s", i+1, s.SlotID)
			}

			var detailText string
			if s.Exists {
				bossStatus := ""
				if s.BossDefeated {
					bossStatus = " | Boss Ditaklukkan"
				}
				detailText = fmt.Sprintf("Hari ke-%d | Desa %s (%s)%s", s.DayCounter, s.VillageName, s.SaveTime, bossStatus)
			} else {
				detailText = "[ Slot Kosong ]"
			}

			lineText := fmt.Sprintf("%-16s : %s", slotBadge, detailText)
			var row string
			if i == selectedIdx {
				row = styles.ItemHighlight.Render("▶ " + lineText)
			} else {
				row = styles.ItemNormal.Render("  " + lineText)
			}
			slotRows = append(slotRows, row)
		}

		slotRows = append(slotRows, "")
		slotRows = append(slotRows, styles.KeyBadge.Render("Enter")+" Muat Slot  "+styles.KeyBadge.Render("Esc")+" Kembali ke Menu Utama")
		menuContent = strings.Join(slotRows, "\n")
		helpText = styles.ResourceLabel.Render("Gunakan tombol [▲/▼] atau [1-4] untuk memilih slot, [Enter] untuk muat")
	}

	helpPad := (contentWidth - lipgloss.Width(helpText)) / 2
	if helpPad < 0 {
		helpPad = 0
	}
	centeredHelp := strings.Repeat(" ", helpPad) + helpText

	var alertBox string
	if alert != "" {
		alertBox = styles.AlertWarning.Render("[INFO] " + alert)
	}

	elements := []string{
		banner,
		"",
		subHeader,
		divider,
		"",
		menuContent,
		"",
		divider,
		centeredHelp,
	}
	if alertBox != "" {
		elements = append(elements, "", alertBox)
	}

	boxContent := lipgloss.JoinVertical(lipgloss.Left, elements...)
	return styles.ActiveBox.Width(boxWidth).Render(boxContent)
}
