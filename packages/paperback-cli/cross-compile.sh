env GOOS=linux GOARCH=amd64 go build -o ../paperback-cli-linux-amd64/bin/paperback-cli -ldflags "-w"
env GOOS=windows GOARCH=amd64 go build -o ../paperback-cli-windows-amd64/bin/paperback-cli.exe -ldflags "-w"
env GOOS=darwin GOARCH=amd64 go build -o ../paperback-cli-darwin-amd64/bin/paperback-cli -ldflags "-w"
env GOOS=darwin GOARCH=arm64 go build -o ../paperback-cli-darwin-arm64/bin/paperback-cli -ldflags "-w"