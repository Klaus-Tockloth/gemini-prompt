package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProgConfig represents program configuration
type ProgConfig struct {
	GeminiAPIKey                    string  `yaml:"GeminiAPIKey"`
	GeminiAiModel                   string  `yaml:"GeminiAiModel"`
	GeminiCandidateCount            int32   `yaml:"GeminiCandidateCount"`
	GeminiMaxOutputTokens           int32   `yaml:"GeminiMaxOutputTokens"`
	GeminiTemperature               float32 `yaml:"GeminiTemperature"`
	GeminiTopP                      float32 `yaml:"GeminiTopP"`
	GeminiTopK                      int32   `yaml:"GeminiTopK"`
	GeminiSystemInstruction         string  `yaml:"GeminiSystemInstruction"`
	GeminiMaxWaitTimeFileProcessing int     `yaml:"GeminiMaxWaitTimeFileProcessing"`
	//
	MarkdownPromptResponseFile       string `yaml:"MarkdownPromptResponseFile"`
	MarkdownOutput                   bool   `yaml:"MarkdownOutput"`
	MarkdownOutputApplication        string
	MarkdownOutputApplicationMacOS   string `yaml:"MarkdownOutputApplicationMacOS"`
	MarkdownOutputApplicationLinux   string `yaml:"MarkdownOutputApplicationLinux"`
	MarkdownOutputApplicationWindows string `yaml:"MarkdownOutputApplicationWindows"`
	MarkdownOutputApplicationOther   string `yaml:"MarkdownOutputApplicationOther"`
	MarkdownHistory                  bool   `yaml:"MarkdownHistory"`
	MarkdownHistoryDirectory         string `yaml:"MarkdownHistoryDirectory"`
	//
	AnsiRendering          bool                `yaml:"AnsiRendering"`
	AnsiPromptResponseFile string              `yaml:"AnsiPromptResponseFile"`
	AnsiOutput             bool                `yaml:"AnsiOutput"`
	AnsiHistory            bool                `yaml:"AnsiHistory"`
	AnsiHistoryDirectory   string              `yaml:"AnsiHistoryDirectory"`
	AnsiReplaceColors      []map[string]string `yaml:"AnsiReplaceColors"`
	//
	HTMLRendering                bool   `yaml:"HTMLRendering"`
	HTMLPromptResponseFile       string `yaml:"HTMLPromptResponseFile"`
	HTMLOutput                   bool   `yaml:"HTMLOutput"`
	HTMLOutputApplication        string
	HTMLOutputApplicationMacOS   string              `yaml:"HTMLOutputApplicationMacOS"`
	HTMLOutputApplicationLinux   string              `yaml:"HTMLOutputApplicationLinux"`
	HTMLOutputApplicationWindows string              `yaml:"HTMLOutputApplicationWindows"`
	HTMLOutputApplicationOther   string              `yaml:"HTMLOutputApplicationOther"`
	HTMLHistory                  bool                `yaml:"HTMLHistory"`
	HTMLHistoryDirectory         string              `yaml:"HTMLHistoryDirectory"`
	HTMLMaxLengthTitle           int                 `yaml:"HTMLMaxLengthTitle"`
	HTMLReplaceElements          []map[string]string `yaml:"HTMLReplaceElements"`
	HTMLHeader                   string              `yaml:"HTMLHeader"`
	HTMLFooter                   string              `yaml:"HTMLFooter"`
	//
	InputFromTerminal  bool   `yaml:"InputFromTerminal"`
	InputFromFile      bool   `yaml:"InputFromFile"`
	InputFile          string `yaml:"InputFile"`
	InputFromLocalhost bool   `yaml:"InputFromLocalhost"`
	InputLocalhostPort int    `yaml:"InputLocalhostPort"`
	//
	NotifyPrompt                     bool `yaml:"NotifyPrompt"`
	NotifyPromptApplication          string
	NotifyPromptApplicationMacOS     string `yaml:"NotifyPromptApplicationMacOS"`
	NotifyPromptApplicationLinux     string `yaml:"NotifyPromptApplicationLinux"`
	NotifyPromptApplicationWindows   string `yaml:"NotifyPromptApplicationWindows"`
	NotifyPromptApplicationOther     string `yaml:"NotifyPromptApplicationOther"`
	NotifyResponse                   bool   `yaml:"NotifyResponse"`
	NotifyResponseApplication        string
	NotifyResponseApplicationMacOS   string `yaml:"NotifyResponseApplicationMacOS"`
	NotifyResponseApplicationLinux   string `yaml:"NotifyResponseApplicationLinux"`
	NotifyResponseApplicationWindows string `yaml:"NotifyResponseApplicationWindows"`
	NotifyResponseApplicationOther   string `yaml:"NotifyResponseApplicationOther"`
	//
	HistoryFilenameSchema            string `yaml:"HistoryFilenameSchema"`
	HistoryFilenameAddPrefix         bool   `yaml:"HistoryFilenameAddPrefix"`
	HistoryFilenameAddPostfix        bool   `yaml:"HistoryFilenameAddPostfix"`
	HistoryFilenameExtensionMarkdown string `yaml:"HistoryFilenameExtensionMarkdown"`
	HistoryFilenameExtensionAnsi     string `yaml:"HistoryFilenameExtensionAnsi"`
	HistoryFilenameExtensionHTML     string `yaml:"HistoryFilenameExtensionHTML"`
	HistoryMaxFilenameLength         int    `yaml:"HistoryMaxFilenameLength"`
	//
	GeneralInternetProxy string `yaml:"GeneralInternetProxy"`
}

