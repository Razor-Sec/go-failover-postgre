# go-failover-postgre
Tested on RHEL BASED & Postgre 12

# example usage :
```bash
./golang_failover --host1 10.8.60.191 --host2 10.8.60.192 --port1 5434 --port2 5434 --user1 postgres --user2 postgres --password1 devopskeren --password2 devopskeren --localuser postgres --localpass devopskeren --localdata "/var/lib/pgsql/12/data" --localpg "/usr/pgsql-12/bin/pg_ctl" --localport 5434 --localhost localhost
```
