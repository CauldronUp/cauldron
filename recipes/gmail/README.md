# Gmail

Emulates the Gmail API (v1), for local development and tests.

**16 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A message listing carries no message. `messages.list` returns only identifiers and thread identifiers -- no subject, no sender, no snippet, no date -- so rendering an inbox from it is impossible; the fix is one fetch per message, which is where the quota actually goes. Nearly every first Gmail integration is written expecting a subject to be sitting right there.

A handful of other shapes surprise on first contact. Headers are an array of name/value pairs rather than an object, so there's no `headers.Subject` to reach for, only a list to search, with names in whatever case the sending client used. A multipart message has nothing in `payload.body` at all -- the content lives in `payload.parts`, base64url encoded, a different alphabet from ordinary base64. And there's no read flag: unread is the presence of the `UNREAD` label, so a message is read precisely when a label is missing.

Trash isn't delete either: `messages.trash` just adds the `TRASH` label and the message is still there and still fetchable, while `messages.delete` is permanent and skips the bin entirely -- worth knowing before writing a cleanup routine.

## Sources

- Documentation: https://developers.google.com/gmail/api/reference/rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gmail     # run it
cauldron verify gmail -v # check every claim
```
