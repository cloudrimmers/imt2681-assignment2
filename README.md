## Assignment 3 - imt2681 Cloud Programming

### Participants
| Name | Studentno |
|-----|----|
| Jonas J. Solsvik | 473193 |
| Halvor | |
| Jone | |


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


### Install

1. Create `.env` file in root directory, *example:*
```
MONGODB_URI=mongodb://localhost
MONGODB_NAME=test
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

### High level system overview





