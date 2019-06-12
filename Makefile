all:
	echo "Building for Linux x86_64..."
	CGO_ENABLED=0 go build
	upx -9 jimbot-go
clean:
	rm -f jimbot-go
