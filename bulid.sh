CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o convertBase64ForLinux main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o convertBase64ForMacos main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o convertBase64ForM1 main.go