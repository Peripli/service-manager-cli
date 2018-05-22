package configuration

import (
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConfiguration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Configuration test", func() {

	Describe("New SM Configuration", func() {
		Context("when cfg file is provided", func() {
			It("should save to this file and load the same client config", func() {
				configuration, err := NewSMConfiguration("test_config.json")

				configuration.Save(&smclient.ClientConfig{URL: "http://sm.com", User: "admin", Token: "token"})
				clientConfig, errLoad := configuration.Load()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(errLoad).ShouldNot(HaveOccurred())
				Expect(*clientConfig).To(Equal(smclient.ClientConfig{URL: "http://sm.com", User: "admin", Token: "token"}))
			})
		})
	})

})
