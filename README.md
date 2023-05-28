# pad-archiver

Archives the list of URLs provided on the standard input.

```plain
Usage of /bin/pad-archiver:
  -C string
        git directory (default "/repo")
  -password string
        password
  -push
        push repository to remote
  -url string
        url of remote
  -username string
        username
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
