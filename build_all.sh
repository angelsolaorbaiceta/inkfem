# Builds the CLI for all operating systems
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/inkfem_linux_amd64 inkfem.go
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/inkfem_linux_arm64 inkfem.go

GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/inkfem_darwin_amd64 inkfem.go
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/inkfem_darwin_arm64 inkfem.go

GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/inkfem_windows_amd64.exe inkfem.go