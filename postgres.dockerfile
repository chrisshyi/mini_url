FROM postgres:12.4

ENV POSTGRES_USER="mini_url"
ENV POSTGRES_PASSWORD="mini_pass"

COPY ./sql/postgres_docker.sql /docker-entrypoint-initdb.d