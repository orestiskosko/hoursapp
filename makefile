templ:
	@templ generate --watch

run:
	@air -c .air.toml

tailwind:
	@npx tailwindcss -i ./input.css -o ./public/output.css --watch
