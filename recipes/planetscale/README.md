# PlanetScale

Emulates the PlanetScale API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The field called `state` on a PlanetScale deploy request is not the state of the deployment -- it is whether the review is open or closed, while a completely different field, `deployment_state`, tracks whether the schema change actually shipped. A request that was opened and closed without ever deploying reads `state: closed`, identical to one whose migration ran an hour ago, so a dashboard built around the field with the more obvious name is confidently wrong in the safe-looking direction.

The identifier a client would naturally store also cannot address the resource: a deploy request has a globally unique `id`, but every path takes its `number` instead, which is scoped per database -- so `1` exists in every database an organization owns, and storing the id gets a 404 that looks like a deleted record. A deploy request also outlives the branch it came from; `branch` still holds the branch name after that branch is deleted, and only a separate `branch_deleted` boolean says it is gone.

The ten endpoints that actually move a deploy through its lifecycle, deploy, apply, complete, cancel, revert, force-cutover, are not modelled, since this format cannot express a route that advances a state machine -- what is served instead is the pair of confusing fields those endpoints move between, which is the point of the Recipe.

## Sources

- Documentation: https://api-docs.planetscale.com/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve planetscale     # run it
cauldron verify planetscale -v # check every claim
```
