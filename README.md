# inkFEM

2D structural analysis program using the Finite Element Method.

## Usage

Structures are defined following the [_.inkfem_ input file format](./io/README.md).
To calculate a structure defined in a _.inkfem_ file:

```bash
$ inkfem solve -i path/to/structure.inkfem
```

Where the argument `-i` should point at the definition file.
If the structural analysis process doesn't encounter any error, it'll produce a solution file inside the same directory as the input file with the same name but _.inkfemsol_ extension.
In the example above, this would be _structure.inkfemsol_.

To also write the sliced (preprocessed) structure to a file, you can provide the `-p` flag:

```bash
$ inkfem solve -i path/to/structure.inkfem -p
```

This will generate an additional file with the _.inkfempre_ extension containing the information about how the structure has been sliced into finite elements.

### Available Flags

| Flag      | Type     | Description                                                            | Required | Default |
| --------- | -------- | ---------------------------------------------------------------------- | -------- | ------- |
| `-i`      | `string` | path to the input file                                                 | yes      | -       |
| `-v`      | `bool`   | use verbose output, including elapsed times                            | no       | `false` |
| `-p`      | `bool`   | save the preprocessed structure into a `.inkfempre` file               | no       | `false` |
| `-mat`    | `bool`   | save the system of equation's matrix as _.png_ image                   | no       | `false` |
| `-safe`   | `bool`   | perform some extra safety checks before proceeding with the resolution | no       | `false` |
| `-error`  | `float`  | maximum displacement error allowed in the resolution                   | no       | `1e-5`  |
| `-weight` | `bool`   | include the own weight of the bars                                     | no       | `false` |

## Build & Test

To build the `inkfem` binary:

```sh
$ ./build.sh
```

or alternatively:

```sh
$ go build inkfem.go
```

This creates the `inkfem` binary at the project's top level.
See the [Usage](#usage) section below to learn how to execute the binary program.

To run the tests:

```sh
$ ./test.sh
```

or alternatively:

```sh
$ go test ./...
```

## Docs

- [Go Modules](https://go.dev/doc/modules/managing-dependencies)
- [Cobra CLI](https://github.com/spf13/cobra)

## Code Structure

The code is split into four main packages:

- [structure](./structure/README.md): defines the structure model
- [preprocess](./preprocess/README.md): implements the preprocessing or slicing of the structure
- [process](./process/README.md): implements the processing of a sliced/preprocessed structure
- [io](./io/README.md): reading from `.inkfem` files and writing to `.inkfempre`and `.inkfemsol` files
- [cmd](): the commands available to the CLI
