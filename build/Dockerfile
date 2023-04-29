FROM scratch

ENV APP=/usr/local/bin/merkle-tree \
    USER_UID=1101 \
    USER_NAME=merkle-tree

# install engine binary
COPY build/_output/bin/merkle-tree ${APP}
USER ${USER_UID}:${USER_UID}
CMD [ "./${APP}" ]
