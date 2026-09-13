package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Palette definitions - Dark fantasy / Cyber-medieval CLI aesthetic
var (
	ColorGold     = lipgloss.Color("#E5C07B")
	ColorAmber    = lipgloss.Color("#D19A66")
	ColorGreen    = lipgloss.Color("#98C379")
	ColorEmerald  = lipgloss.Color("#50FA7B")
	ColorRed      = lipgloss.Color("#E06C75")
	ColorCrimson  = lipgloss.Color("#FF5555")
	ColorBlue     = lipgloss.Color("#61AFEF")
	ColorCyan     = lipgloss.Color("#56B6C2")
	ColorPurple   = lipgloss.Color("#C678DD")
	ColorGray     = lipgloss.Color("#ABB2BF")
	ColorDim      = lipgloss.Color("#5C6370")
	ColorDark     = lipgloss.Color("#282C34")
	ColorDarker   = lipgloss.Color("#1E222A")
	ColorBorder   = lipgloss.Color("#4B5263")
	ColorActive   = lipgloss.Color("#61AFEF")
	ColorWarning  = lipgloss.Color("#E5C07B")
	ColorError    = lipgloss.Color("#E06C75")
)

// Typography & element styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGold)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorCyan)

	SeasonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAmber)

	DayCounterStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPurple)

	KeyBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan).
			Background(ColorDarker).
			Padding(0, 1)

	KeyBadgeHighlight = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorBlue).
				Padding(0, 1)

	ResourceLabel = lipgloss.NewStyle().
			Foreground(ColorGray)

	ResourceVal = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	ResourceGold = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGold)

	ResourceWood = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAmber)

	ResourceStone = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGray)

	ResourceFood = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGreen)

	DefenseStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBlue)

	LogItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DCDFE4")).
			PaddingLeft(1)

	AlertSuccess = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorEmerald).
			Padding(0, 1)

	AlertWarning = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning).
			Padding(0, 1)

	AlertError = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorError).
			Padding(0, 1)

	// Box Containers
	BaseBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	ActiveBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorActive).
			Padding(0, 1)

	ItemHighlight = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2C313C")).
			Padding(0, 1)

	ItemNormal = lipgloss.NewStyle().
			Foreground(ColorGray).
			Padding(0, 1)
)

// CleanText removes variation selector 16 to maintain consistent terminal width
func CleanText(s string) string {
	return strings.ReplaceAll(s, "\ufe0f", "")
}
