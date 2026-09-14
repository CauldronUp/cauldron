# llamacloud

Emulates the LlamaCloud pipeline listing for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document LlamaCloud serves without a credential at [`api.cloud.llamaindex.ai/api/openapi.json`](https://api.cloud.llamaindex.ai/api/openapi.json), and struck live on 2026-09-13 with no header, with a wrong key, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A wrong key is explained as a geography problem.**

```json
{"detail":"Invalid API Key. Please check your region https://developers.llamaindex.ai/python/cloud/general/regions."}
```

The key was wrong; the sentence tells you to check which continent you are calling, and links to a page about it. A missing key gets FastAPI's own `{"detail":"Not authenticated"}` instead, so the two failures share a field and nothing else.

**The document declares no servers at all.** `servers` is absent, so a generated client has an operation for every one of the 147 paths and no base URL to send any of them to.

**The listing declares two responses and neither is a failure a caller will meet.** `search_pipelines_api_v1_pipelines_get` lists a 200 and a 422, on an API that answers 401 to everything anonymous and 404 to a mistyped path — so `cauldron drift` reports every failure this Recipe serves as unbacked, and each report is right.

**Three paths end in a colon and a verb.** `/api/v1/beta/agent-data/:search`, `/:aggregate` and `/:delete`, all POSTs — and `DELETE /api/v1/beta/agent-data/{item_id}` exists beside them, so one resource has two ways to be deleted and one of them is a POST to a path named delete. Three paths out of a hundred and forty-seven use that convention.

**Thirty-three of the paths are under `/beta/`.** Agent data, batch processing, configurations, directories, attachments and data sinks are all beta, in the same document and under the same version number as everything else.

**Three fields for the embedding configuration, and the required one has no description.** A pipeline carries `embedding_model_config_id` ("The ID of the EmbeddingModelConfig this pipeline is using"), `embedding_model_config` ("The embedding model configuration for this pipeline") and `embedding_config` — and `embedding_config`, the one with no description, is the one in `required`.

**Seven configuration fields and two suffixes for the idea.** `embedding_model_config`, `embedding_config`, `sparse_model_config`, `transform_config`, `metadata_config`, `preset_retrieval_parameters` and `llama_parse_parameters`: five `_config` and two `_parameters`, on one record.

**A field that only means something for half the records.** `managed_pipeline_id` is "The ID of the ManagedPipeline this playground pipeline is linked to" — a description that presumes `pipeline_type` is `PLAYGROUND`, on a record whose type is "Either PLAYGROUND or MANAGED".

**And the timestamps are optional.** `created_at` and `updated_at` are nullable and absent from `required`, so a pipeline may not know when it was made, while `id`, `name`, `project_id` and `embedding_config` must be there.

**`config_hash` is "Hashes for the configuration of the pipeline"** — plural, under a singular name. The listing takes a `session` query parameter. And the operation is `search_pipelines_api_v1_pipelines_get` — the handler name with the path and the verb appended — so a generated client's method is called `searchPipelinesApiV1PipelinesGet`.

## Sources

- [`api.cloud.llamaindex.ai/api/openapi.json`](https://api.cloud.llamaindex.ai/api/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.cloud.llamaindex.ai`, struck 2026-09-13 with no header, a wrong key, an unrouted path, and a wrong method.

## Modelling limits

- **One route of a hundred and forty-seven.** The pipeline listing. Jobs, files, projects, data sources, data sinks, retrievers, extraction agents, chat apps, evaluations and the whole `/beta/` surface are the rest.
- **The nested configuration objects are abridged.** `embedding_config`, `transform_config` and `preset_retrieval_parameters` are unions of many shapes in the document; the fixture carries one member of each, with the field names intact.
- **The success fixture is document-derived.** Listing pipelines needs a real key; the records here are `Pipeline`'s own properties with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** LlamaCloud is reached through `llama-cloud` or `llama-index` on PyPI, or a plain HTTP call carrying a bearer token, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
