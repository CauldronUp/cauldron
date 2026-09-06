# cohere

Emulates the cohere API (v2), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Two were struck live against api.cohere.com on 2026-09-05, and corrected the message: this file had "invalid api token" for every auth failure. Cohere actually sends "no api key supplied" for an absent header and, for a present, wrong one, "Incorrect API key provided: *******. You can find your API key at https://dashboard.cohere.com/api-keys." -- redacting the key rather than restating that it was wrong.

## What this Recipe found

Cohere reports two token counts on every response that don't add up to each other: billed_units counts what you pay for, tokens counts what the model actually saw, and they diverge because system prompts, tool definitions and cached content are seen but billed differently or not at all. A client that counts its own prompt with a tokeniser and multiplies by a price is doing arithmetic on a number Cohere never used.

A finish reason is also not an error -- COMPLETE, MAX_TOKENS, STOP_SEQUENCE, TOOL_CALL and ERROR_LIMIT all arrive with a 200 and a full body, and MAX_TOKENS means the answer was cut off mid-sentence, which for structured output means invalid JSON failing somewhere far downstream from here. A tool call is a finish reason rather than a field to test for, too: the message has no text when the model called a tool, so code that reads the text and moves on gets an empty string and reports it as the model's answer.

Embeddings come back in the order sent with no identifiers attached -- the array index is the only thing tying a vector to its input, so filtering or reordering inputs anywhere in the pipeline silently pairs the wrong vector with the wrong document. No model actually runs here; which finish reason and token counts a response carries is what the fixture puts there, deliberately, since one of the four reasons is nearly impossible to produce on demand against the real API.

## Sources

- Documentation: https://docs.cohere.com/reference/about
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cohere     # run it
cauldron verify cohere -v # check every claim
```
