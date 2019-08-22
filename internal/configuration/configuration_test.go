package configuration

import (
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/Peripli/service-manager-cli/pkg/auth"
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
				fs := afero.NewMemMapFs()
				viperEnv := viper.New()
				viperEnv.SetFs(fs)
				configuration, err := NewSMConfiguration(viperEnv, configPath)

				timeNow, _ := time.Parse(time.RFC1123Z, time.Now().Format(time.RFC1123Z))
				settings := Settings{
					URL:  "http://sm.com",
					User: "admin",
					Token: auth.Token{
						AccessToken: "token",
						ExpiresIn:   timeNow,
					},
					AuthFlow: auth.PasswordGrant,
				}

				configuration.Save(&settings)

				viperEnv = viper.New()
				viperEnv.SetFs(fs)
				configuration, err = NewSMConfiguration(viperEnv, configPath)
				clientConfig, errLoad := configuration.Load()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(errLoad).ShouldNot(HaveOccurred())
				Expect(*clientConfig).To(Equal(settings))
			})
		})
	})

})
