package main

import (
	"os"
	"strings"

	markdown "github.com/Klaus-Tockloth/go-term-markdown"
	"golang.org/x/term"
)

/*
renderMarkdown2Ansi renders markdown to ansi.
*/
func renderMarkdown2Ansi(md string) string {
	terminalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	terminalData := markdown.Render(md, terminalWidth, 0)

	// replace ANSI colors in terminal data
	terminalDataModified := string(terminalData)

	for _, item := range progConfig.AnsiReplaceColors {
		for key, value := range item {
			terminalDataModified = strings.ReplaceAll(terminalDataModified, key, value)
		}
	}

	return terminalDataModified
}
