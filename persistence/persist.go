package persist

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
	log "websiteCopier/logger"
)

const (
	RAND_BYTES_LENGTH = 8
	PERMISSIONS       = 0644
	OUTPUT_PATH       = "./output/"
)

type Save interface {
	SaveToFile(context.Context, chan []byte, *sync.WaitGroup)
}

type TextFileSaver struct{}

func (t *TextFileSaver) SaveToFile(ctx context.Context, input chan []byte, wg *sync.WaitGroup) {
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
			if err := os.WriteFile(filename, content, PERMISSIONS); err != nil {
				log.Error("Error WriteFile=", err)
			} else {
				// log.Info("Saved filename:", filename)
			}
		}
	}
}

func (t *TextFileSaver) generateRandomFilename() string {
	randBytes := make([]byte, RAND_BYTES_LENGTH)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes) + ".txt"
}
