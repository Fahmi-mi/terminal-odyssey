package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

type statOption struct {
	Name        string
	DisplayName string
	Description string
	Benefit     string
	Value       int
}

// RenderTrainingView renders the character stat progression and training screen
func RenderTrainingView(e *engine.Engine, selectedIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	tgLvl := e.Village.Buildings[settlement.BuildingTrainingGround]
	capVal := 10 + (tgLvl * 5)
	p := e.Player

	titleText := fmt.Sprintf("PUSAT LATIHAN KARAKTER (LEVEL %d)", tgLvl)
	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	mightBonus := p.Stats.Might - 10
	if mightBonus < 0 {
		mightBonus = 0
	}

	// Header columns
	leftLines := []string{
		styles.SubtitleStyle.Render("[ LOGISTIK & FASILITAS ]"),
		fmt.Sprintf("  Kas Emas   : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", e.Village.Treasury))),
		fmt.Sprintf("  Ransum     : %s", styles.ResourceFood.Render(fmt.Sprintf("%d Ransum", e.Village.Rations))),
		fmt.Sprintf("  Kapasitas  : %d Slot Ransel", p.MaxBackpack),
		fmt.Sprintf("  Batas Stat : %s", styles.DefenseStyle.Render(fmt.Sprintf("Maks %d (Lvl %d)", capVal, tgLvl))),
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ EFEK ATRIBUT PETUALANG ]"),
		fmt.Sprintf("  Nama       : %s", styles.ResourceVal.Render(p.Name)),
		fmt.Sprintf("  Darah HP   : %s (+%d HP)", styles.ResourceVal.Render(fmt.Sprintf("%d/%d", p.HP, p.MaxHP)), (p.Stats.Resolve-10)*5),
		fmt.Sprintf("  Bonus ATK  : +%d Serangan Fisik", mightBonus),
		fmt.Sprintf("  Duel/Crit  : +%d Init | +%.1f%% Crit", p.Stats.Agility, float64(p.Stats.Agility)*0.5),
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

	// Stat options
	stats := []statOption{
		{
			Name:        "Might",
			DisplayName: "Might (Kekuatan Fisik)",
			Description: "Meningkatkan daya hancur serangan fisik senjata dan menambah slot ransel bawaan",
			Benefit:     "+1 Bonus Damage Serangan Fisik, +Slot Ransel tiap tingkatan",
			Value:       p.Stats.Might,
		},
		{
			Name:        "Agility",
			DisplayName: "Agility (Kelincahan)",
			Description: "Meningkatkan inisiatif kecepatan bertarung, peluang kritikal, dan peluang kabur",
			Benefit:     "+1 Inisiatif Duel, +0.5% Peluang Serangan Kritikal",
			Value:       p.Stats.Agility,
		},
		{
			Name:        "Resolve",
			DisplayName: "Resolve (Ketahanan Mental)",
			Description: "Meningkatkan kapasitas darah maksimal karakter dan ketahanan terhadap teror mental",
			Benefit:     "+5 Batas Maksimal Darah (HP), Ketahanan Stres",
			Value:       p.Stats.Resolve,
		},
		{
			Name:        "Ingenuity",
			DisplayName: "Ingenuity (Kecerdikan)",
			Description: "Memperbesar peluang memecahkan teka-teki aksara altar kuno dan tawar-menawar pasar",
			Benefit:     "+Peluang Berhasil Meneliti Altar Misteri & Diskon Pasar",
			Value:       p.Stats.Ingenuity,
		},
	}

	if selectedIdx < 0 {
		selectedIdx = 0
	}
	if selectedIdx >= len(stats) {
		selectedIdx = len(stats) - 1
	}

	var contentLines []string
	contentLines = append(contentLines, styles.SubtitleStyle.Render("[ DAFTAR ATRIBUT PELATIHAN ]"))

	itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)

	for i, s := range stats {
		goldCost := s.Value * 5
		rationCost := 2 + (s.Value - 10)
		if rationCost < 2 {
			rationCost = 2
		}

		isAtCap := s.Value >= capVal
		canAfford := e.Village.Treasury >= goldCost && e.Village.Rations >= rationCost && !isAtCap

		cursor := "  "
		if i == selectedIdx {
			cursor = "▶ "
		}

		itemText := fmt.Sprintf("%s%-26s : Nilai %2d -> %-2d | Biaya: %2d Gold, %d Ransum",
			cursor, s.DisplayName, s.Value, s.Value+1, goldCost, rationCost)
		if isAtCap {
			itemText = fmt.Sprintf("%s%-26s : Nilai %2d        | Batas Maksimum (Lv.%d)", cursor, s.DisplayName, s.Value, tgLvl)
		}

		var renderedLine string
		if i == selectedIdx {
			renderedLine = styles.ItemHighlight.Render(itemText)
		} else if canAfford {
			renderedLine = itemAvailable.Render(itemText)
		} else {
			renderedLine = styles.ItemNormal.Render(itemText)
		}

		contentLines = append(contentLines, renderedLine)
	}

	contentLines = append(contentLines, "")

	// Detail of selected stat
	sel := stats[selectedIdx]
	contentLines = append(contentLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN: %s ]", strings.ToUpper(sel.Name))))

	descWrapped := styles.WrapText(sel.Description, contentWidth-4)
	for _, dw := range descWrapped {
		contentLines = append(contentLines, styles.LogItemStyle.Render(dw))
	}

	contentLines = append(contentLines, fmt.Sprintf("  Manfaat : %s", styles.ItemHighlight.Render(sel.Benefit)))

	goldCost := sel.Value * 5
	rationCost := 2 + (sel.Value - 10)
	if rationCost < 2 {
		rationCost = 2
	}

	isAtCap := sel.Value >= capVal
	if isAtCap {
		contentLines = append(contentLines, styles.AlertError.Render(fmt.Sprintf("  Status  : Atribut mencapai batas Pusat Latihan Level %d (Tingkatkan fasilitas di desa)", tgLvl)))
	} else if e.Village.Treasury < goldCost || e.Village.Rations < rationCost {
		contentLines = append(contentLines, styles.AlertWarning.Render(fmt.Sprintf("  Status  : Sumber daya kurang (Butuh %d Gold, %d Ransum)", goldCost, rationCost)))
	} else {
		contentLines = append(contentLines, styles.AlertSuccess.Render(fmt.Sprintf("  Status  : Siap Dilatih (Biaya: %d Gold, %d Ransum)", goldCost, rationCost)))
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
		fmt.Sprintf("%s Latih Atribut", styles.KeyBadge.Render("ENTER")),
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
