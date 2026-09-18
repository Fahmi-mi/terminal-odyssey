package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderVictoryView displays narrative triumph screen when endgame goals are met
func RenderVictoryView(e *engine.Engine, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	header := styles.TitleStyle.Render("TRIUMPH OF OAKHAVEN : KEMENANGAN AGUNG TERCAPAI")
	divider := strings.Repeat("─", contentWidth)

	triumphBanner := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorGold).
		Background(styles.ColorDarker).
		Padding(0, 1).
		Render("[ CAHAYA KEMENANGAN MENERANGI SELURUH LEMBAH OAKHAVEN ]")

	storyLines := []string{
		"Katakombe terdalam telah runtuh dalam keheningan abadi setelah kejatuhan Penguasa Abyssal",
		"Suara raungan monster yang dahulu meneror penduduk desa kini berganti dengan denting palu pandai besi,",
		"tawa ceria para pekerja di ladang, dan langkah tegap para penjaga benteng yang kokoh berdiri",
		"",
		fmt.Sprintf("Dari sebuah perkemahan darurat terpencil, %s telah menjelma menjadi benteng peradaban megah", e.Village.Name),
		"yang disegani oleh para kafilah dagang dari seluruh penjuru benua",
	}
	storyContent := styles.LogItemStyle.Render(strings.Join(storyLines, "\n"))

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	leftStats := []string{
		styles.SubtitleStyle.Render("[ PENCAPAIAN PERADABAN ]"),
		fmt.Sprintf("  Waktu Tempuh   : Hari ke-%d (%s)", e.DayCounter, e.CurrentSeason.String()),
		fmt.Sprintf("  Balai Desa     : Level %d", e.Village.Buildings[settlement.BuildingTownHall]),
		fmt.Sprintf("  Total Warga    : %d Jiwa", e.Village.Settlers),
		fmt.Sprintf("  Kas Emas Desa  : %d Gold", e.Village.Treasury),
	}

	rightStats := []string{
		styles.SubtitleStyle.Render("[ KEKUATAN & PERTAHANAN ]"),
		fmt.Sprintf("  Nilai Benteng  : %d Poin Pertahanan", e.Village.DefenseVal),
		fmt.Sprintf("  Serbuan Digagalkan: %d Kali Pertahanan Sukses", e.Siege.TotalSiegesRepelled),
		fmt.Sprintf("  Status Boss    : Penguasa Katakombe Ditaklukkan"),
		fmt.Sprintf("  Rombongan Aktif: %d Rekan Setia", len(e.Player.Party)),
	}

	leftCol := lipgloss.NewStyle().Width(halfWidth).Render(strings.Join(leftStats, "\n"))
	rightCol := lipgloss.NewStyle().Width(rightWidth).Render(strings.Join(rightStats, "\n"))
	statsGrid := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	actions := styles.KeyBadge.Render("Enter") + " Lanjut Bermain (Mode Sandbox Bebas)  " +
		styles.KeyBadge.Render("S") + " Simpan Permainan  " +
		styles.KeyBadge.Render("Q") + " Keluar ke Menu Utama"

	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		divider,
		triumphBanner,
		"",
		storyContent,
		"",
		divider,
		statsGrid,
		"",
		divider,
		actions,
	)

	return styles.ActiveBox.Width(boxWidth).Render(boxContent)
}
