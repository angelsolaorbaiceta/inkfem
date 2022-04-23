# Usage

Structures are defined following the _.inkfem_ input file format.
To calculate a structure defined in a _.inkfem_ file:

```bash
$ inkfem solve path/to/structure.inkfem
```

If the structural analysis process doesn't encounter any error, it'll produce a solution file inside the same directory as the input file with the same name but _.inkfemsol_ extension.
In the example above, this would be _structure.inkfemsol_.

To also write the sliced (preprocessed) structure to a file, you can provide the `-p` flag:

```bash
$ inkfem solve path/to/structure.inkfem -p
```

This will generate an additional file with the _.inkfempre_ extension containing the information about how the structure has been sliced into finite elements.

### Available Flags

| Flag                 | Type    | Description                                                            | Required | Default |
| -------------------- | ------- | ---------------------------------------------------------------------- | -------- | ------- |
| `verbose` or `-v`    | `bool`  | use verbose output, including elapsed times                            | no       | `false` |
| `preprocess` or `-p` | `bool`  | save the preprocessed structure into a `.inkfempre` file               | no       | `false` |
| `safe` or `-s`       | `bool`  | perform some extra safety checks before proceeding with the resolution | no       | `false` |
| `error` or `-e`      | `float` | maximum displacement error allowed in the resolution                   | no       | `1e-5`  |
| `weight` or `-w`     | `bool`  | include the own weight of the bars                                     | no       | `false` |
