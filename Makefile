all:
	CGO_ENABLED=0 go build
	upx -9 jimbot-go
clean:
	rm -f jimbot-go
