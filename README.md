# Forum Diapason

## Prérequis

### Go 1.25+

```bash
VERSION=$(curl -s "https://go.dev/VERSION?m=text" | head -1)
curl -LO "https://go.dev/dl/${VERSION}.linux-amd64.tar.gz"
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf ${VERSION}.linux-amd64.tar.gz
rm ${VERSION}.linux-amd64.tar.gz
```

Ajouter dans `~/.bashrc` :

```bash
export PATH=/usr/local/go/bin:$PATH
```

### GCC (requis par go-sqlite3)

```bash
sudo apt install gcc
```

### Air — live reload

```bash
go install github.com/air-verse/air@latest
```

Vérifier que `~/go/bin` est dans le PATH (`air` doit être accessible) :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Tailwind CLI — à télécharger à la racine du projet

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64 && mv tailwindcss-linux-x64 tailwindcss
```

**Après avoir téléchargé Tailwind** (ou après chaque `git pull`), générer le CSS :

```bash
./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --content "./frontend/**/*.html,./frontend/**/*.js"
```

> `make dev` le regénère automatiquement à chaque démarrage.

---

## Lancer le projet

```bash
make dev
```

Build pour la prod :

```bash
make build
./forum-diapason
```

---

## Structure

```text
forum-diapason/
├── main.go                  # serveur HTTP + routes
├── go.mod
├── Makefile
│
├── database/
│   └── database.go          # connexion SQLite + migrations
│
├── models/
│   ├── user.go
│   └── post.go
│
├── handlers/                # handlers HTTP (un fichier par feature)
│   ├── handlers.go
│   ├── auth.go
│   └── post.go
│
├── services/                # logique métier (un fichier par feature)
│   ├── auth_service.go
│   └── post_service.go
│
├── utils/
│   └── utils.go
│
└── frontend/
    ├── components/          # composants réutilisables (header, footer…)
    ├── pages/               # une page HTML par route
    ├── js/
    │   └── app.js
    └── css/
        ├── input.css        # source Tailwind (à modifier)
        └── styles.css       # généré — ne pas commiter
```
