TAILWIND_CONTENT = "./frontend/**/*.html,./frontend/**/*.js"

dev:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --watch --content $(TAILWIND_CONTENT) &
	@$(shell go env GOPATH)/bin/air

build:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --minify --content $(TAILWIND_CONTENT)
	@go build -o forum-diapason .
