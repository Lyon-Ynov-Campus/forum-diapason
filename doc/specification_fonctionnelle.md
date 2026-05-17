# Forum Diapason | Spécification Fonctionnelle

**École :** Ynov Campus Lyon  
**Promotion :** B1 Informatique 2025/2026  
**Durée du projet :** 27/04/2026 au 17/05/2026  
**Dernière mise à jour :** 17/05/2026  
**Rédigé par :** Michel LEVINE  
**Statut :** Livré

---

<details>
<summary>Sommaire</summary>

- [1. Introduction](#1-introduction)
  - [1.1. Présentation du projet](#11-présentation-du-projet)
  - [1.2. Périmètre](#12-périmètre)
  - [1.3. Parties prenantes](#13-parties-prenantes)
  - [1.4. Glossaire](#14-glossaire)
- [2. Fonctionnalités](#2-fonctionnalités)
  - [2.1. Comptes utilisateurs](#21-comptes-utilisateurs)
  - [2.2. Posts](#22-posts)
  - [2.3. Commentaires](#23-commentaires)
  - [2.4. Likes](#24-likes)
  - [2.5. Profil utilisateur](#25-profil-utilisateur)
  - [2.6. Recherche et filtres](#26-recherche-et-filtres)
  - [2.7. Navigation](#27-navigation)
- [3. Cas d'utilisation](#3-cas-dutilisation)
- [4. Maquettes et wireframes](#4-maquettes-et-wireframes)
  - [4.1. Charte graphique](#41-charte-graphique)
  - [4.2. Wireframe](#42-wireframe)
  - [4.3. Maquette haute fidélité](#43-maquette-haute-fidélité)
  - [4.4. Pages couvertes](#44-pages-couvertes)
- [5. Exigences non fonctionnelles](#5-exigences-non-fonctionnelles)
- [6. Conclusion](#7-conclusion)

</details>

---

## 1. Introduction

### 1.1. Présentation du projet

Forum Diapason est un forum web dédié à la musique, inspiré de Reddit. Le projet est réalisé dans le cadre de la formation B1 Informatique à Ynov Campus Lyon, dans des conditions professionnelles simulées.

L'objectif est de proposer une plateforme communautaire où les passionnés de musique peuvent créer un compte, partager leurs découvertes, avis et discussions musicales sous forme de posts, interagir avec les autres membres via des commentaires et des réactions, et naviguer sur un fil d'actualité. Le site est déployé et accessible en ligne.

| | |
|--|--|
| **Dépôt GitHub** | [forum-diapason](https://github.com/Lyon-Ynov-Campus/forum-diapason) |
| **URL de déploiement** | https://forum.prettyflacko.fr |

---

### 1.2. Périmètre

**Ce que le projet couvre :**

- Création de compte, connexion et déconnexion
- Réinitialisation du mot de passe par email
- Modification et suppression du profil utilisateur (avatar, bio, mot de passe)
- Suppression du compte
- Publication, affichage et suppression de posts avec images
- Système de tags sur les posts
- Système de commentaires avec réponses imbriquées
- Système de likes sur les posts
- Consultation du profil d'autres utilisateurs
- Recherche fulltext (titre, contenu, tags) avec filtres de tri
- Suggestions d'utilisateurs dans la barre de recherche
- Interface responsive avec dark mode
- Déploiement de l'application sur Azure

**Ce qui est hors périmètre :**

- Notifications en temps réel
- Messagerie privée entre utilisateurs
- Connexion via un service tiers (Google, GitHub, etc.)
- Système de modération ou de rôles administrateur

---

### 1.3. Parties prenantes

| Rôle | Nom | Attentes |
|------|-----|----------|
| Client | Mentor Kilian MOUN | Un forum fonctionnel livré dans les délais, avec une interface soignée et une gestion de projet suivie sur Jira |
| Équipe | Joaquin COZETTE, PHAM Huy, Michel LEVINE, Andry NOMENAFITIA | Concevoir et livrer l'application dans les délais |
| Utilisateurs finaux | Grand public | Un forum simple à prendre en main, rapide et agréable à utiliser |

---

### 1.4. Glossaire

| Terme | Définition |
|-------|------------|
| **Forum** | Plateforme web permettant à des utilisateurs d'échanger via des posts et des commentaires |
| **Post** | Message publié par un utilisateur, visible par tous les visiteurs |
| **Commentaire** | Réponse textuelle d'un utilisateur à un post |
| **Like** | Action exprimant un avis positif sur un post, réversible |
| **Tag** | Étiquette thématique associée à un post pour faciliter la recherche |
| **Fil d'actualité** | Page principale affichant les posts par ordre chronologique décroissant |
| **Session cookie** | Cookie HttpOnly stockant l'identifiant de session côté client, utilisé pour maintenir la connexion |
| **API REST** | Interface de programmation permettant la communication entre le frontend et le backend via des requêtes HTTP |
| **MVP** | Minimum Viable Product — version minimale fonctionnelle du produit |
| **SQLite** | Base de données légère utilisée pour stocker les données de l'application |
| **ACI** | Azure Container Instances — service Azure pour déployer des conteneurs Docker |

---

## 2. Fonctionnalités

### 2.1. Comptes utilisateurs

#### Inscription

Un visiteur peut créer un compte en renseignant son nom, un pseudo unique, une adresse email et un mot de passe. Si le pseudo ou l'email est déjà utilisé, un message d'erreur est affiché. Le mot de passe est stocké de manière sécurisée (haché avec bcrypt, jamais en clair). Une fois le compte créé, l'utilisateur est connecté automatiquement.

#### Connexion et déconnexion

L'utilisateur se connecte avec son email ou son pseudo et son mot de passe. Une session est créée et maintenue via un cookie HttpOnly. En cas d'identifiants incorrects, un message d'erreur s'affiche. L'utilisateur peut se déconnecter à tout moment, ce qui détruit la session.

#### Réinitialisation du mot de passe

Un utilisateur peut demander la réinitialisation de son mot de passe en renseignant son adresse email. Un lien de réinitialisation valable 1 heure lui est envoyé par email. Ce lien lui permet de définir un nouveau mot de passe.

#### Modification du profil

L'utilisateur connecté peut modifier son nom, son pseudo, son email, son avatar et son mot de passe depuis sa page de paramètres.

#### Suppression du compte

L'utilisateur peut supprimer son compte de façon définitive depuis les paramètres. Une confirmation avec le mot de passe est demandée. Une fois le compte supprimé, toutes ses données sont effacées et il est redirigé vers la page d'accueil.

---

### 2.2. Posts

#### Création d'un post

Un utilisateur connecté peut créer un post en renseignant un titre, un contenu, des tags et optionnellement une image. L'image peut être recadrée avant envoi. Elle est compressée côté client si elle dépasse 2 Mo. Un post avec un titre ou un contenu vide ne peut pas être publié.

#### Affichage des posts

Tous les visiteurs, connectés ou non, peuvent consulter les posts. La page d'accueil affiche un fil de posts avec titre, auteur, date, tags, image et compteurs. Une colonne "top posts" présente les 6 posts les plus likés.

#### Page d'un post

Un clic sur un post ouvre sa page dédiée, qui affiche le contenu complet, l'image, les commentaires imbriqués et le bouton de réaction.

#### Suppression d'un post

L'auteur d'un post peut le supprimer. Un post supprimé disparaît immédiatement du fil.

---

### 2.3. Commentaires

Un utilisateur connecté peut commenter un post et répondre aux commentaires existants. Les commentaires sont affichés du plus ancien au plus récent. L'auteur d'un commentaire peut le supprimer.

---

### 2.4. Likes

Un utilisateur connecté peut liker un post. Il ne peut liker un même post qu'une seule fois. Il peut retirer son like. Le compteur de likes est visible par tous.

---

### 2.5. Profil utilisateur

Chaque utilisateur dispose d'une page de profil consultable. Elle affiche l'avatar, le pseudo et la liste des posts publiés. Sur son propre profil, l'utilisateur voit les options de modification.

---

### 2.6. Recherche et filtres

La barre de recherche permet de filtrer les posts par mot-clé (titre, contenu, tags). Les filtres permettent de trier par récent, populaire ou tendance. Des suggestions d'utilisateurs apparaissent en temps réel lors de la saisie.

---

### 2.7. Navigation

Toutes les pages partagent un header commun avec barre de recherche, menu de navigation et toggle dark mode. Un bouton paramètres donne accès aux options du compte.

---

## 3. Cas d'utilisation

### CU-01 — Créer un compte

| | |
|--|--|
| **Acteur** | Visiteur non connecté |
| **Déclencheur** | Il clique sur "Register" |
| **Scénario** | Il remplit le formulaire (nom, pseudo, email, mot de passe), valide. Le système crée le compte et connecte l'utilisateur. |
| **Échec** | Si le pseudo ou l'email est déjà pris, un message d'erreur s'affiche. |

---

### CU-02 — Se connecter

| | |
|--|--|
| **Acteur** | Visiteur non connecté |
| **Déclencheur** | Il clique sur "Sign in" |
| **Scénario** | Il saisit son email ou pseudo et son mot de passe. Le système crée une session et redirige vers le fil. |
| **Échec** | Si les identifiants sont incorrects, un message d'erreur s'affiche. |

---

### CU-03 — Publier un post

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il clique sur le formulaire de création de post |
| **Scénario** | Il saisit titre, contenu, tags, image optionnelle et valide. Le post apparaît en tête du fil. |
| **Échec** | Si le titre ou le contenu est vide, une erreur s'affiche. |

---

### CU-04 — Commenter un post

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il ouvre un post et saisit un commentaire |
| **Scénario** | Il rédige son commentaire et valide. Il apparaît immédiatement sous le post. |
| **Échec** | Si le champ est vide, rien n'est envoyé. |

---

### CU-05 — Liker un post

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il clique sur le bouton like |
| **Scénario** | Le compteur s'incrémente et le bouton indique visuellement que le like est actif. Un second clic retire le like. |
| **Échec** | Si l'utilisateur n'est pas connecté, le like n'est pas enregistré. |

---

### CU-06 — Modifier ou supprimer son compte

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il accède aux paramètres via l'icône ⚙️ |
| **Scénario (modification)** | Il renseigne les nouvelles informations et sauvegarde. |
| **Scénario (suppression)** | Il confirme avec son mot de passe. Le compte et toutes les données sont supprimés. |
| **Échec** | Si le mot de passe est incorrect, un message d'erreur s'affiche. |

---

## 4. Maquettes et wireframes

### 4.1. Charte graphique

| Élément | Détail |
|---------|--------|
| Palette | Noir (#000), blanc (#fff), gris (#e5e7eb) |
| Typographie | System UI sans-serif |
| Style | Minimaliste, inspiré des forums typographiques |
| Thème | Clair et sombre (dark mode) |

---

### 4.2. Wireframe

Le wireframe a été réalisé sur Excalidraw et couvre les pages principales du site.

Lien : https://excalidraw.com/#room=9090ab9a1cfda3075dfa,BKUC3vwf2gLmmUpQgZKVyQ

---

### 4.3. Maquette haute fidélité

La maquette a été réalisée sur Figma par PHAM Huy et ANDRY Nomenafitia.

Lien : https://www.figma.com/design/E4SGgSbexa4pVr7Nh7Mnqw/maquette-diapason?node-id=40-4072&t=qFlOZp2EpggMAXBz-0

---

### 4.4. Pages couvertes

| Page | Description |
|------|-------------|
| Accueil | Fil de posts, barre de recherche, top posts, création de post |
| Connexion | Formulaire d'authentification |
| Inscription | Formulaire de création de compte |
| Page d'un post | Contenu complet, image, commentaires, like |
| Profil utilisateur | Avatar, pseudo, liste des posts |
| Paramètres | Modal : modification du compte, suppression |
| Mot de passe oublié | Formulaire de demande de reset |
| Réinitialisation | Formulaire de nouveau mot de passe |

---

## 5. Exigences non fonctionnelles

| Catégorie | Exigence |
|-----------|----------|
| Ergonomie | L'interface doit être intuitive et utilisable sans formation préalable |
| Responsivité | Le site doit fonctionner correctement sur desktop et mobile |
| Performance | Les pages doivent se charger rapidement |
| Sécurité | Mots de passe hachés avec bcrypt, sessions cookie HttpOnly, protection CORS |
| Déploiement | Application accessible depuis https://forum.prettyflacko.fr |
| Compatibilité | Le site fonctionne sur les navigateurs modernes (Chrome, Firefox, Edge) |
| Persistance | Les données (base SQLite, uploads) persistent entre les redémarrages via Azure File Share |

---

## 6. Conclusion

Forum Diapason offre un espace d'échange complet pour les passionnés de musique. Le projet a été livré avec l'ensemble des fonctionnalités prévues plus plusieurs ajouts : recherche fulltext, tags, images sur les posts, réinitialisation de mot de passe par email, dark mode et déploiement Azure avec HTTPS.

---

*Document rédigé par Michel LEVINE. Dernière révision : 17/05/2026.*
