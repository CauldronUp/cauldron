# Deepgram

Emulates the Deepgram API (v1), for local development and tests.

**8 conformance cases, 2 checked against the live API.**

Two were struck live against api.deepgram.com on 2026-09-05: the wrong Bearer prefix, a bare key with no scheme at all, and no header whatsoever all come back byte-identical 401 err_code INVALID_AUTH, "Invalid credentials." -- and this file's claim held exactly.

## What this Recipe found

Deepgram and AssemblyAI do the same job and disagree about almost everything, which is why it's worth having both Recipes rather than either alone: Deepgram's pre-recorded endpoint is synchronous and returns the transcript directly, AssemblyAI's is asynchronous and returns an id; Deepgram reports word timings in seconds as floats, AssemblyAI in whole milliseconds; Deepgram buries the transcript four levels down at results.channels[0].alternatives[0].transcript, AssemblyAI puts it at the top. Code written against one and pointed at the other finds nothing at every level.

Both of those channel and alternative arrays are almost always length one, so indexing zero works right up until someone enables multichannel and the second channel gets silently dropped. Failures use err_code and err_msg rather than code and message, so a client reading .message finds nothing, and the credential goes in Authorization with a Token prefix -- not Bearer, not bare, a third spelling that fails as though the key were simply wrong.

This covers the synchronous pre-recorded endpoint only; the callback (asynchronous) form and the WebSocket live-streaming API aren't modelled.

## Sources

- Documentation: https://developers.deepgram.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve deepgram     # run it
cauldron verify deepgram -v # check every claim
```
