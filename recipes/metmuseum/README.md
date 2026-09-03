# Met Museum

Emulates the Met Museum API (metmuseum), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Museum's, where `artistGender` is `"Female"` or
the empty string. There is no `"Male"`. A man is recorded as the absence of a
value -- and the absence of a value is what twenty other fields on the same
record use to mean "this does not apply". `culture`, `period`, `dynasty`,
`reign`, `city`, `country`, `locus` and `excavation` are all `""` on the van
Gogh, so `artistGender`'s empty string sits in a crowd of empty strings that
mean something else entirely. Anything counting women in the collection can do
it; anything counting men cannot tell them from the unattributed, the anonymous,
and the records nobody has filled in.

**And absence is spelled `""` rather than `null` on twenty-one of fifty-seven
fields**, so `if (object.culture)` is falsy, `object.culture === null` is false,
and `"culture" in object` is true -- the three ordinary ways of asking whether a
field is set disagree with each other. Three year fields carry two types:
`accessionYear` is the string `"1993"`, `artistBeginDate` the string `"1853"`,
and `objectBeginDate` the number `1889`. The same `dimensions` field uses U+00D7
on one record and the letter `x` on another. A missing image is `""` where a URL
goes, so assigning it to an `img` src asks the page for its own address. And
`GalleryNumber` is the only capitalised field name on a record that also shouts
`artistWikidata_URL` and `AAT_URL` in the middle: fifty-seven fields, three
naming conventions.

## Sources

- Documentation: https://metmuseum.github.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve metmuseum     # run it
cauldron verify metmuseum -v # check every claim
```
