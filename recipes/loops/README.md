# loops

Emulates the Loops email API for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-06.**

Written against Loops' documentation at `loops.so/docs/api` and struck live against `app.loops.so` on 2026-09-06 with no credential and then with a deliberately invalid one.

## What this Recipe found

**You cannot see a 404 without putting something in the Authorization header.**

Struck live, the same unknown path twice:

```
GET /api/v1/cauldron-nope        (no header)
401 {"success":false,"message":"Invalid API key","error":"Invalid API key"}

GET /api/v1/cauldron-nope        Authorization: Bearer <anything at all>
404 {"message":"Not found"}
```

An unknown path is reported as a credential failure until a credential is present — and it does not have to be a *valid* one. Any string flips the answer.

So a developer with a typo in a URL and no key configured yet is told their key is wrong, which sends them to the wrong file. The two mistakes people make first are reported as one, and the report names the one they did not make.

**The two failures have different envelopes.** The auth failure is `{success, message, error}`. The not-found is `{message}` alone — no `success`, no `error`. A client reading `body.success === false` to decide whether a call failed reads `undefined` on the 404 and treats it as fine.

**`message` and `error` carry the same string.** `"Invalid API key"` twice, in one body, in two fields. Three keys, two facts — and there is no code anywhere, so neither field is the one to branch on.

**A missing key and a wrong key are the same response**, both `"Invalid API key"`. So the sentence is accurate about one of the two cases it covers, and inaccurate about the one where nothing was sent at all.

**Success has no envelope.** `contacts/find` answers a bare array. So the shape of a response depends on whether it worked, and a client has to know which before it can parse: `body.message` on a failure, `body[0].email` on a success, and `body.success` on one failure but not the other.

**"No such contact" is a 200 with an empty array**, not a 404 — the right decision for a search, and it means the only 404 this surface produces is the one that cannot be reached without a credential.

## Detection

`loops` names `app.loops.so` 26 times in its published archive and is mapped. It is the official JavaScript SDK.

## Modelling limits

- **One route.** Contact search. Contact creation and update, transactional sends, events, lists and custom properties each want their own evidence.
- **No `spec:`.** Loops publishes a rendered documentation site.
