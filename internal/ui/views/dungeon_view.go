package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/dungeon"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// renderDungeonHeader renders the unified player profile and equipment header
func renderDungeonHeader(exp *dungeon.Expedition, titleText string, boxWidth, contentWidth int) string {
	p := exp.Player
	w := p.EquippedWeapon
	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	leftLines := []string{
		styles.SubtitleStyle.Render("[ PROFIL & ATRIBUT ]"),
		fmt.Sprintf("  Nama      : %s", styles.ResourceVal.Render(p.Name)),
		fmt.Sprintf("  Darah (HP): %s", styles.ResourceVal.Render(fmt.Sprintf("%d / %d", p.HP, p.MaxHP))),
		fmt.Sprintf("  Kewarasan : %s", styles.ResourceWood.Render(fmt.Sprintf("%d / %d", p.Sanity, p.MaxSanity))),
		fmt.Sprintf("  Might     : %-3d  Agility  : %d", p.Stats.Might, p.Stats.Agility),
		fmt.Sprintf("  Resolve   : %-3d  Ingenuity: %d", p.Stats.Resolve, p.Stats.Ingenuity),
		fmt.Sprintf("  Kapasitas : %d Slot Ransel", p.MaxBackpack),
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ PERLENGKAPAN & LOGISTIK ]"),
		fmt.Sprintf("  Senjata   : %s", styles.DefenseStyle.Render(w.Name)),
		fmt.Sprintf("  Tipe/ATK  : %s (%d-%d ATK)", w.WeaponType, w.BaseDamage[0], w.BaseDamage[1]),
		fmt.Sprintf("  Kritikal  : %.0f%% | Inisiatif: %d", w.CritRate*100, w.Initiative),
		fmt.Sprintf("  Ketahanan : %d/%d (%s)", w.Durability, w.MaxDura, w.SpecialAffix),
		fmt.Sprintf("  Obor/Bekal: %s | %s",
			styles.ResourceWood.Render(fmt.Sprintf("%d%%", exp.Torch)),
			styles.ResourceFood.Render(fmt.Sprintf("%d Ransum", exp.Rations))),
		fmt.Sprintf("  Kas Temuan: %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", exp.GoldFound))),
	}

	leftCol := lipgloss.NewStyle().Width(halfWidth).Render(strings.Join(leftLines, "\n"))
	rightCol := lipgloss.NewStyle().Width(rightWidth).Render(strings.Join(rightLines, "\n"))
	playerGrid := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	return styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			divider,
			playerGrid,
		),
	)
}

