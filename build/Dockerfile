FROM alpine:3.8

ENV OPERATOR=/usr/local/bin/minecraft-operator-go \
    USER_UID=1001 \
    USER_NAME=minecraft-operator-go

# install operator binary
COPY build/_output/bin/minecraft-operator-go ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
