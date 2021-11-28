# inkFEM

Package for 2D structural analysis using the Finite Element Method.

## Usage

Structures are defined following the [`.inkfem` input file format](./io/README.md).
To calculate a structure defined in a `.inkfem` file:

```bash
$ inkfem -i path/to/structure.inkfem
```

Where the argument `-i` should point at the definition file.
If the structural analysis process doesn't encounter any error, it'll produce a solution file inside the same directory as the input file with the same name but `.inkfemsol` extension.
In the example above, this would be _structure.inkfemsol_.

To also write the sliced (preprocessed) structure to a file, you can provide the `-p` flag:

```bash
$ inkfem -i path/to/structure.inkfem -p
```

This will generate an additional file with the `.inkfempre` extension containing the information about how the structure has been sliced into finite elements.

### Available Flags

- `-i (string)`: input file path
- `-v (bool)`: uses verbose output, including the time taken for each operation
- `-p (bool)`: to save the preprocessed structure into a `.inkfempre` file
- `-mat (bool)`: to save the system of equation's matrix as png image
- `-safe (bool)`: to perform some extra safety checks before proceeding with the resolution
- `-error (float64)`: to choose the maximum displacement error allowed in the resolution
- `-weight (bool)`: to include the own weight of the bars

## Build

To build the binary, execute the `build.sh` script:

```bash
$ ./build.sh
```

This will produce the `inkfem` binary in the project's root directory.

## Test

To run all the tests, execute the `test.sh` script:

```bash
$ ./test.sh
```

You can also run the tests inside a particular package like so:

```bash
$ go test ./process
```

## Docs

- [Go Modules](https://go.dev/doc/modules/managing-dependencies)

## Code Structure

The code is split into four main packages:

- [structure](./structure/README.md): defines the structure model
- [preprocess](./preprocess/README.md): implements the preprocessing or slicing of the structure
- [process](./process/README.md): implements the processing of a sliced/preprocessed structure
- [io](./io/README.md): reading from `.inkfem` files and writing to `.inkfempre`and `.inkfemsol` files
