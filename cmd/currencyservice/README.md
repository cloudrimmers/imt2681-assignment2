## Description

### Currencyservice 
A RestAPI service providing currency convertion rates.

### Endpoints

| Route  | Req type | Req example | Resp type | Resp example | 
|-----| ------| ------- | ------- | ------|
| POST /currency/latest/  | application/json | { baseCurrency:"EUR", targetCurrency:"NOR"} | text | 1.55 | 
 

### .env file example
```
PORT=1234
MONGODB_NAME=testdb
MONGODB_URI=mongodb://localhost
```
