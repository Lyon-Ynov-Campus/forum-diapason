# Forum Diapason 🎵

Un forum musical moderne inspiré de Reddit, construit avec Go et JavaScript vanilla.

## 🚀 Fonctionnalités

- **Publications musicales** : consultation des posts et création en cours de stabilisation
- **Système de likes** : like/unlike dynamique côté interface
- **Authentification JWT** : inscription/connexion fonctionnelles
- **Interface moderne** : design inspiré de Reddit avec un thème musical
- **API REST** : backend Go avec SQLite
- **Responsive** : interface adaptée desktop/mobile

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
│   ├── cmd/
│   │   ├── api/             # Serveur backend (API)
│   │   ├── web/             # Serveur frontend (fichiers statiques)
│   │   └── dev/             # Lance les 2 serveurs ensemble
│   ├── infrastructure/
│   │   ├── database/        # Connexion et migrations SQLite
│   │   └── http/handlers/   # Gestionnaires HTTP
│   ├── internal/
│   │   ├── entities/        # Modèles de données
│   │   ├── repositories/    # Accès données (SQLite)
│   │   └── usecases/        # Logique métier
│   ├── pkg/
│   │   ├── config/          # Configuration
│   │   └── logger/          # Logging
│   ├── go.mod
│   └── Makefile
├── frontend/
│   ├── index.html
│   ├── css/styles.css
│   ├── js/app.js
│   └── package.json
├── docs/                    # Documentation
├── .env.example
├── .gitignore
└── README.md
```

## 🛠️ Installation et démarrage

### Lancement rapide (depuis la racine)

Depuis `forum-diapason/`, vous pouvez lancer le projet sans changer de dossier :

```bash
go run ./backend/cmd/dev
```

Cette commande lance en parallèle :
- serveur backend API sur `http://localhost:8080`
- serveur frontend statique sur `http://localhost:3000`

### Prérequis
- Go 1.21+
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
   # Aucune dépendance Node.js nécessaire pour servir le frontend en local
   ```

2. **Démarrer le serveur de développement**
   ```bash
   go run ./backend/cmd/web
   ```

   Ouvrez `http://localhost:3000` dans votre navigateur

## 🔧 Configuration

### Variables d'environnement

Créez un fichier `.env` dans le dossier `backend/` en vous basant sur `backend/.env.example` :

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

> Statut: certains endpoints sont en cours de finalisation (auth middleware et protections).

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
- `GET /api/users/{userId}/posts` - Posts d'un utilisateur (**en attente d'implémentation complète**)

## 🧪 Tests

### Backend
```bash
cd backend
go test ./...
```

Statut actuel: la commande fonctionne mais il n'y a pas encore de fichiers de tests automatisés.

### Frontend
```bash
cd frontend
npm test
```

Statut actuel: script placeholder (tests automatisés frontend à implémenter).

## 🚢 Déploiement

### Backend
```bash
make build
./bin/forum-api
```

### Frontend
```bash
go run ./backend/cmd/web
```

Statut actuel: le frontend est servi comme fichiers statiques via un serveur Go dédié.

## 🧭 Processus projet et gouvernance

Cette section décrit notre manière de travailler dans des conditions professionnelles, tout en restant réalistes par rapport à l'état actuel du projet.

### Objectif du processus

Notre processus sert trois objectifs en même temps :
- **Pilotage interne** : savoir qui fait quoi, quand, et avec quel niveau de priorité
- **Visibilité client/PO** : montrer une progression claire et structurée
- **Onboarding** : permettre à un nouveau membre de comprendre rapidement notre organisation

### Cadre d'équipe

- **Chef de projet** : Cozette Joaquin
- **Membres** : PHAM Quoc Huy, Andry Nomenafitia, Levine-Riff Michel
- **Validation avant présentation** : revue collective de l'équipe, avec validation finale prioritairement assurée par le chef de projet
- **Répartition des responsabilités** : pas de rôles secondaires fixes pour le moment (organisation simple et collaborative)

### Outils de coordination

- **Jira** : gestion des tickets et suivi de l'état des tâches
- **Discord** : communication rapide et partage de documents
- **Travail en présentiel** : coordination principale pendant les sessions de cours

### Cycle de travail (version actuelle)

1. Définition des objectifs au début du cours (et consolidation hebdomadaire)
2. Priorisation collective des tâches
3. Arbitrage final des priorités par le chef de projet
4. Exécution technique (design/implémentation/tests de base)
5. Vérification d'équipe (et surtout chef de projet) puis passage en terminé

