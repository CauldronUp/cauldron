# TfL

Emulates the TfL API (tfl), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The severity scale runs backwards.** `statusSeverity`
10 is `"Good Service"` and 6 is `"Severe Delays"`, with `"Part Suspended"` at 3:
the bigger number is the better railway. Every other severity scale a client has
met counts upward into trouble -- syslog, HTTP, CVSS, log levels -- so sorting
descending to put the worst lines first puts the healthiest lines first instead,
silently, on data that looks sorted.

**And every object carries its .NET type and assembly.** `"$type":
"Tfl.Api.Presentation.Entities.Line, Tfl.Api.Presentation.Entities"` --
namespace, class, comma, assembly, on every object at every depth, beginning
with a sigil so it needs a subscript in most languages. A `LineStatus` carries
`"created": "0001-01-01T00:00:00"`, which is .NET's `DateTime.MinValue`
serialised straight onto the wire because nothing set it, beside a `Line` whose
own `created` is a real 2026 date. A failure says what went wrong four times --
`exceptionType` as a .NET class name, the status as a number, the status as a
string, and a sentence -- with a timestamp to seven fractional digits. And a
path with nothing behind it answers `"Resource not found:
http://api:8001/NoSuchEndpoint"`, a service name and a port from inside the
deployment.

## Sources

- Documentation: https://api-portal.tfl.gov.uk/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tfl     # run it
cauldron verify tfl -v # check every claim
```
