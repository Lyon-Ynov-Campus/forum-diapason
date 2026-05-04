# Forum Diapason 🎵

Un forum musical moderne inspiré de Reddit, construit avec Go et JavaScript vanilla.

## 🚀 Fonctionnalités

- **Publications musicales** : Partagez vos découvertes, avis et discussions
- **Système de likes** : Aimez les publications qui vous plaisent
- **Authentification sécurisée** : Inscription/connexion avec JWT
- **Interface moderne** : Design inspiré de Reddit avec un thème musical
- **API REST** : Backend robuste en Go avec SQLite
- **Responsive** : Fonctionne sur desktop et mobile

## 🏗️ Architecture

### Backend (Go)
- **Framework** : Gorilla Mux pour le routing
- **Base de données** : SQLite avec migrations automatiques
- **Authentification** : JWT tokens
- **Architecture** : Clean Architecture (entités, repositories, usecases, handlers)

### Frontend (Vanilla JS)
- **HTML5/CSS3** : Interface moderne et responsive
- **JavaScript ES6+** : Logique frontend sans frameworks
- **API REST** : Communication avec le backend

## 📁 Structure du projet

```
forum-diapason/
├── backend/
│   ├── cmd/api/           # Point d'entrée de l'API
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── entities/  # Modèles de données
│   │   │   └── repositories/ # Interfaces de données
│   │   └── infrastructure/
│   │       ├── database/  # Connexion et migrations
│   │       └── http/handlers/ # Gestionnaires HTTP
│   ├── pkg/
│   │   ├── config/        # Configuration
│   │   └── logger/        # Logging
│   ├── usecases/          # Logique métier
│   ├── go.mod
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── html/          # Templates HTML
│   │   ├── css/           # Styles
│   │   └── js/            # JavaScript
│   └── package.json
├── docs/                  # Documentation
├── tests/                 # Tests
├── migrations/            # Scripts de migration
├── .gitignore
└── README.md
```

## 🛠️ Installation et démarrage

### Prérequis
- Go 1.21+
- Node.js 16+ (pour les outils frontend)
- SQLite3

### Backend

1. **Installer les dépendances**
   ```bash
   cd backend
   go mod download
   ```

2. **Démarrer le serveur**
   ```bash
   make run
   # ou
   go run cmd/api/main.go
   ```

   Le serveur démarre sur `http://localhost:8080`

### Frontend

1. **Installer les dépendances**
   ```bash
   cd frontend
   npm install
   ```

2. **Démarrer le serveur de développement**
   ```bash
   npm run dev
   ```

   Ouvrez `http://localhost:3000` dans votre navigateur

## 🔧 Configuration

### Variables d'environnement

Créez un fichier `.env` dans le dossier `backend/` :

```env
PORT=8080
DB_DRIVER=sqlite3
DB_FILE=./forum.db
JWT_SECRET=votre-cle-secrete-très-longue-et-complexe
```

### Base de données

Les migrations s'exécutent automatiquement au démarrage. Le schéma inclut :
- `users` : Utilisateurs du forum
- `posts` : Publications
- `comments` : Commentaires (à implémenter)
- `likes` : Système de likes

## 📡 API Endpoints

### Authentification
- `POST /api/auth/register` - Inscription
- `POST /api/auth/login` - Connexion

### Publications
- `GET /api/posts` - Lister les publications
- `POST /api/posts` - Créer une publication
- `GET /api/posts/{id}` - Obtenir une publication
- `PUT /api/posts/{id}` - Modifier une publication
- `DELETE /api/posts/{id}` - Supprimer une publication
- `POST /api/posts/{id}/like` - Aimer une publication
- `DELETE /api/posts/{id}/like` - Ne plus aimer

### Utilisateurs
- `GET /api/users/profile` - Profil utilisateur
- `PUT /api/users/profile` - Modifier le profil
- `PUT /api/users/change-password` - Changer le mot de passe

## 🧪 Tests

### Backend
```bash
cd backend
go test ./...
```

### Frontend
```bash
cd frontend
npm test
```

## 🚢 Déploiement

### Backend
```bash
make build
./bin/forum-api
```

### Frontend
```bash
cd frontend
npm run build
# Servir les fichiers statiques
```

## 🤝 Contribution

1. Fork le projet
2. Créez une branche feature (`git checkout -b feature/AmazingFeature`)
3. Committez vos changements (`git commit -m 'Add some AmazingFeature'`)
4. Pushez vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrez une Pull Request

## 📝 Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.

## 🙏 Remerciements

- Inspiré par Reddit et les communautés musicales
- Icônes par [Font Awesome](https://fontawesome.com/)
- Framework HTTP Go par [Gorilla](https://github.com/gorilla/)

---

**Forum Diapason** - Partagez votre passion pour la musique ! 🎶
