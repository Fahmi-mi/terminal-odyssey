package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderTownView renders the central village hub screen
func RenderTownView(e *engine.Engine, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	v := e.Village
	maxL, maxS, maxR := v.StorageCap()

	// 1. Header Box (Settlement Info & Resources)
	title := styles.TitleStyle.Render(fmt.Sprintf("DESA: %s", strings.ToUpper(v.Name)))
	dayStr := styles.DayCounterStyle.Render(fmt.Sprintf("[HARI KE-%d]", e.DayCounter))
	seasonStr := styles.SeasonStyle.Render(fmt.Sprintf("Musim: %s", e.CurrentSeason.String()))

	leftHeader := fmt.Sprintf("%s  %s", title, dayStr)
	rightHeader := seasonStr
	headerSpacing := contentWidth - lipgloss.Width(leftHeader) - lipgloss.Width(rightHeader)
	if headerSpacing < 2 {
		headerSpacing = 2
	}
	headerLine := leftHeader + strings.Repeat(" ", headerSpacing) + rightHeader

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	leftLines := []string{
		styles.SubtitleStyle.Render("[ WARGA & PEKERJA ]"),
		fmt.Sprintf("  Populasi     : %s (%d Nganggur)", styles.ResourceVal.Render(fmt.Sprintf("%d/%d", v.Settlers, v.MaxSettlers())), v.UnassignedSettlers()),
		fmt.Sprintf("  Petani       : %s", styles.ResourceFood.Render(fmt.Sprintf("%d Orang", v.Workers.Farmers))),
		fmt.Sprintf("  Penebang     : %s", styles.ResourceWood.Render(fmt.Sprintf("%d Orang", v.Workers.Lumberjacks))),
		fmt.Sprintf("  Penambang    : %s", styles.ResourceStone.Render(fmt.Sprintf("%d Orang", v.Workers.Miners))),
		fmt.Sprintf("  Pandai Besi  : %s", styles.ResourceVal.Render(fmt.Sprintf("%d Orang", v.Workers.Blacksmiths))),
		fmt.Sprintf("  Garda        : %s", styles.DefenseStyle.Render(fmt.Sprintf("%d Orang", v.Workers.Militia))),
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ LOGISTIK & PERTAHANAN ]"),
		fmt.Sprintf("  Pertahanan   : %s", styles.DefenseStyle.Render(fmt.Sprintf("%d Poin", v.DefenseVal))),
		fmt.Sprintf("  Kas Emas     : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", v.Treasury))),
		fmt.Sprintf("  Kayu         : %s", styles.ResourceWood.Render(fmt.Sprintf("%d / %d", v.Lumber, maxL))),
		fmt.Sprintf("  Batu         : %s", styles.ResourceStone.Render(fmt.Sprintf("%d / %d", v.Stone, maxS))),
		fmt.Sprintf("  Ransum       : %s", styles.ResourceFood.Render(fmt.Sprintf("%d / %d", v.Rations, maxR))),
		"",
	}

	leftCol := lipgloss.NewStyle().Width(halfWidth).Render(strings.Join(leftLines, "\n"))
	rightCol := lipgloss.NewStyle().Width(rightWidth).Render(strings.Join(rightLines, "\n"))
	detailGrid := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	divider := strings.Repeat("─", contentWidth)

	headerBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerLine,
			divider,
			detailGrid,
		),
	)

	// 2. Daily Log Box
	logHeader := styles.SubtitleStyle.Render("LAPORAN HARIAN:")
	var logLines []string
	if len(e.DailyLogs) == 0 {
		logLines = append(logLines, styles.LogItemStyle.Render("Tidak ada peristiwa penting hari ini"))
	} else {
		for _, log := range e.DailyLogs {
			logLines = append(logLines, styles.LogItemStyle.Render(log))
		}
	}

	logContent := lipgloss.JoinVertical(lipgloss.Left, logLines...)
	logBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			logHeader,
			logContent,
		),
	)

	// 3. Action Menu Box
	actionHeader := styles.TitleStyle.Render("TINDAKAN TERSEDIA:")

	actions := []string{
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("1"), "Pasar & Perdagangan Komoditas"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("2"), "Pembangunan & Peningkatan Fasilitas"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("3"), "Bengkel Pandai Besi (Tempa Senjata & Zirah)"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("4"), "Pusat Latihan (Tingkatkan Stat Karakter)"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("5"), "Kedai Minum (Rekrut Pendamping & Rumor)"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("6"), "Siapkan Ekspedisi Katakombe Bawah Tanah"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("W"), "Alokasi & Penugasan Pekerja"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("D"), "Lewati Hari (Jalankan Siklus Produksi Harian)"),
		fmt.Sprintf("%s %s", styles.KeyBadge.Render("Q"), "Keluar Permainan"),
	}

	actionBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			actionHeader,
			strings.Join(actions, "\n"),
		),
	)

	// 4. Alert Banner (if any)
	var alertBox string
	if e.StatusAlert != "" {
		alertBox = styles.AlertSuccess.Render("[INFO] " + e.StatusAlert)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBox,
		logBox,
		actionBox,
		alertBox,
	)
}
