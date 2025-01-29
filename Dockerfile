FROM postgres:10-alpine

VOLUME [ "/data" ]
COPY ./dat/*.sql /docker-entrypoint-initdb.d/
# docker run -d --rm \
# --name local-pg \
# -e POSTGRES_PASSWORD=postgres \
# -p 5432:5432 \
# -e PGDATA=/var/lib/postgresql/data/pgdata \
# -v /Users/frankmoley/.local/docker/data:/var/lib/postgresql/data \
# postgres
