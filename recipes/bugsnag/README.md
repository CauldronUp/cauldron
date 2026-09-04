# Bugsnag

Emulates the Bugsnag API (v2), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An Error and an Event are deliberately separate objects in Bugsnag, and almost everyone integrating with the Data API conflates them: an Error is a group carrying only counts (events, users, both plain numbers) with no stack trace anywhere on it, while the stack trace, request, user and breadcrumbs all live on the Event -- getting any detail means a second request per error. fixed is also not terminal: an error marked fixed reopens automatically the next time it happens in a later release, using the same field that said fixed, so a report generated Monday is wrong by Tuesday with nothing edited.

severity and unhandled answer different questions -- an unhandled crash can be severity info and a handled exception can be severity error, because one describes how the code failed and the other how much somebody cares -- so branching on the wrong one pages the wrong person. The event count is also sampled above a threshold and presented as a plain integer, so two errors with the same count haven't necessarily happened the same number of times.

This models the Data API only; the Notify API that actually gets events into Bugsnag is a different host with a different credential and isn't covered.

## Sources

- Documentation: https://bugsnagapiv2.docs.apiary.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bugsnag     # run it
cauldron verify bugsnag -v # check every claim
```
