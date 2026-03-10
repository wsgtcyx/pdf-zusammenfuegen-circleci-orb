package main

import (
	"flag"
	"fmt"
	"os"

	pdfzusmerge "github.com/wsgtcyx/pdf-zusammenfuegen-circleci-orb"
)

func main() {
	var (
		outputPath string
		bookmarks  bool
		divider    bool
		noOptimize bool
	)

	flag.StringVar(&outputPath, "o", "", "Ausgabedatei fuer das zusammengefuegte PDF")
	flag.BoolVar(&bookmarks, "bookmarks", false, "Lesezeichen fuer die Eingabedateien erzeugen")
	flag.BoolVar(&divider, "divider", false, "Trennseiten zwischen den Eingabedateien einfuegen")
	flag.BoolVar(&noOptimize, "no-optimize", false, "PDF nach dem Merge nicht optimieren")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "PDF zusammenfuegen mit pdfzus-merge\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Verwendung:\n  pdfzus-merge -o merged.pdf a.pdf b.pdf [c.pdf ...]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Optionen:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "Fehler: Bitte geben Sie mit -o eine Ausgabedatei an.")
		flag.Usage()
		os.Exit(2)
	}

	opts := pdfzusmerge.Options{
		Bookmarks: bookmarks,
		Divider:   divider,
		Optimize:  !noOptimize,
	}

	if err := pdfzusmerge.MergeFiles(outputPath, flag.Args(), opts); err != nil {
		fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Erfolg: %s wurde erstellt.\n", outputPath)
}
