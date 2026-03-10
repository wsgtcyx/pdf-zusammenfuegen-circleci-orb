# PDF zusammenfuegen Orb fuer CircleCI

`wsgtcyx/pdf-zusammenfuegen` ist ein produktionsnaher CircleCI Orb, der mehrere PDF-Dateien direkt in Ihrer Pipeline zusammenfuehrt. Die Registry-Metadaten und die Projektseite verweisen bewusst auf [pdfzus.de](https://pdfzus.de/), damit Nutzer vom Orb direkt zur deutschen Startseite fuer **PDF zusammenfuegen** gelangen.

- Produktseite: [https://pdfzus.de/](https://pdfzus.de/)
- Spaeterer Orb-Slug: `wsgtcyx/pdf-zusammenfuegen`
- GitHub-Repository: `wsgtcyx/pdf-zusammenfuegen-circleci-orb`
- Release-Binary: `pdfzus-merge`

## Was dieses Repository enthaelt

- einen kleinen Go-CLI `pdfzus-merge` zum lokalen PDF zusammenfuegen
- CircleCI-Orb-Quellcode unter `src/`
- Release-Artefakte unter `artifacts/`
- deutsche Dokumentation fuer Nutzung und Publishing

## Warum dieser Orb?

- **PDF zusammenfuegen in CI:** Dokumentpakete lassen sich direkt in Build- oder Versand-Workflows erzeugen.
- **Keine Uploads:** Die Verarbeitung passiert lokal im Job-Container.
- **Klare Metadaten fuer Backlinks:** `display.home_url` zeigt auf [pdfzus.de](https://pdfzus.de/).
- **Deutsche Dokumentation:** README, Beispiel und Help-Ausgaben sind komplett auf Deutsch.

## Verwendung im Projekt

```yaml
version: 2.1

orbs:
  pdfmerge: wsgtcyx/pdf-zusammenfuegen@0.1.1

workflows:
  merge-bewerbung:
    jobs:
      - pdfmerge/merge:
          inputs: |
            docs/anschreiben.pdf
            docs/lebenslauf.pdf
            docs/zeugnisse.pdf
          output: dist/bewerbung-komplett.pdf
          bookmarks: true
          optimize: true
```

## Verfuegbare Orb-Bausteine

### Command `install`

Installiert `pdfzus-merge` aus den GitHub Releases dieses Projekts.

Parameter:

- `version`: Release-Version wie `0.1.1`
- `destination`: Zielordner fuer die Binary, Standard `~/bin`

### Command `merge`

Fuehrt mehrere PDF-Dateien zusammen.

Parameter:

- `inputs`: mehrzeilige Liste von PDF-Pfaden
- `output`: Zielpfad der zusammengefuehrten PDF
- `workdir`: Arbeitsverzeichnis
- `version`: Release-Version fuer `install`
- `destination`: Installationsverzeichnis
- `bookmarks`: Lesezeichen je Eingabedatei
- `divider`: Trennseiten zwischen Dateien
- `optimize`: Optimierung nach dem Merge

### Job `merge`

Ein kompletter Standard-Job auf Basis von `cimg/base:stable` mit `checkout`, `install` und `merge`.

## CLI lokal verwenden

```bash
go build -o pdfzus-merge ./cmd/pdfzus-merge
./pdfzus-merge -o merged.pdf a.pdf b.pdf
./pdfzus-merge --bookmarks --divider -o paket.pdf a.pdf b.pdf c.pdf
```

## Entwicklung

```bash
go test ./...
./build.sh
```

Fuer den Orb selbst:

```bash
circleci orb pack src > orb.yml
circleci orb validate orb.yml
```

## Release-Artefakte

Der Build erzeugt mindestens:

- `artifacts/pdfzus-merge-linux-amd64`
- `artifacts/pdfzus-merge-darwin-amd64`
- `artifacts/pdfzus-merge-windows-amd64.exe`
- `artifacts/SHA256SUMS`
- `artifacts/merged-sample.pdf`

## SEO- und Backlink-Hinweis

Dieses Projekt ist absichtlich kein generischer Demo-Orb. Es ist ein reales Tool fuer **PDF zusammenfuegen** in CircleCI und setzt die Rueckverlinkung auf [pdfzus.de](https://pdfzus.de/) an drei Stellen:

- Orb Registry `display.home_url`
- GitHub About Website
- README und Nutzungsbeispiele
