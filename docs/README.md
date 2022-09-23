# TEx-Ras

1. [TODO](#todo)

## TODO

* instance
  * list
  * validate name
  * register
  * cd instance dir
* init instance
* train instance
* api
  * rasa.Config !global
  * version
  * exec
  * init
  * train
* run
  * instance
  * all
* copy instance
* destroy instance
* update instance config/data
  * enable/disable instance
  * nlu mode true/false
  * training data
  * domain data
  * config
  * actions
  * channel connectors
  * (markers)
* implement external
  * action server management
  * tracker store
  * event broker
  * model store
  * lock store
  * nlg
  * logging facility (with syslogd docker or syslog listener in the mgmt process)
* proxy endpoint
* client/chat endpoint?

* locking when register/init/copy/destroy instance
* set rasa library log levels az in <https://rasa.com/docs/rasa/command-line-interface#log-level>
* validate httpStaticRoot if httpStaticEnabled on start
* implement Logger.Printf(), set fasthttp.Server() Logger.
* self signed cert creation in init for proxy and main (check init.sh)
