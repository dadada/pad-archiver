# pad-archiver

Archives the list of URLs provided on the standard input.

```shell
go build
pad-archiver < pads.txt
```

## Examples

## Using in GitLab CI

```
Archive pads:
  image: ghcr.io/dadada/pad-archiver
  rules:
    - if: $CI_PIPELINE_SOURCE == "schedule"
  script:
    - --push --url "${CI_PROJECT_URL}.git" --username gitlab-ci-token --password "${CI_ACCESS_TOKEN}" < pads.txt
```

## As a GitHub action

*TODO*
