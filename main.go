package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	color.HiMagenta(`
      ___.                   __                
___  _\_ |__  __ __  _______/  |_  ___________ 
\  \/ /| __ \|  |  \/  ___/\   __\/ __ \_  __ \
 \   / | \_\ \  |  /\___ \  |  | \  ___/|  | \/
  \_/  |___  /____//____  > |__|  \___  >__|   
           \/           \/            \/       
`)
	color.Yellow("\t\t\t\t\t\t\t\tby Rai from Vsec 👾\n")
	color.HiMagenta("VBuster - Simple Web Vulnerability Scanner")
	color.HiCyan("Note: Please use this program in accordance with laws and ethical guidelines.")
	color.HiMagenta("-------------------------------------------------------------")
	var _input string
	reader := bufio.NewReader(os.Stdin)
	color.HiMagenta("Type 'q' to exit.")
	fmt.Println(" ")
	for {
		fmt.Print("[36mVBuster > [0m")
		line, err := reader.ReadString('\n')
		if err != nil {
			color.Red("Input error: %v", err)
			continue
		}
		line = strings.TrimSpace(line)
		_input = line
		if _input == "q" {
			color.HiMagenta("Exiting VBuster. Goodbye!")
			return
		} else if _input == "bust" {
			color.Red("Error: 'bust' command requires additional parameters. Use 'bust <target_url> <wordlist>'.")
			continue
		} else if _input == "" {
			color.Red("Error: No input provided. Please enter a valid command.")
			continue
		} else if _input == "help" {
			color.HiMagenta("Available commands:\n- q: Exit the program\n- help: Show this help message\n- bust <target_url> <wordlist>: Start scanning with the specified target URL and wordlist")
			continue
		} else if len(_input) > 5 && _input[:5] == "bust " {
			// Split input by spaces
			args := make([]string, 0)
			for _, v := range splitBySpace(_input) {
				if v != "" {
					args = append(args, v)
				}
			}
			if len(args) != 3 {
				color.Red("Usage: bust <target_url> <wordlist>")
				continue
			}
			target := args[1]
			wordlist := args[2]
			file, err := os.Open(wordlist)
			if err != nil {
				color.Red("Wordlist could not be opened: %v", err)
				continue
			}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				path := scanner.Text()
				if path == "" {
					continue
				}
				url := target + "/" + path
				resp, err := http.Get(url)
				if err != nil {
					color.Red("Request failed for %s: %v", url, err)
					continue
				}
				if resp.StatusCode == 200 {
					color.HiGreen("Found: %s", url)
					color.HiBlue("Do you want to save this URL to a txt file? (y/n)")
					var saveInput string
					fmt.Scanln(&saveInput)
					if strings.ToLower(saveInput) == "y" {
						color.HiBlue("Enter the filename to save the URL:")
						var filename string
						fmt.Scanln(&filename)
						fileToSave, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							color.Red("Error opening file for saving: %v", err)
						} else {
							_, err = fileToSave.WriteString(url + "\n")
							if err != nil {
								color.Red("Error writing to file: %v", err)
							} else {
								color.HiGreen("URL saved to %s", filename)
							}
							fileToSave.Close()
						}
					} else {
						color.HiBlue("URL not saved.")
					}
				} else {
					color.Red("Not found (%d): %s", resp.StatusCode, url)
				}
				resp.Body.Close()
			}
			file.Close()
			if err := scanner.Err(); err != nil {
				color.Red("Wordlist read error: %v", err)
			}
			continue
		} else {
			color.Red("Unknown command: %s. Type 'help' for available commands.", _input)
			continue
		}
	}

}

func splitBySpace(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
