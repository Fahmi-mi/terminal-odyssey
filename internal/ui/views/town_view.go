package views

import (
	"fmt"
	"sort"
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
		fmt.Sprintf("  Populasi     : %s (%d Luang)", styles.ResourceVal.Render(fmt.Sprintf("%d/%d", v.Settlers, v.MaxSettlers())), v.UnassignedSettlers()),
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

	p := e.Player
	w := p.EquippedWeapon

	playerName := p.Name
	maxPlayerNameLen := halfWidth - 14
	if lipgloss.Width(playerName) > maxPlayerNameLen && maxPlayerNameLen > 3 {
		playerName = playerName[:maxPlayerNameLen-3] + "..."
	}

	mightBonus := p.Stats.Might - 10
	if mightBonus < 0 {
		mightBonus = 0
	}
	critPct := (w.CritRate + float64(p.Stats.Agility)*0.005) * 100
	totalInit := w.Initiative + p.Stats.Agility

	playerLeftLines := []string{
		styles.SubtitleStyle.Render("[ PROFIL & ATRIBUT PETUALANG ]"),
		fmt.Sprintf("  Nama      : %s", styles.ResourceVal.Render(playerName)),
		fmt.Sprintf("  Darah (HP): %s", styles.ResourceVal.Render(fmt.Sprintf("%d / %d", p.HP, p.MaxHP))),
		fmt.Sprintf("  Kewarasan : %s", styles.ResourceWood.Render(fmt.Sprintf("%d / %d", p.Sanity, p.MaxSanity))),
		fmt.Sprintf("  Might: %-2d (+%d ATK) Agi: %-2d (+%d)", p.Stats.Might, mightBonus, p.Stats.Agility, p.Stats.Agility),
		fmt.Sprintf("  Resolve: %-2d (+%d HP)  Ing: %-2d", p.Stats.Resolve, (p.Stats.Resolve-10)*5, p.Stats.Ingenuity),
		fmt.Sprintf("  Kapasitas : %d Slot Ransel", p.MaxBackpack),
	}

	weaponName := w.Name
	maxWeaponNameLen := rightWidth - 14
	if lipgloss.Width(weaponName) > maxWeaponNameLen && maxWeaponNameLen > 3 {
		weaponName = weaponName[:maxWeaponNameLen-3] + "..."
	}

	conditionStr := "Siap Bertualang"
	if p.HP < p.MaxHP {
		conditionStr = "Pemulihan (+20 HP/Hari)"
	}

	playerRightLines := []string{
		styles.SubtitleStyle.Render("[ PERLENGKAPAN & ROMBONGAN ]"),
		fmt.Sprintf("  Senjata   : %s", styles.DefenseStyle.Render(weaponName)),
		fmt.Sprintf("  Tipe/ATK  : %s (%d-%d ATK)", w.WeaponType, w.BaseDamage[0]+mightBonus, w.BaseDamage[1]+mightBonus),
		fmt.Sprintf("  Kritikal  : %.1f%% | Init: %d", critPct, totalInit),
		fmt.Sprintf("  Ketahanan : %d/%d (%s)", w.Durability, w.MaxDura, w.SpecialAffix),
		fmt.Sprintf("  Kondisi   : %s", styles.ResourceVal.Render(conditionStr)),
		fmt.Sprintf("  Rombongan : %d Rekan | %d Ramuan", len(p.Party), p.TotalPotions()),
	}

	playerLeftCol := lipgloss.NewStyle().Width(halfWidth).Render(strings.Join(playerLeftLines, "\n"))
	playerRightCol := lipgloss.NewStyle().Width(rightWidth).Render(strings.Join(playerRightLines, "\n"))
	playerGrid := lipgloss.JoinHorizontal(lipgloss.Top, playerLeftCol, playerRightCol)

	headerBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerLine,
			divider,
			detailGrid,
			divider,
			playerGrid,
		),
	)

	// 2. Daily Log Box (Capped to at most 2 priority entries for compact layout)
	logTitle := "LAPORAN HARIAN:"
	if len(e.DailyLogs) > 2 {
		logTitle = fmt.Sprintf("LAPORAN HARIAN (PRIORITAS - 2 DARI %d PERISTIWA):", len(e.DailyLogs))
		if lipgloss.Width(logTitle) > contentWidth {
			logTitle = "LAPORAN HARIAN (2 CATATAN UTAMA):"
		}
	}
	logHeader := styles.SubtitleStyle.Render(logTitle)

	var logLines []string
	if len(e.DailyLogs) == 0 {
		logLines = append(logLines, styles.LogItemStyle.Render("Tidak ada peristiwa penting hari ini"))
	} else {
		displayLogs := selectPriorityDailyLogs(e.DailyLogs, 2)
		for _, log := range displayLogs {
			displayLog := log
			maxLen := contentWidth - 3
			if lipgloss.Width(displayLog) > maxLen {
				displayLog = displayLog[:maxLen-3] + "..."
			}
			logLines = append(logLines, styles.LogItemStyle.Render(displayLog))
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

	// 3. Action Menu Box (2-Column Grid Layout)
	actionHeader := styles.TitleStyle.Render("TINDAKAN TERSEDIA:")

	type actionItem struct {
		key, label string
	}

	leftActions := []actionItem{
		{"1", "Pasar & Komoditas"},
		{"2", "Pembangunan Fasilitas"},
		{"3", "Bengkel Pandai Besi"},
		{"4", "Pusat Latihan Karakter"},
		{"5", "Kedai Minum (Pendamping)"},
	}

	rightActions := []actionItem{
		{"6", "Ekspedisi Katakombe"},
		{"7", "Laboratorium Alkimia"},
		{"W", "Alokasi Pekerja"},
		{"D", "Lewati Hari (Produksi)"},
		{"Q", "Keluar Permainan"},
	}

	var actionRows []string
	for i := 0; i < len(leftActions); i++ {
		l := leftActions[i]
		r := rightActions[i]

		leftStr := fmt.Sprintf("%s %s", styles.KeyBadge.Render(l.key), l.label)
		rightStr := fmt.Sprintf("%s %s", styles.KeyBadge.Render(r.key), r.label)

		lCol := lipgloss.NewStyle().Width(halfWidth).Render(leftStr)
		rCol := lipgloss.NewStyle().Width(rightWidth).Render(rightStr)
		actionRows = append(actionRows, lipgloss.JoinHorizontal(lipgloss.Top, lCol, rCol))
	}

	actionBox := styles.BaseBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			actionHeader,
			strings.Join(actionRows, "\n"),
		),
	)

	// 4. Alert Banner (if any)
	var alertBox string
	if e.StatusAlert != "" {
		alertBox = styles.AlertSuccess.Render("[INFO] " + e.StatusAlert)
	}

	elements := []string{headerBox, logBox, actionBox}
	if alertBox != "" {
		elements = append(elements, alertBox)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		elements...,
	)
}

