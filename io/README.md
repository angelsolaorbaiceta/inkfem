# Input File Format

The structure input file defines the structure to be analyzed.
The calculation engine is unit-agnostic, so the units used as input are also the output units.
Units need, nevertheless, to be congruent.

The input file should be a plain-text file with the `.inkfem` extension.
Inside the file, the first line should include the header:

```
inkfem vM.m
```

where `M` and `m` are the major and minor versions of the binary used to compute the structure.
The current version is `1.1`:

```
inkfem v1.1
```

Then go the definition sections:

- `nodes`: the structure nodes, referred by id
- `sections`: the element's sections, referred by name
- `materials`: the element's materials, referred by name
- `loads`: the loads applied to the nodes and elements
- `bars`: the structure bars (linear resistant elements), referred by id

The sections can appear in any order.

## The Nodes

The nodes are defined under the header:

```
|nodes|
```

Each node is defined following the format:

```
<id> -> <xCoord> <yCoord> {[dx dy rz]}
```

where:

- _id_: the node's unique id
- _xCoord_: the node's position x-coordinate
- _yCoord_: the node's position y-coordinate
- _{dx dy rz}_: set of externally constrained degrees of freedom

### Examples

Node with id 23, at position `(120, 450)` and no external constraints:

```
23 -> 120.0 450.0 { }
```

Node with id 48, at position `(300, 50)` and the displacement in the x and y directions externally constrained:

```
48 -> 300.0 50.0 { dx dy }
```

## The Materials

The materials are defined under the header:

```
|materials|
```

Each material is defined following the format:

```
<name> -> <density> <young> <shear> <poisson> <yield> <ultimate>
```

where:

- _name_: the material's unique name
- _density_: the material's density
- _young_: the material's Young or elasticity modulus
- _shear_: the material's shear modulus
- _poisson_: the material's poisson ratio
- _yield_: the material's yield strength
- _ultimate_: the material's ultimate strength

### Examples

Standard steel 275:

```
'steel_275' -> 0.00000785 21000000.0 8100000.0 0.3 27500.0 43000.0
```

## The Sections

The sections are defined under the header:

```
|sections|
```

Each section is defined following the format:

```
<name> -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
```

where:

- _name_: the section's unique name
- _area_: the section's cross section area
- _iStrong_: the strong axis' moment of inertia
- _iWeak_: the weak axis' moment of inertia
- _sStrong_: the strong axis' section modulus
- _sWeak_: the weak axis' section modulus

### Examples

Standard European IPE-100 section:

```
'ipe_100' -> 10.3 171.0 15.92 34.2 5.79
```

## The Loads

The loads are defined under the header:

```
|loads|
```

There are two types of loads:

- Distributed
- Concentrated

Both types are applied to bars.
To define a concentrated load on a node, choose a bar which contains the node, and add the concentrated load to that bar at position `t = 0` (start node) or `t = 1` (end node).

**Distributed** loads are defined following the format:

```
<term> <reference><type> <barId> <tStart> <valueStart> <tEnd> <valueEnd>
```

where:

- _term_: is either:
  - `fx`: force in the x-axis direction
  - `fy`: force in the y-axis direction
  - `mz`: moment about the z-axis
- _reference_: the reference frame in which the load is defined. Can be:
  - `l`: reference frame **local** to the bar
  - `g`: **global** reference frame
- _type_: must be `d` to signify this is a distributed load
- _barId_: The id of the bar where the load is applied
- _tStart_: the load start position in the bar's directrix (`0 <= t <= 1`)
- _valueStart_: the value for the load at `tStart`
- _tEnd_: the load end position in the bar's directrix (`tStart <= t <= 1`)
- _valueEnd_: the value for the load at `tEnd`

Distributed loads are always linear: they have a start and end value, and those values are linearly interpolated.
The current implementation doesn't allow any other kind of distributed load interpolation.

### Examples (Distributed)

A distributed force in the bar's local y-axis direction, applied to a bar with id 4, starting at `t = 0` with value `-50` and ending at `t = 1` with value `-75`.

```
fy ld 4 0.0 -50.0 1.0 -75.0
```

A distributed moment about the global z-axis, applied to a bar with id 12, starting at `t = 0.25` with value `100` and ending at `t = 0.75` with value `200`.

```
mz gd 12 0.25 100 0.75 200
```

**Concentrated** loads are defined following the format:

```
<term> <reference><type> <elementId> <t> <value>
```

where:

- _term_: is either:
  - `fx`: force in the x-axis direction
  - `fy`: force in the y-axis direction
  - `mz`: moment about the z-axis
- _reference_: the reference frame in which the load is defined. Can be:
  - `l`: reference frame **local** to the bar
  - `g`: **global** reference frame
- _type_: must be `c` to signify this is a concentrated load
- _elementId_: The if of the bar where the load is applied
- _t_: the load position in the bar's directrix (`0 <= t <= 1`)
- _value_: the load's value

### Examples (Concentrated)

A concentrated force in the element with id 11 local y-axis direction, at position `t = 0` (applied in the start node), with value `-70`.

```
fy lc 11 0.0 -70.0
```

## The Bars

The bars are defined under the header:

```
|bars|
```

Each bar is defined following the format:

## Input File Example

Here's a complete input file example:

```
inkfem v1.1

|nodes|
1 -> 0 0 {dx dy}
2 -> 200 300 {}
3 -> 400 0 {}
4 -> 600 300 {}
5 -> 800 0 {dx dy}

|materials|
'mat_A' -> 1.0 1.0 1.0 1.0 1.0 1.0

|sections|
'sec_A' -> 1.0 1.0 1.0 1.0 1.0

|loads|
fy ld 4 0.0 -50.0 1.0 -75.0

|bars|
1 -> 1{dx dy rz} 2{dx dy rz} 'mat_A' 'sec_A'
2 -> 1{dx dy rz} 3{dx dy rz} 'mat_A' 'sec_A'
3 -> 2{dx dy rz} 3{dx dy rz} 'mat_A' 'sec_A'
4 -> 2{dx dy rz} 4{dx dy rz} 'mat_A' 'sec_A'
5 -> 3{dx dy rz} 4{dx dy rz} 'mat_A' 'sec_A'
6 -> 3{dx dy rz} 5{dx dy rz} 'mat_A' 'sec_A'
7 -> 4{dx dy rz} 5{dx dy rz} 'mat_A' 'sec_A'
```

# Preprocess File Format

The preprocessed structure is saved into a `.inkfempre` file if the `-p` flag is passed to inkfem.
The file's template is defined in [preprocess.template.txt](./templates/preprocess.template.txt).

# Solution File Format

The solution structure is saved into a `.inkfemsol` file.
The file's template is defined in [solution.template.txt](./templates/solution.template.txt).
