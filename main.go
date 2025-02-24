/*
Purpose:
- gemini prompt

Description:
- Prompt Google Gemini AI and display the response.

Releases:
- v0.1.0 - 2025/02/20: initial release
- v0.1.1 - 2025/02/23: fixed: nil pointer dereference in processResponse()
- v0.2.0 - 2025/02/24: added: 'system instruction' to prompt output, internet proxy support

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
	progVersion = "v0.2.0"
	progDate    = "2025/02/24"
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

	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Printf("error [%v] getting working directory\n", err)
		os.Exit(1)
	}

	candidates := flag.Int("candidates", -1, "specifies number of AI responses (overwrites YAML config)")
	temperature := flag.Float64("temperature", -1.0, "specifies variation range of AI responses (overwrites YAML config)")
	dryrun := flag.Bool("dryrun", false, "only print list of files given via command line")
	uploads := flag.String("uploads", "", "name of list with files to upload to AI (one file per line)")
	dir, _ := filepath.Split(os.Args[0])
	defaultConfigFile := dir + progName + ".yaml"
	config := flag.String("config", defaultConfigFile, "name of YAML config file")
	models := flag.Bool("models", false, "show all AI Gemini models and terminate")
	maxtokens := flag.Int("maxtokens", -1, "max output tokens (useful to force short content)")

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
		showAvailableGeminiModels()
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

	showConfiguration()

	initializeProgram()

	// overwrite YAML config values with cli parameters
	if *candidates > 0 {
		progConfig.GeminiCandidateCount = int32(*candidates)
	}
	if *temperature > -1.0 {
		progConfig.GeminiTemperature = float32(*temperature)
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

	// request terminal width
	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("error [%v] at term.GetSize()", err)
		os.Exit(1)
	}

	// configure AI model parameters
	if progConfig.GeminiCandidateCount > -1 {
		geminiModel.SetCandidateCount(progConfig.GeminiCandidateCount)
	}
	if progConfig.GeminiMaxOutputTokens > -1 {
		geminiModel.SetMaxOutputTokens(progConfig.GeminiMaxOutputTokens)
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
	if progConfig.GeminiSystemInstruction != "" {
		geminiModel.SystemInstruction = genai.NewUserContent(genai.Text(progConfig.GeminiSystemInstruction))
	}

	fmt.Printf("\nAI model information:\n")
	fmt.Printf("  Name              : %v\n", modelInfo.Name)
	fmt.Printf("  BaseModelID       : %v\n", modelInfo.BaseModelID)
	fmt.Printf("  Version           : %v\n", modelInfo.Version)
	fmt.Printf("  DisplayName       : %v\n", modelInfo.DisplayName)
	fmt.Printf("  Description       : %v\n", wrapString(modelInfo.Description, terminalWidth, 22))
	fmt.Printf("  InputTokenLimit   : %v\n", modelInfo.InputTokenLimit)
	fmt.Printf("  OutputTokenLimit  : %v\n", modelInfo.OutputTokenLimit)
	fmt.Printf("  Supported Methods : %v\n", strings.Join(modelInfo.SupportedGenerationMethods, ", "))
	fmt.Printf("  Temperature       : %v\n", modelInfo.Temperature)
	fmt.Printf("  MaxTemperature    : %v\n", *modelInfo.MaxTemperature)
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

	// define prompt channel
	promptChannel := make(chan string)

	// set up signal handling for shutdown (e.g. Ctrl-C)
	shutdownTrigger := make(chan os.Signal, 1)
	signal.Notify(shutdownTrigger, syscall.SIGINT)  // kill -SIGINT pid -> interrupt
	signal.Notify(shutdownTrigger, syscall.SIGTERM) // kill -SIGTERM pid -> terminated

	fmt.Printf("\nProgram termination:\n")
	fmt.Printf("  Press CTRL-C to terminate this program.\n\n")

	go func() {
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
	}()

	// list of input possibilities
	inputPossibilities := []string{}

	// input from keyboard
	if progConfig.InputFromTerminal {
		go readPromptFromKeyboard(promptChannel)
		inputPossibilities = append(inputPossibilities, "Terminal")
	}

	// input from file
	if progConfig.InputFromFile {
		if !fileExists(progConfig.InputFile) {
			file, err := os.Create(progConfig.InputFile)
			if err != nil {
				fmt.Printf("error [%v] creating input prompt text file\n", err)
				return
			}
			file.Close()
		}
		go readPromptFromFile(progConfig.InputFile, promptChannel)
		inputPossibilities = append(inputPossibilities, "File")
	}

	// input from localhost
	if progConfig.InputFromLocalhost {
		addr := fmt.Sprintf("localhost:%d", progConfig.InputLocalhostPort)
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

	// prompt Google Gemini AI
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

		// build prompt with all parts (text and files)
		promptParts := []genai.Part{}
		promptParts = append(promptParts, genai.Text(prompt))
		for _, uploadedFile := range uploadedFiles {
			promptParts = append(promptParts, genai.FileData{URI: uploadedFile.URI})
		}

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
		commandLine = fmt.Sprintf(progConfig.HtmlOutputApplication, progConfig.HtmlPromptResponseFile)
		_ = buildHtmlPage(prompt, progConfig.HtmlPromptResponseFile, progConfig.HtmlPromptResponseFile)

		// copy html file to history
		if progConfig.HtmlHistory {
			htmlDestinationFile := buildDestinationFilename(now, prompt, progConfig.HistoryFilenameExtensionHtml)
			htmlDestinationPathFile := filepath.Join(workingDirectory, progConfig.HtmlHistoryDirectory, htmlDestinationFile)
			copyFile(progConfig.HtmlPromptResponseFile, htmlDestinationPathFile)
			commandLine = fmt.Sprintf(progConfig.HtmlOutputApplication, "\""+htmlDestinationPathFile+"\"")
		}

		// open html page in application
		if progConfig.HtmlOutput {
			err := runCommand(commandLine)
			if err != nil {
				fmt.Printf("error [%v] at runCommand()\n", err)
			}
		}
	}
}
