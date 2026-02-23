FROM mysql:8.3.0

# Copy SQL initialization scripts to the initialization directory
COPY ./point/init.sql /docker-entrypoint-initdb.d/01-point.sql
COPY ./store/init.sql /docker-entrypoint-initdb.d/02-store.sql
COPY ./grandall.sql /docker-entrypoint-initdb.d/03-grandall.sql

EXPOSE 3306
