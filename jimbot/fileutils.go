package jimbot

import (
	"bufio"
	"encoding/json" // new import
	"fmt"           // new import
	"io"
	"log"
	"net/http"
	"os"
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

// Updated UpdateConfig : Update a key in config.json without overwriting existing values
func UpdateConfig(pattern string, withStr string) error {
	// Open and decode config.json into a map
	file, err := os.Open("config.json")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	var configMap map[string]interface{}
	if err = json.NewDecoder(file).Decode(&configMap); err != nil {
		log.Println(err)
		return err
	}

	// Update the key if it exists
	if _, ok := configMap[pattern]; ok {
		configMap[pattern] = withStr
	} else {
		log.Printf("[-] Key %s not found", pattern)
		return fmt.Errorf("key not found: %s", pattern)
	}

	// Write the updated map back to config.json
	outFile, err := os.Create("config.json")
	if err != nil {
		log.Println(err)
		return err
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ") // format with indentation
	if err = encoder.Encode(configMap); err != nil {
		log.Println(err)
		return err
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
