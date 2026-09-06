# Agora

Emulates the Agora API (dev/v1), for local development and tests.

**17 conformance cases, 9 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **collapse runs the other way** from its neighbour's:
five broken shapes fold into rejected, and only a truly absent header stands
apart.

## Sources

- Documentation: https://docs.agora.io/en/api-reference/api-ref/console/solutions-agora-console-rest-api
- Machine-readable description: https://docs.agora.io/openapi/rtc/channel-management.en.yaml, last checked 2026-09-02
  `cauldron drift agora` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve agora     # run it
cauldron verify agora -v # check every claim
```
