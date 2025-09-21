Enhance go's default test coverage UI.

Generate test coverage like normal:

```bash
$ go test -coverprofile=cover.prof ./...
$ go tool cover -html=cover.prof -o cover.html
```

Then serve with testviz:

```bash
$ go install github.com/cjea/testviz@latest
$ testviz cover.html > viz.html && open viz.html
```
