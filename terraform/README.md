# terraform config

## init

Make sure a `.tfbackend` file exists with the PostgreSQL connection string defined. File name should be something like `config.pg.tfbackend`. Example contents:

```terraform
conn_str = "postgresql://webster:verysecurepassword@somepostgresqlserver:34567/defaultdb?sslmode=verify-full"
```

Then initialize terraform using the backend file

```sh
terraform init -backend-config="config.pg.tfbackend"
```
