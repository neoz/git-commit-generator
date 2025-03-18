package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Command line flags
type Flags struct {
	OnlyMessage bool
	Verbose     bool
	Help        bool
	Update      bool
}

func main() {
	// Parse command line arguments
	flags := parseFlags()

	// Display help if requested
	if flags.Help {
		displayHelp()
		return
	}

	// Update Ollama model if requested
	if flags.Update {
		updateOllamaModel(flags)
	}

	// Welcome Header
	if !flags.OnlyMessage {
		displayWelcomeHeader()
	}

	// Get staged diff
	diff, err := getStagedDiff()
	if err != nil {
		fmt.Printf("\033[1;31m⚠️ Error getting staged changes: %v\033[0m\n", err)
		return
	}

	if diff == "" {
		fmt.Println("\033[1;31m⚠️ No changes detected. Please stage your changes first.\033[0m")
		return
	}

	// Display diff in a box if not in only-message mode
	if !flags.OnlyMessage {
		displayDiff(diff)
	}

	// Ask for additional context (optional)
	globalContext := ""
	if !flags.OnlyMessage {
		globalContext = requestContext()
	}

	// Split diff into chunks
	chunks := splitDiff(diff)

	// Variable to aggregate micro commit messages
	finalMicroMessages := ""

	// Display header for micro commits reasoning if not in only-message mode
	if !flags.OnlyMessage {
		fmt.Println("\n\033[1;34m──────────────────────────────────────────────────────────────\033[0m")
		fmt.Println("\033[1;34mReasoning\033[0m")
		fmt.Println("\033[1;34m──────────────────────────────────────────────────────────────\033[0m")
	}

	totalChunks := len(chunks)

	// Process each chunk
	for i, chunk := range chunks {
		chunkIndex := i + 1

		if !flags.OnlyMessage {
			fmt.Println(" ")
		}

		// Verbose output: Show chunk diff
		if flags.Verbose {
			fmt.Printf("\033[1;36mChunk %d/%d:\033[0m\n", chunkIndex, totalChunks)
			printColorizedDiff(chunk)
			fmt.Println("\033[1;34m----------------\033[0m")
		}

		// Prepare input for Ollama with global context if provided
		inputForChunk := chunk
		if globalContext != "" {
			inputForChunk = chunk + "\n\nUser Extra Context Input: " + globalContext
		}

		// Get micro message for the chunk
		_, chunkMessage := processChunk(inputForChunk, flags.OnlyMessage)

		// Verbose output: Show generated micro message
		if flags.Verbose {
			fmt.Printf("\033[1;36mGenerated Micro Message for Chunk %d:\033[0m\n", chunkIndex)
			fmt.Println(chunkMessage)
			fmt.Println("\033[1;34m----------------\033[0m")
		}

		// Aggregate micro message for final commit
		finalMicroMessages += "\n" + chunkMessage + "\n----\n"
	}

	ignoreOption := ""
	if flags.OnlyMessage {
		ignoreOption = "c"
	}

	generateAndProcessFinalMessage(finalMicroMessages, globalContext, flags, ignoreOption)
}

// Parse command line flags
func parseFlags() Flags {
	flags := Flags{}

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--only-message":
			flags.OnlyMessage = true
		case "--verbose":
			flags.Verbose = true
		case "-h", "--help":
			flags.Help = true
		case "--update":
			flags.Update = true
		}
	}

	return flags
}

// Display help message
func displayHelp() {
	fmt.Println("\033[1;34mGit Commit Message Generator\033[0m")
	fmt.Println("\033[1;34m======================\033[0m")
	fmt.Println("This script generates intelligent git commit messages based on staged changes.")
	fmt.Println("\n\033[1mUsage:\033[0m")
	fmt.Println("  ./git-commit-generator [options]")
	fmt.Println("\n\033[1mOptions:\033[0m")
	fmt.Println("  --only-message    Output only the final commit message without UI")
	fmt.Println("  --verbose         Print detailed steps including chunks and diffs")
	fmt.Println("  -h, --help        Display this help message")
	fmt.Println("  --update          Update the Ollama model before running")
	fmt.Println("\n\033[1mDescription:\033[0m")
	fmt.Println("  - Analyzes staged git changes (git diff --staged)")
	fmt.Println("  - Splits changes into chunks for better analysis")
	fmt.Println("  - Generates micro commit messages for each chunk")
	fmt.Println("  - Combines them into a final cohesive commit message")
	fmt.Println("  - Allows review and editing before committing")
}

