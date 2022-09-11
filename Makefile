linux:
	GOOS=linux GOARCH=amd64 go build -o bin/wagumi-sbt-client_linux-amd64 .
win:
	GOOS=windows GOARCH=amd64 go build -o bin/wagumi-sbt-client_windows-amd64.exe .