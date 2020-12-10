module github.com/Peripli/service-manager-cli

go 1.13

require (
	github.com/InVisionApp/go-health v2.1.0+incompatible // indirect
	github.com/InVisionApp/go-logger v1.0.1
	github.com/Peripli/service-manager v0.16.11
	github.com/antlr/antlr4 v0.0.0-20200119161855-7a3f40bc341d
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/uuid v3.1.0+incompatible
	github.com/golang/protobuf v1.3.4
	github.com/hashicorp/hcl v1.0.0
	github.com/hpcloud/tail v1.0.0
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onrik/logrus v0.4.1
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/pelletier/go-toml v1.6.0
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/subosito/gotenv v1.2.0
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/match v1.0.1
	github.com/tidwall/pretty v1.0.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527
	golang.org/x/text v0.3.2
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/appengine v1.6.5
	gopkg.in/ini.v1 v1.54.0
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 v2.2.8
)

replace gopkg.in/fsnotify.v1 v1.4.9 => github.com/fsnotify/fsnotify v1.4.9
