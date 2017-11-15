## Description

The bot is the link between:
- slack
- dialogflow
- currencyservice

### Endpoints
| Route  | Req type | Req example | Resp type | Resp example |
|-----| ------| ------- | ------- | ------|
| POST /  | application/x-www-form-urlencoded | text="What is 1 EUR in USD" | application/json | { text:"EUR->USD = 1.5", username:"rimbot" } |


### .env file example
```
PORT=5000
```
