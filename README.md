# Web micro-services in Go

## Task

Test how web micro-services can be written/developed in Go. 

Features to check

- [x] provide RESTful API for a fake "Todo" service
- [x] deploy service to Cloud Foundry
- [x] provide WebSocket endpoint that reverses the string passed
- [x] use RabbitMQ with publish/subscribe; implement something like _when I `POST` at a REST endpoint, inform through a RabbitMQ queue on this ORM change_
- [ ] check out `gorm` for ORM support
- [ ] use PostgreSQL as a database backend
- [ ] deploy micro-service and link it with RabbitMQ and MySQL services to Cloud Foundry: so we change the used database backend when deploying - how do we do it?
- [ ] use OAuth2 and authenticate and authorize against _Cloud Foundry's_ UAA    


## Components used

On my machine I used

- Windows `10.0.17134.137`
- Erlang `10.0.1`
- RabbitMQ `3.7.7`
- Go `1.10.3 x64`
- cf CLI `cf version 6.36.2+18ceab10f.2018-05-16`
- PCF Dev `version 0.30.0 (CLI: 850ae45, OVA: 0.549.0)`
- PCF Dev has a `p-rabbitmq` service deployed under the name `rabbitmq`
- GoLand as IDE 


## AMQP integration
- everytime a `Todo` is added or updated using `POST /todo/{id}` end-point, this `Todo` is piped into a Go channel `chan Todo` called `todoChanged` 
- AMQP connection is established and queue `orm:todos` is created in the `go` routine `Connect2AMQPAndSetupQueue` that is executed at start
- in its body an endless loop waits for new `Todo` that arrive in the `todoChanged` `chan` and finally publishes them in the queue `orm:todos` upon arrival



## Lessons learned

- how `GOPATH` works
- how `dep` works to store dependencies and resolve them to a `vendor` directory - needed to deploy external dependencies when staging in Cloud Foundry
- how the `manifest.yml` needs to written (`env` e.g. needs to contain `GOPACKAGENAME`, `command` contains the name of the binary to execute)
- Cloud Foundry buildpacks
    - how to use a remote buildpack in `cf push` with `-b https://....` 
    - how to update a buildpack with `cf update-buildpack` in a PCF Dev deployment
- `for i,v := range slice` returns the index `i` - for updating the element - and a copy `v` of the element - that can not be changed!  
- we can write "execute an operation with error result" and "check for actual error" as one-liner in go using `if err := call(); err != nil {...}`
- WebSocket can simply be tested out of any recent browser:
    - open the URL were our service runs, e.g. `localhost:8080`
    - press `F12` to open the developer console
    - enter `
var ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = function(data) { console.log(data); }    
`
    - send stuff entering `ws.send("שלום")` and check the console for the response
- if a REST call changed a `Todo` we can not use this `Todo` directly in the `AMQP` part as this will cause an _Access violation_ ; use a _Go routine_ instead to pipe the `Todo`
- if a `rabbitmq` service is bound to an application in Cloud Foundry, the `VCAP_SERVICES` environment variable will contain a service definition that can be identified by its `amqp` tag and whose `credentials.uri` field contains the full `amqp://...` URI
