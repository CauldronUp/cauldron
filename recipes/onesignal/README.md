# OneSignal

Emulates the OneSignal API (v1), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against onesignal.com, no account and no key. No credential at all is refused with a different sentence than this file declared. More surprising: a client doing standard HTTP Basic -- the key, base64 encoded with a colon -- does not read as "no credential" to OneSignal's parser at all. It reads as a deprecated v1 user token, and answers 400 with two error messages instead of 401 with one, which is what this file's own existing case for exactly that request had asserted and never checked.

## What this Recipe found

A OneSignal send with nobody to actually send to still answers HTTP 200, with an `errors` array buried in the body -- the request was fine and the notification reached nobody, so code that branches on status code alone counts it as delivered. This is exactly what happens when a segment is misspelled or empty, which makes it the failure most worth reproducing here.

`recipients` in the response is how many devices were targeted, not how many actually received anything -- delivery is asynchronous and its real outcome only shows up later via webhook or the dashboard, so this number is a ceiling, not a result. The REST API key also goes after `Basic ` completely unencoded, no base64 and no colon, despite the header name, so a client that does the standard Basic-auth thing produces a header OneSignal rejects. A device that has uninstalled the app is still a subscription record with an id and a negative `notification_types` value, so counting subscriptions is not the same as counting an actual audience.

## Sources

- Documentation: https://documentation.onesignal.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve onesignal     # run it
cauldron verify onesignal -v # check every claim
```
