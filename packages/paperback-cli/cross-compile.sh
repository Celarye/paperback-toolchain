env GOOS=linux GOARCH=amd64 go build -o ./bin/paperback-cli_linux_amd64 -ldflags "-w"
env GOOS=windows GOARCH=amd64 go build -o ./bin/paperback-cli_windows_amd64 -ldflags "-w"
env GOOS=darwin GOARCH=amd64 go build -o ./bin/paperback-cli_darwin_amd64 -ldflags "-w"
env GOOS=darwin GOARCH=arm64 go build -o ./bin/paperback-cli_darwin_arm64 -ldflags "-w"