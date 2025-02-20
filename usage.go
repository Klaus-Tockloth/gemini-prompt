package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

/*
printUsage prints the usage of this program.
*/
func printUsage() {
	fmt.Printf("\nUsage:\n")
	fmt.Printf("  %s [options] [files]\n", progName)

	fmt.Printf("\nExamples:\n")
	fmt.Printf("  %s\n", progName)
	fmt.Printf("  %s *.go README.md\n", progName)
	fmt.Printf("  %s -dryrun -uploads ganymed-project-files.txt\n", progName)

	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()

	fmt.Printf("\nNotes:\n")
	fmt.Printf("  - This program allows you to prompt 'Google Gemini AI' and\n")
	fmt.Printf("    seamlessly integrate its capabilities into your workflow.\n")
	fmt.Printf("  - You can submit prompts via the following input channels:\n")
	fmt.Printf("    Terminal, File, localhost\n")
	fmt.Printf("  - Output is available in the following formats:\n")
	fmt.Printf("    Markdown (Editor), HTML (Browser), Ansi (Terminal)\n")
	fmt.Printf("  - Each prompt is self-contained (no chat).\n")
	fmt.Printf("  - Specified files are transmitted to 'Google Gemini AI',\n")
	fmt.Printf("    allowing prompts to reference their contents.\n")
	fmt.Printf("  - The program offers many configuration options.\n")
	fmt.Printf("  - The presentation of the outputs can be customized.\n")

	fmt.Printf("\nDisclaimer:\n")
	fmt.Printf("  This application is for evaluating the concept of integrating and using AI in\n")
	fmt.Printf("  a personalized work environment. All v0.x versions require a Gemini API key,\n")
	fmt.Printf("  enabling free and limited use of 'Google Gemini AI'.\n")
	fmt.Printf("  A Gemini API key is associated with a personal Google account, allowing trial\n")
	fmt.Printf("  use of 'Google Gemini AI'. The Gemini API key is not intended for permanent or\n")
	fmt.Printf("  extensive use. From v1.0 onwards, the application will switch to using a regular\n")
	fmt.Printf("  Google Gemini account.\n")

	fmt.Printf("\nTerms of service apply to 'Google Gemini AI':\n")
	fmt.Printf("  The Google Terms of Service (policies.google.com/terms) and the Generative AI\n")
	fmt.Printf("  Prohibited Use Policy (policies.google.com/terms/generative-ai/use-policy) apply\n")
	fmt.Printf("  to 'Google Gemini AI' service. Visit the Gemini Apps Privacy Hub\n")
	fmt.Printf("  (support.google.com/gemini?p=privacy_help) to learn more about how Google uses\n")
	fmt.Printf("  your Gemini Apps data. See also the Gemini Apps FAQ (gemini.google.com/faq).\n")

	fmt.Printf("\nNotes concerning the freely available version of 'Google Gemini AI'\n")
	fmt.Printf("  - All input data will be used by Google to improve 'Gemini AI'.\n")
	fmt.Printf("  - Therefore, do not process any private or confidential data.\n")

	fmt.Printf("\nRequired:\n")
	fmt.Printf("  - Obtain a personal 'Gemini API Key' from Google.\n")
	fmt.Printf("  - Configure the API Key in your program environment:\n")
	fmt.Printf("    macOS, Linux : export GEMINI_API_KEY=Your-API-Key\n")
	fmt.Printf("    Windows      : set GEMINI_API_KEY Your-API-Key\n")

	fmt.Printf("\nTip:\n")
	fmt.Printf("  In practice, a browser is useful for both creating prompts and presenting\n")
	fmt.Printf("  the output. The simple 'prompt-input.html' webpage can be used for creating\n")
	fmt.Printf("  and sending prompts to 'localhost'.\n")

	fmt.Printf("\n")
	os.Exit(1)
}

/*
showAvailableGeminiModels requests and shows all currently available and possible Gemini AI models.
*/
func showAvailableGeminiModels() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(progConfig.GeminiAPIKey))
	if err != nil {
		fmt.Printf("error [%v] creating AI client\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nPossible Gemini AI models (currently determined):\n")

	iter := client.ListModels(ctx)
	for {
		modelInfo, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("error [%v] iterating AI models\n", err)
			break
		}

		// only show models that support 'generateContent'
		hasGenerateContent := false
		for _, method := range modelInfo.SupportedGenerationMethods {
			if method == "generateContent" {
				hasGenerateContent = true
				break
			}
		}
		if !hasGenerateContent {
			continue
		}

		fmt.Printf("\n%v (version: %v, in: %v, out: %v)\n", strings.TrimPrefix(modelInfo.Name, "models/"),
			modelInfo.Version, modelInfo.InputTokenLimit, modelInfo.OutputTokenLimit)
		fmt.Printf("%v\n", modelInfo.Description)
	}
	client.Close()

	fmt.Printf("\n")
}
