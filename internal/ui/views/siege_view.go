package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderSiegeView displays dramatic battle report when raiders assault the village
func RenderSiegeView(e *engine.Engine, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	header := styles.TitleStyle.Render("LAPORAN SERBUAN GERBANG PEMUKIMAN")
	daySub := styles.DayCounterStyle.Render(fmt.Sprintf("Hari ke-%d | Lokasi: %s", e.DayCounter, e.Village.Name))
	divider := strings.Repeat("─", contentWidth)

	res := e.Siege.LastResult
	if res == nil {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			daySub,
			divider,
			"Tidak ada catatan pertempuran terkini",
			divider,
			styles.KeyBadge.Render("Enter")+" Kembali ke Desa",
		)
		return styles.ActiveBox.Width(boxWidth).Render(content)
	}

	var statusBanner string
	if res.Victory {
		statusBanner = lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorEmerald).
			Background(styles.ColorDarker).
			Padding(0, 1).
			Render("[ KEMENANGAN BERTAHAN : SERBUAN MUSUH BERHASIL DIPUKUL MUNDUR ]")
	} else {
		statusBanner = lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorCrimson).
			Background(styles.ColorDarker).
			Padding(0, 1).
			Render("[ PERTAHANAN JEBOL : PENYERBU MENEROBOS GERBANG DESA ]")
	}

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	leftLines := []string{
		styles.SubtitleStyle.Render("[ KEKUATAN PENYERBU ]"),
		fmt.Sprintf("  Kelompok Musuh : %s", styles.ResourceVal.Render(res.RaiderName)),
		fmt.Sprintf("  Kekuatan Serbu : %s Poin", styles.AlertError.Render(fmt.Sprintf("%d", res.AssaultPower))),
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ PERTAHANAN GERBANG ]"),
		fmt.Sprintf("  Garda & Benteng: %s", styles.DefenseStyle.Render(fmt.Sprintf("%d Poin", res.DefensePower))),
		fmt.Sprintf("  Status Gerbang : %s", map[bool]string{true: styles.ResourceFood.Render("Kokoh Bertahan"), false: styles.AlertError.Render("Tembus Terjarah")}[res.Victory]),
	}

	leftCol := lipgloss.NewStyle().Width(halfWidth).Render(strings.Join(leftLines, "\n"))
	rightCol := lipgloss.NewStyle().Width(rightWidth).Render(strings.Join(rightLines, "\n"))
	grid := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	var impactLines []string
	impactLines = append(impactLines, styles.SubtitleStyle.Render("[ DAMPAK & KONSEKUENSI PERTEMPURAN ]"))

	if res.Victory {
		impactLines = append(impactLines, fmt.Sprintf("  Rampasan Kas Emas : +%s Gold", styles.ResourceGold.Render(fmt.Sprintf("%d", res.GoldDifference))))
		if res.SettlersGained > 0 {
			impactLines = append(impactLines, fmt.Sprintf("  Tawanan Warga     : +%d Orang (Menetap sebagai penduduk baru)", res.SettlersGained))
		}
		impactLines = append(impactLines, "  Kondisi Warga     : Seluruh penduduk dan lumbung perbekalan selamat")
	} else {
		impactLines = append(impactLines, fmt.Sprintf("  Emas Terjarah     : %s Gold", styles.AlertError.Render(fmt.Sprintf("%d", res.GoldDifference))))
		if res.RationsLost > 0 {
			impactLines = append(impactLines, fmt.Sprintf("  Ransum Dirampas   : -%d Pasokan", res.RationsLost))
		}
		if res.LumberLost > 0 || res.StoneLost > 0 {
			impactLines = append(impactLines, fmt.Sprintf("  Material Gudang   : -%d Kayu, -%d Batu", res.LumberLost, res.StoneLost))
		}
		if res.SettlersLost > 0 {
			impactLines = append(impactLines, fmt.Sprintf("  Korban Jiwa       : %s Penduduk gugur dalam penyerbuan", styles.AlertError.Render(fmt.Sprintf("%d", res.SettlersLost))))
		}
	}
	impactContent := strings.Join(impactLines, "\n")

	narration := styles.LogItemStyle.Render(res.Log)

	navHelp := styles.KeyBadge.Render("Enter / Spasi") + " Melanjutkan Langkah ke Desa"

	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		daySub,
		divider,
		statusBanner,
		"",
		grid,
		"",
		divider,
		impactContent,
		"",
		narration,
		divider,
		navHelp,
	)

	return styles.ActiveBox.Width(boxWidth).Render(boxContent)
}
