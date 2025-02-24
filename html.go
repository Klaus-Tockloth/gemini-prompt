package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/aquilax/truncate"
)

/*
renderMarkdown2Html renders markdown to html.
*/
func renderMarkdown2Html(md string) string {
	var buf bytes.Buffer
	err := markdownParser.Convert([]byte(md), &buf)
	if err != nil {
		fmt.Printf("error [%v] at markdownParser.Convert()", err)
	}

	// replace Html elements
	htmlDataModified := string(buf.String())
	for _, item := range progConfig.HtmlReplaceElements {
		for key, value := range item {
			htmlDataModified = strings.ReplaceAll(htmlDataModified, key, value)
		}
	}

	return htmlDataModified
}

/*
buildHtmlPage builds html with header, body and footer.
*/
func buildHtmlPage(prompt, source, destination string) error {
	htmlBody, err := os.ReadFile(source)
	if err != nil {
		fmt.Printf("error [%v] at os.ReadFile()", err)
		return err
	}

	title := strings.ReplaceAll(prompt, "\r\n", " ")
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\t", " ")

	title = truncate.Truncate(title, progConfig.HtmlMaxLengthTitle, "...", truncate.PositionEnd)
	htmlHeader := fmt.Sprintf(progConfig.HtmlHeader, title)
	htmlFooter := progConfig.HtmlFooter

	// build html page
	htmlPage := htmlHeader + string(htmlBody) + htmlFooter

	// write html to file
	err = os.WriteFile(destination, []byte(htmlPage), 0666)
	if err != nil {
		fmt.Printf("error [%v] at os.WriteFile()", err)
		return err
	}

	return nil
}
