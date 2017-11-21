

## Assignment 3 - IMT2681 Cloud Programming 

### Participants

| Name             | Studentno |
| ---------------- | --------- |
| Jonas J. Solsvik | 473193    |
| Halvor B. Smedås | 473196    |
| Jone             | 473181    |

### Highlights

_Parts of our assignment we're particularly happy about_



- Managing multiple .env files through a configuration script
- Building minimal Docker images
- 4 images orchestrated by Docker Compose
- Deploying simultaneously on both Open Stack and Heroku across 3 domains
- Strict separation of application and library code.
- Stripping down and simplifying assignment 2 code to bare minimum - No unnecessary code in the project


- Support for <Amount> in queries - "Rimbot, in greenbills, how much is 100 New Zealand $"

- Additional, "non-conventional" currency identifiers as showed above, greenbills = USD

- Support for all currency symbols (kr, $, ¥, ₪, £, zł, etc.) of currencies present in https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html

- Randomized bot names, making the experience  of using Rimbot more fun

  ​


### Project executeables (eventually docker containers)
| Name            | Description                              |
| --------------- | ---------------------------------------- |
| currencyservice | RestAPI web service, serves data from the mongoDB |
| fixerworker     | Batch process, fetches data from fixer.io, stores in mongoDB |
| rimbot          | A slack bot you can talk to, uses Dialogflow.com and currencyservice |

### External cloud services
| Name           | Description                              |
| -------------- | ---------------------------------------- |
| api.fixer.io   | historical currency rates provider       |
| mlab.com       | mongoDB provider                         |
| dialogflow.com | natural language processing with machine learning |
| slack.com      | chat and collaboration service           |
##Installation and Setup

### 	Locally

1. Generate the required (for local usage) .env files by running `./script/env-generator.sh`

2. Configure the project by defining the env-variables for the three services now available in `./cmd/<app>/.env`

3. Setup Procfiles:

   Procfile:

```Procfile
web: currencyservice local
clock: fixerworker local
```

​	Procfile_Bot

```Procfile
web: rimbot local
```

4. Run install script `./script/install.sh`

5. Run local heroku instances (one with each Procfile, and the respective env files)

   For the currencyservice and fixerworker Procfile, use the .env file of currencyservice.

```
heroku local -f Procfile -e ./cmd/currencyservice/.env
```

​	in another terminal window:

```
heroku local -f Procfile_Bot -e ./cmd/rimbot/.env
```

​	alternatively

```
rimbot local
```

The bot service is now running alongside it's dependents. Using _Postman_ or similar tools now lets you post queries in a url-encoded form through the service. Slack itself can't take use of your localhost bot service however, so to fully take use of this service you will need to deploy it on

### Remotely using Heroku

1. a) If you do not have a running heroku app already: (needs to be done for each of the two services)

```
heroku create
git remote -v
```

5. b) Use existing heroku app/s by adding remotes in git
```
heroku git:remote <app_name>
```

6. Set Heroku config variables (PORT may be omitted, it's controlled by Heroku anyways)
```
source ./cmd/<app>/.env
heroku config:set VAR1=$VAR1
heroku config:set VAR2=$VAR2
...
```

7. Configure `Procfile` based on which app is being deployed (using the different heroku remotes you created in 1) 

   For currencyservice & fixerworker:

   ```procfile
   web: currencyservice heroku
   clock: fixerworker heroku
   ```

   For rimbot

   ```procfile
   web: rimbot heroku
   ```

8. Push changes and build re-build Heroku app
```
git add . && git commit -m "a message"
git push <heroku_remote> master
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
  - Procfile
  - .dockerignore
  - .gitignore
  - docker-compose.yml
```