### Workflow de tickets (Jira)

Colonnes utilisées actuellement :
- **En cours**
- **Bug à fixe**
- **Terminé**

Règle de passage :
- Une tâche est considérée **terminée** lorsqu'elle est vérifiée par l'équipe, avec validation principale par le chef de projet.

### Processus Git (actuel)

- **Stratégie de branches** : `main` / `dev` / `feature/*`
- **Commits** : usage recommandé des Conventional Commits
- **Pull Requests** : recommandées mais non systématiques à ce stade
- **Intégration** : merge des `feature/*` vers `dev`
- **Checks automatiques** : non mis en place pour le moment
- **Politique hotfix** : non formalisée pour l'instant

### Processus de développement

#### Frontend/UI
- Le design est en cours sur Figma
- L'implémentation se fait ensuite par fonctionnalité
- Les retours sont intégrés au fil de l'eau

#### Backend/API
- API Go connectée à SQLite
- Développement orienté fonctionnalités prioritaires (authentification et espace utilisateur en premier)
- Sécurité déjà engagée : hash des mots de passe (mécanismes avancés anti-abus à renforcer)

### Qualité et validation (phase actuelle)

Le projet est encore en phase de lancement. Le niveau de qualité attendu aujourd'hui est :
- un socle fonctionnel stable pour les premiers parcours utilisateur,
- une organisation claire de l'équipe,
- un suivi des bugs dans Jira.

Le processus de "recette" (validation formelle avant démo) sera progressivement renforcé avec une checklist dédiée.

### Environnements et déploiement

- **Environnement courant** : local
- **Base de données** : SQLite (via API backend)
- **Déploiement** : en réflexion
- **Décision de mise en ligne** : validation collective de l'équipe

### Priorisation fonctionnelle (MVP vs suite)

#### Priorités MVP (maintenant)
Les fonctionnalités suivantes sont traitées comme prioritaires pour une base exploitable :
- authentification (inscription/connexion),
- gestion du profil utilisateur (minimum utile),
- création et consultation de posts,
- base des interactions (likes/commentaires selon avancement),
- structure générale de l'application (home, navigation, architecture API + DB).

#### Fonctionnalités suivantes (après socle)
- tri/filtrage/recherche avancés,
- section réseau complète (profils et exploration),
- durcissement qualité (tests plus systématiques),
- stratégie de déploiement et exploitation.

### Suivi d'avancement des processus

> Mettre à jour ce tableau au début de chaque cours.

| Processus | État | Responsable principal | Dernière mise à jour | Prochaine action |
|---|---|---|---|---|
| Gouvernance d'équipe | En cours | Cozette Joaquin | À renseigner | Maintenir le rythme de coordination |
| Gestion Jira | En cours | Équipe + CP | À renseigner | Continuer la mise à jour des tickets |
| Workflow Git | En cours | Équipe | À renseigner | Rendre les PR plus systématiques |
| Développement Auth/Utilisateur | En cours | Équipe | À renseigner | Stabiliser les premiers parcours |
| Gestion des bugs | En cours | Équipe + CP | À renseigner | Traiter les tickets "Bug à fixe" |
| Recette formelle avant démo | À venir | Cozette Joaquin | À renseigner | Définir une checklist de validation |
| Déploiement | À venir | Équipe | À renseigner | Choisir une stratégie de déploiement |

### Règles de mise à jour de cette section

- Cette documentation est **vivante** : elle évolue tout au long du projet.
- Mise à jour prévue **au début de chaque cours**.
- Responsable de mise à jour : **Cozette Joaquin**.
- En cas de changement d'organisation, mettre à jour en priorité :
  1. les rôles,
  2. les statuts du tableau,
  3. les priorités MVP.

## 🤝 Contribution

1. Fork le projet
2. Créez une branche feature (`git checkout -b feature/AmazingFeature`)
3. Committez vos changements (`git commit -m 'Add some AmazingFeature'`)
4. Pushez vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrez une Pull Request

## 📝 Licence

Licence en cours de définition (fichier `LICENSE` non ajouté pour le moment).

## 🙏 Remerciements

- Inspiré par Reddit et les communautés musicales
- Icônes par [Font Awesome](https://fontawesome.com/)
- Framework HTTP Go par [Gorilla](https://github.com/gorilla/)

---

**Forum Diapason** - Partagez votre passion pour la musique ! 🎶
