FROM alpine/git

RUN apk add --no-cache curl

ADD update /bin/update