// progConfig contains program configuration
var progConfig = ProgConfig{GeminiCandidateCount: -1, GeminiTemperature: -1.0, GeminiTopP: -1.0, GeminiTopK: -1, GeminiMaxOutputTokens: -1}

/*
loadConfiguration loads program configuration from yaml file.
*/
func loadConfiguration(configFile string) error {
	operatingSystem := runtime.GOOS

	source, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error [%w] reading configuration file", err)
	}
	err = yaml.Unmarshal(source, &progConfig)
	if err != nil {
		return fmt.Errorf("error [%w] unmarshalling configuration file", err)
	}

	// gemini
	if progConfig.GeminiAPIKey == "" {
		return fmt.Errorf("empty GeminiAPIKey not allowed")
	}
	if progConfig.GeminiAiModel == "" {
		return fmt.Errorf("empty GeminiAiModel not allowed")
	}
	if progConfig.GeminiCandidateCount <= 0 {
		return fmt.Errorf("empty GeminiCandidateCount not allowed")
	}

	// markdown
	if progConfig.MarkdownPromptResponseFile == "" {
		return fmt.Errorf("empty MarkdownPromptResponseFile not allowed")
	}
	switch operatingSystem {
	case "darwin":
		progConfig.MarkdownOutputApplication = progConfig.MarkdownOutputApplicationMacOS
	case "linux":
		progConfig.MarkdownOutputApplication = progConfig.MarkdownOutputApplicationLinux
	case "windows":
		progConfig.MarkdownOutputApplication = progConfig.MarkdownOutputApplicationWindows
	default:
		progConfig.MarkdownOutputApplication = progConfig.MarkdownOutputApplicationOther
	}
	if progConfig.MarkdownOutput && progConfig.MarkdownOutputApplication == "" {
		return fmt.Errorf("empty operating system specific MarkdownOutputApplication not allowed")
	}
	if progConfig.MarkdownHistory && progConfig.MarkdownHistoryDirectory == "" {
		return fmt.Errorf("empty MarkdownHistoryDirectory not allowed")
	}

	// ansi
	if progConfig.AnsiRendering && progConfig.AnsiPromptResponseFile == "" {
		return fmt.Errorf("empty AnsiPromptResponseFile not allowed")
	}
	if progConfig.AnsiHistory && progConfig.AnsiHistoryDirectory == "" {
		return fmt.Errorf("empty AnsiHistoryDirectory not allowed")
	}

	// html
	if progConfig.HTMLRendering && progConfig.HTMLPromptResponseFile == "" {
		return fmt.Errorf("empty HTMLPromptResponseFile not allowed")
	}
	switch operatingSystem {
	case "darwin":
		progConfig.HTMLOutputApplication = progConfig.HTMLOutputApplicationMacOS
	case "linux":
		progConfig.HTMLOutputApplication = progConfig.HTMLOutputApplicationLinux
	case "windows":
		progConfig.HTMLOutputApplication = progConfig.HTMLOutputApplicationWindows
	default:
		progConfig.HTMLOutputApplication = progConfig.HTMLOutputApplicationOther
	}
	if progConfig.HTMLOutput && progConfig.HTMLOutputApplication == "" {
		return fmt.Errorf("empty operating system specific HTMLOutputApplication not allowed")
	}
	if progConfig.HTMLHistory && progConfig.HTMLHistoryDirectory == "" {
		return fmt.Errorf("empty HTMLHistoryDirectory not allowed")
	}

	// input
	if progConfig.InputFromFile && progConfig.InputFile == "" {
		return fmt.Errorf("empty InputFile not allowed")
	}

	// notification
	switch operatingSystem {
	case "darwin":
		progConfig.NotifyPromptApplication = progConfig.NotifyPromptApplicationMacOS
	case "linux":
		progConfig.NotifyPromptApplication = progConfig.NotifyPromptApplicationLinux
	case "windows":
		progConfig.NotifyPromptApplication = progConfig.NotifyPromptApplicationWindows
	default:
		progConfig.NotifyPromptApplication = progConfig.NotifyPromptApplicationOther
	}
	if progConfig.NotifyPrompt && progConfig.NotifyPromptApplication == "" {
		return fmt.Errorf("empty operating system specific NotifyPromptApplication not allowed")
	}
	switch operatingSystem {
	case "darwin":
		progConfig.NotifyResponseApplication = progConfig.NotifyResponseApplicationMacOS
	case "linux":
		progConfig.NotifyResponseApplication = progConfig.NotifyResponseApplicationLinux
	case "windows":
		progConfig.NotifyResponseApplication = progConfig.NotifyResponseApplicationWindows
	default:
		progConfig.NotifyResponseApplication = progConfig.NotifyResponseApplicationOther
	}
	if progConfig.NotifyResponse && progConfig.NotifyResponseApplication == "" {
		return fmt.Errorf("empty operating system specific NotifyResponseApplication not allowed")
	}

	// history
	progConfig.HistoryFilenameSchema = strings.ToLower(progConfig.HistoryFilenameSchema)
	switch progConfig.HistoryFilenameSchema {
	case "timestamp":
	case "prompt":
	default:
		return fmt.Errorf("unsupported history filename schema")
	}
	if progConfig.HistoryMaxFilenameLength > 255 {
		return fmt.Errorf("max length of history filename show not be greater than 255")
	}

	// get api-key (password)
	progConfig.GeminiAPIKey, err = getPassword(progConfig.GeminiAPIKey)
	if err != nil {
		return fmt.Errorf("error [%w] getting api-key", err)
	}

	// get internet proxy (password)
	if progConfig.GeneralInternetProxy != "" {
		progConfig.GeneralInternetProxy, err = getPassword(progConfig.GeneralInternetProxy)
		if err != nil {
			return fmt.Errorf("error [%w] getting internet proxy", err)
		}
	}

	return nil
}

