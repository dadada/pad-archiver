FROM alpine/git

RUN apk add --no-cache curl

ADD ./pad-archiver /usr/bin/pad-archiver

# Override ENTRYPOINT of alpine/git
ENTRYPOINT /bin/sh
