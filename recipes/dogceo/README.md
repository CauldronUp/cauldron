# Dog CEO

Emulates the Dog CEO API (dogceo), for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

CEO's, where `message` is a URL, or an array of strings, or an
error sentence. The random-image endpoint answers a URL under it, the sub-breed
list answers an array of names under it, and a breed that does not exist answers
a sentence under it -- and two of the three are strings, so a client cannot tell
the image from the error by type. The only way to know which arrived is to read
`status` first, or to notice that one of them starts with `https`. Code that
assigns `message` straight to an `img` src renders a broken image whose alt text
is the error.

**And the object's shape depends on whether it worked.** A success carries
`message` and `status`; a failure carries `message`, `status` and `code` -- so
the status line is restated twice, once as a word and once as a number, and only
on the way that failed. The sub-breeds are bare names rather than paths or
objects, so a client holding the list has to remember which breed it asked about
to build the next URL.

## Sources

- Documentation: https://dog.ceo/dog-api/documentation/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dogceo     # run it
cauldron verify dogceo -v # check every claim
```
