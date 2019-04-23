package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("ParseQuery test", func() {

	input := [][]string{{"description = description with multiple     spaces"},
		{"description = description with operators: [in = != eqornil gt lt in notin]"},
		{`description = description with \`},
		{`description in [description with "quotes"||description with \]`},
		{"type = type", `description = description with "quotes"`}}
	output := []string{"description+%3D+description+with+multiple+++++spaces",
		"description+%3D+description+with+operators%3A+%5Bin+%3D+%21%3D+eqornil+gt+lt+in+notin%5D",
		"description+%3D+description+with+%5C",
		"description+in+%5Bdescription+with+%22quotes%22%7C%7Cdescription+with+%5C%5D",
		"type+%3D+type|description+%3D+description+with+%22quotes%22"}

	Context("when queries are provided", func() {
		It("should url encode and join them", func() {
			for i := range input {
				Expect(ParseQuery(input[i])).To(Equal(output[i]))
			}
		})
	})

})
