# Go Structure

## Description

The package for creating and modifying a structure in runtime

This package allows:

- to change tags to fields of external structure.
- to build structure.
- to merge two structure.
- to fill structure from another structure or map
- to fill another structure or map from structure

## Import

```golang
structure "github.com/irbgeo/go-structure"
```

## Usage

### Work with external structure

1. Create structure interface from external structure

```golang
    b, _ := structure.New(new(pkg.Structure))
```

2. Change tags in structure

```golang
    b.ChangeTags(getNewTag)

    func getNewTag(fieldName, fieldTag string) string {
        return strings.ToLower(fieldName)
    }
```

3. Fill structure from yaml

```golang
    err = yaml.Unmarshal([]byte("content for pkg.Structure"), b.Struct())
```

4. Fill structure from another struct

```golang
    var dst pkg.Structure
    err = b.AssignFrom(&dst)
```

5. Save structure data into external structure

```golang
    var dst pkg.Structure
    err = b.SaveInto(&dst)
```

### Build structure

1. Create structure builder

```golang
    builder := structure.NewBuilder()
```

2. Add fields to structure

```golang
    builder.AddField("Field1", "example-string", `yaml:"field1"`)
	builder.AddField("Field2", 1, `yaml:"field2"`)
	builder.AddField("Field3", false, `yaml:"field3"`)
```

3. Build structure interface

```golang
    st := builder.Build()
```

### Merge two structs

Merge two different structs.

```golang
	src := testSrcStructure{
		Field1: "src-field1",
		Field2: 1,
		Field3: 1,
		Field5: "src-field5",
	}

	dst := testDstStructure{
		Field1: "dst-field1",
		Field2: "dst-field2",
		Field3: 2,
		Field4: "dst-field4",
	}

    err := structure.Merge(&dst, &src)

    // result dst
    // testDstStructure{
	// 	Field1: "src-field1",
	// 	Field2: "dst-field2",
	// 	Field3: 1,
	// 	Field4: "dst-field4",
	// }
```

### Save structure to map

```golang
	src := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}


	dst := make(map[string]any)
	err := structure.SaveStructToMap(dst, &src)

    // result dst
    // map[string]any{
	// 	"Field1": "test-value",
	// 	"Field2": 3,
	// 	"Field3": true,
	// }
```

### Assign structure from map

```golang
	src := map[string]any{
		"Field1": "test-value",
		"Field2": 3,
		"Field3": true,
	}

	var dst testStructureWithoutTag
	err := structure.AssignStructFromMap(&dst, src)

    // result dst
    // testStructureWithoutTag{
	// 	Field1: "test-value",
	// 	Field2: 3,
	// 	Field3: true,
	// }
```

## TODO

    - Benchmarks
    - to add tag to substructs
