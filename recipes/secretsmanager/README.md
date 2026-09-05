# AWS Secrets Manager

Emulates the AWS Secrets Manager API (2017-10-17), for local development and tests.

**15 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs real AWS credentials. The credential check itself was verified directly against secretsmanager.us-east-1.amazonaws.com, unsigned, on 2026-09-05.

## What this Recipe found

Checked live: a request with no Authorization header at all answers 400 `UnrecognizedClientException`, "The security token included in the request is invalid." -- the same sentence this file already had, but paired with the wrong status and code. The 403 `InvalidClientTokenId` this file had instead is a real AWS error too, just for a different failure (a signed request naming a revoked access key) that this probe could not produce without one.

A secret's ARN isn't derivable from its name: Secrets Manager appends six random characters, so `prod/db` becomes an ARN ending `prod/db-AbCdEf`, and code (or an IAM policy) that builds ARNs by hand matches nothing. `SecretString` and `SecretBinary` are alternatives, never both, so `JSON.parse(s.SecretString)` throws on undefined rather than reporting a type problem.

Deleting a secret only schedules it: it keeps its ARN and value, gains a `DeletedDate`, and reading it fails with `InvalidRequestException` rather than `ResourceNotFoundException` -- code that catches "not found" doesn't catch this, and the secret is still billed and still blocking its own name for up to thirty days. During rotation there are two versions, and `AWSPENDING` is the new value while `AWSCURRENT` is still the old one, so a read mid-rotation returns exactly what nobody expects.

## Sources

- Documentation: https://docs.aws.amazon.com/secretsmanager/latest/apireference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve secretsmanager     # run it
cauldron verify secretsmanager -v # check every claim
```
