SHELL := /bin/bash

PYTHON := $(shell command -v python3 2>/dev/null || command -v python 2>/dev/null)

ifndef PYTHON
$(error Python is required to run the frontend static server)
endif

.PHONY: run-backend run-frontend run-all

run-backend:
	@echo "Lancement du backend sur http://localhost:8080"
	@go run main.go

run-frontend:
	@echo "Lancement du frontend sur http://localhost:3000"
	@$(PYTHON) -m http.server 3000 --directory frontend

run-all:
	@echo "Lancement du backend et du frontend..."
	@$(MAKE) run-backend &\
	 sleep 1; \
	 $(MAKE) run-frontend
