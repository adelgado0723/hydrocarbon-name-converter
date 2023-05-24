package converter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverter(t *testing.T) {
	t.Run("returns empty input error on blank string", func(t *testing.T) {
		iupacName, err := Convert("")
		assert.Equal(t, "", iupacName)
		assert.Equal(t, err.Error(), "Cannot set data from empty SMILES string")
	})

	t.Run("returns invalid string error if input is not composed of 'C', '(' and ')' )", func(t *testing.T) {
		iupacName, err := Convert("12334sdhjhkkjsd")
		assert.Equal(t, "", iupacName)
		assert.Equal(t, err.Error(), "SMILES string contains invalid characters")
	})

	t.Run("returns invalid string error if input has unbalanced parentheses", func(t *testing.T) {
		iupacName, err := Convert("C(C")
		assert.Equal(t, "", iupacName)
		assert.Equal(t, err.Error(), "Invalid input string. Parentheses are not balanced.")
	})

	t.Run("converts c to methane", func(t *testing.T) {
		iupacName, err := Convert("c")
		assert.Nil(t, err)
		assert.Equal(t, "methane", iupacName)
	})

	t.Run("converts ` c ` to methane", func(t *testing.T) {
		iupacName, err := Convert(" c ")
		assert.Nil(t, err)
		assert.Equal(t, "methane", iupacName)
	})

	t.Run("converts C to methane", func(t *testing.T) {
		iupacName, err := Convert("C")
		assert.Nil(t, err)
		assert.Equal(t, "methane", iupacName)
	})

	t.Run("converts CC to ethane", func(t *testing.T) {
		iupacName, err := Convert("CC")
		assert.Nil(t, err)
		assert.Equal(t, "ethane", iupacName)
	})

	t.Run("converts CCC to propane", func(t *testing.T) {
		iupacName, err := Convert("CCC")
		assert.Nil(t, err)
		assert.Equal(t, "propane", iupacName)
	})

	t.Run("converts CCCC to butane", func(t *testing.T) {
		iupacName, err := Convert("CCCC")
		assert.Nil(t, err)
		assert.Equal(t, "butane", iupacName)
	})

	t.Run("converts CCCCC to pentane", func(t *testing.T) {
		iupacName, err := Convert("CCCCC")
		assert.Nil(t, err)
		assert.Equal(t, "pentane", iupacName)
	})

	t.Run("converts CCCCCC to hexane", func(t *testing.T) {
		iupacName, err := Convert("CCCCCC")
		assert.Nil(t, err)
		assert.Equal(t, "hexane", iupacName)
	})

	t.Run("converts CCCCCCC to heptane", func(t *testing.T) {
		iupacName, err := Convert("CCCCCCC")
		assert.Nil(t, err)
		assert.Equal(t, "heptane", iupacName)
	})

	t.Run("converts CCCCCCCC to octane", func(t *testing.T) {
		iupacName, err := Convert("CCCCCCCC")
		assert.Nil(t, err)
		assert.Equal(t, "octane", iupacName)
	})

	t.Run("converts CCCCCCCCC to nonane", func(t *testing.T) {
		iupacName, err := Convert("CCCCCCCCC")
		assert.Nil(t, err)
		assert.Equal(t, "nonane", iupacName)
	})

	t.Run("converts CCCCCCCCCC to decane", func(t *testing.T) {
		iupacName, err := Convert("CCCCCCCCCC")
		assert.Nil(t, err)
		assert.Equal(t, "decane", iupacName)
	})

	t.Run("returns invalid arg error when carbon chain is longer than ten", func(t *testing.T) {
		iupacName, err := Convert("CCCCCCCCCCC")
		assert.Equal(t, "", iupacName)
		assert.Equal(t, "Cannot determine name for hydrocarbon of length 11", err.Error())
	})

	t.Run("returns invalid arg error when carbon chain with branches is longer than ten", func(t *testing.T) {
		iupacName, err := Convert("CCCCCCCC(CCCCCCCCCCCCCCC)CC")
		assert.Equal(t, "", iupacName)
		assert.Equal(t, "Cannot determine carbon molecule for chain longer than 10", err.Error())
	})

	t.Run("converts CC(C)CC to 2-methylbutane", func(t *testing.T) {
		iupacName, err := Convert("CC(C)CC")
		assert.Nil(t, err)
		assert.Equal(t, "2-methylbutane", iupacName)
	})

	t.Run("handles multiple branches of the same length", func(t *testing.T) {
		iupacName, err := Convert("CC(C)C(C)C")
		assert.Nil(t, err)
		assert.Equal(t, "2,3-dimethylbutane", iupacName)
	})

	t.Run("handles multiple branches of differing lengths", func(t *testing.T) {
		iupacName, err := Convert("CC(C)C(CC)C")
		assert.Nil(t, err)
		assert.Equal(t, "3-ethyl-2-methylbutane", iupacName)
	})

	t.Run(
		"handles multiple branches of differing lengths including multiple of same length",
		func(t *testing.T) {
			iupacName, err := Convert("CC(C)C(CC)C(C)C")
			assert.Nil(t, err)
			assert.Equal(t, "3-ethyl-2,4-dimethylpentane", iupacName)
		},
	)

	t.Run("handles multiple branches off a single carbon atom", func(t *testing.T) {
		iupacName, err := Convert("CC(C)(C)CC")
		fmt.Println(err)
		fmt.Println(iupacName)
		assert.Nil(t, err)
		assert.Equal(t, "2-dimethylbutane", iupacName)
	})

}
