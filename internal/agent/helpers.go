package agent

import (
	"fmt"

	"github.com/pterm/pterm"
)

func runStep(showProgress bool, current int, total int, label string, run func() error) error {
	if !showProgress {
		return run()
	}

	title := fmt.Sprintf("[%d/%d] %s...", current, total, label)

	spinner, err := pterm.DefaultSpinner.Start(title)
	if err != nil {
		return err
	}
	spinner.SuccessPrinter = &pterm.PrefixPrinter{
		MessageStyle: pterm.NewStyle(pterm.FgGreen),
		Prefix: pterm.Prefix{
			Text:  "✓",
			Style: pterm.NewStyle(pterm.FgGreen),
		},
	}
	spinner.FailPrinter = &pterm.PrefixPrinter{
		MessageStyle: pterm.NewStyle(pterm.FgRed),
		Prefix: pterm.Prefix{
			Text:  "✗",
			Style: pterm.NewStyle(pterm.FgRed),
		},
	}

	if err := run(); err != nil {
		spinner.Fail(fmt.Sprintf("[%d/%d] %s failed", current, total, label))

		return err
	}

	spinner.Success(fmt.Sprintf("[%d/%d] %s", current, total, label))

	return nil
}
