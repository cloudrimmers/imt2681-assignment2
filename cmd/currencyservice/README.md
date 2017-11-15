## Description

<<<<<<< HEAD
### Currencyservice
- A RestAPI service providing currency convertion rates.
=======
### Currencyservice 
A RestAPI service providing currency convertion rates.
>>>>>>> f519090512307ffcffe2f783ac99772b187e5aa9

### Endpoints

<<<<<<< HEAD
| URL  | Params | Return type | Return example |
|-----| ------| ------- | ------- |
|GET /latest  |  Query: base, target | text | 1.55 |

=======
| Route  | Req type | Req example | Resp type | Resp example | 
|-----| ------| ------- | ------- | ------|
| POST /currency/latest/  | application/json | { baseCurrency:"EUR", targetCurrency:"NOR"} | text | 1.55 | 
 
>>>>>>> f519090512307ffcffe2f783ac99772b187e5aa9

### .env file example
```
<<<<<<< HEAD
PORT=5000
MONGODB_NAME=assignment3
=======
PORT=1234
MONGODB_NAME=testdb
>>>>>>> f519090512307ffcffe2f783ac99772b187e5aa9
MONGODB_URI=mongodb://localhost
```
