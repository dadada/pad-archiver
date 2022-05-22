FROM alpine/git

ADD ./pad-archiver /usr/bin/pad-archiver

# Override ENTRYPOINT of alpine/git
ENTRYPOINT /bin/sh
