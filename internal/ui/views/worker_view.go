package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

type WorkerRoleItem struct {
	Role settlement.WorkerRole
	Desc string
}

var RoleItems = []WorkerRoleItem{
	{Role: settlement.RoleFarmer, Desc: "Memanen ransum pangan agar warga tidak kelaparan"},
	{Role: settlement.RoleLumberjack, Desc: "Memasok kayu bangunan & bahan bakar musim dingin"},
	{Role: settlement.RoleMiner, Desc: "Menggali batu benteng & fasilitas pertahanan"},
	{Role: settlement.RoleBlacksmith, Desc: "Menempa senjata tempur (Butuh Bengkel Lv.1)"},
	{Role: settlement.RoleMilitia, Desc: "Menjaga benteng dari bandit (+15 Pertahanan)"},
}

// RenderWorkerView renders the worker management screen
func RenderWorkerView(e *engine.Engine, selectedIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	v := e.Village

	title := styles.TitleStyle.Render("ALOKASI & PENUGASAN TENAGA KERJA")
	unassignedStr := styles.ResourceVal.Render(fmt.Sprintf("%d Orang", v.UnassignedSettlers()))
	popSummary := fmt.Sprintf("Warga Menganggur: %s | Total Populasi: %d/%d", unassignedStr, v.Settlers, v.MaxSettlers())
	divider := strings.Repeat("─", contentWidth)

	headerBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			divider,
			popSummary,
		),
	)

	// Roles List
	var lines []string
	for i, item := range RoleItems {
		var count int
		switch item.Role {
		case settlement.RoleFarmer:
			count = v.Workers.Farmers
		case settlement.RoleLumberjack:
			count = v.Workers.Lumberjacks
		case settlement.RoleMiner:
			count = v.Workers.Miners
		case settlement.RoleBlacksmith:
			count = v.Workers.Blacksmiths
		case settlement.RoleMilitia:
			count = v.Workers.Militia
		}

		cursor := "  "
		isSelected := i == selectedIdx
		if isSelected {
			cursor = "▶ "
		}

		roleName := fmt.Sprintf("%-12s", string(item.Role))
		countBadge := fmt.Sprintf("[%d]", count)

		rowText := fmt.Sprintf("%s%s %-4s  -- %s", cursor, roleName, countBadge, item.Desc)

		if isSelected {
			lines = append(lines, styles.ItemHighlight.Render(rowText))
		} else {
			lines = append(lines, styles.ItemNormal.Render(rowText))
		}
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lines...,
		),
	)

	// Controls guide
	controls := fmt.Sprintf("%s Pilih  |  %s Tambah (+1)  |  %s Kurangi (-1)  |  %s Kembali",
		styles.KeyBadge.Render("↑/↓"),
		styles.KeyBadge.Render("→/+"),
		styles.KeyBadge.Render("←/-"),
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
