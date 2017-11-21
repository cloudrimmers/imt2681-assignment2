#!/bin/sh

#
# DOCKER RUN AND BUILD
#
docker_run(){
	name=$1
	port=$2
	echo
	echo Running docker image of $name:latest ...
	echo
	docker run -it \
			   --rm \
			   -p :$port:$port \
			   --name $name \
			   --env-file cmd/$name/.env \
			   --env PORT=$port	\
			   $name:latest
}

docker_build(){
	name=$1
	echo 
	echo  Building docker image of $name:latest ..
	echo
	docker build --tag $name:latest \
			     --file ./cmd/$name/Dockerfile \
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
		docker_run currencyservice 5001
  		;;
	"csbuild")
		docker_build currencyservice
  		;;
  	"fwrun")
		docker_run fixerworker 5002
  		;;
	"fwbuild")
		docker_build fixerworker
  		;;
  	"rbrun")
		docker_run rimbot 5000	  
  		;;
	"rbbuild")
  		docker_build rimbot
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