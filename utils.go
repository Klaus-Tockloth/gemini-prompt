package main

import (
	"bufio"
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/flytam/filenamify"
	"github.com/gofrs/uuid"
	"github.com/mitchellh/go-wordwrap"
)

/*
fileExists checks if file already exists.
*/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/*
wrapString wraps (long) string for better readability.
*/
func wrapString(message string, width int, ident int) string {
	wrapped := wordwrap.WrapString(message, uint(width-ident))
	wrapped = strings.ReplaceAll(wrapped, "\n", "\n"+strings.Repeat(" ", ident))
	return wrapped
}

/*
copyFile copies source file to destination file.
*/
func copyFile(sourceFile, destinationFile string) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Printf("error [%v] at os.ReadFile()\n", err)
		return
	}

	err = os.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Printf("error [%v] at os.WriteFile()\n", err)
		return
	}
}

/*
pluralize pluralizes term (simply) by adding 's'.
*/
func pluralize(count int, singular string) string {
	if count == 1 {
		return singular
	}
	return singular + "s"
}

/*
runCommand runs a command or program.
*/
func runCommand(commandLine string) error {
	parsedArgs := splitCommandLine(commandLine)
	cmd := exec.Command(parsedArgs[0], parsedArgs[1:]...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error [%v] executing command [%v]\n", err, commandLine)
	}
	return err
}

/*
splitCommandLine splits command line into components (program, options, parameters).
*/
func splitCommandLine(commandLine string) []string {
	var args []string
	var inQuote bool
	var quoteType rune // ' or "
	var currentArg strings.Builder

	for _, r := range commandLine {
		switch {
		case r == '"' || r == '\'':
			if inQuote {
				if quoteType == r {
					inQuote = false
					args = append(args, currentArg.String())
					currentArg.Reset()
				} else {
					// Inside a quotation mark a different type is found, so treat it as part of the argument.
					currentArg.WriteRune(r)
				}
			} else {
				inQuote = true
				quoteType = r
			}
		case r == ' ' && !inQuote:
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		default:
			currentArg.WriteRune(r)
		}
	}

	// add remaining argument, if any
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

/*
getPassword gets password from pass phrase argument.
*/
func getPassword(passPhrase string) (string, error) {
	items := strings.SplitN(passPhrase, ":", 2)
	source := strings.ToLower(items[0])
	password := ""

	if len(items) != 2 {
		return "", fmt.Errorf("unable to split pass phrase argument into 'source:password'")
	}

	switch source {
	case "pass":
		password = items[1]
	case "env":
		password = os.Getenv(items[1])
		if password == "" {
			return "", fmt.Errorf("password empty or env variable [%s] not found", items[1])
		}
	case "file":
		// read password file
		lines, err := slurpFile(items[1])
		if err != nil || len(lines) == 0 {
			return "", fmt.Errorf("unable to read password from file, error = [%w], file = [%v]", err, items[1])
		}
		password = lines[0]
	default:
		return "", fmt.Errorf("invalid password source (not 'pass:', 'env:' or 'file:')")
	}

	return password, nil
}

/*
slurpFile slurps all lines of a text file into a slice of strings.
*/
func slurpFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, nil
}

/*
getMimeType gets mime type for given file.
*/
func getMimeType(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		// empty file = default mime type
		if err.Error() == "EOF" {
			return "application/octet-stream", nil
		}
		return "", err
	}

	mimeType := http.DetectContentType(buffer)

	// if content type detection fails, resolve mime type based on file extension
	if mimeType == "application/octet-stream" {
		ext := filepath.Ext(filename)
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			// unknown mime type = default mime type
			mimeType = "application/octet-stream"
		}
	}

	return mimeType, nil
}

/*
promptToFilename derives valid filename from given prompt, prefix, postfix and extension.
*/
func promptToFilename(prompt string, maxLength int, prefix, postfix, extension string) string {
	var filename string
	var err error

	// length correction for 'core' name
	maxLength -= 2 // "[" + "]"
	if prefix != "" {
		maxLength -= len(prefix) + 1
	}
	if postfix != "" {
		maxLength -= 1 + len(postfix)
	}
	if extension != "" {
		maxLength -= 1 + len(extension)
	}

	// replace problematic characters with visuell similar runes
	prompt = strings.ReplaceAll(prompt, "?", "ʔ")  // glottal stop
	prompt = strings.ReplaceAll(prompt, ":", "ː")  // triangular colon
	prompt = strings.ReplaceAll(prompt, "/", "∕")  // division slash
	prompt = strings.ReplaceAll(prompt, "\\", "＼") // fullwidth reverse solidus
	prompt = strings.ReplaceAll(prompt, "*", "⁎")  // low asterisk
	prompt = strings.ReplaceAll(prompt, "|", "¦")  // broken bar
	prompt = strings.ReplaceAll(prompt, "<", "‹")  // single left-pointing angle quotation mark
	prompt = strings.ReplaceAll(prompt, ">", "›")  // single right-pointing angle quotation mark
	prompt = strings.ReplaceAll(prompt, "\"", "”") // right double quotation mark
	prompt = strings.ReplaceAll(prompt, ".", "․")  // one dot leader

	filename, err = filenamify.Filenamify(prompt, filenamify.Options{Replacement: " ", MaxLength: maxLength})
	if err != nil {
		fmt.Printf("error [%v] at filenamify.Filenamify()\n", err)
		uuid4, _ := uuid.NewV4()
		filename = uuid4.String()
	}
	filename = "[" + filename + "]"

	if prefix != "" {
		filename = prefix + "." + filename
	}
	if postfix != "" {
		filename += "." + postfix
	}
	if extension != "" {
		filename += "." + extension
	}

	return filename
}

/*
buildDestinationFilename builds destination filename based on given parameters and program configuration.
*/
func buildDestinationFilename(now time.Time, prompt, extension string) string {
	formatLayout := "20060102-150405"
	timestamp := now.Format(formatLayout)

	destinationFilename := ""
	switch progConfig.HistoryFilenameSchema {
	case "prompt":
		prefix := ""
		if progConfig.HistoryFilenameAddPrefix {
			prefix = timestamp
		}
		postfix := ""
		if progConfig.HistoryFilenameAddPostfix {
			postfix = timestamp
		}
		destinationFilename = promptToFilename(prompt, progConfig.HistoryMaxFilenameLength, prefix, postfix, extension)
	case "timestamp":
		destinationFilename = timestamp
		if extension != "" {
			destinationFilename += "." + extension
		}
	}
	return destinationFilename
}
