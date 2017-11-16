## Assignment 3 - imt2681 Cloud Programming 

*!!UNDER CONSTRUCTION!! - 14.11.17*

### Participants
| Name | Studentno |
|-----|----|
| Jonas J. Solsvik | 473193 |
| Halvor B. Smedås | 473196 |
| Jone | 473181 |


### Project executeables (eventually docker containers)
| Name | Description |
| -----| -----------|
| currencyservice | RestAPI web service, serves data from the mongoDB |
| fixerworker  | Batch process, fetches data from fixer.io, stores in mongoDB | 
| rimbot | A slack bot you can talk to, uses Dialogflow.com and currencyservice | 

### External cloud services
| Name | Description | 
| ----| ------- |
| api.fixer.io | historical currency rates provider | 
| mlab.com | mongoDB provider | 
| dialogflow.com | natural language processing with machine learning | 
| slack.com | chat and collaboration service


### Install and setup environment

1. Create `.env` file in root directory, *example:*
```
export PORT=5000
export MONGODB_URI=mongodb://localhost
export MONGODB_NAME=test
export ACCESS_TOKEN=k2lj3h43lk2jh432klh4   # clara's 
```

2. Create `Procfile` for heroku in root directory, *example:*
```
web: currencyservice
clock: fixerworker
```

3. Run install script
```
./script/install.sh
```

4. Run local heroku instance
```
heroku local
```

5. a) If you do not have a running heroku app already:
```
heroku create
git remote -v
```

5. b) Use existing heroku app by adding it as a remote in git
```
heroku git:remote <app_name>
```

6. Set heroku cloud environment variables
```
heroku config:set MONGODB_URI=mongodb://xxxxxxx
heroku config:set MONGODB_NAME=xxxxx
```

7. Push changes and build re-build Heroku app
```
git add . && git commit -m "a message"
git push heroku master
```

### Install dependencies on Ubuntu 17.10
```
snap install mongo33
snap install go
go get github.com/gorilla/mux
go get github.com/subosito/gotenv
go get gopkg.in/mgo.v2
go get github.com/kardianos/govendor

# Dev dependencies
apt-get install git
snap install heroku 
snap install docker
```

### Directory structure
```yml
root:
  cmd:
    currencyservice:
      app:
        - app.go
        - app_test.go

      - main.go
      - Dockerfile
      - .env

    fixerclock:
      app:
        - app.go
        - app_test.go

      - main.go
      - Dockerfile
      - .env
    
    rimbot:
      app:
        - app.go
        - app_test.go
      
      - main.go
      - Dockerfile
      - .env

  vendor:
    - vendor.json

  - .dockerignore
  - .gitignore
  - docker-compose.yml

```





