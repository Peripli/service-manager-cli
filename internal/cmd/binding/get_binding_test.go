package binding

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Get binding command test", func() {
	var client *smclientfakes.FakeClient
	var command *GetBindingCmd
	var buffer *bytes.Buffer
	binding := types.ServiceBinding{
		Name:              "binding1",
		ServiceInstanceID: "1",
		ID:                "id1",
	}
	binding2 := types.ServiceBinding{
		Name:              "binding1",
		ServiceInstanceID: "2",
		ID:                "id2",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		client.ListBindingsReturns(&types.ServiceBindings{ServiceBindings: []types.ServiceBinding{binding}}, nil)
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewGetBindingCmd(context)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no binding name is provided", func() {
		It("should return error", func() {
			client.GetBindingByIDReturns(&binding, nil)
			err := executeWithArgs("")

			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when more than one binding with same name exists", func() {
		var response *types.ServiceBindings
		BeforeEach(func() {
			response = &types.ServiceBindings{ServiceBindings: []types.ServiceBinding{binding, binding2}, Vertical: true}
			client.ListBindingsReturns(response, nil)
		})

		It("should return both bindings", func() {
			client.GetBindingByIDReturnsOnCall(0, &binding, nil)
			client.GetBindingByIDReturnsOnCall(1, &binding2, nil)
			err := executeWithArgs("binding1")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(buffer.String()).To(ContainSubstring(response.TableData().String()))
		})
	})

	Context("when no known binding name is provided", func() {
		It("should return no binding", func() {
			client.ListBindingsReturns(&types.ServiceBindings{}, nil)
			err := executeWithArgs("unknown")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("No binding found with name: unknown"))
		})
	})

	Context("when binding with name is found", func() {
		It("should return its data", func() {
			client.GetBindingByIDReturns(&binding, nil)
			err := executeWithArgs("binding1")

			Expect(err).ShouldNot(HaveOccurred())
			result := &types.ServiceBindings{ServiceBindings: []types.ServiceBinding{binding}, Vertical:true}
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})
})
