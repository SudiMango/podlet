package util

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04")
}

func PrintTable(w io.Writer, header []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	fmt.Fprintln(tw, strings.Join(header, "\t"))
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	return tw.Flush()
}

func PrintFields(w io.Writer, labelsAndValues ...string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	for i := 0; i+1 < len(labelsAndValues); i += 2 {
		fmt.Fprintf(tw, "%s:\t%s\n", labelsAndValues[i], labelsAndValues[i+1])
	}

	return tw.Flush()
}
