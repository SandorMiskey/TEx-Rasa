# TEx-Ras

1. [TODO](#todo)

## TODO

* train instance
  * lock
* api
  * storage/engine/instances/
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
  * api
* copy instance
* destroy instance
* update instance config/data
  * enable/disable instance
  * nlu mode true/false
  * change port
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
    * direct rasa output there, like in `rasa run --use-syslog --syslog-(address|port|protocol)
* proxy endpoint
* client/chat endpoint?

---

* Unix syslog delivery error
* api locking
  * when register/init/copy/destroy instance
  * when run instance (sync.WaitGroup?)
* storage to submodule
* log levels everywhere
* set rasa library log levels as in <https://rasa.com/docs/rasa/command-line-interface#log-level>
* validate httpStaticRoot if httpStaticEnabled on start
* self signed cert creation in init for proxy and main (check init.sh)
