FROM BASEIMAGE

ARG ARCH
ARG TINI_VERSION

ADD app-postgres-db /usr/local/bin/app-postgres-db

USER 1001
ENTRYPOINT ["app-postgres-db"]
