# cloudconvert

Emulates the CloudConvert jobs API for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-13.**

Struck live against `api.cloudconvert.com` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist. Record shapes come from CloudConvert's own Node client, which is open.

## What this Recipe found

**Two failures, two formatters.**

```
/v2/jobs           401 {"message":"Unauthenticated.","code":"UNAUTHENTICATED"}
/v2/cauldron-nope  404 {
                           "message": "The route v2/cauldron-nope could not be found.",
                           "code": "NOT_FOUND"
                       }
```

The 401 is compact and the 404 is pretty-printed with four spaces of indent. Same API, same two keys, two encoders — visible only in the bytes, which is where a recorded-response test looks and nothing else does.

**The 404 names the route without its leading slash.** `v2/cauldron-nope`, not `/v2/cauldron-nope` — the path as the router stored it rather than as the caller sent it. A client matching the message against the URL it requested will not find it.

**`Unauthenticated.`** — with the full stop — is the framework's default sentence, and it answers a missing credential and a wrong one alike, so the two cannot be told apart.

**A task changes shape depending on where you read it.** From the client's own types:

```ts
interface JobTask extends Omit<Task, 'job_id' | 'status'> { name; status }
```

A task fetched on its own has `job_id` and a `TaskStatus`. The same task inside a job has neither — it has a `name` and a `JobTaskStatus`, which is `TaskStatus` plus `queued`. One object, two shapes, and one of its statuses has an extra value, decided by which endpoint returned it.

And the dependency map keys tasks by the `name` they only have inside a job, so `depends_on_tasks` is unreadable against a task fetched on its own.

**The job's status is typed as a task's.** `Job.status` is `TaskStatus`, and the same file separately defines `JobStatus = 'processing' | 'finished' | 'error'` — a type the `Job` interface never uses. The client publishes a type for a job's status and types the job's status as something else.

**Every task carries a `message` and a `code` even when nothing went wrong.** Both are `string | null` on the record, so a successful task's error fields are present and null, and a client testing for their presence finds them.

**And two fields are optional *and* nullable.** `retry_of_task_id?: string | null` and `retries?: string[] | null` — absent, null, or a value: three states where two would do.

## Sources

- Live: `api.cloudconvert.com`, struck 2026-09-13.
- [`JobsResource.ts`](https://github.com/cloudconvert/cloudconvert-node/blob/master/lib/JobsResource.ts) — `Job`, `JobTask` and the unused `JobStatus`.
- [`TasksResource.ts`](https://github.com/cloudconvert/cloudconvert-node/blob/master/lib/TasksResource.ts) — `Task`, and the optional-and-nullable pair.

## Modelling limits

- **One route.** Listing jobs. Tasks, imports, exports, conversions, webhooks, signed URLs and the user surface each want their own evidence.
- **The whitespace difference is recorded, not served.** Cauldron serialises every body the same way, so the two failures here are byte-comparable and live they are not.
- **`retry_of_task_id` and `retries` are not modelled.** They are optional and nullable on a task, and nothing this Recipe does produces a retried one.
- **No `spec:`.** CloudConvert's API reference is a single-page application that answers the same 210KB of HTML at every path, so there is no document to fetch. The record shape comes from the Node client, which is the artefact a description would have produced.
- **Nothing is mapped in detection.** The `cloudconvert` packages on npm and Packagist are real clients, but CloudConvert is most often reached through its Zapier and Make integrations or a signed upload URL, neither of which names a dependency. Checked 2026-09-13.
