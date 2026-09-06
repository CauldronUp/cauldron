# TheMealDB

Emulates the TheMealDB API (themealdb), for local development and tests.

**10 conformance cases, 9 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

`Meals` is an array, or `null`, or a string. A
lookup that finds something answers `{"meals": [ {...} ]}`; one that finds
nothing answers `{"meals": null}`; and one with no identifier at all answers
`{"meals": "Invalid ID"}`. One key, three JSON types, and HTTP 200 on all three.
`body.meals.length` works on the first, throws on the second, and returns 10 on
the third -- because `"Invalid ID"` is ten characters long -- so a client
checking `if (body.meals)` treats the error message as a truthy result and
iterates it character by character.

**And forty of the record's fifty-four fields are one array flattened into
columns.** `strIngredient1` through `strIngredient20` and `strMeasure1` through
`strMeasure20`, so a nine-ingredient recipe carries eleven empty slots -- and
they are not empty the same way. Slots 10 to 15 are the empty string and slots
16 to 20 are `null`, in one numbered series, split at an index nothing explains.
`strTags` is a comma-joined string where an array belongs. `strInstructions`
carries CRLF inside a JSON string. And `strArea` is `"Japanese"` beside
`strCountry`'s `"Japan"`: the adjective and the noun, in two fields, on one
record.

## Sources

- Documentation: https://www.themealdb.com/api.php
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve themealdb     # run it
cauldron verify themealdb -v # check every claim
```
