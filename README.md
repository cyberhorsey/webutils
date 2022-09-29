# webutils

## Usage

To run unit tests:

```bash
go test ./... -v
```

To run linter:

```bash
go vet ./...
# command returns exit code 1 if code is not correctly formatted - https://circleci.com/blog/enforce-build-standards/
'! go fmt ./... 2>&1 | read'
```

To run race detector:

```bash
go test -race ./...
```

To test code coverage (will output report as `.html` file):

```bash
go test -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -html coverage.out -o coverage.html
go tool cover -func=coverage.out
```
