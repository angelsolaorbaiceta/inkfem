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

## Code Structure

The code is split into four main packages:

- [structure](./structure/README.md): defines the structure model
- [preprocess](./preprocess/README.md): implements the preprocessing or slicing of the structure
- [process](./process/README.md): implements the processing of a sliced/preprocessed structure
- [io](./io/README.md): reading from `.inkfem` files and writing to `.inkfempre`and `.inkfemsol` files
