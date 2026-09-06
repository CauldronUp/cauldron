# Trello

Emulates the Trello API (1), for local development and tests.

**15 conformance cases, 3 checked against the live API on 2026-09-05.**

Trello has no sandbox, so the board and card cases still cite documentation. The credential shape needed no account at all, and checking it live found this Recipe's own claim backwards.

## What writing this Recipe changed

This is the clearest case in the collection for reproducing something rather
than improving on it. Trello's credentials travel in the query string, which is
a bad idea, and its failures are plain text rather than JSON.

A fake that accepted a header instead would hide the fact that the key ends up
in access logs. One that answered in JSON would hide that calling `.json()` on a
real error throws. Both would be more pleasant and less true.

## What checking it live found

A present, wrong key answers 401 "invalid key" on every route tried, which this Recipe already claimed correctly. But no key or token at all is not refused: GET and PUT on a board both proceed to the ordinary lookup and 404 on an id that does not resolve, and a POST with nothing sent got far enough to complain about a missing body parameter instead of ever mentioning a credential. A case here had claimed 401 for exactly this GET request and had never been run; it is fixed now, with the route marked `public: when-absent`. The pattern looks Recipe-wide from GET, PUT and POST all agreeing, and this file does not yet apply it everywhere that would need it -- an open gap rather than one closed by assumption. A path nothing declares is a separate answer again: Express's own bare 404, resolved before any credential is judged.

## Sources

- Documentation: https://developer.atlassian.com/cloud/trello/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve trello     # run it
cauldron verify trello -v # check every claim
```
