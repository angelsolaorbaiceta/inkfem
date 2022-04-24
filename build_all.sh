# Builds the CLI for all operating systems

rm -rf bin/

build_for_os() {
  echo Building for $1 OS, with $2 arch
  
  if [ $1 = "windows" ]; then
    FILE_PATH=bin/inkfem.exe
  else
    FILE_PATH=bin/inkfem
  fi

  SHASUM_PATH=bin/$1_$2_sha256sum.txt
  ZIP_PATH=bin/$1_$2.gz

  GOOS=$1 GOARCH=$2 go build -ldflags "-s -w" -o $FILE_PATH inkfem.go
  shasum -a 256 $FILE_PATH > $SHASUM_PATH
  gzip -9c $FILE_PATH > $ZIP_PATH
  rm $FILE_PATH
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