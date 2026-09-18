package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/dungeon"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// formatRoomTypeBadge formats badges with distinctive colors and no padding
func formatRoomTypeBadge(roomType string) string {
	switch roomType {
	case dungeon.RoomTypeCombat:
		return lipgloss.NewStyle().Bold(true).Foreground(styles.ColorRed).Render("[Pertarungan]")
	case dungeon.RoomTypeTreasure:
		return lipgloss.NewStyle().Bold(true).Foreground(styles.ColorGold).Render("[Peti Karun]")
	case dungeon.RoomTypeRest:
		return lipgloss.NewStyle().Bold(true).Foreground(styles.ColorGreen).Render("[Suaka Istirahat]")
	case dungeon.RoomTypeMystery:
		return lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPurple).Render("[Misteri Aksara]")
	case dungeon.RoomTypeExit:
		return lipgloss.NewStyle().Bold(true).Foreground(styles.ColorEmerald).Render("[Tangga Keluar]")
	default:
		return "[Ruangan]"
	}
}

// renderDungeonHeader renders the unified player profile and equipment header
func renderDungeonHeader(exp *dungeon.Expedition, titleText string, boxWidth, contentWidth int) string {
	p := exp.Player
	w := p.EquippedWeapon
	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	mightBonus := p.Stats.Might - 10
	if mightBonus < 0 {
		mightBonus = 0
	}
	critPct := (w.CritRate + float64(p.Stats.Agility)*0.005) * 100
	totalInit := w.Initiative + p.Stats.Agility

	leftLines := []string{
		styles.SubtitleStyle.Render("[ PROFIL & ATRIBUT ]"),
		fmt.Sprintf("  Nama      : %s", styles.ResourceVal.Render(p.Name)),
		fmt.Sprintf("  Darah (HP): %s", styles.ResourceVal.Render(fmt.Sprintf("%d / %d", p.HP, p.MaxHP))),
		fmt.Sprintf("  Kewarasan : %s", styles.ResourceWood.Render(fmt.Sprintf("%d / %d", p.Sanity, p.MaxSanity))),
		fmt.Sprintf("  Might: %-2d (+%d ATK)  Agi: %-2d (+%d)", p.Stats.Might, mightBonus, p.Stats.Agility, p.Stats.Agility),
		fmt.Sprintf("  Resolve: %-2d (+%d HP)  Ing: %-2d", p.Stats.Resolve, (p.Stats.Resolve-10)*5, p.Stats.Ingenuity),
		fmt.Sprintf("  Rombongan : %d Rekan | %d Ramuan", len(p.Party), p.TotalPotions()),
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ PERLENGKAPAN & LOGISTIK ]"),
		fmt.Sprintf("  Senjata   : %s", styles.DefenseStyle.Render(w.Name)),
		fmt.Sprintf("  Tipe/ATK  : %s (%d-%d ATK)", w.WeaponType, w.BaseDamage[0]+mightBonus, w.BaseDamage[1]+mightBonus),
		fmt.Sprintf("  Kritikal  : %.1f%% | Init: %d", critPct, totalInit),
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
	totalDepths := exp.TotalDepths
	if totalDepths == 0 {
		totalDepths = len(exp.Rooms)
	}
	depth := room.Depth
	if depth == 0 {
		depth = room.Index
	}
	roomTitle := fmt.Sprintf("RUANG %d/%d: %s", depth, totalDepths, strings.ToUpper(room.Def.Title))
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
				contentLines = append(contentLines, styles.AlertWarning.Render("  Catatan: Anda sempat mundur, musuh masih terluka"))
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
		contentLines = append(contentLines, styles.ItemHighlight.Render("  Tekan ENTER untuk menyelesaikan ekspedisi dan membawa pulang seluruh jarahan"))
	}

	// Branching corridors display if room is resolved
	choices := exp.NextRoomChoices()
	if room.IsResolved && len(choices) > 1 {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, styles.SubtitleStyle.Render("[ PILIHAN PERCABANGAN LORONG ]"))
		for i, c := range choices {
			typeBadge := formatRoomTypeBadge(c.Def.Type)
			choiceLine := fmt.Sprintf("  [%d] %-12s : %s %s", i+1, c.BranchName, typeBadge, c.Def.Title)
			if lipgloss.Width(choiceLine) > contentWidth-4 {
				wrapped := styles.WrapText(choiceLine, contentWidth-4)
				choiceLine = wrapped[0]
			}
			contentLines = append(contentLines, choiceLine)
		}
	} else if room.IsResolved && len(choices) == 1 {
		contentLines = append(contentLines, "")
		nextRoom := choices[0]
		typeBadge := formatRoomTypeBadge(nextRoom.Def.Type)
		choiceLine := fmt.Sprintf("  Jalur berikutnya: %s %s", typeBadge, nextRoom.Def.Title)
		if lipgloss.Width(choiceLine) > contentWidth-4 {
			wrapped := styles.WrapText(choiceLine, contentWidth-4)
			choiceLine = wrapped[0]
		}
		contentLines = append(contentLines, choiceLine)
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
			wrapped := styles.WrapText(log, contentWidth-4)
			for _, wl := range wrapped {
				contentLines = append(contentLines, fmt.Sprintf("  %s", wl))
			}
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

	if room.Def.Type == dungeon.RoomTypeExit {
		controlActions = append(controlActions,
			fmt.Sprintf("%s Selesai & Bawa Jarahan", styles.KeyBadge.Render("ENTER")),
			fmt.Sprintf("%s Mundur", styles.KeyBadge.Render("ESC")),
		)
		if exp.Rations > 0 && exp.Player.HP < exp.Player.MaxHP {
			controlActions = append(controlActions, fmt.Sprintf("%s Makan", styles.KeyBadge.Render("M")))
		}
	} else if room.Def.Type == dungeon.RoomTypeCombat && !room.IsResolved {
		actionText := "Bertarung"
		if room.Enemy != nil && room.Enemy.HP < room.Enemy.MaxHP {
			actionText = "Serang Lagi"
		}
		controlActions = append(controlActions,
			fmt.Sprintf("%s %s", styles.KeyBadge.Render("ENTER"), actionText),
			fmt.Sprintf("%s Makan", styles.KeyBadge.Render("M")),
			fmt.Sprintf("%s Mundur ke Desa", styles.KeyBadge.Render("ESC")),
		)
	} else if !room.IsResolved {
		switch room.Def.Type {
		case dungeon.RoomTypeTreasure:
			controlActions = append(controlActions, fmt.Sprintf("%s Buka Peti", styles.KeyBadge.Render("ENTER")))
		case dungeon.RoomTypeRest:
			controlActions = append(controlActions, fmt.Sprintf("%s Beristirahat", styles.KeyBadge.Render("ENTER")))
		case dungeon.RoomTypeMystery:
			controlActions = append(controlActions, fmt.Sprintf("%s Teliti Altar", styles.KeyBadge.Render("ENTER")))
		}
		controlActions = append(controlActions,
			fmt.Sprintf("%s Makan", styles.KeyBadge.Render("M")),
			fmt.Sprintf("%s Obor", styles.KeyBadge.Render("O")),
		)
		if exp.Player.TotalPotions() > 0 {
			controlActions = append(controlActions, fmt.Sprintf("%s Ramuan", styles.KeyBadge.Render("P")))
		}
		controlActions = append(controlActions, fmt.Sprintf("%s Mundur", styles.KeyBadge.Render("ESC")))
	} else {
		// Room is resolved
		if len(choices) >= 2 {
			controlActions = append(controlActions,
				fmt.Sprintf("%s %s", styles.KeyBadge.Render("1"), choices[0].BranchName),
				fmt.Sprintf("%s %s", styles.KeyBadge.Render("2"), choices[1].BranchName),
			)
		} else {
			controlActions = append(controlActions, fmt.Sprintf("%s Lanjut Melangkah", styles.KeyBadge.Render("ENTER")))
		}
		controlActions = append(controlActions,
			fmt.Sprintf("%s Makan", styles.KeyBadge.Render("M")),
			fmt.Sprintf("%s Obor", styles.KeyBadge.Render("O")),
		)
		if exp.Player.TotalPotions() > 0 {
			controlActions = append(controlActions, fmt.Sprintf("%s Ramuan", styles.KeyBadge.Render("P")))
		}
		controlActions = append(controlActions, fmt.Sprintf("%s Mundur", styles.KeyBadge.Render("ESC")))
	}

	controlsText := strings.Join(controlActions, " | ")
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
