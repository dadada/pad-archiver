# pad-archiver

Archives your pad. Use in your project that contains `pads.txt`, containing a list of URLs (one per line) with `.gitlab-ci.yml` like so:

```
include:
  - project: 'fginfo/pad-archiver'
    file: 'lib/gitlab-ci.yml'
```

The including CI configuration has to provide a `CI_ACCESS_TOKEN` that can push to the repo.
