package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderBlacksmithView renders the modular weapon crafting and repair screen
func RenderBlacksmithView(e *engine.Engine, selectedIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	bsLvl := e.Village.Buildings[settlement.BuildingBlacksmith]
	titleText := fmt.Sprintf("BENGKEL PANDAI BESI: TEMPA SENJATA (LEVEL %d)", bsLvl)
	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	// Header columns
	leftLines := []string{
		styles.SubtitleStyle.Render("[ LOGISTIK & MATERIAL DESA ]"),
		fmt.Sprintf("  Kas Emas : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", e.Village.Treasury))),
		fmt.Sprintf("  Kayu     : %s", styles.ResourceWood.Render(fmt.Sprintf("%d Kayu", e.Village.Lumber))),
		fmt.Sprintf("  Batu     : %s", styles.ResourceStone.Render(fmt.Sprintf("%d Batu", e.Village.Stone))),
	}

	w := e.Player.EquippedWeapon
	repairCostGold := ((w.MaxDura - w.Durability) * 2) / 5
	if repairCostGold < 5 && w.Durability < w.MaxDura {
		repairCostGold = 5
	}
	repairCostStone := ((w.MaxDura - w.Durability) * 1) / 5
	if repairCostStone < 2 && w.Durability < w.MaxDura {
		repairCostStone = 2
	}

	repairStatus := styles.AlertSuccess.Render("Maksimal (100%)")
	if w.Durability < w.MaxDura {
		repairStatus = styles.AlertWarning.Render(fmt.Sprintf("Butuh %d G, %d Batu", repairCostGold, repairCostStone))
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ SENJATA TERPASANG ]"),
		fmt.Sprintf("  Nama     : %s (%s)", styles.DefenseStyle.Render(w.Name), w.WeaponType),
		fmt.Sprintf("  Stat/ATK : %d-%d ATK | %.0f%% Crit | Init %d", w.BaseDamage[0], w.BaseDamage[1], w.CritRate*100, w.Initiative),
		fmt.Sprintf("  Kondisi  : %d/%d Durabilitas (%s)", w.Durability, w.MaxDura, repairStatus),
	}

	leftCol := lipgloss.NewStyle().Width(halfWidth).Render(strings.Join(leftLines, "\n"))
	rightCol := lipgloss.NewStyle().Width(rightWidth).Render(strings.Join(rightLines, "\n"))
	headerGrid := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	headerBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			divider,
			headerGrid,
		),
	)

	// Content Box: Recipe List
	recipes, err := data.LoadRecipeDefs()
	if err != nil {
		return styles.AlertError.Render("[!] Gagal membaca katalog resep senjata")
	}

	if selectedIdx < 0 {
		selectedIdx = 0
	}
	if selectedIdx >= len(recipes) {
		selectedIdx = len(recipes) - 1
	}

	var contentLines []string
	contentLines = append(contentLines, styles.SubtitleStyle.Render("[ KATALOG RESEP TEMPA ]"))

	itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)

	nameWidth := 28
	for _, r := range recipes {
		if len(r.Name) > nameWidth {
			nameWidth = len(r.Name)
		}
	}

	for i, r := range recipes {
		canCraft := bsLvl >= r.BlacksmithLevel &&
			e.Village.Lumber >= r.WoodCost &&
			e.Village.Stone >= r.StoneCost &&
			e.Village.Treasury >= r.GoldCost

		cursor := "  "
		if i == selectedIdx {
			cursor = "> "
		}

		dmgText := fmt.Sprintf("%d-%d ATK", r.MinDamage, r.MaxDamage)
		itemText := fmt.Sprintf("%s%-*s | %-7s | %-9s | Bengkel Lv.%d", cursor, nameWidth, r.Name, r.Type, dmgText, r.BlacksmithLevel)

		var renderedLine string
		if i == selectedIdx {
			renderedLine = styles.ItemHighlight.Render(itemText)
		} else if canCraft {
			renderedLine = itemAvailable.Render(itemText)
		} else {
			renderedLine = styles.ItemNormal.Render(itemText)
		}

		contentLines = append(contentLines, renderedLine)
	}

	contentLines = append(contentLines, "")

	// Detail of selected recipe
	sel := recipes[selectedIdx]
	contentLines = append(contentLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN: %s ]", strings.ToUpper(sel.Name))))

	descWrapped := styles.WrapText(sel.Description, contentWidth-4)
	for _, dw := range descWrapped {
		contentLines = append(contentLines, styles.LogItemStyle.Render(dw))
	}

	statLine := fmt.Sprintf("  Tipe: %-7s | ATK: %d-%d | Crit: %.0f%% | Inisiatif: %d | Durabilitas: %d",
		sel.Type, sel.MinDamage, sel.MaxDamage, sel.CritRate*100, sel.Initiative, sel.Durability)
	contentLines = append(contentLines, statLine)

	if sel.SpecialAffix != "" {
		contentLines = append(contentLines, fmt.Sprintf("  Efek Khusus : %s", styles.ItemHighlight.Render(sel.SpecialAffix)))
	}

	costLine := fmt.Sprintf("  Kebutuhan   : Lvl %d Bengkel | %d Kayu | %d Batu | %d Gold",
		sel.BlacksmithLevel, sel.WoodCost, sel.StoneCost, sel.GoldCost)
	contentLines = append(contentLines, costLine)

	// Craftability status
	if bsLvl < sel.BlacksmithLevel {
		contentLines = append(contentLines, styles.AlertError.Render(fmt.Sprintf("  Status      : Terkunci (Butuh Bengkel Pandai Besi Level %d)", sel.BlacksmithLevel)))
	} else if e.Village.Lumber < sel.WoodCost || e.Village.Stone < sel.StoneCost || e.Village.Treasury < sel.GoldCost {
		contentLines = append(contentLines, styles.AlertWarning.Render("  Status      : Material atau Kas Emas belum mencukupi"))
	} else {
		contentLines = append(contentLines, styles.AlertSuccess.Render("  Status      : Siap Ditempa"))
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			contentLines...,
		),
	)

	// Bottom Box (Controls)
	controls := []string{
		fmt.Sprintf("%s Pilih", styles.KeyBadge.Render("↑/↓")),
		fmt.Sprintf("%s Tempa", styles.KeyBadge.Render("ENTER")),
		fmt.Sprintf("%s Reparasi", styles.KeyBadge.Render("R")),
		fmt.Sprintf("%s Kembali", styles.KeyBadge.Render("ESC")),
	}

	controlsText := strings.Join(controls, " | ")
	bottomBox := styles.BaseBox.Width(boxWidth).Render(controlsText)

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
