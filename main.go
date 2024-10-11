package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var paths []string

func isExecutable(mode os.FileMode) bool {
	return mode.Perm()&(1<<(uint(7-2))) != 0
}

func ShowError() {
	errorMessage := "An error has occurred\n"
	_, err := os.Stderr.Write([]byte(errorMessage))
	if err != nil {
		return
	}
}

func HasRedirection(tokens ...string) (pos int, valid bool) {
	i := -1
	for idx, val := range tokens {
		if val == ">" && i == -1 {
			i = idx
		} else if val == ">" && i != -1 {
			return -1, false
		}
	}

	if i != -1 && i+1 != len(tokens)-1 {
		return -1, false
	}

	return i, true
}

func ExternalCommand(wg *sync.WaitGroup, tokens ...string) {
	defer wg.Done()
	//check if this command has redirection command and it's valid
	pos, valid := HasRedirection(tokens...)
	if !valid || (pos != -1 && pos == len(tokens)-1) {
		ShowError()
		return
	}

	if pos == -1 {
		pos = len(tokens)
	}

	exist, path := SearchCommand(tokens[0])
	if exist == true {
		args := tokens[1:pos]
		cmd := exec.Command(path, args...)
		if pos != len(tokens) {
			outFile, err := os.Create(tokens[pos+1])
			if err != nil {
				ShowError()
				return
			}
			defer outFile.Close()

			cmd.Stdout = outFile
			cmd.Stderr = outFile

		}

		if pos == len(tokens) {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err := cmd.Start()
		if err != nil {
			ShowError()
			return
		}
		err = cmd.Wait()
	} else {
		ShowError()
	}
}

func SearchCommand(command string) (bool, string) {
	// Iterate over paths and check if the file exist and it executable
	for _, val := range paths {
		path := val + "/" + command

		info, err := os.Stat(path)
		if err == nil {
			mode := info.Mode()
			if isExecutable(mode) {
				path = strings.Replace(path, "/bin/", "", 1)
				return true, path
			}
		}

	}
	return false, ""
}

func CdCommand(tokens ...string) {
	if len(tokens) == 0 || len(tokens) > 1 {
		ShowError()
		return
	}

	err := os.Chdir(tokens[0])
	if err != nil {
		ShowError()
		return
	}

}

func ExitCommand(tokens ...string) {
	if len(tokens) > 0 {
		ShowError()
		return
	}
	os.Exit(0)
}

func PathCommand(tokens ...string) {
	if len(tokens) == 0 {
		clear(paths)
	} else {
		paths = tokens
	}
}

func ExecuteCommand(cmd [][]string) {

	//	Check if the command is built-in command
	//	if yes -> execute it
	//	if no -> search on paths and check if the command exist or not
	var wg sync.WaitGroup
	for _, tokens := range cmd {
		if len(tokens) == 0 {
			continue
		}

		switch tokens[0] {
		case "cd":
			CdCommand(tokens[1:]...)
		case "exit":
			ExitCommand(tokens[1:]...)
		case "path":
			PathCommand(tokens[1:]...)
		default:
			wg.Add(1)
			go ExternalCommand(&wg, tokens...)
		}
	}
	wg.Wait()
}

func IsValid(command string) bool {
	input := strings.ReplaceAll(command, " ", "")
	if input[0] == '>' || input[len(input)-1] == '>' {
		return false
	}

	for i := 0; i < len(input)-1; i++ {
		if (input[i] == '&' || input[i] == '>') && (input[i+1] == '>' || input[i+1] == '&') {
			return false
		}
	}

	return true
}

func ParallelCommand(input ...string) [][]string {
	var tokens = make([][]string, 1)
	pos := 0
	for _, val := range input {

		if val == "&" {
			tokens = append(tokens, []string{})
			pos++
		} else {
			str := ""
			for _, c := range val {
				if c != '&' {
					str += string(c)
				} else {
					tokens[pos] = append(tokens[pos], str)
					tokens = append(tokens, []string{})
					str = ""
					pos++
				}
			}

			if str != "" {
				tokens[pos] = append(tokens[pos], str)
			}
		}
	}
	return tokens
}

func ParseToken(input ...string) [][]string {

	tokens := ParallelCommand(input...)

	for idx, token := range tokens {
		var newToken []string
		for _, val := range token {
			if val == ">" {
				newToken = append(newToken, val)
			} else {
				str := ""
				for _, c := range val {
					if c != '>' {
						str += string(c)
					} else {
						newToken = append(newToken, str)
						newToken = append(newToken, string(c))
						str = ""
					}
				}
				if str != "" {
					newToken = append(newToken, str)
				}
			}
			tokens[idx] = newToken
		}
	}

	return tokens
}

func BatchMode(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		ShowError()
		os.Exit(1)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		command := scanner.Text()

		command = strings.TrimSpace(command)
		if len(command) == 0 {
			continue
		}
		if !IsValid(command) {
			ShowError()
			continue
		}
		tokens := ParseToken(strings.Split(command, " ")...)
		ExecuteCommand(tokens)

	}

	if err := scanner.Err(); err != nil {
		ShowError()
	}

}

func main() {
	reader := bufio.NewReader(os.Stdin)
	paths = append(paths, "/bin")

	if len(os.Args) > 2 {
		ShowError()
		os.Exit(1)
		return
	}

	if len(os.Args) > 1 {
		fileName := os.Args[1]
		BatchMode(fileName)
		return
	}

	for {
		fmt.Print("wish> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			ShowError()
			continue
		}

		input = strings.TrimSpace(input)
		if len(input) == 0 {
			continue
		}
		if !IsValid(input) {
			ShowError()
			continue
		}
		tokens := ParseToken(strings.Split(input, " ")...)
		ExecuteCommand(tokens)

	}
}
