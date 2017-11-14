#!/bin/sh

currencyservice_run(){
	echo
	echo Running docker container of currencyservice:latest ...
	echo
	docker run -it \
			   --rm \
			   -p 127.0.0.1:5001:5001 \
			   --name currencyservice \
			   --env-file cmd/currencyservice/.env \
			   --link db \
			   --env MONGODB_URI=db \
			   currencyservice:latest
}

currencyservice_build(){
	echo 
	echo  Building docker image of currencyservice:latest ..
	echo
	$(govendor sync)
	docker build --tag currencyservice:latest \
			     --file ./cmd/currencyservice/Dockerfile \
		         .
}

mongo_run() {
	echo
	echo Running docker image of mongo:latest ...
	echo
	docker run -it \
			   --rm \
			   --name db \
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
  	"mgorun")
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