![structure inkfem](https://github.com/user-attachments/assets/11e7b870-5d96-4cf3-a2a3-a6ab3a09e097)# inkFEM

An open-source, 2D structural analysis CLI that implements the Finite Element Method to calculate, generate and plot structures made of linear bars.

## Quick Tutorial

Structures are defined as plain-text files following the [_.inkfem_ input file format](./io/README.md).
Let's define a simple structure made of two columns, a beam, and a vertical distributed load over the beam like the following:

```
                                                           
                 qy = -100 N/cm                           
                ┌─────┬─────┬─────┬─────┐                  
                │     │     │     │     │                  
                │     ▼     ▼     ▼     │                  
 nodeC (0, 200) O───────────────────────O nodeD (300, 200) 
                │                       │                              
                │                       │                        
                │                       │                  
                │                       │                  
                │nodeA (0, 0)           │nodeB (300, 0)    
              ──O──                   ──O──                
```

Define it in a file called _structure.inkfem_ like so:

```
inkfem v1.1

|nodes|
nodeA -> 0.0    0.0    {dx dy rz}
nodeB -> 300.0  0.0    {dx dy rz}
nodeC -> 0.0    200.0  {}
nodeD -> 300.0  200.0  {}

|materials|
'steel' -> 1.0 20000000 1.0 1.0 25000 40000

|sections|
'ipe_120' -> 14 318 28 53 9

|loads|
fy ld beam 0.0 -100.0 1.00 -100.00

|bars|
# Columns
col1 -> nodeA{dx dy rz} nodeC{dx dy rz} 'steel' 'ipe_120'
col2 -> nodeB{dx dy rz} nodeD{dx dy rz} 'steel' 'ipe_120'
# Beam
beam -> nodeC{dx dy rz} nodeD{dx dy rz} 'steel' 'ipe_120'
```

> [!NOTE]
> To understand how to define structures using the _.inkfem_ file format, read the [specs here](./io/README.md).

You can plot the structure using the `plot` command:

```bash
$ inkfem plot path/to/structure.inkfem --scale 1.0 --dark
```


To solve the structure defined in the _structure.inkfem_ file:

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

## Build & Test

To build the `inkfem` binary:

```sh
$ make build
```

This creates the `inkfem` binary at the project's top level.
See the [Usage](#usage) section below to learn how to execute the binary program.

To run the tests:

```sh
$ make test
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
- [plot](): drawing SVG files from the `.inkfem`, `.inkfempre`and `.inkfemsol` files
- [cmd](): the commands available to the CLI
