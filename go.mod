module gitlab.com/passelecasque/propolis

go 1.13

require (
	github.com/anacrolix/log v0.3.1-0.20191001111012-13cede988bcd
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	gitlab.com/catastrophic/assistance v0.36.2
	gitlab.com/passelecasque/obstruction v0.12.2
)

//replace gitlab.com/catastrophic/assistance => ../../catastrophic/assistance
replace gitlab.com/passelecasque/propolis-bot => ../propolis-bot
