module github.com/Peripli/service-manager-cli

go 1.13

require (
	github.com/Peripli/service-manager v0.18.7-0.20210121153144-b52d0698b1be
	github.com/antlr/antlr4 v0.0.0-20210114010855-d34d2e1c271a // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/onrik/logrus v0.8.0 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/afero v1.5.1
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/oauth2 v0.0.0-20210113205817-d3ed898aa8a3
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace gopkg.in/fsnotify.v1 v1.4.9 => github.com/fsnotify/fsnotify v1.4.9
