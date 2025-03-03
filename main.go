/*
Purpose:
- gemini prompt

Description:
- Prompt Google Gemini AI and display the response.

Releases:
  - v0.1.0 - 2025/02/20: initial release
  - v0.1.1 - 2025/02/23: fixed: nil pointer dereference in processResponse()
  - v0.2.0 - 2025/02/24: added: 'system instruction' to prompt output, internet proxy support
  - v0.3.0 - 2025/03/02: options '-topk' and 'topp' added, general improvements, refactoring
  - v0.3.1 - 2025/03/03: fixed: nil pointer dereference in printAIModelInfo()
    method filter removed in showAvailableGeminiModels()

Copyright:
- © 2025 | Klaus Tockloth

License:
- MIT License

Contact:
- klaus.tockloth@googlemail.com

Remarks:
- none

Links:
- https://pkg.go.dev/github.com/google/generative-ai-go/genai
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/aquilax/truncate"
	"github.com/google/generative-ai-go/genai"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/term"
	"google.golang.org/api/option"
)

// general program info
var (
	progName    = strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(filepath.Base(os.Args[0])))
	progVersion = "v0.3.1"
	progDate    = "2025/03/03"
	progPurpose = "gemini prompt"
	progInfo    = "Prompt Google Gemini AI and display the response."
)

// processing timestamp
var (
	startProcessing  time.Time
	finishProcessing time.Time
)

// markdown to html parser
var markdownParser goldmark.Markdown

// gemini AI model information
var modelInfo *genai.ModelInfo

// files uploaded to gemini AI model
var uploadedFiles []*genai.File

/*
main starts this program.
*/
func main() {
	var err error

	fmt.Printf("\nProgram:\n")
	fmt.Printf("  Name    : %s\n", progName)
	fmt.Printf("  Release : %s - %s\n", progVersion, progDate)
	fmt.Printf("  Purpose : %s\n", progPurpose)
	fmt.Printf("  Info    : %s\n", progInfo)

	// request terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("error [%v] at term.GetSize()", err)
		os.Exit(1)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Printf("error [%v] getting working directory\n", err)
		os.Exit(1)
	}

	candidates := flag.Int("candidates", -1, "specifies number of AI responses (overwrites YAML config)")
	temperature := flag.Float64("temperature", -1.0, "specifies variation range of AI responses (overwrites YAML config)")
	topp := flag.Float64("topp", -1.0, "maximum cumulative probability of tokens to consider when sampling (overwrites YAML config)")
	topk := flag.Int("topk", -1, "maximum number of tokens to consider when sampling (overwrites YAML config)")
	maxtokens := flag.Int("maxtokens", -1, "max output tokens (useful to force short content, overwrites YAML config)")
	dryrun := flag.Bool("dryrun", false, "only print list of files given via command line")
	uploads := flag.String("uploads", "", "name of list with files to upload to AI (one file per line)")
	dir, _ := filepath.Split(os.Args[0])
	defaultConfigFile := dir + progName + ".yaml"
	config := flag.String("config", defaultConfigFile, "name of YAML config file")
	models := flag.Bool("models", false, "show all AI Gemini models and terminate")

	flag.Usage = printUsage
	flag.Parse()

	if !fileExists(*config) {
		writeConfig()
	}
	if !fileExists("./assets") {
		err = os.Mkdir("./assets", 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("error [%v] at os.Mkdir()\n", err)
			os.Exit(1)
		}
		writeAssets(".")
	}

	if !fileExists("./prompt-input.html") {
		writePromptInput()
	}

	err = loadConfiguration(*config)
	if err != nil {
		fmt.Printf("error [%v] loading configuration\n", err)
		os.Exit(1)
	}

	if *models {
		showAvailableGeminiModels(terminalWidth)
		os.Exit(1)
	}

	var uploadFiles []string
	if *uploads != "" {
		uploadFiles, err = slurpFile(*uploads)
		if err != nil {
			fmt.Printf("error [%v] reading list of files to upload to AI\n", err)
		}
	}
	commandlineFiles := flag.Args()
	allFiles := uploadFiles
	allFiles = append(allFiles, commandlineFiles...)

	if *dryrun {
		fmt.Printf("\nFiles given via command line:\n")
		for _, file := range allFiles {
			mimeType, err := getMimeType(file)
			info := "ok"
			if err != nil {
				info = "error"
			}
			if mimeType == "application/octet-stream" {
				info = "warn"
			}
			if err == nil {
				fmt.Printf("  %-5s  %-32.32s  %s\n", info, mimeType, file)
			} else {
				fmt.Printf("  %-5s  %s\n", info, err)
			}
		}
		fmt.Printf("\n")
		return
	}

	// show configuration
	showConfiguration()

	// initialize this program
	initializeProgram()

	// overwrite YAML config values with cli parameters
	if *candidates > 0 {
		progConfig.GeminiCandidateCount = int32(*candidates)
	}
	if *temperature > -1.0 {
		progConfig.GeminiTemperature = float32(*temperature)
	}
	if *topp > -1.0 {
		progConfig.GeminiTopP = float32(*topp)
	}
	if *topk > -1 {
		progConfig.GeminiTopK = int32(*topk)
	}
	if *maxtokens > 0 {
		progConfig.GeminiMaxOutputTokens = int32(*maxtokens)
	}

	// create markdown parser
	markdownParser = goldmark.New(goldmark.WithExtensions(extension.GFM))

	// create AI client
	var client *genai.Client
	ctx := context.Background()
	if progConfig.GeneralInternetProxy != "" {
		// indirect internet connection: client -> proxy -> internet
		httpClient := &http.Client{Transport: &ProxyRoundTripper{
			APIKey:   progConfig.GeminiAPIKey,
			ProxyURL: progConfig.GeneralInternetProxy,
		}}
		// option.WithAPIKey() shouldn't be necessary because the key is set in ProxyRoundTripper
		// but without the option, NewClient() attempts to authenticate via Google Cloud SDK (ADC)
		client, err = genai.NewClient(ctx, option.WithAPIKey(progConfig.GeminiAPIKey), option.WithHTTPClient(httpClient))
	} else {
		// direct internet connection: client -> internet
		client, err = genai.NewClient(ctx, option.WithAPIKey(progConfig.GeminiAPIKey))
	}
	if err != nil {
		fmt.Printf("error [%v] creating AI client\n", err)
		os.Exit(1)
	}

	// upload all files given from command line
	uploadedFiles, err = uploadFilesToGemini(ctx, client, allFiles)
	if err != nil {
		fmt.Printf("error [%v] uploading files\n", err)
		return
	}

	// define Gemini AI model
	geminiModel := client.GenerativeModel(progConfig.GeminiAiModel)

	// get model infos
	modelInfo, err = geminiModel.Info(ctx)
	if err != nil {
		fmt.Printf("error [%v] getting AI model information\n", err)
		return
	}

	// configure AI model parameters
	if progConfig.GeminiCandidateCount > -1 {
		geminiModel.SetCandidateCount(progConfig.GeminiCandidateCount)
	}
	if progConfig.GeminiTemperature > -1.0 {
		geminiModel.SetTemperature(progConfig.GeminiTemperature)
	}
	if progConfig.GeminiTopP > -1.0 {
		geminiModel.SetTopP(progConfig.GeminiTopP)
	}
	if progConfig.GeminiTopK > -1 {
		geminiModel.SetTopK(progConfig.GeminiTopK)
	}
	if progConfig.GeminiMaxOutputTokens > -1 {
		geminiModel.SetMaxOutputTokens(progConfig.GeminiMaxOutputTokens)
	}
	if progConfig.GeminiSystemInstruction != "" {
		geminiModel.SystemInstruction = genai.NewUserContent(genai.Text(progConfig.GeminiSystemInstruction))
	}

	// print AI model information
	printAIModelInfo(geminiModel, modelInfo, terminalWidth)

	// define prompt channel
	promptChannel := make(chan string)

	// set up signal handling for shutdown (e.g. Ctrl-C)
	shutdownTrigger := make(chan os.Signal, 1)
	signal.Notify(shutdownTrigger, syscall.SIGINT)  // kill -SIGINT pid -> interrupt
	signal.Notify(shutdownTrigger, syscall.SIGTERM) // kill -SIGTERM pid -> terminated

	fmt.Printf("\nProgram termination:\n")
	fmt.Printf("  Press CTRL-C to terminate this program.\n\n")

	// start graceful shutdown handler
	go handleShutdown(ctx, shutdownTrigger, client, uploadedFiles)

	// start input readers
	inputPossibilities := startInputReaders(promptChannel, progConfig)

	// main loop: 'Prompt Google Gemini AI'
	for {
		fmt.Printf("Waiting for input from %s ...\n", strings.Join(inputPossibilities, ", "))

		// read prompt from channel
		prompt := <-promptChannel
		prompt = strings.TrimSpace(prompt)

		now := time.Now()
		if progConfig.NotifyPrompt {
			_ = runCommand(progConfig.NotifyPromptApplication)
		}
		fmt.Printf("%02d:%02d:%02d: Processing prompt ...\n", now.Hour(), now.Minute(), now.Second())
		processPrompt(prompt)

		// build prompt with all parts (files and text)
		promptParts := []genai.Part{}
		for _, uploadedFile := range uploadedFiles {
			promptParts = append(promptParts, genai.FileData{URI: uploadedFile.URI})
		}
		promptParts = append(promptParts, genai.Text(prompt))

		// generate content
		startProcessing = time.Now()
		var resp *genai.GenerateContentResponse
		resp, err = geminiModel.GenerateContent(ctx, promptParts...)
		if err != nil {
			fmt.Printf("error [%v] generating content\n", err)
		}
		finishProcessing = time.Now()

		now = finishProcessing
		fmt.Printf("%02d:%02d:%02d: Processing response ...\n", now.Hour(), now.Minute(), now.Second())
		processResponse(resp, err)

		// trigger response notification
		if progConfig.NotifyResponse {
			_ = runCommand(progConfig.NotifyResponseApplication)
		}

		// print prompt and response to terminal
		if progConfig.AnsiOutput {
			printPromptResponseToTerminal()
		}

		// copy ansi file to history
		if progConfig.AnsiHistory {
			ansiDestinationFile := buildDestinationFilename(now, prompt, progConfig.HistoryFilenameExtensionAnsi)
			ansiDestinationPathFile := filepath.Join(workingDirectory, progConfig.AnsiHistoryDirectory, ansiDestinationFile)
			copyFile(progConfig.AnsiPromptResponseFile, ansiDestinationPathFile)
		}

		// markdown prompt and response file: nothing to do
		commandLine := fmt.Sprintf(progConfig.MarkdownOutputApplication, progConfig.MarkdownPromptResponseFile)

		// copy markdown file to history
		if progConfig.MarkdownHistory {
			markdownDestinationFile := buildDestinationFilename(now, prompt, progConfig.HistoryFilenameExtensionMarkdown)
			markdownDestinationPathFile := filepath.Join(workingDirectory, progConfig.MarkdownHistoryDirectory, markdownDestinationFile)
			copyFile(progConfig.MarkdownPromptResponseFile, markdownDestinationPathFile)
			commandLine = fmt.Sprintf(progConfig.MarkdownOutputApplication, "\""+markdownDestinationPathFile+"\"")
		}

		// open markdown document in application
		if progConfig.MarkdownOutput {
			err := runCommand(commandLine)
			if err != nil {
				fmt.Printf("error [%v] at runCommand()\n", err)
			}
		}

		// build prompt and response html page
		commandLine = fmt.Sprintf(progConfig.HTMLOutputApplication, progConfig.HTMLPromptResponseFile)
		_ = buildHTMLPage(prompt, progConfig.HTMLPromptResponseFile, progConfig.HTMLPromptResponseFile)

		// copy html file to history
		if progConfig.HTMLHistory {
			htmlDestinationFile := buildDestinationFilename(now, prompt, progConfig.HistoryFilenameExtensionHTML)
			htmlDestinationPathFile := filepath.Join(workingDirectory, progConfig.HTMLHistoryDirectory, htmlDestinationFile)
			copyFile(progConfig.HTMLPromptResponseFile, htmlDestinationPathFile)
			commandLine = fmt.Sprintf(progConfig.HTMLOutputApplication, "\""+htmlDestinationPathFile+"\"")
		}

		// open html page in application
		if progConfig.HTMLOutput {
			err := runCommand(commandLine)
			if err != nil {
				fmt.Printf("error [%v] at runCommand()\n", err)
			}
		}
	}
}