// Update Ollama model
func updateOllamaModel(flags Flags) {
	fmt.Println("\033[1;33mUpdating Ollama model...\033[0m")

	cmd := exec.Command("ollama", "pull", "tavernari/git-commit-message:merge_commits")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("ollama", "pull", "tavernari/git-commit-message:reasoning")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Println("\033[1;32mModel updated successfully!\033[0m")
}

// Display welcome header
func displayWelcomeHeader() {
	fmt.Println("\033[1;34m╔═════════════════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;34m║                 Git Commit Message Generator                    ║\033[0m")
	fmt.Println("\033[1;34m╚═════════════════════════════════════════════════════════════════╝\033[0m")
}

// Get staged diff from git
func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Display diff in a colored box
func displayDiff(diff string) {
	fmt.Println("\n\033[1;34m┌─────────────────────────────────────────────────────────────────┐\033[0m")
	fmt.Println("\033[1;34m│ Diff                                                            │\033[0m")
	fmt.Println("\033[1;34m└─────────────────────────────────────────────────────────────────┘\033[0m")
	printColorizedDiff(diff)
	fmt.Println("\n\033[1;34m──────────────────────────────────────────────────────────────\033[0m")
}

// Request additional context from user
func requestContext() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\033[1;33mProvide additional context for the commit (optional, press Enter to skip):\033[0m")
	context, _ := reader.ReadString('\n')
	return strings.TrimSpace(context)
}

// Split diff into chunks
func splitDiff(diff string) []string {
	chunks := []string{}
	scanner := bufio.NewScanner(strings.NewReader(diff))

	var currentChunk strings.Builder
	diffGitPattern := regexp.MustCompile(`^diff --git`)

	for scanner.Scan() {
		line := scanner.Text()

		if diffGitPattern.MatchString(line) && currentChunk.Len() > 0 {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}

		currentChunk.WriteString(line)
		currentChunk.WriteString("\n")
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// Process a chunk to get a micro commit message
func processChunk(chunk string, onlyMessage bool) (string, string) {
	cmd := exec.Command("ollama", "run", "tavernari/git-commit-message:reasoning", chunk)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("\033[1;31mError creating pipe: %v\033[0m\n", err)
		return "", ""
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("\033[1;31mError starting command: %v\033[0m\n", err)
		return "", ""
	}

	section := "normal"
	var chunkReasoning strings.Builder
	var chunkMessage strings.Builder

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "<reasoning>") {
			section = "reasoning"
			continue
		}
		if strings.Contains(line, "</reasoning>") {
			section = "commit"
			continue
		}

		switch section {
		case "reasoning":
			if !onlyMessage {
				typeEffect(line, "", "  ")
			}
			chunkReasoning.WriteString(line)
			chunkReasoning.WriteString("\n")
		case "commit":
			chunkMessage.WriteString(line)
			chunkMessage.WriteString("\n")
		}
	}

	cmd.Wait()

	return chunkReasoning.String(), strings.TrimSpace(chunkMessage.String())
}

