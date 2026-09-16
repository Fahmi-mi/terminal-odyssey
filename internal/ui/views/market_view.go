package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/styles"
)

// RenderMarketView renders the commodity trading and caravan dispatch screen
func RenderMarketView(e *engine.Engine, selectedCommodityIdx int, activeTab int, selectedRouteIdx int, width int) string {
	boxWidth := width - 4
	if boxWidth < 74 {
		boxWidth = 74
	}
	if boxWidth > 96 {
		boxWidth = 96
	}
	contentWidth := boxWidth - 2

	v := e.Village
	p := e.Player
	maxL, maxS, maxR := v.StorageCap()

	postLvl := v.Buildings[settlement.BuildingCaravanPost]
	titleText := "PASAR DESA & JARINGAN PERDAGANGAN REGIONAL"
	title := styles.TitleStyle.Render(titleText)
	divider := strings.Repeat("─", contentWidth)

	leftWidth := 28
	rightWidth := contentWidth - leftWidth

	// Header Box
	leftLines := []string{
		styles.SubtitleStyle.Render("[ KAS & LOGISTIK DESA ]"),
		fmt.Sprintf("  Kas Emas  : %s", styles.ResourceGold.Render(fmt.Sprintf("%d Gold", v.Treasury))),
		fmt.Sprintf("  Kayu      : %s", styles.ResourceWood.Render(fmt.Sprintf("%d / %d", v.Lumber, maxL))),
		fmt.Sprintf("  Batu      : %s", styles.ResourceStone.Render(fmt.Sprintf("%d / %d", v.Stone, maxS))),
		fmt.Sprintf("  Ransum    : %s", styles.ResourceFood.Render(fmt.Sprintf("%d / %d", v.Rations, maxR))),
	}

	diff := p.Stats.Ingenuity - 10
	if diff < 0 {
		diff = 0
	}
	buyDisc := diff * 1
	if buyDisc > 15 {
		buyDisc = 15
	}
	sellBonus := int(float64(diff) * 1.5)
	if sellBonus > 22 {
		sellBonus = 22
	}

	eventShort := e.Market.MarketEvent
	if lipgloss.Width(eventShort) > rightWidth-14 && rightWidth > 17 {
		eventShort = eventShort[:rightWidth-17] + "..."
	}

	rightLines := []string{
		styles.SubtitleStyle.Render("[ STATUS PERDAGANGAN ]"),
		fmt.Sprintf("  Petualang : %s (Ingenuity: %d)", styles.ResourceVal.Render(p.Name), p.Stats.Ingenuity),
		fmt.Sprintf("  Tawar     : %s", styles.ItemHighlight.Copy().Padding(0, 0).Render(fmt.Sprintf("Diskon -%d%% | Margin +%d%%", buyDisc, sellBonus))),
		fmt.Sprintf("  Musim     : %s", styles.SeasonStyle.Render(e.CurrentSeason.String())),
		fmt.Sprintf("  Kondisi   : %s", styles.ResourceVal.Render(eventShort)),
	}

	leftCol := lipgloss.NewStyle().Width(leftWidth).Render(strings.Join(leftLines, "\n"))
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

	var contentLines []string
	var controls []string

	if activeTab == 0 {
		// Tab 0: Pasar Lokal
		tabHeader := styles.ItemHighlight.Render("[1] PASAR LOKAL") + "    " + styles.ItemNormal.Render("[2/TAB] POS KAFILAH")
		contentLines = append(contentLines, tabHeader)
		contentLines = append(contentLines, "")

		items := e.Market.Items
		if selectedCommodityIdx < 0 {
			selectedCommodityIdx = 0
		}
		if selectedCommodityIdx >= len(items) {
			selectedCommodityIdx = len(items) - 1
		}

		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorCyan).Padding(0, 1)
		dividerStyle := lipgloss.NewStyle().Foreground(styles.ColorDim).Padding(0, 1)

		tableHeader := fmt.Sprintf("  %-17s | %-8s | %-11s | %-5s | %-5s | %-9s",
			"Komoditas", "Kategori", "Stok Desa", "Beli", "Jual", "Tren")
		tableDivider := strings.Repeat("─", 72)
		contentLines = append(contentLines, headerStyle.Render(tableHeader))
		contentLines = append(contentLines, dividerStyle.Render(tableDivider))

		itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)

		for i, it := range items {
			cursor := "  "
			if i == selectedCommodityIdx {
				cursor = "▶ "
			}

			stock := v.GetCommodityStock(it.Def.ID)
			effectiveBuy := e.Market.GetEffectiveBuyPrice(it.Def.ID, p.Stats.Ingenuity)
			effectiveSell := e.Market.GetEffectiveSellPrice(it.Def.ID, p.Stats.Ingenuity)

			nameStr := it.Def.Name
			if len(nameStr) > 17 {
				nameStr = nameStr[:15] + ".."
			}

			stockText := fmt.Sprintf("%3d %s", stock, it.Def.Unit)
			if len(stockText) > 11 {
				stockText = stockText[:11]
			}

			buyStr := fmt.Sprintf("%3d G", effectiveBuy)
			sellStr := fmt.Sprintf("%3d G", effectiveSell)

			trendStr := it.Trend

			rowText := fmt.Sprintf("%s%-17s | %-8s | %-11s | %-5s | %-5s | %-9s",
				cursor,
				nameStr,
				it.Def.Category,
				stockText,
				buyStr,
				sellStr,
				trendStr,
			)

			if i == selectedCommodityIdx {
				contentLines = append(contentLines, styles.ItemHighlight.Render(rowText))
			} else if stock > 0 || v.Treasury >= effectiveBuy {
				contentLines = append(contentLines, itemAvailable.Render(rowText))
			} else {
				contentLines = append(contentLines, styles.ItemNormal.Render(rowText))
			}
		}

		contentLines = append(contentLines, "")

		// Detail of selected commodity
		sel := items[selectedCommodityIdx]
		contentLines = append(contentLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN: %s ]", strings.ToUpper(sel.Def.Name))))

		descWrapped := styles.WrapText(sel.Def.Description, contentWidth-4)
		for _, dw := range descWrapped {
			contentLines = append(contentLines, styles.LogItemStyle.Render(dw))
		}

		effBuy := e.Market.GetEffectiveBuyPrice(sel.Def.ID, p.Stats.Ingenuity)
		effSell := e.Market.GetEffectiveSellPrice(sel.Def.ID, p.Stats.Ingenuity)
		stock := v.GetCommodityStock(sel.Def.ID)

		statLine := fmt.Sprintf("  Kategori       : %-8s | Satuan: %-7s | Stok Desa: %d %s",
			sel.Def.Category, sel.Def.Unit, stock, sel.Def.Unit)
		contentLines = append(contentLines, statLine)

		priceLine := fmt.Sprintf("  Harga Beli     : %d Gold (-%d%% Diskon) | Harga Jual: %d Gold (+%d%% Bonus)",
			effBuy, buyDisc, effSell, sellBonus)
		contentLines = append(contentLines, priceLine)

		canBuy := v.Treasury >= effBuy
		if sel.Def.ID == "lumber" && v.Lumber >= maxL {
			canBuy = false
		} else if sel.Def.ID == "stone" && v.Stone >= maxS {
			canBuy = false
		} else if sel.Def.ID == "rations" && v.Rations >= maxR {
			canBuy = false
		}

		var statusVal string
		if canBuy && stock > 0 {
			statusVal = styles.AlertSuccess.Copy().Padding(0, 0).Render("Siap Ditransaksikan (Tekan [B] Beli, [S] Jual)")
		} else if canBuy {
			statusVal = styles.AlertSuccess.Copy().Padding(0, 0).Render("Siap Dibeli (Tekan [B] untuk Membeli)")
		} else if stock > 0 {
			statusVal = styles.AlertWarning.Copy().Padding(0, 0).Render("Siap Dijual (Tekan [S] untuk Menjual)")
		} else {
			statusVal = styles.AlertError.Copy().Padding(0, 0).Render("Kas Emas tidak cukup & stok desa kosong")
		}
		contentLines = append(contentLines, fmt.Sprintf("  Status         : %s", statusVal))

		controls = []string{
			fmt.Sprintf("%s Kafilah", styles.KeyBadge.Render("TAB")),
			fmt.Sprintf("%s Pilih", styles.KeyBadge.Render("↑/↓")),
			fmt.Sprintf("%s Beli", styles.KeyBadge.Render("B")),
			fmt.Sprintf("%s Jual", styles.KeyBadge.Render("S")),
			fmt.Sprintf("%s Kembali", styles.KeyBadge.Render("ESC")),
		}
	} else {
		// Tab 1: Pos Kafilah
		tabHeader := styles.ItemNormal.Render("[1/TAB] PASAR LOKAL") + "    " + styles.ItemHighlight.Render("[2] POS KAFILAH")
		contentLines = append(contentLines, tabHeader)

		if postLvl < 1 {
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, styles.AlertWarning.Render("  [!] POS KAFILAH BELUM DIDIRIKAN DI DESA"))
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, styles.LogItemStyle.Render("  Pos Kafilah membuka ekspedisi gerobak dagang ke kota-kota lain,"))
			contentLines = append(contentLines, styles.LogItemStyle.Render("  memungkinkan arbitrase regional dan perolehan komoditas mewah."))
			contentLines = append(contentLines, "")
			contentLines = append(contentLines, styles.ResourceLabel.Render("  Silakan dirikan Pos Kafilah melalui menu [2] Pembangunan Fasilitas Desa"))
			contentLines = append(contentLines, "")

			controls = []string{
				fmt.Sprintf("%s Pasar Lokal", styles.KeyBadge.Render("TAB")),
				fmt.Sprintf("%s Kembali", styles.KeyBadge.Render("ESC")),
			}
		} else {
			routes := e.Caravans.Routes
			if selectedRouteIdx < 0 {
				selectedRouteIdx = 0
			}
			if selectedRouteIdx >= len(routes) {
				selectedRouteIdx = len(routes) - 1
			}

			maxArmada := e.Caravans.MaxConcurrentCaravans(postLvl)
			activeCount := len(e.Caravans.ActiveCaravans)
			militiaCount := v.Workers.Militia
			militiaBonus := militiaCount * 3

			contentLines = append(contentLines, "")
			contentLines = append(contentLines, styles.ResourceVal.Render(
				fmt.Sprintf("  Armada Dagang: %d/%d Berjalan | Pos Kafilah Lv.%d | Proteksi Milisi: +%d%%",
					activeCount, maxArmada, postLvl, militiaBonus)))
			contentLines = append(contentLines, "")

			headerStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorCyan).Padding(0, 1)
			dividerStyle := lipgloss.NewStyle().Foreground(styles.ColorDim).Padding(0, 1)

			tableHeader := fmt.Sprintf("  %-20s | %-5s | %-15s | %-6s | %-12s",
				"Rute Tujuan", "Level", "Modal & Kargo", "Waktu", "Status")
			tableDivider := strings.Repeat("─", 72)
			contentLines = append(contentLines, headerStyle.Render(tableHeader))
			contentLines = append(contentLines, dividerStyle.Render(tableDivider))

			itemInFlight := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorBlue).Padding(0, 1)
			itemAvailable := lipgloss.NewStyle().Foreground(styles.ColorGreen).Padding(0, 1)

			for i, r := range routes {
				cursor := "  "
				if i == selectedRouteIdx {
					cursor = "▶ "
				}

				isInFlight := false
				daysLeft := 0
				for _, ac := range e.Caravans.ActiveCaravans {
					if ac.Route.ID == r.ID {
						isInFlight = true
						daysLeft = ac.DaysRemaining
						break
					}
				}

				travelDays := r.BaseDays
				if e.CurrentSeason == settlement.SeasonWinter {
					travelDays = r.BaseDays * 2
				}

				reqCommodityShort := formatCommodityShort(r.CargoCommodity)
				cargoStr := fmt.Sprintf("%dG, %d %s", r.GoldInvestment, r.CargoAmount, reqCommodityShort)
				if len(cargoStr) > 15 {
					cargoStr = cargoStr[:15]
				}

				statusStr := "Siap Kirim"
				if isInFlight {
					statusStr = fmt.Sprintf("Sisa %d Hari", daysLeft)
				} else if postLvl < r.RequiredLevel {
					statusStr = fmt.Sprintf("Kunci Lv. %d", r.RequiredLevel)
				}

				lvlStr := fmt.Sprintf("Lv. %d", r.RequiredLevel)
				timeStr := fmt.Sprintf("%d Hari", travelDays)

				rowText := fmt.Sprintf("%s%-20s | %-5s | %-15s | %-6s | %-12s",
					cursor,
					r.Name,
					lvlStr,
					cargoStr,
					timeStr,
					statusStr,
				)

				if i == selectedRouteIdx {
					contentLines = append(contentLines, styles.ItemHighlight.Render(rowText))
				} else if isInFlight {
					contentLines = append(contentLines, itemInFlight.Render(rowText))
				} else if postLvl >= r.RequiredLevel {
					contentLines = append(contentLines, itemAvailable.Render(rowText))
				} else {
					contentLines = append(contentLines, styles.ItemNormal.Render(rowText))
				}
			}

			contentLines = append(contentLines, "")

			// Detail of selected route
			selRoute := routes[selectedRouteIdx]
			contentLines = append(contentLines, styles.SubtitleStyle.Render(fmt.Sprintf("[ RINCIAN RUTE: %s ]", strings.ToUpper(selRoute.Name))))

			descWrapped := styles.WrapText(selRoute.Description, contentWidth-4)
			for _, dw := range descWrapped {
				contentLines = append(contentLines, styles.LogItemStyle.Render(dw))
			}

			selCargoName := formatCommodityName(selRoute.CargoCommodity)
			bonusName := formatCommodityName(selRoute.BonusItemID)

			modalText := fmt.Sprintf("  Modal & Kargo  : %d Gold, %d %s",
				selRoute.GoldInvestment, selRoute.CargoAmount, selCargoName)
			contentLines = append(contentLines, modalText)

			seasonInfo := "Kondisi Normal"
			travelDays := selRoute.BaseDays
			if e.CurrentSeason == settlement.SeasonWinter {
				travelDays = selRoute.BaseDays * 2
				seasonInfo = "Musim Dingin x2"
			}
			waktuText := fmt.Sprintf("  Waktu Tempuh   : %d Hari (%s)", travelDays, seasonInfo)
			contentLines = append(contentLines, waktuText)

			rewardText := fmt.Sprintf("  Estimasi Hasil : %d - %d Gold + %d %s",
				selRoute.RewardGoldMin, selRoute.RewardGoldMax, selRoute.BonusItemAmount, bonusName)
			contentLines = append(contentLines, rewardText)

			effectiveRisk := selRoute.AmbushRisk - militiaBonus
			if effectiveRisk < 0 {
				effectiveRisk = 0
			}
			riskText := fmt.Sprintf("  Risiko Bandit  : %d%% (Basis: %d%%, Diredam Milisi: -%d%%)",
				effectiveRisk, selRoute.AmbushRisk, militiaBonus)
			contentLines = append(contentLines, riskText)

			isInFlight := false
			daysLeft := 0
			for _, ac := range e.Caravans.ActiveCaravans {
				if ac.Route.ID == selRoute.ID {
					isInFlight = true
					daysLeft = ac.DaysRemaining
					break
				}
			}

			var statusVal string
			if isInFlight {
				statusVal = styles.AlertWarning.Copy().Padding(0, 0).Render(
					fmt.Sprintf("Sedang menempuh perjalanan (Kembali dalam %d hari)", daysLeft))
			} else if postLvl < selRoute.RequiredLevel {
				statusVal = styles.AlertError.Copy().Padding(0, 0).Render(
					fmt.Sprintf("Terkunci (Perlu Pos Kafilah Level %d)", selRoute.RequiredLevel))
			} else {
				errDispatch := e.Caravans.CanDispatch(selRoute.ID, postLvl, v.Treasury, v.Commodities, v.Lumber, v.Rations)
				if errDispatch != nil {
					errMsg := errDispatch.Error()
					maxErrLen := contentWidth - 36
					if lipgloss.Width(errMsg) > maxErrLen {
						errMsg = errMsg[:maxErrLen-3] + "..."
					}
					statusVal = styles.AlertWarning.Copy().Padding(0, 0).Render(
						fmt.Sprintf("Syarat belum lengkap (%s)", errMsg))
				} else {
					statusVal = styles.AlertSuccess.Copy().Padding(0, 0).Render(
						"Siap Berangkat (Tekan [ENTER] untuk Berangkatkan)")
				}
			}
			contentLines = append(contentLines, fmt.Sprintf("  Status Armada  : %s", statusVal))

			controls = []string{
				fmt.Sprintf("%s Pasar Lokal", styles.KeyBadge.Render("TAB")),
				fmt.Sprintf("%s Pilih", styles.KeyBadge.Render("↑/↓")),
				fmt.Sprintf("%s Berangkatkan", styles.KeyBadge.Render("ENTER")),
				fmt.Sprintf("%s Kembali", styles.KeyBadge.Render("ESC")),
			}
		}
	}

	contentBox := styles.ActiveBox.Width(boxWidth).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			contentLines...,
		),
	)

	controlsText := strings.Join(controls, " | ")
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

func formatCommodityName(id string) string {
	switch id {
	case "lumber":
		return "Kayu Bangunan"
	case "stone":
		return "Batu Tambang"
	case "rations":
		return "Ransum Kering"
	case "iron_ore":
		return "Bijih Besi Murni"
	case "herbal_salve":
		return "Salep Herbal"
	case "spices":
		return "Rempah Lembah"
	case "caravan_silk":
		return "Sutra Musafir"
	case "ancient_relic":
		return "Relik Kuno"
	default:
		return id
	}
}

func formatCommodityShort(id string) string {
	switch id {
	case "lumber":
		return "Kayu"
	case "stone":
		return "Batu"
	case "rations":
		return "Ransum"
	case "iron_ore":
		return "Besi"
	case "herbal_salve":
		return "Salep"
	case "spices":
		return "Rempah"
	case "caravan_silk":
		return "Sutra"
	case "ancient_relic":
		return "Relik"
	default:
		return id
	}
}
