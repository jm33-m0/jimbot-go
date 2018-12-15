package jimbot

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// WriteStringToFile : write or append line to file
func WriteStringToFile(path string, text string, overwrite bool) error {
	var err error
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if overwrite {
		f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = f.WriteString(text + "\n")
	if err != nil {
		return err
	}
	return nil
}

// FileToLines : Read lines from a text file
func FileToLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return linesFromReader(f)
}

// UpdateConfig : Grep a line and replace it with a given string
func UpdateConfig(pattern string, withStr string) error {
	lines, err := FileToLines("config.txt")
	if err != nil {
		log.Println(err)
		return err
	}
	err = os.Remove("config.txt")
	if err != nil {
		log.Println(err)
	}

	for _, line := range lines {
		if strings.HasPrefix(line, pattern) {
			greeting := strings.Split(line, ": ")[1]
			log.Printf("Replacing %s with %s", greeting, withStr)
			line = withStr
		}

		err = WriteStringToFile("config.txt", line, false)
		if err != nil {
			log.Printf("Error updating config: %s", err.Error())
		}
	}

	return nil
}

// DownloadFile : Download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		err = out.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