// Generate and process the final commit message
func generateAndProcessFinalMessage(finalMicroMessages, globalContext string, flags Flags, ignoreOption string) {
	finalInput := fmt.Sprintf(`
### Commits:
%s

### Extra user input context:
%s
`, finalMicroMessages, globalContext)

	if flags.Verbose {
		fmt.Println("\033[1;36mFinal Input to Ollama:\033[0m")
		fmt.Println(finalInput)
	}

	fmt.Println("")

	// Try generating final commit message with retries
	maxRetries := 10
	retryCount := 0
	finalCommitMessage := ""

	for finalCommitMessage == "" && retryCount < maxRetries {
		_, commitMessage := getFinalCommitMessage(finalInput, flags.OnlyMessage)
		finalCommitMessage = commitMessage

		if finalCommitMessage == "" {
			retryCount++
			if retryCount < maxRetries {
				if flags.Verbose {
					fmt.Println("\033[1;31m❌ Failed to generate a commit message. Retrying...\033[0m")
				}
				time.Sleep(1 * time.Second)
			} else {
				fmt.Println("\033[1;31m❌ Failed to generate a commit message.\033[0m")
				ignoreOption = "c, e"
				choice := displayOptions(ignoreOption)
				processChoice(choice, finalCommitMessage, flags, ignoreOption, finalMicroMessages, globalContext)
				return
			}
		}
	}

	// Verbose output
	if flags.Verbose {
		fmt.Println("\033[1;36mFinal Input to Ollama:\033[0m")
		fmt.Println(finalInput)
		fmt.Println("\033[1;34m----------------\033[0m")
		fmt.Println("\033[1;36mGenerated Final Commit Message:\033[0m")
		fmt.Println(finalCommitMessage)
		fmt.Println("\033[1;34m----------------\033[0m")
	}

	if flags.OnlyMessage {
		fmt.Println(finalCommitMessage)
		return
	}

	// Display proposed commit message with typing effect
	fmt.Println("\n\033[1;34m──────────────────────────────────────────────────────────────\033[0m")
	fmt.Println("\033[1;34mFinal Commit Message\033[0m")
	fmt.Println("\033[1;34m──────────────────────────────────────────────────────────────\033[0m")

	scanner := bufio.NewScanner(strings.NewReader(finalCommitMessage))
	for scanner.Scan() {
		typeEffect(scanner.Text(), "", "")
	}

	// Present options to the user
	choice := displayOptions(ignoreOption)
	processChoice(choice, finalCommitMessage, flags, ignoreOption, finalMicroMessages, globalContext)
}

// Get final commit message from Ollama
func getFinalCommitMessage(finalInput string, onlyMessage bool) (string, string) {
	cmd := exec.Command("ollama", "run", "tavernari/git-commit-message:merge_commits", finalInput)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("\033[1;31mError creating pipe: %v\033[0m\n", err)
		return "", ""
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("\033[1;31mError starting command: %v\033[0m\n", err)
		return "", ""
	}

	section := "normal"
	var reasoningOutput strings.Builder
	var commitMessage strings.Builder

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "<reasoning>") {
			section = "reasoning"
			continue
		}
		if strings.Contains(line, "</reasoning>") {
			section = "commit"
			continue
		}

		switch section {
		case "reasoning":
			if !onlyMessage {
				typeEffect(line, "", "  ")
			}
			reasoningOutput.WriteString(line)
			reasoningOutput.WriteString("\n")
		case "commit":
			commitMessage.WriteString(line)
			commitMessage.WriteString("\n")
		}
	}

	cmd.Wait()

	return reasoningOutput.String(), strings.TrimSpace(commitMessage.String())
}

