# authentik

Emulates the authentik identity-provider API, for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against authentik's own generated OpenAPI 3.0.3 document, published in its own repository — `goauthentik/authentik`, `schema.yml`, 613 paths, version 2026.11.0-rc1 — read on 2026-09-06. authentik is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**A null relation means "you did not ask", not "there are none".**

A group carries four pairs of fields: `parents`/`parents_obj`, `users`/`users_obj`, `children`/`children_obj`, `roles`/`roles_obj`. The bare name is an array of identifiers and takes writes; the `_obj` name is read-only and holds expanded records.

Whether the expanded ones are filled is decided by four query parameters, and they do not agree about their defaults:

| parameter | default |
|---|---|
| `include_users` | **true** |
| `include_parents` | false |
| `include_children` | false |
| `include_inherited_roles` | false |

Three of the four `_obj` fields are declared `nullable: true` **and** listed in `required`, so they are always present and may be null. So `group.users_obj` is a populated array by default and `group.parents_obj` is null on the same response — not because the group has no parents, but because nobody passed `include_parents=true`.

A client reading `group.parents_obj.length` throws. One reading `if (group.parents_obj)` decides the group is a root. Both are wrong in the same direction, and the group may have parents.

`roles_obj` is required and **not** nullable while `inherited_roles_obj` is required and nullable — in the same object, next to each other.

**Two identifiers, and the path takes the long one.** Every group carries `pk`, a UUID string, and `num_pk`, an integer — both required. The path is `/core/groups/{group_uuid}/`, so `pk` addresses a group and `num_pk` is Django's own integer primary key, exposed and useless for fetching. The same trap repeats one level down: a group's `users` array holds integers, so members are referenced by the identifier a user's own path does not take.

**The paging block has no URLs and no nulls.** `pagination` is `{next, previous, count, current, total_pages, start_index, end_index}` — and `next` and `previous` are `type: number`, both in `required`. Not URLs, not nullable, not optional. So the usual end-of-listing test — `next == null`, or the key being absent — never fires, because the key is always there holding a number. What that number is on the last page the document does not say.

**Every listing carries a required field whose schema says nothing.** `PaginatedGroupList` requires `pagination`, `results` and `autocomplete`. `Autocomplete` is declared, in full, as `{"type": "object", "additionalProperties": {}}` — an object with no described properties. A generated client gets a required field it cannot type, on every paginated listing in a 613-path API.

**`is_superuser` is on the group and is optional.** Its own description: "Users added to this group will be superusers." The most consequential boolean in an identity provider, and it is not in the schema's `required` list — so a response may omit it, and `undefined` is falsy. Code that warns when `group.is_superuser` is true warns about nothing for a group whose response left the key out.

**One of the four declared credential schemes is not a registered one.** `authentik_device_auth` declares `type: http` with `scheme: bearer+agent`, which is not in the IANA Authentication Scheme registry that OpenAPI's own wording says the value SHOULD be. It appears on 33 operations, always beside plain `bearer`, so it is an alternative rather than a requirement — but a generator that maps `type: http` onto a library's built-in bearer helper produces the wrong header for it.

**No 401 is declared.** `GET /core/groups/{group_uuid}/` declares 400 and 403 and no 401 at all, so a client branching on 401 to refresh a token never refreshes.

## Modelling limits

- **Nothing here is verified against a live API.** authentik is self-hosted and there is no public instance.
- **Groups, listed and fetched.** 613 paths is an identity provider: users, applications, providers, sources, flows, stages, policies, outposts, tokens and the whole endpoints surface each want their own evidence.
- **`pagination.next` and `pagination.previous` are served as page numbers**, which is what the document types them as. What the provider puts there on the last page is not stated and is not guessed at here.
- **`autocomplete` is served as an empty object**, because an object with no described properties is what the document declares. Filling it would be inventing a shape.
- **The `include_*` parameters are modelled for users and parents**, which is enough to show that a null relation means "not asked for". Children and inherited roles are declared and left null.
- **`bearer+agent` is recorded and not served.** The credential checked here is the plain bearer that every one of those 33 operations also accepts.
