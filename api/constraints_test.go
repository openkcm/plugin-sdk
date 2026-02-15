package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openkcm/plugin-sdk/api"
)

func TestConstraints(t *testing.T) {
	t.Run("exactly one", func(t *testing.T) {
		testConstraint(t, api.ExactlyOne(),
			"expected exactly 1 but got 0",
			"expected exactly 1 but got 2",
		)
	})

	t.Run("maybe one", func(t *testing.T) {
		testConstraint(t, api.MaybeOne(),
			"",
			"expected at most 1 but got 2",
		)
	})

	t.Run("at least one", func(t *testing.T) {
		testConstraint(t, api.AtLeastOne(),
			"expected at least 1 but got 0",
			"",
		)
	})

	t.Run("zero or more", func(t *testing.T) {
		testConstraint(t, api.ZeroOrMore(),
			"",
			"",
		)
	})
}

func testConstraint(t *testing.T, constraints api.Constraints, zeroError, twoError string) {
	testCheck(t, constraints, 0, zeroError)
	testCheck(t, constraints, 1, "")
	testCheck(t, constraints, 2, twoError)
}

func testCheck(t *testing.T, constraints api.Constraints, count int, expectedErr string) {
	err := constraints.Check(count)
	if expectedErr == "" {
		assert.NoError(t, err)
	} else {
		assert.EqualError(t, err, expectedErr)
	}
}
