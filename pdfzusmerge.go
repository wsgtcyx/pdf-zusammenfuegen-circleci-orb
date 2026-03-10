package pdfzusmerge

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Options steuert das Verhalten beim Zusammenfuegen.
type Options struct {
	Bookmarks bool
	Divider   bool
	Optimize  bool
}

// MergeFiles fuegt mindestens zwei lokale PDF-Dateien zu einer neuen Datei zusammen.
func MergeFiles(outputPath string, inputPaths []string, opts Options) error {
	outputPath = strings.TrimSpace(outputPath)
	if outputPath == "" {
		return errors.New("der Ausgabepfad darf nicht leer sein")
	}

	if len(inputPaths) < 2 {
		return errors.New("mindestens zwei PDF-Dateien sind erforderlich")
	}

	normalizedInputs := make([]string, 0, len(inputPaths))
	for _, inputPath := range inputPaths {
		normalizedPath, err := validateInputFile(inputPath)
		if err != nil {
			return err
		}
		normalizedInputs = append(normalizedInputs, normalizedPath)
	}

	if !strings.EqualFold(filepath.Ext(outputPath), ".pdf") {
		return fmt.Errorf("die Ausgabedatei muss auf .pdf enden: %s", outputPath)
	}

	conf := model.NewDefaultConfiguration()
	conf.CreateBookmarks = opts.Bookmarks
	conf.Optimize = opts.Optimize
	conf.OptimizeBeforeWriting = opts.Optimize

	if err := api.MergeCreateFile(normalizedInputs, outputPath, opts.Divider, conf); err != nil {
		return fmt.Errorf("pdfs konnten nicht zusammengefuegt werden: %w", err)
	}

	return nil
}

func validateInputFile(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("ein Eingabepfad ist leer")
	}

	if !strings.EqualFold(filepath.Ext(path), ".pdf") {
		return "", fmt.Errorf("nur PDF-Dateien sind erlaubt: %s", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("datei nicht gefunden: %s", path)
		}
		return "", fmt.Errorf("datei konnte nicht gelesen werden: %s", path)
	}

	if info.IsDir() {
		return "", fmt.Errorf("eingabepfad ist ein Verzeichnis und keine PDF-Datei: %s", path)
	}

	return path, nil
}
