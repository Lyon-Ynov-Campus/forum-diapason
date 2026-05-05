# Forum Diapason | Spécification Fonctionnelle

**École :** Ynov Campus Lyon  
**Promotion :** B1 Informatique 2025/2026  
**Durée du projet :** 27/04/2026 au 11/05/2026  
**Dernière mise à jour :** 08/05/2026  
**Rédigé par :** Michel LEVINE  
**Statut :** En cours

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
  - [2.6. Navigation](#26-navigation)
- [3. Cas d'utilisation](#3-cas-dutilisation)
- [4. Maquettes et wireframes](#4-maquettes-et-wireframes)
  - [4.1. Charte graphique](#41-charte-graphique)
  - [4.2. Wireframe](#42-wireframe)
  - [4.3. Maquette haute fidélité](#43-maquette-haute-fidélité)
  - [4.4. Pages couvertes](#44-pages-couvertes)
- [5. Exigences non fonctionnelles](#5-exigences-non-fonctionnelles)
- [6. Améliorations futures](#6-améliorations-futures)
- [7. Conclusion](#7-conclusion)

</details>

---

## 1. Introduction

### 1.1. Présentation du projet

Forum Diapason est un forum web dédié à la musique, inspiré de Reddit. Le projet est réalisé dans le cadre de la formation B1 Informatique à Ynov Campus Lyon, dans des conditions professionnelles simulées.

L'objectif est de proposer une plateforme communautaire où les passionnés de musique peuvent créer un compte, partager leurs découvertes, avis et discussions musicales sous forme de posts, interagir avec les autres membres via des commentaires et des réactions, et naviguer sur un fil d'actualité. Le site doit être déployé et accessible en ligne à la date de la soutenance.

| | |
|--|--|
| **Dépôt GitHub** | [forum-diapason](https://github.com/Lyon-Ynov-Campus/forum-diapason) |
| **URL de déploiement** | *(à compléter)* |

---

### 1.2. Périmètre

**Ce que le projet couvre :**

- Création de compte, connexion et déconnexion
- Modification et suppression du profil utilisateur
- Publication, affichage et suppression de posts
- Système de commentaires sur les posts
- Système de likes sur les posts
- Consultation du profil d'autres utilisateurs
- Interface responsive, accessible sur desktop et mobile
- Déploiement continu de l'application

**Ce qui est hors périmètre :**

- Notifications en temps réel
- Messagerie privée entre utilisateurs
- Connexion via un service tiers (Google, GitHub, etc.)
- Système de modération ou de rôles administrateur
- Pièces jointes sur les posts (images, fichiers)
- Recherche par mot-clé 
- Page découverte 

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
| **Fil d'actualité** | Page principale affichant les posts par ordre chronologique décroissant |
| **JWT** | JSON Web Token — mécanisme d'authentification sécurisé utilisé pour maintenir la session côté client |
| **API REST** | Interface de programmation permettant la communication entre le frontend et le backend via des requêtes HTTP |
| **MVP** | Minimum Viable Product — version minimale fonctionnelle du produit |
| **Déploiement continu** | Mise en ligne automatique de l'application à chaque mise à jour du code sur la branche principale |
| **SQLite** | Base de données légère utilisée pour stocker les données de l'application |
| **Gorilla Mux** | Bibliothèque Golang utilisée pour gérer le routage de l'API |

---

## 2. Fonctionnalités

### 2.1. Comptes utilisateurs

#### Inscription

Un visiteur peut créer un compte en renseignant un identifiant unique et un mot de passe. Si l'identifiant est déjà utilisé, un message d'erreur lui est affiché. Le mot de passe est stocké de manière sécurisée (haché, jamais en clair). Une fois le compte créé, l'utilisateur est connecté automatiquement.

#### Connexion et déconnexion

L'utilisateur se connecte avec son identifiant et son mot de passe. Une session est créée via un token JWT et maintenue côté client. En cas d'identifiants incorrects, un message d'erreur s'affiche. L'utilisateur peut se déconnecter à tout moment, ce qui détruit la session.

#### Modification du profil

L'utilisateur connecté peut modifier son identifiant ou son mot de passe depuis sa page de paramètres. La modification nécessite de confirmer le mot de passe actuel.

#### Suppression du compte

L'utilisateur peut supprimer son compte de façon définitive. Une confirmation est demandée avant la suppression. Une fois le compte supprimé, l'utilisateur est redirigé vers la page d'accueil.

---

### 2.2. Posts

#### Création d'un post

Un utilisateur connecté peut créer un post en renseignant un titre et un contenu textuel. Les posts peuvent porter sur n'importe quel sujet musical : découvertes d'artistes, avis sur des albums, discussions sur des genres, partage de concerts, etc. Le post est horodaté et associé à son auteur. Un post avec un titre ou un contenu vide ne peut pas être publié.

#### Affichage des posts

Tous les visiteurs, connectés ou non, peuvent consulter les posts. La page d'accueil affiche un fil de posts du plus récent au plus ancien. Chaque post y affiche son titre, son auteur, sa date de publication, le nombre de likes et le nombre de commentaires.

#### Page d'un post

Un clic sur un post ouvre sa page dédiée, qui affiche le contenu complet, les commentaires et le bouton de réaction.

#### Suppression d'un post

L'auteur d'un post peut le supprimer. Un post supprimé disparaît du fil et de sa page.

---

### 2.3. Commentaires

Un utilisateur connecté peut commenter un post. Les commentaires sont affichés du plus ancien au plus récent. L'auteur d'un commentaire peut le supprimer. Un champ vide ne peut pas être soumis.

---

### 2.4. Likes

Un utilisateur connecté peut liker un post. Il ne peut liker un même post qu'une seule fois. Il peut retirer son like en cliquant à nouveau. Le compteur de likes est visible par tous les visiteurs.

---

### 2.5. Profil utilisateur

Chaque utilisateur dispose d'une page de profil consultable par les autres membres. Elle affiche l'identifiant, la date d'inscription et la liste des posts publiés. Sur son propre profil, l'utilisateur connecté voit également les options de modification et de suppression de compte.

---

### 2.6. Navigation

Toutes les pages partagent une barre de navigation et un pied de page communs. La navbar affiche le nom du forum, un lien vers le fil principal et, selon l'état de connexion, soit un lien vers le profil avec un bouton de déconnexion, soit un lien vers la page de connexion.

---

## 3. Cas d'utilisation

### CU-01 — Créer un compte

| | |
|--|--|
| **Acteur** | Visiteur non connecté |
| **Déclencheur** | Il clique sur "S'inscrire" |
| **Scénario** | Il remplit le formulaire avec un identifiant et un mot de passe, puis valide. Le système vérifie la disponibilité de l'identifiant, crée le compte et connecte l'utilisateur automatiquement. |
| **Échec** | Si l'identifiant est déjà pris, un message d'erreur s'affiche et le formulaire reste ouvert. |

---

### CU-02 — Se connecter

| | |
|--|--|
| **Acteur** | Visiteur non connecté |
| **Déclencheur** | Il clique sur "Se connecter" |
| **Scénario** | Il saisit son identifiant et son mot de passe. Le système vérifie les informations, génère un token JWT et redirige vers le fil principal. |
| **Échec** | Si les identifiants sont incorrects, un message d'erreur s'affiche. |

---

### CU-03 — Publier un post

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il clique sur "Nouveau post" |
| **Scénario** | Il saisit un titre et un contenu musical, puis valide. Le post est publié et apparaît en tête du fil. |
| **Échec** | Si le titre ou le contenu est vide, le formulaire affiche une erreur. |

---

### CU-04 — Commenter un post

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il ouvre un post et saisit un commentaire |
| **Scénario** | Il rédige son commentaire dans le champ prévu et valide. Le commentaire apparaît immédiatement sous le post. |
| **Échec** | Si le champ est vide, le formulaire affiche une erreur. |

---

### CU-05 — Liker un post

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il clique sur le bouton like d'un post |
| **Scénario** | Le compteur s'incrémente et le bouton indique visuellement que le like est actif. Un second clic retire le like et décrémente le compteur. |
| **Échec** | Si l'utilisateur n'est pas connecté, il est redirigé vers la page de connexion. |

---

### CU-06 — Modifier ou supprimer son compte

| | |
|--|--|
| **Acteur** | Utilisateur connecté |
| **Déclencheur** | Il accède à sa page de paramètres |
| **Scénario (modification)** | Il renseigne les nouvelles informations et confirme avec son mot de passe actuel. Les modifications sont enregistrées. |
| **Scénario (suppression)** | Il clique sur "Supprimer mon compte" et confirme l'action. Le compte est supprimé et il est redirigé vers l'accueil. |
| **Échec** | Si le mot de passe actuel est incorrect, un message d'erreur s'affiche. |

---

## 4. Maquettes et wireframes

### 4.1. Charte graphique

La charte graphique de Forum Diapason est en cours de définition par l'équipe design. Elle sera alignée sur la thématique musicale du forum.

| Élément | Statut |
|---------|--------|
| Palette de couleurs | *(à compléter)* |
| Typographie | *(à compléter)* |
| Logo | *(à compléter)* |

---

### 4.2. Wireframe

Le wireframe a été réalisé sur Excalidraw et couvre les pages principales du site.

*image*

Lien : https://excalidraw.com/#room=9090ab9a1cfda3075dfa,BKUC3vwf2gLmmUpQgZKVyQ

---

### 4.3. Maquette haute fidélité

La maquette est en cours de réalisation sur Figma par PHAM Huy et ANDRY Nomenafitia.

*image*

Lien Figma : *(à compléter)*

---

### 4.4. Pages couvertes

| Page | Description |
|------|-------------|
| Accueil | Fil de posts musicaux, navbar, footer, bouton de création de post |
| Connexion et inscription | Formulaires d'authentification |
| Page d'un post | Contenu complet, commentaires, bouton like |
| Profil utilisateur | Informations du compte, liste des posts publiés |
| Paramètres | Modification et suppression du compte |

---

## 5. Exigences non fonctionnelles

| Catégorie | Exigence |
|-----------|----------|
| Ergonomie | L'interface doit être intuitive et utilisable sans formation préalable |
| Responsivité | Le site doit fonctionner correctement sur desktop et mobile |
| Performance | Les pages doivent se charger en moins de 3 secondes en conditions normales |
| Sécurité | Les mots de passe sont hachés, l'authentification repose sur des tokens JWT |
| Déploiement | L'application doit être accessible depuis une URL publique le jour de la soutenance |
| Compatibilité | Le site doit fonctionner sur les navigateurs modernes (Chrome, Firefox, Edge) |

---

## 6. Améliorations futures

Ces fonctionnalités ne font pas partie du MVP mais pourraient être envisagées dans une version ultérieure du projet.

**Recherche par mot-clé**
Permettre aux utilisateurs de rechercher des posts par titre ou contenu, avec filtrage par popularité ou date.

**Page découverte**
Une page dédiée mettant en avant les posts les plus likés ou les discussions les plus actives, pour favoriser la découverte de contenu musical.

**Catégories et tags musicaux**
Associer des tags aux posts (genre musical, artiste, album, instrument) pour organiser le contenu et faciliter la navigation.

**Pièces jointes**
Permettre l'ajout d'images ou de liens dans les posts, par exemple pour partager des pochettes d'albums ou des clips.

**Système de notifications**
Notifier l'utilisateur lorsqu'un commentaire est posté sur l'un de ses posts ou lorsqu'il reçoit un like.

---

## 7. Conclusion

Forum Diapason vise à offrir un espace d'échange simple et agréable pour les passionnés de musique. En se concentrant sur les fonctionnalités essentielles — comptes, posts, commentaires et likes — le projet répond aux exigences du cahier des charges tout en restant réalisable dans le délai imparti de deux semaines.

Ce document constitue la référence fonctionnelle du projet pour l'ensemble de l'équipe. Il sera mis à jour au fur et à mesure de l'avancement, notamment pour intégrer les captures de la maquette Figma et les éléments de charte graphique une fois finalisés.

---

*Document rédigé par Michel LEVINE. Dernière révision : 08/05/2026.*
