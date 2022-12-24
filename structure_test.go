package structure_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	structure "github.com/irbgeo/go-structure"
)

var testContent = `field1: test-value
field2: 3
field3: true
`

type testStructureWithoutTag struct {
	FieLD1 string
	FielD2 int
	FIeld3 bool
}

func TestUnmarshal(t *testing.T) {
	b, err := structure.New(new(testStructureWithoutTag))
	require.NoError(t, err)

	b.AddTags(getTag)

	err = yaml.Unmarshal([]byte(testContent), b.Struct())
	require.NoError(t, err)

	expectedStruct := testStructureWithoutTag{
		FieLD1: "test-value",
		FielD2: 3,
		FIeld3: true,
	}

	var actualStruct testStructureWithoutTag
	err = b.SaveInto(&actualStruct)
	require.NoError(t, err)
	require.Equal(t, expectedStruct, actualStruct)
}

func TestMarshal(t *testing.T) {
	b, err := structure.New(new(testStructureWithoutTag))
	require.NoError(t, err)

	b.AddTags(getTag)

	expectedStruct := testStructureWithoutTag{
		FieLD1: "test-value",
		FielD2: 3,
		FIeld3: true,
	}

	err = b.AssignFrom(expectedStruct)
	require.NoError(t, err)

	actualContent, err := yaml.Marshal(b.Struct())
	require.NoError(t, err)
	require.Equal(t, testContent, string(actualContent))
}

func TestMerge(t *testing.T) {
	type testSrcStructure struct {
		Field1 string
		Field2 int
		Field3 int
		Field5 string
	}

	src := testSrcStructure{
		Field1: "src-field1",
		Field2: 1,
		Field3: 1,
		Field5: "src-field5",
	}

	type testDstStructure struct {
		Field1 string
		Field2 string
		Field3 int
		Field4 string
	}

	dst := testDstStructure{
		Field1: "dst-field1",
		Field2: "dst-field2",
		Field3: 2,
		Field4: "dst-field4",
	}

	expected := testDstStructure{
		Field1: "src-field1",
		Field2: "dst-field2",
		Field3: 1,
		Field4: "dst-field4",
	}

	err := structure.Merge(&dst, &src)
	require.NoError(t, err)
	require.Equal(t, expected, dst)
}

func getTag(fieldName string) string {
	return strings.ToLower(fieldName)
}
