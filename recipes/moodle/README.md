# Moodle

Emulates the Moodle API (rest), for local development and tests.

**15 conformance cases, 9 checked against the live API on 2026-09-03.**

## What this Recipe found

**Every failure is a success with an exception inside**
-- nine ways to fail, including an invented operation name, all answering 200.

## Sources

- Documentation: https://moodledev.io/docs/apis/subsystems/external
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve moodle     # run it
cauldron verify moodle -v # check every claim
```
