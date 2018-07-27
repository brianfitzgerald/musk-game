build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/createRoom createRoom/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dealCard dealCard/main.go
