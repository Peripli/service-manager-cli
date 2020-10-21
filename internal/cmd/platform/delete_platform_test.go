package platform

import (
	"github.com/Peripli/service-manager/pkg/util"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Delete platforms command test", func() {

	var client *smclientfakes.FakeClient
	var command *DeletePlatformCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeletePlatformCmd(context, promptBuffer)

		platforms := &types.Platforms{}
		platforms.Platforms = []types.Platform{{ID: "1234", Name: "platform-name"}, {ID: "456", Name: "platform2"}}
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Describe("delete platform", func() {

		When("existing platform is being deleted forcefully", func() {
			It("should list success message", func() {
				client.DeletePlatformsReturns(nil)
				err := executeWithArgs([]string{"platform-name", "-f"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Platform(s) successfully deleted."))
			})
		})

		When("existing platform is being deleted", func() {
			It("should list success message when confirmed", func() {
				client.DeletePlatformsReturns(nil)
				promptBuffer.WriteString("y")
				err := executeWithArgs([]string{"platform-name"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Platform(s) successfully deleted."))
			})

			It("should print delete declined when declined", func() {
				client.DeletePlatformsReturns(nil)
				promptBuffer.WriteString("n")
				err := executeWithArgs([]string{"platform-name"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Delete declined"))
			})
		})

		When("generic parameter flag is used", func() {
			It("should pass it to SM", func() {
				client.DeletePlatformsReturns(nil)
				promptBuffer.WriteString("y")
				param := "parameterKey=parameterValue"
				err := executeWithArgs([]string{"platform-name", "--param", param})
				Expect(err).ShouldNot(HaveOccurred())

				args := client.DeletePlatformsArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf(param))
				Expect(args.FieldQuery).To(ConsistOf("name eq 'platform-name'"))
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})

		When("non-existing platform is being deleted", func() {
			It("should return message", func() {
				body := ioutil.NopCloser(bytes.NewReader([]byte("")))
				expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusNotFound})
				client.DeletePlatformsReturns(expectedError)
				err := executeWithArgs([]string{"non-existing-name", "-f"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Platform(s) not found."))
			})
		})

		When("SM returns error", func() {
			It("should return error message", func() {
				body := ioutil.NopCloser(bytes.NewReader([]byte("")))
				expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusInternalServerError})
				client.DeletePlatformsReturns(expectedError)
				err := executeWithArgs([]string{"name", "-f"})

				Expect(err).Should(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Could not delete platform(s). Reason:"))

			})
		})

		When("no arguments are provided", func() {
			It("should print required arguments", func() {
				client.DeletePlatformsReturns(nil)
				err := executeWithArgs([]string{})

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError("single [name] is required"))
			})
		})
	})

	Describe("cascade delete platform", func() {

		platform1 := types.Platform{
			Name:       "platform1",
			ID:         "id1",
		}

		When("platform exists", func() {
			It("should print cascade delete successfully scheduled", func() {
				client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{platform1}}, nil)
				location := "/v1/platforms/id1/operations/1a3e795d-819c-4661-89b5-344adb2ec26a"

				client.CascadeDeletePlatformReturns(location, nil)
				err := executeWithArgs([]string{platform1.Name, "--cascade-delete", "-f"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Cascade delete successfully scheduled"))
			})
		})

		When("platform does not exist", func() {
			It("should print platform(s) not found", func() {
				client.ListPlatformsReturns(&types.Platforms{}, nil)
				err := executeWithArgs([]string{"non-existing-name", "--cascade-delete", "-f"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Platform(s) not found."))
			})
		})

		When("error is returned", func() {
			It("should print could not cascade delete platform", func() {
				client.ListPlatformsReturns(&types.Platforms{Platforms: []types.Platform{platform1}}, nil)
				body := ioutil.NopCloser(bytes.NewReader([]byte("")))
				expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusInternalServerError})
				client.CascadeDeletePlatformReturns("", expectedError)
				promptBuffer.WriteString("y")
				_ = executeWithArgs([]string{platform1.Name, "--cascade-delete"})

				Expect(buffer.String()).To(ContainSubstring("Could not cascade-delete platform"))
			})
		})

		When("no platform name is provided", func() {
			It("should print required arguments", func() {
				err := executeWithArgs([]string{"", "--cascade-delete"})

				Expect(err).Should(HaveOccurred())
				Expect(err).To(MatchError("single [name] is required"))
			})
		})
	})

})
