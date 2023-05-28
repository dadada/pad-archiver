# pad-archiver

Archives the list of URLs provided on the standard input.

```plain
Usage of pad-archiver:
  -C string
        The directory containing the git repository in which to archive the pads. (default "$PWS")
  -password string
        The password for authenticating to the remote. Can also be specified via the environment variable GIT_PASSWORD.
  -push
        Push the changes to the remote specified by remoteUrl.
  -url string
        URL to push changes to.
  -username string
        The username for authenticating to the remote.
```

```shell
go build
pad-archiver < pads.txt
```

## Examples

## Using in GitLab CI

```yaml
Archive pads:
  image: ghcr.io/dadada/pad-archiver
  rules:
    - if: $CI_PIPELINE_SOURCE == "schedule"
  script:
    - --push --url "${CI_PROJECT_URL}.git" --username gitlab-ci-token --password "${CI_ACCESS_TOKEN}" < pads.txt
```

## As a GitHub action

*TODO*

## Using Container Image

```shell
podman run --rm --mount type=bind,source=./.,destination=/repo -i --workdir /repo ghcr.io/dadada/pad-archiver < pads.txt
```
