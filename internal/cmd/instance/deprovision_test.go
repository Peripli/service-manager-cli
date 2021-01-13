package instance

import (
	"fmt"
	"github.com/Peripli/service-manager/pkg/web"
	"io/ioutil"
	"net/http"

	"github.com/Peripli/service-manager/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Deprovision command test", func() {
	var client *smclientfakes.FakeClient
	var command *DeprovisionCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer
	var instances *types.ServiceInstances

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewDeprovisionCmd(context, promptBuffer)

		instances = &types.ServiceInstances{}
		instances.ServiceInstances = []types.ServiceInstance{{ID: "1234", Name: "instance-name"}}
	})

	JustBeforeEach(func() {
		client.ListInstancesReturns(instances, nil)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing instance is being deleted forcefully", func() {
		It("should list success message", func() {
			client.DeprovisionReturns("", nil)
			err := executeWithArgs("instance-name", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Instance successfully deleted."))
		})
	})

	Context("when existing instance is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.DeprovisionReturns("", nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("instance-name")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Instance successfully deleted."))
		})

		It("should print delete declined when declined", func() {
			client.DeprovisionReturns("", nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs("instance-name")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when generic parameter flag is used", func() {
		It("should pass it to SM", func() {
			client.DeprovisionReturns("", nil)
			promptBuffer.WriteString("y")
			param := "parameterKey=parameterValue"
			err := executeWithArgs("instance-name", "--param", param)
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.DeprovisionArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param, "async=true"))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when forceDelete flag is used", func() {
		It("should pass it to SM", func() {
			client.DeprovisionReturns("", nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("instance-name", "--force-delete")
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.DeprovisionArgsForCall(0)

			cascadeParam := fmt.Sprintf("%s=%s", web.QueryParamCascade, "true")
			forceParam := fmt.Sprintf("%s=%s", web.QueryParamForce, "true")
			asyncParam := "async=true"
			Expect(args.GeneralParams).To(ConsistOf(cascadeParam, forceParam, asyncParam))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})
	Context("when force-delete parameter not in use", func() {
		It("should not pass the parameters (force & cascade) to SM", func() {
			client.DeprovisionReturns("", nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("instance-name", "--force")
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.DeprovisionArgsForCall(0)

			cascadeParam := fmt.Sprintf("%s=%s", web.QueryParamCascade, "true")
			forceParam := fmt.Sprintf("%s=%s", web.QueryParamForce, "true")
			Expect(args.GeneralParams).ToNot(ConsistOf(cascadeParam, forceParam))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("With sync flag", func() {
		It("should pass it to SM", func() {
			client.DeprovisionReturns("", nil)
			promptBuffer.WriteString("y")

			err := executeWithArgs("instance-name", "--mode", "sync")
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.DeprovisionArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf("async=false"))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when non-existing instances are being deleted", func() {
		BeforeEach(func() {
			instances = &types.ServiceInstances{}
		})
		It("should return message", func() {
			err := executeWithArgs("non-existing-name", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Instance not found"))
		})
	})

	Context("when more than one instance with given name found", func() {
		BeforeEach(func() {
			instances.ServiceInstances = append(instances.ServiceInstances, types.ServiceInstance{ID: "456", Name: "instance-name"})
		})
		It("should return message", func() {
			err := executeWithArgs("instance-name", "-f")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(SatisfyAll(ContainSubstring("more than one service instance with name"), ContainSubstring("found. Use --id flag to specify id of the instance to be deleted")))
		})
	})

	Context("when SM returns error", func() {
		It("should return error message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusInternalServerError})
			client.DeprovisionReturns("", expectedError)
			err := executeWithArgs("name", "-f")

			Expect(err).Should(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Could not delete service instance. Reason:"))

		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.DeprovisionReturns("", nil)
			err := executeWithArgs()

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("single [name] is required"))
		})
	})
})
