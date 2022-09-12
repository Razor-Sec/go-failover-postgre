#!/bin/bash
go build && ./golang_failover --host1 10.8.60.191 --host2 10.8.60.192 --port1 5434 --port2 5438 --user1 postgres --user2 postgres --password1 redhat123 --password2 redhat123 --localuser postgres --localpass redhat123 --localdata "/usr/pgsql-12/bin/pg_ctl" 


