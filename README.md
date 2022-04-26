# pad-archiver

Archives your pad. Use in your project that contains `pads.txt`, containing a list of URLs (one per line) with `.gitlab-ci.yml` like so:

```
include:
  - project: 'fginfo/pad-archiver'
    file: 'lib/gitlab-ci.yml'
```

The project including the CI configuration has to provide the variable `CI_ACCESS_TOKEN`. It must contain  a project access token that can push to the repo.
