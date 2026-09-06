# deps.dev

Emulates the deps.dev API (v3), for local development and tests.

**13 conformance cases, 12 checked against the live API on 2026-08-28.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Capitalising a package name gets you a different
package and a 200. On npm the name is case-sensitive, so
`/v3/systems/npm/packages/express` answers with 288 versions and
`/v3/systems/npm/packages/Express` answers with three -- a real package, three
releases in 2016, all deprecated, whose deprecation text is the only thing in
the response that says anything is wrong. `packageKey.name` echoes `Express`
back as though it were right and `isDefault` is set on one of the three, so a
client that title-cased a name reads a nine-year-stale version history and gets
no error at all. **The other ecosystems do not agree**: on PyPI the same request
normalises instead, and `/v3/systems/pypi/packages/Django` answers with
`packageKey.name: "django"` -- the name you asked for silently replaced with the
one it resolved to. One API, one path, two opposite readings of the same
mistake.

The rest of it is the same kind of quiet. The system is echoed in upper case
(`NPM`, `PYPI`) where the path took lower. A version that is not deprecated
carries `deprecatedReason: ""` -- present and empty, the opposite of the npm
registry above, whose `deprecated` is absent rather than false. The dependency
graph is nodes and edges with integer indices rather than a tree, the nodes are
in **alphabetical** order rather than dependency order, and the package itself
is node 0 with relation `SELF`, so the array is one longer than the dependency
count and filtering it makes every edge point somewhere else. A successful graph
says so with `error: ""`. And the failures are bare plain text: seventeen bytes
reading `package not found`, `text/plain`, no trailing newline, no JSON --
which is also what an ecosystem that does not exist gets, naming the wrong one
of the two things in the URL.

## Sources

- Documentation: https://docs.deps.dev/api/v3/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve depsdev     # run it
cauldron verify depsdev -v # check every claim
```
