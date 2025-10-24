package persist

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	log "search_engine/crawler/logger"
	"search_engine/crawler/metrics"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	RAND_BYTES_LENGTH = 8
	PERMISSIONS       = 0644
	OUTPUT_PATH       = "./output/"
)

type Content struct {
	URL     string
	Payload []byte
}

type TextFileSaver struct{}
type InvertedIndex struct {
	Metrics *metrics.Metrics
}

type Save interface {
	SaveToFile(context.Context, chan Content, *sync.WaitGroup)
}

func (t *TextFileSaver) SaveToFile(ctx context.Context, input chan Content, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a directory to store output files
	if err := os.MkdirAll(OUTPUT_PATH, os.ModePerm); err != nil {
		log.Error("Error creating output directory=", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Warn("Stopping SaveToFile due to ctx cancelled")
			return
		case content, ok := <-input:
			if !ok {
				return
			}

			filename := OUTPUT_PATH + t.generateRandomFilename()
			if err := os.WriteFile(filename, content.Payload, PERMISSIONS); err != nil {
				log.Error("Error WriteFile=", err)
			} else {
				// log.Info("Saved filename:", filename)
			}
		}
	}
}

func (t *InvertedIndex) SaveToFile(ctx context.Context, input chan Content, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	// Create a directory to store output files
	if err := os.MkdirAll(OUTPUT_PATH, os.ModePerm); err != nil {
		log.Error("Error creating output directory=", err)
		return
	}

	wordRegex := regexp.MustCompile(`\b\w+\b`)
	isAlpha := regexp.MustCompile(`^[a-zA-Z]+$`)

	for c := range input {
		select {
		case <-ctx.Done():
			return
		default:
			text := extractText(c.Payload)
			words := wordRegex.FindAllString(strings.ToLower(text), -1)

			for _, word := range words {
				if !isAlpha.MatchString(word) {
					continue
				}
				wordFile := filepath.Join(OUTPUT_PATH, word+".txt")

				// Open file in append mode, create if not exists
				f, err := os.OpenFile(wordFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					t.Metrics.IncrementFailure()
					continue
				}

				t.Metrics.IncrementSuccess()
				t.Metrics.AddTotalTime(time.Since(start))

				writer := bufio.NewWriter(f)
				writer.WriteString(c.URL + "\n")
				writer.Flush()
				f.Close()
			}
		}
	}
}

func (t *TextFileSaver) generateRandomFilename() string {
	randBytes := make([]byte, RAND_BYTES_LENGTH)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes) + ".txt"
}

func extractText(htmlContent []byte) string {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return ""
	}
	var buf strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				buf.WriteString(text + " ")
			}
		}
		// skip script and style tags
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return buf.String()
}
