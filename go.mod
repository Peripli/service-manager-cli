module github.com/Peripli/service-manager-cli

go 1.13

require (
	github.com/Peripli/service-manager v0.23.2
	//github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20210521184019-c5ad59b459ec // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/onrik/logrus v0.9.0 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.10.3
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.0
	github.com/tidwall/gjson v1.9.3
	github.com/tidwall/sjson v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/oauth2 v0.0.0-20210615190721-d04028783cf1
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace gopkg.in/fsnotify.v1 v1.4.9 => github.com/fsnotify/fsnotify v1.4.9
