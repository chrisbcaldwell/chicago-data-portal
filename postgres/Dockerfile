FROM postgres:alpine
ENV POSTGRES_DB=chicago
ENV POSTGRES_USER="postgres"
ENV POSTGRES_PASSWORD="italianbeef"
COPY create_tables.sql /docker-entrypoint-initdb.d/