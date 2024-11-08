package ws

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Ensure MessageFromFiles implements MessageHandler
var _ MessageHandler = (*MessageFromFiles)(nil)

type MessageFromFiles struct {
	dirPath string
}

func NewMessageFromFiles(dirPath string) MessageFromFiles {
	return MessageFromFiles{
		dirPath: dirPath,
	}
}

func (r MessageFromFiles) Handle(ctx context.Context, _ Message, connClient Connection, _ Connection) error {
	entries, err := os.ReadDir(r.dirPath)
	if err != nil {
		return err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		files = append(files, e.Name())
	}

	if len(files) == 0 {
		return nil
	}

	t0, err := parseTime(files[0])
	if err != nil {
		return err
	}

	startTime := time.Now()

	for _, f := range files {
		if err := context.Cause(ctx); err != nil {
			return err
		}

		ti, err := parseTime(f)
		if err != nil {
			return err
		}

		message, err := ReadFromFile(filepath.Join(r.dirPath, f))
		if err != nil {
			return err
		}

		interval := ti.Sub(t0)
		time.Sleep(time.Until(startTime.Add(interval)))

		connClient.Write(ctx, message)
	}
	return nil
}

func parseTime(filename string) (time.Time, error) {
	f, _ := strings.CutSuffix(filename, ".yaml")
	return time.Parse(time.RFC3339Nano, f)
}