// RenderDungeonView renders the dungeon room exploration screen
func RenderDungeonView(e *engine.Engine, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	exp := e.ActiveExpedition
	if exp == nil {
		return styles.AlertError.Render("[!] Tidak ada ekspedisi yang sedang aktif")
	}

	room := exp.CurrentRoom()

	// 1. Header Box (Unified with player profile, stats, weapon, and supplies)
	roomTitle := fmt.Sprintf("RUANG %d/%d: %s", room.Index, len(exp.Rooms), strings.ToUpper(room.Def.Title))
	headerBox := renderDungeonHeader(exp, roomTitle, boxWidth, contentWidth)

	// 2. Room Content Box
	var contentLines []string

	// Atmosphere Narrative (wrapped to fit contentWidth)
	descLines := styles.WrapText(room.Def.Description, contentWidth-4)
	for _, dl := range descLines {
		contentLines = append(contentLines, styles.LogItemStyle.Render(dl))
	}
	contentLines = append(contentLines, "")

	// Room status details
	switch room.Def.Type {
	case dungeon.RoomTypeCombat:
		contentLines = append(contentLines, styles.SubtitleStyle.Render("[ INTEL RUANGAN ]"))
		if room.IsResolved {
			contentLines = append(contentLines, styles.ResourceLabel.Render("  Status : Musuh di ruangan ini telah ditumpas habis"))
		} else if room.Enemy != nil {
			contentLines = append(contentLines, fmt.Sprintf("  Musuh  : %s (HP %d/%d, %d-%d ATK)",
				styles.AlertError.Render(room.Enemy.Name), room.Enemy.HP, room.Enemy.MaxHP, room.Enemy.MinDamage, room.Enemy.MaxDamage))
			if room.Enemy.HP < room.Enemy.MaxHP {
				contentLines = append(contentLines, styles.AlertWarning.Render("  Catatan: Anda sempat mundur, musuh masih terluka dan siap diserang lagi"))
			} else {
				contentLines = append(contentLines, styles.ResourceLabel.Render(fmt.Sprintf("  Bahaya : %s", room.Enemy.Description)))
			}
		}
	case dungeon.RoomTypeTreasure:
		contentLines = append(contentLines, styles.SubtitleStyle.Render("[ PETI TEMUAN ]"))
		if room.IsResolved {
			contentLines = append(contentLines, styles.ResourceLabel.Render(fmt.Sprintf("  Status : %s", room.ResolutionLog)))
		} else {
			contentLines = append(contentLines, styles.ItemHighlight.Render("  Peti kayu berbalut besi tempa menanti untuk dibuka"))
		}
	case dungeon.RoomTypeRest:
		contentLines = append(contentLines, styles.SubtitleStyle.Render("[ SUAKA AMAN ]"))
		if room.IsResolved {
			contentLines = append(contentLines, styles.ResourceLabel.Render(fmt.Sprintf("  Status : %s", room.ResolutionLog)))
		} else {
			contentLines = append(contentLines, styles.ItemHighlight.Render("  Mata air jernih dan sisa perapian hangat (+35% HP, +30% Obor)"))
		}
	case dungeon.RoomTypeMystery:
		contentLines = append(contentLines, styles.SubtitleStyle.Render("[ MISTERI AKSARA ]"))
		if room.IsResolved {
			contentLines = append(contentLines, styles.ResourceLabel.Render(fmt.Sprintf("  Status : %s", room.ResolutionLog)))
		} else {
			contentLines = append(contentLines, styles.ItemHighlight.Render("  Altar obsidian berukir rune (Uji Kecerdikan / Ingenuity)"))
		}
	case dungeon.RoomTypeExit:
		contentLines = append(contentLines, styles.SubtitleStyle.Render("[ TANGGA KELUAR ]"))
		contentLines = append(contentLines, styles.AlertSuccess.Render("  Jalur evakuasi menuju permukaan desa terbuka lebar"))
	}

	contentLines = append(contentLines, "")

	// Expedition Travel Logs
	contentLines = append(contentLines, styles.SubtitleStyle.Render("[ CATATAN PERJALANAN ]"))
	if len(exp.Logs) == 0 {
		contentLines = append(contentLines, styles.ResourceLabel.Render("  Belum ada catatan baru"))
	} else {
		startIdx := 0
		if len(exp.Logs) > 3 {
			startIdx = len(exp.Logs) - 3
		}
		for _, log := range exp.Logs[startIdx:] {
			contentLines = append(contentLines, fmt.Sprintf("  %s", log))
		}
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			contentLines...,
		),
	)

	// 3. Controls / Action Guide Box
	var controlActions []string

	if room.Def.Type == dungeon.RoomTypeCombat && !room.IsResolved {
		actionText := "Bertarung"
		if room.Enemy != nil && room.Enemy.HP < room.Enemy.MaxHP {
			actionText = "Serang Lagi"
		}
		controlActions = append(controlActions,
			fmt.Sprintf("%s %s", styles.KeyBadge.Render("ENTER"), actionText),
			fmt.Sprintf("%s Makan", styles.KeyBadge.Render("M")),
			fmt.Sprintf("%s Mundur ke Desa", styles.KeyBadge.Render("ESC")),
		)
	} else {
		// Room specific interaction
		if !room.IsResolved {
			switch room.Def.Type {
			case dungeon.RoomTypeTreasure:
				controlActions = append(controlActions, fmt.Sprintf("%s Buka Peti", styles.KeyBadge.Render("ENTER")))
			case dungeon.RoomTypeRest:
				controlActions = append(controlActions, fmt.Sprintf("%s Beristirahat", styles.KeyBadge.Render("ENTER")))
			case dungeon.RoomTypeMystery:
				controlActions = append(controlActions, fmt.Sprintf("%s Teliti Altar", styles.KeyBadge.Render("ENTER")))
			}
		}

		if room.Def.Type == dungeon.RoomTypeExit {
			controlActions = append(controlActions, fmt.Sprintf("%s Selesai & Bawa Jarahan", styles.KeyBadge.Render("ENTER")))
		} else {
			controlActions = append(controlActions, fmt.Sprintf("%s Lanjut Melangkah", styles.KeyBadge.Render("SPACE")))
		}

		controlActions = append(controlActions,
			fmt.Sprintf("%s Makan", styles.KeyBadge.Render("M")),
			fmt.Sprintf("%s Obor", styles.KeyBadge.Render("O")),
			fmt.Sprintf("%s Mundur", styles.KeyBadge.Render("ESC")),
		)
	}

	controlsText := strings.Join(controlActions, "  |  ")
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
