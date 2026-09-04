# qdrant

Emulates the qdrant API (v1), for local development and tests.

**19 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Qdrant's `status` field is a string on success and an object on failure -- every success carries `"status": "ok"`, every failure carries `"status": {"error": "..."}`, same field name, two different JSON types. In a typed language that is an outright deserialization failure on whichever path gets written second; in an untyped one it is quieter and worse, since `response.status.error` is simply undefined on success, which reads exactly like no error at all.

A write answers acknowledged, not completed -- `wait` defaults to false, so the response only means the operation reached the write-ahead log, not that anything is actually searchable yet, and upsert-then-query is the very first thing anyone writes against this API and it is a race. Two point counts disagree on purpose too: `points_count` is how many points exist and `indexed_vectors_count` is how many are actually in the index, and the gap between them is exactly the set of points you can count but cannot find. Collection health also is not binary -- green, yellow, grey, and red are all valid, where yellow means optimization is running and grey means it is simply possible but untriggered, so code waiting for green before querying waits on collections that were already fine.

The direction a similarity score sorts in is not in the response at all -- cosine and dot products sort high-is-near, while Euclidean and Manhattan distances sort low-is-near, and which one applies lives on the collection's separate config. A threshold like `score > 0.8` is a sensible floor on one collection and a nonsense filter on the next.

## Sources

- Documentation: https://api.qdrant.tech/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve qdrant     # run it
cauldron verify qdrant -v # check every claim
```
