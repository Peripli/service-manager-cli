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
	"io/ioutil"
	"github.com/Peripli/service-manager/pkg/util"
	"net/http"
	"github.com/spf13/cobra"
	"encoding/json"
)

var _ = Describe("update instance command test", func() {
	var client *smclientfakes.FakeClient
	var command *UpdateCmd
	var buffer *bytes.Buffer
	var instances *types.ServiceInstances
	var instance *types.ServiceInstance
	var updatedIntance *types.ServiceInstance
	var plan *types.ServicePlan
	var plans *types.ServicePlans
	var errGetInstance error

	const instanceId = "instanceid"
	const planId = "planId"
	const newPlanId = "plandidnew"
	const offeringId = "offeringId"
	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdateInstanceCmd(context)
		errGetInstance = nil
	})

	validAsyncUpdateExecution := func(location string, args ...string) *cobra.Command {
		var planName, instanceName, newInstanceName, instanceParams, serviceInstanceId string
		instanceName = args[0]
		for i, value := range args {
			switch value {
			case "--plan":
				planName = args[i+1]
			case "--new-name":
				newInstanceName = args[i+1]
			case "--instance-params":
				instanceParams = args[i+1]
			case "id":
				serviceInstanceId = args[i+1]
			}
		}
		instance = &types.ServiceInstance{
			ID:            instanceId,
			Name:          instanceName,
			ServicePlanID: newPlanId,
		}

		instances = &types.ServiceInstances{
			ServiceInstances: []types.ServiceInstance{types.ServiceInstance{
				ID:            instanceId,
				Name:          instanceName,
				ServicePlanID: newPlanId,
			},
			},
		}
		plan = &types.ServicePlan{
			ID:                planId,
			ServiceOfferingID: offeringId,
			CatalogName:       "large",
		}

		plans = &types.ServicePlans{
			ServicePlans: []types.ServicePlan{
				{ID: newPlanId,
					ServiceOfferingID: offeringId,
					CatalogName:       planName},
			},
		}
		if newInstanceName == "" {
			newInstanceName = instanceName
		}
		updatedIntance = &types.ServiceInstance{
			ID:            instanceId,
			Name:          newInstanceName,
			ServicePlanID: newPlanId,
			Parameters:    json.RawMessage(instanceParams),
		}
		if serviceInstanceId != "" {
			updatedIntance.ID = serviceInstanceId
		}

		operation := &types.Operation{
			State: "in progress",
		}
		client.StatusReturns(operation, nil)
		client.ListInstancesReturns(instances, nil)
		client.GetPlanByIDReturns(plan, nil)
		client.GetInstanceByIDReturns(instance, errGetInstance)
		client.ListPlansReturns(plans, nil)
		client.UpdateInstanceReturns(updatedIntance, location, nil)
		piCmd := command.Prepare(cmd.SmPrepare)
		piCmd.SetArgs(args)
		Expect(piCmd.Execute()).ToNot(HaveOccurred())

		return piCmd
	}

	validSyncUpdateExecution := func(args ...string) *cobra.Command {
		return validAsyncUpdateExecution("", append(args, "--mode", "sync")...)
	}

	invalidUpdateInstanceCommandExecution := func(args ...string) error {
		trCmd := command.Prepare(cmd.SmPrepare)
		trCmd.SetArgs(args)
		return trCmd.Execute()
	}

	Describe("valid request", func() {
		Context("async call", func() {
			It("should print status command", func() {
				validAsyncUpdateExecution("location",
					"myinstancename", "--plan", "small", "--new-name", "new-service-instance-name", "--instance-params", "{\"color\":\"red\"}")
				Expect(buffer.String()).To(ContainSubstring(`smctl status location`))
			})
		})

		Context("sync call", func() {
			It("should print object", func() {
				validSyncUpdateExecution("myinstancename", "--plan", "small", "--new-name", "new-service-instance-name", "--instance-params", "{\"color\":\"red\"}")
				Expect(command.instanceName).To(Equal("myinstancename"))
				Expect(command.instance.Name).To(Equal("new-service-instance-name"))
				Expect(command.instance.ServicePlanID).To(Equal(newPlanId))
				Expect(command.instance.Parameters).To(Equal(json.RawMessage("{\"color\":\"red\"}")))
				tableOutputExpected := updatedIntance.TableData().String()
				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})
		})
		Context("sync call by id", func() {
			It("should print object", func() {
				validSyncUpdateExecution("myinstancename", "--id", "serviceinstanceid", "--plan", "small", "--new-name", "new-service-instance-name", "--instance-params", "{\"color\":\"red\"}")
				Expect(command.instance.ID).To(Equal("serviceinstanceid"))
				Expect(command.instanceName).To(Equal("myinstancename"))
				Expect(command.instance.Name).To(Equal("new-service-instance-name"))
				Expect(command.instance.ServicePlanID).To(Equal(newPlanId))
				Expect(command.instance.Parameters).To(Equal(json.RawMessage("{\"color\":\"red\"}")))
				tableOutputExpected := updatedIntance.TableData().String()
				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})
		})

	})

	Describe("invalid request", func() {
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
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(cmd.NO_INSTANCES_FOUND, "instance-name")))
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
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(cmd.FOUND_TOO_MANY_INSTANCES, "instance-name","update")))
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

			Context("update instance", func() {
				BeforeEach(func() {
					plans = &types.ServicePlans{
						ServicePlans: []types.ServicePlan{
							{ID: "plandid",
								ServiceOfferingID: "offeringid",
								CatalogName:       "small"},
						},
					}
				})
				Context("With http response error from http client", func() {
					It("should return error's description", func() {
						body := ioutil.NopCloser(bytes.NewReader([]byte("HTTP response error")))
						expectedError := util.HandleResponseError(&http.Response{Body: body})
						client.UpdateInstanceReturns(nil, "", expectedError)
						err := invalidUpdateInstanceCommandExecution("instance-name", "plan", "small")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("HTTP response error"))
					})
				})

				Context("With invalid output format", func() {
					It("should return error", func() {
						invFormat := "invalid-format"
						err := invalidUpdateInstanceCommandExecution("validName", "--plan", "small", "--output", invFormat)
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(Equal("unknown output: " + invFormat))
					})
				})
			})

		})

	})
})
