package configuration

import (
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

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
				configuration, err := New(viperEnv, configPath)

				data := struct {
					Name string
					Age  int
				}{
					Name: "test",
					Age:  10,
				}
				configuration.Save("key", &data)

				var loadIn struct {
					Name string
					Age  int
				}
				errLoad := configuration.UnmarshalKey("key", &loadIn)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(errLoad).ShouldNot(HaveOccurred())
				Expect(loadIn).To(Equal(data))
			})
		})
	})

})
