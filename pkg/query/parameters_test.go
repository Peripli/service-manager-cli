package query_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/Peripli/service-manager-cli/pkg/query"
)

func TestParameters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Query Parameters test", func() {
	DescribeTable("Encode",
		func(p *query.Parameters, queryStr string) {
			Expect(p.Encode()).To(Equal(queryStr))
		},

		Entry("Special characters in FieldQuery are properly encoded",
			&query.Parameters{FieldQuery: []string{`key = a in [123] & b != x\y | c = "c" | d = 'd'`}},
			"fieldQuery=key+%3D+a+in+%5B123%5D+%26+b+%21%3D+x%5Cy+%7C+c+%3D+%22c%22+%7C+d+%3D+%27d%27"),
		Entry("Special characters in LabelQuery are properly encoded",
			&query.Parameters{LabelQuery: []string{`key = a in [123] & b != x\y | c = "c" | d = 'd'`}},
			"labelQuery=key+%3D+a+in+%5B123%5D+%26+b+%21%3D+x%5Cy+%7C+c+%3D+%22c%22+%7C+d+%3D+%27d%27"),
		Entry("Special characters in GeneralParams are properly encoded",
			&query.Parameters{GeneralParams: []string{`key = a in [123] & b != x\y | c = "c" | d = 'd'`}},
			"key+=+a+in+%5B123%5D+%26+b+%21%3D+x%5Cy+%7C+c+%3D+%22c%22+%7C+d+%3D+%27d%27"),

		Entry("Multiple values for FieldQuery are properly encoded",
			&query.Parameters{FieldQuery: []string{`a = 1`, `b = 2`}},
			"fieldQuery=a+%3D+1%7Cb+%3D+2"),
		Entry("Multiple values for LabelQuery are properly encoded",
			&query.Parameters{LabelQuery: []string{`a = 1`, `b = 2`}},
			"labelQuery=a+%3D+1%7Cb+%3D+2"),
		Entry("Multiple values for GeneralParams are properly encoded",
			&query.Parameters{GeneralParams: []string{`a = 1`, `b = 2`}},
			"a+=+1&b+=+2"),

		Entry("GeneralParams value without = is properly encoded",
			&query.Parameters{GeneralParams: []string{"abc"}},
			"abc="),

		Entry("Values for all paramters are properly encoded sorted by key",
			&query.Parameters{
				FieldQuery:    []string{"a = 1"},
				LabelQuery:    []string{"b = 2"},
				GeneralParams: []string{"c = 3", "x = 9"},
			},
			"c+=+3&fieldQuery=a+%3D+1&labelQuery=b+%3D+2&x+=+9"),
	)
})
