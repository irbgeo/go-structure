# Struct Tag Builder

## Description

This package allows:

- to add tags to an external structure.
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

var src pkg.Structure
err = b.SaveInto(&actualStruct)


...

```
