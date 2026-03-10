package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	pdfzusmerge "github.com/wsgtcyx/pdf-zusammenfuegen-circleci-orb"
)

func main() {
	var outputDir string

	flag.StringVar(&outputDir, "output-dir", "artifacts", "Zielverzeichnis fuer Beispielartefakte")
	flag.Parse()

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		exitf("Ausgabeverzeichnis konnte nicht erstellt werden: %v", err)
	}

	first := filepath.Join(outputDir, "beispiel-1.pdf")
	second := filepath.Join(outputDir, "beispiel-2.pdf")
	merged := filepath.Join(outputDir, "merged-sample.pdf")

	if err := writeMinimalPDF(first, 210, 297); err != nil {
		exitf("beispiel-1.pdf konnte nicht geschrieben werden: %v", err)
	}
	if err := writeMinimalPDF(second, 297, 210); err != nil {
		exitf("beispiel-2.pdf konnte nicht geschrieben werden: %v", err)
	}

	if err := pdfzusmerge.MergeFiles(merged, []string{first, second}, pdfzusmerge.Options{
		Bookmarks: true,
		Divider:   false,
		Optimize:  true,
	}); err != nil {
		exitf("merged-sample.pdf konnte nicht erzeugt werden: %v", err)
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func writeMinimalPDF(path string, width, height int) error {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Count 1 /Kids [3 0 R] >>",
		fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %d %d] /Resources << >> /Contents 4 0 R >>", width, height),
		"<< /Length 0 >>\nstream\n\nendstream",
	}

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")

	offsets := make([]int, len(objects)+1)
	for index, object := range objects {
		offsets[index+1] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", index+1, object)
	}

	xrefOffset := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n", len(objects)+1)
	buf.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offset)
	}

	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefOffset)

	return os.WriteFile(path, buf.Bytes(), 0o644)
}
