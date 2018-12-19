package curl

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"testing"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
)

func TestCurlCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Curl command test", func() {

	var client *smclientfakes.FakeClient
	var command *Cmd
	var buffer *bytes.Buffer
	var fs afero.Fs

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	setCallReturns := func(status int, body []byte, err error) {
		client.CallReturns(fakeResponse(status, body), err)
	}

	assertLastCall := func(expectedMethod, expectedPath string, expectedBody []byte, expectedOutput string, expectedError error) {
		cmdArgs := []string{expectedPath, "-X", expectedMethod, "-d", string(expectedBody)}
		err := executeWithArgs(cmdArgs)
		if expectedError != nil {
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(expectedError.Error()))
			return
		}
		Expect(buffer.String()).To(Equal(expectedOutput))

		lastCallIndex := client.CallCallCount() - 1
		method, path, reader := client.CallArgsForCall(lastCallIndex)
		Expect(method).To(Equal(expectedMethod))
		Expect(path).To(Equal(expectedPath))
		if expectedMethod != http.MethodGet {
			body, err := ioutil.ReadAll(reader)
			Expect(err).To(BeNil())
			if expectedBody[0] != byte('@') {
				Expect(body).To(Equal(expectedBody))
			}
		}
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		fs = afero.NewMemMapFs()
		command = NewCurlCmd(context, fs)
	})

	Context("when curl with path only", func() {
		It("should do GET by default", func() {
			expectedOutput := `{"brokers":[]}`
			setCallReturns(200, []byte(expectedOutput), nil)
			err := executeWithArgs([]string{"/v1/service_brokers"})
			Expect(err).To(BeNil())
			Expect(buffer.String()).To(Equal(expectedOutput))
		})
	})

	Context("when curl with path, method and body", func() {
		It("should do GET method", func() {
			setCallReturns(200, nil, nil)
			assertLastCall(http.MethodGet, "/v1/service_brokers", nil, "", nil)
		})

		It("should do PATCH method", func() {
			body := []byte(`{"name": "test-broker"}`)
			setCallReturns(201, body, nil)
			assertLastCall(http.MethodPost, "/v1/service_brokers", body, `{"name": "test-broker"}`, nil)
		})

		Context("when body is file", func() {
			It("should do POST with file data", func() {
				filename, _ := filepath.Abs("test.txt")
				f, err := fs.Create(filename)
				Expect(err).To(BeNil())
				content := `{"name":"test"}`
				f.Write([]byte(content))

				setCallReturns(201, []byte(content), nil)
				assertLastCall(http.MethodPost, "/v1/service_brokers", []byte(`@test.txt`), content, nil)
			})
		})
	})

	Context("when call errors", func() {
		It("should handle error", func() {
			err := errors.New("problem during call")
			setCallReturns(0, nil, err)
			assertLastCall(http.MethodGet, "/v1/service_brokers", nil, "", err)
		})
	})
})

func fakeResponse(statusCode int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
}