// Display options to the user and get choice
func displayOptions(ignoreOption string) string {
	fmt.Println("\n\033[1;34mOptions:\033[0m")

	if !strings.Contains(ignoreOption, "c") {
		fmt.Println("  \033[1m(c)\033[0m Commit with this message")
	}
	if !strings.Contains(ignoreOption, "e") {
		fmt.Println("  \033[1m(e)\033[0m Edit this message")
	}
	if !strings.Contains(ignoreOption, "g") {
		fmt.Println("  \033[1m(g)\033[0m Generate again with some context")
	}
	if !strings.Contains(ignoreOption, "d") {
		fmt.Println("  \033[1m(d)\033[0m Discard")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\033[1;33mChoose: \033[0m")
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

// Process user choice
func processChoice(choice, finalCommitMessage string, flags Flags, ignoreOption, finalMicroMessages, globalContext string) {
	switch choice {
	case "c":
		if !strings.Contains(ignoreOption, "c") {
			fmt.Println("\033[1;33mCommitting with the following message:\033[0m")
			fmt.Println(finalCommitMessage)

			cmd := exec.Command("git", "commit", "-m", finalCommitMessage)
			err := cmd.Run()

			if err == nil {
				fmt.Println("\033[1;32m✅ Commit created successfully!\033[0m")
			} else {
				fmt.Println("\033[1;31m❌ Commit failed. Please check the errors above.\033[0m")
				newChoice := displayOptions(ignoreOption)
				processChoice(newChoice, finalCommitMessage, flags, ignoreOption, finalMicroMessages, globalContext)
			}
		} else {
			fmt.Println("\033[1;31mInvalid option. Aborting.\033[0m")
			os.Exit(1)
		}

	case "e":
		if !strings.Contains(ignoreOption, "e") {
			fmt.Println("\033[1;33mOpening editor to edit the commit message...\033[0m")

			// Create temp file
			tempFile, err := ioutil.TempFile("", "commit-msg-")
			if err != nil {
				fmt.Printf("\033[1;31mError creating temp file: %v\033[0m\n", err)
				return
			}
			defer os.Remove(tempFile.Name())

			// Write message to temp file
			_, err = tempFile.WriteString(finalCommitMessage)
			tempFile.Close()
			if err != nil {
				fmt.Printf("\033[1;31mError writing to temp file: %v\033[0m\n", err)
				return
			}

			// Open editor
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nano"
			}

			cmd := exec.Command(editor, tempFile.Name())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				fmt.Printf("\033[1;31mError with editor: %v\033[0m\n", err)
				return
			}

			// Read edited message
			editedMessage, err := ioutil.ReadFile(tempFile.Name())
			if err != nil {
				fmt.Printf("\033[1;31mError reading edited message: %v\033[0m\n", err)
				return
			}

			updatedMessage := string(editedMessage)

			fmt.Println("\n\033[1;34m──────────────────────────────────────────────────────────────\033[0m")
			fmt.Println("\033[1;34mUpdated Commit Message\033[0m")
			fmt.Println("\033[1;34m──────────────────────────────────────────────────────────────\033[0m")
			fmt.Println(updatedMessage)

			newChoice := displayOptions(ignoreOption)
			processChoice(newChoice, updatedMessage, flags, ignoreOption, finalMicroMessages, globalContext)
		} else {
			fmt.Println("\033[1;31mInvalid option. Aborting.\033[0m")
			os.Exit(1)
		}

	case "g":
		if !strings.Contains(ignoreOption, "g") {
			fmt.Println("\033[1;33mGenerating again with some context...\033[0m")
			newContext := requestContext()
			generateAndProcessFinalMessage(finalMicroMessages, newContext, flags, ignoreOption)
		} else {
			fmt.Println("\033[1;31mInvalid option. Aborting.\033[0m")
			os.Exit(1)
		}

	case "d":
		fmt.Println("\033[1;31m❌ Commit discarded.\033[0m")

	default:
		fmt.Println("\033[1;31mInvalid option. Aborting.\033[0m")
	}
}

// Print colorized diff
func printColorizedDiff(diff string) {
	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "+") {
			fmt.Printf("\033[32m%s\033[0m\n", line) // Green for additions
		} else if strings.HasPrefix(line, "-") {
			fmt.Printf("\033[31m%s\033[0m\n", line) // Red for deletions
		} else {
			fmt.Printf("\033[90m%s\033[0m\n", line) // Gray for context
		}
	}
}

// Simulate typing effect
func typeEffect(text, color, indent string) {
	// Remove ANSI color codes for calculation
	cleanText := removeANSIColorCodes(text)

	fmt.Print(indent, color)
	for _, char := range cleanText {
		fmt.Print(string(char))
		time.Sleep(2 * time.Millisecond)
	}
	fmt.Println("\033[0m")
}

// Remove ANSI color codes from text
func removeANSIColorCodes(text string) string {
	re := regexp.MustCompile(`\033\[[0-9;]*m`)
	return re.ReplaceAllString(text, "")
}
