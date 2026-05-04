# Forum Diapason | Charte de projet

**École :** Ynov Campus Lyon  
**Promotion :** B1 Informatique 2025/2026  
**Durée du projet :** 27/04/2026 au 11/05/2026  
**Dernière mise à jour :** 04/05/2026  
**Rédigé par :** Michel LEVINE

---

## 1. Définition du projet

Forum Diapason est un forum web sur la thématique musicale, inspiré de Reddit. L'idée est de proposer un espace où les utilisateurs peuvent créer un compte, publier des posts, commenter et réagir aux contenus des autres membres.

Le site doit être accessible au grand public et entièrement déployé pour la soutenance du 11 mai 2026.

---

## 2. Objectif

### Ce que le projet inclut

- Un système de comptes avec inscription, connexion sécurisée via JWT, modification et suppression de profil
- Un système de posts avec création, affichage et fil d'actualité
- Un système de réactions (likes) et de commentaires sur les posts
- Une API RESTful en Golang (Gorilla Mux) qui fait le lien entre la base de données SQLite et le frontend
- Le déploiement continu de l'application
- La documentation de gestion du projet (charte et rapports hebdomadaires)
- Une spécification fonctionnelle et une spécification technique *(optionnelles, allégées à la demande du client)*

### Ce que le projet n'inclut pas

- Les notifications en temps réel
- La messagerie privée entre utilisateurs
- La connexion via un service tiers comme Google ou GitHub

---

## 3. Parties prenantes

| Rôle | Nom | Attentes |
|------|-----|----------|
| Client | Mentor *(à compléter)* | Un forum fonctionnel livré dans les délais, avec une gestion de projet suivie sur Jira |
| Utilisateurs finaux | Grand public | Un site ergonomique, rapide et agréable à utiliser depuis un navigateur |

---

## 4. Équipe et responsabilités

| Nom | Domaines principaux | Critères de performance |
|-----|---------------------|------------------------|
| Joaquin Cozette | Gestion de projet, suivi Jira, wireframe | Tâches à jour sur Jira, wireframe livré avant le 30/04 |
| PHAM Huy | Design (maquette Figma), développement front-end | Maquette finalisée en semaine 1, pages conformes au design |
| Michel LEVINE | Développement back-end, API Golang, base de données SQLite, déploiement, documentation | API fonctionnelle, JWT opérationnel, site déployé, documents rendus |
| ANDRY Nomenafitia | Développement front-end, CSS, documentation utilisateur | Pages implémentées, design responsive, documentation rédigée |

---

## 5. Organisation du projet

Le projet est organisé en deux sprints d'une semaine, alignés sur les deux Weeklies notés. On utilise une approche Agile avec un daily de 15 minutes chaque matin pour faire le point. Le mentor joue le rôle de client et participe aux deux revues hebdomadaires.

Les tâches sont toutes suivies dans Jira sur le Tableau Sprint 1.

**Sprint 1 (27/04 au 07/05) :** Mise en place du projet, design, wireframes, maquettes et premières briques de développement.

**Sprint 2 (07/05 au 11/05) :** Authentification JWT, gestion du profil, posts, commentaires, likes, déploiement et documentation.

---

## 6. Dates clés

| Date | Échéance | Livrables attendus |
|------|----------|--------------------|
| 27/04/2026 | Lancement | Configuration Jira, diagramme de Gantt, dépôt GitHub |
| 30/04/2026 | Point de planification | Découpage des tâches, wireframe Excalidraw, début de la maquette Figma |
| 07/05/2026 | Weekly 1 | Authentification, modification et suppression de profil, page d'accueil avec header, footer et template |
| 11/05/2026 | Weekly 2 et Soutenance | Forum complet déployé, présentation technique au client |

---

## 7. Livrables

**Techniques :**
- Code source hébergé sur GitHub : [forum-diapason](https://github.com/Lyon-Ynov-Campus/forum-diapason)
- Application web déployée *(URL à compléter)*
- Base de données SQLite

**Documentation :**
- Charte de projet
- Rapports hebdomadaires x2
- Spécification fonctionnelle *(optionnelle)*
- Spécification technique *(optionnelle)*
- Documentation utilisateur (FAQ)

---

## 8. Ressources

| Ressource | Détail |
|-----------|--------|
| Budget | 0 € |
| Équipe | 4 membres |
| Durée | 11 jours ouvrés |
| Outils | GitHub Classroom, Jira, Excalidraw, Figma, VS Code |
| Stack technique | Golang avec Gorilla Mux, HTML/CSS/JS vanilla, SQLite, JWT |

---

## 9. Risques identifiés

| Risque | Probabilité | Impact | Mitigation |
|--------|-------------|--------|------------|
| Délai très serré pour un projet full-stack | Élevée | Élevé | Prioriser le MVP et reporter les fonctionnalités secondaires si nécessaire |
| Complexité du back-end (Golang, SQLite, JWT, Clean Architecture) | Moyenne | Élevé | S'appuyer sur les ressources fournies et démarrer le back-end tôt |
| Problèmes de déploiement en dernière minute | Moyenne | Élevé | Tester le déploiement dès le début du Sprint 2 |
| Ajout de fonctionnalités non prioritaires avant que le cœur soit terminé | Moyenne | Moyen | Traiter les fonctionnalités principales d'abord, les extras seulement si le temps le permet |

---

## 10. Définition du "terminé"

On considère qu'une fonctionnalité est terminée quand :
1. Elle est implémentée et testée manuellement
2. Elle est fusionnée dans la branche principale sur GitHub
3. Le ticket Jira correspondant est passé en "Terminé"
4. Elle est visible et fonctionnelle sur la version déployée

---

*Document rédigé par Michel LEVINE. Dernière révision : 07/05/2026.*
