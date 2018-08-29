# eosforce-db
sync eosforce to db

## Usage

install:

```shell
go get -u -v github.com/fanyang1988/eosforce-db
cd $GOPATH/src/github.com/fanyang1988/eosforce-db
go get -u -v ./...
cd eosc2db 
go build
```

use `--help` to show usage:

```shell
./eosc2db --help
eosc2db is tool to dump data from eosforce to db
Source code is available at: https://github.com/fanyang1988/eosforce-db

Usage:
  eosc2db [command]

Available Commands:
  help        Help about any command
  listen      listen eosforce data to log
  watch       watch eosforce then sync data to db

Flags:
      --db string           db (default "test1")
      --db-address string   db address (default "127.0.0.1:5432")
      --db-passwd string    db password (default "123456")
      --db-user string      db user (default "pgfy")
  -h, --help                help for eosc2db
      --log-cfg string      seelog cfg file (default "./log_cfg.xml")

Use "eosc2db [command] --help" for more information about a command.

```

**Listen eosforce Actions to Log**

usage: 

```
eosc2db listen [eosforce http api url] [eosforce p2p address]
```


example:

```shell
./eosc2db listen  http://p1.eosforce.cn:8888 111.231.190.233:8001
```
