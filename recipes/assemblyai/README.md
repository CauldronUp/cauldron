# AssemblyAI

Emulates the AssemblyAI API (v2), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A failed AssemblyAI transcription arrives with HTTP 200 -- status is "error" and the reason sits in the error field, so code that branches on the HTTP status sees a success and code that reads text gets null. A queued transcript's text is also null rather than empty or absent, so a check for the key's presence passes and the value is unusable regardless. And utterances is null -- not empty, null -- on a perfectly successful completed transcript if speaker labels weren't requested, so a diarisation feature that forgets the flag gets a clean response with nothing in it.

Word timings are in milliseconds; treating them as seconds compresses an hour of audio into three and a half seconds of subtitles. Nothing here actually transcribes or progresses on its own -- a fixture's status is whatever it was given, which is the point, since the whole reason to emulate this is to remove the wait.

## Sources

- Documentation: https://www.assemblyai.com/docs/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve assemblyai     # run it
cauldron verify assemblyai -v # check every claim
```
