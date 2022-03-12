# web-api

# Makefile 

# go migrate 

# sqlc 

# testify 

# viper

# mockgen

# PASETO 

# Swagger

# Docker 

For docker change this line in `app.env` file: 

```
DB_SOURCE=postgresql://root:secret@localhost:5432/web-api?sslmode=disable
```

by this one : 

```
DB_SOURCE=postgresql://root:secret@postgres:5432/web-api?sslmode=disable
```