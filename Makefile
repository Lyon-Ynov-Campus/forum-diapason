TAILWIND_CONTENT = "./frontend/**/*.html,./frontend/**/*.js"
GOPATH_BIN       = $(shell go env GOPATH)/bin

dev:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --watch --content $(TAILWIND_CONTENT) &
	@go run ./api/ &
	@$(GOPATH_BIN)/air

dev-api:
	@go run ./api/

build:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --minify --content $(TAILWIND_CONTENT)
	@go build -o forum-diapason .
	@go build -o forum-api ./api/
