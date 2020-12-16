package binding

import (
	"fmt"
	"github.com/Peripli/service-manager/pkg/util"
	"github.com/Peripli/service-manager/pkg/web"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Unbind command test", func() {
	var client *smclientfakes.FakeClient
	var command *UnbindCmd
	var buffer *bytes.Buffer
	var promptBuffer *bytes.Buffer
	var instances *types.ServiceInstances
	var bindings *types.ServiceBindings

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		promptBuffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUnbindCmd(context, promptBuffer)

		instances = &types.ServiceInstances{}
		instances.ServiceInstances = []types.ServiceInstance{{ID: "1234", Name: "instance-name"}}
		bindings = &types.ServiceBindings{}
		bindings.ServiceBindings = []types.ServiceBinding{{ID: "id", Name: "binding-name", ServiceInstanceID: "1234"}}
	})

	JustBeforeEach(func() {
		client.ListInstancesReturns(instances, nil)
		client.ListBindingsReturns(bindings, nil)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when existing binding is being deleted forcefully", func() {
		It("should list success message", func() {
			client.UnbindReturns("", nil)
			err := executeWithArgs("instance-name", "binding-name", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Binding successfully deleted."))
		})
	})
	Context("when existing binding is being deleted", func() {
		It("should list success message when confirmed", func() {
			client.UnbindReturns("", nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("instance-name", "binding-name")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Service Binding successfully deleted."))
		})

		It("should print delete declined when declined", func() {
			client.UnbindReturns("", nil)
			promptBuffer.WriteString("n")
			err := executeWithArgs("instance-name", "binding-name")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Delete declined"))
		})
	})

	Context("when generic parameter flag is used", func() {
		It("should pass it to SM", func() {
			client.UnbindReturns("", nil)
			promptBuffer.WriteString("y")
			param := "parameterKey=parameterValue"
			err := executeWithArgs("instance-name", "binding-name", "--param", param)
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.UnbindArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf(param, "async=true"))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when purge parameter flag is used", func() {
		It("should pass it to SM with force=true and cascade=true", func() {
			client.UnbindReturns("", nil)
			promptBuffer.WriteString("y")
			err := executeWithArgs("instance-name", "binding-name", "--purge")
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.UnbindArgsForCall(0)

			cascadeParam := fmt.Sprintf("%s=%s", web.QueryParamCascade, "true")
			forceParam := fmt.Sprintf("%s=%s", web.QueryParamForce, "true")
			asyncParam := "async=true"
			Expect(args.GeneralParams).To(ConsistOf(cascadeParam, forceParam, asyncParam))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("With sync flag", func() {
		It("should pass it to SM", func() {
			client.UnbindReturns("", nil)
			promptBuffer.WriteString("y")

			err := executeWithArgs("instance-name", "binding-name", "--mode", "sync")
			Expect(err).ShouldNot(HaveOccurred())

			_, args := client.UnbindArgsForCall(0)

			Expect(args.GeneralParams).To(ConsistOf("async=false"))
			Expect(args.FieldQuery).To(BeEmpty())
			Expect(args.LabelQuery).To(BeEmpty())
		})
	})

	Context("when non-existing bindings are being deleted", func() {
		BeforeEach(func() {
			bindings = &types.ServiceBindings{}
		})
		It("should return message", func() {
			err := executeWithArgs("instance-name", "non-existing-name", "-f")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(SatisfyAll(ContainSubstring("Service Binding"), ContainSubstring("not found")))
		})
	})

	Context("when more than one instance with given name found", func() {
		BeforeEach(func() {
			instances.ServiceInstances = append(instances.ServiceInstances, types.ServiceInstance{ID: "456", Name: "instance-name"})
		})
		It("should return message", func() {
			err := executeWithArgs("instance-name", "binding-name", "-f")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(SatisfyAll(ContainSubstring("more than one service instance with name"), ContainSubstring("found. Use --id flag to specify id of the binding to be deleted")))
		})
	})

	Context("when SM returns error", func() {
		It("should return error message", func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("")))
			expectedError := util.HandleResponseError(&http.Response{Body: body, StatusCode: http.StatusInternalServerError})
			client.UnbindReturns("", expectedError)
			err := executeWithArgs("instance-name", "binding-name", "-f")

			Expect(err).Should(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("Could not delete service binding. Reason:"))

		})
	})

	Context("when no arguments are provided", func() {
		It("should print required arguments", func() {
			client.UnbindReturns("", nil)
			err := executeWithArgs()

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("instance and binding names are required"))
		})
	})
})
