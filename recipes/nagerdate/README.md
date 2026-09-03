# Nager.Date

Emulates the Nager.Date API (v3), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

`Fixed` is false on Christmas Day. The
twenty-fifth of December is as fixed as a date gets, and the field that says
whether a holiday falls on the same date every year says no -- as it does for New
Year's Day, and for every other holiday: across Canada, the United States, the
United Kingdom and Germany, eighty-one holidays in 2026, `fixed` is false
eighty-one times. `launchYear` is null on all eighty-one too. Both are fields
with one value, so branching on either is branching on a constant.

**And one date carries six holidays, none of them national.** The third of August
2026 in Canada is the Civic Holiday, British Columbia Day, Heritage Day, New
Brunswick Day, Natal Day and Saskatchewan Day -- six entries, one date, thirteen
provinces and territories between them, every one `global: false`. So
`holidays.find(h => h.date === today)` returns whichever the array ordered first,
and calling it "the holiday" is wrong in five provinces out of six. The field
naming those subdivisions is called `counties`, holds Canadian provinces and
German states, and is `null` rather than empty when the holiday is national. And
the two failures share an envelope and disagree about the rest of it: an unknown
country carries `detail`, an unsupported year carries `errors`, and the `title`
on the second is a framework's default validation sentence with the real message
two levels down inside an array.

## Sources

- Documentation: https://date.nager.at/swagger/index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nagerdate     # run it
cauldron verify nagerdate -v # check every claim
```
