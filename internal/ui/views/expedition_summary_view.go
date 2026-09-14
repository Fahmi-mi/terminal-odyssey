package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderExpeditionSummaryView renders the post-expedition debrief screen
func RenderExpeditionSummaryView(e *engine.Engine, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}

	summary := e.LastExpeditionSummary
	if summary == nil {
		return styles.AlertError.Render("[!] Tidak ada ringkasan ekspedisi yang tersedia")
	}

	// 1. Header Box
	title := styles.TitleStyle.Render("LAPORAN RINGKASAN EKSPEDISI KATAKOMBE")
	headerBox := styles.BaseBox.Width(boxWidth).Render(title)

	// 2. Summary Content Box
	var lines []string

	if summary.WasEvacuated {
		lines = append(lines, styles.AlertSuccess.Render("[ EKSPEDISI SELESAI: BERHASIL DIEVAKUASI ]"))
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("  Emas Disetor ke Kas Desa  : %s", styles.ResourceGold.Render(fmt.Sprintf("+%d Gold", summary.GoldEarned))))
		lines = append(lines, fmt.Sprintf("  Kayu Disimpan di Gudang   : %s", styles.ResourceWood.Render(fmt.Sprintf("+%d Kayu", summary.LumberEarned))))
		lines = append(lines, fmt.Sprintf("  Batu Disimpan di Gudang   : %s", styles.ResourceStone.Render(fmt.Sprintf("+%d Batu", summary.StoneEarned))))
		lines = append(lines, fmt.Sprintf("  Musuh Berhasil Ditumpas   : %s", styles.ResourceVal.Render(fmt.Sprintf("%d Monster", summary.EnemiesDefeated))))
		lines = append(lines, fmt.Sprintf("  Ruangan Dijelajahi        : %s", styles.ResourceVal.Render(fmt.Sprintf("%d / %d Ruang", summary.RoomsExplored, summary.TotalRooms))))
		lines = append(lines, "")
		lines = append(lines, styles.ResourceLabel.Render("  Seluruh sumber daya jarahan berhasil diamankan ke kas dan lumbung desa"))
	} else {
		lines = append(lines, styles.AlertError.Render("[ EKSPEDISI GAGAL: DIEVAKUASI DALAM KONDISI KRITIS ]"))
		lines = append(lines, "")
		lines = append(lines, styles.ResourceLabel.Render("  Seluruh jarahan emas dan perbekalan di dalam ransel hilang"))
		lines = append(lines, styles.ResourceLabel.Render("  Warga mengevakuasi tubuh Anda ke balai desa untuk pemulihan darurat"))
		lines = append(lines, styles.ResourceLabel.Render("  1 hari telah berlalu di desa selama proses perawatan medis"))
		lines = append(lines, fmt.Sprintf("  Kondisi Terkini           : Siuman dengan %s", styles.ResourceVal.Render(fmt.Sprintf("%d HP", e.Player.HP))))
		lines = append(lines, fmt.Sprintf("  Musuh Sempat Ditumpas     : %s", styles.ResourceVal.Render(fmt.Sprintf("%d Monster", summary.EnemiesDefeated))))
		lines = append(lines, fmt.Sprintf("  Ruangan Terakhir Dicapai  : %s", styles.ResourceVal.Render(fmt.Sprintf("Ruang %d / %d", summary.RoomsExplored, summary.TotalRooms))))
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lines...,
		),
	)

	// 3. Bottom Controls Box
	controls := fmt.Sprintf("%s Kembali ke Pusat Desa",
		styles.KeyBadgeHighlight.Render("ENTER"))
	bottomBox := styles.BaseBox.Width(boxWidth).Render(controls)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBox,
		contentBox,
		bottomBox,
	)
}
