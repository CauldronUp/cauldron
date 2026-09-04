# Vimeo

Emulates the Vimeo API (3.4), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The privacy value `unlisted` is described as a fact about search, not about access -- Vimeo's own field description reads "Not searchable from vimeo.com," sitting in a list where every other value (`anybody`, `contacts`, `nobody`, `password`, `users`) is a rule about who may watch. It restricts nobody. The document's other privacy enum, for a user's default, describes the identical word correctly: "Anybody can view the user's videos if they have a link." One value, two descriptions, and only the one attached to the field people actually set is the misleading one.

There are also two fields both called `status` that answer different questions -- one is playability (`playable`, `restricted`), the other is upload/transcode state -- and a video has no `id` field at all, only a `uri` like `/videos/12345` that the number has to be cut off the end of.

## Sources

- Documentation: https://developer.vimeo.com/api/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vimeo     # run it
cauldron verify vimeo -v # check every claim
```
