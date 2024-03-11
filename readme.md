# ENV VAR

## SHOPIFYx
Copy variables below to your `.bashrc` or `.zshrc`

```
export DB_NAME=postgres
export DB_PORT=5432
export DB_HOST=localhost
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export PROMETHEUS_ADDRESS=comingsoon
export JWT_SECRET=secretjwt
export BCRYPT_SALT=8 # jangan pake 8 di prod! pake > 10
export S3_ID=comingsoon
export S3_SECRET_KEY=comingsoon
export S3_BASE_URL=commingsoon
```

## SHOPIFYx LOCAL MIGRATIONS
### GOLANG-MIGRATE
Please install https://github.com/golang-migrate/migrate
### UP
```
sh scripts/migrate_up_local.sh
```

### DOWN
```
sh scripts/migrate_down_local.sh
```