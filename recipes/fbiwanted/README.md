# FBI Wanted

Emulates the FBI Wanted API (fbiwanted), for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Wanted's, where **the reward is zero and the sentence beside it
says twenty-five thousand dollars.** `reward_min` and `reward_max` are both `0`,
and `reward_text` reads "The FBI is offering a reward of up to $25,000 for
information leading to the identification, arrest, and conviction ...". The
numbers are the machine-readable reward and the prose is the real one, so
anything sorting by `reward_max` puts a twenty-five-thousand-dollar case last
and anything filtering on `reward_max > 0` does not show it at all.

**And thirty-three of the record's fifty-four fields are null.** Every physical
descriptor -- `sex`, `hair`, `eyes`, `build`, `complexion`, `weight`, `race`,
`scars_and_marks`, `place_of_birth` -- because the subject has not been
identified, on a record that still carries a `warning_message` reading "SHOULD
BE CONSIDERED ARMED AND DANGEROUS". Absence is spelled two ways among fields of
the same kind: eight lists are `null` and one, `coordinates`, is `[]`. `detail`
is an array on one failure and a string on the other. Two timestamps in one
record disagree about time zones. And `description` packs three facts into one
string with carriage returns between them.

## Sources

- Documentation: https://www.fbi.gov/wanted/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fbiwanted     # run it
cauldron verify fbiwanted -v # check every claim
```
