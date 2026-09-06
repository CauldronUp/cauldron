# AssemblyAI

Emulates the AssemblyAI API (v2), for local development and tests.

**12 conformance cases, 2 checked against the live API.**

Two were struck live against api.assemblyai.com on 2026-09-05: no Authorization header at all, and the wrong-scheme Bearer-prefixed value the Recipe already modelled. Both come back byte-identical 401. This file had the message as "API token missing or invalid"; AssemblyAI actually sends "API token missing/invalid", with a slash.

## What this Recipe found

A failed AssemblyAI transcription arrives with HTTP 200 -- status is "error" and the reason sits in the error field, so code that branches on the HTTP status sees a success and code that reads text gets null. A queued transcript's text is also null rather than empty or absent, so a check for the key's presence passes and the value is unusable regardless. And utterances is null -- not empty, null -- on a perfectly successful completed transcript if speaker labels weren't requested, so a diarisation feature that forgets the flag gets a clean response with nothing in it.

Word timings are in milliseconds; treating them as seconds compresses an hour of audio into three and a half seconds of subtitles. Nothing here actually transcribes or progresses on its own -- a fixture's status is whatever it was given, which is the point, since the whole reason to emulate this is to remove the wait.

## Sources

- Documentation: https://www.assemblyai.com/docs/api-reference
- Machine-readable description: https://www.assemblyai.com/openapi.json, last checked 2026-09-05
  `cauldron drift assemblyai` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve assemblyai     # run it
cauldron verify assemblyai -v # check every claim
```
