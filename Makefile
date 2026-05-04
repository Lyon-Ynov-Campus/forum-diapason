dev:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --watch --content "./frontend/**/*.html" &
	@go run .

build:
	@./tailwindcss -i ./frontend/css/input.css -o ./frontend/css/styles.css --minify --content "./frontend/**/*.html"
	@go build -o forum-diapason .
