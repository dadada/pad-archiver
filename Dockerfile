FROM alpine/git

RUN apk add --no-cache curl

ADD ./update /usr/bin/update

# Override ENTRYPOINT of alpine/git
ENTRYPOINT /bin/sh
