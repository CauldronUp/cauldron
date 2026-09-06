# Keycloak

Emulates the Keycloak API (admin), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Creating a user tells you nothing about the user you created. The documented response for `POST /admin/realms/{realm}/users` is bare "201 Created" with no body and no schema, and the actual id goes into a `Location` header the document never mentions. Every client generated straight from the spec is left holding a username and no id, in an API where every subsequent call needs the id -- the only way to find it is to read an undocumented header or list every user and search for the one you think you just made.

`enabled` also doesn't mean able to log in: a user can have `enabled: true` and still be completely locked out if `requiredActions` contains something like `UPDATE_PASSWORD`, and provisioning code that sets `enabled` and moves on creates accounts nobody can actually use. Search is prefix-based by default and the punctuation is the query language -- a bare substring matches nothing, `*foo*` does infix matching, and quoting a value does exact matching, three modes encoded entirely in asterisks and quotation marks, with the least useful one as the default. And `first`, in the same parameter list, means both "the pagination offset" and, per a different parameter's own description, a field name to match on exactly -- one word carrying two unrelated meanings on the same endpoint.

## Sources

- Documentation: https://www.keycloak.org/docs-api/latest/rest-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve keycloak     # run it
cauldron verify keycloak -v # check every claim
```
