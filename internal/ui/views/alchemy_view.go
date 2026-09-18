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

// RenderAlchemyView renders the apothecary brewing laboratory screen
func RenderAlchemyView(e *engine.Engine, selectedRecipeIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	labLvl := e.Village.Buildings[settlement.BuildingApothecary]
	titleText := fmt.Sprintf("LABORATORIUM ALKIMIA: RACIK RAMUAN & MINYAK (LEVEL %d)", labLvl)
	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	// 1. Header Box: Materials & Pouch
	herbs := e.Village.GetCommodityStock("herbal_salve")
	beastBlood := e.Village.GetCommodityStock("beast_blood")
	shadowEss := e.Village.GetCommodityStock("shadow_essence")

	leftLines := []string{
		styles.SubtitleStyle.Render("[ BAHAN BAKU PEMUKIMAN ]"),
		fmt.Sprintf("  Kas Emas       : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", e.Village.Treasury))),
		fmt.Sprintf("  Herba Salve    : %s", styles.ResourceWood.Render(fmt.Sprintf("%d Unit", herbs))),
		fmt.Sprintf("  Darah Monster  : %s", styles.ResourceStone.Render(fmt.Sprintf("%d Unit", beastBlood))),
		fmt.Sprintf("  Esensi Bayangan: %s", styles.ResourceVal.Render(fmt.Sprintf("%d Unit", shadowEss))),
	}

	p := e.Player
	ingDiff := p.Stats.Ingenuity - 10
	critChance := 0
	if ingDiff > 0 {
		critChance = ingDiff * 4
		if critChance > 35 {
			critChance = 35
		}
	}

	salepCnt := p.GetPotionCount("salep_pemulih")
	oborCnt := p.GetPotionCount("minyak_obor")
	racunCnt := p.GetPotionCount("penawar_racun")
	tonikCnt := p.GetPotionCount("tonik_penenang")
	eliksirCnt := p.GetPotionCount("eliksir_kekuatan")

	rightLines := []string{
		styles.SubtitleStyle.Render("[ KANTONG RAMUAN PETUALANG ]"),
		fmt.Sprintf("  Salep / Obor   : %d Salep | %d Obor", salepCnt, oborCnt),
		fmt.Sprintf("  Penawar/Tonik  : %d Penawar | %d Tonik", racunCnt, tonikCnt),
		fmt.Sprintf("  Eliksir Tempur : %d Eliksir", eliksirCnt),
		fmt.Sprintf("  Bonus Racikan  : %d Ing (+%d%% Ganda)", p.Stats.Ingenuity, critChance),
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

	// 2. Recipe Catalog
	recipes := e.Alchemy.Recipes
	if selectedRecipeIdx < 0 {
		selectedRecipeIdx = 0
	}
	if selectedRecipeIdx >= len(recipes) {
		selectedRecipeIdx = len(recipes) - 1
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorCyan).Padding(0, 1)
	dividerStyle := lipgloss.NewStyle().Foreground(styles.ColorDim).Padding(0, 1)

	tableHeader := fmt.Sprintf("  %-20s | %-9s | %-20s | %-10s",
		"NAMA RESEP RAMUAN", "FASILITAS", "BIAYA & BAHAN", "STATUS")
	tableDivider := strings.Repeat("─", 70)
	recipeLines := []string{
		styles.SubtitleStyle.Render("[ KATALOG RESEP ALKIMIA ]"),
		"",
		headerStyle.Render(tableHeader),
		dividerStyle.Render(tableDivider),
	}

	itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)

	for i, r := range recipes {
		cursor := "  "
		isSelected := i == selectedRecipeIdx
		if isSelected {
			cursor = "▶ "
		}

		lvlReq := fmt.Sprintf("Lab Lv.%d", r.ApothecaryLevel)
		costStr := formatRecipeIngredients(r)
		if len(costStr) > 20 {
			costStr = costStr[:17] + "..."
		}

		canBrewErr := e.Alchemy.CanBrew(r.ID, labLvl, e.Village.Treasury, e.Village.Commodities, e.Village.Lumber, e.Village.Stone, e.Village.Rations)
		statusStr := "Siap Racik"
		if labLvl < r.ApothecaryLevel {
			statusStr = "Terkunci"
		} else if canBrewErr != nil {
			statusStr = "Kurang"
		}

		nameStr := r.Name
		if len(nameStr) > 20 {
			nameStr = nameStr[:17] + "..."
		}

		if len(statusStr) > 10 {
			statusStr = statusStr[:10]
		}

		rowText := fmt.Sprintf("%s%-20s | %-9s | %-20s | %-10s",
			cursor, nameStr, lvlReq, costStr, statusStr)

		if isSelected {
			recipeLines = append(recipeLines, styles.ItemHighlight.Render(rowText))
		} else if labLvl < r.ApothecaryLevel {
			recipeLines = append(recipeLines, styles.ItemNormal.Render(rowText))
		} else if canBrewErr == nil {
			recipeLines = append(recipeLines, itemAvailable.Render(rowText))
		} else {
			recipeLines = append(recipeLines, styles.ItemNormal.Render(rowText))
		}
	}

	catalogContent := lipgloss.JoinVertical(lipgloss.Left, recipeLines...)
	catalogBox := styles.ActiveBox.Width(boxWidth).Render(catalogContent)

	// 3. Selected Recipe Details Box
	var detailLines []string
	if selectedRecipeIdx >= 0 && selectedRecipeIdx < len(recipes) {
		sel := recipes[selectedRecipeIdx]
		effectStr := ""
		switch sel.EffectType {
		case "heal_hp":
			effectStr = fmt.Sprintf("Memulihkan +%d HP seketika", sel.EffectValue)
		case "refill_torch":
			effectStr = fmt.Sprintf("Menambah +%d%% cahaya obor eksplorasi", sel.EffectValue)
		case "cure_debuff":
			effectStr = fmt.Sprintf("Menetralkan racun & pendarahan (+%d HP)", sel.EffectValue)
		case "restore_sanity":
			effectStr = fmt.Sprintf("Memulihkan +%d Sanity kewarasan", sel.EffectValue)
		case "buff_atk":
			effectStr = fmt.Sprintf("Meningkatkan Serangan +%d ATK saat pertempuran", sel.EffectValue)
		default:
			effectStr = fmt.Sprintf("%s (%d)", sel.EffectType, sel.EffectValue)
		}

		detailLines = append(detailLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN: %s ]", strings.ToUpper(sel.Name))))

		descLines := styles.WrapText(sel.Description, contentWidth-16)
		if len(descLines) > 0 {
			detailLines = append(detailLines, fmt.Sprintf("  Deskripsi : %s", styles.ResourceLabel.Render(descLines[0])))
			for _, dl := range descLines[1:] {
				detailLines = append(detailLines, fmt.Sprintf("              %s", styles.ResourceLabel.Render(dl)))
			}
		}

		detailLines = append(detailLines, fmt.Sprintf("  Khasiat   : %s", styles.ResourceVal.Render(effectStr)))
		detailLines = append(detailLines, fmt.Sprintf("  Kebutuhan : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", sel.GoldCost)))+formatDetailedRequirements(e, sel))

		canBrewErr := e.Alchemy.CanBrew(sel.ID, labLvl, e.Village.Treasury, e.Village.Commodities, e.Village.Lumber, e.Village.Stone, e.Village.Rations)
		if labLvl < sel.ApothecaryLevel {
			detailLines = append(detailLines, fmt.Sprintf("  Status    : %s", styles.AlertError.Render(fmt.Sprintf("Laboratorium Alkimia Level %d diperlukan", sel.ApothecaryLevel))))
		} else if canBrewErr != nil {
			detailLines = append(detailLines, fmt.Sprintf("  Status    : %s", styles.AlertWarning.Render(canBrewErr.Error())))
		} else {
			detailLines = append(detailLines, fmt.Sprintf("  Status    : %s", styles.AlertSuccess.Render("Bahan lengkap, siap diproduksi")))
		}
	}

	detailBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left, detailLines...),
	)

	// 4. Alert Banner (if any)
	var alertBox string
	if e.StatusAlert != "" {
		alertBox = styles.AlertSuccess.Render("[INFO] " + e.StatusAlert)
	}

	// 5. Bottom Controls Box
	controls := fmt.Sprintf("%s Pilih Resep   %s Racik Ramuan   %s Kembali ke Desa",
		styles.KeyBadgeHighlight.Render("↑/↓"),
		styles.KeyBadgeHighlight.Render("ENTER"),
		styles.KeyBadge.Render("ESC"),
	)
	bottomBox := styles.BaseBox.Width(boxWidth).Render(controls)

	elements := []string{headerBox, catalogBox, detailBox}
	if alertBox != "" {
		elements = append(elements, alertBox)
	}
	elements = append(elements, bottomBox)

	return lipgloss.JoinVertical(lipgloss.Left, elements...)
}

