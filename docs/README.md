# TEx-Ras

1. [TODO](#todo)

## TODO

* instance
  * list
  * validate name
  * register
  * instanceRootValidate
  * instance lock/unlock/validate string vs. interface
* instances to submodule
* init instance
* train instance
* api
  * (rasa|instance).Config !global
  * output encoders (json, yaml, raw string, spew)
  * instanceList
  * instanceRegister
  * rasaVersion
  * rasaExec
  * rasaInit
    * buffer output
  * train
    * buffer output
* run
  * rasa run instance
  * instance run foobar
  * instance run (all)
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
  * logging facility
    * direct logs there (conditional, flags for that)
    * direct rasa output there, like in `rasa run -- use-syslog --syslog-(address|port|protocol)
* proxy endpoint
* client/chat endpoint?

* api locking
  * when register/init/copy/destroy instance
  * when run instance (sync.WaitGroup?)
* set rasa library log levels az in <https://rasa.com/docs/rasa/command-line-interface#log-level>
* validate httpStaticRoot if httpStaticEnabled on start
* implement Logger.Printf(), set fasthttp.Server() Logger.
* self signed cert creation in init for proxy and main (check init.sh)
