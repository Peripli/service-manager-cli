package instance

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

var _ = Describe("Get instance command test", func() {
	var client *smclientfakes.FakeClient
	var command *GetInstanceCmd
	var buffer *bytes.Buffer
	instance := types.ServiceInstance{
		Name:       "instance1",
		PlatformID: "platformID1",
		ID:         "id1",
	}
	instance2 := types.ServiceInstance{
		Name:       "instance1",
		PlatformID: "platformID2",
		ID:         "id2",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		client.ListInstancesReturns(&types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance}}, nil)
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewGetInstanceCmd(context)
	})

	executeWithArgs := func(args ...string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no instance name is provided", func() {
		It("should return error", func() {
			client.GetInstanceByIDReturns(&instance, nil)
			err := executeWithArgs("")

			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when more than one instance with same name exists", func() {
		var response *types.ServiceInstances
		BeforeEach(func() {
			response = &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance, instance2}, Vertical: true}
			client.ListInstancesReturns(response, nil)
		})

		It("should return both instances", func() {
			client.GetInstanceByIDReturnsOnCall(0, &instance, nil)
			client.GetInstanceByIDReturnsOnCall(1, &instance2, nil)
			err := executeWithArgs("instance1")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(buffer.String()).To(ContainSubstring(response.TableData().String()))
		})
	})

	Context("when no known instance name is provided", func() {
		It("should return no instance", func() {
			client.ListInstancesReturns(&types.ServiceInstances{}, nil)
			err := executeWithArgs("unknown")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("No instance found with name: unknown"))
		})
	})

	Context("when instance with name is found", func() {
		It("should return its data", func() {
			client.GetInstanceByIDReturns(&instance, nil)
			err := executeWithArgs("instance1")

			Expect(err).ShouldNot(HaveOccurred())
			result := &types.ServiceInstances{ServiceInstances: []types.ServiceInstance{instance}, Vertical: true}
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})
})
