FROM mysql:8.3.0

# Copy SQL initialization scripts to the initialization directory
COPY ./grandall.sql /docker-entrypoint-initdb.d/01-grandall.sql
COPY ./store/init.sql /docker-entrypoint-initdb.d/02-store.sql
COPY ./point/init.sql /docker-entrypoint-initdb.d/03-point.sql

EXPOSE 3306