/*
printAIModelInfo prints AI model information to the console.
*/
func printAIModelInfo(geminiModel *genai.GenerativeModel, modelInfo *genai.ModelInfo, terminalWidth int) {
	// calculate words from tokens
	inputTokenLimitWordsLower := float64(modelInfo.InputTokenLimit) / 100.0 * 60.0
	inputTokenLimitWordsLower = math.Floor(inputTokenLimitWordsLower/100.0) * 100.0
	inputTokenLimitWordsUpper := float64(modelInfo.InputTokenLimit) / 100.0 * 80.0
	inputTokenLimitWordsUpper = math.Floor(inputTokenLimitWordsUpper/100.0) * 100.0
	outputTokenLimitWordsLower := float64(modelInfo.OutputTokenLimit) / 100.0 * 60.0
	outputTokenLimitWordsLower = math.Floor(outputTokenLimitWordsLower/100.0) * 100.0
	outputTokenLimitWordsUpper := float64(modelInfo.OutputTokenLimit) / 100.0 * 80.0
	outputTokenLimitWordsUpper = math.Floor(outputTokenLimitWordsUpper/100.0) * 100.0

	fmt.Printf("\nAI model information:\n")
	fmt.Printf("  Name              : %v\n", modelInfo.Name)
	fmt.Printf("  BaseModelID       : %v\n", modelInfo.BaseModelID)
	fmt.Printf("  Version           : %v\n", modelInfo.Version)
	fmt.Printf("  DisplayName       : %v\n", modelInfo.DisplayName)
	fmt.Printf("  Description       : %v\n", wrapString(modelInfo.Description, terminalWidth, 22))
	fmt.Printf("  InputTokenLimit   : %v (approx. %.0f-%.0f english words)\n", modelInfo.InputTokenLimit, inputTokenLimitWordsLower, inputTokenLimitWordsUpper)
	fmt.Printf("  OutputTokenLimit  : %v (approx. %.0f-%.0f english words)\n", modelInfo.OutputTokenLimit, outputTokenLimitWordsLower, outputTokenLimitWordsUpper)
	fmt.Printf("  Supported Methods : %v\n", strings.Join(modelInfo.SupportedGenerationMethods, ", "))
	fmt.Printf("  Temperature       : %v\n", modelInfo.Temperature)
	if modelInfo.MaxTemperature != nil {
		fmt.Printf("  MaxTemperature    : %v\n", *modelInfo.MaxTemperature)
	}
	fmt.Printf("  TopP              : %v\n", modelInfo.TopP)
	fmt.Printf("  TopK              : %v\n", modelInfo.TopK)

	fmt.Printf("\nUser defined AI model configuration:\n")
	if geminiModel.GenerationConfig.CandidateCount != nil {
		fmt.Printf("  CandidateCount    : %v\n", *geminiModel.GenerationConfig.CandidateCount)
	}
	if geminiModel.GenerationConfig.MaxOutputTokens != nil {
		fmt.Printf("  MaxOutputTokens   : %v\n", *geminiModel.GenerationConfig.MaxOutputTokens)
	}
	if geminiModel.GenerationConfig.Temperature != nil {
		fmt.Printf("  Temperature       : %v\n", *geminiModel.GenerationConfig.Temperature)
	}
	if geminiModel.GenerationConfig.TopP != nil {
		fmt.Printf("  TopP              : %v\n", *geminiModel.GenerationConfig.TopP)
	}
	if geminiModel.GenerationConfig.TopK != nil {
		fmt.Printf("  TopK              : %v\n", *geminiModel.GenerationConfig.TopK)
	}
	if progConfig.GeminiSystemInstruction != "" {
		truncatedSystemInstruction := truncate.Truncate(progConfig.GeminiSystemInstruction, 96, "...", truncate.PositionMiddle)
		fmt.Printf("  SystemInstruction : %v\n", truncatedSystemInstruction)
	}
}