// selectPriorityDailyLogs selects at most maxCount logs sorted by importance
func selectPriorityDailyLogs(logs []string, maxCount int) []string {
	if len(logs) <= maxCount {
		return logs
	}

	getPriority := func(log string) int {
		// Priority 1: Critical threats, famine, deaths, freezing
		if strings.HasPrefix(log, "[!]") || strings.Contains(log, "KELAPARAN") || strings.Contains(log, "gugur") || strings.Contains(log, "kedinginan") {
			return 1
		}
		// Priority 2: Major events like season transitions
		if strings.HasPrefix(log, "[*]") || strings.Contains(log, "PERUBAHAN MUSIM") {
			return 2
		}
		// Priority 3: Milestones & positive settlement growth
		if strings.Contains(log, "pengembara") || strings.Contains(log, "menetap") || strings.Contains(log, "ditingkatkan") {
			return 3
		}
		// Priority 4: Consolidated daily production
		if strings.Contains(log, "Produksi Harian") {
			return 4
		}
		// Priority 5: Atmospheric and other informational logs
		return 5
	}

	type scoredLog struct {
		index    int
		priority int
		text     string
	}

	scored := make([]scoredLog, len(logs))
	for i, l := range logs {
		scored[i] = scoredLog{
			index:    i,
			priority: getPriority(l),
			text:     l,
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].priority != scored[j].priority {
			return scored[i].priority < scored[j].priority
		}
		return scored[i].index < scored[j].index
	})

	result := make([]string, maxCount)
	for i := 0; i < maxCount; i++ {
		result[i] = scored[i].text
	}
	return result
}

