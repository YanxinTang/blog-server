#!/bin/bash

set -e

MIGRATIONS="migrations"

cd $PWD

./migrate create -ext sql -dir $MIGRATIONS -seq $1