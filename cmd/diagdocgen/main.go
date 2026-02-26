// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

// package main implements command diagdocgen that generates diagnostics markdown registry from catalog.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/woozymasta/rvcfg"
)

type stageSection struct {
	title string
	items []rvcfg.DiagnosticSpec
}

type tableWidths struct {
	code     int
	severity int
	summary  int
}

func main() {
	out := flag.String("out", "DIAGNOSTICS.md", "output markdown file path")
	flag.Parse()

	content := renderDiagnosticsMarkdown(rvcfg.DiagnosticCatalog())
	if err := os.WriteFile(*out, []byte(content), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "diagdocgen: write %s: %v\n", *out, err)
		os.Exit(1)
	}
}

func renderDiagnosticsMarkdown(catalog []rvcfg.DiagnosticSpec) string {
	sections := []stageSection{
		{title: "Lexer", items: make([]rvcfg.DiagnosticSpec, 0)},
		{title: "Parser", items: make([]rvcfg.DiagnosticSpec, 0)},
		{title: "Processor", items: make([]rvcfg.DiagnosticSpec, 0)},
	}

	for _, spec := range catalog {
		switch spec.Stage {
		case rvcfg.StageLex:
			sections[0].items = append(sections[0].items, spec)
		case rvcfg.StageParse:
			sections[1].items = append(sections[1].items, spec)
		case rvcfg.StagePreprocess:
			sections[2].items = append(sections[2].items, spec)
		}
	}

	var b strings.Builder
	b.WriteString("# Diagnostics Registry\n\n")

	for _, section := range sections {
		widths := calcTableWidths(section.items)

		b.WriteString("## ")
		b.WriteString(section.title)
		b.WriteString("\n\n")
		b.WriteString(renderRow("Code", "Severity", "Summary", widths))
		b.WriteString(renderAlignRow(widths))

		for _, spec := range section.items {
			b.WriteString(renderRow(
				fmt.Sprintf("`%s`", spec.Code),
				severityMarkdown(spec.Severity),
				spec.Summary,
				widths,
			))
		}

		b.WriteString("\n")
	}

	return b.String()
}

func severityMarkdown(severity rvcfg.Severity) string {
	switch severity {
	case rvcfg.SeverityWarning:
		return "**warning**"
	case rvcfg.SeverityError:
		return "**error**"
	default:
		return fmt.Sprintf("**%s**", severity)
	}
}

func calcTableWidths(items []rvcfg.DiagnosticSpec) tableWidths {
	widths := tableWidths{
		code:     len("Code"),
		severity: len("Severity"),
		summary:  len("Summary"),
	}

	for _, item := range items {
		codeCell := fmt.Sprintf("`%s`", item.Code)
		severityCell := severityMarkdown(item.Severity)
		summaryCell := item.Summary

		if len(codeCell) > widths.code {
			widths.code = len(codeCell)
		}

		if len(severityCell) > widths.severity {
			widths.severity = len(severityCell)
		}

		if len(summaryCell) > widths.summary {
			widths.summary = len(summaryCell)
		}
	}

	return widths
}

func renderRow(code string, severity string, summary string, widths tableWidths) string {
	return fmt.Sprintf(
		"| %s | %s | %s |\n",
		padRight(code, widths.code),
		padRight(severity, widths.severity),
		padRight(summary, widths.summary),
	)
}

func renderAlignRow(widths tableWidths) string {
	return fmt.Sprintf(
		"| %s | %s | %s |\n",
		alignCenter(widths.code),
		alignCenter(widths.severity),
		alignCenter(widths.summary),
	)
}

func alignCenter(width int) string {
	if width < 3 {
		width = 3
	}

	return ":" + strings.Repeat("-", width-2) + ":"
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}

	return value + strings.Repeat(" ", width-len(value))
}
