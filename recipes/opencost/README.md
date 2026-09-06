# opencost

Emulates the OpenCost Kubernetes cost API, for local development and tests.

**8 conformance cases, none checked against a live API.**

Written against OpenCost's own OpenAPI 3.0.1 document, published in its own repository — `opencost/opencost`, `docs/swagger.json`, 5 paths, version 1.0.0 — read on 2026-09-06. OpenCost runs inside a Kubernetes cluster and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**`data` is an array of maps, and `data[0]` is not a record.**

```yaml
data:
  type: array
  items:
    type: object
    additionalProperties:
      $ref: Allocation
```

`data` is a list whose elements are objects keyed by allocation name. Each element is one **time-step** — `step` slices the window — and each step is a map from the thing you aggregated by to its costs.

So `data[0]` is a step, not an allocation; `data[0]["kube-system/coredns"]` is the allocation. `data.map(a => a.totalCost)` yields `undefined` for every element, `data.length` counts time-steps rather than workloads, and a client that flattens the array into records loses the time dimension it was flattening.

The keys of those maps are whatever `aggregate` asked for — `namespace`, `controller`, `label:app` — so **the shape of the response depends on a query parameter**, and no field in the body names the aggregation that produced it.

**One parameter, four grammars.** `window` is documented as accepting "words like `today`, `lastweek`; durations like `30m`, `7d`; RFC3339 date pairs; or Unix timestamps." Four syntaxes in one string with no discriminator, so a client assembling one from user input has to know which grammar it is in before it can validate it. The response normalises what the request did not: each allocation carries `window.start` and `window.end` whichever form asked for it.

**The document's own error text names a parameter it does not declare.** The 400 on `/allocation` is described as "missing or invalid parameters (e.g. `window`, `step`, `accumulateBy`)". The operation declares `window`, `aggregate`, `step` and `accumulate`. There is no `accumulateBy` — so a reader debugging a 400 goes looking for a parameter that does not exist.

**A status in the body beside a code that is not the HTTP status.** `AllocationResponse` carries `code` (integer) and `status` (string) alongside `data` — the same shape [casdoor](../casdoor) has, from a completely unrelated project.

**Six numbers on one record and none is typed as money.** Five costs and a duration in minutes, all plain JSON numbers, so a client summing them accumulates binary floating-point error across every workload in a cluster. `pvCost: 0` is a real cost, not an absence.

**No authentication of any kind, and here that is the design.** OpenCost runs as a service inside a cluster and is reached through the Kubernetes API server or a port-forward. This is the fourth document in the collection declaring no `securitySchemes` — after [argocd](../argocd) (an omission), [cilium](../cilium) (correct — a unix socket) and [exoscale](../exoscale) (a signature OpenAPI cannot describe) — and it is correct for the same reason Cilium's is.

## Modelling limits

- **Nothing here is verified against a live API.** OpenCost runs inside a cluster.
- **The allocation query only.** Five paths, and the other four — assets, cloud costs and the two inference-cost endpoints — each want their own evidence.
- **The response is served raw.** An array whose elements are maps has no representation among this format's list styles: `bare` is an array of records, `map` is one object keyed by identifier, and this is a *list of those objects*. Serving it raw reproduces the shape exactly and is the honest option; the gap is recorded in the backlog.
- **Because it is raw, the fixture is fixed** — seeding different records does not change what this route answers. The `allocation` resource is declared to document the record shape, not to drive the response.
