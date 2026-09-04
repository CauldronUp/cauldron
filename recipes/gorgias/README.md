# gorgias

Emulates the gorgias API (v1), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A ticket and its messages are two separately paginated endpoints, and a customer replying is an event that can land between reading them. Fetch the ticket and see four messages, then fetch the messages and get five, because the customer replied in between -- there's no version, no etag, no consistent snapshot, just a count on one object and an array on another read at two different moments.

A ticket's `channel` and its messages' channels are separate fields that can disagree -- a ticket opened by email whose agent replied over chat keeps the email channel on the ticket while its messages show both. And `from_agent` is the only thing distinguishing a customer's words from an agent's: it's a bare boolean, so an automated integration reply reads as true and counts as an agent response in any metric built on it. A closed ticket also reopens silently when the customer replies -- the status flips back to open, and the only trace left behind is a `closed_datetime` sitting in the past on a ticket that's marked open.

## Sources

- Documentation: https://developers.gorgias.com/reference/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gorgias     # run it
cauldron verify gorgias -v # check every claim
```
