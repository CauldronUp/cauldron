# Go proxy

Emulates the Go proxy API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-27.**

## What this Recipe found

Proxy's, which refuses the module path in the URL and
then hands it back in the body. A module path is case-sensitive and the proxy
protocol will not carry capitals, so `github.com/BurntSushi/toml` -- the
spelling in every `go.mod` -- is fetched as `github.com/!burnt!sushi/toml`, and
asking for the author's spelling answers **404 with the sentence "bad request"**.
The escaped form then returns `Origin.URL: "https://github.com/BurntSushi/toml"`,
capitals restored, from the server that would not take them. The JSON keys are
Go struct field names -- `Version`, `Time`, `Origin`, `VCS`, `URL`, `Ref`,
`Hash` -- because the type is marshalled straight out of `cmd/go` with no tags.
And a module nobody has published answers with the proxy's own shell failure:
the `git ls-remote` it ran, a path inside its cache on its own disk, `exit
status 128`, git's "terminal prompts disabled", and two lines of advice.

## Sources

- Documentation: https://go.dev/ref/mod#goproxy-protocol
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve goproxy     # run it
cauldron verify goproxy -v # check every claim
```
