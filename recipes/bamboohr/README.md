# BambooHR

Emulates the BambooHR API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

BambooHR's own docs state the trap directly: a terminated employee and one who never existed produce the identical absence from GET /employees/directory, because that endpoint excludes inactive and former employees entirely -- "absence from this response... does not mean no such employee exists." The directory's fieldset is also fixed by company configuration and architecturally cannot include status, employmentStatus or terminationDate at all; GET /employees/{id} has those fields but returns nothing by default unless the caller names them explicitly.

The permissions model generalizes the same trap to any field: a caller who lacks permission gets null values with the suppressed names listed in _restrictedFields on the newer List Employees endpoint, but on GET /employees/{id} an unrequested field and a permission-denied field are indistinguishable -- both are simply missing, with no list to tell them apart. And BambooHR's own OpenAPI schema disagrees with itself about whether an employee id is a string or an integer, in three different property definitions in the same document.

Live probing (struck against real, nonexistent subdomains) found that tenant resolution happens before anything else -- routing, method, credential -- and returns an identical empty 404 regardless, which is the only thing this Recipe could actually verify live; every field-level claim is otherwise sourced to BambooHR's own published documentation.

## Sources

- Documentation: https://documentation.bamboohr.com/reference/get-employees-directory
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bamboohr     # run it
cauldron verify bamboohr -v # check every claim
```
