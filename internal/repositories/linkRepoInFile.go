package repositories

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/iHamsin/practicum-shortener-service/internal/util"
)

// LinkRepoRAM -.
type linksRepoInFile struct {
	file     os.File
	lastUUID int
}

// Link -.
type linkItem struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// New -.
func NewLinksRepoFile(file os.File) (*linksRepoInFile, error) {

	// получаем последнюю строку через tail @todo работает только для nix
	c := exec.Command("tail", "-1", file.Name())
	output, _ := c.Output()
	var lastLink linkItem
	_ = json.Unmarshal(output, &lastLink)

	return &linksRepoInFile{file, lastLink.UUID}, nil
}

// Insert -.
func (r *linksRepoInFile) Insert(originalURL string) (string, error) {
	shortURL := util.RandomString(8)

	r.lastUUID++
	var link = linkItem{r.lastUUID, shortURL, originalURL}

	jsonLink, jsonEncodeError := json.Marshal(link)
	if jsonEncodeError != nil {
		return "", jsonEncodeError
	}

	_, fileWriteError := r.file.WriteString(string(jsonLink) + "\n")
	if fileWriteError != nil {
		return "", fileWriteError
	}

	return shortURL, nil
}

// GetByCode -.
func (r *linksRepoInFile) GetByCode(shortURL string) (string, error) {

	scanner := bufio.NewScanner(&r.file)
	_, cursorResetError := r.file.Seek(0, io.SeekStart)
	if cursorResetError != nil {
		return "", cursorResetError
	}
	// optionally, resize scanner's capacity for lines over 64K, see next example
	var nextLink linkItem
	for scanner.Scan() {
		_ = json.Unmarshal(scanner.Bytes(), &nextLink)
		if nextLink.ShortURL == shortURL {
			return nextLink.OriginalURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return "", errors.New("link not found")
}
