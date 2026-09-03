# Open Notify

Emulates the Open Notify API (opennotify), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Notify's, where **there is no HTTPS at all**. Port 443 does not
answer: `http://api.open-notify.org/iss-now.json` returns 200 and the same URL
over TLS times out, so no page served over HTTPS can call this API from a
browser. The only published client hardcodes `http://` as its default domain,
because there is no other one to reach.

**And every failure declares JSON and sends none.** A 404 and a 405 both answer
`Content-Type: application/json; charset=UTF-8` with a body of zero bytes, so a
client that trusts the content type calls `.json()` on an empty string and
throws -- having been told by the response itself that a document was coming.
The endpoint the documentation still describes fails differently again:
`/iss-pass.json` is gone and answers nginx's own HTML 404, `text/html`, 169
bytes, naming the version it runs. Two 404s in two content types on one API, and
which one arrives depends on how far in the request got. The position is two
strings beside a number -- `iss_position.longitude` and `iss_position.latitude`
are quoted while the `timestamp` next to them is an integer -- `message` reads
`"success"` on a 200, and `number` is `people.length` sent alongside the array,
on an endpoint where `craft` is not always `"ISS"` because Tiangong appears in
it too.

## Sources

- Documentation: http://open-notify.org/Open-Notify-API/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve opennotify     # run it
cauldron verify opennotify -v # check every claim
```
