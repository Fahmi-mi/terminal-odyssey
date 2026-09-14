package styles

import (
	"fmt"
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

// RenderHPBar renders an ASCII health bar using = and -
func RenderHPBar(current, maxVal, barLen int) string {
	if maxVal <= 0 {
		maxVal = 1
	}
	if current < 0 {
		current = 0
	}
	if current > maxVal {
		current = maxVal
	}

	filled := (current * barLen) / maxVal
	empty := barLen - filled

	filledStr := strings.Repeat("=", filled)
	emptyStr := strings.Repeat("-", empty)

	colorStyle := AlertSuccess
	ratio := float64(current) / float64(maxVal)
	if ratio <= 0.35 {
		colorStyle = AlertError
	} else if ratio <= 0.65 {
		colorStyle = AlertWarning
	}

	return fmt.Sprintf("[%s%s] %d/%d HP", colorStyle.Render(filledStr), ResourceLabel.Render(emptyStr), current, maxVal)
}

// WrapText splits a long text string into multiple lines of at most maxWidth characters
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := words[0]

	for _, w := range words[1:] {
		if lipgloss.Width(currentLine)+1+lipgloss.Width(w) <= maxWidth {
			currentLine += " " + w
		} else {
			lines = append(lines, currentLine)
			currentLine = w
		}
	}
	lines = append(lines, currentLine)
	return lines
}

