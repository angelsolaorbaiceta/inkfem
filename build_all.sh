# Builds the CLI for all operating systems
rm -rf bin/

build_for_os() {
  echo Building for $1 OS, with $2 arch
  
  FILE_PATH=bin/inkfem_$1_$2
  GOOS=$1 GOARCH=$2 go build -ldflags "-s -w" -o $FILE_PATH inkfem.go
  shasum -a 256 $FILE_PATH >> bin/sha256sums.txt
  gzip -9 $FILE_PATH
}

# Linux
build_for_os linux amd64
build_for_os linux arm64

# OSX
build_for_os darwin amd64
build_for_os darwin arm64

# Windows
build_for_os windows amd64

echo Done!