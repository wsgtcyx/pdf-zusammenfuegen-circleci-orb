package pdfzusmerge

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func TestMergeFilesMergesTwoPDFs(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputA := filepath.Join(dir, "a.pdf")
	inputB := filepath.Join(dir, "b.pdf")
	output := filepath.Join(dir, "merged.pdf")

	writeMinimalPDF(t, inputA, 200, 400)
	writeMinimalPDF(t, inputB, 300, 400)

	err := MergeFiles(output, []string{inputA, inputB}, Options{Optimize: true})
	if err != nil {
		t.Fatalf("MergeFiles returned error: %v", err)
	}

	pageCount, err := api.PageCountFile(output)
	if err != nil {
		t.Fatalf("PageCountFile returned error: %v", err)
	}

	if pageCount != 2 {
		t.Fatalf("expected 2 pages, got %d", pageCount)
	}

	if err := api.ValidateFile(output, nil); err != nil {
		t.Fatalf("ValidateFile returned error: %v", err)
	}
}

func TestMergeFilesPreservesInputOrder(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputA := filepath.Join(dir, "first.pdf")
	inputB := filepath.Join(dir, "second.pdf")
	output := filepath.Join(dir, "merged.pdf")

	writeMinimalPDF(t, inputA, 210, 400)
	writeMinimalPDF(t, inputB, 420, 400)

	err := MergeFiles(output, []string{inputA, inputB}, Options{Optimize: false})
	if err != nil {
		t.Fatalf("MergeFiles returned error: %v", err)
	}

	dims, err := api.PageDimsFile(output)
	if err != nil {
		t.Fatalf("PageDimsFile returned error: %v", err)
	}

	if len(dims) != 2 {
		t.Fatalf("expected 2 page dimensions, got %d", len(dims))
	}

	if int(dims[0].Width) != 210 || int(dims[1].Width) != 420 {
		t.Fatalf("unexpected page order: got widths %.0f then %.0f", dims[0].Width, dims[1].Width)
	}
}

func TestMergeFilesRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	validPDF := filepath.Join(dir, "a.pdf")
	output := filepath.Join(dir, "merged.pdf")
	textFile := filepath.Join(dir, "not-a-pdf.txt")

	writeMinimalPDF(t, validPDF, 200, 400)

	if err := os.WriteFile(textFile, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	testCases := []struct {
		name   string
		inputs []string
		want   string
	}{
		{
			name:   "too-few-files",
			inputs: []string{validPDF},
			want:   "mindestens zwei PDF-Dateien",
		},
		{
			name:   "missing-file",
			inputs: []string{validPDF, filepath.Join(dir, "missing.pdf")},
			want:   "datei nicht gefunden",
		},
		{
			name:   "wrong-extension",
			inputs: []string{validPDF, textFile},
			want:   "nur PDF-Dateien sind erlaubt",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := MergeFiles(output, tc.inputs, Options{})
			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error to contain %q, got %q", tc.want, err.Error())
			}
		})
	}
}

func TestCLIProducesMergedPDF(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	inputA := filepath.Join(dir, "a.pdf")
	inputB := filepath.Join(dir, "b.pdf")
	output := filepath.Join(dir, "cli-merged.pdf")
	binaryPath := filepath.Join(dir, "pdfzus-merge")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	writeMinimalPDF(t, inputA, 200, 400)
	writeMinimalPDF(t, inputB, 300, 400)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/pdfzus-merge")
	buildCmd.Dir = "."
	buildCmd.Env = os.Environ()
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(buildOutput))
	}

	runCmd := exec.Command(binaryPath, "-o", output, inputA, inputB)
	runCmd.Dir = "."
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cli execution failed: %v\n%s", err, string(runOutput))
	}

	if _, err := os.Stat(output); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}

	pageCount, err := api.PageCountFile(output)
	if err != nil {
		t.Fatalf("PageCountFile returned error: %v", err)
	}

	if pageCount != 2 {
		t.Fatalf("expected 2 pages, got %d", pageCount)
	}
}

func writeMinimalPDF(t *testing.T, path string, width, height int) {
	t.Helper()

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

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
