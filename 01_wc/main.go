package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type WcCommand struct {
	WithByteCount      bool
	WithWordCount      bool
	WithLineCount      bool
	WithCharacterCount bool
	FilePath           string
	Inputs             []*WcInput
}

func (c *WcCommand) hasCommands() bool {
	return c.WithByteCount || c.WithCharacterCount || c.WithLineCount || c.WithWordCount
}

type WcInput struct {
	Type     string
	FileName string
	Reader   io.Reader
}

func main() {
	runner()
}

func runner() {
	args := os.Args[1:]
	command, err := parseArgs(args)
	if err != nil {
		fmt.Printf("Error occurred while parsing :: %s", err.Error())
		return
	}
	exec(command)
	if err != nil {
		fmt.Printf("Error occurred while executing :: %s", err.Error())
		return
	}
}

func exec(command *WcCommand) error {
	for _, input := range command.Inputs {
		fmt.Printf("Type :: %s\n", input.Type)
		content, err := io.ReadAll(input.Reader)
		if err != nil {
			return err
		}
		if command.WithWordCount {
			fmt.Printf("Word Count :: %d %s \n", len(strings.Fields(string(content))), input.FileName)
		}
		if command.WithLineCount {
			fmt.Printf("Line Count :: %d %s \n", strings.Count(string(content), "\n"), input.FileName)
		}
		if command.WithByteCount {
			fmt.Printf("Byte Count :: %d %s \n", len(content), input.FileName)
		}
		if command.WithCharacterCount {
			str := string(content)
			runeCount := 0
			for range str {
				runeCount++
			}
			fmt.Printf("Character Count :: %d %s \n", runeCount, input.FileName)
		}
	}

	return nil
}

func parseArgs(args []string) (*WcCommand, error) {
	command := &WcCommand{}
	i := 0

	for i < len(args) && strings.HasPrefix(args[i], "-") {
		commandString := strings.Split(args[i], "")
		for idx, c := range commandString {
			if idx == 0 {
				continue
			}
			switch c {
			case "c":
				command.WithByteCount = true
			case "w":
				command.WithWordCount = true
			case "l":
				command.WithLineCount = true
			case "m":
				command.WithCharacterCount = true
			default:
				return nil, fmt.Errorf("command not found %s", c)
			}
		}
		i++
	}

	if !command.hasCommands() {
		// Use defaults
		command.WithByteCount = true
		command.WithLineCount = true
		command.WithWordCount = true
	}

	if i == len(args) {
		command.Inputs = []*WcInput{
			{
				Type:   "StdIn",
				Reader: os.Stdin,
			},
		}

		return command, nil
	}

	command.Inputs = make([]*WcInput, 0)

	for i < len(args) {
		fileName := args[i]
		data, err := os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		command.Inputs = append(command.Inputs, &WcInput{
			Type:     "File",
			Reader:   bytes.NewReader(data),
			FileName: fileName,
		})
		i++
	}

	return command, nil
}
