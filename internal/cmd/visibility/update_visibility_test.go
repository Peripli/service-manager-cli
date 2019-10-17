package visibility

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func TestUpdateVisibilityCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Update visibility command test", func() {

	var client *smclientfakes.FakeClient
	var command *UpdateVisibilityCmd
	var buffer *bytes.Buffer
	var visibility *types.Visibility

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdateVisibilityCmd(context)
	})

	validUpdateVisibilityExecution := func(args ...string) error {
		visibility = &types.Visibility{
			ID:            "visibilityID",
			PlatformID:    "platformID",
			ServicePlanID: "planID",
		}
		Expect(json.Unmarshal([]byte(args[1]), visibility)).ToNot(HaveOccurred())
		client.UpdateVisibilityReturns(visibility, nil)
		uvCmd := command.Prepare(cmd.SmPrepare)
		uvCmd.SetArgs(args)
		return uvCmd.Execute()
	}

	invalidUpdateVisibilityExecution := func(args ...string) error {
		uvCmd := command.Prepare(cmd.SmPrepare)
		uvCmd.SetArgs(args)
		return uvCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("visibility should be updated", func() {

				err := validUpdateVisibilityExecution("id", `{"platform_id":"newPlatformID"}`)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(visibility.TableData().String()))
			})

			It("argument values should be as expected", func() {
				err := validUpdateVisibilityExecution("id", `{"platform_id":"newPlatformID"}`)

				id, v, _ := client.UpdateVisibilityArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id).To(Equal("id"))
				Expect(v).To(Equal(&types.Visibility{PlatformID: "newPlatformID"}))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				err := validUpdateVisibilityExecution("id", `{"platform_id":"newPlatformID"}`, "--output", "json")

				jsonByte, _ := json.MarshalIndent(visibility, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				err := validUpdateVisibilityExecution("platformId", `{"platform_id":"newPlatformID"}`, "--output", "yaml")

				yamlByte, _ := yaml.Marshal(visibility)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})

		Context("With generic parameter flag provided", func() {
			It("should pass it to SM", func() {
				err := validUpdateVisibilityExecution("platformId", `{"platform_id":"newPlatformID"}`, "--param", "paramKey=paramValue")
				Expect(err).ShouldNot(HaveOccurred())

				_, _, args := client.UpdateVisibilityArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("paramKey=paramValue"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With missing arguments", func() {
			It("Should return error missing id", func() {
				err := invalidUpdateVisibilityExecution()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("id is required"))
			})
			It("Should return error missing json", func() {
				err := invalidUpdateVisibilityExecution("id")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nothing to update. Visibility JSON is not provided"))
			})
		})
	})

	Context("With error from http client", func() {
		It("Should return error", func() {
			expectedErr := errors.New("http client error")
			client.UpdateVisibilityReturns(nil, expectedErr)

			err := invalidUpdateVisibilityExecution("id", `{"platform_id":"newPlatformID"}`)

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr.Error()))
		})
	})

	Context("With invalid output format", func() {
		It("should return error", func() {
			invFormat := "invalid-format"
			err := invalidUpdateVisibilityExecution("id", `{"platform_id":"newPlatformID"}`, "--output", invFormat)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown output: " + invFormat))
		})
	})
})