/*
showConfiguration shows loaded configuration.
*/
func showConfiguration() {
	// general notes
	fmt.Printf("\nNotes concerning the freely available version of 'Google Gemini AI':\n")
	fmt.Printf("  See the help page for the 'Google Gemini AI' terms of service.\n")
	fmt.Printf("  All input data will be used by Google to improve 'Gemini AI'.\n")
	fmt.Printf("  Therefore, do not process any private or confidential data.\n")

	fmt.Printf("\nInput from:\n")
	if progConfig.InputFromTerminal {
		fmt.Printf("  Terminal  : yes\n")
	}
	if progConfig.InputFromFile {
		fmt.Printf("  File      : %v\n", progConfig.InputFile)
	}
	if progConfig.InputFromLocalhost {
		fmt.Printf("  localhost : %v (port)\n", progConfig.InputLocalhostPort)
	}

	fmt.Printf("\nRendering:\n")
	fmt.Printf("  Markdown : %v\n", progConfig.MarkdownPromptResponseFile)
	if progConfig.AnsiRendering {
		fmt.Printf("  Ansi     : %v\n", progConfig.AnsiPromptResponseFile)
	}
	if progConfig.HTMLRendering {
		fmt.Printf("  HTML     : %v\n", progConfig.HTMLPromptResponseFile)
	}

	fmt.Printf("\nHistory:\n")
	if progConfig.MarkdownHistory {
		fmt.Printf("  Markdown : %v\n", progConfig.MarkdownHistoryDirectory)
	}
	if progConfig.AnsiHistory {
		fmt.Printf("  Ansi     : %v\n", progConfig.AnsiHistoryDirectory)
	}
	if progConfig.HTMLHistory {
		fmt.Printf("  HTML     : %v\n", progConfig.HTMLHistoryDirectory)
	}

	fmt.Printf("\nOutput:\n")
	if progConfig.AnsiOutput {
		fmt.Printf("  Terminal : yes\n")
	}
	if progConfig.MarkdownOutput {
		fmt.Printf("  Markdown : execute application\n")
	}
	if progConfig.HTMLOutput {
		fmt.Printf("  HTML     : execute application\n")
	}
}

/*
initializeProgram initializes this program.
*/
func initializeProgram() {
	var err error

	// create history directories
	if progConfig.MarkdownHistory {
		err = os.Mkdir(progConfig.MarkdownHistoryDirectory, 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("error [%v] at os.Mkdir()\n", err)
			os.Exit(1)
		}
	}
	if progConfig.AnsiHistory {
		err = os.Mkdir(progConfig.AnsiHistoryDirectory, 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("error [%v] at os.Mkdir()\n", err)
			os.Exit(1)
		}
	}
	if progConfig.HTMLHistory {
		err = os.Mkdir(progConfig.HTMLHistoryDirectory, 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("error [%v] at os.Mkdir()\n", err)
			os.Exit(1)
		}
		err = os.Mkdir(progConfig.HTMLHistoryDirectory+"/assets", 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("error [%v] at os.Mkdir()\n", err)
			os.Exit(1)
		}
		writeAssets(progConfig.HTMLHistoryDirectory)
	}
}
