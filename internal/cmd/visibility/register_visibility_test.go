package visibility

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRegisterVisibilityCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Register visibility command test", func() {

	var client *smclientfakes.FakeClient
	var command *RegisterVisibilityCmd
	var buffer *bytes.Buffer
	var visibility *types.Visibility

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewRegisterVisibilityCmd(context)
	})

	validRegisterVisibilityExecution := func(args []string) error {
		visibility = &types.Visibility{
			ID:            "visibilityID",
			PlatformID:    args[0],
			ServicePlanID: args[1],
		}
		client.RegisterVisibilityReturns(visibility, nil)
		rvCmd := command.Prepare(cmd.SmPrepare)
		rvCmd.SetArgs(args)
		return rvCmd.Execute()
	}

	invalidRegisterVisibilityExecution := func(args []string) error {
		rvCmd := command.Prepare(cmd.SmPrepare)
		rvCmd.SetArgs(args)
		return rvCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("visibility should be registered", func() {
				err := validRegisterVisibilityExecution([]string{"platformId", "planId"})
				tableOutputExpected := visibility.TableData().String()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("argument values should be as expected", func() {
				err := validRegisterVisibilityExecution([]string{"platformId", "planId"})

				v := client.RegisterVisibilityArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(v.PlatformID).To(Equal("platformId"))
				Expect(v.ServicePlanID).To(Equal("planId"))
			})
		})

		Context("With id flag provided", func() {
			It("visibility id should be as expected", func() {
				args := []string{"platformId", "planId", "--id", "id"}

				err := validRegisterVisibilityExecution(args)
				v := client.RegisterVisibilityArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(v.ID).To(Equal("id"))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				err := validRegisterVisibilityExecution([]string{"platformId", "planId", "--output", "json"})

				jsonByte, _ := json.MarshalIndent(visibility, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				err := validRegisterVisibilityExecution([]string{"platformId", "planId", "--output", "yaml"})

				yamlByte, _ := yaml.Marshal(visibility)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With not enough arguments provided", func() {
			It("Should return error", func() {
				err := invalidRegisterVisibilityExecution([]string{"platformId"})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("platform-id and service-plan-id required but not provided"))
			})
		})

		Context("With error from http client", func() {
			It("Should return error", func() {
				expectedErr := errors.New("http client error")
				client.RegisterVisibilityReturns(nil, expectedErr)

				err := invalidRegisterVisibilityExecution([]string{"platformId", "planId"})

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError(expectedErr.Error()))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidRegisterVisibilityExecution([]string{"platformId", "planId", "--output", invFormat})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown output: " + invFormat))
			})
		})
	})
})
