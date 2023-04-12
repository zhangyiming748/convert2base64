CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o forLinux main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o forMacos main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o forM1 main.go