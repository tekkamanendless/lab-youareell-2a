package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/tekkamanendless/lab-youareell-2a/internal/commandlineprocessor"
	"github.com/tekkamanendless/lab-youareell-2a/internal/youareellclient"
)

func main() {
	interactive := flag.Bool("interactive", false, "Enable the interactive console.")
	verbose := flag.Bool("verbose", false, "Enable verbose logging.")
	flag.Parse()

	args := flag.Args()

	client := &youareellclient.Client{
		Debug: *verbose,
	}
	if !*interactive {
		err := commandlineprocessor.Process(client, args)
		if err != nil {
			fmt.Printf("Error: [%T] %v\n", err, err)
			os.Exit(1)
		}
		return
	}

	// This will do a decent job of splitting a string into tokens, respecting quoted parts.
	// It's not magic, but it's workable.
	// See: https://stackoverflow.com/questions/47489745/splitting-a-string-at-space-except-inside-quotation-marks
	tokenizer := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			break
		}
		tokens := tokenizer.FindAllString(line, -1)
		if *verbose {
			fmt.Printf("[Tokens: %v]\n", tokens)
		}
		err = commandlineprocessor.Process(client, tokens)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
	}
}
