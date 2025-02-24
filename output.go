package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
)

/*
printPromptResponseToTerminal prints prompt and response to ansi terminal.
*/
func printPromptResponseToTerminal() {
	data, err := os.ReadFile(progConfig.AnsiPromptResponseFile)
	if err != nil {
		fmt.Printf("error [%v] at os.ReadFile()\n", err)
		return
	}
	os.Stdout.Write(data)
}

/*
processPrompt processes (user input) prompt.
*/
func processPrompt(prompt string) {
	var promptString strings.Builder

	// text part of prompt
	promptString.WriteString("***\n")
	promptString.WriteString("**Prompt to Gemini:**\n")
	promptString.WriteString("\n```plaintext\n")
	promptString.WriteString(prompt)
	promptString.WriteString("\n```\n")
	promptString.WriteString("\n***\n")

	// system instructions part of prompt
	if progConfig.GeminiSystemInstruction != "" {
		promptString.WriteString("**System Instruction to Gemini:**\n")
		promptString.WriteString("\n```plaintext\n")
		promptString.WriteString(progConfig.GeminiSystemInstruction)
		promptString.WriteString("\n```\n")
		promptString.WriteString("\n***\n")
	}

	// data part of prompt
	if len(uploadedFiles) > 0 {
		promptString.WriteString("**Data referenced by the Prompt:**\n")
		promptString.WriteString("\n```plaintext\n")
		for _, uploadedFile := range uploadedFiles {
			if uploadedFile.State == genai.FileStateActive {
				promptString.WriteString(fmt.Sprintf("%s (%s, %.1f KiB, %s)\n",
					uploadedFile.DisplayName, uploadedFile.UpdateTime.Format("20060102-150405"),
					float64(uploadedFile.SizeBytes)/1024.0, uploadedFile.MIMEType))
			}
		}
		promptString.WriteString("```\n")
		promptString.WriteString("\n***\n")
	}

	// write prompt to current markdown request/response file
	err := os.WriteFile(progConfig.MarkdownPromptResponseFile, []byte(promptString.String()), 0666)
	if err != nil {
		fmt.Printf("error [%v] at os.WriteFile()\n", err)
		return
	}

	// render prompt as ansi
	ansiData := promptString.String()
	if progConfig.AnsiRendering {
		ansiData = renderMarkdown2Ansi(promptString.String())
	}

	// write prompt to current ansi request/response file
	err = os.WriteFile(progConfig.AnsiPromptResponseFile, []byte(ansiData), 0666)
	if err != nil {
		fmt.Printf("error [%v] at os.WriteFile()\n", err)
		return
	}

	// render prompt as html
	htmlData := promptString.String()
	if progConfig.HtmlRendering {
		htmlData = renderMarkdown2Html(promptString.String())
	}

	// write prompt to current html request/response file
	err = os.WriteFile(progConfig.HtmlPromptResponseFile, []byte(htmlData), 0666)
	if err != nil {
		fmt.Printf("error [%v] at os.WriteFile()\n", err)
		return
	}
}

