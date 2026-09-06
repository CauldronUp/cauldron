# OSV.dev

Emulates the OSV.dev API (v1), for local development and tests.

**14 conformance cases, 13 checked against the live API on 2026-08-27.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

One vulnerability is three records that disagree.
The SQL injection fixed in Django 2.2.28, 3.2.13 and 4.0.4 is in the database
under three ids from three sources, each naming the others in `aliases`:
`GHSA-2gwj-7jmv-h26r` has a summary, a severity and three `affected` entries;
`PYSEC-2022-190` has no summary and no severity; `CVE-2022-28346` has neither
and **no `package` at all** -- no name, no ecosystem, no purl, and a `GIT` range
whose events are commit hashes. So the shape of the answer depends on which name
you looked it up by, and the CVE, the name a scanner is most likely to hold, is
the emptiest of the three. They are not even on one schema version: the CVE
record says `1.9.0` and the other two say `1.7.3`.

**And the records disagree about which version fixed it.** GitHub's splits the
three release branches into three entries with two events each, ascending;
PyPA's packs all three into one entry with six events, descending. So
`affected[0].ranges[0].events[1].fixed` is `2.2.28` on one and `4.0.4` on the
other, for the same vulnerability in the same ecosystem. A query counts each
advisory once per source, too: django 3.2.0 answers 63 vulns of which 26 are
pairs of the same advisory, so there are 37, and the array is sorted by id so
the halves of a pair sit forty entries apart. `versions` is sorted as text --
`2.2.10` before `2.2.2` -- so its last entry is `2.2.9` while the highest
affected release is `2.2.27`. No vulnerabilities is a 200 whose body is `{}`
rather than `{"vulns": []}`. And `code` is a gRPC status on one failure and an
HTTP status on another: `5` beside a 404, `3` beside a 400, and `400` beside a
400 when the body will not parse.

## Sources

- Documentation: https://google.github.io/osv.dev/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve osvdev     # run it
cauldron verify osvdev -v # check every claim
```
