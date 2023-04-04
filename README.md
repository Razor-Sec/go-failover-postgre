# go-failover-postgre
Tested on RHEL BASED & Postgre 12

# example usage :
```bash
./golang_failover --host1 10.8.60.191 --host2 10.8.60.192 --port1 5434 --port2 5434 --user1 postgres --user2 postgres --password1 <pg password> --password2 <pg password> --localuser postgres --localpass <pg password> --localdata "/var/lib/pgsql/12/data" --localpg "/usr/pgsql-12/bin/pg_ctl" --localport 5434 --localhost localhost
```

# run as systemd
```bash
[Unit]
Description=GoFailover Service
After=network.target

[Service]
Type=simple
User=postgres
Group=postgres
PermissionsStartOnly=true
WorkingDirectory=/opt/golang_failover
ExecStart=/opt/golang_failover/golang_failover --host1 10.8.60.191 --host2 10.8.60.192 --port1 5432 --port2 5432 --user1 postgres --user2 postgres --password1 redhat123 --password2 redhat123 --localuser postgres --localpass redhat123 --localdata "/usr/pgsql-12/bin/pg_ctl"
ExecStartPre=/bin/chown postgres /opt/golang_failover

[Install]
WantedBy=multi-user.target
```
