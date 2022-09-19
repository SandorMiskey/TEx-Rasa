# TEx-Rasa/mgmt

1. [TODO](#todo)

## TODO

* get agents dir form cli/env
* create instance/domain
* edit/show instance/domain data
  * training data
  * domain data
  * config
  * actions
  * channel connectors
  * (markers)
* training endpoint
* start/stop agent process
  * nlu only mode
  * (multiple agent process)
* client/chat endpoint
* create/destroy instance or agent
* implement external
  * action server management
  * tracker store
  * event broker
  * model store
  * lock store
  * nlg
  * logging facility (with syslogd docker or syslog listener in the mgmt process)
* validate httpStaticRoot if httpStaticEnabled on start
* implement Logger.Printf(), set fasthttp.Server() Logger.
* self signed cert creation in init for proxy and main (check init.sh)
