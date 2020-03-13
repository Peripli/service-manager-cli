package broker

import (
	"encoding/json"
	"errors"
	"testing"

	"gopkg.in/yaml.v2"

	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestUpdateBrokerCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Update broker command test", func() {

	var client *smclientfakes.FakeClient
	var command *UpdateBrokerCmd
	var buffer *bytes.Buffer
	var broker types.Broker

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewUpdateBrokerCmd(context)
	})

	validAsyncUpdateBrokerExecution := func(location string, args ...string) error {
		broker = types.Broker{
			Name:        "broker1",
			ID:          "id",
			URL:         "http://broker1.com",
			Description: "description",
		}
		brokers := &types.Brokers{Brokers: []types.Broker{broker}}
		client.ListBrokersReturns(brokers, nil)
		_ = json.Unmarshal([]byte(args[1]), broker)
		client.UpdateBrokerReturns(&broker, location, nil)
		ubCmd := command.Prepare(cmd.SmPrepare)
		ubCmd.SetArgs(args)
		return ubCmd.Execute()
	}

	validSyncUpdateBrokerExecution := func(args ...string) error {
		return validAsyncUpdateBrokerExecution("", args...)
	}

	invalidUpdateBrokerExecution := func(args ...string) error {
		ubCmd := command.Prepare(cmd.SmPrepare)
		ubCmd.SetArgs(args)
		return ubCmd.Execute()
	}

	Describe("Valid request", func() {
		Context("With necessary arguments provided", func() {
			It("broker should be updated", func() {

				err := validSyncUpdateBrokerExecution("broker1", `{"description":"newDescription"}`)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(broker.TableData().String()))
			})

			It("should print location when updated asynchronously", func() {
				Expect(validAsyncUpdateBrokerExecution("location", "broker1", `{"description":"newDescription"}`)).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring(`smctl poll location`))
			})

			It("argument values should be as expected", func() {
				err := validSyncUpdateBrokerExecution("broker1", `{"description":"newDescription"}`)

				id, broker, _ := client.UpdateBrokerArgsForCall(0)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id).To(Equal("id"))
				Expect(broker).To(Equal(&types.Broker{Description: "newDescription"}))
			})
		})

		Context("With json format flag", func() {
			It("should be printed in json format", func() {
				err := validSyncUpdateBrokerExecution("broker1", `{"description":"newDescription"}`, "--output", "json")

				jsonByte, _ := json.MarshalIndent(broker, "", "  ")
				jsonOutputExpected := string(jsonByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(jsonOutputExpected))
			})
		})

		Context("With yaml format flag", func() {
			It("should be printed in yaml format", func() {
				err := validSyncUpdateBrokerExecution("broker1", `{"description":"newDescription"}`, "--output", "yaml")

				yamlByte, _ := yaml.Marshal(broker)
				yamlOutputExpected := string(yamlByte) + "\n"

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(Equal(yamlOutputExpected))
			})
		})

		Context("With generic parameter flag provided", func() {
			It("should pass it to SM", func() {
				err := validSyncUpdateBrokerExecution("broker1", `{"description":"newDescription"}`, "--param", "paramKey=paramValue")
				Expect(err).ShouldNot(HaveOccurred())

				_, _, args := client.UpdateBrokerArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("paramKey=paramValue", "async=false"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})

		Context("With async flag", func() {
			It("should pass it to SM", func() {
				err := validSyncUpdateBrokerExecution("broker1", `{"description":"newDescription"}`, "--async")
				Expect(err).ShouldNot(HaveOccurred())

				_, _, args := client.UpdateBrokerArgsForCall(0)

				Expect(args.GeneralParams).To(ConsistOf("async=true"))
				Expect(args.FieldQuery).To(BeEmpty())
				Expect(args.LabelQuery).To(BeEmpty())
			})
		})
	})

	Describe("Invalid request", func() {
		Context("With missing arguments", func() {
			It("Should return error missing name", func() {
				err := invalidUpdateBrokerExecution([]string{}...)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("[name] is required"))
			})
			It("Should return error missing json", func() {
				err := invalidUpdateBrokerExecution("broker1")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nothing to update. Broker JSON is not provided"))
			})
		})
	})

	Context("When non existing broker updated", func() {
		It("should return error", func() {
			client.ListBrokersReturns(&types.Brokers{}, nil)

			err := invalidUpdateBrokerExecution("broker1", `{"description":"newDescription"}`)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("broker with name broker1 not found"))
		})
	})

	Context("With error from http client", func() {
		It("Should return error", func() {
			expectedErr := errors.New("http client error")
			brokers := &types.Brokers{Brokers: []types.Broker{broker}}
			client.ListBrokersReturns(brokers, nil)
			client.UpdateBrokerReturns(nil, "", expectedErr)

			err := invalidUpdateBrokerExecution("broker1", `{"description":"newDescription"}`)

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expectedErr.Error()))
		})
	})

	Context("With invalid output format", func() {
		It("should return error", func() {
			invFormat := "invalid-format"
			err := invalidUpdateBrokerExecution("broker1", `{"description":"newDescription"}`, "--output", invFormat)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown output: " + invFormat))
		})
	})

	Context("when json is invalid", func() {
		It("should throw error", func() {
			err := invalidUpdateBrokerExecution("broker1", "{name: none}")

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("broker JSON is invalid"))
		})
	})

})
