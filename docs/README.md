# TEx-Ras

1. [TODO](#todo)

## TODO

* init instance
* copy instance
* destroy instance
* list instances
* api for init/destroy/list
* update instance config/data
  * enable/disable instance
  * nlu mode true/false
  * training data
  * domain data
  * config
  * actions
  * channel connectors
  * (markers)
* train instance
* start stop instance(s)
  * specific instance
  * all instances
* api
  * start enabled
  * rasa api gw
* client/chat endpoint
* implement external
  * action server management
  * tracker store
  * event broker
  * model store
  * lock store
  * nlg
  * logging facility (with syslogd docker or syslog listener in the mgmt process)

* get agents dir form cli/env
* set rasa specific log levels az in <https://rasa.com/docs/rasa/command-line-interface#log-level>
* validate httpStaticRoot if httpStaticEnabled on start
* implement Logger.Printf(), set fasthttp.Server() Logger.
* self signed cert creation in init for proxy and main (check init.sh)
