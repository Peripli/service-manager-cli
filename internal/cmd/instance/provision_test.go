package instance

import (
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"

	"bytes"
	"errors"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/spf13/cobra"
)

const OfferingID = "offering_id"
const PlanID = "plan_id"

var _ = Describe("Provision Command test", func() {
	var client *smclientfakes.FakeClient
	var command *ProvisionCmd
	var buffer *bytes.Buffer

	var offerings *types.ServiceOfferings
	var plans *types.ServicePlans
	var instance *types.ServiceInstance

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewProvisionCmd(context)
	})

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

	invalidProvisionCommandExecution := func(args ...string) error {
		rpcCmd := command.Prepare(cmd.SmPrepare)
		rpcCmd.SetArgs(args)
		return rpcCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("should be registered synchronously", func() {
				validSyncProvisionExecution("instance-name", "offering-name", "plan-name")

				tableOutputExpected := instance.TableData().String()

				Expect(buffer.String()).To(ContainSubstring(tableOutputExpected))
			})

			It("should print location when registered asynchronously", func() {
				validAsyncProvisionExecution("location", "instance-name", "offering-name", "plan-name")

				Expect(buffer.String()).To(ContainSubstring(`smctl status location`))
			})

			It("Argument values should be as expected", func() {
				validSyncProvisionExecution("instance-name", "offering-name", "plan-name")

				Expect(command.instance.Name).To(Equal("instance-name"))
				Expect(command.instance.ServiceID).To(Equal(OfferingID))
				Expect(command.instance.ServicePlanID).To(Equal(PlanID))
			})
		})

		Context("With json output flag", func() {
			It("should be printed in json output format", func() {
				validSyncProvisionExecution("instance-name", "offering-name", "plan-name", "--output", "json")

				jsonByte, _ := json.MarshalIndent(instance, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml output flag", func() {
			It("should be printed in yaml output format", func() {
				validSyncProvisionExecution("instance-name", "offering-name", "plan-name", "--output", "yaml")

				yamlByte, _ := yaml.Marshal(instance)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})

		Context("With generic param flag", func() {
			It("should pass it to SM", func() {
				validSyncProvisionExecution("instance-name", "offering-name", "plan-name", "--param", "paramKey=paramValue")

				_, args := client.ProvisionArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("paramKey=paramValue", "async=false"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})

		Context("With sync flag", func() {
			It("should pass it to SM", func() {
				validSyncProvisionExecution("instance-name", "offering-name", "plan-name")

				_, args := client.ProvisionArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("async=false"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})
	})

	Describe("Invalid request", func() {
		var offerings *types.ServiceOfferings
		var plans *types.ServicePlans
		var brokers *types.Brokers

		BeforeEach(func() {
			offerings = &types.ServiceOfferings{
				ServiceOfferings: []types.ServiceOffering{
					{ID: OfferingID, Name: "offering-name", BrokerID: "broker-id"},
				},
			}
			plans = &types.ServicePlans{
				ServicePlans: []types.ServicePlan{
					{ID: PlanID, Name: "plan-name"},
				},
			}
			brokers = &types.Brokers{
				Brokers: []types.Broker{
					{ID: "broker-id", Name: "broker-name"},
				},
			}
		})

		JustBeforeEach(func() {
			client.ListOfferingsReturns(offerings, nil)
			client.ListPlansReturns(plans, nil)
			client.ListBrokersReturns(brokers, nil)
		})

		Context("With not enough arguments provided", func() {
			It("should return error", func() {
				err := invalidProvisionCommandExecution("validName", "offering-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name, offering and plan are required"))
			})
		})

		Context("When offering not found", func() {
			BeforeEach(func() {
				offerings = &types.ServiceOfferings{}
			})

			It("should return error", func() {
				err := invalidProvisionCommandExecution("validName", "offering-name", "plan-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(SatisfyAll(ContainSubstring("service offering with name"), ContainSubstring("not found")))
			})
		})

		Context("When more than one offering with same name found", func() {
			BeforeEach(func() {
				offerings.ServiceOfferings = append(offerings.ServiceOfferings, types.ServiceOffering{ID: OfferingID, Name: "offering-name"})
			})

			Context("and broker name is not provided", func() {
				It("should return error", func() {
					err := invalidProvisionCommandExecution("validName", "offering-name", "plan-name")

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(SatisfyAll(ContainSubstring("more than one service offering with name"), ContainSubstring("found. Use -b flag to specify broker name")))
				})
			})

			Context("and broker name is provided", func() {
				It("should distinguish between offerings", func() {
					client.ProvisionReturns(&types.ServiceInstance{}, "", nil)

					err := invalidProvisionCommandExecution("validName", "offering-name", "plan-name", "-b", "broker-name")

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		Context("With error from http client", func() {
			It("should return error", func() {
				client.ProvisionReturns(nil, "", errors.New("Http Client Error"))

				err := invalidProvisionCommandExecution("validName", "offering-name", "plan-name")

				Expect(err).To(MatchError("Http Client Error"))
			})
		})

		Context("With http response error from http client", func() {
			It("should return error's description", func() {
				body := ioutil.NopCloser(bytes.NewReader([]byte("HTTP response error")))
				expectedError := util.HandleResponseError(&http.Response{Body: body})
				client.ProvisionReturns(nil, "", expectedError)

				err := invalidProvisionCommandExecution("validName", "offering-name", "plan-name")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("HTTP response error"))
			})
		})

		Context("With invalid output format", func() {
			It("should return error", func() {
				invFormat := "invalid-format"
				err := invalidProvisionCommandExecution("validName", "offering-name", "plan-name", "--output", invFormat)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("unknown output: " + invFormat))
			})
		})
	})
})
