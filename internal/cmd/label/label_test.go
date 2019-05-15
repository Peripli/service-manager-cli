package label

import (
	"bytes"
	"errors"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	resperrors "github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/Peripli/service-manager/pkg/query"
	"github.com/Peripli/service-manager/pkg/web"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestLabelCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Label Command test", func() {

	var client *smclientfakes.FakeClient
	var command *Cmd
	var buffer *bytes.Buffer
	var labelChanges *types.LabelChanges

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewLabelCmd(context)
	})

	validLabelExecution := func(args ...string) error {
		client.LabelReturns(nil)
		lc := command.Prepare(cmd.SmPrepare)
		lc.SetArgs(args)
		return lc.Execute()
	}

	invalidLabelExecution := func(args ...string) error {
		lc := command.Prepare(cmd.SmPrepare)
		lc.SetArgs(args)
		return lc.Execute()
	}

	Describe("Valid Invocation", func() {
		Context("with valid arguments provided", func() {
			It("should label resource successfully", func() {
				err := validLabelExecution("platform", "id", "add", "key", "--val", "val1", "--val", "val2", "--val", "val3")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(buffer.String()).To(ContainSubstring("Resource labeled successfully!"))
			})

			It("should pass arguments properly", func() {
				labelChanges = &types.LabelChanges{
					LabelChanges: []*query.LabelChange{{Key: "key", Operation: "add_values", Values: []string{"val1", "val2", "val3"}}},
				}
				err := validLabelExecution("platform", "id", "add-values", "key", "--val", "val1","--val", "val2","--val", "val3")

				Expect(err).ShouldNot(HaveOccurred())
				resourcePath, id, changes := client.LabelArgsForCall(0)
				Expect(resourcePath).To(Equal(web.PlatformsURL))
				Expect(id).To(Equal("id"))
				Expect(changes).To(Equal(labelChanges))
			})
		})
	})

	Describe("Invalid invocation", func() {
		Context("with less than 4 required arguments", func() {
			It("should return error", func() {
				err := invalidLabelExecution()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("resource type, id, operation, key and values are required"))
			})

			It("should not call SM", func() {
				err := invalidLabelExecution()
				c := client.LabelCallCount()
				Expect(err).Should(HaveOccurred())
				Expect(c).To(Equal(0))
			})
		})

		Context("with more than 4 arguments", func() {
			It("should return error", func() {
				err := invalidLabelExecution("platform", "id", "add-values", "key", "--val", "value", "redundant")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("too much arguments, in case you have whitespaces in some of the arguments consider enclosig it with single quotes"))
			})

			It("should not call SM", func() {
				err := invalidLabelExecution("platform", "id", "add-values", "key", "--val", "value", "redundant")
				c := client.LabelCallCount()
				Expect(err).Should(HaveOccurred())
				Expect(c).To(Equal(0))
			})
		})

		Context("with unknown resource provided", func() {
			It("should return error", func() {
				err := invalidLabelExecution("invalid resource", "id", "add", "key", "--val", "value")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown resource"))
			})
			It("should not call SM", func() {
				err := invalidLabelExecution("invalid resource", "id", "add", "key", "--val", "value")
				c := client.LabelCallCount()
				Expect(err).Should(HaveOccurred())
				Expect(c).To(Equal(0))
			})
		})

		Context("with unknown operation provided", func() {
			It("should return error", func() {
				err := invalidLabelExecution("platform", "id", "invalid", "key", "--val", "value")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown operation"))
			})
			It("should not call SM", func() {
				err := invalidLabelExecution("platform", "id", "invalid", "key", "--val", "value")
				c := client.LabelCallCount()
				Expect(err).Should(HaveOccurred())
				Expect(c).To(Equal(0))
			})
		})

		Context("with error from http client", func() {
			It("should return error", func() {
				expectedErr := errors.New("http client error")
				client.LabelReturns(expectedErr)

				err := invalidLabelExecution("platform", "id", "add", "key", "--val", "value")

				Expect(err).To(MatchError(expectedErr))
			})
		})

		Context("with http response error from http client", func() {
			It("should return error's description", func() {
				description := "HTTP response error"
				client.LabelReturns(resperrors.ResponseError{Description: description})
				err := invalidLabelExecution("platform", "id", "add", "key", "--val", "value")
				Expect(err).To(MatchError(description))

			})
		})

	})

})
