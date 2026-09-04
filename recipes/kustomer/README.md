# kustomer

Emulates the kustomer API (v1), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

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
