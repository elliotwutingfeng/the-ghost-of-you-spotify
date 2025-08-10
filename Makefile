build:
	go build -o ./dist/theghostofyouspotify ./main.go
	chmod +x ./dist/theghostofyouspotify

setup:
	cp --update=none .env.txt .env
