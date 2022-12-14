build:
	set GOOS=linux&& go build -o main main.go
	set GOARCH=amd64
	set CGO_ENABLED=0
	zip -o main.zip main
	sam deploy