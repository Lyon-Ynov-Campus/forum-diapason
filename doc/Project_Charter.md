# Forum Diapason | Charte de projet

**École :** Ynov Campus Lyon  
**Promotion :** B1 Informatique 2025/2026  
**Durée du projet :** 27/04/2026 au 17/05/2026  
**Dernière mise à jour :** 17/05/2026  
**Rédigé par :** Michel LEVINE

---

## 1. Définition du projet

Forum Diapason est un forum web sur la thématique musicale, inspiré de Reddit. L'idée est de proposer un espace où les utilisateurs peuvent créer un compte, publier des posts, commenter et réagir aux contenus des autres membres.

Le site est accessible en ligne à l'adresse https://forum.prettyflacko.fr.

---

## 2. Objectif

### Ce que le projet inclut

- Un système de comptes avec inscription, connexion sécurisée par session cookie, modification et suppression de profil
- Réinitialisation de mot de passe par email (SMTP)
- Un système de posts avec titre, contenu, tags et image (recadrage + compression)
- Un système de réactions (likes) et de commentaires avec réponses imbriquées
- Une recherche fulltext avec filtres de tri
- Un dark mode persistant
- Une API RESTful en Golang (stdlib `net/http`) connectée à une base SQLite
- Le déploiement de l'application sur Azure Container Instances avec Cloudflare
- La documentation de gestion du projet (charte et rapports hebdomadaires)
- Une spécification fonctionnelle

### Ce que le projet n'inclut pas

- Les notifications en temps réel
- La messagerie privée entre utilisateurs
- La connexion via un service tiers comme Google ou GitHub
- Un système de modération ou de rôles administrateur

---

## 3. Parties prenantes

| Rôle | Nom | Attentes |
|------|-----|----------|
| Client | Mentor Kilian MOUN | Un forum fonctionnel livré dans les délais, avec une gestion de projet suivie sur Jira |
| Utilisateurs finaux | Grand public | Un site ergonomique, rapide et agréable à utiliser depuis un navigateur |

---

## 4. Équipe et responsabilités

| Nom | Domaines principaux | Critères de performance |
|-----|---------------------|------------------------|
| Joaquin COZETTE | Gestion de projet, suivi Jira, wireframe | Tâches à jour sur Jira, wireframe livré avant le 30/04 |
| Quochuy PHAM | Design (maquette Figma), développement front-end | Maquette finalisée en semaine 1, pages conformes au design |
| Michel LEVINE | Développement back-end, API Golang, base de données SQLite, déploiement, documentation | API fonctionnelle, sessions opérationnelles, site déployé, documents rendus |
| Nomenafitia ANDRY | Développement front-end, CSS, documentation utilisateur | Pages implémentées, design responsive, documentation rédigée |

---

## 5. Organisation du projet

Le projet est organisé en deux sprints d'une semaine, alignés sur les deux Weeklies notés. On utilise une approche Agile avec un daily de 15 minutes chaque matin pour faire le point. Le mentor joue le rôle de client et participe aux deux revues hebdomadaires.

Les tâches sont toutes suivies dans Jira sur le Tableau Sprint 1.

**Sprint 1 (27/04 au 07/05) :** Mise en place du projet, design, wireframes, maquettes et premières briques de développement.

**Sprint 2 (07/05 au 17/05) :** Authentification, gestion du profil, posts, commentaires, likes, recherche, déploiement et documentation.

---

## 6. Dates clés

| Date | Échéance | Livrables attendus |
|------|----------|--------------------|
| 27/04/2026 | Lancement | Configuration Jira, diagramme de Gantt, dépôt GitHub |
| 30/04/2026 | Point de planification | Découpage des tâches, wireframe Excalidraw, début de la maquette Figma |
| 07/05/2026 | Weekly 1 | Authentification, modification et suppression de profil, page d'accueil avec header, footer et template |
| 17/05/2026 | Soutenance | Forum complet déployé, présentation technique au client |

---

## 7. Livrables

**Techniques :**
- Code source hébergé sur GitHub : [forum-diapason](https://github.com/Lyon-Ynov-Campus/forum-diapason)
- Application web déployée : https://forum.prettyflacko.fr
- Base de données SQLite persistante (Azure File Share)
- Images Docker publiées sur Azure Container Registry

**Documentation :**
- Charte de projet
- Rapports hebdomadaires x2
- Spécification fonctionnelle

---

## 8. Ressources

| Ressource | Détail |
|-----------|--------|
| Équipe | 4 membres |
| Durée | 15 jours ouvrés |
| Outils | GitHub Classroom, Jira, Excalidraw, Figma, VS Code, Docker, Azure |
| Stack technique | Golang stdlib, HTML/CSS/JS vanilla, Tailwind CSS v4, SQLite, sessions cookie |

---

## 9. Risques identifiés

| Risque | Probabilité | Impact | Statut |
|--------|-------------|--------|--------|
| Délai très serré pour un projet full-stack | Élevée | Élevé | Géré — MVP livré dans les délais |
| Complexité du back-end (Golang, SQLite, architecture 2 serveurs) | Moyenne | Élevé | Résolu |
| Problèmes de déploiement en dernière minute | Moyenne | Élevé | Résolu — déployé sur Azure ACI + Cloudflare |
| Ajout de fonctionnalités non prioritaires avant que le cœur soit terminé | Moyenne | Moyen | Géré — priorisation MVP respectée |

---

## 10. Définition du "terminé"

On considère qu'une fonctionnalité est terminée quand :
1. Elle est implémentée et testée manuellement
2. Elle est fusionnée dans la branche principale sur GitHub
3. Le ticket Jira correspondant est passé en "Terminé"
4. Elle est visible et fonctionnelle sur la version déployée

---

*Document rédigé par Michel LEVINE. Dernière révision : 17/05/2026.*
