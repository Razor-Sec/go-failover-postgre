#!/bin/bash
GOOS=linux GOARCH=amd64 go build && ./golang_failover --host1 127.0.0.1 --host2 127.0.0.1 --port1 5434 --port2 5438 --user1 postgres --user2 postgres --password1 redhat123 --password2 redhat123 --localuser postgres --localpass redhat123


