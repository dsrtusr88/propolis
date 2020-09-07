module gitlab.com/passelecasque/propolis

go 1.13

require (
	github.com/dgraph-io/badger/v2 v2.2007.1
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/go-chat-bot/bot v0.0.0-20200527181414-ef71c72a524a
	github.com/gorilla/mux v1.6.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/thoj/go-ircevent v0.0.0-20190807115034-8e7ce4b5a1eb
	gitlab.com/catastrophic/assistance v0.38.7
	gitlab.com/passelecasque/obstruction v0.12.2
	gopkg.in/yaml.v2 v2.2.8
)

// replace gitlab.com/catastrophic/assistance => ../../catastrophic/assistance
