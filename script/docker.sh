#!/bin/sh

#
# RIMBOT | DOCKER BUILD AND RUN
#
rimbot_run(){
	echo
	echo Running docker image of rimbot:latest ...
	echo
	docker run -it \
			   --rm \
			   -p :5000:5000 \
			   --name rimbot \
			   --env-file cmd/rimbot/.env \
			   rimbot:latest
}

rimbot_build(){
	echo 
	echo  Building docker image of rimbot:latest ..
	echo
	docker build --tag rimbot:latest \
			     --file ./cmd/rimbot/Dockerfile \
		         .
}

#
# CURRENCYSERVICE | DOCKER BUILD AND RUN
#
currencyservice_run(){
	echo
	echo Running docker image of currencyservice:latest ...
	echo
	docker run -it \
			   --rm \
			   -p :5001:5001 \
			   --name currencyservice \
			   --env-file cmd/currencyservice/.env \
			   currencyservice:latest
}

currencyservice_build(){
	echo 
	echo  Building docker image of currencyservice:latest ..
	echo
	docker build --tag currencyservice:latest \
			     --file ./cmd/currencyservice/Dockerfile \
		         .
}

#
# FIXERWORKER | DOCKER BUILD AND RUN
#
fixerworker_run(){
	echo
	echo Running docker image of fixerworker:latest ...
	echo
	docker run -it \
			   --rm \
			   -p :5002:5002 \
			   --name fixerworker \
			   --env-file cmd/fixerworker/.env \
			   fixerworker:latest
}

fixerworker_build(){
	echo 
	echo  Building docker image of fixerworker:latest ..
	echo
	docker build --tag fixerworker:latest \
			     --file ./cmd/fixerworker/Dockerfile \
		         .
}


#
# MONGO | DOCKER RUN
#
mongo_run() {
	echo
	echo Running docker image of mongo:latest ...
	echo
	docker run -it \
			   --rm \
			   --name db \
			   -p :27017:27017 \
		       --mount source=dbdata,target=/data/db \
			   mongo:latest
}


remove_all_containers() {
	docker rm $(sudo docker ps -a -q)
}

kill_all_containers() {
	docker kill $(sudo docker ps -q)	
}

clean_images() {
	sudo docker rmi \
		$(sudo docker images \
		| awk '$1 == "<none>"' \
		| awk '{print $3;}')
}

# echo ARGS: $1 $2 $3
case $1 in
	"csrun")
		currencyservice_run
  		;;
	"csbuild")
  		currencyservice_build
  		;;
  	"fwrun")
		fixerworker_run
  		;;
	"fwbuild")
  		fixerworker_build
  		;;
  	"rbrun")
		rimbot_run
  		;;
	"rbbuild")
  		rimbot_build
  		;;
  	"mgrun")
		mongo_run
		;;
  	"rmall")
		remove_all_containers
		;;
	"killall")
		kill_all_containers
		;;
	"cleanimg")
		clean_images
		;;
	*)
		echo "default"
		exit 1
esac