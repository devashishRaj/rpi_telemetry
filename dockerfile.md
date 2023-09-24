#### Docker file present in repo requires a ".env" file to pass information like login credentials , paths , example of .env file is :

```

# .env
POSTGRES_DB="xyz"
POSTGRES_USER="xyz"
POSTGRES_PASSWORD="xyz"
GF_DATABASE_TYPE="postgres"
GF_DATABASE_HOST="postgres"
GF_DATABASE_NAME="xyz"
GF_DATABASE_USER="xyz"
GF_DATABASE_PASSWORD="xyz"
SCHEMA_PATH="<path to schema.sql file>"
DATABASE_PATH="<path to postgres backup folder on host"
GRAFANA_VOLUME="<path to grafana backup folder on host>"

```
