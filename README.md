# Forum Diapason

Forum musical permettant de partager des découvertes, des avis et des discussions autour de la musique.

> Pour le déploiement en production (Azure), voir la branche `azure-deploy`.

---

## Stack

- **Backend** : Go stdlib (`net/http`, `html/template`)
- **Base de données** : SQLite via `go-sqlite3` (CGO)
- **Frontend** : HTML, JavaScript vanilla, Tailwind CSS v4 (CLI standalone)
- **Auth** : Sessions cookie (`SameSite=Lax`, `HttpOnly`)
- **Email** : SMTP (reset de mot de passe)

---

## Fonctionnalités

- Inscription / connexion / déconnexion
- Réinitialisation du mot de passe par email
- Paramètres du compte (modification, suppression)
- Posts avec titre, contenu, tags et image (recadrage + compression côté client)
- Likes, commentaires et réponses aux commentaires
- Page profil avec avatar, bio et liste des posts
- Recherche fulltext (titre, contenu, tags) et filtres (tri, tags)
- Suggestions d'utilisateurs dans la barre de recherche
- Follows entre utilisateurs
- Dark mode persistant

---

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

Vérifier que `~/go/bin` est dans le PATH :
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Tailwind CSS CLI

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64 && mv tailwindcss-linux-x64 tailwindcss
```

---

## Configuration

Copier `.env.example` en `.env` et remplir les valeurs :
```bash
cp .env.example .env
```

```env
PORT=8080
API_PORT=8081
DB_FILE=./forum.db
COOKIE_SECURE=false
COOKIE_DOMAIN=
FRONTEND_ORIGIN=http://localhost:8080
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=
SMTP_PASS=
SMTP_FROM=
```

> Sans SMTP configuré, les liens de reset password s'affichent dans les logs du serveur API.

---

## Lancer le projet

```bash
make dev
```

Lance en parallèle : Tailwind (watch), serveur API (port 8081) et serveur front avec live reload via Air (port 8080).

Ouvrir http://localhost:8080 dans le navigateur.

---

## Structure

```
forum-diapason/
├── main.go                      # entrée du serveur front
├── go.mod
├── Makefile
├── .env.example
│
├── api/                         # serveur API
│   ├── main.go
│   ├── auth.go
│   ├── routers.go
│   ├── profile.go
│   └── helpers.go
│
├── database/
│   └── database.go              # connexion SQLite + migrations
│
├── models/
│   ├── user.go
│   └── post.go
│
├── handlers/                    # handlers du serveur front
│   ├── handlers.go
│   ├── auth.go
│   ├── page.go
│   ├── profile.go
│   └── search.go
│
├── services/
│   ├── auth_service.go
│   ├── post_service.go
│   └── user_service.go
│
├── utils/
│   ├── utils.go                 # sessions, cookies
│   └── mail.go                  # envoi SMTP
│
└── frontend/
    ├── components/              # composants Go template réutilisables
    ├── pages/                   # une page par route
    ├── js/
    │   ├── app.js               # logique globale, modals, auth
    │   ├── posts.js             # chargement et rendu des posts
    │   ├── post.js              # page détail d'un post
    │   └── profile.js           # page profil
    └── css/
        ├── input.css            # source Tailwind (à modifier)
        └── styles.css           # généré — ne pas commiter
```
