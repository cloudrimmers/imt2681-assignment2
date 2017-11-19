## Description

### Currencyservice
- A RestAPI service providing currency convertion rates.
=======
### Currencyservice 
A RestAPI service providing currency convertion rates.

### Endpoints

| URL  | Params | Return type | Return example |
|-----| ------| ------- | ------- |
|GET /latest  |  Query: base, target | text | 1.55 |

=======
| Route  | Req type | Req example | Resp type | Resp example | 
|-----| ------| ------- | ------- | ------|
| POST /currency/latest/  | application/json | { baseCurrency:"EUR", targetCurrency:"NOR"} | text | 1.55 | 
 

### .env file example
```
PORT=5000
MONGODB_NAME=assignment3
=======
PORT=1234
MONGODB_NAME=testdb
MONGODB_URI=mongodb://localhost
```
