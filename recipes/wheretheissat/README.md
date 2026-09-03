# Where the ISS at

Emulates the Where the ISS at API (wheretheissat), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

The ISS at's, where **the same field is a number on one endpoint
and a string on another.** `/v1/satellites/25544` sends `"latitude":
26.268487855719` and `/v1/coordinates/37.79,-122.39` sends `"latitude":
"37.795517"` -- one API, one field name, two types. `id` does it too: the
satellite carries `25544` and its own two-line element set carries `"25544"`.

**And the units are declared after the numbers they govern.** `altitude` and
`velocity` mean different things depending on a query parameter, and the only
thing that says which is a field at the end of the object: 419.69504466844
kilometres by default and 260.21728290923 miles with `?units=miles`, so a
parser reading fields in order has the numbers before it has the scale. Beside
them, `daynum` is a Julian date -- 2461282.8361574, counting from 4713 BC --
sharing an object with a Unix `timestamp` of 1788077044 and nothing to tell
them apart. The failure puts the status last, `{"error": "satellite not found",
"status": 404}`, and a path with no controller behind it says exactly that:
`"Invalid controller specified (nosuchthing)"`.

## Sources

- Documentation: https://wheretheiss.at/w/developer
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wheretheissat     # run it
cauldron verify wheretheissat -v # check every claim
```
