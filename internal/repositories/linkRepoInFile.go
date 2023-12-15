package repositories

import (
	"bufio"
	"context"
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
func (r *linksRepoInFile) InsertLink(ctx context.Context, originalURL string) (string, error) {
	shortURL := util.RandomString(cfg.ShortCodeLength)

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

// BatchInsert -.
func (r *linksRepoInFile) BatchInsertLink(ctx context.Context, links []string) ([]string, error) {
	result := make([]string, len(links))

	for i, link := range links {
		result[i] = util.RandomString(cfg.ShortCodeLength)
		// insert
		r.lastUUID++
		var newLine = linkItem{r.lastUUID, result[i], link}
		jsonLink, jsonEncodeError := json.Marshal(newLine)
		if jsonEncodeError != nil {
			return nil, jsonEncodeError
		}
		_, fileWriteError := r.file.WriteString(string(jsonLink) + "\n")
		if fileWriteError != nil {
			return nil, fileWriteError
		}
	}

	return result, nil
}

// GetByCode -.
func (r *linksRepoInFile) GetLinkByCode(ctx context.Context, shortURL string) (string, error) {

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

// Close -.
func (r *linksRepoInFile) Close() {
	r.file.Close()
}

// Check -.
func (r *linksRepoInFile) Check() error {
	return nil
}
