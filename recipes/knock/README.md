# knock

Emulates the knock API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Triggering a workflow answers with just a `workflow_run_id` -- not the messages, not the recipients, not even a count. The messages don't exist yet at the moment of that 200; Knock resolves who gets notified, on which channels, filtered through their preferences, entirely after the call returns, so the only way to learn what actually happened is to go look separately.

A message that was never sent is still a message: `status: "not_sent"` with a reason sits in the same list as delivered ones, so code that just counts messages counts intentions rather than outcomes. Preferences are also layered -- a workspace default, a per-tenant set, and a per-user set -- and the effective answer is Knock's own merge of all three, not any single layer, so reading a user's stated preferences tells you what they chose, not what will actually happen. And seen, read and archived are three independent timestamps: a message can be archived without ever being read, and an unread count that treats archived as read is wrong in exactly the direction users notice.

## Sources

- Documentation: https://docs.knock.app/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve knock     # run it
cauldron verify knock -v # check every claim
```