/*
handleShutdown handles program termination signals.
*/
func handleShutdown(ctx context.Context, shutdownTrigger chan os.Signal, client *genai.Client, uploadedFiles []*genai.File) {
	<-shutdownTrigger
	fmt.Printf("\nShutdown signal received. Exiting gracefully ...\n")

	// cleanup/delete all uploaded files before program termination
	for _, uploadedFile := range uploadedFiles {
		err := client.DeleteFile(ctx, uploadedFile.Name)
		fmt.Printf("deleting uploaded remote file [%v]\n", uploadedFile.DisplayName)
		if err != nil {
			fmt.Printf("error [%v] deleting uploaded file\n", err)
		}
	}

	fmt.Printf("Closing Gemini AI client ...\n")
	err := client.Close()
	if err != nil {
		fmt.Printf("error [%v] closing Gemini AI client\n", err)
	}

	fmt.Printf("Done\n")
	os.Exit(0)
}

/*
startInputReaders starts input readers based on the configuration.
*/
func startInputReaders(promptChannel chan string, config ProgConfig) []string {
	inputPossibilities := []string{}

	// input from keyboard
	if config.InputFromTerminal {
		go readPromptFromKeyboard(promptChannel)
		inputPossibilities = append(inputPossibilities, "Terminal")
	}

	// input from file
	if config.InputFromFile {
		if !fileExists(config.InputFile) {
			file, err := os.Create(config.InputFile)
			if err != nil {
				fmt.Printf("error [%v] creating input prompt text file\n", err)
				return inputPossibilities
			}
			file.Close()
		}
		go readPromptFromFile(config.InputFile, promptChannel)
		inputPossibilities = append(inputPossibilities, "File")
	}

	// input from localhost
	if config.InputFromLocalhost {
		addr := fmt.Sprintf("localhost:%d", config.InputLocalhostPort)
		go func() {
			http.HandleFunc("/", readPromptFromLocalhost(promptChannel))
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				fmt.Printf("error [%v] starting internal webserver\n", err)
				return
			}
		}()
		inputPossibilities = append(inputPossibilities, addr)
	}

	return inputPossibilities
}
