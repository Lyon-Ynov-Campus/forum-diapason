.PHONY: help run-all

help:
	@echo "Forum Diapason root Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  make run-all   Run backend and frontend together"

run-all:
	@echo "Starting backend and frontend..."
	@bash -lc 'set -e; ROOT="$$(pwd)"; cd "$$ROOT/backend"; go run ./cmd/api & BACKEND=$$!; cd "$$ROOT/frontend"; npm run dev & FRONTEND=$$!; trap "kill $$BACKEND $$FRONTEND" INT TERM; wait $$BACKEND $$FRONTEND'
