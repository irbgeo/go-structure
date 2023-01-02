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
	Field1 string
	Field2 int
	Field3 bool
}

func getNewTag(fieldName, fieldTag string) string {
	return strings.ToLower(fieldName)
}

func TestSaveIntoMap(t *testing.T) {
	s, err := structure.New(new(testStructureWithoutTag))
	require.NoError(t, err)

	s.ChangeTags(getNewTag)

	valueStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}

	err = s.AssignFrom(valueStruct)
	require.NoError(t, err)

	actualMap := make(map[string]any)

	expectedMap := map[string]any{
		"Field1": "test-value",
		"Field2": 3,
		"Field3": true,
	}
	err = s.SaveInto(actualMap)
	require.NoError(t, err)
	require.Equal(t, expectedMap, actualMap)
}

func TestSaveStructIntoMap(t *testing.T) {
	valueStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}
	expectedMap := map[string]any{
		"Field1": "test-value",
		"Field2": 3,
		"Field3": true,
	}

	actualMap := make(map[string]any)
	err := structure.SaveStructToMap(actualMap, &valueStruct)
	require.NoError(t, err)
	require.Equal(t, expectedMap, actualMap)
}

func TestAssignFromMap(t *testing.T) {
	s, err := structure.New(new(testStructureWithoutTag))
	require.NoError(t, err)

	s.ChangeTags(getNewTag)

	expectedStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}

	valueMap := map[string]any{
		"Field1": "test-value",
		"Field2": 3,
		"Field3": true,
	}

	err = s.AssignFrom(valueMap)
	require.NoError(t, err)

	var actualStruct testStructureWithoutTag

	err = s.SaveInto(&actualStruct)
	require.NoError(t, err)
	require.Equal(t, expectedStruct, actualStruct)
}

func TestAssignStructFromMap(t *testing.T) {
	expectedStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}

	valueMap := map[string]any{
		"Field1": "test-value",
		"Field2": 3,
		"Field3": true,
	}

	var actualStruct testStructureWithoutTag
	err := structure.AssignStructFromMap(&actualStruct, valueMap)
	require.NoError(t, err)
	require.Equal(t, expectedStruct, actualStruct)
}

func TestUnmarshal(t *testing.T) {
	s, err := structure.New(new(testStructureWithoutTag))
	require.NoError(t, err)

	s.ChangeTags(getNewTag)

	err = yaml.Unmarshal([]byte(testContent), s.Struct())
	require.NoError(t, err)

	expectedStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}

	var actualStruct testStructureWithoutTag
	err = s.SaveInto(&actualStruct)
	require.NoError(t, err)
	require.Equal(t, expectedStruct, actualStruct)
}

func TestMarshal(t *testing.T) {
	s, err := structure.New(new(testStructureWithoutTag))
	require.NoError(t, err)

	s.ChangeTags(getNewTag)

	expectedStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}

	err = s.AssignFrom(expectedStruct)
	require.NoError(t, err)

	actualContent, err := yaml.Marshal(s.Struct())
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

func TestBuilder(t *testing.T) {
	builder := structure.NewBuilder()

	builder.AddField("Field1", "example-string", `yaml:"field1"`)
	builder.AddField("Field2", 1, `yaml:"field2"`)
	builder.AddField("Field3", false, `yaml:"field3"`)

	st := builder.Build()

	err := yaml.Unmarshal([]byte(testContent), st.Struct())
	require.NoError(t, err)

	expectedStruct := testStructureWithoutTag{
		Field1: "test-value",
		Field2: 3,
		Field3: true,
	}

	var actualStruct testStructureWithoutTag
	err = st.SaveInto(&actualStruct)
	require.NoError(t, err)
	require.Equal(t, expectedStruct, actualStruct)
}
