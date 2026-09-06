# Guardian

Emulates the Guardian API (v1), for local development and tests.

**16 conformance cases, 14 checked against the live API on 2026-09-01.** The unchecked one is the paging case: it sends the two parameter names this Recipe declares, which is a claim about the provider read from its own description rather than struck against it.

## What this Recipe found

**A parameter that looks like a filter adds an
object.** `show-fields` does not narrow the response; naming any of `headline`,
`byline` or `wordcount` attaches a whole `fields` object that exists under no
other circumstances, so `result.fields.headline` is `undefined` on every search
that forgot it. Inside, `wordcount` is the string `"571"`. An unrecognised field
name is neither refused nor reported. Everything sits under a `response`
wrapper -- `body.results` is `undefined` -- except the 401, which drops the
wrapper entirely, so the one shape a client can rely on is absent exactly when
it is needed.

## Sources

- Documentation: https://open-platforms.theguardian.com/documentation/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve guardian     # run it
cauldron verify guardian -v # check every claim
```
