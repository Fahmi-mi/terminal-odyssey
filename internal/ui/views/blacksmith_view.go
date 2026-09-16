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

// RenderBlacksmithView renders the modular weapon crafting and armory screen
func RenderBlacksmithView(e *engine.Engine, selectedRecipeIdx int, activeTab int, selectedWeaponIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	bsLvl := e.Village.Buildings[settlement.BuildingBlacksmith]
	titleText := fmt.Sprintf("BENGKEL PANDAI BESI: TEMPA & GUDANG SENJATA (LEVEL %d)", bsLvl)
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
	mightBonus := e.Player.Stats.Might - 10
	if mightBonus < 0 {
		mightBonus = 0
	}
	critPct := (w.CritRate + float64(e.Player.Stats.Agility)*0.005) * 100
	totalInit := w.Initiative + e.Player.Stats.Agility

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
		fmt.Sprintf("  Stat/ATK : %d-%d ATK | %.1f%% Crit | Init %d", w.BaseDamage[0]+mightBonus, w.BaseDamage[1]+mightBonus, critPct, totalInit),
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

	recipes, err := data.LoadRecipeDefs()
	if err != nil {
		return styles.AlertError.Render("[!] Gagal membaca katalog resep senjata")
	}

	itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)
	var contentLines []string
	var controls []string

	if activeTab == 0 {
		// Tab 0: Katalog Resep Tempa
		tabHeader := styles.ItemHighlight.Render("[1] RESEP TEMPA") + "    " + styles.ItemNormal.Render("[2/TAB] GUDANG SENJATA ("+fmt.Sprintf("%d", len(e.Player.OwnedWeapons))+")")
		contentLines = append(contentLines, tabHeader)

		if selectedRecipeIdx < 0 {
			selectedRecipeIdx = 0
		}
		if selectedRecipeIdx >= len(recipes) {
			selectedRecipeIdx = len(recipes) - 1
		}

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
			if i == selectedRecipeIdx {
				cursor = "▶ "
			}

			dmgText := fmt.Sprintf("%d-%d ATK", r.MinDamage, r.MaxDamage)
			itemText := fmt.Sprintf("%s%-*s | %-7s | %-9s | Bengkel Lv.%d", cursor, nameWidth, r.Name, r.Type, dmgText, r.BlacksmithLevel)

			var renderedLine string
			if i == selectedRecipeIdx {
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
		sel := recipes[selectedRecipeIdx]
		contentLines = append(contentLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN: %s ]", strings.ToUpper(sel.Name))))

		descWrapped := styles.WrapText(sel.Description, contentWidth-4)
		for _, dw := range descWrapped {
			contentLines = append(contentLines, styles.LogItemStyle.Render(dw))
		}

		selCritPct := (sel.CritRate + float64(e.Player.Stats.Agility)*0.005) * 100
		selInit := sel.Initiative + e.Player.Stats.Agility
		statLine := fmt.Sprintf("  Tipe: %-7s | ATK: %d-%d (+%d Might) | Crit: %.1f%% | Init: %d",
			sel.Type, sel.MinDamage+mightBonus, sel.MaxDamage+mightBonus, mightBonus, selCritPct, selInit)
		contentLines = append(contentLines, statLine)

		if sel.SpecialAffix != "" {
			contentLines = append(contentLines, fmt.Sprintf("  Efek Khusus : %s", styles.ItemHighlight.Render(sel.SpecialAffix)))
		}

		costLine := fmt.Sprintf("  Kebutuhan   : Lvl %d Bengkel | %d Kayu | %d Batu | %d Gold",
			sel.BlacksmithLevel, sel.WoodCost, sel.StoneCost, sel.GoldCost)
		contentLines = append(contentLines, costLine)

		// Craftability / ownership status
		if e.Player.OwnsWeapon(sel.ID) {
			contentLines = append(contentLines, styles.AlertSuccess.Render("  Status      : Dimiliki di Gudang (Buka tab [TAB] Gudang untuk Pasang)"))
		} else if bsLvl < sel.BlacksmithLevel {
			contentLines = append(contentLines, styles.AlertError.Render(fmt.Sprintf("  Status      : Terkunci (Butuh Bengkel Pandai Besi Level %d)", sel.BlacksmithLevel)))
		} else if e.Village.Lumber < sel.WoodCost || e.Village.Stone < sel.StoneCost || e.Village.Treasury < sel.GoldCost {
			contentLines = append(contentLines, styles.AlertWarning.Render("  Status      : Material atau Kas Emas belum mencukupi"))
		} else {
			contentLines = append(contentLines, styles.AlertSuccess.Render("  Status      : Siap Ditempa"))
		}

		controls = []string{
			fmt.Sprintf("%s Gudang", styles.KeyBadge.Render("TAB")),
			fmt.Sprintf("%s Pilih", styles.KeyBadge.Render("↑/↓")),
			fmt.Sprintf("%s Tempa", styles.KeyBadge.Render("ENTER")),
			fmt.Sprintf("%s Reparasi", styles.KeyBadge.Render("R")),
			fmt.Sprintf("%s Keluar", styles.KeyBadge.Render("ESC")),
		}
	} else {
		// Tab 1: Gudang Senjata (Armory)
		tabHeader := styles.ItemNormal.Render("[1/TAB] RESEP TEMPA") + "    " + styles.ItemHighlight.Render("[2] GUDANG SENJATA ("+fmt.Sprintf("%d", len(e.Player.OwnedWeapons))+")")
		contentLines = append(contentLines, tabHeader)

		owned := e.Player.OwnedWeapons
		if selectedWeaponIdx < 0 {
			selectedWeaponIdx = 0
		}
		if selectedWeaponIdx >= len(owned) {
			selectedWeaponIdx = len(owned) - 1
		}

		for i, ow := range owned {
			cursor := "  "
			if i == selectedWeaponIdx {
				cursor = "▶ "
			}

			isEquipped := ow.ID == e.Player.EquippedWeapon.ID
			statusTag := "SIMPAN "
			if isEquipped {
				statusTag = "DIPAKAI"
			}

			dmgText := fmt.Sprintf("%2d-%2d ATK", ow.BaseDamage[0]+mightBonus, ow.BaseDamage[1]+mightBonus)
			duraText := fmt.Sprintf("%2d/%2d Dur", ow.Durability, ow.MaxDura)

			displayName := ow.Name
			if len(displayName) > 23 {
				displayName = displayName[:20] + "..."
			}

			lineText := fmt.Sprintf("%s%-23s | %-7s | %-9s | %s | %s", cursor, displayName, ow.WeaponType, dmgText, duraText, statusTag)

			var renderedLine string
			if i == selectedWeaponIdx {
				renderedLine = styles.ItemHighlight.Render(lineText)
			} else if isEquipped {
				renderedLine = itemAvailable.Render(lineText)
			} else {
				renderedLine = styles.ItemNormal.Render(lineText)
			}

			contentLines = append(contentLines, renderedLine)
		}

		contentLines = append(contentLines, "")

		// Detail of selected owned weapon
		sel := owned[selectedWeaponIdx]
		contentLines = append(contentLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN: %s ]", strings.ToUpper(sel.Name))))

		desc := "Senjata perlengkapan milik petualang"
		for _, r := range recipes {
			if r.ID == sel.ID {
				desc = r.Description
				break
			}
		}
		if sel.ID == "rusty_sword" {
			desc = "Pedang tua peninggalan masa lampau dengan bilah bergerigi karat namun tetap kokoh"
		}

		descWrapped := styles.WrapText(desc, contentWidth-4)
		for _, dw := range descWrapped {
			contentLines = append(contentLines, styles.LogItemStyle.Render(dw))
		}

		selCritPct := (sel.CritRate + float64(e.Player.Stats.Agility)*0.005) * 100
		selInit := sel.Initiative + e.Player.Stats.Agility
		statLine := fmt.Sprintf("  Tipe: %-7s | ATK: %d-%d (+%d Might) | Crit: %.1f%% | Init: %d",
			sel.WeaponType, sel.BaseDamage[0]+mightBonus, sel.BaseDamage[1]+mightBonus, mightBonus, selCritPct, selInit)
		contentLines = append(contentLines, statLine)

		if sel.SpecialAffix != "" {
			contentLines = append(contentLines, fmt.Sprintf("  Efek Khusus : %s", styles.ItemHighlight.Render(sel.SpecialAffix)))
		}

		condLine := fmt.Sprintf("  Kondisi     : %d/%d Durabilitas", sel.Durability, sel.MaxDura)
		contentLines = append(contentLines, condLine)

		if sel.ID == e.Player.EquippedWeapon.ID {
			contentLines = append(contentLines, styles.AlertSuccess.Render("  Status      : Senjata Aktif Sedang Digunakan"))
		} else {
			contentLines = append(contentLines, styles.AlertWarning.Render("  Status      : Tersimpan di Gudang (Tekan [ENTER] untuk Memasang Senjata Ini)"))
		}

		controls = []string{
			fmt.Sprintf("%s Tempa", styles.KeyBadge.Render("TAB")),
			fmt.Sprintf("%s Pilih", styles.KeyBadge.Render("↑/↓")),
			fmt.Sprintf("%s Pasang", styles.KeyBadge.Render("ENTER")),
			fmt.Sprintf("%s Reparasi", styles.KeyBadge.Render("R")),
			fmt.Sprintf("%s Keluar", styles.KeyBadge.Render("ESC")),
		}
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			contentLines...,
		),
	)

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
