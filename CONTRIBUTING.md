# Contributing to Cauldron

Cauldron is early. The architecture is still moving, so the most useful thing you can do right now is open an issue before writing a large change.

## The most valuable contribution

**A Recipe for an API you actually use** — especially one whose sandbox requires partner approval, costs money, or takes weeks to obtain. Those are exactly the providers no vendor will build for, and exactly the ones people are blocked on.

The person who owns the weird internal system is the only person who can model it faithfully. That's the whole reason this project is open source.

## Development

Requires Go 1.26 or later. There are no third-party dependencies, and adding one needs a good reason.

```bash
go build -o cauldron ./cmd/cauldron   # build
go test ./...                         # test
go vet ./...                          # vet
gofmt -l .                            # must print nothing
```

## What we're strict about

**Detection never guesses.** If you add a package-to-Recipe mapping, it goes in the explicit table in `internal/detect/providers.go` with the exact package name. No prefix matching, no fuzzy heuristics. A wrong guess sends somebody chasing a bug that doesn't exist, which is worse than not detecting at all.

**Determinism.** Recipes must not read the wall clock, generate unseeded randomness, or touch the network. If a test can fail on a Tuesday, it isn't finished.

**Honest output.** No command may print a success message for work it didn't do. If something isn't implemented, say so in the output.

**Tests come with behaviour.** New detection rules need a case in `internal/detect/detect_test.go`. New commands need a case in `internal/cli/cli_test.go`.

## Commit messages

Explain what changed and why. The why matters more — the diff already shows the what.

## Licence

By contributing you agree that your contributions are licensed under [Apache-2.0](LICENSE).
