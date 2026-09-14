package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderCombatView renders the tactical turn-based combat screen
func RenderCombatView(e *engine.Engine, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	exp := e.ActiveExpedition
	if exp == nil || exp.ActiveCombat == nil {
		return styles.AlertError.Render("[!] Tidak ada sesi pertarungan aktif")
	}

	session := exp.ActiveCombat
	room := exp.CurrentRoom()

	// 1. Header Box (Unified with player profile, stats, weapon, and supplies)
	combatTitle := fmt.Sprintf("PERTEMPURAN TAKTIS - RUANG %d: %s (GILIRAN %d)", room.Index, strings.ToUpper(room.Def.Title), session.TurnCount)
	if lipgloss.Width(combatTitle) > contentWidth {
		combatTitle = combatTitle[:contentWidth]
	}
	headerBox := renderDungeonHeader(exp, combatTitle, boxWidth, contentWidth)

	// 2. Battle Arena Content Box
	var arenaLines []string
	divider := strings.Repeat("─", contentWidth)

	// Enemy Status
	arenaLines = append(arenaLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ MUSUH: %s ]", strings.ToUpper(session.Enemy.Name))))
	arenaLines = append(arenaLines, fmt.Sprintf("  Status HP  : %s", styles.RenderHPBar(session.Enemy.HP, session.Enemy.MaxHP, 16)))
	arenaLines = append(arenaLines, fmt.Sprintf("  Daya Serang: %d - %d ATK   |   Pertahanan: %d Poin   |   Inisiatif: %d",
		session.Enemy.MinDamage, session.Enemy.MaxDamage, session.Enemy.Defense, session.Enemy.Initiative))
	if session.Enemy.Description != "" {
		descLines := styles.WrapText(fmt.Sprintf("Ciri/Bahaya: %s", session.Enemy.Description), contentWidth-4)
		for _, dl := range descLines {
			arenaLines = append(arenaLines, fmt.Sprintf("  %s", styles.ResourceLabel.Render(dl)))
		}
	}
	arenaLines = append(arenaLines, divider)

	// Combat Logs (latest 4 entries)
	arenaLines = append(arenaLines, styles.SubtitleStyle.Render("[ LOG PERTEMPURAN ]"))
	if len(session.Logs) == 0 {
		arenaLines = append(arenaLines, styles.ResourceLabel.Render("  Pertarungan dimulai, tentukan tindakan Anda"))
	} else {
		startIdx := 0
		if len(session.Logs) > 4 {
			startIdx = len(session.Logs) - 4
		}
		for _, log := range session.Logs[startIdx:] {
			wrapped := styles.WrapText(log, contentWidth-4)
			for _, wl := range wrapped {
				arenaLines = append(arenaLines, fmt.Sprintf("  %s", wl))
			}
		}
	}

	arenaBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			arenaLines...,
		),
	)

	// 3. Combat Controls Box
	var controlsText string
	if !session.IsOver {
		controlsText = fmt.Sprintf("%s Serang  |  %s Bertahan  |  %s Makan Ransum  |  %s Kabur",
			styles.KeyBadge.Render("1"),
			styles.KeyBadge.Render("2"),
			styles.KeyBadge.Render("3"),
			styles.KeyBadge.Render("4"),
		)
	} else if session.Won {
		controlsText = fmt.Sprintf("%s Menang! Lanjutkan Eksplorasi Ruangan",
			styles.KeyBadgeHighlight.Render("ENTER"))
	} else if session.Fled {
		controlsText = fmt.Sprintf("%s Berhasil Kabur! Kembali ke Ruangan",
			styles.KeyBadgeHighlight.Render("ENTER"))
	} else {
		controlsText = fmt.Sprintf("%s Karakter Tumbang! Dievakuasi Kembali ke Desa",
			styles.KeyBadgeHighlight.Render("ENTER"))
	}

	bottomBox := styles.BaseBox.Width(boxWidth).Render(controlsText)

	var alertBox string
	if e.StatusAlert != "" {
		alertBox = styles.AlertWarning.Render("[INFO] " + e.StatusAlert)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBox,
		arenaBox,
		bottomBox,
		alertBox,
	)
}
