#!/bin/sh

#
# --help
#
case $1 in
	"--help")
		echo '
------ ARGUMENTS ------
# $1 = ACCESS_TOKEN (Clara)
# $2 = RIMBOT_PORT
# $3 = CURRENCY_URI
# $4 = CURRENCY_PORT
# $5 = MONGODB_URI
# $6 = MONGODB_NAME
'
		exit 0
		;;
esac

#
# rimbot/.env
#
WORKFILE=./cmd/rimbot/.env
env="
ACCESS_TOKEN=$1\n
PORT=$2		\n
CURRENCY_URI=$3\n"
echo $env > $WORKFILE


#
# currencyservice/.env
#
WORKFILE=./cmd/currencyservice/.env
env="
PORT=$4\n
MONGODB_URI=$5\n
MONGODB_NAME=$6\n"
echo $env > $WORKFILE

#
# fixerworker/.env
#
WORKFILE=./cmd/fixerworker/.env
env="
MONGODB_URI=$5\n
MONGODB_NAME=$6\n"
echo $env > $WORKFILE


#
# .env
#
WORKFILE=.env
env="
CURRENCY_PORT=$4\n
RIMBOT_PORT=$2\n
MONGO_PORT=27017\n"
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

echo "------- global .env ---------"
cat .env
