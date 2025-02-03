all:
	echo "Building for Linux x86_64..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o jimbot-go cmd/jimbot/main.go
clean:
	rm -f jimbot-go
