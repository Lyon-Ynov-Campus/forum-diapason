# Forum Diapason — Déploiement Azure

Forum musical déployé sur Azure Container Instances avec Cloudflare.

**URL de production** : https://forum.prettyflacko.fr

> Pour le développement local, voir la branche `main`.

---

## Stack

- **Backend** : Go stdlib (`net/http`, `html/template`)
- **Base de données** : SQLite via `go-sqlite3` (CGO)
- **Frontend** : HTML, JavaScript vanilla, Tailwind CSS v4
- **Auth** : Sessions cookie (`SameSite=Lax`, `HttpOnly`)
- **Email** : SMTP (reset de mot de passe)
- **Déploiement** : Azure Container Instances + Azure File Share + Cloudflare

---

## Architecture de déploiement

```
Cloudflare (HTTPS + Origin Rules)
       │
       ├── forum.prettyflacko.fr → ACI forum-front :8080
       └── api.prettyflacko.fr   → ACI forum-api   :8081
                                          │
                                 Azure File Share
                                 ├── forumdata   → /app/data   (SQLite)
                                 └── forumpublic → /app/public (uploads)
```

Les 2 containers partagent la même base SQLite via le File Share `forumdata`.

---

## Prérequis Azure

- Azure CLI (`az`)
- Docker
- Azure Container Registry : `forumdiapason.azurecr.io`
- Azure Resource Group : `forum-diapason`
- Azure Storage Account : `forumdiapasonstorage`
  - File Share `forumdata` (base de données)
  - File Share `forumpublic` (uploads avatars et images)

---

## Build et push des images

```bash
az acr login --name forumdiapason

docker build --no-cache -f Dockerfile.front -t forumdiapason.azurecr.io/forum-front:latest .
docker build --no-cache -f Dockerfile.api   -t forumdiapason.azurecr.io/forum-api:latest .

docker push forumdiapason.azurecr.io/forum-front:latest
docker push forumdiapason.azurecr.io/forum-api:latest
```

---

## Déploiement des containers

Partir des templates `aci-front.example.yaml` et `aci-api.example.yaml`, remplacer les `PLACEHOLDER_*` par les vraies valeurs (credentials ACR + clé Storage), puis :

```bash
az container delete -g forum-diapason -n forum-front --yes
az container delete -g forum-diapason -n forum-api --yes

az container create -g forum-diapason --file aci-front.yaml
az container create -g forum-diapason --file aci-api.yaml
```

---

## Variables d'environnement

### forum-front

| Variable | Valeur prod |
|----------|-------------|
| `PORT` | `8080` |
| `DB_FILE` | `/app/data/forum.db` |
| `COOKIE_SECURE` | `false` |
| `COOKIE_DOMAIN` | `.prettyflacko.fr` |

### forum-api

| Variable | Valeur prod |
|----------|-------------|
| `API_PORT` | `8081` |
| `DB_FILE` | `/app/data/forum.db` |
| `FRONTEND_ORIGIN` | `https://forum.prettyflacko.fr` |
| `COOKIE_SECURE` | `false` |
| `COOKIE_DOMAIN` | `.prettyflacko.fr` |
| `SMTP_HOST` | `smtp.gmail.com` |
| `SMTP_PORT` | `587` |
| `SMTP_USER` | *(compte Gmail)* |
| `SMTP_PASS` | *(app password 16 car)* |
| `SMTP_FROM` | `contact@prettyflacko.fr` |

---

## Configuration Cloudflare

**DNS** (proxied 🟠) :
| Type | Name | Cible |
|------|------|-------|
| CNAME | `forum` | `forum-diapason.francecentral.azurecontainer.io` |
| CNAME | `api` | `forum-diapason-api.francecentral.azurecontainer.io` |

**SSL/TLS** : mode `Flexible`

**Origin Rules** :
| Règle | Condition | Action |
|-------|-----------|--------|
| `front-port` | Hostname = `forum.prettyflacko.fr` | Port → `8080` |
| `api-port` | Hostname = `api.prettyflacko.fr` | Port → `8081` |

---

## Vérification

```bash
# Santé de l'API
curl -i https://api.prettyflacko.fr/api/auth/me
# → HTTP/2 401 {"error":"non authentifié"}

# Logs des containers
az container logs -g forum-diapason -n forum-front
az container logs -g forum-diapason -n forum-api
```
