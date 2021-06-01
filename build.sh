# builds the release version of the inkfem binary: strips debug and symbol info.
go build -ldflags "-s -w" inkfem.go
