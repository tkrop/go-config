package reflect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkrop/go-testing/test"
)

// TestToComplexPanic tests the defensive error path in toComplex
// for types that should never be passed through the public API.
func TestToComplexPanic(t *testing.T) {
	// Given
	parsed := "test"
	fieldType := reflect.TypeOf(parsed)
	defer test.Recover(t, NewErrTagWalker("unsupported type",
		"<unknown>", fieldType, nil))

	// When
	result, err := toComplex(parsed, fieldType)

	// Then
	assert.Nil(t, result)
	assert.NoError(t, err)
}
