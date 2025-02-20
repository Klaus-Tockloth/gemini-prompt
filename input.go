package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
readPromptFromKeyboard reads prompt from keyboard (Stdin).
*/
func readPromptFromKeyboard(promptChannel chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		promptData, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error [%v] at reader.ReadString()", err)
			return
		}
		if promptData == "\n" || promptData == "\r\n" {
			continue
		}

		// read prompt from given text file (e.g. "<<<MyQuery.txt" or "<<< MyQuery.txt")
		var fileData []byte
		if strings.HasPrefix(promptData, "<<<") {
			filename := strings.TrimSpace(strings.TrimPrefix(promptData, "<<<"))
			fileData, err = os.ReadFile(filename)
			if err != nil {
				fmt.Printf("error [%v] at os.ReadFile()\n", err)
				continue
			}
			if len(fileData) > 0 {
				promptChannel <- string(fileData)
			}
		} else {
			promptChannel <- promptData
		}
	}
}

/*
readPromptFromFile reads prompt (user input) from named file.
*/
func readPromptFromFile(filePath string, promptChannel chan string) {
	currentStat, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("error [%v] at os.Stat()", err)
	}
	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("error [%v] at os.Stat()", err)
		}
		if stat.Size() != currentStat.Size() || stat.ModTime() != currentStat.ModTime() {
			promptData, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("error [%v] at os.ReadFile()", err)
			}
			if len(promptData) > 0 {
				promptChannel <- string(promptData)
			}
			currentStat = stat
		}
		time.Sleep(500 * time.Millisecond)
	}
}

/*
readPromptFromLocalhost reads prompt (user input) from localhost.
*/
func readPromptFromLocalhost(promptChannel chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request body", http.StatusBadRequest)
			fmt.Printf("error [%v] reading request body\n", err)
			return
		} else {
			if len(body) == 0 {
				http.Error(w, "prompt empty", http.StatusBadRequest)
				return
			} else {
				promptChannel <- string(body)
			}
		}
		defer r.Body.Close()

		fmt.Fprintln(w, "prompt received")
	}
}
