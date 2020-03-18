package test

import (
	"encoding/json"
	"fmt"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"net/http"
	"net/http/httptest"
	"testing"

	cliquery "github.com/Peripli/service-manager-cli/pkg/query"

	"github.com/Peripli/service-manager-cli/pkg/types"
	smtypes "github.com/Peripli/service-manager/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeAuthClient struct {
	AccessToken string
	requestURI  string
}

func (c *FakeAuthClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	c.requestURI = req.URL.RequestURI()
	return http.DefaultClient.Do(req)
}

type HandlerDetails struct {
	Method             string
	Path               string
	ResponseBody       []byte
	ResponseStatusCode int
	Headers            map[string]string
}

func TestSMClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var (
	client         smclient.Client
	handlerDetails []HandlerDetails
	validToken     = "valid-token"
	invalidToken   = "invalid-token"
	smServer       *httptest.Server
	fakeAuthClient *FakeAuthClient
	params         *cliquery.Parameters

	broker = &types.Broker{
		ID:          "broker-id",
		Name:        "test-broker",
		URL:         "http://test-url.com",
		Credentials: &types.Credentials{Basic: types.Basic{User: "test user", Password: "test password"}},
	}

	initialOffering = &types.ServiceOffering{
		ID:          "offeringID",
		Name:        "initial-offering",
		Description: "Some description",
		BrokerID:    "id",
	}

	plan = &types.ServicePlan{
		ID:                "planID",
		Name:              "plan-1",
		Description:       "Sample Plan",
		ServiceOfferingID: "offeringID",
	}

	resultOffering = &types.ServiceOffering{
		ID:          "offeringID",
		Name:        "initial-offering",
		Description: "Some description",
		Plans:       []types.ServicePlan{*plan},
		BrokerID:    "id",
		BrokerName:  "test-broker",
	}

	instance = &types.ServiceInstance{
		ID:            "instanceID",
		Name:          "instance1",
		ServicePlanID: "service_plan_id",
		PlatformID:    "platform_id",
		Context:       json.RawMessage("{}"),
	}

	binding = &types.ServiceBinding{
		ID:                "instanceID",
		Name:              "instance1",
		ServiceInstanceID: "service_instance_id",
	}

	operation = &types.Operation{
		ID:         "operation-id",
		Type:       "create",
		State:      "failed",
		ResourceID: "broker-id",
	}

	labelChanges = &types.LabelChanges{
		LabelChanges: []*smtypes.LabelChange{
			{Key: "key", Operation: smtypes.LabelOperation("add"), Values: []string{"val1", "val2"}},
		},
	}
)

var createSMHandler = func() http.Handler {
	mux := http.NewServeMux()
	for i := range handlerDetails {
		v := handlerDetails[i]
		mux.HandleFunc(v.Path, func(response http.ResponseWriter, req *http.Request) {
			if v.Method != req.Method {
				return
			}
			for key, value := range v.Headers {
				response.Header().Set(key, value)
			}
			authorization := req.Header.Get("Authorization")
			if authorization != "Bearer "+validToken {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write([]byte(""))
				return
			}
			response.WriteHeader(v.ResponseStatusCode)
			response.Write(v.ResponseBody)
		})
	}
	return mux
}

var verifyErrorMsg = func(errorMsg, path string, body []byte, statusCode int) {
	Expect(errorMsg).To(ContainSubstring(smclient.BuildURL(smServer.URL+path, params)))
	Expect(errorMsg).To(ContainSubstring(string(body)))
	Expect(errorMsg).To(ContainSubstring(fmt.Sprintf("StatusCode: %d", statusCode)))
}

var _ = BeforeEach(func() {
	params = &cliquery.Parameters{
		GeneralParams: []string{"key=value"},
	}
})

var _ = AfterEach(func() {
	Expect(fakeAuthClient.requestURI).Should(ContainSubstring("key=value"), fmt.Sprintf("Request URI %s should contain ?key=value", fakeAuthClient.requestURI))
})

var _ = JustBeforeEach(func() {
	smServer = httptest.NewServer(createSMHandler())
	fakeAuthClient = &FakeAuthClient{AccessToken: validToken}
	client = smclient.NewClient(fakeAuthClient, smServer.URL)
})
