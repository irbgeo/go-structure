# Struct Tag Builder

## Description

A package for creating and modifying a structure in runtime

This package allows:

- to build structure.
- to add tags to structure.
- to merge two structure.

## Usage

```golang
import(
    ...
    structure "github.com/irbgeo/go-structure"
    "external/pkg"
    ...
)

...
func getTag(fieldName string) string {
	return strings.ToLower(fieldName)
}
...
b, _ := structure.New(new(pkg.Structure))S
b.AddTags(getTag)

err = yaml.Unmarshal([]byte("content for pkg.Structure"), b.Struct())

var dst pkg.Structure
err = b.SaveInto(&dst)
...

```

Other examples in tests
