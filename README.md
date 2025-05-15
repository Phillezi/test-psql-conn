# Docker image to test psql connection inside pod network

## Configuration

| Key       | Description           | Default     |
| --------- | --------------------- | ----------- |
| `DB_HOST` | The database hostname | `localhost` |
| `DB_PORT` | The database port     | `5432`      |
| `DB_USER` | The database user     | `postgres`  |
| `DB_PASS` | The database password | `password`  |
| `DB_NAME` | The database name     | `postgres`  |

## How to use

### In docker

If you want to test if you can connect to a postgres db locally with docker you can try this,

```bash
docker run --rm -it \
        -e DB_USER=myuser \
        -e DB_PASS=mypassword \
        -e DB_NAME=mydb \
        -p 8080:8080 \
        ghcr.io/phillezi/test-psql-conn:latest
```

> [!NOTE]
> The above command will try to connect to a postgres db called `mydb` with the user `myuser` and password `mypassword` hosted on the host machine (172.17.0.1:5432).
