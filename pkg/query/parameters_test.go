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

		Entry("Empty parameters are properly encoded",
			&query.Parameters{},
			""),
		Entry("nil parameters are properly encoded",
			nil,
			""),

		Entry("Special characters in FieldQuery are properly encoded",
			&query.Parameters{FieldQuery: []string{`key = a in (123) & b != x\y and c = "c" and d = 'd'`}},
			"fieldQuery=key+%3D+a+in+%28123%29+%26+b+%21%3D+x%5Cy+and+c+%3D+%22c%22+and+d+%3D+%27d%27"),
		Entry("Special characters in LabelQuery are properly encoded",
			&query.Parameters{LabelQuery: []string{`key = a in (123) & b != x\y and c = "c" and d = 'd'`}},
			"labelQuery=key+%3D+a+in+%28123%29+%26+b+%21%3D+x%5Cy+and+c+%3D+%22c%22+and+d+%3D+%27d%27"),
		Entry("Special characters in GeneralParams are properly encoded",
			&query.Parameters{GeneralParams: []string{`key = a in (123) & b != x\y and c = "c" and d = 'd'`}},
			"key+=+a+in+%28123%29+%26+b+%21%3D+x%5Cy+and+c+%3D+%22c%22+and+d+%3D+%27d%27"),

		Entry("Multiple values for FieldQuery are properly encoded",
			&query.Parameters{FieldQuery: []string{`a = 1`, `b = 2`}},
			"fieldQuery=a+%3D+1+and+b+%3D+2"),
		Entry("Multiple values for LabelQuery are properly encoded",
			&query.Parameters{LabelQuery: []string{`a = 1`, `b = 2`}},
			"labelQuery=a+%3D+1+and+b+%3D+2"),
		Entry("Multiple values for GeneralParams are properly encoded",
			&query.Parameters{GeneralParams: []string{`a = 1`, `b = 2`}},
			"a+=+1&b+=+2"),

		Entry("GeneralParams value without = is properly encoded",
			&query.Parameters{GeneralParams: []string{"abc"}},
			"abc="),

		Entry("Values for all parameters are properly encoded sorted by key",
			&query.Parameters{
				FieldQuery:    []string{"a = 1"},
				LabelQuery:    []string{"b = 2"},
				GeneralParams: []string{"c = 3", "x = 9"},
			},
			"c+=+3&fieldQuery=a+%3D+1&labelQuery=b+%3D+2&x+=+9"),
	)
})
