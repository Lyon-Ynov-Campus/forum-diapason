# Forum Diapason

## Prérequis

- [Go](https://go.dev/dl/) 1.21+
- GCC (requis par go-sqlite3) — `sudo apt install gcc` sur Ubuntu
- [Air](https://github.com/air-verse/air) — live reload

```bash
go install github.com/air-verse/air@latest
```

- **Tailwind CLI** — à télécharger à la racine du projet

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64 && mv tailwindcss-linux-x64 tailwindcss
```

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

```
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
        └── input.css        # source Tailwind (styles.css est généré)
```
