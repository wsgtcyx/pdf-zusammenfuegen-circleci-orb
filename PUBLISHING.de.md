# Publishing fuer GitHub und CircleCI Orbs

Dieses Projekt wird als oeffentliches Repository `wsgtcyx/pdf-zusammenfuegen-circleci-orb` veroeffentlicht. Der eigentliche Orb-Slug lautet `wsgtcyx/pdf-zusammenfuegen`.

## Voraussetzungen

- funktionierendes Go lokal
- `gh` CLI mit Login fuer `wsgtcyx`
- `circleci` CLI mit gueltigem Personal Token
- CircleCI Namespace `wsgtcyx` ist in Ihrem Konto verknuepft

## Lokale Pruefung

```bash
cd reference-repos/circleci-orb-pdf-zusammenfuegen
go test ./...
./build.sh
circleci orb pack src > orb.yml
circleci orb validate orb.yml
./artifacts/pdfzus-merge-linux-amd64 --help
```

## GitHub-Repository anlegen

```bash
git init
git branch -M main
git add .
git commit -m "init: circleci orb fuer pdf zusammenfuegen"
gh repo create wsgtcyx/pdf-zusammenfuegen-circleci-orb --public --source=. --remote=origin --push
gh repo edit wsgtcyx/pdf-zusammenfuegen-circleci-orb --homepage "https://pdfzus.de/" --description "CircleCI Orb zum PDF zusammenfuegen ohne Uploads"
```

## Release 0.1.0 erstellen

```bash
git tag v0.1.0
git push origin main --tags
gh release create v0.1.0 artifacts/* --repo wsgtcyx/pdf-zusammenfuegen-circleci-orb --title "v0.1.0" --notes "Erste oeffentliche Version des CircleCI Orbs zum PDF zusammenfuegen."
```

## CircleCI Orb veroeffentlichen

Falls der Orb noch nicht existiert:

```bash
circleci orb create wsgtcyx/pdf-zusammenfuegen
```

Dann erst eine Dev-Version publizieren:

```bash
circleci orb publish orb.yml wsgtcyx/pdf-zusammenfuegen@dev:first
```

Die Dev-Version mit einer echten Pipeline pruefen und danach promoten:

```bash
circleci orb publish promote wsgtcyx/pdf-zusammenfuegen@dev:first patch
```

## Nachkontrolle

- GitHub-Repository ist oeffentlich
- GitHub About zeigt `https://pdfzus.de/`
- Release `v0.1.0` enthaelt die drei Binaries und `SHA256SUMS`
- Orb Registry zeigt `display.home_url` auf `https://pdfzus.de/`
- Orb Example ist auf Deutsch und enthaelt `PDF zusammenfuegen`