func formatRecipeIngredients(r data.AlchemyRecipeDef) string {
	parts := []string{fmt.Sprintf("%d G", r.GoldCost)}
	for id, amt := range r.Ingredients {
		displayName := id
		switch id {
		case "herbal_salve":
			displayName = "Herba"
		case "beast_blood":
			displayName = "Darah"
		case "shadow_essence":
			displayName = "Esensi"
		case "stone":
			displayName = "Batu"
		case "spices":
			displayName = "Rempah"
		case "iron_ore":
			displayName = "Besi"
		case "lumber":
			displayName = "Kayu"
		case "rations":
			displayName = "Ransum"
		}
		parts = append(parts, fmt.Sprintf("%d %s", amt, displayName))
	}
	return strings.Join(parts, ", ")
}

func formatDetailedRequirements(e *engine.Engine, r data.AlchemyRecipeDef) string {
	var parts []string
	for id, amt := range r.Ingredients {
		displayName := id
		avail := 0
		switch id {
		case "herbal_salve":
			displayName = "Herba Salve"
			avail = e.Village.GetCommodityStock(id)
		case "beast_blood":
			displayName = "Darah Monster"
			avail = e.Village.GetCommodityStock(id)
		case "shadow_essence":
			displayName = "Esensi Bayangan"
			avail = e.Village.GetCommodityStock(id)
		case "stone":
			displayName = "Batu"
			avail = e.Village.Stone
		case "spices":
			displayName = "Rempah Aromatik"
			avail = e.Village.GetCommodityStock(id)
		case "iron_ore":
			displayName = "Biji Besi"
			avail = e.Village.GetCommodityStock(id)
		case "lumber":
			displayName = "Kayu"
			avail = e.Village.Lumber
		case "rations":
			displayName = "Ransum"
			avail = e.Village.Rations
		}

		parts = append(parts, fmt.Sprintf("%s (%d/%d)", displayName, avail, amt))
	}
	if len(parts) == 0 {
		return ""
	}
	return " | " + strings.Join(parts, ", ")
}
