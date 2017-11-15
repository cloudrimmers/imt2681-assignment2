#!/bin/sh

# $1 = PORT Rimbot
# $2 = PORT Currencyservice
# $3 = DBNAME
# $4 = DBURI

#
# rimbot/.env
#
WORKFILE=./cmd/rimbot/.env
touch $WORKFILE
env="PORT=$1\n"
echo $env > $WORKFILE


#
# currencyservice/.env
#
WORKFILE=./cmd/currencyservice/.env
touch $WORKFILE
env="PORT=$2\nMONGODB_NAME=$3\nMONGODB_URI=$4\n" 
echo $env > $WORKFILE

#
# fixerworker/.env
#
WORKFILE=./cmd/fixerworker/.env
touch $WORKFILE
env="MONGODB_NAME=$3\nMONGODB_URI=$4\n"
echo $env > $WORKFILE
