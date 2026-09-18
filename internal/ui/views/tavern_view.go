package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderTavernView renders the tavern recruitment and rumor screen
func RenderTavernView(e *engine.Engine, selectedMercIdx int, activeTab int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	tavLvl := e.Village.Buildings[settlement.BuildingTavern]
	titleText := fmt.Sprintf("KEDAI MINUM: REKRUTMEN PENDAMPING & RUMOR (LEVEL %d)", tavLvl)
	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	halfWidth := contentWidth / 2
	rightWidth := contentWidth - halfWidth

	// 1. Header Box: Party capacity and settlement supplies
	maxParty := e.Tavern.MaxPartySize(tavLvl)
	partyCount := len(e.Player.Party)

	partyStr := fmt.Sprintf("%d / %d Rekan", partyCount, maxParty)
	partyStyle := styles.ResourceVal
	if partyCount >= maxParty && maxParty > 0 {
		partyStyle = styles.AlertSuccess
	}

	leftLines := []string{
		styles.SubtitleStyle.Render("[ STATUS ROMBONGAN PETUALANG ]"),
		fmt.Sprintf("  Kapasitas Tim  : %s", partyStyle.Render(partyStr)),
		fmt.Sprintf("  Fasilitas Kedai: Level %d", tavLvl),
	}
	if partyCount == 0 {
		leftLines = append(leftLines, styles.ResourceLabel.Render("  Rekan Aktif    : Belum ada rekan"))
	} else {
		for _, comp := range e.Player.Party {
			compName := comp.Name
			if len(compName) > 12 {
				compName = compName[:10] + ".."
			}
			leftLines = append(leftLines, fmt.Sprintf("  * %s (%s, %d HP)", compName, comp.Role, comp.HP))
		}
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ LOGISTIK & KAS PEMUKIMAN ]"),
		fmt.Sprintf("  Kas Emas Desa  : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", e.Village.Treasury))),
		fmt.Sprintf("  Stok Ransum    : %s", styles.ResourceWood.Render(fmt.Sprintf("%d Ransum", e.Village.Rations))),
		fmt.Sprintf("  Kewarasan Jiwa : %s", styles.ResourceWood.Render(fmt.Sprintf("%d / %d Sanity", e.Player.Sanity, e.Player.MaxSanity))),
		fmt.Sprintf("  Darah Petualang: %s", styles.ResourceVal.Render(fmt.Sprintf("%d / %d HP", e.Player.HP, e.Player.MaxHP))),
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

	// Tab switcher
	var tabBar string
	if activeTab == 0 {
		tabBar = styles.ItemHighlight.Render("[1] REKRUTMEN PENDAMPING") + "    " + styles.ItemNormal.Render("[2/TAB] RUMOR & SANTAP KEDAI")
	} else {
		tabBar = styles.ItemNormal.Render("[1/TAB] REKRUTMEN PENDAMPING") + "    " + styles.ItemHighlight.Render("[2] RUMOR & SANTAP KEDAI")
	}

	var mainBox string
	var detailBox string

	if activeTab == 0 {
		// Tab 0: Mercenary Recruitment Table
		mercs := e.Tavern.Mercenaries
		if selectedMercIdx < 0 {
			selectedMercIdx = 0
		}
		if selectedMercIdx >= len(mercs) {
			selectedMercIdx = len(mercs) - 1
		}

		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorCyan).Padding(0, 1)
		dividerStyle := lipgloss.NewStyle().Foreground(styles.ColorDim).Padding(0, 1)

		tableHeader := fmt.Sprintf("  %-17s | %-10s | %-11s | %-9s | %-9s",
			"NAMA REKAN", "PERAN", "BIAYA SEWA", "UPAH JARAH", "STATUS")
		tableDivider := strings.Repeat("─", 70)

		tableLines := []string{
			tabBar,
			"",
			headerStyle.Render(tableHeader),
			dividerStyle.Render(tableDivider),
		}

		itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)
		itemHired := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorBlue).Padding(0, 1)

		for i, m := range mercs {
			cursor := "  "
			isSelected := i == selectedMercIdx
			if isSelected {
				cursor = "▶ "
			}

			statusStr := "Tersedia"
			if m.IsHired {
				statusStr = "Bergabung"
			} else if tavLvl < 1 {
				statusStr = "Terkunci"
			} else if partyCount >= maxParty {
				statusStr = "Penuh"
			} else if e.Village.Treasury < m.Def.HireCost {
				statusStr = "Kurang"
			}

			nameStr := m.Def.Name
			if len(nameStr) > 17 {
				nameStr = nameStr[:14] + "..."
			}

			hireCostStr := fmt.Sprintf("%d Gold", m.Def.HireCost)
			cutStr := fmt.Sprintf("%d%% Emas", m.Def.CutPercent)

			if len(statusStr) > 9 {
				statusStr = statusStr[:9]
			}

			rowText := fmt.Sprintf("%s%-17s | %-10s | %-11s | %-9s | %-9s",
				cursor, nameStr, m.Def.RoleDisplay, hireCostStr, cutStr, statusStr)

			if isSelected {
				tableLines = append(tableLines, styles.ItemHighlight.Render(rowText))
			} else if m.IsHired {
				tableLines = append(tableLines, itemHired.Render(rowText))
			} else if tavLvl >= 1 && partyCount < maxParty && e.Village.Treasury >= m.Def.HireCost {
				tableLines = append(tableLines, itemAvailable.Render(rowText))
			} else {
				tableLines = append(tableLines, styles.ItemNormal.Render(rowText))
			}
		}

		mainBox = styles.ActiveBox.Width(boxWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left, tableLines...),
		)

		// Detail Box for selected mercenary
		if selectedMercIdx >= 0 && selectedMercIdx < len(mercs) {
			sel := mercs[selectedMercIdx]
			var dLines []string
			dLines = append(dLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN REKAN: %s (%s) ]", strings.ToUpper(sel.Def.Name), strings.ToUpper(sel.Def.RoleDisplay))))

			descLines := styles.WrapText(sel.Def.Description, contentWidth-19)
			if len(descLines) > 0 {
				dLines = append(dLines, fmt.Sprintf("  Latar Belakang : %s", descLines[0]))
				for _, dl := range descLines[1:] {
					dLines = append(dLines, fmt.Sprintf("                   %s", dl))
				}
			}

			perkLines := styles.WrapText(sel.Def.PerkDesc, contentWidth-19)
			if len(perkLines) > 0 {
				dLines = append(dLines, fmt.Sprintf("  Keahlian Khusus: %s", styles.ResourceVal.Render(perkLines[0])))
				for _, pl := range perkLines[1:] {
					dLines = append(dLines, fmt.Sprintf("                   %s", styles.ResourceVal.Render(pl)))
				}
			}

			dLines = append(dLines, fmt.Sprintf("  Ketahanan Jiwa : %s | %s | %s",
				styles.ResourceVal.Render(fmt.Sprintf("%d HP", sel.Def.HP)),
				styles.ResourceGold.Render(fmt.Sprintf("Biaya Sewa %d Gold", sel.Def.HireCost)),
				styles.ResourceWood.Render(fmt.Sprintf("Potongan Jarahan %d%%", sel.Def.CutPercent)),
			))

			if sel.IsHired {
				dLines = append(dLines, fmt.Sprintf("  Aksi Rekrutmen : %s", styles.AlertWarning.Render("Tekan [ENTER] untuk memberhentikan (lepas) dari tim")))
			} else {
				canHireErr := e.Tavern.CanHire(sel.Def.ID, tavLvl, e.Village.Treasury, partyCount)
				if canHireErr != nil {
					dLines = append(dLines, fmt.Sprintf("  Aksi Rekrutmen : %s", styles.AlertError.Render(canHireErr.Error())))
				} else {
					dLines = append(dLines, fmt.Sprintf("  Aksi Rekrutmen : %s", styles.AlertSuccess.Render("Tekan [ENTER] untuk merekrut ke dalam rombongan")))
				}
			}

			detailBox = styles.BaseBox.Width(boxWidth).Render(
				lipgloss.JoinVertical(lipgloss.Left, dLines...),
			)
		}
	} else {
		// Tab 1: Atmosphere & Tavern Dining
		var tavernLines []string
		tavernLines = append(tavernLines, tabBar)
		tavernLines = append(tavernLines, "")
		tavernLines = append(tavernLines, styles.SubtitleStyle.Render("[ SUASANA & KABAR ANGIN DI KEDAI ]"))
		tavernLines = append(tavernLines, fmt.Sprintf("  \"%s\"", styles.ResourceVal.Render(e.Tavern.GetRandomRumor())))
		tavernLines = append(tavernLines, "")
		tavernLines = append(tavernLines, divider)
		tavernLines = append(tavernLines, "")
		tavernLines = append(tavernLines, styles.SubtitleStyle.Render("[ SANTAPAN HANGAT & KETENANGAN JIWA ]"))
		tavernLines = append(tavernLines, "  Pesan seporsi hidangan daging panggang dan bir gandum hangat untuk meredakan")
		tavernLines = append(tavernLines, "  kengerian lorong bawah tanah serta memulihkan kesehatan tubuh petualang.")
		tavernLines = append(tavernLines, "")
		tavernLines = append(tavernLines, fmt.Sprintf("  Biaya Hidangan : %s dan %s", styles.ResourceGold.Render("10 Gold"), styles.ResourceWood.Render("1 Ransum")))
		tavernLines = append(tavernLines, fmt.Sprintf("  Efek Pemulihan : %s dan %s", styles.ResourceWood.Render("+20 Sanity"), styles.ResourceVal.Render("+15 HP")))
		tavernLines = append(tavernLines, "")

		restErr := e.Tavern.CanRest(e.Village.Treasury, e.Village.Rations)
		if restErr != nil {
			tavernLines = append(tavernLines, fmt.Sprintf("  Kondisi Pesanan: %s", styles.AlertError.Render(restErr.Error())))
		} else {
			tavernLines = append(tavernLines, fmt.Sprintf("  Kondisi Pesanan: %s", styles.AlertSuccess.Render("Bahan mencukupi! Tekan [ENTER] untuk santap & istirahat")))
		}

		mainBox = styles.ActiveBox.Width(boxWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left, tavernLines...),
		)
	}

	// Alert Banner
	var alertBox string
	if e.StatusAlert != "" {
		alertBox = styles.AlertSuccess.Render("[INFO] " + e.StatusAlert)
	}

	// Bottom Controls Box
	var controls string
	if activeTab == 0 {
		controls = fmt.Sprintf("%s Tab   %s Pilih   %s Rekrut/Lepas   %s Kembali",
			styles.KeyBadgeHighlight.Render("TAB"),
			styles.KeyBadgeHighlight.Render("↑/↓"),
			styles.KeyBadgeHighlight.Render("ENTER"),
			styles.KeyBadge.Render("ESC"),
		)
	} else {
		controls = fmt.Sprintf("%s Tab   %s Pesan Santapan   %s Kembali",
			styles.KeyBadgeHighlight.Render("TAB"),
			styles.KeyBadgeHighlight.Render("ENTER"),
			styles.KeyBadge.Render("ESC"),
		)
	}
	bottomBox := styles.BaseBox.Width(boxWidth).Render(controls)

	elements := []string{headerBox, mainBox}
	if detailBox != "" {
		elements = append(elements, detailBox)
	}
	if alertBox != "" {
		elements = append(elements, alertBox)
	}
	elements = append(elements, bottomBox)

	return lipgloss.JoinVertical(lipgloss.Left, elements...)
}
