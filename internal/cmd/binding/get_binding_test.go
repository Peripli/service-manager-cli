package binding

import (
	"github.com/Peripli/service-manager-cli/internal/output"
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

	instance1 := types.ServiceInstance{
		ID:   "1",
		Name: "instance-name1",
	}

	instance2 := types.ServiceInstance{
		ID:   "2",
		Name: "instance-name2",
	}
	binding := types.ServiceBinding{
		Name:                "binding1",
		ServiceInstanceID:   "1",
		ServiceInstanceName: "instance-name1",
		ID:                  "id1",
	}
	binding2 := types.ServiceBinding{
		Name:                "binding1",
		ServiceInstanceID:   "2",
		ServiceInstanceName: "instance-name2",
		ID:                  "id2",
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

	Describe("Get service binding", func() {
		BeforeEach(func() {
			client.GetInstanceByIDReturnsOnCall(0, &instance1, nil)
			client.GetInstanceByIDReturnsOnCall(1, &instance2, nil)
		})
		When("no binding name is provided", func() {
			It("should return error", func() {
				client.GetBindingByIDReturns(&binding, nil)
				err := executeWithArgs("")

				Expect(err).Should(HaveOccurred())
			})
		})

		When("more than one binding with same name exists", func() {
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

		When("no known binding name is provided", func() {
			It("should return no binding", func() {
				client.ListBindingsReturns(&types.ServiceBindings{}, nil)
				err := executeWithArgs("unknown")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("No binding found with name: unknown"))
			})
		})

		When("binding with name is found", func() {
			It("should return its data", func() {
				client.GetBindingByIDReturns(&binding, nil)
				err := executeWithArgs("binding1")

				Expect(err).ShouldNot(HaveOccurred())
				result := &types.ServiceBindings{ServiceBindings: []types.ServiceBinding{binding}, Vertical: true}
				Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
			})
		})
	})

	Describe("Get service binding parameters", func() {
		bindingParameters1 := map[string]interface{}{"param1":"value1","param2":"value2"}
		bindingParameters2 := make(map[string]interface{})

		When("no binding name is provided", func() {
			It("should return error", func() {
				err := executeWithArgs("", "--show-binding-params")

				Expect(err).Should(HaveOccurred())
			})
		})

		When("no known binding name is provided", func() {
			It("should print no binding found", func() {
				client.ListBindingsReturns(&types.ServiceBindings{}, nil)
				err := executeWithArgs("unknown", "--show-binding-params")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("No binding found with name: unknown"))
			})
		})

		When("there is binding with this name with parameters", func() {
			It("should print parameters", func() {
				client.GetBindingParametersReturns(bindingParameters1, nil)
				err := executeWithArgs("binding1", "--show-binding-params")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(output.PrintParameters(bindingParameters1)))
			})
		})

		When("there is instance with this name without parameters", func() {
			It("should print no parameters", func() {
				client.GetBindingParametersReturns(bindingParameters2, nil)
				err := executeWithArgs("binding1", "--show-binding-params")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("No configuration parameters are set for service binding id: %s", binding.ID))
			})
		})

		When("two bindings with same name exists," +
			"one with parameters and the second without parameters", func() {
			var response *types.ServiceBindings
			BeforeEach(func() {
				response = &types.ServiceBindings{ServiceBindings: []types.ServiceBinding{binding, binding2}, Vertical: true}
				client.ListBindingsReturns(response, nil)
			})

			It("should print parameters for both bindings", func() {
				client.GetBindingParametersReturnsOnCall(0, bindingParameters1 , nil)
				client.GetBindingParametersReturnsOnCall(1, bindingParameters2, nil)
				err := executeWithArgs("binding1", "--show-binding-params")
				Expect(err).ShouldNot(HaveOccurred())

				Expect(buffer.String()).To(ContainSubstring(output.PrintParameters(bindingParameters1)))
				Expect(buffer.String()).To(ContainSubstring("No configuration parameters are set for service binding id: %s", binding2.ID))
			})
		})
	})
})
