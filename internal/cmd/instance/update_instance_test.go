package instance

import (
	"bytes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"errors"
)

var _ = Describe("update instance command test", func() {
	var client *smclientfakes.FakeClient
	var command *UpdateCmd
	var buffer *bytes.Buffer
	/*
		var offerings *types.ServiceOfferings
		var plans *types.ServicePlans
		var instance *types.ServiceInstance
	*/
	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdateInstanceCmd(context)
	})
	/*
		validAsyncProvisionExecution := func(location string, args ...string) *cobra.Command {
			offerings = &types.ServiceOfferings{
				ServiceOfferings: []types.ServiceOffering{
					{ID: OfferingID, Name: args[1]},
				},
			}
			plans = &types.ServicePlans{
				ServicePlans: []types.ServicePlan{
					{ID: PlanID, Name: args[2]},
				},
			}
			instance = &types.ServiceInstance{
				Name: args[0],
			}
			operation := &types.Operation{
				State: "in progress",
			}
			client.StatusReturns(operation, nil)
			client.ListOfferingsReturns(offerings, nil)
			client.ListPlansReturns(plans, nil)
			client.ProvisionReturns(instance, location, nil)

			piCmd := command.Prepare(cmd.SmPrepare)
			piCmd.SetArgs(args)
			Expect(piCmd.Execute()).ToNot(HaveOccurred())

			return piCmd
		}

		validSyncProvisionExecution := func(args ...string) *cobra.Command {
			return validAsyncProvisionExecution("", append(args, "--mode", "sync")...)
		}

	*/
	invalidUpdateInstanceCommandExecution := func(args ...string) error {
		trCmd := command.Prepare(cmd.SmPrepare)
		trCmd.SetArgs(args)
		return trCmd.Execute()
	}

	Describe("Invalid request", func() {
		var instances *types.ServiceInstances
		var instance *types.ServiceInstance
		var plan *types.ServicePlan
		var plans *types.ServicePlans
		var errGetInstance error
		JustBeforeEach(func() {
			client.ListInstancesReturns(instances, nil)
			client.GetPlanByIDReturns(plan, nil)
			client.GetInstanceByIDReturns(instance, errGetInstance)
			client.ListPlansReturns(plans, nil)
		})

		Context("when service instance not found", func() {
			BeforeEach(func() {
				instances = &types.ServiceInstances{}
			})

			Context("by name", func() {
				It("should return an error", func() {
					err := invalidUpdateInstanceCommandExecution("instance-name", "--new-name", "new name")
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("no instances found with name %s", "instance-name")))
				})

			})
			Context("by id", func() {
				BeforeEach(func() {
					errGetInstance = errors.New("errore occured")
				})
				It("should return an error", func() {
					err := invalidUpdateInstanceCommandExecution("instance-name", "--new-name", "new name", "--id", "service-instance-id")
					Expect(err.Error()).To(ContainSubstring("errore occured"))
				})

			})

		})

		Context("when more than once instance found", func() {
			BeforeEach(func() {
				client.ListInstancesReturnsOnCall(0, &types.ServiceInstances{
					ServiceInstances: []types.ServiceInstance{
						types.ServiceInstance{
							Name: "instance-name",
						},
						types.ServiceInstance{
							Name: "instance-name",
						},
					},
				}, nil)
			})
			It("should return an error", func() {
				err := invalidUpdateInstanceCommandExecution("instance-name", "--new-name", "new name")
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("more than 1 instance found with name %s", "instance-name")))
			})

		})
		Context("update plan", func() {
			BeforeEach(func() {
				instances = &types.ServiceInstances{
					ServiceInstances: []types.ServiceInstance{
						{ID: "instanceid", Name: "instance"},
					},
				}
				plan = &types.ServicePlan{
					ID:                "plandid",
					ServiceOfferingID: "offeringid",
					CatalogName:       "large",
				}

			})
			Context("when service plan is not found", func() {
				BeforeEach(func() {
					plans = &types.ServicePlans{}
				})
				It("should return an error", func() {
					err := invalidUpdateInstanceCommandExecution("instance-name", "--new-name", "new name", "--plan", plan.CatalogName)
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("service plan with name %s for offering with id %s not found", plan.CatalogName, plan.ServiceOfferingID)))
				})

			})
			Context("when more than one service plan found", func() {
				BeforeEach(func() {
					plans = &types.ServicePlans{
						ServicePlans: []types.ServicePlan{
							{ID: "plandid",
								ServiceOfferingID: "offeringid",
								CatalogName:       "large"},
							{ID: "plandid2",
								ServiceOfferingID: "offeringid",
								CatalogName:       "large"},
						},
					}
				})

				It("should return an error", func() {
					err := invalidUpdateInstanceCommandExecution("instance-name", "--new-name", "new name", "--plan", plan.CatalogName)
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("exactly one service plan with name %s for offering with id %s expected", plan.CatalogName, plan.ServiceOfferingID)))
				})
			})

		})

	})
})
