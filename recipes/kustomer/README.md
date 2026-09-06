# kustomer

Emulates the kustomer API (v1), for local development and tests.

**20 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.kustomerapp.com, no account and no key -- a missing credential and a wrong one both answer the identical body. This file's own existing case had the sentence travelling under `detail`; the real failure carries it under `title` and has no `detail` field at all, so the case was passing against a shape the API does not send. Fixed, and the fix is scoped to this one failure rather than moved to the shared envelope, since the other failures in this file carrying `detail` have not been checked live.

## What this Recipe found

Kustomer models support as a customer's timeline, not a queue of tickets, and that one decision explains almost everything an integration written for a ticket-shaped helpdesk gets wrong. A conversation belongs to a customer and collects messages across every channel they've ever used, so a customer who emailed in March and chatted in August may be one conversation or two depending on workspace configuration -- "one conversation is one issue" is an assumption the product never promises.

Everything is JSON:API, so `conversation.subject` finds nothing; it's `conversation.attributes.subject`, on every object. An identifier also appears twice under two names -- `id` is Kustomer's, `externalId` is whatever the upstream system called it, and only one of them works in a path. And assignment lives on the conversation, not on any individual message: there's an `assignedUsers` array on the conversation and a `createdBy` on the message, and they answer different questions entirely, so a conversation can be snoozed and still assigned, or open and sitting in nobody's queue.

## Sources

- Documentation: https://developer.kustomer.com/kustomer-api-docs/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kustomer     # run it
cauldron verify kustomer -v # check every claim
```
