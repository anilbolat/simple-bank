# a simple bank app

# requirements

1. create and manage account
   - owner, balance, currency

2. record all balance changes
    - create an account entry for each change

3. money transfer transaction
   - perform money transfer between 2 accounts consistently within a transaction


# db
### options
    - database/sql
    - gorm
    - sqlx
    - sqlc (https://sqlc.dev/)
    - gorp (https://github.com/go-gorp/gorp)

#### why is sqlc used ?
    - Very fast & easy to use
    - Automatic code generation
    - Catch SQL query errors before generating code
    - Full support Postgres. MySQL is experimental.

#### sql stmts
    - https://dbdiagram.io/
    - simple-bank.sql
    - sqlc.yaml
 


    



