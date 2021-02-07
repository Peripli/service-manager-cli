package test

import (
	"encoding/json"
	cliquery "github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/Peripli/service-manager/pkg/web"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Plans test", func() {
	plan_supports_any_env := &types.ServicePlan{
		ID:                "plan_supports_any_env",
		Name:              "plan_supports_any_env",
		Description:       "no metadata plan", // supports any env
		ServiceOfferingID: "of1",
	}
	plan_supports_cf_env1 := &types.ServicePlan{
		ID:                "plan_supports_cf_env1",
		Name:              "plan_supports_cf_env1",
		Description:       "plan with cf support",
		ServiceOfferingID: "of2",
		Metadata:          []byte("{\"supportedPlatforms\":[\"cloudfoundry\"]}"),
	}
	plan_supports_cf_env2 := &types.ServicePlan{
		ID:                "plan_supports_cf_env2",
		Name:              "plan_supports_cf_env2",
		Description:       "plan with cf support",
		ServiceOfferingID: "of3",
		Metadata:          []byte("{\"supportedPlatforms\":[\"cloudfoundry\"]}"),
	}
	plan_supports_k8s := &types.ServicePlan{
		ID:                "plan_supports_k8s",
		Name:              "plan_supports_k8s",
		Description:       "plan with k8s support",
		ServiceOfferingID: "of3",
		Metadata:          []byte("{\"supportedPlatforms\":[\"kubernetes\"]}"),
	}
	plan_supports_k8s_and_cf := &types.ServicePlan{
		ID:                "plan_supports_k8s_and_cf",
		Name:              "plan_supports_k8s_and_cf",
		Description:       "plan with cf & k8s support",
		ServiceOfferingID: "of4",
		Metadata:          []byte("{\"supportedPlatforms\":[\"kubernetes\",\"cloudfoundry\"]}"),
	}
	plan_supports_test_env := &types.ServicePlan{
		ID:                "plan_supports_test_env",
		Name:              "plan_supports_test_env",
		Description:       "plan with test env support",
		ServiceOfferingID: "of5",
		Metadata:          []byte("{\"supportedPlatforms\":[\"test-env\"]}"),
	}

	Context("when plans with explicit and any supported platform provided", func() {
		BeforeEach(func() {
			handlerDetails = []HandlerDetails{
				{
					Method:             http.MethodGet,
					Path:               web.ServicePlansURL,
					ResponseStatusCode: http.StatusOK,
					ResponseBodyProvider: func(req *http.Request) []byte {
						response, _ := json.Marshal(types.ServicePlans{
							ServicePlans: []types.ServicePlan{
								*plan_supports_any_env,
								*plan_supports_cf_env1,
								*plan_supports_cf_env2,
								*plan_supports_k8s,
								*plan_supports_k8s_and_cf,
								*plan_supports_test_env,
							},
						})
						return response
					}},
			}
		})

		It("should return all plans if any environment requested", func() {
			result, err := client.ListPlans(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "any",
			})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(result.ServicePlans).To(HaveLen(6))
			Expect(result.ServicePlans[0]).To(Equal(*plan_supports_any_env))
			Expect(result.ServicePlans[1]).To(Equal(*plan_supports_cf_env1))
			Expect(result.ServicePlans[2]).To(Equal(*plan_supports_cf_env2))
			Expect(result.ServicePlans[3]).To(Equal(*plan_supports_k8s))
			Expect(result.ServicePlans[4]).To(Equal(*plan_supports_k8s_and_cf))
			Expect(result.ServicePlans[5]).To(Equal(*plan_supports_test_env))
		})
		It("should return cloudfoundry plans if cloudfoundry environment requested", func() {
			result, err := client.ListPlans(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "cloudfoundry",
			})

			Expect(err).ShouldNot(HaveOccurred())

			Expect(result.ServicePlans).To(HaveLen(4))
			Expect(result.ServicePlans[0]).To(Equal(*plan_supports_any_env))
			Expect(result.ServicePlans[1]).To(Equal(*plan_supports_cf_env1))
			Expect(result.ServicePlans[2]).To(Equal(*plan_supports_cf_env2))
			Expect(result.ServicePlans[3]).To(Equal(*plan_supports_k8s_and_cf))
		})
		It("should return kubernetes plans if kubernetes environment requested", func() {
			result, err := client.ListPlans(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "kubernetes",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServicePlans).To(HaveLen(3))
			Expect(result.ServicePlans[0]).To(Equal(*plan_supports_any_env))
			Expect(result.ServicePlans[1]).To(Equal(*plan_supports_k8s))
			Expect(result.ServicePlans[2]).To(Equal(*plan_supports_k8s_and_cf))
		})
		It("should return test-env if test-env environment requested", func() {
			result, err := client.ListPlans(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "test-env",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServicePlans).To(HaveLen(2))
			Expect(result.ServicePlans[0]).To(Equal(*plan_supports_any_env))
			Expect(result.ServicePlans[1]).To(Equal(*plan_supports_test_env))
		})
		It("should return plans with any env support if non existing environment requested", func() {
			result, err := client.ListPlans(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "non-existing",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServicePlans).To(HaveLen(1))
			Expect(result.ServicePlans[0]).To(Equal(*plan_supports_any_env))
		})
	})
	Context("when plans with explicitly supported environments provided only (no plans with any platform support)", func() {
		BeforeEach(func() {
			handlerDetails = []HandlerDetails{
				{
					Method:             http.MethodGet,
					Path:               web.ServicePlansURL,
					ResponseStatusCode: http.StatusOK,
					ResponseBodyProvider: func(req *http.Request) []byte {
						response, _ := json.Marshal(types.ServicePlans{
							ServicePlans: []types.ServicePlan{
								*plan_supports_cf_env1,
								*plan_supports_cf_env2,
								*plan_supports_k8s,
								*plan_supports_k8s_and_cf,
								*plan_supports_test_env,
							},
						})
						return response
					}},
			}
		})
		It("should return empty plan response if non existing environment requested", func() {
			result, err := client.ListPlans(&cliquery.Parameters{
				GeneralParams: []string{"key=value"},
				Environment:   "non-existing",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result.ServicePlans).To(HaveLen(0))
		})
	})

})
