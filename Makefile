TAILWIND_CONTENT = "./frontend/**/*.html,./frontend/**/*.js"

dev:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --watch --content $(TAILWIND_CONTENT) & \
	TW_PID=$$!; \
	trap "kill $$TW_PID 2>/dev/null" EXIT INT TERM; \
	$(shell go env GOPATH)/bin/air

build:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --minify --content $(TAILWIND_CONTENT)
	@go build -o forum-diapason .
