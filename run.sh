#!/usr/bin/env bash

set -e
set -o pipefail

initdb -D /usr/lib/postgresql/data/
pg_ctl -D /usr/lib/postgresql/data/ -l /home/postgres/psql.log start