/*
processResponse processes response from AI model.
*/
func processResponse(resp *genai.GenerateContentResponse, err error) {
	var responseString strings.Builder

	if err == nil {
		// print response candidate(s)
		for i, candidate := range resp.Candidates {
			if len(resp.Candidates) > 1 {
				responseString.WriteString(fmt.Sprintf("**Response from Gemini (Candidate #%d):**\n\n", (i + 1)))
			} else {
				responseString.WriteString("**Response from Gemini:**\n\n")
			}
			if candidate.Content == nil {
				responseString.WriteString("No content available in this candidate.\n")
				continue
			}
			for j, part := range candidate.Content.Parts {
				if len(candidate.Content.Parts) > 1 {
					responseString.WriteString(fmt.Sprintf("\nPart #%d:\n", j+1))
				}
				switch p := part.(type) {
				case genai.Text:
					responseString.WriteString(fmt.Sprintf("%s\n", p))
				case genai.FileData:
					responseString.WriteString(fmt.Sprintf("File Data: URI=%s, MIME=%s\n", p.URI, p.MIMEType))
				default:
					responseString.WriteString(fmt.Sprintf("Unsupported part type: %T\n", part))
				}
			}
			responseString.WriteString("\n")

			// build list of text citation source URIs
			citationSourceURIs := []string{}
			if candidate.CitationMetadata != nil {
				for _, citationSource := range candidate.CitationMetadata.CitationSources {
					if citationSource.URI != nil {
						citationSourceURIs = append(citationSourceURIs, (fmt.Sprintf("%v", *citationSource.URI)))
					}
				}
			}

			// show text citation source URIs
			if len(citationSourceURIs) > 0 {
				responseString.WriteString("\n***\n")
				responseString.WriteString(fmt.Sprintf("Text Citation %s:\n\n", pluralize(len(citationSourceURIs), "Source")))
				for _, citationSourceURI := range citationSourceURIs {
					responseString.WriteString(fmt.Sprintf("* [%s](%s)\n", citationSourceURI, citationSourceURI))
				}
			}

			// build list of code citation licenses
			citationSourceLicenses := []string{}
			if candidate.CitationMetadata != nil {
				for _, citationSource := range candidate.CitationMetadata.CitationSources {
					if citationSource.License != "" {
						citationSourceLicenses = append(citationSourceLicenses, citationSource.License)
					}
				}
			}

			// show code citation licenses (needs revision, output never seen)
			if len(citationSourceLicenses) > 0 {
				responseString.WriteString("\n***\n")
				responseString.WriteString(fmt.Sprintf("Code Citation %s:\n\n", pluralize(len(citationSourceLicenses), "License")))
				for _, citationSourceLicense := range citationSourceLicenses {
					responseString.WriteString(fmt.Sprintf("* %s\n", citationSourceLicense))
				}
			}

			// show why the model stopped generating tokens (content) (needs revision, output never seen)
			if candidate.FinishReason != genai.FinishReasonStop {
				responseString.WriteString("\n***\n")
				responseString.WriteString(fmt.Sprintf("Model stopped generating tokens (content) with reason [%s].\n", candidate.FinishReason.String()))
			}

			responseString.WriteString("\n***\n")
		}
	} else {
		// handle response error
		responseString.WriteString("**Error Response from Gemini:**\n\n")
		responseString.WriteString(err.Error())
		responseString.WriteString("\n***\n")
	}

	// print response metadata
	responseString.WriteString("```plaintext\n")
	responseString.WriteString(fmt.Sprintf("AI model   : %v (version %v)\n", strings.TrimPrefix(modelInfo.Name, "models/"), modelInfo.Version))
	responseString.WriteString(fmt.Sprintf("Generated  : %v\n", finishProcessing.Format(time.RFC850)))

	duration := finishProcessing.Sub(startProcessing)
	if err == nil {
		responseString.WriteString(fmt.Sprintf("Processing : %.1f secs for %d %s\n", duration.Seconds(),
			len(resp.Candidates), pluralize(len(resp.Candidates), "candidate")))
	} else {
		responseString.WriteString(fmt.Sprintf("Processing : %.1f secs resulting in error\n", duration.Seconds()))
	}

	if err == nil {
		if resp.UsageMetadata != nil {
			responseString.WriteString(fmt.Sprintf("Tokens     : %v (in: %v, out: %v)\n",
				resp.UsageMetadata.TotalTokenCount, resp.UsageMetadata.PromptTokenCount, resp.UsageMetadata.CandidatesTokenCount))
		}
		/* needs revision, output never seen */
		if resp.PromptFeedback != nil {
			responseString.WriteString(fmt.Sprintf("Blocked    : %v\n", resp.PromptFeedback.BlockReason.String()))
		}
	}

	responseString.WriteString("```\n")
	responseString.WriteString("\n***\n")

	// append response string to current markdown request/response file
	currentFileMarkdown, err := os.OpenFile(progConfig.MarkdownPromptResponseFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error [%v] at os.OpenFile()\n", err)
		return
	}
	defer currentFileMarkdown.Close()
	fmt.Fprint(currentFileMarkdown, responseString.String())

	// render markdown response as ansi
	ansiData := responseString.String()
	if progConfig.AnsiRendering {
		ansiData = renderMarkdown2Ansi(responseString.String())
	}

	// append response string to current ansi request/response file
	currentFileAnsi, err := os.OpenFile(progConfig.AnsiPromptResponseFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error [%v] at os.OpenFile()\n", err)
		return
	}
	defer currentFileAnsi.Close()
	fmt.Fprint(currentFileAnsi, ansiData)

	// render markdown response as html
	htmlData := responseString.String()
	if progConfig.HtmlRendering {
		htmlData = renderMarkdown2Html(responseString.String())
	}

	// append response string to current html request/response file
	currentFileHtml, err := os.OpenFile(progConfig.HtmlPromptResponseFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error [%v] at os.OpenFile()\n", err)
		return
	}
	defer currentFileHtml.Close()
	fmt.Fprint(currentFileHtml, htmlData)
}
