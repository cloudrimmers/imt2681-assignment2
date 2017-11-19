#!/bin/sh

# $1 = PORT Rimbot
# $2 = ACCESS_TOKEN (Clara)
# $3 = PORT Currencyservice
# $4 = DBNAME
# $5 = DBURI
# $6 = CURRENCY_URI

#
# rimbot/.env
#
WORKFILE=./cmd/rimbot/.env
touch $WORKFILE
env="PORT=$1\nACCESS_TOKEN=$2\nCURRENCY_URI=$6"
echo $env > $WORKFILE


#
# currencyservice/.env
#
WORKFILE=./cmd/currencyservice/.env
touch $WORKFILE
env="PORT=$3\nMONGODB_NAME=$4\nMONGODB_URI=$5\n"
echo $env > $WORKFILE

#
# fixerworker/.env
#
WORKFILE=./cmd/fixerworker/.env
touch $WORKFILE
env="MONGODB_NAME=$4\nMONGODB_URI=$5\n"
echo $env > $WORKFILE



#
# Confidentely cat stuff
#
echo "\n------- rimbot .env --------------"
cat ./cmd/rimbot/.env

echo "------- currencyservice .env -----"
cat ./cmd/currencyservice/.env

echo "------- fixerworker .env ---------"
cat ./cmd/fixerworker/.env
