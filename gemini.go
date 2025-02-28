package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
)

/*
uploadFilesToGemini uploads all files given from command line.
*/
func uploadFilesToGemini(ctx context.Context, client *genai.Client, clFiles []string) ([]*genai.File, error) {
	var err error

	files := []*genai.File{}
	if len(clFiles) == 0 {
		return files, nil
	}

	fmt.Printf("\nFile uploads:\n")
	for _, filename := range clFiles {
		fmt.Printf("  %s ... ", filename)
		if !fileExists(filename) {
			fmt.Printf("error: file does't exist\n")
			continue
		}

		// display name = max 512 characters
		uploadOptions := genai.UploadFileOptions{}
		uploadOptions.DisplayName = filename

		file, err := client.UploadFileFromPath(ctx, filename, &uploadOptions)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Printf("uploaded (%.1f KiB, %s)\n", float64(file.SizeBytes)/1024.0, file.MIMEType)
			files = append(files, file)
		}
	}

	// ensure that all remote (uploaded) files have state "active" (e.g. videos need to be processed)
	fmt.Printf("\nRemote file states:\n")
	maxWaitDuration := progConfig.GeminiMaxWaitTimeFileProcessing // seconds
	for i, file := range files {
		currentWaitDuration := 0
		tmpFile := file
		for tmpFile.State == genai.FileStateProcessing {
			tmpFile, err = client.GetFile(ctx, file.Name)
			if err != nil {
				fmt.Printf("  error [%s] getting state for remote file [%s]\n", err, file.DisplayName)
				break
			}
			if tmpFile.State == genai.FileStateActive {
				break
			}
			if tmpFile.State != genai.FileStateProcessing {
				break
			}
			time.Sleep(2 * time.Second)
			currentWaitDuration += 2
			if currentWaitDuration >= maxWaitDuration {
				break
			}
			fmt.Printf(".")
		}
		files[i] = tmpFile
		fmt.Printf("\r  %s ... %s\n", tmpFile.DisplayName, tmpFile.State.String())
	}

	return files, nil
}
