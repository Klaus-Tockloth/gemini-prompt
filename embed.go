package main

import (
	_ "embed"
	"log"
	"os"
)

//go:embed gemini-prompt.yaml
var geminiPromptYaml []byte

func writeConfig() {
	filename := "gemini-prompt.yaml"
	err := os.WriteFile(filename, geminiPromptYaml, 0666)
	if err != nil {
		log.Fatalf("embed: error [%v] at os.WriteFile(), file = [%s]", err, filename)
	}
}

//go:embed prompt-input.html
var geminiPromptInputHTML []byte

func writePromptInput() {
	filename := "prompt-input.html"
	err := os.WriteFile(filename, geminiPromptInputHTML, 0666)
	if err != nil {
		log.Fatalf("embed: error [%v] at os.WriteFile(), file = [%s]", err, filename)
	}
}

//go:embed assets/gemini-prompt.css
var assetsGeminiPromptCSS []byte

//go:embed assets/gemini-prompt-303030.svg
var assetsGeminiPrompt303030Svg []byte

//go:embed assets/gemini-prompt-ebebeb.svg
var assetsGeminiPromptEbebebSvg []byte

//go:embed assets/copy-to-clipboard.js
var assetsCopyToClipboardJs []byte

func writeAssets(basepath string) {
	filename := basepath + "/assets/gemini-prompt.css"
	err := os.WriteFile(filename, assetsGeminiPromptCSS, 0666)
	if err != nil {
		log.Fatalf("embed: error [%v] at os.WriteFile(), file = [%s]", err, filename)
	}

	filename = basepath + "/assets/gemini-prompt-303030.svg"
	err = os.WriteFile(filename, assetsGeminiPrompt303030Svg, 0666)
	if err != nil {
		log.Fatalf("embed: error [%v] at os.WriteFile(), file = [%s]", err, filename)
	}

	filename = basepath + "/assets/gemini-prompt-ebebeb.svg"
	err = os.WriteFile(filename, assetsGeminiPromptEbebebSvg, 0666)
	if err != nil {
		log.Fatalf("embed: error [%v] at os.WriteFile(), file = [%s]", err, filename)
	}

	filename = basepath + "/assets/copy-to-clipboard.js"
	err = os.WriteFile(filename, assetsCopyToClipboardJs, 0666)
	if err != nil {
		log.Fatalf("embed: error [%v] at os.WriteFile(), file = [%s]", err, filename)
	}
}
