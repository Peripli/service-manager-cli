package configuration

import (
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

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

	var configPath string

	BeforeSuite(func() {
		var err error
		configPath, err = filepath.Abs("test_config.json")
		if err != nil {
			panic(err)
		}
	})

	Describe("New SM Configuration", func() {
		Context("when cfg file is provided", func() {
			It("should save to this file and load the same client config", func() {
				viperEnv := viper.New()
				viperEnv.SetFs(afero.NewMemMapFs())
				configuration, err := NewSMConfiguration(viperEnv, configPath)
				configuration.Save(&smclient.ClientConfig{URL: "http://sm.com", User: "admin", Token: "token"})

				clientConfig, errLoad := configuration.Load()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(errLoad).ShouldNot(HaveOccurred())
				Expect(*clientConfig).To(Equal(smclient.ClientConfig{URL: "http://sm.com", User: "admin", Token: "token"}))
			})
		})
	})

})
