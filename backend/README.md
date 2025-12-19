# Go template

A minimal template for Go that includes JWT and PostgreSQL.

## Run docker
```bash
cd docker
docker compose up -d
```

## Run app
```bash
cd backend
export CONFIG_PATH=./config/dev.yaml
make run
```
