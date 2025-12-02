FROM liquibase/liquibase:5.0.1

USER liquibase
RUN lpm add mysql --global

USER root
WORKDIR /liquibase/changelog

# Copy master changeset
# COPY changeset-master.yaml /liquibase/changelog/changeset-master.yaml
# ADD changelogs /liquibase/changelog/changelogs
COPY db/ ./

# Copy test changeset and test directories
# COPY changeset-test.yaml /liquibase/changelog/changeset-test.yaml
# ADD changelogs-test /liquibase/changelog/changelogs-test

USER liquibase

ENTRYPOINT ["liquibase"]
CMD ["update"]