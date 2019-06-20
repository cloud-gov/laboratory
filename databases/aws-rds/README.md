## aws-rds

Quick and dirty Go program used in database broker smoke tests.

### Usage

Set up a manifest file like this:

```yaml
applications:
  - name: aws-rds-smoke-test
    buildpacks:
      - go_buildpack
    env:
      DB_TYPE: oracle | mysql | postgres
      ENABLE_FUNCTIONS: true | false
      GOVERSION: go1.12
      GOPACKAGENAME: aws-rds
      SERVICE_NAME: # name of service instance to smoke test
      CGO_CFLAGS: -I /app/code/vendor/include/oracle/
    services:
      - # name of service instance to smoke test
```

1. `cf push aws-rds-smoke-test -f manifest.yml --no-start`
1. `cf create-service ...`
1. `cf bind-service aws-rds-smoke-test <si>`
1. `while [ $? -ne 0 ]; do !!; done` <- keeps trying to bind until the DB is provisioned, then exits
1. `cf aws-rds-smoke-test set-env DB_TYPE "postgresl"`
1. `cf start aws-rds-smoke-test`
1. If the app starts successfully, your brokered database service was able to be written to.

### Notes

This tool vendors some Oracle binaries, which are licensed separately. You can find them under the `include/oracle` library, and the license at `include/oracle/BASIC_LICENSE`.

It's also important to note that `go mod -vendor` will not copy over C source and header files, so if you run `go mod -vendor`, you will lose those. If you do accidentally run it, never fear, just download and run this tool: https://github.com/nomad-software/vend.