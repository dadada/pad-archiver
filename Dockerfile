FROM alpine/git@sha256:8f4173e730f0ae6df38e35695120ab77a0c3e0593d34b6cbe7ee585497f61013

ADD ./pad-archiver /pad-archiver

# Override ENTRYPOINT of alpine/git
ENTRYPOINT /bin/sh
