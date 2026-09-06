# GitLab

Emulates the GitLab API (v4), for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-08-23.**

## What this Recipe found

It answers for a public project. A merge request is
addressed by its per-project iid and not by its global id, and putting one
where the other belongs finds nothing -- or, on a busy instance, finds a
different merge request. A missing project answers one key and one only,
`{"message":"404 Project Not Found"}`, with the status repeated inside the
sentence and nothing called `error` or `errors`. A refused credential answers
the same way, `{"message":"401 Unauthorized"}`, so the status appears twice
and neither copy is a number.

The other two corrected this Recipe rather than confirming it. It had a closed
merge request answering with no `merged_at` and a merged one with no
`closed_at`. GitLab sends both dates on every merge request and nulls the one
that did not happen, so the key is there and its value is not. That is the
difference between `"merged_at" in mr` being false and being true: code asking
whether the key exists got one answer here and the other from GitLab, passed
locally, and read every merge request as merged.

## What writing this Recipe changed

Its fixture carries the third state that breaks integrations: a merge request
that is closed rather than merged. Neither open nor merged, and code branching
on two outcomes has nowhere to put it.

## Sources

- Documentation: https://docs.gitlab.com/ee/api/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gitlab     # run it
cauldron verify gitlab -v # check every claim
```
