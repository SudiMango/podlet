package util

import (
	"strings"

	"github.com/pterm/pterm"
)

var (
	Tick  = marker("✓", pterm.FgGreen)
	Cross = marker("✗", pterm.FgRed)
)

func PrintDivider() {
	pterm.FgGray.Println(strings.Repeat("─", 40))
}

func marker(symbol string, color pterm.Color) pterm.PrefixPrinter {
	style := pterm.NewStyle(color)
	return pterm.PrefixPrinter{MessageStyle: style, Prefix: pterm.Prefix{Text: symbol, Style: style}}
}
