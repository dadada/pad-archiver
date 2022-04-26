# pad-archiver

Archives the list of URLs provided on the standard input.

```shell
update < pads.txt
```

The CI-config in `lib/gitlab-ci.yml` can be used in your project by including the following at the top of your project's `.gitlab-ci.yml`.

```yaml
include:
  - project: 'fginfo/pad-archiver'
    file: 'lib/gitlab-ci.yml'
```

The project including the CI configuration has to provide the variable `CI_ACCESS_TOKEN`. The variable must contain  a project access token that can push to commits to your project.
