linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/wagumi-sbt-client_linux-amd64 .
win:
	GOOS=windows GOARCH=amd64 go build -o bin/wagumi-sbt-client_windows-amd64.exe .
mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/wagumi-sbt-client_mac-amd64 .
mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/wagumi-sbt-client_mac-arm64 .
