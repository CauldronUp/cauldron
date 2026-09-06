# Provider backlog

The providers queued for Recipes, roughly in the order they are worth doing.
This is a working list, not a promise, and it exists so the queue survives
between sessions rather than living in somebody's head.

Two rules apply to everything here:

1. **A Recipe has to earn its place.** The test is whether the provider has
   behaviour worth reproducing — a third state nobody branches on, a shape that
   breaks a client, an asynchronous outcome, a failure that only exists in
   production. Linear and Attio were assessed and left out rather than
   approximated, and both were built later when the case for them was made
   rather than assumed -- which is the rule working in both directions, and
   why this paragraph names them still. A thin Recipe that raises the count is
   worse than an honest gap.
2. **Every modelling limit gets stated in the Recipe header**, not discovered
   later. Cauldron does not execute SOQL, GROQ or Algolia's search, and each of
   those Recipes says so.

## Cloud infrastructure

| Provider | Why |
|---|---|
| AWS S3 | Buckets, objects, presigned URLs, permissions, multipart upload |
| ~~AWS SQS~~ | Shipped. Visibility timeout, redelivery, dead-letter queues |
| AWS SNS | Topics, subscriptions, delivery failures |
| ~~Justworks~~ | **Assessed and refused: there is nothing to model that is not Cloudflare.** Every path on every Justworks subdomain -- `www`, `api`, `help` -- sits behind an active managed bot challenge (`Cf-Mitigated: challenge`, the "Just a moment..." page) that stops any non-browser client before it reaches Justworks' own servers. Not a connectivity problem: `status.justworks.com` resolves and answers normally from the same place. There are no public developer docs, no OpenAPI document, no repositories under the company's GitHub organisation, and no third-party client on any registry. A Recipe here would describe a bot wall, and this project does not work around those -- see `SECURITY` on why not. What would reopen this is Justworks publishing an API, not somebody trying harder |
| ~~Substack~~ | **Assessed and refused: there is no API to model.** Probed live with no account. `substack.com/api/v1/*` is real and answers, and it is the web client's own internals -- four incompatible failure envelopes across five endpoints (`403` plain text, a `401` JSON body with an HTML `<a>` tag inside the message meant for React to render, a full single-page-application shell for both a real path and an invented one, and a `404` with an empty body and no `Content-Type` at all), authenticated by session cookie, with no key any integration could hold. The one thing Substack calls a Developer API is approval-gated behind a seven-to-ten-day wait, publishes no technical reference of any kind -- no paths, no auth header name, no example JSON, only a legal document naming categories of data -- and returns a public subscriber **count** on a creator profile, never a subscriber record. The question its group was written to ask, whether a subscriber is a person or a subscription, has no vendor surface to answer it against |
| ~~AWS SES~~ | Shipped. Accepted-not-delivered, the invisible suppression list |
| ~~AWS DynamoDB~~ | Shipped. Typed attributes, omitted Items, table states |
| ~~AWS Secrets Manager~~ | Shipped. Undeducible ARNs, rotation stages, scheduled deletion |
| ~~Google Cloud Storage~~ | Shipped, written against the discovery document Google publishes. The folders come back in a different array from the files: there are no directories in a bucket, only names with slashes, and of delimiter -- "Objects whose names, aside from the prefix, contain delimiter will have their name, truncated after the delimiter, returned in prefixes. Duplicate prefixes are omitted." So a bucket of ten thousand files answers with two items and a handful of prefixes, and code reading items has a complete, accurate array that is not the answer to its question. And the paging counts both arrays at once -- maxResults is "items plus prefixes" and "fewer total results may be returned than requested" -- so the loop everybody writes stops early on the first page where two objects shared a folder. Also: the id includes the generation, so overwriting a file changes it and name is the stable one; size, generation and metageneration are digits in strings; and being deleted is two states, timeDeleted and softDeleteTime. Fetching one object by name is stated and not served -- an object name goes into the path URL-encoded and Go decodes it before the router sees it |
| ~~Google Pub/Sub~~ | Shipped. Base64 bodies, ack deadlines, delivery attempts |
| Azure Blob Storage | Containers, blobs, SAS and auth behaviour |
| Cloudflare R2 | Sits beside the existing Cloudflare Recipe |
| ~~Vultr~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`meta.total` on `/v2/regions` is the page size plus one** -- 3 for `per_page=2`, 33 for `per_page=32`, confirmed across six page sizes -- a one-page lookahead, while the other three catalogues on the same host report the real figure from the same field of the same envelope. Also pinned: ids that disagree with themselves, `os` and `applications` as JSON numbers beside `regions` and `plans` as strings; a base64 cursor that decodes to something legible (`bmV4dF9fV0FX` is `next__WAW`); and unknown ids collapsing into malformed-path 404s, which is honest here because these catalogues have no id-addressed routes at all. Stated and not served: past the true end Vultr reports `total: 0`, which the lookahead count cannot reproduce because it counts the whole matching set rather than what remains after the cursor; and `meta.links.next` is always present and empty on a final page, which no field can declare |

The queueing ones are the most valuable and the least like anything shipped so
far: at-least-once delivery, a message that comes back after a visibility
timeout, and a dead-letter queue are behaviours no fake reproduces by accident.

**XML is the blocker for several of these.** S3, SNS and Azure Blob answer in
XML, not JSON, and so do UPS, FedEx and USPS further down. Avalara and FedEx were on
that list and should not have been. AvaTax v2 is REST and JSON and the XML one
is the SOAP API it replaced; FedEx's developer platform is REST and JSON too,
and both now ship as Recipes.

The pattern that hid it: the old and new clients live side by side, and the
old one is usually the more installed. avalara/avatax is the SOAP client and
avalara/avataxclient is the REST one; php-fedex-api-wrapper wraps the SOAP web
services and has over a million installs, while the REST wrappers have tens of
thousands.

UPS and USPS were checked at the same time and the note still holds for them,
on this evidence: the widely used PHP client for UPS requires ext-simplexml,
and the USPS ones wrap Web Tools, which is the XML API. Both vendors do
document newer JSON APIs, so the ecosystem rather than the vendor is what was
measured here, and that is worth re-checking rather than trusting. Cauldron
serves JSON, so those need either XML response support in the format or an
honest exclusion. That decision is worth making once, deliberately, rather than
one Recipe at a time. The AWS services that use the JSON protocol — SQS,
DynamoDB, Secrets Manager, SES v2 — are unaffected and can go first.

## Payments and money

| Provider | Why |
|---|---|
| ~~PayPal~~ | Shipped. APPROVED is not paid, the links array is the flow, a capture can succeed without settling, and the fee comes out of the middle |
| Braintree | Transactions, customers, payment methods |
| ~~Mollie~~ | Shipped. The webhook is an id and nothing else, open is not pending, the checkout link disappears |
| Wise Business | Transfers, recipients, balances, settlement states |
| GoCardless | Mandates, direct debit, failed collections |
| ~~Brex~~ | Shipped. Money as an object, pending cards, expenses against transactions |
| ~~Ramp~~ | Shipped. Cents without saying so, pending holds, declines in the list |
| ~~Mercury~~ | Shipped. Direction in the sign, two balances, failed transfers in the list |

## Commerce

| Provider | Why |
|---|---|
| ~~Shopify GraphQL Admin API~~ | Shipped. Separate from the REST Recipe, not a replacement. Answers 200 when it refuses, and in three different shapes |
| ~~BigCommerce~~ | Shipped. Two APIs in one store: V3 wraps, V2 does not; a price is a number and the order total is a string |
| ~~Magento / Adobe Commerce~~ | Shipped. searchCriteria brackets; a product is addressed by sku and an order by entity_id; error messages are templates |
| ~~Etsy~~ | Shipped. Two credentials on every call; money is an integer over a divisor; an order is a receipt |
| ~~eBay~~ | Shipped. A dozen APIs under one host that disagree: the collection key and the total field both change per endpoint |
| ~~Amazon Selling Partner API~~ | Shipped, Orders only. Buyer name, email and address need a second token and vanish silently without it |

## Shipping and logistics

| Provider | Why |
|---|---|
| ~~ShipStation~~ | Shipped. Every date is missing its timezone, and the instant is in an account timezone the API never states |
| ~~EasyPost~~ | Shipped. A rate is a value you were given, not a price you can look up; an address can be stored and undeliverable |
| ~~ShipEngine~~ | Shipped. A rate request half-fails behind a 200: some carriers quote, some refuse, both in one body |
| ~~AfterShip~~ | Shipped. A tracking number does not identify a parcel; Delivered is not terminal |
| ~~Easyship~~ | Shipped, rating. The cheapest rate is cheapest because the customer pays the duty on the doorstep |
| ~~USPS~~ | Shipped. Three credential verdicts, and a gateway seam where a path outside the registered prefix never reaches the authorizer |
| ~~UPS~~ | Shipped. The token endpoint separates absent from wrong and the API behind it does not |
| ~~FedEx~~ | Shipped, tracking. An unknown tracking number answers 200 with the error two arrays deep, because one call can ask about thirty parcels |
| ~~DHL~~ | Shipped, Express tracking. The UTC offset is a field beside the timestamp rather than inside it |

## Productivity and identity

| Provider | Why |
|---|---|
| Microsoft Graph | Mail, calendars, users, files. Enormous surface area |
| ~~Google Calendar~~ | Shipped. Listings return the series not the occurrence, all-day events have no time, cancelled instances are almost empty |
| ~~Gmail~~ | Shipped. A listing carries no message, headers are an array, no read flag, trash is not delete |
| ~~Google Drive~~ | Shipped, written against the discovery document Google publishes. The listing tells you when it might be wrong, in a field nobody reads: incompleteSearch is "Whether the search process was incomplete. If true, then some search results might be missing, since all documents were not searched" -- a 200, a files array, a working page token, and a boolean saying the answer may be short, so a sync job terminates cleanly having seen an unknown fraction. Also: Drive is not a filesystem and says so, because a name "isn't necessarily unique within a folder"; parents is an array documented as holding at most one; and shared drives are invisible until you ask twice, since includeItemsFromAllDrives and supportsAllDrives both default to false. Stated and not served: the q query language, four deprecated parameters that still sit beside their replacements, and a page token that "should be discarded, and pagination should be restarted from the first page" if it is ever rejected |
| Microsoft OneDrive | Files, sharing, async operations |
| Microsoft Teams | Channels, messages, members |
| GitHub Actions | Could extend the GitHub Recipe, probably its own |
| ~~WorkOS~~ | Shipped. Inactive-but-present users, per-IdP raw attributes, draft connections |
| ~~Stytch~~ | Shipped. Per-factor verification, session factors, invited members |
| ~~Descope~~ | Shipped, with three cases checked against the live API and five drafted from its description. **Its two commonest error codes are undocumented**: `E011007` for a missing bearer header and `E011008` for an unparseable one -- the two mistakes a caller meets first -- appear nowhere in Descope's own error reference, which documents E011001 through E011004 and stops. The enchanted-link flow is the shape this collection collects: signing in answers **200 with no session at all**, only a pending reference, and a client treating that as a login proceeds with nothing. Three ways of getting the follow-up poll wrong collapse into one code whose own text admits it -- "probably expired by time or already clicked on". Its published description declares **no error responses anywhere**, verified by searching the whole document, so every status on a documented failure is inferred and the header says so |
| ~~FusionAuth~~ | Shipped, written against the OpenAPI document FusionAuth publishes. A header that is optional until somebody else changes something: X-FusionAuth-TenantId is declared required: false on nearly every operation, and its own description says "Only required when there is more than one tenant and the API key is not tenant-scoped". Both of those are administrative acts in a different system, so code written against a single-tenant install works, ships, and breaks the day a second tenant appears -- and a developer's sandbox is single-tenant, which is precisely the state in which this cannot be found. This Recipe seeds two. Also: an account's password can go from fine to breached without the account being touched, insertInstant and lastUpdateInstant identical while breachedPasswordStatus changes; active is a boolean and expiry is a date and neither mentions the other; and timestamps are called insertInstant and lastUpdateInstant, in milliseconds, so mapping by convention finds neither |
| ~~Keycloak~~ | Shipped, the Admin REST API, written against the OpenAPI document Keycloak publishes. Creating a user tells you nothing about the user you created: POST /admin/realms/{realm}/users is documented as "201 Created" and nothing else -- no content, no schema, no headers -- so the document describes no response body for the endpoint that makes the central object of the API, and does not say the identifier went into a Location header either. Also: enabled does not mean able to log in, because a requiredActions entry stops a login dead while enabled stays true; search is prefix-based by default with *foo* for infix and quotes for exact, so the punctuation is the API; and first is the pagination offset while the exact parameter's own description names "first" as a field to match on. The search punctuation is stated and not served |

## Analytics and flags

| Provider | Why |
|---|---|
| ~~LaunchDarkly~~ | Shipped. Per-environment state, variations as indices |
| ~~PostHog~~ | Shipped. A property is nested under properties and a flag is a string, a boolean or false, all in the same field |
| ~~Mixpanel~~ | Shipped. The blocker recorded below -- array request bodies -- was lifted as a category rather than for this provider, and serving it sharpened the finding: an array is accepted only when three fields are present, and never for their validity |
| ~~Amplitude~~ | Shipped, and the one ingest of the three that checks the key: an unrecognised one is refused by name rather than silently accepted |

## Models and inference

| Provider | Why |
|---|---|
| ~~OpenAI~~ | Shipped, written against the OpenAPI document OpenAI publishes and generates its own SDKs from. The refusal is not in the answer: content and refusal are two sibling nullable strings on the message, only one is ever filled in, and a declined request is a 200 with finish_reason stop -- so the obvious read gets null and logs an outage while the model answered in the field beside it. Also: completion_tokens counts tokens that never arrive ("like reasoning tokens, these tokens are still counted in the total completion tokens for purposes of billing"), max_completion_tokens is spent on reasoning before a word is emitted, max_tokens is deprecated and "not compatible with o-series models", store defaults to false so the completion you just made cannot be fetched, and content_filter is a stopping reason rather than an error. Streaming and the Responses API are stated as not modelled |
| ~~Anthropic~~ | Shipped. One route resolves an identifier's shape and existence before the credential, and its sibling on the same host does not |
| ~~Google Gemini~~ | Shipped, and the section near the end of this file that recorded it as unservable is now the story of how it was unblocked. A blocked prompt is a 200 with the candidates taken away -- promptFeedback.blockReason is "If set, the prompt was blocked and no candidates are returned" -- so candidates[0] throws where OpenAI's refusal hands back a null. Serving the contrast needed one path to answer two shapes chosen by the request body, which selects could not do, so selects_body was added beside it as a separate field. Also: an empty finishReason means the model "has not stopped generating tokens", from one schema serving the streaming and non-streaming calls; totalTokenCount is documented as "prompt + thoughts + response candidates"; and promptTokenCount "is still the total effective prompt size" when content is cached |
| ~~Pinecone~~ | Shipped. The response to "make me an index" is an address somewhere else. Pinecone is two APIs with two base URLs: the control plane at `api.pinecone.io` creates and lists indexes, and the data plane -- query, upsert, fetch, which is every operation an application actually performs -- lives at the index's own host, which you cannot know until the control plane tells you. The data plane's own OpenAPI document says so by being unable to name its server: `url: https://{index_host}` with the variable's `default: unknown`, so a client generated from it and left unconfigured posts to `https://unknown/query`. `host` is required on every IndexModel because there is no other way to find the thing you just made. Also pinned: `ready` and `state` are two fields and the vendor's own example has them disagreeing -- `{"ready": true, "state": "ScalingUpPodSize"}`, ready and not in the state called Ready, out of nine states; a create answers 201 with `host` filled in and `ready` false, so the address exists before the index does; a delete is a **202** ("The request to delete the index has been accepted") and `Terminating` is one of the nine states, so the index you deleted answers describe for a while afterwards; a sparse index has no `dimension` at all rather than a null one; and the error vocabulary is gRPC's -- the seventeen canonical status codes with FORBIDDEN, UNPROCESSABLE_ENTITY and PAYMENT_REQUIRED bolted on, which is why **`OK` is a documented value of a field called "The error code"** -- with the HTTP status repeated inside the body. `X-Pinecone-Api-Version` is `required: true` on every operation in both documents. Stated and not served: the data plane itself, because it is a second base URL and this emulator is one address, and serving `/query` here would teach the exact mistake the Recipe describes |
| ~~Replicate~~ | Shipped. A created prediction has no output property at all rather than a null one, succeeded is not the same as produced something, a cold start is a minute with no signal but a boot_time, and the output is a link to a file deleted after an hour |
| ~~Hugging Face~~ | Shipped, and written entirely from live responses: no token, no account, every case dated. **A model that does not exist answers 401, not 404** -- `{"error": "Invalid username or password."}` with a `WWW-Authenticate` header -- so a typo reads as a credentials problem and an anonymous client cannot tell a wrong name from a private one. **Two identifiers per record and the primary-key-shaped one opens nothing**: every listing carries `_id`, a Mongo ObjectId, beside `id`, the real slug, and fetching by `_id` gives the same 401 as an invented name. Also pinned: an invalid `sort` whose body says `✖` (U+2716) while the `X-Error-Message` header says a plain ASCII `*` for the same failure; a percent-encoded slash explicitly refused, closing npm's escape hatch for scoped names; `limit` silently clamped to 1000 whether asked for 0, -1, 999999 or nothing; and `trendingScore` vanishing from every record the moment an explicit `sort` is given. Stated and not served: a real two-segment repo cannot be routed, because the part after the slash appears in no field of the wire format and inventing one would put it on the wire |
| ~~ElevenLabs~~ | Shipped, written against the OpenAPI document ElevenLabs serves from its own API host. The credential is optional on every endpoint that needs it: there is no securitySchemes block anywhere, xi-api-key is declared as an ordinary header parameter on 383 operations with required: false and a string-or-null schema, and its description hedges -- "required by most endpoints". Of 387 operations, 385 document exactly one failure (422) and two document a 401, so a generated client has a branch for a malformed request nearly everywhere and none for a wrong key almost anywhere. Also: the main endpoint of the product answers audio/mpeg rather than JSON; the history "does not include" music and SFX and they "cannot currently be retrieved via the API"; a generation's cost is the difference between two meter readings on the history item; and allowed_to_extend_character_limit is deprecated in favour of an inequality, `max_credit_limit_extension != 0`, on a field whose value may be the string "unlimited" |
| ~~Deepgram~~ | Shipped. Four-level nesting, seconds against milliseconds, err_code |
| ~~AssemblyAI~~ | Shipped. Failure at HTTP 200, present-and-null fields, ms against s |

Streaming is the open question here. Cauldron serves whole responses, and
server-sent events are a different shape. Whether that is modelled, or stated
as a gap, needs deciding before the first of these ships rather than after.

## Data and messaging platforms

| Provider | Why |
|---|---|
| ~~Supabase~~ | Shipped. Auth, storage, database REST and realtime |
| ~~Redis Cloud~~ | Shipped, and not what this row imagined: the control plane, not the data plane. Half a credential answers a bare nginx 500 |
| ~~Upstash~~ | Shipped. Three different requests to the same endpoint with no valid credential answer three different 401s -- an empty body with no Content-Type, a JSON-Web-Token parser's own error text for a bearer token, and `{"error": "Unauthorized"}` for the documented Basic scheme -- and this format can only reproduce one of them at a time. Also: `region` is a one-member enum, `"global"`, so it says nothing about where a database actually runs, which is `primary_region` instead; `db_acl_enabled` is a string `"true"`/`"false"` beside ordinary booleans on the same object; and a database delete answers 200 with the bare JSON string `"OK"`, not an object |
| ~~MongoDB Atlas~~ | Shipped, the Administration API, written against the OpenAPI document MongoDB publishes. The version of the API you are talking to is a date inside a content type: GET /groups/{groupId}/clusters documents one 200 with three content types under it -- vnd.atlas.2023-01-01+json, 2023-02-01+json and 2024-08-05+json -- resolving to three separately named schemas, all current. And the difference is not cosmetic: the legacy view carries mongoURI, srvAddress, diskSizeGB and providerSettings and the newest carries none of them, so the field an application connects with is present under one date and absent under the next, same URL, same credentials, same instant. Also: the schema names carry the date (ClusterDescription20240805), failure is better documented than success (more operations describe a 401 than a 200), and the credential is HTTP Digest -- challenge-response, on a cloud API written in the 2020s. Digest is stated and not served, and so is the Content-Type echo |
| ~~Neon~~ | Shipped. Branches, databases, endpoints |
| ~~PlanetScale~~ | Shipped, written against the Swagger 2.0 document PlanetScale serves from its own API host. The field called state is not the state of the deployment: a deploy request carries state -- "Whether the deploy request is open or closed" -- beside deployment_state, "The deployment state of the deploy request", so a request abandoned after review and one whose migration ran an hour ago both read closed. Also: the id is "The ID of the deploy request" and every path takes the number instead, which is per database, so the globally unique identifier addresses nothing; a request outlives the branch it came from, still naming it while branch_deleted says it is gone; and next_page is "null when this is the last page" where Confluence omits the field entirely. The ten deploy-lifecycle endpoints are stated as not modelled, because each advances a state machine this format cannot express |
| ~~CockroachDB Cloud~~ | Shipped as `cockroachdb`. The state enum has no word for deleting |
| ~~Turso~~ | Shipped. **It answers a missing header with a JWT parsing error** -- "token contains an invalid number of segments" when no token was sent at all, so the message describes the shape of a credential the caller never supplied, and a junk string answers identically. Its create is a different failure from Temporal's: not an unfinished operation but a **truncated projection**, three fields where a read gives eight, and among the missing is `group`, which the create request was required to supply. There is no state field in the schema to be pending. Also pinned: `DbId` is a real UUID that nothing addresses by, since every path takes the name; a delete answers a bare string under the key a create wraps an object under; and **no 405 exists on this host at all**, where SingleStore's Recipe records a real one with an `Allow` header |
| Confluent Cloud | Kafka topics, schemas, consumers |
| Aiven | Managed databases, Kafka, service lifecycle |
| CloudAMQP | RabbitMQ instances and queues |

## Messaging and notifications

| Provider | Why |
|---|---|
| ~~Resend~~ | Shipped. last_event as the only history, restricted keys |
| Amazon Pinpoint | Messaging and delivery events |
| ~~OneSignal~~ | Shipped. 200 with errors, recipients as a ceiling, dead devices |
| ~~Pusher~~ | Shipped. Keyed channel map, empty publish response, text errors |
| ~~Ably~~ | Shipped. 204 publish, newest-first history, capability as a string |
| Bird | SMS, WhatsApp, conversations |
| ~~Sinch~~ | Shipped. **Its documented API host has no public address**: `api.sinch.com` resolves through public DNS to three addresses on 10.65.0.0/16 -- private, unroutable -- and every connection times out, while the SMS product actually answers on `us.sms.api.sinch.com`. Where it does answer, a missing token and a wrong one are identical: 401, zero bytes, no `WWW-Authenticate`. The batch model is the finding: a send answers with the batch and **no status field anywhere**, for the batch or any recipient, and the default delivery report is aggregated counts by status code with no phone numbers in it -- so the default answer to "did my messages arrive" is a histogram. The batch's `id` is renamed `batch_id` on that report, the same rename Bandwidth performs from the other direction |
| ~~Telnyx~~ | Shipped. Per-recipient status, cost after the fact, silent MMS |
| ~~Infobip~~ | Shipped, and it finishes a five-way comparison the SMS Recipes have been building. One request to many recipients gives back nothing shared at all from Plivo, one bare id from Bandwidth and Sinch, and N separate requests from Twilio -- **Infobip alone answers one entry per recipient in the first response, each already carrying a structured status object**, and its own status reference documents two recipients in one request coming back in different groups, one pending and one already rejected. Its failures are less generous: a missing credential and a well-formed wrong one are identical, and a GET where a POST belongs answers **404 rather than 405**, naming the path it did not find, so a wrong method is reported as a wrong URL. Recorded and out of scope: its successor endpoint at `/sms/3/messages` uses a completely different error envelope from this one -- one vendor, two generations, no shared field |
| WhatsApp Business Platform | Messages, templates, delivery states |

## Documents and signatures

| Provider | Why |
|---|---|
| ~~Documenso~~ | Shipped. Document against recipient status, tokens, CC recipients |
| ~~Dropbox Sign~~ | Shipped, written against the OpenAPI document Dropbox Sign publishes and generates its own SDKs from. There is no status: a signature request carries is_complete, is_declined and has_error, three independent booleans, so is_complete: false is four situations wearing one face -- nobody has signed, somebody declined, something broke, or two of three are done. A dashboard rendering complete-or-pending shows a declined contract as pending for ever. Also: the person who signs need not be the one you asked (reassigned_by, reassignment_reason and reassigned_from on the signature that replaced yours); every signer carries an error of its own beside the request's; and a test request that "has no legal value" is in the listing with the contracts, which takes account_id, page, page_size and query and nothing else. The most surprising thing about this API is stated and not served: the document defines EventCallbackResponse as text/plain defaulting to "Hello API Event Received", which is the body your callback endpoint must answer with -- every other provider here reads your status code, this one reads your prose |
| Adobe Acrobat Sign | Agreements and the signing lifecycle |
| Ironclad | Contracts and approval workflows |
| DocSpring | PDF generation and templates |

## Tax, accounting and payroll

| Provider | Why |
|---|---|
| ~~TaxJar~~ | Shipped. Decimal rates, present-and-zero jurisdictions, nexus |
| ~~Avalara~~ | Shipped, AvaTax REST v2. A quote and a record are the same call with one word different, and only one of them is ever filed |
| Stripe Tax | Probably a specialised extension of the Stripe Recipe |
| FreshBooks | Invoices, clients, expenses |
| Sage | Accounting, contacts, invoices |
| Zoho Books | Invoices, payments, contacts |
| Wave | Accounting and invoicing |
| ~~Bill.com~~ | Shipped. Failure at HTTP 200, isActive as "1"/"2", scheduled payments |
| ~~Gusto~~ | Shipped. Processed against paid, four dates, unpayable employees |
| ~~Deel~~ | Shipped. in_progress means running, mixed currencies, post-termination invoices |

## Media

| Provider | Why |
|---|---|
| ~~Cloudinary~~ | Shipped. url against secure_url, 200-with-not-found, moderation |
| ~~Ghost~~ | Shipped. Plural-wrapped singles, sent-not-published, Lexical |
| Imgix | Image transformations and delivery |
| ~~Mux~~ | Shipped. Asset against playback id, preparing state, array messages |
| ~~Vimeo~~ | Shipped, written against the OpenAPI document Vimeo publishes at vimeo/openapi. The privacy setting called unlisted is described as a fact about search: the video privacy enum describes six of its seven values as rules about who may watch and the seventh as "Not searchable from vimeo.com" -- while the user privacy enum in the same file says of the same value "Anybody can view the user's videos if they have a link". One value, two descriptions, and only the one attached to the field you do not set tells you who can see the thing. Also: two fields called status answering different questions, one of them "deprecated in favor of `upload` and `transcode`", so the migration from one field is to two; a video with no identifier at all, only a uri to cut a number off; and a privacy value that changes which content types the request may use -- "When this value is `users`, `application/json` is the only valid content type" -- which is stated and not served. Detection found the largest near miss in the collection: vimeo/psalm, the PHP static analyser, at eighty-five million downloads, which is Etsy's phan a second time |
| ~~YouTube Data API~~ | Shipped, written against the discovery document Google generates its own client libraries from. The identifier of a search result is an object, and none of its fields is called id: a video's `id` is a string -- "The ID that YouTube uses to uniquely identify the video" -- while a search result's `id` is a ResourceId carrying `kind`, `videoId`, `channelId` and `playlistId`, each "only present if the resource is" that kind. There is no `id.id`. Code that reads `result.id` gets an object, code that reads `result.id.videoId` gets undefined for every channel in the results, and neither is an error. Also: a search returns five things unless you ask for more, which is small enough that a developer eyeballing the response concludes the search matched five things; three `kind` fields to read one identifier; counts as digit-strings. And the document describes its own results array as pagination -- the description on `items` in SearchListResponse is "Pagination information for token pagination", the sentence belonging to the field beside it, a copy-and-paste mistake sitting in the file Google generates its client libraries from. Quota is stated and not served: it is a 403 whose reason is quotaExceeded rather than a 429, so retry logic written for rate limits does not fire |

## A note on the sections below

Everything from here down was added in one sweep looking for services with
public APIs worth emulating, and it has not been through rule 1 yet. An entry
whose reason starts with **Assess** names a surface rather than a behaviour:
somebody still has to read the documentation and decide whether there is a lie
in there worth reproducing, or whether it belongs in the not-done table
instead. Saying which entries have been assessed and which have not is cheaper
than finding out halfway through writing one.

## Telephony and contact centre

| Provider | Why |
|---|---|
| ~~RingCentral~~ | Shipped. Sent is not delivered, per-group rate limits, deleted messages still in the store |
| Aircall | Assess — calls, contacts, webhook ordering. A call's final state arrives after the call ends, so the webhook and the fetched object disagree for a while |
| Dialpad | Assess — calls, SMS, transcripts |
| OpenPhone | Assess — the API is young enough that shapes have moved between versions, which is worth pinning |
| 8x8 | Assess — messaging and voice |

## Issue tracking, docs and planning

| Provider | Why |
|---|---|
| ~~Jira Cloud~~ | Shipped. Numbered custom fields, administrator-typed statuses, ADF descriptions, a create that returns no issue, and a 410 where the old search used to be |
| ~~Jira Service Management~~ | Shipped, written against the OpenAPI document Atlassian publishes for the Service Desk API. This row's guess was right and understated: elapsedTime is "Duration of the service", and withinCalendarHours means it only counts during working hours, so a ticket raised at five on Friday reports ninety minutes of elapsed time on Monday -- `now - startTime` is the number every dashboard computes and is not the number the breach is decided on. paused stops it entirely. And the field is not in the response unless you ask: sla is "Expandable details", so a plain GET answers 200 with no indication of whether the ticket is in trouble. Also: a request has "zero or more" SLAs each with "zero or 1" ongoing cycle, so nothing answers "is this breaching"; breachTime is populated on cycles that never breached ("would have breached in case of non-breached completed cycle"); paused is on the ongoing cycle and not the completed one, so the record that a clock stopped does not survive it; and every instant arrives four times, one of them a `jira` field defined as "ISO8601 format but extended with milliseconds" |
| ~~Confluence Cloud~~ | Shipped, written against the OpenAPI 3.0.3 document Atlassian publishes. This row's claim is true and is the second half. The first half is that asking for a page does not get you the page: body-format has no default, and "the representation will be available under a response field of the same name under the `body` field" only if you named one -- so an ordinary fetch returns a title, a status, a version, an author and nothing of what the page says, at 200. And the format you name becomes the key, so the same sentence is at body.storage.value for one request and body.atlas_doc_format.value for the next. Also: no partial update (id, status, title, body and version all required, so every updater must read before it writes); two documented ways to lose a draft, neither an error; the default listing includes archived pages; spaceId is accepted on an update and does nothing; and the id is int64 in the path and a string in every response |
| ~~Smartsheet~~ | Shipped, and written entirely from live responses: no token, no account, every case dated. **Four routes disagree about what a missing credential is.** No credential on a collection is 403 `errorCode 1004`; a merely wrong token is 401 `errorCode 1002`, verified with three different garbage tokens; and `/sheets/{id}` and `/workspaces/{id}` answer a missing credential with an entirely empty 401 -- zero bytes, no Content-Type, no code -- while `/users/{id}` and every collection answer 403 with a body. Same product, same missing header, and whether you get a diagnosis depends on which noun you asked for. Also pinned: a malformed path as 404 `errorCode 1006`, and a wrong method as 405 `errorCode 1122` with the offending method named inside the message prose rather than in a field. Stated and not served: Smartsheet pretty-prints with a space before the colon and announces `application/json;charset=UTF-8` in that casing, neither reproduced byte-for-byte |
| Wrike | Assess — tasks, folders, custom fields |
| Height | Assess — tasks and lists |
| Productboard | Assess — features, notes, insights |
| Canny | Assess — posts, votes, changelog |

## Banking rails and money movement

Same standing rule as Mercury, Brex, Bill.com, Gusto and Deel: an emulator for
an API that moves money models the reads and leaves the transfers alone, and
the header says so.

| Provider | Why |
|---|---|
| ~~Increase~~ | Shipped. A return is a new object days later and the transfer stays unmarked |
| ~~Modern Treasury~~ | Shipped. Every amount appears twice and cancels, and one direction means the opposite sign |
| ~~Column~~ | Shipped. A notification of change is not a failure, and R01 and R07 are two different problems |
| ~~Dwolla~~ | Shipped. A create answers 201 with no body and the id only in a Location header, nothing has an id field at all, and micro-deposit verification has three attempts before the funding source is permanently unverifiable |
| Unit | Assess — accounts, cards, authorisations. An authorisation is not a transaction and the two have separate ids |
| ~~Marqeta~~ | Shipped. An authorization and its clearing are two transactions for different amounts; three balances and only one is spendable; the PAN never leaves. JIT funding is not modelled and the Recipe says why |
| ~~Checkout.com~~ | Shipped. Every credential failure is 401 with zero bytes and no Content-Type |
| Authorize.Net | Assess — the XML-shaped API and its own result codes, which are not HTTP statuses |
| ~~Klarna~~ | Shipped. Its own reference tells integrators not to rely on the status field it publishes |
| Affirm | Assess — checkout, capture, partial refunds |
| ~~RevenueCat~~ | Shipped. Assess — mobile subscriptions and entitlements, where the entitlement is the thing an app reads and the subscription is the thing that renews, and they can disagree for a whole billing period |
| Firebase Auth | Assess — identity, and the emulator Google already ships for it is the question: a Recipe has to be better than the official one to earn its place |
| Coinbase Commerce | Assess — charges that expire, underpayment and overpayment as distinct outcomes |
| Circle | Assess — USDC transfers and their settlement states |

## Feature flags and experimentation

| Provider | Why |
|---|---|
| Statsig | Assess — gates, experiments, exposure logging |
| ~~GrowthBook~~ | Shipped. The field everybody branches on has no declared values and the field beside it does: an Experiment requires both `type` and `status`, `type` is `enum: [standard, multi-armed-bandit, holdout]`, and `status` is `{"type": "string"}` and nothing else -- no enum, no example, no description, in a 2.7MB document. `resultSummary` then carries a **second** bare string called `status` meaning something different, on the same object. Also pinned: **six fields to answer one question**, since a listing's envelope carries `limit`, `offset`, `count`, `total`, `hasMore` and `nextOffset` -- `count` is this page and `total` is everything, `hasMore` is a boolean, and `nextOffset` is `anyOf [integer, null]` so it is present and null on the last page; two identifiers where the useful one is not the id, since `trackingKey` is what reaches the analytics warehouse (the third provider here with that shape, after Grafana's `uid` and ConfigCat's `key`); an `owner` field holding "The userId of the owner (or raw owner name/email for legacy records)" beside an `ownerEmail` sent only "when the owner can be resolved to a known user", so neither field alone says which case you are in; and three version numbers -- `hashVersion` as `anyOf [const 1, const 2]`, plus `bucketVersion` and `minBucketVersion` -- that between them decide which variation a user sees. Detection repeats ConfigCat's split and settles it as a rule about the category: a flag platform has a delivery surface SDKs read and a management API everything else uses. `@growthbook/growthbook` has 3.49M downloads and the PHP SDK 3.92M, and neither calls this API; exactly one package does, and says so -- `growthbook`, "Command-line interface for the GrowthBook REST API", at 110k. Where the split is by host, as here and at ConfigCat, it is two Recipes or none; Unleash's client and admin APIs share a host, which is why one Recipe covered both there |
| ~~Flagsmith~~ | Shipped, and it **answers three failures in three formats with no field in common**: a missing environment key gets a printed Python tuple -- parentheses, single quotes, a comma -- served under a JSON content type, so `.json()` throws while the header promises otherwise; a wrong key gets a **404 with zero bytes**, a not-found status for a credential problem; and an unrouted path gets a bare top-level array holding one object, so `body.title` is undefined and a client must index first. Missing and wrong are both modelled here, which ten other Recipes could not manage, because they run through different mechanisms -- one checks the header is present, the other that it matches |
| ~~Split~~ | Shipped, with seven cases checked against the live API and two drafted from Split's own per-page description fragments. **Its failure sentence echoes the last five characters of the key you sent**, unmasked, inside the message explaining that the key was rejected -- and the mask before them is a fixed run of x and X regardless of the real key's length, so the width says nothing while the tail says something. That is on api.split.io; sdk.split.io does the opposite and tells nothing apart, answering a missing credential and five different wrong ones byte-identically with `{"code":400,"message":"","details":null,"transactionId":null}` -- four fields, three of them carrying nothing, and the message an empty string rather than absent. `transactionId` alone takes four states across the two hosts: null, absent, empty string, and a real value. A third host, auth.split.io, separates missing from wrong by status code and is left unrouted rather than half-modelled |
| ~~ConfigCat~~ | Shipped, and the row's third item is now a separate piece of work -- see below. **The Public Management API** is what ships. Two paths that differ by one character answer with different documents: `/v1/settings/{settingKeyOrId}/value` gives `SettingValueModel` -- `value`, `rolloutRules`, `rolloutPercentageItems` -- and `/v2/...` gives `SettingFormulaModel` -- `defaultValue`, `targetingRules`, `percentageEvaluationAttribute`. The field holding the value is called `value` in one and `defaultValue` in the other, and it is wrapped in the second; the two rule arrays become one. A client that upgrades the URL and not the reader gets undefined from a 200. Also pinned: **two credentials in one request doing different jobs** -- Basic in `Authorization` says who you are, `X-CONFIGCAT-SDKKEY` says which config and environment you mean, and only 6 of the API's 76 operations take the second; **one path segment holding two identifier spaces**, since `settingKeyOrId` is "The key or id of the Setting" and a setting has both a human-chosen `key` and an int32 `settingId`, so a flag keyed "42" collides with setting number 42; "may I save this" having four answers on v2 (`readOnly`, `approveRequired`, `canBypassApproval`, `reasonRequired`) and one on v1, so the same flag looks editable on one path and gated on the other; and a `value` whose JSON type follows `settingType` (boolean, string, int, double), with `isJson` beside it for the strings that are documents. **The config-delivery surface is the remaining work and it is where the users are**: ConfigCat SDKs fetch a static configuration file from the CDN and evaluate locally, never touching `api.configcat.com`, and `configcat/configcat-client` has 1.29M downloads against about 5k for the clients that speak the Management API -- a ratio of roughly 260 to 1. Different host, different shape, different job |
| ~~Unleash~~ | Shipped, and the row was right that the APIs are different shapes -- it is worse than that. Three APIs expose the same flag, each schema has an `enabled` boolean, and in none of them does it mean "this flag is on". `frontendApiFeatureSchema.enabled` is documented as **"Always set to `true`."** -- the Frontend API returns only the flags that already evaluated on, so the field is a constant and the real answer is whether the object is in the array at all, which makes `toggles.find(t => t.name === x).enabled` a question whose answer is always yes when it does not throw. `clientFeatureSchema.enabled` is one input: "This is ANDed with the evaluation results of the strategies list", and the strategies are themselves "evaluated and ORed together", so the rule is an AND over an OR and it exists nowhere except as two clauses in two field descriptions. And the Admin API's `enabled` is a summary of an `environments` array whose members disagree with it -- on in development, off in production, one boolean at the top. Also pinned: the two SDK APIs do not agree on what to call the collection (`toggles` on one, `features` on the other, same flags, same server, because Unleash renamed the concept and one endpoint moved); the response format carries its own version number in the body, where version 2 "includes segments as a separate array" so an evaluation that reads strategies without joining segments applies a different rule while looking complete; and `description` and `variants` arrive as null rather than absent. Two authentication schemes share one header -- `apiKey` raw and `bearerToken` behind "Bearer " -- and only the raw form is served, which is pinned as a refusal so the gap is visible. Detection found the shape next door to YouTube's simile: matched on an imperative, where the biggest npm result for "unleash" is `precinct` at twenty million downloads, a dependency parser described as "Unleash the detectives" |

## Observability and product analytics

| Provider | Why |
|---|---|
| ~~Bugsnag~~ | Shipped. An error is not its occurrences, and the counts are on the error |
| Honeybadger | Assess — faults and notices |
| Grafana Cloud | Assess — the stack-management API, which is a second surface on top of the one that now ships. **The Grafana HTTP API is shipped** as `grafana`, written against `api-merged.json` in `grafana/grafana`: two unique identifiers for one dashboard, and the one called `id` is the one you cannot use. A save's required fields are `["status", "title", "version", "id", "uid", "url"]` and the two identifiers carry the same sentence -- "The unique identifier (id)" and "The unique identifier (uid)" -- for an int64 and a string, while the only path that fetches a dashboard is `/api/dashboards/uid/{uid}`. So the integer is required in every save response and can address nothing; it is the instance's row number, and a deploy that stored it stored the identifier that will not survive the move. The document says the same about folders outright: `folderId` is "Deprecated: use FolderUID instead", beside `folderUid`. Also pinned: a field called `title` whose description is **"Slug The slug of the dashboard."** with the example "my-dashboard" -- a struct comment one field out of place, on the field a UI shows back to the user; the `version` that guards the next save living in `meta` rather than on the dashboard, which is free-form JSON, with a **412** declared on the save; five permission booleans (`canAdmin`, `canDelete`, `canEdit`, `canSave`, `canStar`) and a sixth field, `provisioned`, that none of them accounts for; and a search that answers with a bare array, no envelope and no count, whose hits can be dashboards already in the bin. Detection is the thinnest in the collection and the reason is the finding: eighteen npm results for "grafana" and not one calls this API -- the whole `@grafana/` scope is plugin tooling, and the biggest Packagist numbers are Loki, a different product from the same vendor |
| ~~Honeycomb~~ | Shipped. **One endpoint answers `problem+json` and six do not** -- `/1/auth` sends RFC 7807 while boards, datasets, columns, triggers, markers and the events ingest answer the same credential failure as a plain `{"error": ...}`. Written partly to find out whether RFC 7807 needed a new error style: it did not, a flat style with the message under `title` reproduces it exactly |
| ~~Better Stack~~ | Shipped. Its 404 hardcodes the word GET whatever method was sent, and it is the only one of its group with both a last-checked time and paused as its own status word |
| ~~Heap~~ | Shipped. The application identifier authenticates nothing -- one that was never registered gets the identical success a real one would |
| LogRocket | Assess — sessions and issues |
| ~~New Relic~~ | Shipped, with five cases checked against the live API and five drafted from its published description. **Its own description would generate a broken client**: the spec declares the applications list with the same schema as the single fetch, promising `{"application": {...}}` where the wire sends `{"applications": [...]}`, so generated code reads `response.application.id` and gets undefined on every list call. Also pinned: a wrong key echoed back inside the failure body, and `health_status` as the string "gray" beside `reporting` as a boolean -- two fields stating one fact in two types. NerdGraph was probed for the 200-with-errors behaviour this collection collects and does **not** do it unauthenticated: it 401s before GraphQL executes |

## Hosting, deployment and package registries

| Provider | Why |
|---|---|
| ~~Netlify~~ | Shipped. A deploy id exists long before the site is live |
| ~~Render~~ | Shipped, written against `render-public-api-1.json`, which is not at the well-known `render.com/openapi.json` -- that URL is a different product, an AI docs-search API, and the real description sits three redirects deeper at `api-docs.render.com`. A listing does not answer with the records: `serviceList` and `deployList` are both arrays of `{cursor, service}`/`{cursor, deploy}`, so the record is one level down under its own name and a client reading `response[0].id` gets undefined. Triggering a deploy answers 201 from an eleven-value status enum, long before any of them is `live`. And three boolean-shaped fields on one resource are spelled three ways -- `suspended` is `"suspended"`/`"not_suspended"`, `autoDeploy` is `"yes"`/`"no"`, `notifyOnFail` is a three-way `"default"`/`"notify"`/`"ignore"` -- and none of them is `true` or `false` |
| ~~Fly.io~~ | Shipped, and the comparison this row asked for is in the file: the Machines API is served and the older platform API is recorded as a stated gap, because it is GraphQL at a different host and is a second Recipe rather than more of this one. The headline is that `started` does not mean your application is up -- a machine has a `state`, everybody reads it, and three other fields on the same object independently contradict it. `host_status` can be `unreachable` (started, on a host Fly cannot currently talk to -- not stopped, not failed, not anything `state` has a word for), `cordoned` can be true (started and taking no new traffic, on purpose, and nothing in `state` says so), and `checks` can be `critical` (the process is running and failing the health check that exists to describe exactly that). Four answers to "is this machine up" on one object, disagreeing by design |
| ~~Heroku~~ | Shipped. Assess — the API still uses `Accept: application/vnd.heroku+json; version=3`, so a missing header is a different response rather than an error |
| ~~Linode~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. Written the same day as the Vultr Recipe and meant to be read beside it -- **two competitors, one job, nothing in common at the envelope**. Linode wraps everything under a generic `data` with `page`/`pages`/`results`; Vultr names each collection after itself and pages by cursor. Their unauthenticated failures share no field at all: a `{"errors": [{"reason"}]}` list against a flat object with a status inside it. Linode is at least consistent about ids, which are strings everywhere. Also pinned: a GPU type whose `price.monthly` is null while its hourly price is a number, and timestamps that are not RFC 3339 -- `2025-10-01T04:00:00`, no zone, no `Z`, so a strict parser rejects them and a loose one guesses. Stated and not served: `/images` and `/linode/kernels` are addressed by two-segment ids like `linode/almalinux10`, and the router matches an exact segment count, so both are list-only here |
| ~~Hetzner Cloud~~ | Shipped, and this row asked for the sentence the Recipe opens with. Powering off a server does not answer with the server: it answers with an action -- a job, with a status of `running`, a progress percentage and a `finished` timestamp that is null. The machine is still on. Whether it will go off is a second request, to a different endpoint, about a different object, and there is no version of the call that waits. Every mutation in this API is that shape -- attach a volume, change a type, rebuild, enable rescue, request a console, each answers 201 with a job -- so nothing that changes anything answers with the thing it changed, and a test that powers a server off and asserts it is off passes against a mock, passes against a fixture, and is a race against the real thing. Also pinned: an action that fails long after its successful 201, an `error` that is an object rather than a sentence, a `progress` that reaches 100 on a job still running, nine server statuses of which "not running" is eight, `locked` as a field separate from `status` refusing the next mutation with a 423, and a root password returned exactly once |
| ~~Scaleway~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A path that does not exist answers 401** -- byte-identical to a real authenticated endpoint -- because the credential gate runs before anything checks whether the route exists, so an anonymous client cannot tell a typo from a permissions problem anywhere on the authenticated surface. The public catalogue gives three different failures for what a caller experiences as one mistake: an invalid zone is 400, a misspelled product is 404 with a `type`, and the same path minus the zone segment is 404 without one. **The regional question has an unglamorous answer**: the catalogue shape does not vary by zone at all -- identical seventeen keys everywhere -- what varies is membership, and the smaller zone is a strict subset, confirmed by diffing full catalogues. Stated and not served: Scaleway's Link header carries `rel="first"` and `rel="last"` and spells the back-relation `"previous"` where this format emits `"prev"`, so the previous link stays off rather than going out byte-wrong |
| ~~Docker Hub~~ | Shipped. The rate limit is counted per IP, not per token |
| ~~npm registry~~ | Shipped. Unpublished versions leave a tombstone, and dist-tags are a flat map |
| ~~PyPI~~ | Shipped, and the row was right that upload is a separate protocol -- the JSON API is what ships, and it is the first Recipe in a while whose **every case was checked against the live API**, because the endpoint is public and needs no credential. Three download counters, on every package, that always say minus one: `"downloads": {"last_day": -1, "last_month": -1, "last_week": -1}`. Warehouse stopped serving download statistics through this API and kept the field, so the shape is intact, the keys are present, the types are right and the numbers are a sentinel -- and every file in `urls` carries a fourth one, also minus one. Also pinned: **the biggest field in the response disappears one path segment along** (`/pypi/{project}/json` carries `releases`, 163 versions for requests, and `/pypi/{project}/{version}/json` does not carry it at all); a version can exist with nothing to install, since three of requests' releases map to empty arrays, so counting the keys of `releases` counts versions nobody can download; **one key in a file entry is hyphenated** -- `core-metadata`, beside `comment_text`, `md5_digest`, `python_version` and `upload_time_iso_8601` -- so destructuring works for every sibling and not for that one; the same hash arriving twice (`md5_digest` at the top and `digests.md5` inside) and the same instant arriving twice at two precisions (`upload_time` with no timezone, `upload_time_iso_8601` with milliseconds and a Z); `has_sig`, a field for GPG signatures PyPI stopped accepting in 2023; and two licence fields where only the older is filled in. Detection carries the largest number this collection has excluded, and it is not a homonym: `composer/installers` at **148 million** is the top Packagist result for "pypi" and is a Composer library installer. Coverage here is structurally thin rather than accidentally so -- the obvious clients of PyPI are Python packages, and this detector reads Composer, npm and Go modules |
| ~~crates.io~~ | Shipped, and written entirely from live responses: no OpenAPI, no credential, so every case carries the date it was checked. **Four fields say which version is the latest and on a real crate they disagree.** A crate carries `max_version`, `newest_version`, `max_stable_version` and `default_version`; on `serde` all four read 1.0.229, which is the ordinary case and the reason nobody notices. On `rand`: `max_version` 0.10.2, **`newest_version` 0.8.8**, `max_stable_version` 0.10.2, `default_version` 0.10.2 -- newest means most recently published and max means highest by semver, so a patch to an old line published on 25 August against a 0.10.2 from 2 July leaves the field named newest two minor lines behind. Also pinned: **a request with no `User-Agent` is a 403** that has nothing to do with permission ("we ask that your user agent actually identify your bot, and not just report the HTTP client library you're using") on a registry that needs no credential at all, so the first thing a hand-rolled client meets is a refusal about attribution; `id` is the crate's name on a crate and an integer on a version, in one response; one link in the `links` object is `null` while its five siblings are paths; and `versions` means whole records at the top level and integer ids inside `crate`. Detection carries two near-miss kinds at once -- a homonym one letter away, since every Packagist result for "crates" is **CrateDB at crate.io**, and a kind this collection had not recorded: **matched on a README badge**, where `ts-gettext-extractor` and the Tauri plugins surface because their READMEs carry a shields.io badge whose URL contains crates.io |
| ~~RubyGems~~ | Shipped, and written entirely from live responses: no OpenAPI, no credential for reads, every case carries the date it was checked. **The flag that decides whether publishing needs two-factor authentication is the string `"true"`** -- a gemspec's metadata is declared as a map of strings to strings and the API hands it back as written, so `if (metadata.rubygems_mfa_required)` is true for `"true"` and equally true for `"false"`, because both are non-empty strings. **And a gem nobody has published answers in plain text**: the entire 404 body is `This rubygem could not be found.` -- no braces, no key, no JSON, from a path ending in `.json`, so a client calling `.json()` throws a parse error rather than learning anything about the gem. Also pinned: `authors` is prose and `licenses` is a list, so nokogiri's four authors arrive as one comma-joined string and splitting on ", " is both the only way back and wrong for any name containing a comma; a dependency requirement is a sentence, `{"name": "actioncable", "requirements": "= 8.1.3.1"}`, so resolving anything means parsing it; the gemspec's URIs are promoted to the top level and left in place, agreeing three times out of four, with `homepage_uri` at the top and absent from `metadata` so neither copy is authoritative; two checksums for two different things (`sha` for the .gem, `spec_sha` for the serialised gemspec); and `platform` is the literal word `ruby` when there is no platform. Detection found another new near-miss kind, and it owns the biggest number on the page: **matched on a semantic** -- `@snyk/ruby-semver` at 42k is "a node-semver compatible API with RubyGems semantics", a JavaScript implementation of how this registry compares versions, which knows more about it than any client here and never sends it a request |
| ~~Repology~~ | Shipped, and written entirely from live responses: no OpenAPI, no credential, every case dated. **Asking for a project nobody packages is a 200 with an empty array** -- every other registry in this collection answers 404, so `if (!response.ok) throw` never fires and a typo is indistinguishable from a project that exists and is packaged nowhere. It is defensible (Repology indexes how others package things, and "nobody packages this" is a real answer) and it is the answer a client is least likely to have written code for. **And only five of the fourteen fields are on every entry**: one project's response is a bare array with an entry per packaging -- 806 for `curl` -- where `repo`, `visiblename`, `version`, `origversion` and `status` appear 806 times, `srcname` 794, `vulnerable` 632, `licenses` 568, `summary` 496, `categories` 462, `subrepo` 404, `binname` 347, `maintainers` 339 and `binnames` 158. Entries run from six keys to thirteen in eight distinct shapes; `srcname`, the field anybody reaches for first, is missing from twelve; `vulnerable` is absent from 174, which is not the same as false; and `binname` and `binnames` are the same idea singular and plural with **no entry carrying both**, so reading either alone misses half the repositories. Also pinned: six statuses at once for one project (377 outdated, 255 legacy, 163 newest, 7 rolling, 3 devel, 1 `incorrect`), the index's own judgement travelling in the data as an ordinary entry, and `origversion` -- the more specific of the two version fields -- being the one that can be null. A missing `User-Agent` is a 403, the second registry here to refuse for attribution rather than permission after crates.io. Nothing on either registry calls it |
| ~~Homebrew~~ | Shipped, and written entirely from live responses: no OpenAPI, no credential, every case dated. **An object called `versions` holds two strings and a boolean**, and one of the three is a version: `{"stable": "8.21.0", "head": "HEAD", "bottle": true}`. `head` is a git ref spelled the same way on every formula that has one, and `bottle` is a boolean saying whether a prebuilt binary exists -- so `Object.values(f.versions)` is `["8.21.0", "HEAD", true]` and rendering it as a version list shows a user the word `true`. **And four of its fields describe a computer that is not involved**: `installed` is `[]`, `pinned` and `outdated` are `false`, `linked_keg` is `null`, on every formula the API serves, because it is a static document on a CDN with no machine behind it -- while the same schema is what `brew info --json` prints locally, where those four are the whole point. Also pinned: a missing formula answers a **full HTML document** with `Content-Type: text/html` from a path ending in `.json`, which is the third way this collection has seen a `.json` path refuse (PyPI a JSON object, RubyGems a bare sentence, Homebrew a web page); Ruby symbols arriving as strings with their colon intact, `":provided_by_macos"` and `":any"`, so comparing against the bare word fails; install counts keyed by *command line*, where `analytics.install["30d"]` is `{"curl": 45041, "curl --HEAD": 20}`; and bottle files keyed by platform where two of the "macOS release" names are `arm64_linux` and `x86_64_linux`. Nothing calls this API from npm or Packagist, and the homonym is beer |
| ~~Packagist~~ | Shipped, and written entirely from live responses: no OpenAPI, no credential, every case dated. **The versions array is a chain of diffs, and reading one entry gives you the wrong answer.** `GET /p2/monolog/monolog.json` returns 87 entries: the first has 21 keys, the second 8, the third 7 -- because every entry after the first carries **only the fields that differ from the entry before it**. So `packages["monolog/monolog"][1].license` is undefined, not because 3.9.0 has no licence but because it has the same one as 3.10.0. **And the chain runs newest first**, so the deltas apply backwards in version order and a reader who assumes chronological order replays the whole thing the wrong way round -- getting an answer for every version, all of them wrong, none of them throwing. The only signal is a sibling key at the top of the document, `"minified": "composer/2.0"`; nothing inside an entry says it is a patch. **A field that goes away is the string `"__unset"`** -- six of monolog's entries carry one -- so a merge finds seven characters where it expected an object. Also pinned: a missing package is a **bare JSON string served as `text/html`** (`"404 not found, no packages here"`, with the quotes, so it parses as JSON and is not an object) -- a fifth distinct not-found shape in this collection; two timestamps 17 seconds apart for one release (`time` is the tag, `published-time` is when Packagist saw it); a `version_normalized` with four components against a three-component `version`; and a `dist.shasum` that is the empty string on every entry, so the field that would verify a download is present, typed correctly, and carries nothing. Detection's careful exclusion is `private-packagist/api-client`: the same company's commercial product at packagist.com with an API of its own, sharing a vendor prefix as well as a brand |
| ~~Hex~~ | Shipped, and written entirely from live responses: no OpenAPI, no credential for reads, every case dated. **The API hands you the dependency line to paste, and the three it hands you disagree about what to depend on.** Every package carries a `configs` object with a ready-made declaration per build tool; for `phoenix` 1.8.13 they are `mix.exs: {:phoenix, "~> 1.8"}`, `rebar.config: {phoenix, "1.8.13"}` and `erlang.mk: dep_phoenix = hex 1.8.13`. The first is a range that accepts any 1.8.x and will pick up 1.8.14 on the next resolve; the other two are exact pins. So one object tells an Elixir project to float within a minor line and an Erlang project to freeze on a patch, and the only thing deciding which policy you get is which key you copied -- neither wrong for its tool, and the choice made for you by a field most people read as three spellings of one answer. Also pinned: **four download counters and one has no window** (`all`, `day` and `week` say what they measure; `recent` does not, and there is no `month` to infer from); two fields for the latest version agreeing today and parting company the moment a pre-release lands, which is crates.io's four a second time and marks it as a category rather than one provider's quirk; a `meta.maintainers` kept from an older shape and always empty, while `owners` holds the actual usernames; an advisory with **three identifiers** for one finding (Hex's own id plus the CVE and the GHSA as aliases); a release entry that is a pointer rather than a release; and `retirements` as an empty *object* beside `security_advisories` as an array, so the two things a client checks before installing are shaped differently from each other. The 404 is clean JSON repeating its own status -- the sixth distinct not-found shape in this collection |
| ~~NuGet~~ | Shipped, and written entirely from live responses: no credential, every case dated. **The same package at two base URLs, and one of them is missing seventy-seven of its versions.** The service index advertises the registration resource four times under four versions of its *type*, pointing at three base URLs; resolve the unversioned `RegistrationsBaseUrl` -- the obvious thing to do -- and you get `registration5-semver1`, which for `Microsoft.Extensions.Logging` holds 98 versions against `semver2`'s 175. The missing ones are the SemVer 2.0.0 versions, `10.0.0-preview.1.25080.5` through `11.0.0-preview.7.26381.103` -- every preview of the current major and the one before it -- because SemVer 1 cannot express a dotted prerelease identifier, so those packages are omitted entirely. No error, no flag, no field saying anything was left out: a tool listing available versions is quietly two majors out of date. **And whether the index arrives whole or in pieces is decided by how many versions the package has**, not by which endpoint you asked: `newtonsoft.json` (84) is inlined on both base URLs with fragment `@id`s pointing back into the document you already hold; `serilog` (600) is split on both, with real URLs to fetch; and `Microsoft.Extensions.Logging` is inlined at 98 and split at 175, so the same package on the same day is both shapes. Nothing announces which you got — the only way to tell is to look for `items` and find it missing. Also pinned: `@type` as an array on the index (`["catalog:CatalogRoot", "PackageRegistration", "catalog:Permalink"]`) and a plain string on every object beneath it; a leaf carrying three URLs to three different things, one of them on a fourth host prefix; and a 404 that is **Azure Blob Storage's XML** -- `<Error><Code>BlobNotFound</Code>...` with `Content-Type: application/xml` -- so asking a package registry for a missing package reports that a blob does not exist, and the storage layer names itself while doing it. That is the seventh distinct not-found shape in this collection |
| ~~Go module proxy~~ | Shipped, and written entirely from live responses: no credential, every case dated. **The proxy refuses the module path in the URL and then hands it back in the body.** A Go module path is case-sensitive and the proxy protocol will not carry capitals -- each one is escaped to `!` plus its lowercase -- so the module written in every `go.mod` as `github.com/BurntSushi/toml` is fetched as `github.com/!burnt!sushi/toml`, and asking for the author's spelling gives `404  bad request: invalid escaped module path "github.com/BurntSushi/toml"`. Two things there: the status is 404 and the sentence says *bad request*, which are different claims about whose fault it is; and the escaped form, once accepted, answers with `Origin.URL: "https://github.com/BurntSushi/toml"` -- the capitals restored, in the body, by the server that would not take them in the path. Also pinned: **the JSON keys are Go struct field names**, `Version`, `Time`, `Origin` and inside it `VCS`, `URL`, `Ref`, `Hash`, PascalCase with an all-caps acronym because the type is marshalled straight out of `cmd/go` with no tags -- the only registry here whose wire format is somebody's internal struct; `Ref` being optional, present on `gorilla/mux` and absent from `BurntSushi/toml`; and a module nobody has published answering with **the proxy's own shell failure** -- the `git ls-remote` it ran, a path inside its cache on its own disk, `exit status 128`, git's "terminal prompts disabled" and two lines of advice. That is the eighth distinct not-found shape in this collection and the only one naming a directory on the server. Stated and not served: `/@v/list`, which answers plain text with one version per line and **in no order at all**, so anything wanting the newest must parse and compare every line. Detection found the sharpest matched-on-a-semantic yet -- `golang.org/x/mod` implements `module.EscapePath`, the exact rule above, and never opens a connection |
| ~~Maven Central~~ | Shipped, and written entirely from live responses: no credential, every case dated. **The search API is Apache Solr's response, handed over unedited.** There is no envelope of Maven Central's own -- what comes back is what the search engine produced, including a header describing the engine's own work. A successful search reports `"status": 0`, which is Solr's convention and not an HTTP status, sitting in a body that arrived with a 200; beside it are `QTime` and `params`, the request read back including `fl: id,g,a,v,p,ec,timestamp,tags` and `sort: score desc,timestamp desc,g asc,a asc,v desc` -- the internal field list and sort expression the service asks Solr for, published on the public API in every response. **And the fields are single letters**: `g`, `a`, `v`, `p`, `ec` for groupId, artifactId, version, packaging and the extensions that exist, with nothing in the response saying which is which -- and `p` is `bundle` rather than `jar`, because packaging is what the POM declared. Also pinned: `start` twice in one document in two types, `0` where the service computed it and `""` where the caller's value is echoed as text; `timestamp` as thirteen-digit epoch milliseconds, so a client reading seconds lands 57 000 years out; a search matching nothing answering **200 with status 0 and numFound 0**, which makes a missing artifact and a mistyped query the same response; and a query Solr will not parse coming back as plain text that stops mid-sentence -- `Solr returned 400, msg:` and nothing after the colon, an error relayed by a service that did not read it. Stated and not served: the Solr query language is not interpreted, `params` echoes whatever query was sent, and `repo1.maven.org` is a separate host serving jars and XML. Detection found the word meaning three things on one registry -- `bariew/maven` is an Israeli invoicing service and `sukohi/maven` manages an FAQ -- and matched-on-a-semantic for the fourth time, in `mvn-artifact-name-parser` |
| ~~OSV.dev~~ | Shipped, and written entirely from live responses: no credential, every case dated. **One vulnerability is three records, and they disagree.** The SQL injection fixed in Django 2.2.28, 3.2.13 and 4.0.4 is in the database three times, from three sources, each naming the others in `aliases` and none of them a redirect: `GHSA-2gwj-7jmv-h26r` carries a summary, a severity and three `affected` entries; `PYSEC-2022-190` carries no summary, no severity and no `database_specific`; `CVE-2022-28346` carries neither and **no `package` at all** -- no name, no ecosystem, no purl -- with a `GIT` range whose events are commit hashes and the version numbers only inside `database_specific.extracted_events`, a namespace the schema says is not portable. So `vuln.severity[0].score` works on one record and throws on the other two, and the CVE id -- the name a scanner is most likely to hold -- is the emptiest of the three. They are not on one schema version either: `1.9.0` against `1.7.3`, three publication dates spanning two days, and nine fractional digits on one `modified` against six on the others. **And they disagree about which version fixed it**: GitHub's record splits the three branches into three entries with two events each ascending, PyPA's packs all three into one entry with six events descending, so `affected[0].ranges[0].events[1].fixed` is `2.2.28` on one and `4.0.4` on the other. Also pinned: a query counts each advisory once per source, so django 3.2.0 answers 63 vulns of which 26 are pairs and 37 are distinct, with the halves forty entries apart because the array is sorted by id; `versions` sorted as text, `2.2.10` before `2.2.2`, so the last entry is `2.2.9` and the highest affected release is `2.2.27`; `severity` as CVSS vector strings with no number in the field while `database_specific.severity` is the word `CRITICAL`; `querybatch` answering with `{id, modified}` stubs under the same key `vulns` that carries whole advisories on `/v1/query`, and with `{}` for a query that matched nothing; and **`code` being a gRPC status on one failure and an HTTP status on another** -- `5` beside a 404, `3` beside a 400. Stated and not served: the third error shape, `{"message":"...","code":400}` with the keys reversed and a trailing newline, which is only reachable by sending a body that is not JSON. Detection found a new near-miss kind -- the vendor ships its database as a file, so `@renovatebot/osv-offline` contains no `osv.dev` URL at all -- and matched-on-a-semantic for the fifth time in `github.com/ossf/osv-schema`, which defines the format and imports `net/http` nowhere |
| ~~deps.dev~~ | Shipped, and written entirely from live responses: no credential, every case dated. **Capitalising a package name gets you a different package, and a 200.** On npm the name is case-sensitive, so `/v3/systems/npm/packages/express` answers 288 versions with 5.2.1 as the default and `/v3/systems/npm/packages/Express` answers three, from 2016, all deprecated. `Express` is a real package and its deprecation text -- "Package unsupported. Please use the express package (all lowercase) instead." -- is the only thing in the whole response that says anything is wrong: `packageKey.name` echoes `Express` back as though it were right, `isDefault` is set on one of the three, `advisoryKeys` is empty so a vulnerability scan of the wrong package comes back clean, and nothing else about the two responses differs in shape. **And the ecosystems disagree about the fix**: on PyPI the same request normalises instead, so `/v3/systems/pypi/packages/Django` answers with `packageKey.name: "django"`, the name you asked for silently replaced. Also pinned: the system echoed in upper case, `NPM` and `PYPI`, where the path took lower; `deprecatedReason: ""` present and empty on a version that is not deprecated, the exact opposite of the npm registry Recipe here, whose `deprecated` is absent rather than false; the dependency graph as **nodes and edges with integer indices** rather than a tree, so filtering the nodes array makes every edge point somewhere else; those nodes in **alphabetical order** rather than dependency order, so for `accepts` 1.3.8 the array is accepts, mime-db, mime-types, negotiator and `nodes[1]` is INDIRECT while `nodes[2]` and `nodes[3]` are DIRECT; the package itself as node 0 with relation `SELF`, one longer than the dependency count; `error: ""` as the way a successful graph says nothing went wrong; and **failures as bare plain text** -- seventeen bytes reading `package not found`, `text/plain`, no trailing newline -- which is also what an ecosystem that does not exist gets, naming the wrong one of the two things in the URL. Stated and not served: the version route's name scope, because `express` and `Express` both have a 3.0.1, four years apart, and this format's identifier space is one per resource -- the same limit recorded on the Go module proxy, and this time the headline is the reason for it. Detection found matched-on-a-semantic for the sixth time and the first time **inside the mapped module**: `deps.dev` holds both `deps.dev/api/v3` and `deps.dev/util/semver`, which imports `net/http` nowhere |
| ~~Open-Meteo~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The timestamps carry no offset, and the same string means two different hours.** With `timezone=America/Toronto` the hourly times read `2026-08-28T00:00` and the first is 17.5 °C; without the parameter they read `2026-08-28T00:00` and the first is 21.7 °C, because the default zone is GMT. There is no `Z`, no `+00:00` and no seconds, so `new Date()` on either parses as local time in the reader's own zone -- a third answer -- and being right means reading `utc_offset_seconds` from elsewhere in the document and applying it by hand. **And the coordinates come back moved**: a request for 43.65, -79.38 answers 43.646603, -79.38272, the nearest cell of the model, so code that round-trips what it sent is querying a different point next time. Also pinned: the hourly data as **parallel arrays** whose only relationship is the index, so sorting or filtering one destroys the association with no error at all; the units in a separate object where `hourly_units.time` is `"iso8601"`, a serialisation format called a unit for a field that has none; a second comma-separated coordinate changing the top-level type **from an object to an array**, with `location_id` then absent on the first element rather than 0, so anything keying by it loses the first location; `timezone_abbreviation` as `"GMT-4"` rather than `EDT`; and failures naming a field `error` that is only ever `true` -- a success carries no such key -- beside a `reason` that is **a Swift decoder's stringified error**, carrying the service's internal generic type signature onto the public API, while a forgotten `latitude` is reported as two lists disagreeing in length. Stated and not served: the forecast is served rather than computed, the arrays are trimmed to four of twenty-four hours, and `generationtime_ms` changes on every request. Detection found **the sharpest near miss in the collection**: the npm package `openmeteo` is Open-Meteo's own client and it sets `format=flatbuffers` unconditionally, so the vendor's own SDK never receives the JSON the vendor documents. `@openmeteo/sdk` is the generated schema for it, and `@openmeteo/weather-map-layer` is a third wire format again. Neither is mapped |
| ~~USGS earthquake catalogue~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Asking for a page removes the count.** The `metadata` block changes shape with the request: without a limit it carries `count`, with one it carries `limit` and `offset` and no `count` at all. A client that pages and reads `metadata.count` to know when to stop finds `undefined` on exactly the requests where it is paging, and only there. **And `offset` is one-based** -- a request that sent no offset comes back saying `offset: 1` -- so treating it as a zero-based cursor and adding the page size skips a record on every page after the first. Also pinned: `coordinates` as **`[longitude, latitude, depth]`**, GeoJSON's order and the opposite of the way it is said aloud, with the third element depth in kilometres rather than altitude, so reading it as a latitude-longitude pair puts the epicentre in the wrong hemisphere and nothing errors; `ids`, `sources` and `types` as comma-delimited strings that **begin and end with a comma**, so `",us7000tcf3,usauto7000tcf3,attkg5cs,"` splits into five elements for three identifiers; the key `type` appearing **four times in one document with four vocabularies** -- `FeatureCollection`, `Feature`, `Point`, `earthquake`; `mag` as a JSON number, so 5.0 arrives as `5` and prints as "5" while the `title` the service built beside it says `M 5.0`; `tz`, `felt`, `cdi`, `mmi` and `alert` present and null rather than absent, and carrying values on the event that has them; `metadata.status: 200`, the HTTP status restated inside the body; and **failures in plain text to a request that said `format=geojson`** -- a multi-line human-readable report, `text/plain;charset=UTF-8` with no space after the semicolon, reading the request back with its ampersand as `&amp;`, an HTML entity in a plain text body. Stated and not served: the catalogue is served rather than queried, the paged route is selected by `limit=2` rather than by the presence of a limit, and only `format=geojson` is modelled. Detection found one host serving three URL families where only one is this API: `dat-usgs-earthquakes` reads the pre-baked summary feeds, which take no query parameters and therefore **always** carry `metadata.count`, so the behaviour above cannot happen there at all |
| ~~Frankfurter~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Ask for a Sunday and you get Friday, with a 200 and no word about it.** `GET /v1/2026-08-23` answers `{"date": "2026-08-21"}`: the ECB publishes on working days, so a weekend has no fixing and the service falls back to the most recent one. The only thing in the response that says the date moved is the `date` field, which a client that already knows what it asked for has no reason to read -- so a two-day-old rate arrives labelled as the day you wanted, and over Easter that is four days. **And the base is never among the rates**: `base=USD` answers with EUR, GBP and JPY and no USD, so `rates[base]` is undefined and a loop converting each currency skips the identity; naming it explicitly does not help, because `symbols=USD,EUR` answers with EUR alone. Also pinned: a range changing the shape under the same key names -- `date` becomes `start_date`/`end_date` and `rates` nests one level deeper, so `rates.EUR` works on one and is undefined on the other; the range you get being the range that **has data**, so asking to 2026-12-31 answers `end_date` 2026-08-27, silently clamped; the days between simply absent, three keys for five days, so iterating a range and indexing `rates[date]` finds nothing two days in seven and more around a holiday; `amount` as `1.0`, a float multiplier that is 1.0 in almost every response and easy to read as a rate; and **two 404 shapes for the same words** -- a currency that does not exist and a date before the series both answer `{"message":"not found"}`, while a path that does not exist answers `{"status":404,"message":"not found"}`, so branching on `body.status` finds a number for the mistake you can see in your own source and `undefined` for the one you cannot. Stated and not served: the rates are served rather than computed; `/v1/currencies` answers a bare object keyed by currency code, which is data in the key positions rather than a collection of records; and `amount` is `1.0` on the wire and `1` here, because a JSON number that happens to be integral loses its decimal point through a Go float64 and this format cannot say "a number that keeps its point". Detection found two hosts, two version prefixes and packages written for a version that is not there: `api.frankfurter.app` 301s everything to `api.frankfurter.dev` with `/v1` added, while `@pontx/frankfurter-v2` and `go-finance` target a `/v2` where every path answers 404 |
| ~~Nominatim~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Nothing found is an empty array on one path and an error object on the other, and both are 200.** `GET /search` with no matches answers `[]`; `GET /reverse?lat=0&lon=0` with no match answers `{"error": "Unable to geocode"}` -- same service, same `format=jsonv2`, same status, and the reverse one puts an error key where an object of results should be, so `response.lat` is undefined and `response.ok` is true. **And search answers with an array while reverse answers with a bare object**, not an array of one, so the two endpoints of one geocoder cannot share a parser and the difference is not in the format parameter. Also pinned: a request without a `User-Agent` refused in **plain text** -- 403 and one line pointing at the usage policy -- from an endpoint that asked for JSON; `lat` and `lon` as **strings**, so arithmetic concatenates and comparison sorts lexically; `boundingbox` as four strings in the order **south, north, west, east**, which is neither GeoJSON's `[west, south, east, north]` nor a pair of coordinates in sequence; the keys of `address` depending on the place, with one carrying an administrative level **inside its own name** (`ISO3166-2-lvl4`, where the 4 is Canada's subdivision depth and is 3 or 6 elsewhere), so the key to read is not knowable in advance; `importance` in scientific notation, `8.130235751883753e-05`, which anything printing into a template writes as `8.130235751883753e-5`; and `place_id` being local to the installation, so the stable identity is `osm_type` plus `osm_id` -- two fields. One more that is not a shape but explains a lot of wrong-looking results: reverse geocoding the CN Tower's own coordinates answers with the restaurant inside it, because the nearest node beats the containing way. Stated and not served: the gazetteer is served rather than searched, the usage policy's rate limit is not a response shape, and only `format=jsonv2` is modelled. Detection found three kinds of neighbour that are not clients: `@mailwoman/nominatim` is a Nominatim-compatible **server**, `geo-golang/mapquest/nominatim` is the same software hosted by somebody else behind an API key -- so neither the 403 nor the no-credential reading applies there -- and `tile_proxy` is the same project's other service, raster tiles rather than JSON |
| ~~Open Library~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The answer is keyed by the string you sent, and a miss is that key's absence.** `GET /api/books?bibkeys=ISBN:0451526538` answers `{"ISBN:0451526538": {...}}`, so reading it means rebuilding the query string -- prefix, colon and all -- to index your own response; and a lookup that matched nothing answers `{}` with a 200, not a null under the key and not a 404, so the difference between found and not found is whether a key you constructed yourself is there. **And the same field name carries a reference on one endpoint and a resolved object on the other**: the canonical document answers `"authors": [{"key": "/authors/OL18319A"}]` and `/api/books`, for the same book, answers `[{"url": ..., "name": "Mark Twain"}]` -- one a pointer that costs another request, and nothing but the shape to say which you hold. Also pinned: `/isbn/{isbn}.json` answering a **302 when it finds the book** and a 404 when it does not, so the redirect is the happy path and a client that does not follow redirects gets an empty body on success; that 404 being a **full HTML page** served as `text/html` from a path ending in `.json`; timestamps as typed objects, `{"type": "/type/datetime", "value": "2008-04-01T03:28:50.625462"}`, with microseconds and no timezone, so `record.created` is an object; `type` as a reference as well; `revision` and `latest_revision` as two fields carrying one number; every identifier as an array including `"openlibrary": ["OL1017798M"]`, which cannot repeat; and **two schemes in one document** -- `url` and `authors[].url` are `http://` while `subjects[].url` and the covers are `https://`, so the record links are the insecure ones. Stated and not served: the catalogue is served rather than searched, the `bibkeys` key is declared as one of its possible values, and the `/isbn` redirect and its 404 are served for one ISBN each rather than as a rule. Detection found the ordinary-word collision at its limit: searching Packagist for "open library" returns PHPWord, openspout, bootstrap-icons, geophp, casbin and nelexa/zip -- eight libraries that are open source and none about this provider |
| ~~Wikipedia REST~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A page that does not exist is reported as an internal error.** `GET /api/rest_v1/page/summary/Zzzznotarealpagexyz` answers `{"status": 404, "type": "Internal error"}` and that is the whole body -- two fields, nothing naming the page, and a `type` that says the same thing for a title nobody created as for a real failure inside the service, so a client logging it logs the wrong cause and a client branching on it cannot tell a 404 it should expect from a 500 it should alert on. **And the title arrives five times in four fields, two of them HTML**: `title`, `displaytitle`, and `titles` with `canonical`, `normalized` and `display` inside it, where `displaytitle` and `titles.display` are the same markup twice. Also pinned: a **disambiguation page answering 200** with an `extract` that is a flattened list -- `"Mercury most commonly refers to:Mercury (planet), the closest planet to the Sun\nMercury (element), a chemical element"`, no space after the colon and newlines between the entries, so anything rendering `extract` as prose renders exactly that; the spelling you asked for being **nowhere in the answer**, so asking for `toronto` gives `Toronto` in every field that could have recorded the correction; `originalimage.source` ending `3840px-...jpg` while `originalimage.width` says 6632, so laying out from the declared size and loading the URL gets a picture forty per cent smaller; those image URLs carrying `utm_source`, `utm_campaign` and `utm_content` **analytics parameters**; `content_urls.desktop` and `content_urls.mobile` sharing four keys of which three are byte-identical, so a client that picks one branch is right three times in four by accident; `revision` as the string `"1370998502"` where `pageid` is the number 64646; and the main namespace's `text` being the empty string rather than "Main" or null. Stated and not served: title normalisation and redirect following are real behaviours asserted from recorded responses rather than applied, the extracts are trimmed except the disambiguation one, and only `/page/summary` is modelled. Detection found **two APIs on one host and almost everything using the other one**: MediaWiki's Action API at `/w/api.php` predates this by a decade, so node-wikifetch, symfony/ai-wikipedia-tool and easy-wiki all reach wikipedia.org and none of them reach anything served here |
| ~~Hacker News~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A thing that is not there is the four bytes `null`, with a 200.** Not an empty object, not a 404, not an error field: the JSON literal, sent as `application/json`, from the item endpoint and the user endpoint alike. `response.ok` is true, `.json()` succeeds, and the value it hands back is the one that means "no value" in the language you are writing in -- so `(await res.json()).title` throws a TypeError on a 200 with no status anywhere to have branched on first. **And two endpoints have no envelope at all**: `/v0/maxitem.json` answers a bare integer, `49475165`, and that is the whole document; `/v0/topstories.json` answers a bare array of five hundred integers. Not objects -- identifiers, so rendering a front page is one request for the list and five hundred for the items, and there is no batch endpoint anywhere in the API. Also pinned: `id` as a **number on an item and the username on a user**; `text` as **HTML inside a JSON string** with no plain-text companion, so anything rendering it strips tags and anything not rendering it shows them; `time` as epoch **seconds**, ten digits where most of this collection is thirteen, so a client guessing the unit from the field name puts a 2011 comment in 1970; one resource serving five kinds of thing where a story has `title`, `url`, `score` and `descendants` and a comment has `parent` and `text` and none of those; and `kids` being the direct replies while `descendants` is the total, so a story with `descendants: 3` has four entries in `kids` and neither field is the length of the other. Stated and not served: `topstories` carries three of its five hundred identifiers, because the length is the point and cannot be held in a fixture; the bare values are served as declared bytes with a JSON content type, because this format describes records and neither a number nor an array of numbers is one. Detection found the generic name belonging to somebody else's API of the same data: the npm package called `hacker-news-api` is a client for `hn.algolia.com` and holds no `firebaseio.com` at all |
| ~~PokeAPI~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The description is still formatted for a Game Boy.** `flavor_text` is `"A strange seed was\nplanted on its\nback at birth.\fThe plant sprouts\nand grows with\nthis POKéMON."` and there are three things in that one string: the newlines are hard wraps at the original 1996 screen's column width, so re-flowing them breaks mid-phrase; the `\f` is a **form feed**, U+000C, the page break from the handheld's text box, in a JSON string today; and the name carries a lowercase é between two capitals, because the font in those games had no capital É. None of it is escaped, flagged or duplicated in a cleaned-up field, so a client that upper-cases the string produces something the source never said. **And there are twenty-eight English descriptions, not one**: `flavor_text_entries` is a flat array of a hundred and two, one per game version per language, so "the English description" is a filter returning twenty-eight answers and picking the first is picking Red from 1996. Also pinned: a miss as the **bare plain text `Not Found`** with a 404, from an API whose every success is JSON; a listing row as a `name` and a `url` with **no identifier at all**, so the id exists only as a path segment inside that URL and anything more than the name costs one request per row; `previous` and `next` both always present and each null at its own end, so `"next" in body` is true on every page and says nothing; the lookup being case-insensitive, so `/pokemon/1`, `/pokemon/bulbasaur` and `/pokemon/BULBASAUR` are one record and the response never says which was asked for; a type being `{slot, type: {name, url}}`, so the name is two levels down and the array index is not the slot; and four of the ten top-level `sprites` keys present and null. Stated and not served: the entries hold two of a hundred and two, the list four of 1351, and `previous` is a constant null -- right on the first page and wrong on any other, because this format has a forward cursor and no companion going back. Detection found the offline-data kind at its largest: `pokeapi-sprites` is 4237 files, every sprite shipped as an asset, and it lands doubly because the sprite URLs this API returns point at raw.githubusercontent.com and were never on pokeapi.co either |
| ~~MusicBrainz~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **It answers XML unless you ask for JSON, and it says so in XML.** `/ws/2/artist/{mbid}` defaults to `<?xml version="1.0" ...>`, so a client that forgets `fmt=json` gets a document its `.json()` throws on; and a client that mistypes the parameter gets **406 in `application/xml`**, carrying the message that the recognised types are application/json and application/xml. The one answer that would tell you how to ask for JSON is the one you cannot parse without already having solved the problem. **And an identifier that does not exist is a 400 about the identifier's format**: `00000000-0000-0000-0000-000000000000` is a perfectly well-formed UUID that nothing uses, and it answers `{"error": "Invalid mbid."}` -- exactly what a string that is not a UUID gets. There is no 404 in this shape, so "no such artist" and "you sent nonsense" are one answer and it is wrong about the first; the two bodies also arrive with their keys in opposite orders. Also pinned: a request without a `User-Agent` answering **403 Forbidden with a message about being throttled**, where the status says you are not allowed, the text says you are going too fast, and the cause is neither -- and that body arriving as `application/json` with **no charset** where every other JSON answer here declares one; the field names being **hyphenated** (`sort-name`, `type-id`, `begin-area`, `iso-3166-1-codes`), so `artist.sort-name` is a subtraction in JavaScript and every one must be bracket-accessed; `life-span` mixing precisions in one object, `begin: "1987"` beside `end: "1994-04-05"`, both strings and neither saying which; and **four spellings of "not applicable" in one document** -- `gender` null, `end-area` null, `ipis` `[]`, and inside `area` a `disambiguation` of `""` beside a `type` of null. Stated and not served: the database is served rather than searched, the one-request-a-second rate limit is not a response shape, `inc=` subtrees are not modelled, and `fmt=yaml` stands for every format that is neither json nor xml. Detection found the headline as a near miss: npm carries two clients for one API, one per wire format, and only the JSON one is mapped |
| ~~Open Food Facts~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A product that is not there is a 200 with `status: 0`.** The HTTP status says it worked and the truth is a number in the body, so `response.ok` is true and there is **no `product` key at all** -- `body.product.product_name` throws a TypeError on the answer a barcode scan is most likely to get. And `status: 0` has more than one cause: `product not found` and `no code or invalid code` share the number, so branching on it cannot tell a barcode that does not exist from one that should not have been sent. **And the code you sent is not always the code that comes back**: asking for `0000000000000`, thirteen zeros, answers `"code": "00000000"` -- eight -- and nothing anywhere records what was asked for. Also pinned: **`energy` is kilojoules**, with `energy: 2252` and `energy-kcal: 539` in one object, so reading the unqualified name and labelling it calories is out by 4.18; every nutriment repeated under four keys (`fat`, `fat_100g`, `fat_value`, `fat_unit`) and energy under twelve; key names mixing separators inside one key, `added-sugars_100g`; `brands` as a **comma-joined string** rather than a list, crowd-sourced, where one of the three on the reference product is not a brand; `quantity` as the empty string on a jar that says 400 g; and `nutriments_estimated` as a parallel object of the same nutrients, computed rather than declared and keyed identically, with nothing inside either saying which is which. Stated and not served: the database is served rather than searched, a product without `?fields=` comes back with several hundred keys and the fixture holds the subset the cases turn on, and only the `world.` subdomain and v2 are modelled. Detection found the counter-example to Open-Meteo: there the vendor's own client spoke a format the vendor's API does not use, and here the vendor publishes an official client in all three ecosystems and every one calls this API. Six MCP servers carry the name, the most for any Recipe here |
| ~~Crossref~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A date is an array inside an array, and it is not always the same length.** `"issued": {"date-parts": [[2013, 7, 31]]}` on one work and `[[1970, 6]]` on another: the outer array exists because the field can hold a range, and the inner one stops wherever the metadata stopped. Nothing anywhere says which precision arrived, so `date-parts[0][2]` is a number on one work and `undefined` on the next, and building a Date from three positional arguments silently reads the month as the day. **And one work has three publication dates at two precisions**: `published` and `published-online` are `[[2013, 7, 31]]` while `published-print` is `[[2013, 8]]`, so "when was this published" has three answers, one of them a month, and the print date is after the online one. Also pinned: `title` as an **array of one** -- and so `container-title`, `short-container-title` and `ISSN`; `created` carrying the same instant three ways, `date-parts`, `date-time` and a `timestamp` in epoch milliseconds, where `published` beside it carries only the first, so two date fields in one record disagree about how much of a date they will give you; the envelope being `status`, `message-type`, `message-version` and `message`, where the version is of the envelope rather than the API, has read `1.0.0` for years, and `status` is the string `"ok"` beside an HTTP 200; a single-work lookup carrying a relevance `score` of 1 for a query it did not make; and a DOI that does not resolve answering the **bare plain text** `Resource not found.` with a 404 and no charset, from an API whose every success is a four-field JSON envelope. Also served: the two halves of a DOI are not independent, so a real suffix under another registrant gets the same prose as a DOI nobody has minted. Stated and not served: the index is served rather than searched, a work's forty-odd top-level keys are trimmed to the subset the cases turn on, and the polite pool is a routing decision rather than a response shape. Detection found "cross-reference" as a technical term in **four unrelated fields** -- content nodes in Neos, objects in a vector database, pandoc-style references in Markdown, and the bibliographic one |
| ~~Zippopotam.us~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The keys have spaces in them.** `"country abbreviation"`, `"post code"`, `"place name"`, `"state abbreviation"` -- not camelCase, not snake_case, not kebab-case: a space, in the key. Nothing here can be destructured, nothing can be reached with a dot, and every access is a bracket with a quoted string in it. It is legal JSON and almost nothing else does it. **And `"CA"` is two different things in one API**: the country abbreviation for Canada, and the state abbreviation for California, in fields whose names differ by one word. Also pinned: `latitude` and `longitude` as **strings**, where the minus sign puts the lexical order nowhere near the numeric one; a postcode that does not exist answering `{}` with a 404, so `.json()` succeeds and every field is undefined and the status is the only signal -- and a country code that does not exist, and a real postcode under the wrong country, answering exactly the same way; `places` holding **two different record types under one key**, with `state` and `state abbreviation` from the postcode lookup and `post code` from the place lookup, and nothing in either saying which arrived; and the reverse lookup repeating `"place name"` at the top level and inside every entry, on the request that named it in the first place. One more that explains a lot of wrong-looking Canadian addresses: a forward sortation area covers several neighbourhoods, so M5V's place name is a hundred and thirty characters of slash-separated district names inside parentheses -- one string, one field, and not a name anybody would print. Stated and not served: the gazetteer is served rather than looked up, and the reverse lookup holds three of the five Beverly Hills postcodes. Detection found the funniest near miss in the collection: `haoteam/hippo-buyer-php-sdk` is a "hippopotamus buyer php sdk" and matches on `-potam-`, the Greek root for river that both names happen to carry |
| ~~TVmaze~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The failure puts the reason in `name` and leaves `message` empty.** A missing show answers `{"name": "Not Found", "message": "", "code": 0, "status": 404}` -- four fields, and the one a client would read is blank. `name` carries the reason phrase, and on every successful response `name` is the show's title, so the same key holds "Under the Dome" on a 200 and "Not Found" on a 404; `code` is 0, which is neither an HTTP status nor an error code; and `status` is the HTTP status restated inside the body. **And "where does it air" is two mutually exclusive fields**: every show carries `network` and `webChannel` and exactly one is an object, so `show.network.name` throws on half the catalogue and the check has to be written both ways round -- and inside them the shapes disagree, the network carrying a country object where the web channel carries `country: null`. Also pinned: `runtime` and `averageRuntime` as two fields where one is often null and nothing says which to prefer, Loki having `runtime: null` beside `averageRuntime: 50` while Under the Dome has 60 in both; `schedule.time` as **the empty string** on a show that has a broadcast day, not null and not absent; `externals` holding three identifiers in **three types** -- `null`, a number, and the string `"tt9140554"` -- so parsing them the same way is impossible; `summary` as HTML in a JSON string with no plain-text companion; and `updated` as epoch seconds, ten digits, beside a `rating` object holding one number and a `weight` integer with no unit. Stated and not served: the catalogue is served rather than searched, and the episode, season, cast and crew endpoints are not modelled. Detection found **the vendor's other API behind the vendor's own paywall**: `@datafire/tvmaze` integrates the TVmaze *user* API, which lives under a key and covers watchlists, while `tvmaze-api-ts` calls itself "a tvmaze scraper and api wrapper" and is both -- half a client of this API and half a client of the website |
| ~~Sunrise-Sunset~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The sunset is before the sunrise.** A request for Toronto answers `"sunrise": "10:35:37 AM"` and `"sunset": "12:01:46 AM"` -- both UTC, both a twelve-hour wall clock with no date and no zone on them. The sunset that evening is a minute past midnight UTC, which is the next calendar day, so parsed against the date that was asked for it lands sixteen hours before the sunrise; anything comparing the two, or drawing a bar between them, gets a negative interval and nothing says the day rolled over. **And the times are UTC while looking like local time**: `"10:35:37 AM"` reads as half past ten in the morning, Toronto's sunrise was 06:35, and the only thing that says otherwise is `tzid: "UTC"` -- a sibling of `results` rather than a field inside it, so code holding the results object alone cannot know. Also pinned: a latitude of **999** answering `status: "OK"` with a full set of times, because there is no validation of the coordinate at all; a failure setting `results` to the **empty string** where a success has an object, so `typeof results` is "string" on one and "object" on the other; `day_length` as `"13:26:09"`, a duration in the same colon-separated shape as the times beside it and the only field in the object that is not a clock reading; and **the polar day reported as a day of zero length** -- at 78 north on the solstice every time is `"12:00:01 AM"` and `day_length` is `"00:00:00"`, with `status: "OK"`, so the one latitude where the question is hard is answered exactly backwards. Stated and not served: the ephemeris is served rather than computed, and `formatted=0` switches every field to ISO 8601 with an offset, which is a different document from the same path -- the formatted form is the default. Detection found **matched-on-a-semantic covering the entire neighbourhood** for the first time: the sun's position is a closed-form calculation, so suncalc, suncalc-php, go-sunrise and a dozen more compute the answer instead of asking for it, both Go packages named exactly after the provider compute rather than fetch, and there is no Go or PHP client of this API anywhere |
| ~~NHTSA vPIC~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Four errors arrive in one string, and one of them is 400 on a 200.** `ErrorCode` is `"6,7,11,400"` -- a comma-separated list of numbers inside a string, four codes for one VIN, and the `400` looks like an HTTP status and is not, because the response is a 200 and always is. `ErrorText` is those four joined with `"; "` and one of them contains its own semicolon, so splitting on the separator gives five pieces for four errors. **And the failure is a success with 148 empty strings in it**: `Count` is 1, the `Results` array holds one object, and 148 of its 154 fields are `""` -- nothing absent and nothing null, so a VIN the service could not read looks exactly like one it read and found nothing about. Also pinned: **every value a string**, with `ModelYear` as `"2003"`, `EngineCylinders` as `"6"` and `DisplacementL` as `"2.998832712"`, a three-litre engine to nine decimal places as text; `Message` as a 250-character **disclaimer** on every successful response, the same paragraph every time, explaining that a missing value does not mean the feature is absent; `SearchCriteria` as the request in prose, `"VIN(s): 1HGCM82633A004352"`, with a label and parentheses; the **error fields populated on a success**, where `ErrorCode` is the string `"0"` and `ErrorText` restates it as `0 - VIN decoded clean...`; and `SuggestedVIN` on a bad VIN being `"N!TAV!N"` -- the input with its invalid characters replaced by exclamation marks, which is not a VIN and cannot be submitted anywhere, while the actual correction sits in `AdditionalErrorText` as `Invalid character(s): 2:O, 6:I.` And **the one failure status on the whole API** is for a path that does not exist, whose message quotes the URI as the backend saw it -- `backend-vpic-api.nhtsa.dot.gov`, a host the caller never named -- under a `message` key spelled in lower case where every success spells it `Message`. Stated and not served: a result carries 154 keys and the fixture holds fourteen, and without `?format=json` the same path answers XML -- the same shape recorded for MusicBrainz. Detection found matched-on-a-semantic and the offline-data kind at once: a VIN is defined by ISO 3779, so `sunrise/vin` at 152,000 installs decodes it locally and `@cardog/corgi` advertises "fast, offline VIN decoding" and holds no nhtsa.dot.gov at all |
| ~~Wikidata REST~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The English label is not there, and the name is under a code that is not a language.** Douglas Adams's item carries 75 labels and `en` is not one of them -- nor `de`, `fr`, `es`, `it`, `nl`, `pt`, `pl` or `sv` -- because every Latin-script language was consolidated into a single entry under `mul`, ISO 639-3 for "multiple languages" and not a code anybody's browser will send. Asking the label endpoint directly confirms the absence: `/labels/en` is a 404 and `/labels/mul` is `"Douglas Adams"`. So `item.labels[userLang] ?? item.labels.en`, which is what every client writes, shows nothing at all -- while `item.descriptions.en` beside it works, because the descriptions were not consolidated. Two sibling objects, keyed the same way, disagreeing about which languages exist. **And every key in `statements` is an opaque number**: 312 of them on this item, and learning that `P569` means "date of birth" is another request per property. Also pinned: a value as `{"type": "value", "content": ...}` where `content` is a **string for an item reference and an object for a date**, decided by a `data_type` on the property beside it; a date as three fields one of which is a URL -- `time` with a **leading plus** on the year, an integer `precision` nothing decodes, and a `calendarmodel` of `http://www.wikidata.org/entity/Q1985727` that must be resolved to learn the calendar is Gregorian; a statement identifier of `Q42$F078E5B3-...`, the item, a **dollar sign** and a UUID, which has to be URL-encoded to be sent anywhere -- and whose item prefix is not consistently cased, 285 beginning `Q42$` and 27 beginning `q42$`; and a missing item and a missing label being **the same body** with only `context.resource_type` different. Stated and not served: the knowledge base is served rather than queried, SPARQL is a different service on a different host, and an item's 75 labels, 121 descriptions, 312 statements and 132 sitelinks are trimmed to a handful each. Detection found a kind not seen before -- the npm package called `wikidata` is a **placeholder**, two files whose description says so, and nothing has arrived -- alongside Packagist's whole first page being the vendor's own PHP value-object model, `diff/diff` at 1.7 million installs included, which opens no connection at all |
| ~~GBIF species match~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A name it could not find comes back at a hundred per cent confidence.** `?name=Zzzznotaspeciesxyz` answers `{"confidence": 100, "matchType": "NONE", "synonym": false}` with a 200 -- three fields, and `confidence` is the highest the scale goes on the one answer that found nothing, so code ranking results by confidence puts the failures first. Every real match scores lower: exact is 99, fuzzy 95, a genus standing in for a species 94. **And `synonym` is present only when nothing was found**: the no-match carries `synonym: false`, and the response whose `status` actually is `SYNONYM` does not carry the field at all -- the one key named for the concept appears exactly where the concept does not apply. Also pinned: a **deprecated name matching EXACTly**, since `Felis concolor` is what cougars were called until 1993 and asking for it answers `matchType: "EXACT"` at 98 with `status: "SYNONYM"`, exact because the string exactly matched a name nobody should use; that response carrying **two keys for two different taxa**, where `usageKey` is the synonym and `acceptedUsageKey` and `speciesKey` are the accepted species, so the field with the plainest name is the one not to store; a misspelled species **silently becoming a genus**, with `Puma notaspecies` answering `matchType: "HIGHERRANK"` at 94, `rank: "GENUS"` and `scientificName: "Puma Jardine, 1834"`, which a client will print as though it were the species asked for; and a typo corrected without any record of the input, so nothing in the response can be shown as "did you mean". And the one failure status here is a **leaked Java exception**: `/v1/species/notarealendpoint` answers 400 with the bare text `For input string: "notarealendpoint"`, because `/v1/species/{key}` takes an integer key and the segment was handed to a number parser whose message became the body. Stated and not served: a request with no `name` parameter at all, which answers the same 200 with `confidence: 100` plus a `note: "No name given"` -- so a forgotten parameter and a name nobody uses are the same answer, one of them with prose attached; and the backbone taxonomy is served rather than matched, and the occurrence and dataset endpoints the classification keys point into are not modelled. Detection found a new way to collide -- not a shared word but a shared **edit distance**: Packagist's three highest-ranked results for "gbif" are a PHP GIF codec at 34.9 million installs, a payments library at 17.2 million and an Apple sign-in parser at 2.2 million, fifty-four million installs and not one about biodiversity |
| ~~Nager.Date~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`fixed` is false on Christmas Day.** The twenty-fifth of December is as fixed as a date gets, and the field that says whether a holiday falls on the same date every year says no. It says no for New Year's Day too, and for every other holiday: across Canada, the United States, the United Kingdom and Germany -- eighty-one holidays in 2026 -- `fixed` is false eighty-one times, and `launchYear` is null on all eighty-one. Two booleans-and-a-number with one value each, so branching on either is branching on a constant. **And one date carries six holidays, none of them national**: the third of August 2026 in Canada is the Civic Holiday, British Columbia Day, Heritage Day, New Brunswick Day, Natal Day and Saskatchewan Day -- six entries, one date, thirteen provinces and territories between them, and every one `global: false`, so `holidays.find(h => h.date === today)` returns whichever the array ordered first and calling it "the holiday" is wrong in five provinces out of six. Also pinned: the field naming a subdivision being called **`counties`** while holding Canadian provinces, German states and US states, and being **`null` rather than empty** when the holiday is national, so the type is array-or-null and `.length` throws on the national ones; twenty-three of Canada's thirty-one holidays not being `global`; `localName` and `name` as separate fields that are identical on all thirty-one, differing only where a country has more than one official language; `types` as an array whose only value in a national list is `Public`; and **two failures sharing an envelope and disagreeing about the rest of it** -- an unknown country carries `detail`, an unsupported year carries `errors`, and the `title` on the second is a framework's default validation sentence with the real message two levels down inside an array under a key named after the parameter. Stated and not served: the calendar is served rather than computed, and only `PublicHolidays` is modelled. Detection found the Sunrise-Sunset shape again -- a holiday calendar is a rule set, so `date-holidays` and its country-specific cousins compute rather than ask -- plus two packages written for versions that are not there: `@pontx/nager-date` is a "v4 SDK" and `/api/v4/` answers 404 |
| ~~Wayback Machine~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Ask whether a page existed in 1990 and it answers yes, from 2002.** `closest` means closest available, `available` is `true`, `status` is `"200"`, and nothing in the response says the snapshot is twelve years from the date that was asked for -- the distance is computable only by subtracting two timestamps in two different formats, and no field states it. **And `timestamp` is in the document twice, meaning two different things**: the one inside `closest` is when the snapshot was taken, fourteen digits, and the one at the top level is the date the caller asked for, eight digits, echoed back -- and absent entirely when the caller sent none, so the same key is a request parameter in one place and an answer in the other, at two lengths, one of them conditional. `url` is doubled the same way. Also pinned: `status` as the **string** `"200"`, an HTTP status quoted inside a body that already carried one, so `=== 200` is false and `== 200` is true; a URL nothing has archived answering **`archived_snapshots: {}`** with a 200, an empty object where the hits have `closest`, so reading `.closest.url` throws rather than finding a status; a request with **no `url` parameter answering the plain text `Error: no url parameter` -- with a 200**, in `text/html`, from an endpoint whose successes are JSON; and the archive address being `http` with the URL inside it varying, the same site appearing as `http://example.com/`, `https://example.com/` and `http://example.com:80/` across three snapshots. Stated and not served: the archive is served rather than searched, and only `/wayback/available` is modelled. Detection turned up the sharpest near miss in the collection -- **`internetarchive/virustotalapi`**, published by the Internet Archive and wrapping VirusTotal, so a Composer vendor prefix that is exactly right maps exactly the wrong API -- plus Esri's World Imagery Wayback under the same name on `wayback.maptiles.arcgis.com`, the bare npm name `wayback` meaning revision-tracking for arbitrary data, and four other endpoints on the archive itself: CDX (unanimously, every Go package that turned up), Save Page Now, item metadata, and an undocumented internal |
| ~~National Weather Service~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The path says lat,lon and the body says lon,lat.** `GET /points/39.7456,-97.0892` answers `"coordinates": [-97.0892, 39.7456]` -- the same pair, in the same exchange, reversed, because the path takes latitude first the way a person writes a coordinate and the body puts longitude first the way GeoJSON requires, and nothing in the response says so. **And `type` is in that document three times, in three vocabularies**: `Feature` at the top (GeoJSON's), `properties.@type: wx:Point` (the NWS ontology's), and `properties.type: land` (the NWS's own) -- only the middle one namespaced, plus a fourth on a Feature nested inside a property. Also pinned: **a number with a unit expressed four ways in one forecast** -- a `unitCode` object for elevation, a bare number beside `temperatureUnit` for temperature, the prose `"10 to 15 mph"` (and sometimes the scalar `"15 mph"`) for wind, and a `unitCode` object for a percentage; **`validTimes` as an interval whose right half is a duration**, `2026-08-28T16:00:00+00:00/P7DT9H`, so `Date.parse` gives `NaN`; the envelope being UTC while every period inside it is local; `geometry.coordinates` being a number at depth 1 on `/points` and an array at depth 3 on `/gridpoints`, with the polygon's ring carrying five entries for four corners; **a missing `User-Agent` refused by the CDN rather than the API** -- 403 in HTML, so a client with `application/problem+json` handling parses markup and has no `correlationId` to quote; and **three 404s meaning three different things**, where the accurate one carries the least and a latitude of 999 says only "Not Found" with the reason hidden in a `parameterErrors` array RFC 9457 does not define, whose message for an unknown office runs 840 characters because it inlines all 133 valid office codes into a sentence. Stated and not served: two points rather than the continental grid, three of fourteen periods, and only `/points` and the forecast. Detection mapped clients in all three ecosystems and turned up a new near miss: **`latlon-utils` carries `noaajs`'s description word for word** -- "obtaining data from weather.gov and tidesandcurrents.noaa.gov" -- while its code contains no weather.gov at all, so prose would map it and names do not. Around it: npm's bare `nws` is a static web server, `vrtnws` is Dutch for news, NOAA's tides, buoys, space weather and GFS files are other services, `alerts.weather.gov` is XML and `forecast.weather.gov` is HTML somebody scrapes, `@hebcal/noaa` computes rather than asks, and five near-identical `weathermodule` packages are one course's coursework pointed at Dark Sky |
| ~~Open Trivia DB~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The success key is `results` and the rate limit's key is `result`.** One letter, on one of four answers: `body.results.length` works against a hit, against no results and against an invalid parameter, and finds `undefined` on the fourth -- which is the one that arrives when the caller is going too fast, so the crash appears under load and not before it. **And three of the four are HTTP 200**: no results is 200, an invalid parameter is 200, and `response.ok` is true for every outcome except the rate limit, leaving a number in the body that has no name anywhere in the response as the only way to tell them apart. Also pinned: **an unknown category reported as an empty one** -- `?category=999` answers `response_code` 1, for no results, not 2, for a bad parameter, so a typo in a category id looks exactly like a category that happens to be empty; **HTML entities inside JSON**, `&quot;` and `&#039;` in a format that needs neither, so every string wants an HTML decode after it has already been JSON-decoded; `incorrect_answers` running three entries on a multiple-choice question and one on a boolean with nothing but the sibling `type` to say which; the category filter taking a number while the field holds a name, with nothing in the response linking `9` to `General Knowledge`; and **`encode=base64` encoding the enum values too**, so `type === "multiple"` is false until even the fields a caller would never think to decode have been decoded. The Recipe also found a runtime defect, fixed alongside it: a route's list-envelope field override was written into the Recipe-wide map it was copied from, so one route's `response_code` leaked into every other route for the life of the process. Stated and not served: questions are answered in fixture order rather than drawn at random, the base64 answer is recorded rather than computed, and the API has no pagination to declare -- `amount` caps at fifty and there is no offset or cursor, which is why its five listings sit in the count above. Detection mapped npm and Go and found nothing on Packagist at all: searching it for "opentdb" returns **OpenTBS** five times, a document-templating library one letter away, plus Open Tibia and a package describing itself as "a trivial implementation of timeouts". The sharpest miss is named after the API and does not call it -- `github.com/tamnd/opentdb-cli` builds `opentdb.com/browse.php`, the human website |
| ~~Open Brewery DB~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`brewery_type` says `closed`.** The field naming what kind of brewery this is also answers whether it is one -- fourteen values covering `micro`, `brewpub` and `regional`, which are kinds, plus `closed`, which is a status, `planning`, which is a stage, and `location`, which is neither. 642 are `closed` and 639 `planning`, 1,281 of 11,848, so `by_type=micro` silently drops every closed micro. **And two pairs of fields are the same field twice**: `state` and `state_province` carry the same string on every record, as do `street` and `address_1` -- identical, not similar -- with nothing saying which pair is current and which is kept for compatibility. Also pinned: **`?per_page=201` answering 302 to the homepage** rather than 400 or a clamp, so a client following redirects gets the landing page and a 200; a missing brewery answering **404 as an HTML page** from a JSON API, so a client parsing every response as JSON throws on the parse; the documented **`autocomplete` path being a 301 to `search`**; **`by_state` holding 214 keys most of which are not states** -- ACT is Australian, Argyll Scottish, Auckland a New Zealand region, Aveiro a Portuguese district, all sorted in among Alabama and Arizona; `meta` carrying `page` and `per_page` on a document that is not a page of anything; `phone` running `"+49 9261 628000"` and `"4058160490"` in one column; and the country filter taking `united_states` while the field answers `"United States"`, so a client cannot round-trip its own filter. Stated and not served: three breweries rather than 11,848, `by_state` trimmed from 214 keys to eight, and only the list, the single brewery, search, meta and the two redirects modelled. Detection's near miss is one word wide: **BreweryDB** is a different, commercial API on `api.brewerydb.com` with more clients than this one -- `brewerydb-node` against `openbrewerydb-node`, same suffix, two hosts -- and "brewery" on npm is also a scaffolding tool with no beer in it, while `@umbraculum/brewery-api-client` runs a brewery rather than listing them |
| ~~PoetryDB~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`status` is the number `404` and the string `"405"`.** The same key, in the same API, holding a JSON number on one failure and a JSON string on the other -- so `body.status === 404` catches one and misses the other, and `body.status >= 400` is true for one and false for the other because `"405"` is not a number. **And neither is an HTTP failure**: both arrive 200 with `application/json`, so `response.ok` is true and every retry and error boundary sees a success. Also pinned: **`linecount` as the string `"14"`**, a number that is a string beside a string that is a number; **`poemcount` being a valid input field and an invalid output field**, with one refusal's list of allowed words containing it and the other's not, both 405; a path segment that is not a field being reported as **Method Not Allowed** on a GET that is the only method the endpoint has; and **`/author` answering an object where `/author/{name}` answers a bare array**, two paths one segment apart that cannot share a parser. The Recipe also made a runtime declaration real: `envelope.array` wraps one record in a list and was read only on the wrapped path, so a Recipe asking for a bare array of one got the object back and nothing said otherwise -- Xero and Ghost both wrap, so no shipped Recipe had exposed it. PoetryDB answers `/title/Ozymandias` with `[{...}]`, no key and no envelope, which is now served and now refused in combination with the Recipe-wide success fields, since those are object properties with nowhere to go in a list. Stated and not served: two poems rather than the collection, the roll of authors trimmed from 129 names to six, and output-field projection declared per route rather than modelled combinatorially. Detection's neighbourhood is mostly the Python packaging tool -- `snyk-poetry-lockfile-parser` reads `poetry.lock` against `python-poetry.org` -- plus the **European Commission's Poetry Service**, which is how translation is requested inside the Commission, Chinese divination verses, and a three-hundred-thousand-poem corpus shipped as files |
| ~~Datamuse~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A definition is a part of speech, a tab character and a sentence, in one string.** `"defs": ["n\tA large amount. "]` packs two fields into a JSON string with a control character between them, in a format that has arrays and objects for exactly this -- and every entry ends in a space, so splitting on the tab and rendering index 1 shows trailing whitespace nobody asked for. **And `tags` is three vocabularies in one array of strings**: `["n", "v", "pron:S L UW1 "]` puts bare part-of-speech codes beside a colon-prefixed key-value pair, so getting the pronunciation means scanning the array for a prefix and splitting on a colon, and that value carries a trailing space too. Also pinned: **three queries on one endpoint answering three different field sets**, where `md=` explains only one of the differences -- `?ml=ocean` gets tags and no syllable count, `?rel_rhy=blue` gets a syllable count nobody asked for and no tags, and `md=dpsr` on the rhyme query gets all five; **a bad parameter name, a bad value and a bad sign answering 200, 400 and 500**: `?nosuchparam=x` is an empty array with a 200, so a typo in a name looks exactly like a query that matched nothing; `?max=abc` is a clean 400 naming the parameter; and `?max=-1` is a 500 with an incident id, because the validation that caught `abc` checks the type and not the sign; `score` as an integer with no scale anywhere, 40,041,792 for a synonym and 10,051 for a spelling suggestion; and `/sug` answering two fields where `/words` answers five, so two endpoints on one host do not share a type. The Recipe also fixed the router: a route selected by query parameters scored a flat bonus however many it matched, so two matching routes tied and the request went to whichever was declared first -- silently, with nothing in either route saying it depended on the order. Scoring is per parameter now, so the more specific route wins. Stated and not served: results are answered from a fixture rather than ranked, "slew" is trimmed from thirteen definitions to two, and the unknown-parameter case is served from an empty store, so it shows the shape of that answer rather than proving the parameter caused it. Detection's near miss is a registry habit: **Packagist answers "datamuse" with Monolog and Doctrine**, falling back to popularity when nothing matches. Its one real candidate is left unmapped on purpose -- the repository is gone, so the only evidence is a registry description, which is exactly what `latlon-utils` had. Around it: Wordnik needs a key and, because Swagger was created there, searching for it returns the whole Swagger toolchain still pointing at `api.wordnik.com` |
| ~~Art Institute of Chicago~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **One field of ninety-eight is licensed differently from the rest, and the response says so in prose.** `info.license_text` names `description` in backticks as CC-BY while everything else is CC0, and that English sentence in a sibling key is the only statement of which is which -- nothing structural marks it, so a client republishing the response under one licence is wrong about exactly one column. **And an image URL is three pieces from two levels**: the artwork carries `image_id` and no URL, the base lives in a top-level `config` beside `data` rather than in it, and `/full/843,/0/default.jpg` is in neither. Also pinned: **`date_display` and `date_end` disagreeing** -- Seurat's Grande Jatte reads "1884-86, border added 1888-89" and sorts on 1886, so anything ordering by date puts it earlier than the display says; **the same coordinate three times with the third rounded**, `latitude` and `longitude` at thirteen decimal places and `latlon` joining them as a string at twelve; **`colorfulness` scoring 0 on a pointillist painting** made of coloured dots against 43.0038 for a Franz Marc, so nothing distinguishes "not computed" from "computed, and grey"; **`has_not_been_viewed_much`** as a negated boolean, making the ordinary question a double negative; `alt_classification_ids` being `classification_ids` without its first entry, which the name does not say; and a base64 GIF riding inline in every row of every listing. Stated and not served: two artworks rather than 132,681, each trimmed from 98 fields to the 25 the cases read, and four of the six paging keys left out -- including the pair worth recording, since the envelope reports an `offset` while the `next_url` it hands you asks for a `page`. Detection's near miss is a substring: **`artic` is the start of "article"**, so Packagist answers it with `php-readability`, `graby` and Facebook's Instant Articles SDK, and npm's bare `artic` describes itself as "Artic. It's for articles." It is also a React component library and a CSS grid. The two Go modules are one author's tool under two names, and the same author's `opentdb-cli` -- unmapped in the Open Trivia DB row above -- calls that provider's website instead of its API |
| ~~Postcodes.io~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Four of the fields are named after organisations that no longer exist**: `nhs_ha` (Health Authorities, abolished 2002), `primary_care_trust` (PCTs, 2013), `nuts` (EU regions the UK left in 2020) and `ccg` (Clinical Commissioning Groups, 2022) -- the last sitting next to `icb`, the body that replaced it, holding almost the same string with nothing marking either as historical. **And in Wales the two abolished English fields hold the name of a Welsh body**: both `nhs_ha` and `primary_care_trust` answer "Cardiff and Vale University Health Board", while `icb`, an English structure, says "Wales". Also pinned: **the unversioned `lsoa` being an alias only one postcode can resolve** -- in Westminster `lsoa`, `lsoa11` and `lsoa21` agree and say nothing, and in Cardiff they disagree, which is what shows `lsoa` tracking the 2021 census; **a null field's code being a row of nines**, `admin_county` null against `codes.admin_county` `"E99999999"`, with Cardiff carrying that same sentinel under `codes.icb` while `icb` itself holds "Wales"; **`ccg_id` in three shapes across two records**, `"W2U3Z"` and `"7A4"` against the letter-and-eight-digits everywhere else; `ruc11` and `ruc21` differing in vocabulary rather than vintage, one carrying a parenthetical country prefix and the other not; `date_of_introduction` as `"198001"`, six characters with no separator; a Senedd constituency that is numbered where every other constituency is named; and two coordinate systems on one record, OSGB36 integers beside WGS84 floats, with neither datum named. Two 404s carry two sentences -- "Postcode not found" and "Invalid postcode" -- so the API draws a distinction its status code does not. Stated and not served: two postcodes rather than 1.7 million, each trimmed from 46 fields to 28 and `codes` from 29 entries to 8. Detection's near miss is the vendor's own commercial product: **Ideal Postcodes runs postcodes.io and also sells a paid UK lookup** that needs a key. Around it, "postcode" mostly means South Korea -- the Daum/Kakao address widget outnumbers the UK packages in a search for the bare word -- and three more carry no URL at all, because a postcode is a string with a grammar |
| ~~Metropolitan Museum of Art~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`artistGender` is `"Female"` or the empty string.** There is no `"Male"` -- a man is recorded as the absence of a value, and the absence of a value is what twenty other fields on the same record use to mean "this does not apply". Anything counting women in the collection can do it; anything counting men cannot tell them from the unattributed. **And absence is spelled `""` rather than `null` on twenty-one of fifty-seven fields**, so `if (x)`, `x === null` and `"x" in obj` disagree with each other. Also pinned: **three year fields carrying two types** -- `accessionYear` the string `"1993"`, `artistBeginDate` the string `"1853"`, `objectBeginDate` the number `1889`; **the same `dimensions` field using U+00D7 on one record and the letter `x` on another**, so anything splitting on the character has to know both; a missing image as `""` where a URL goes, which assigned to an `img` src asks the page for its own address; a painting whose `classification` is empty while another painting says `"Paintings"`; and `GalleryNumber` as the only capitalised field name on a record that also shouts `artistWikidata_URL` and `AAT_URL` in the middle -- fifty-seven fields, three naming conventions. Stated and not served: two objects rather than 490,000, each trimmed from 57 fields to 30, and only the object lookup modelled. Detection found nothing at all on Packagist and mostly MCP servers on npm, three of them the same server republished under three scopes. The one worth naming calls itself **"The official TypeScript library for the Met Museum API"** and is a demo package under a personal scope from an SDK-generation company: it does call the API, so it maps, but the word "official" in a description is worth nothing on its own |
| ~~Rick and Morty~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Asking for a character by a name instead of a number is a 500.** `GET /api/character/abc` answers `{"error": "Hey! you must provide an id"}` with a server-fault status, so every 5xx alarm in the client fires on a caller mistake that retrying cannot fix -- while the 404 beside it, for a number that is simply not a character, is the correct status. **And `"unknown"` is the only lowercase value in the enum**: `status` is `"Alive"`, `"Dead"` or `"unknown"`, so the one value meaning "we have no answer" is the one that breaks a switch written from the other two. Also pinned: **one object saying "no value" two ways at once** -- `origin.name` as the string `"unknown"` beside `origin.url` as the empty string, in the same two-key object; `type` as prose or `""`, so a character with no sub-species gets a string rather than a null; **related records as URLs rather than identifiers**, an `episode` array of `.../api/episode/10` and an `origin.url` with no id beside it, so joining means parsing an integer back out of a string; and three failures carrying three sentences across two statuses, only one of which ends in a full stop. Stated and not served: two characters rather than 826, and Rick's 51 episode URLs trimmed to three. Detection's near miss is a fandom rather than a homonym -- the show has more packages about it than clients of its API. Two fetch GIFs from Giphy, `rick-morty-components-lib` is a React component library themed after it that **names `rickandmortyapi.com` in its README while its code calls React's CDN**, and `squanch` is a text mangler named after a catchphrase |
| ~~Deck of Cards~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A failed draw empties the deck and hands you the cards anyway.** Asking a fifty-two card deck for 999 cards draws all fifty-two, sets `remaining` to 0, says `"success": false`, returns them, and answers HTTP 200 -- so a client that checks `success` and discards the body has thrown away fifty-two real cards and left the deck empty, while one that ignores it has fifty-two it did not ask for the right number of. Either way the deck is spent, and the request that spent it is the one the API called a failure. **And the ten is a zero**: codes are two characters, so the ten of spades is `"0S"` while its `value` says `"10"`, and anything parsing a code by its first character reads a nought. Also pinned: `value` as a string for numbers and a word for faces, so sorting a hand puts `"10"` before `"6"` and both before `"ACE"`; **every image in the response twice**, `image` and `images.png` character for character; **a non-numeric count answering Django's default error page in HTML** from a JSON API; and the failure that does mean it -- a deck id that does not exist -- carrying the same `success` and `error` fields with a 404 that agrees with them. Stated and not served: three cards rather than fifty-two, a fixed `deck_id` because the real one is minted at random, and `/api/deck/new/` left out entirely, since an emulator answering it would have to invent an identifier that could never be asked about again. Detection's near miss is that **a deck of cards is a data structure**: nearly every package with the name is a local implementation, including the bare npm name `deckofcards`, which is "a standard deck of playing cards". And a new kind of miss sits beside them -- `api-false-success` is a dataset about APIs that return HTTP 200 on requests that cannot succeed, which is exactly what this Recipe's headline pins: a package about the finding, not about the provider |
| ~~TheMealDB~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`meals` is an array, or `null`, or a string.** A lookup that finds something answers `{"meals": [{...}]}`, one that finds nothing answers `{"meals": null}`, and one with no identifier answers `{"meals": "Invalid ID"}` -- one key, three JSON types, HTTP 200 on all three. `body.meals.length` works, throws, then returns 10, because `"Invalid ID"` is ten characters long, so a client checking `if (body.meals)` iterates the error message character by character. **And forty of the record's fifty-four fields are one array flattened into columns**: `strIngredient1..20` and `strMeasure1..20`, so a nine-ingredient recipe carries eleven empty slots -- **and they are not empty the same way**. Slots 10 to 15 are `""` and slots 16 to 20 are `null`, in one numbered series, split at an index nothing explains. Also pinned: `strTags` as a comma-joined string with no spaces where an array belongs; `strInstructions` carrying **CRLF inside a JSON string**, so splitting on `\n` leaves a carriage return on every line; `strArea` as `"Japanese"` beside `strCountry`'s `"Japan"`, the adjective and the noun in two fields; and every field prefixed `str` except the four that are not, so the convention holds for fifty of fifty-four and a client cannot guess which. Stated and not served: one meal rather than the collection, its instructions trimmed to two sentences, and a bad API key recorded rather than served -- it answers an IIS 404 page in XHTML, and the path would have to carry the key for a Recipe to route on it. Detection's near miss is a substring again: **Packagist reads "themealdb" as "theme"**, answering with a themeable syntax highlighter and a themeable admin panel, because the first five letters of the provider are a word |
| ~~World Bank~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The success is a two-element array and the failure is a one-element array.** A country lookup answers `[{page, pages, per_page, total}, [{...}]]` and a bad code answers `[{message: [{...}]}]`, both HTTP 200 -- so `body[1]` is the collection on success and `undefined` on failure, and the length of the outer array is the only structural signal that anything went wrong. **And `per_page` is a string where its three neighbours are numbers**: four fields describing one page, three integers and one quoted. Also pinned: **`iso2Code` and `iso2code` both on the record**, capital-C at the top and lowercase one level down; **three codes that are not ISO codes** -- `"XU"` for North America, `"XD"` for High income, `"XX"` for Not classified, all from the range the standard reserves for private use, in a field named after it; **a missing value as an object of three empty strings** rather than `null`, so `if (country.adminregion)` is true and `.value` is `""`; and coordinates quoted on a record whose page numbers are not. The Recipe also added a list style: no existing one produced a two-element array, so the shape had to be modelled as something it is not or the provider left out. `tuple` serves it, and it earns its own style because of the failure -- a Recipe that served an object would let a client write `body.data.map(...)` and ship it. Stated and not served: one country rather than 296, the refusal routed on the literal code that was checked because a listing whose scope matches nothing answers an empty collection rather than a differently-shaped document, and omitting `format=json` recorded rather than served, since it answers XML. Detection's near miss is the word "world": Packagist returns Swift codes for "Banks in the world", two PayTR gateways where **World is a Turkish credit card brand**, and `laravel-zip` -- twice -- described as "the world's leading zip utility", which matched on a marketing superlative |
| ~~disease.sh~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Six fields end in `PerOneMillion` and three of them are integers.** `cases`, `deaths` and `tests` are whole numbers; `active`, `recovered` and `critical` carry two decimal places -- one suffix, one obvious meaning, two JSON types, with nothing in the response distinguishing the groups, so typing the family as an integer truncates half of it. **And `countryInfo` carries `_id`**: a field named the way a database names its primary key, underscore and all, on a public API, holding the UN numeric country code rather than an internal row id. Also pinned: `updated` as a **JavaScript millisecond timestamp** in a bare integer, which changes on every request while the figures beside it do not; **coordinates in whole degrees** on a record whose derived statistics carry two decimals, so the least precise numbers describe a place -- and one of them called `long` rather than `longitude`; the reciprocal fields floored to integers, so `oneTestPerPeople` is 1 where the true figure is 0.58 and reads as one test each when it is nearly two; and **a country that does not exist answering JSON while a path that does not exist answers an Express HTML page**, both 404. Stated and not served: one country rather than 231, and figures that no longer move because the dataset stopped being updated. Detection's near miss is the clearest example of the data shipped as a file in this collection: **`simple-country-names` is one static JavaScript module holding disease.sh's country list**, with the vendor's hostname on every flag URL in all two hundred entries, and it never opens a connection. Grepping for the host maps it; reading what the host is doing there does not |
| ~~Dog CEO~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`message` is a URL, or an array of strings, or an error sentence.** The random-image endpoint answers a URL under it, the sub-breed list an array of names, and a missing breed a sentence -- and two of the three are strings, so nothing but reading `status` first tells the image from the error. Code that assigns `message` to an `img` src renders a broken image whose alt text is the error. **And the object's shape depends on whether it worked**: a success carries `message` and `status`, a failure carries `message`, `status` and `code`, so the status line is restated twice -- once as a word and once as a number -- and only on the way that failed. Also pinned: sub-breeds as bare names rather than paths or objects, so a client holding the list has to remember which breed it asked about. Stated and not served: one image and two breeds' sub-breeds rather than the collection, and the random endpoint answering the recorded image. Detection's near miss is a package that is only a type definition: **`@mgrzmil-org/api-types` is "Generated TypeScript types + OpenAPI spec for dog-api" and names dog.ceo nowhere** -- not in code, not in documentation -- so it describes the API's shapes without naming the API, and there is nothing to verify. Beside it, `dog-names` is a list shipped as data, and the rest of what "dog api" returns is the word "api": `fast-diff`, `wrap-ansi`, `slice-ansi` and `mock-axios` |
| ~~Advice Slip~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every response carries two `Cache-Control` headers, and they contradict**: one says `max-age=3600`, the other `max-age=600, private, must-revalidate`, both on the same 200. A recipient joining repeated field lines the way RFC 9110 says gets one directive list with two `max-age`s in it, which the grammar has no answer for -- and that is exactly what `fetch`'s `headers.get` returns. **And three outcomes carry three different top-level keys**: `{"slip": {...}}`, `{"slips": [...]}` and `{"message": {...}}` -- singular, plural, and neither, so there is no key a client can read first to find out what happened. Also pinned: `total_results` as the string `"1"` beside an `id` that is the number 1; **the same slip having two field sets**, since `/advice/1` answers `{id, advice}` and the search answers `{id, advice, date}` for it, so a date exists that one route never shows; and **a miss and an empty search both answering 200**, distinguished only by `"error"` against `"notice"` in a field called `type` inside an object called `message`. The Recipe also fixed the runtime: route headers were written from the create path and the listing path and nowhere else, so a get declaring them had them dropped in silence -- every Recipe using route headers until now declared them on a create, which is the shape SendGrid needed them for. Stated and not served: two slips rather than the collection, and the two `Cache-Control` lines served as one joined field, which is the closest a header map can come and preserves the finding either way. Detection's near miss is a networking standard: **`@serialport/parser-slip-encoder` implements SLIP**, the Serial Line Internet Protocol, which has been framing packets since 1988 |
| ~~SWAPI~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every number is a string, and one of them has a unit in it**: `diameter` is `"10465"`, `population` is `"200000"`, and `gravity` is `"1 standard"` -- so sorting by diameter puts `"10465"` before `"9000"`, and the field that looks most like a number is the one least safely read as one. **And the absence of data is spelled two ways in one record**: the planet at `/planets/28` is named `"unknown"`, and its six numeric fields split three-three -- `rotation_period`, `orbital_period` and `diameter` are `"0"` while `surface_water`, `population` and `gravity` are `"unknown"`, so a client treating `"0"` as a real zero charts a planet with no diameter and one treating `"unknown"` as the sentinel misses the other three. Also pinned: a planet's **name** being the string `"unknown"` rather than null or absent, in the field an interface renders as a title; related records as URLs rather than identifiers, with the empty case an empty array; and a planet that does not exist answering **the static site's own HTML shell with a 200**. Stated and not served: two planets rather than sixty, ten residents trimmed to two, and swapi.info rather than the original swapi.dev. Detection's near miss is the provider's own family tree: **one API, four hostnames, three response shapes**. swapi.co is the original and answers 301; swapi.dev and swapi.info answer the bare record; and **swapi.tech wraps everything three levels deep** with a Mongo identifier beside it and timestamps generated when the record is read -- so `react-swapi` calls a different API wearing the same name and the same field names, and is not mapped. Beside it, **Swapi Finance** is a decentralised exchange, swap plus API, with its own npm scope and a bee for a logo |
| ~~Chuck Norris~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The success and the failure carry timestamps in two different grammars**: a joke's `created_at` is `"2020-01-05 13:42:19.324003"` -- a space instead of a T, six digits of fraction, no timezone -- and an error's `timestamp` is proper ISO 8601, so which one arrives depends on whether the request worked and a client logging when something happened gets a date on its errors and `Invalid Date` on its results. **And the two identifier alphabets do not match**: most jokes carry base64url with hyphens and underscores, some carry lowercase letters and digits only, both twenty-two characters. Also pinned: `updated_at` never differing from `created_at`, so a field that exists to say when something changed has never said it; **a query parameter that is not a category answering 404 about a path that does exist and did answer**; and an empty search as `{"total": 0, "result": []}`, the array called `result`, singular, beside a total that is not. Stated and not served: two jokes rather than the collection, with the random endpoint answering the recorded one. Detection's near miss is that **there are two Chuck Norris APIs**: ICNDb, the older one, is gone, and Packagist has no client for this one at all -- its four results share one description word for word, and the two that were opened both call `icndb.com`. Four forks of one package, pointed at a host that no longer answers. Three more open no connection at all, including the bare npm name `chucknorris` ("Chuck Norris does not need a description") and `chuckscript`, a programming language whose source is written entirely in zeroes |
| ~~Open Notify~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **There is no HTTPS at all** -- port 443 does not answer, so `http://api.open-notify.org/iss-now.json` returns 200 and the same URL over TLS times out, and no page served over HTTPS can call this API from a browser. The only published client hardcodes `http://` as its default domain because there is no other one to reach. **And every failure declares JSON and sends none**: a 404 and a 405 both answer `Content-Type: application/json; charset=UTF-8` with zero bytes, so a client that trusts the header calls `.json()` on an empty string and throws. **The endpoint the documentation still describes fails differently again** -- `/iss-pass.json` is gone and answers nginx's own HTML 404, naming the version it runs, so there are two 404s in two content types on one API and which one arrives depends on how far in the request got. Also pinned: the position as two strings beside a number, `message` reading `"success"` on a 200, and `number` sent alongside the array it counts on an endpoint where `craft` is not always `"ISS"` because Tiangong is in there too. Stated and not served: the station moves, Cauldron pages a listing at ten while the real endpoint sends all twelve at once, and the nginx page carries Cauldron's error content type rather than `text/html`. Detection's near miss is **the exe wrapper at scale**: PeopleInSpace ships as five npm packages, one per CPU architecture, forty to ninety-five megabytes each, every one holding the same Java jar with the string `open-notify` nowhere in it. Packagist has nothing -- searching it for `iss` returns ten packages about *issues*, and its one `open-notify` hit is a WeChat callback SDK |
| ~~Agify~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Three APIs share one budget** -- api.agify.io, api.genderize.io and api.nationalize.io are documented separately and interleaved calls decrement a single counter, one allowance of twenty-five a day, so asking all three about a name costs three of them. **And `X-Rate-Limit-Reset` is seconds until midnight UTC, not a timestamp**: it read `60485` with `60501` seconds left in the UTC day, so a client reading it as a Unix epoch lands in January 1971. Also pinned: a name nobody has heard of answering 200 with `"age": null` rather than 404; present-but-empty and absent being different, where `?name=` is a 200 and omitting the parameter is a 422; and **the 422 carrying no rate-limit headers at all**, on a content type without the `charset=utf-8` every success declares -- the response that says a client is doing something wrong is the one that will not say how much budget is left, and it is free. Stated and not served: the empty-name 200, the sibling hosts, and the counter, which is a real quota against a real host. Detection's near miss is **the vendor's own npm scope**: `@agify/n8n-nodes-leadifyv2` and `v3` are about lead management and contain no API host at all. And **the edit distance** -- searching Packagist for `agify` returns ten packages for Apify, a different company one character away, one of which describes itself as "AI-generated and AI-maintained" |
| ~~Bible API~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The same verse comes back with different whitespace depending on the translation, and not the same kind of whitespace** -- the World English Bible pads `John 3:16` with a leading newline and two trailing ones, the King James with ten leading spaces and one trailing newline, one query parameter apart on the same path. **And the text is sent twice, both copies untrimmed**: the top-level `text` and `verses[0].text` are byte-identical, padding included. Also pinned: **two endpoints sharing no keys at all**, where a reference answers six top-level fields and `/data/web/random` answers two and none of them match; the same concept under two names, `book_name` and `book`; `"Public Domain"` arriving as `translation_note` in one place and `license` in the other; and **two 404s, one of them nginx's own HTML**, because a chapter past the end of a book never reaches the application while an unknown book does. Stated and not served: two translations of seventeen, and a recorded verse on the random endpoint. Detection's near miss is **one letter of top-level domain** -- the npm scope `@bible-api` ships a translation per package, kjv1769 and the Vulgate and a Polish Bible of 1632, and every one of them points at `bible-api.io`. "Bible API" names a category rather than this API: at least seven products answer to it, and three more npm clients carry no host at all, so nothing in them says which one they mean |
| ~~CoinGecko~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The map is keyed by currency and one of the keys is not a currency** -- `/simple/price` answers `{"bitcoin": {"usd": 78148, "eur": 67468, "last_updated_at": 1788075830}}`, so a client iterating the object to list prices gets a currency called `last_updated_at` worth 1,788,075,830. **And anything the API does not recognise is dropped in silence**: `?ids=bitcoin,zzzznotacoin` answers a 200 carrying only bitcoin, with no error and no null to say one was rejected, and asking only for unknown coins answers `{}`. Also pinned: one price an integer and the next a float in the same response, truncated by magnitude rather than typed by field; **two error shapes on one API**, a flat `{"error": "coin not found"}` and a nested `{"status": {"error_code": 429, ...}}`; and `/ping` answering `{"gecko_says": "(V3) To the Moon!"}`, where the only statement of which version replied is the `(V3)` mid-sentence. This is the Recipe that made a map with no wrapper real: `key: "-"` on the map style, the sibling of the bare array, with the same rule refusing Recipe-wide success fields beside it. Stated and not served: prices move, `ids=` is not modelled, and the two empty answers are recorded rather than served. Detection's near miss is **a browser engine** -- `ua-is-frozen` and `universal-user-agent` match on Gecko as in Mozilla's, in the user-agent string every browser still sends -- and **the vendor's own MCP server**, published under the same scope as its official library and still not mapped |
| ~~iTunes Search~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The JSON is served as JavaScript** -- every response carries `Content-Type: text/javascript; charset=utf-8`, success and failure alike, so a browser sniffs its way to the right parser while anything checking the header first either refuses a body it could have read or hands it to a JavaScript evaluator. **And there is no such thing as no results**: `?term=zzzzqqqxyznothing` answers sixty matches, every one an audiobook by "ZZZ Sleep", because the search decomposes the term and finds something for the pronounceable part. Also pinned: **a band is not the first result for its own name** -- `?term=radiohead` answers a feature film, and Radiohead appears only with `entity=musicArtist`; a heterogeneous array where a track and an artist share almost no keys and are told apart only by `wrapperType`; `amgArtistId`, a foreign key into All Media Guide, renamed in 2013; and an id nobody has answering the same empty envelope as a search with no term, so a rejected request and an empty one are one answer. Stated and not served: the real body is wrapped in three newlines each side, so byte-exact comparisons differ by six bytes. Detection's near miss is **an XML namespace** -- Apple's RSS extension for podcasts means feed libraries carry the word without making any request at all -- and Apple's other APIs, which is most of what Packagist returns for it: receipts, iTunes Connect sales figures, and Apple Music behind a developer token |
| ~~EPSS~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **One of the keys has a hyphen in it** -- the envelope is `{"status": "OK", "status-code": 200, ...}`, and `response.status-code` in JavaScript is `response.status` minus a variable called `code`, so the field most likely to be read by code checking whether the request worked evaluates to `NaN` with no error and no warning. **And the status is stated three times**: the HTTP status line, `status`, and `status-code`. Also pinned: the probability as a string to nine decimal places, `"0.999990000"`, wrong for every numeric comparison that does not parse it first; **a malformed CVE answering 200 while omitting the parameter answers 404**, so the wrong request succeeds and the incomplete one fails; and **failing making the data private**, where a 200 says `"access": "public"` and a 404 says `"access": "private"`. Stated and not served: scores move, two of 366,298 are served, and `limit` is clamped rather than refused. Detection's near miss is **the same data without the API** -- `npm-epss-audit` reads the bulk CSV from `epss.cyentia.com` and links `www.first.org` for the documentation. Packagist has no client at all, and what it returns for `epss` is a thermal receipt printer library, where the letters sit inside ESC/POS. And `cctx` exists only to occupy a typo of another package's name |
| ~~Chess.com~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The country is a URL** -- a player's `country` is `"https://api.chess.com/pub/country/US"`, so rendering it puts an API address on the page and getting the two letters costs a second request to a resource that is three fields, one of them that same URL. **And the player has two URLs that are not the same URL**: `"@id"` is the API resource and `"url"` is the web page, on two hosts, with two capitalisations of one username, and `@id` begins with a sigil so it needs a subscript in most languages. Also pinned: `twitch_url` sent twice, once at the top level and once inside `streaming_platforms`, byte-identical; `status` holding a subscription tier, `"premium"`, in a response whose HTTP status is 200; and **a failure code of zero**, shared by a player who does not exist and a path that does not exist, the latter answering `Data provider not found for key "/pub/player/hikaru/notanendpoint".` Stated and not served: follower counts move, and one player and one country are served. Detection's near miss is **the vendor namespace that is not even the vendor's repository** -- Packagist's `chesscom` namespace publishes Honeybadger error monitoring, from someone else's GitHub account -- alongside a shelf of chess libraries that never make a request, and the widest set of fuzzy matches yet for one word, Chester County's ArcGIS FeatureServer among them |
| ~~ip-api~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **TLS is a paid feature, and the message saying so is delivered over TLS** -- `https://ip-api.com/json/1.1.1.1` answers 403 with `"SSL unavailable for this endpoint, order a key"`, so the connection is accepted, the certificate is valid, the handshake completes, and the body says SSL is unavailable. **And every other failure is a 200**: a malformed address and a string that is not an address both answer HTTP 200 with `{"status": "fail", "message": "invalid query"}`, so `response.ok` is true for both and the only real signal is a word. Also pinned: **a format the API does not recognise being the marketing website**, eleven kilobytes of HTML on a 200 where a client called `JSON.parse`; a rate limit of `X-Rl` and `X-Ttl`, two abbreviations and neither standard; and the autonomous system as a number and a name in one string, comma included, beside an `isp` with no full stop and an `org` that is somebody else. Stated and not served: the TLS refusal is armed rather than reached by changing scheme, and `/csv/` and `/xml/` are recorded rather than served. Detection's near miss **corroborates the Recipe**: `react-currency-localizer` advertises "HTTPS-compatible IP geolocation" as a feature and calls `ipapi.co`, one hyphen away. And `ethernet-ip-cip` is an industrial automation protocol, where the IP is Industrial rather than Internet |
| ~~Where the ISS at~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The same field is a number on one endpoint and a string on another** -- `/v1/satellites/25544` sends `"latitude": 26.268487855719` and `/v1/coordinates/...` sends `"latitude": "37.795517"`, and `id` does it too, `25544` on the satellite and `"25544"` on its own element set. **And the units are declared after the numbers they govern**: `altitude` is 419.69504466844 kilometres by default and 260.21728290923 miles with `?units=miles`, and the only thing that says which is the last field in the object, so a parser reading in order has the numbers before it has the scale. Also pinned: **`daynum` as a Julian date beside a Unix timestamp**, two epochs in one object with no unit, no suffix and no name to tell them apart; the failure putting the status last, `{"error": "satellite not found", "status": 404}`; and **a path with no controller behind it saying exactly that** -- `"Invalid controller specified (nosuchthing)"`, internal vocabulary in a message for whoever typed the URL. Stated and not served: the station moves, and Cauldron's encoder does not escape forward slashes the way this one does. Nothing on either registry calls it, so it ships unmapped -- see the row below |
| ~~TfL~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The severity scale runs backwards** -- `statusSeverity` 10 is `"Good Service"` and 6 is `"Severe Delays"`, so the bigger number is the better railway and sorting descending to put the worst lines first puts the healthiest first instead. **And every object carries its .NET type and assembly**: `"$type": "Tfl.Api.Presentation.Entities.Line, Tfl.Api.Presentation.Entities"`, at every depth, beginning with a sigil so it needs a subscript in most languages. Also pinned: **a status carrying the year 1**, `"created": "0001-01-01T00:00:00"`, which is .NET's `DateTime.MinValue` on the wire because nothing set it, beside a `Line` whose own `created` is a real 2026 date; **a failure saying what went wrong four times**, with a .NET class name in the field a client switches on and a timestamp to seven fractional digits; and **a 404 naming an internal host**, `"Resource not found: http://api:8001/NoSuchEndpoint"`. Stated and not served: line status changes by the minute, two lines of eleven are served, and one line can carry two statuses at once. Detection's near miss is **one person's initials** -- Packagist answers "tfl" with nine packages by `tflori` and nothing else -- and on npm the letters mean Truth-Functional Logic: `tfl-lsp` is a proof checker for a formal logic course |
| ~~openFDA~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every value begins with its own field name** -- `"purpose"` is `["Purpose Pain reliever"]`, `"questions"` is `["Questions or comments? Call ..."]` -- because these are scanned drug labels and the section heading came along with the section, so rendering "Purpose: Purpose Pain reliever" is the natural thing to build. **And two of them begin with a different field's name**: `indications_and_usage` starts "Uses" and `storage_and_handling` starts "Other information", so stripping the prefix is not a fix -- the prefix is the heading printed on the box and the key is what a database called it. Also pinned: every field being an array of one, identifiers included; **finding nothing being a 404**, `"No matches found!"` with an exclamation mark, identical byte for byte to what a misspelt field name answers; a path that does not exist being Express's own HTML page; and **U+00BA MASCULINE ORDINAL INDICATOR used as a degree sign**, a character that looks right and matches no search for degrees. Stated and not served: one label of 441 matching aspirin alone, trimmed to opening sentences. Detection's near miss is that **the whole npm neighbourhood is agent tooling** -- every result for the word is an MCP server or a tool for one, with no ordinary client library at all, the first time that has been true here. Packagist answers with two other `open-` projects, OpenFeature and OpenLDAP, and the acronym belongs to more than one country: `laravel-saudi-fda` is a different regulator |
| ~~OpenAlex~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The error message is the schema** -- a filter that does not exist answers 400 with 4,779 characters and 207 field names, comma-separated, inside one sentence in one string, so the failure is larger than most successful responses and finding out which field was wrong means splitting English prose on a colon and then on commas. **And a free API quotes a price**: every response carries `"cost_usd": 0.0001` in its meta, beside an `x_query` that leaks the internal query language. Also pinned: **three failure shapes on one host** -- Flask's default HTML page for a missing work, a lone `message` for a missing entity, `{error, message}` for a bad filter; `title` and `display_name` being byte-identical; and identifiers that are URLs, so getting the bare W-number or DOI means string-splitting an address. Stated and not served: the schema message is trimmed for reading, and one work of 321,811,346 is served. Detection's near miss is **the `open-` prefix**, for the second Recipe running -- openFDA drew OpenFeature and OpenLDAP, and this one draws OpenExchangeRates four times over, beside an events API and a list of countries |
| ~~MBTA~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`type` is a number and a string in the same object** -- `data.attributes.type` is `1`, the GTFS route type, and `data.type` is `"route"`, the JSON:API resource type, two levels apart under one key name, with a third inside every relationship. **And the human-readable text moves keys between failures**: a missing resource puts it under `title` with no `detail`, and a bad filter under `detail` with no `title`, so a client rendering one shows nothing for the other. Also pinned: `status` as a string in the body beside a number on the wire, which is what JSON:API asks for; **the same error body under two content types**, `application/vnd.api+json` from one path and `application/json` from another; colours as six hex digits with no hash, so `style="color: DA291C"` renders nothing; and `direction_names` and `direction_destinations` as two arrays joined only by index. Stated and not served: one route of the network, and relationships carrying identifiers without included resources. Detection's near miss is **the same agency's previous API on its own host** -- three packages reach `api-v3.mbta.com` and three reach `realtime.mbta.com`, the retired v2 -- while the bare npm name `mbta` ("an api for your api") points only at the documentation portal, and Packagist answers the word with metadata libraries and two mutation-testing frameworks |
| ~~Carbon Intensity~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A path that does not exist answers 200 and the body says 400** -- `GET /nosuchpath` returns HTTP 200 carrying `{"error": {"code": "400 Bad Request", ...}}`, so `response.ok` is true for a URL the API does not have, and the code is a status line as a string rather than a code. A date it cannot parse answers a real 400 with the same body. **And the same fuels are named two ways on two endpoints**: `/generation` says `"gas"` and `"imports"`, `/intensity/factors` says `"Gas (Combined Cycle)"`, `"Gas (Open Cycle)"` and three separate imports, and neither vocabulary is a key of the other. Also pinned: timestamps with no seconds, which is the form its own error message asks for; and a `forecast`, an `actual` and an `index` in one object, two numbers and a category derived from one of them by a threshold nothing states. Stated and not served: **the line endings**, because the successes are pretty-printed with CRLF and a trailing space after each opening brace while the failures use a bare carriage return with no line feed at all -- classic-Mac line endings, in JSON, in 2026 -- and the key order, where the two `generationmix` entries at zero per cent have their two keys the other way round. Detection's near miss is **the ordinary word at its largest scale yet**: Packagist answers "carbon" with `nesbot/carbon`, a DateTime library and one of the most installed PHP packages there is, and npm answers it with IBM's Carbon Design System |
| ~~UK Police~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **`url` means three different things** -- a force sends a web address, a crime category sends a slug, and the telephone engagement method sends the empty string, so a client rendering `url` as a link produces a working link, a broken relative one and an empty href depending on which endpoint it came from. **And every link the force owns is plain HTTP**: its website, Facebook, Twitter and YouTube, in a payload delivered over TLS, with the RSS feed the lone `https://` exception. Also pinned: **three content types for three kinds of request, only one of them JSON** -- a force that does not exist is `text/plain` carrying the nine bytes `Not Found`, and a path that does not exist is `text/html` carrying a fragment with no doctype around it; `type` and `title` being the same word on every engagement method; and `stop-and-search` as a hyphenated key a dot cannot reach. Stated and not served: the forward slashes are escaped on the wire, and one force of the forty-odd is served. Detection's near miss is **"police" as a metaphor** -- `case-police` fixes the capitalisation of words like javascript, and `openapi-police` is schema validation -- and **the country code that means another country**: `dictionary-uk` is a Ukrainian spelling dictionary, `uk` being Ukraine's language code where the United Kingdom's is `gb` |
| ~~FHRS~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The same digit means best in one field and not in the other** -- `"RatingValue": "5"` is the top of a scale that counts up, and `scores.ConfidenceInManagement` of `0` is the top of one that counts points against you, in the same object, with the first a string and the second a number. **And absence is spelled three ways in one record**: `Phone` is `""`, `Distance` is `null`, and `AddressLine2` is `""` while `AddressLine3` holds the town, so an empty line is a gap in the middle of an address. Also pinned: **a forgotten `x-api-version` header reported as a path that does not exist**, on a path that works the moment the header is sent; that failure being **a bare JSON string rather than an object**, so `JSON.parse` succeeds and `body.Message` is `undefined` rather than throwing, with the internal handler name inside it -- the version and the path joined by a dot; coordinates as strings where the scores are numbers; and `RatingKey` packing scheme, rating and locale into `"fhrs_5_en-gb"`. Stated and not served: one establishment of hundreds of thousands, and the version header is not enforced, because requiring it would make every case send it and hide the finding. Detection's near miss is that **both words pull in a different cluster** -- "rating" brings star widgets and "hygiene" brings a pull-request linter -- and the acronym belongs to a CSS framework: Packagist answers `fhrs` with Bootstrap and Foundation |
| ~~Valet~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The response documents its own key names inside itself** -- `seriesDetail.FXUSDCAD.dimension.key` is `"d"` and `.name` is `"Date"`, and the observations carry the date under `"d"`, so a client that wants it reads the dimension first and subscripts by whatever it finds. **And the number is two levels down, under a key the caller supplied**: a rate is `obs["FXUSDCAD"].v`, the series name being a JSON key rather than a value and the thing under it an object with one field holding a string. Also pinned: **the same key differing by one letter between endpoints** -- `seriesDetail`, a map keyed by name, against `seriesDetails`, flat, with the name as a field; the terms of use as the first key of every successful body; and a bad parameter answering with a link to a generated OpenAPI operation id, underscores and all, in a URL fragment. Stated and not served: the rate moves, one series of thousands is served, and the CSV format begins with a byte order mark and packs three tables into one file. Detection's near miss is **Laravel Valet, eight packages deep** -- Packagist has nothing at all for "bank of canada" and answers "valet" with a local development environment -- while npm spreads the word to a CSS-in-JS engine and a password prompt, and two packages calling themselves "Official Bank of Canada" exchange rates carry no host, because the rates are bundled rather than fetched |
| ~~FBI Wanted~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The reward is zero and the sentence beside it says twenty-five thousand dollars** -- `reward_min` and `reward_max` are both `0` while `reward_text` offers "up to $25,000", so anything sorting by `reward_max` puts that case last and anything filtering on `reward_max > 0` does not show it at all. **And thirty-three of the record's fifty-four fields are null**, every physical descriptor among them, because the subject has not been identified -- on a record that still carries a `warning_message` reading "SHOULD BE CONSIDERED ARMED AND DANGEROUS". Also pinned: **absence spelled two ways among fields of the same kind**, eight lists `null` and one `[]`; **`detail` as an array on one failure and a string on the other**, the array carrying the validation library's own vocabulary and echoing the offending input back; two timestamps in one record disagreeing about time zones; `description` packing three facts into one string with carriage returns; a language in a field called `name`; and a URL with an `@` in the middle of it. Stated and not served: one record of 1,235. Nothing on either registry calls it, so it ships unmapped -- see the row below |
| ~~Zenodo~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The same identifier appears four times in four forms** -- `id` as a number, `recid` as the same digits quoted, `doi` as those digits prefixed, and `doi_url` as that with a resolver in front -- and a fifth, `conceptrecid`, is one less as a string, with `conceptdoi` prefixing it, so eight of the record's twenty top-level keys are the same two numbers written out differently. **And the record says it is finished three times, in three vocabularies**: `"status": "published"`, `"state": "done"`, `"submitted": true`, with nothing saying which is authoritative. Also pinned: `modified` and `updated` being byte-identical to the microsecond; **an array called `hits` inside an object called `hits`**, so reaching the records is `response.hits.hits`; and two failures using different words for one status, one about a persistent identifier and one the web framework apologising in case you typed the URL by hand. Stated and not served: key order, which the two failures differ in on the wire, and one record of 7,205,905. Detection's near miss is that **four of the six packages naming this API carry no host at all** -- including the bare name `zenodo` -- while Packagist has nothing and answers the word with a vendor prefix, a person's name (`zenorocha/clipboardjs`, one of the most installed front-end libraries there is) and Zend Framework |
| ~~UniProt~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The same acronym is capitalised two ways in two top-level keys of one object** -- `uniProtkbId` with a lowercase `kb` and `uniProtKBCrossReferences` with an uppercase `KB`, four fields apart. **And the failure echoes your URL back with the scheme downgraded**: a request over TLS answers `{"url": "http://rest.uniprot.org/uniprotkb/NOSUCHACC", ...}`, so a client that logs it or retries it drops to plain HTTP without being told. Also pinned: that failure being a **400 rather than a 404** while a path that does not exist is nginx's own HTML page naming the version it runs; `proteinExistence` packing a number and its label into one string, `"1: Evidence at protein level"`; and `entryType` as a sentence with a parenthetical naming the database twice. Stated and not served: the trailing zero on `annotationScore`, because the live API writes `5.0` and Go writes `5` -- JSON has one number type, and the digits were the point. Detection's near miss is **a substring inside a word about protection**: Packagist answers "uniprot" with two Silverstripe spam-protection modules, a WordPress unprotect plugin and a CSRF bundle. The bare npm name is a text-file parser linking the website rather than the REST host, and the rest of npm is one viewer with six adapters carrying no host at all |
| ~~Spaceflight News~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **One record links the same site twice, once securely and once not** -- `url` is `https://spaceflightnow.com/...` and `image_url` is `http://spaceflightnow.com/...`, on the same host, so a page served over TLS renders the link and shows nothing for the picture. **And the API's own 404 page carries a third-party analytics beacon**, a script tag pointing at `static.cloudflareinsights.com/beacon.min.js` inside an error response. Also pinned: **two timestamps at two precisions** in one record, `published_at` to the second and `updated_at` to the microsecond; a `limit` that is not a number being silently replaced with the default, so a client with a broken parameter never learns; absence spelled two ways four fields apart, `socials: null` and `events: []`; and two failures that do not agree on a format, one JSON with the model name capitalised mid-sentence and one HTML. Stated and not served: one article of 35,875, and the news moves. Nothing on either registry calls it, so it ships unmapped -- see the row below |
| ~~ClinicalTrials.gov~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every failure is plain text and one of them is nothing at all** -- a malformed identifier and an unconvertible parameter both answer 400 in `text/plain`, and a path that does not exist answers 404 with no body and no `Content-Type` header, so there is not even a wrong document to parse. The two that speak wrap their parameter names in **Markdown backticks, in responses declared as plain text**, formatted for a renderer that is not there. **And a date says whether it actually happened**: `startDateStruct` carries `"type": "ACTUAL"` and `primaryCompletionDateStruct` carries `"ESTIMATED"`, a fact and a guess told apart by one word a level down, in fields whose names end in a word about the code. Also pinned: `overallStatus` of `"UNKNOWN"` beside a `lastKnownStatus` of `"RECRUITING"`; `statusVerifiedDate` as a year and a month with no day; `phases` as a list of one string where the string is `"NA"`; and a sponsor's identifier carrying a space and a slash, `"30-225 ex 17/18"`. Stated and not served: one study, and three of its thirteen protocol modules. Detection's near miss is **the retired interface** -- `clinical-trials-gov` reaches `clinicaltrials.gov/ct2`, the classic site this API replaced -- and **nine of the twelve npm results are MCP servers**, the openFDA shape again, one of them about Israeli hospitals rather than this registry at all |
| ~~DataCite~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The same key holds a list in one failure and an object in another** -- a missing DOI answers `{"errors": [{"status": "404", "title": "..."}]}` and an unparseable query answers `{"errors": {"title": "parse_exception ..."}}`, so `errors[0].title` works on one and `errors.title` on the other, under a key spelled the same in both. The object one is **318 characters of the query parser's own Lucene grammar**, newlines and token names included. **And the listing says there are 133,915,241 records and 10,000 pages of one**: `totalPages` is ten thousand records' worth whatever the page size, so a client paging to the end reads 10,000 of 133 million. Also pinned: asking for page 10,001 answering 200 with `"page": 10000`; a `page[size]` that is not a number answering 200 with `"data": []` and `"totalPages": 0` beside that unchanged total; **absence spelled three ways in one record**, `container: {}`, `formats: []` and `reason: null`; the same year sent twice as `publicationYear: 2011` and `published: "2011"`; a type sent six times in six spellings where five say dataset and the sixth says `misc`; a size as `"128717 bytes"`, a number and its unit in one string; and an identifier that spans two path segments, so a router declaring `/dois/{id}` matches none of them. Stated and not served: one DOI of 133,915,241, and the counters that move. Detection's near miss is **the vendor's other API** -- `mds.datacite.org` is the older Metadata Store -- beside **the vendor's own namespace four times over**, sharpest at `datacite-rest`, which is named after this API and described, in full, as "React components" |
| ~~Deezer~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every failure is HTTP 200** -- a record that does not exist, a path that does not exist and an empty parameter all answer 200 with an `error` object, so a client checking `response.ok` sees success on every failure this API has. **And one carries a 500 inside the 200**: the code for an empty parameter is `500`, three digits that look exactly like a status, beside neighbours of `600` and `800` that are not HTTP statuses at all. **And the number of matching records depends on how many you ask for** -- the same query answers `total: 179` at `limit=1` and `limit=5`, `177` at `limit=25` and `172` at `limit=100`, so asking for a bigger page makes the collection smaller. Also pinned: a `limit=0` ignored and answered with the default twenty-five; a share link carrying Deezer's own `utm_source`, `utm_medium`, `utm_content` and `utm_term`, so an application rendering a share button propagates their attribution without deciding to; **five picture fields where the first is not a picture**, an API path beside four CDN JPEGs; explicitness sent three times in three encodings in one record; a `title_version` of `""` where absence elsewhere is a missing key; and a message reading `Unknown path components : /nosuchthing`, with the space before the colon that Deezer put there. Stated and not served: one artist and one track of the catalogue, the three totals as route constants because upstream they are not derived from the fixture either, and the preview URL's expiring token, dropped because a signed URL in a repository is a credential. Detection's near miss is **a client that reaches the right host for the one endpoint this Recipe does not have** -- Deezer's login providers call `https://api.deezer.com/user/me` and nothing else, on this host, under a path with no resource here to answer it |
| ~~iNaturalist~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **One field carries two schemes and the only broken one is the record you asked for** -- a single taxon response holds forty-five `wikipedia_url` values, twenty-one `http://` and twenty-four `https://`, and exactly one contains a raw space: `https://en.wikipedia.org/wiki/Apis mellifera`, the top-level record, the only one a client reading `results[0]` will ever touch. Every nested one is escaped. An `<a href>` works because browsers repair it, so the bug survives every manual check. **And every fetch is a search**: `/v1/taxa/47219` answers with a paging envelope, and a thing that does not exist answers 200 with `total_results: 0` rather than 404 -- and the two disagree about the page size, `per_page: 30` when something was found and `per_page: 1` when nothing was. Also pinned: **the one real 404 on the host is Express's own HTML page**, `Cannot GET /v1/nosuchthing` in a `<pre>`, on a host whose every other response is JSON; `wikipedia_summary` as HTML in a field called summary, with an en dash in its range; the ancestry sent twice in two types, as `ancestor_ids` and as a slash-joined string; and `conservation_status: null` beside a `conservation_statuses` holding fourteen. Stated and not served: one taxon with two of its fifteen ancestors and none of its children, only the first sentence of the summary, and the observation total, which rose by fourteen between two requests seconds apart. Detection's near miss is **the dataset rather than the service** -- `node-red-contrib-oscar-plants-classification` classifies images against the iNaturalist corpus and never asks iNaturalist anything |
| ~~RCSB~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Two of its fields are the letters Y and N, and `"N"` is true** -- `has_released_experimental_data` is `"N"` and `pdb_format_compatible` is `"Y"`, single-character strings where a boolean belongs, and the string `"N"` is truthy in every language here, so the yes branch runs on an entry that said no. **And the identifier you get back is not the one you sent**: ask for `4hhb` and the record answers `"rcsb_id": "4HHB"`, and the failure for a nonsense id says `No data found for entryId: NOTANID` in capitals it was never given. Also pinned: **two 404s sharing one host and one key**, the application's `{status, message, link}` against Spring Boot's `{timestamp, status, error, path}` with no `message` at all; the experimental method sent twice as `"X-RAY DIFFRACTION"` and `"X-ray"`; **CIF capitals among snake_case siblings**, `cell.Z_PDB` beside `cell.angle_alpha` and `symmetry.Int_Tables_number`, so a client that lower-cases its keys destroys three of them; a resolution of `[1.74]`, a lone number inside an array; a version split across `major_revision` and `minor_revision`; and a 1984 deposition's punched-card formatting, a block of capitals wrapped at sixty columns inside a JSON string. Stated and not served: one entry of the Protein Data Bank, and the Spring failure's timestamp, which is the moment it was generated upstream and a fixed recorded value here. Nothing on either registry is a client of this API, so it ships unmapped -- see the row below |
| ~~Europe PMC~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Three failures have three shapes and the HTTP status is right on exactly one of them** -- a search with no criteria answers 200 carrying `{"errCode": 404, "errMsg": ...}`, a bad page size does the same, and a path that does not exist answers a real 404 with `Content-Length: 0` and nothing in it, so `response.ok` is true for both failures that explain themselves and false for the one that does not. **And asking for a format that does not exist gets you XML**: `?format=xxx` answers 200 with `application/xml`, telling you in a document your JSON client cannot parse that you should have asked for `"json"`, with a space left before the closing tag. **The same document sends a real boolean in one object and letters in the next** -- the echoed request carries `"synonym": false` while eleven fields on the record beside it are `"Y"` and `"N"`, none of which is falsy anywhere. Also pinned: `journalIssn` as `"0305-1048; 1362-4962; "`, two identifiers in one string with a trailing separator and a trailing space; `pubType` joined the same way with commas inside its items, beside an `authorString` that uses the comma the other way; **an array wrapped in an object under a singular key inside a field named for a list**, `"fullTextIdList": {"fullTextId": [...]}`; and `journalVolume` of `"41"` as a string beside an `issue` of `"database issue"`. Stated and not served: one record of Europe PMC's millions, fetched by a query for its own identifier so the result does not move, and the XML failure served as recorded bytes because Cauldron has no XML writer. Detection's near miss is **the vendor's own npm scope holding something unrelated** -- `@europepmc/express-middleware-minio` stores files in Minio |
| ~~OpenAIRE~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The JSON is XML wearing a costume** -- every text value under `"$"`, every attribute under a key beginning `"@"`, and element names keeping their XML namespaces, so reading a title means `result.metadata["oaf:entity"]["oaf:result"].title[0]["$"]` and no client can use dot notation at all. **And `total` is a number when it is not zero and a string when it is**: one match answers `{"$": 1}`, no match answers `{"$": "0"}` with `"results": null`, so the two paths through the response share neither type. Also pinned: **the useful sentence in the wrong field** -- a page size out of range answers `message: "400 - Illegal argument exception."`, restating the status a client already has, while `exception` holds the only text that says what happened, beside a `code` of `"400"` as a string and a `trace` of `null`; and **a 404 that is Apache Tomcat's own error report, naming its version -- 7.0.68, from 2016 -- twice**, in the title and at the foot, with a `message` field of `<u></u>`, an empty underline element. Stated and not served: one publication of 237,304,810, fetched by its DOI so the result does not move; the twenty-one children of the real `oaf:result`; and the Tomcat page's inline stylesheet. Nothing on either registry is a client of this API, so it ships unmapped -- see the row below |
| ~~ROR~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A path that does not exist says your authentication is missing, on an API that has none** -- `GET /v2/nosuchthing` answers `403 {"message": "Missing Authentication Token"}`, the gateway's default for an unmatched route, sending whoever typed the path wrong to look for an API key that does not exist. It is a 403 rather than a 404, so a client branching on the status treats a typo as a permissions problem. **And the identifier is a whole URL inside the path**: `/v2/organizations/https://ror.org/02mhbdp94` returns a record whose `id` is that same address, so a client that percent-encodes its path segments correctly gets a 404 and one that does not gets the record. Also pinned: **two failure shapes where the one that is not your fault answers 200** -- an invalid id is 404 with `{"errors": ["..."]}`, an array of bare strings so `.errors[0].message` is undefined, and a page outside the range is that same shape at 200; **two schema versions in one record**, `"1.0"` for when it was created and `"2.1"` for when it was last touched; a display name chosen by finding `"ror_display"` inside an array inside an array element; a founding year as a bare integer; and links as typed pairs where the organisation's own website is `http`. Stated and not served: one organisation with two of its names, two links and one relationship, and the `locations` and `external_ids` arrays. Detection's near miss is **the acronym, three ways** -- `@rork8s/ror-resources` is a Kubernetes project, `react-ror` and `luniki/trails` are Ruby on Rails, and six of Packagist's eight results just begin with somebody's name |
| ~~Nobel Prize~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A laureate who does not exist is an array of length one** -- `/2.1/laureate/999999` answers 200 with `[{"meta": {...}}]`, the licence notice and nothing else, so a client checking `.length` is told yes and finds out only by reaching for a field. There is no 404 on that path at all. **And the language maps have keys in them that are not languages**: `city` is `{en, no, se}`, `cityNow` is `{en, no, se, sameAs, latitude, longitude}`, `nameNow` is `{en}` where its sibling has three, and `knownName` has two -- four objects in one record and no two agreeing, so `Object.keys(cityNow)` hands an array of URLs and two coordinates to whatever renders a city. Also pinned: **the record links to a different version of the API that served it**, `links[0].href` pointing at `/2/laureate/745` on a record fetched from `/2.1/`; coordinates as strings keeping a trailing zero each, `"40.825930"`, a precision a JSON number cannot carry; a share of a prize as the string `"1/3"` beside money as a bare integer with no currency named anywhere; a `date` and the `year` of that date, the year as a string; and the licence repeated inside every record rather than sent once. Stated and not served: one laureate of 1,018, with one prize and one affiliation, and the category's Norwegian name, which carries a letter the Recipe file does not. Nothing on either registry is a client of this version, so it ships unmapped -- see the row below |
| ~~PDBe~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The response is keyed by what you asked for, and not as you asked it** -- `/pdbe/api/pdb/entry/summary/4HHB` answers `{"4hhb": [{...}]}`, a map with one key, and that key is your own query lowered, so `res["4HHB"]` is undefined and no field anywhere in the response holds the identifier. **The same record from RCSB does the opposite**: that Recipe pins `rcsb_id: "4HHB"` for a request that said `4hhb`. Two organisations, one deposition, normalisation in opposite directions -- and they disagree on the date too, `"19840307"` here against `"1984-03-07T00:00:00.000+00:00"` there. Also pinned: **three failures in three shapes, one of them a 405 for a GET** -- a missing entry 404 under `message`, a missing path 404 under `detail`, and the same path without an identifier answering `405 Method Not Allowed` on an endpoint that takes GET and only GET; **a JSON key with a slash in it**, `"dna/rna"` beside the `dna` and `rna` it joins; the experimental method as two arrays in two cases, `["x-ray"]` and `["X-ray diffraction"]`; and absence spelled two ways, `deposition_site: null` beside two empty arrays. Serving it needed a new list shape: `entry_style: list` inside a map, because `res["4hhb"].title` is undefined and `res["4hhb"][0].title` is the answer, and the format could previously claim only the first. Stated and not served: one entry and only its summary, of the dozen endpoints PDBe answers about it. Detection's near miss is **the other organisation's viewer** -- `@rcsb/rcsb-molstar` reaches `data.rcsb.org/graphql`, the other protocol, where these four reach this API |
| ~~Jisho~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **There is no failure at all** -- a search that finds nothing and a search with no keyword in it, the one parameter the endpoint exists for, are both answered 200 with an empty list, exactly as a working search is answered 200 with a full one. **And `meta.status` is a field that can only hold one value**: 200 on every answer, because the only non-200 on the host is an HTML page from a path that never reaches this code, so a client checking it is checking a constant that looks like a guard. Also pinned: **sibling objects in one array with different key sets** -- three senses identical in shape except that the third has `sentences` and the first two do not, so mapping over them and reading `.sentences.length` throws on two of three; **`attribution` holding two booleans and a URL**, so `Object.values` gives `[true, true, "http://..."]`; a `parts_of_speech` of `["Wikipedia definition"]`, which says where the definition came from rather than what kind of word it is; `jlpt` of `["jlpt-n5"]`, the field's own name repeated inside its value; a `tags` of `["wanikani7"]`, a different company's learning app and its level, in a dictionary's public tags; and links pinned to a 2012 Wikipedia revision over plain http, with curly quotation marks in the text beside them. Stated and not served: one word of the twenty the query returns, and the Japanese link's path, which arrives with an unescaped Japanese character and is served percent-encoded. Detection's near miss is **the client that scrapes the site instead** -- `jishon` promises "json from jisho" and reaches `jisho.org/search`, the web page, to get JSON this API already returns |
| ~~Stack Exchange~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every failure is 400 and one of them carries a 404 in the body** -- a path that does not exist answers `400` on the wire with `{"error_id": 404, "error_name": "no_method"}` in the payload, so a client switching on the status cannot tell "you asked wrongly" from "there is nothing there", and one switching on `error_id` disagrees with its own transport. **And `error_message` is sometimes a sentence and sometimes just a parameter name**: a bad site says ``No site found for name `nosuchsite` ``, Markdown backticks in a response that is not Markdown, while a bad page size says, in full, `"pagesize"`. Also pinned: **the rate-limit budget in the body of every answer**, `quota_max` and `quota_remaining` beside the results, so the payload differs on every request even when the data does not; **no total anywhere**, only `has_more`, so a pager cannot be built from it; every timestamp as a Unix epoch integer with `last_activity_date` and `last_edit_date` the same number; a licence on every item rather than on the response; and a tag of `c++`, which a client building a URL must escape. Stated and not served: one question, and its counters move -- they are pinned as they stood when checked, because a Recipe that omitted them would describe a question with no score. `quota_remaining` is a constant here and a countdown upstream; no fixture can be right about it. Detection's near miss is **a whole ecosystem one version behind**: `stackexchange-api` reaches `/2.2`, and two Packagist clients say 2.2 in their own descriptions, where this Recipe describes 2.3 |
| ~~Kraken~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **Every success carries an empty error array, and an empty array is true** -- `{"error": [], "result": {...}}`, so `if (res.error) throw` fires on every success and never on a failure, in JavaScript, Ruby and PHP alike, and the only reading that works is `res.error.length`. **And a failure answers 200 with no `result` key at all** -- not null, not empty, absent -- so `res.result.XXBTZUSD` is a TypeError rather than a message anyone could show; one failure is 200 and another 404 in the same shape, both messages a machine code and a sentence joined by a colon. **You ask about `XBTUSD` and are answered about `XXBTZUSD`**, and no field anywhere holds either name -- it is only the key. Also pinned: **nine single-letter fields in three shapes**, `a` and `b` three-entry tuples, six two-entry tuples, and `o` a bare string where its eight neighbours are arrays; **every number a string except the one that is not**, `t` holding two bare integers among seven fields of strings; and trailing zeros kept, `"77000.00000"`, a precision a JSON number could not carry. Stated and not served: one pair of the several hundred, and the prices move every second -- they are recorded as they stood at the verified moment, and the shape is what does not move. Detection's near miss is **the name itself**: "Kraken" is at least five unrelated products, and the exchange holds the bare name on neither registry -- npm's `kraken` reaches `api.kraken.io`, the image optimiser |
| ~~Scryfall~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The prices are five nulls and one string** -- six fields about money, five holding nothing and the sixth holding `"46.91"`, so summing a basket means parsing every value that is not null and finding out which those are means checking all six. Scryfall is a well-built API, which is what makes it worth describing: the traps are decisions rather than carelessness. **Four failures, four typographies, and one of them is minified** -- curly double quotes on one, Markdown backticks on another, `You didn‘t enter anything` on a third where the apostrophe is U+2018, a *left* single quotation mark and the wrong one for a contraction, and plain prose on the fourth. Also pinned: **sibling failures with different key sets**, the minified one carrying `"warnings": null` where its three siblings have no such key, so `err.warnings` is undefined on three and null on one; `object` as a discriminator on everything, `"card"`, `"error"` and `"related_card"`; the card's own `scryfall_uri` ending `?utm_source=api`, so rendering "view on Scryfall" passes on their attribution; a `collector_number` of `"4"` kept as a string; a rarity of `bonus`; and five produced colours sorted alphabetically rather than in the order the game itself always writes them. Stated and not served: one card of tens of thousands and only the fields the cases name; the whitespace, since Cauldron writes JSON its own way; and `cmc`, which arrives as `0.0` and is served as `0` because a zero float cannot carry a decimal through Go's writer -- a mutation about it was dropped for that reason rather than left looking like an unchecked claim. Detection's near miss is **types without a client, twice, one of them in the vendor's own scope** |
| ~~Lichess~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The game counts do not add up** -- 23,723 in total against 11,329 won, 11,345 lost and 1,044 drawn, which is 23,718, on a record where `playing` is 0. Five games are unaccounted for, nothing names them, and there is no fifth field to blame; each number is right on its own and only the sum is wrong, which is why it survives. **And a list of links is separated by a Windows line ending inside a JSON string**, `"github.com/ornicar\r\nmas.to/@thibault"`, so splitting on `\n` leaves a carriage return that renders as nothing and compares as something. Also pinned: **two time units in one record, neither named** -- `createdAt` and `seenAt` in epoch milliseconds beside a `playTime.total` in seconds, so a client dividing the wrong one is out by a factor of a million; **sibling objects with different key sets**, `prov` sent only when true so one rating has it and the next does not, across eighteen game types; **two 404s, one JSON and one the whole website in HTML**, the second opening with a comment about the project being open source; and a colour as the integer 10, with nothing saying what the other nine are. Stated and not served: one user, whose counts move as they play -- the arithmetic gap is the finding rather than the figures -- two of the eighteen perfs, and the HTML page as a fragment rather than its tens of kilobytes. Detection's near miss is **the vendor's own scope, three times over and none of them a client**: a PGN viewer, a WebAssembly engine and the chessboard itself |
| ~~Gutendex~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **A filter it cannot parse is dropped, and you get the whole catalogue** -- `?ids=84` answers one book and `?ids=abc` answers 79,296, every book Project Gutenberg has, at 200, with nothing saying the filter was ignored. **And the invalid parameter is carried faithfully into the next link**, `?ids=abc&page=2`, so a client paging by the link it was given walks the entire catalogue still asking for "abc" -- the parameter preserved exactly where it does nothing and dropped exactly where it would matter. Also pinned: **one array with two naming conventions**, four `bookshelves` entries prefixed `Category: ` and then four that are not, with nothing marking the boundary; **MIME types as JSON keys**, carrying slashes, plus signs and a semicolon-and-space in `"text/plain; charset=utf-8"`, not one of which can be written after a dot; **the same title twice in two capitalisations**, the field disagreeing with the summary quoting it; a Library of Congress subdivision separator arriving as `-- ` inside a subject string; an author written surname-first; and a `copyright` boolean for a fact that is really true, false or unknown. Stated and not served: one book of seventy-nine thousand; the catalogue count, which read 79293 and 79296 within a minute of itself, so what the case pins is that the filter was dropped rather than the figure; and three of seven subjects and six of eight bookshelves, both a contiguous prefix so an index means what it means upstream. Detection's near miss is **the name itself** -- WordPress's block editor is called Gutenberg, so npm answers "project gutenberg api" with six `@wordpress` packages and Packagist answers "gutendex" with six WordPress plugins |
| ~~Wikimedia~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The error message is inside a map keyed by language** -- every failure carries `messageTranslations` and none carries `message`, so `err.message` is undefined on all of them and the sentence is one level down under a code the client must choose; choosing wrongly is silent. **And three failures share no key set at all**: a missing page sends four keys, a missing path sends two with no message in any language -- not an empty map, no map -- and an empty search sends nine, two of which name the same failure and disagree, `error` saying `parameter-validation-failed` where `errorKey` says `missingparam`. Also pinned: **the record links to a different host running the same API**, served from `api.wikimedia.org/core/v1` and pointing at `en.wikipedia.org/w/rest.php/v1`; **two identifiers four orders of magnitude apart**, `id` 9228 for the page beside `latest.id` 1372042770 for the revision, with nothing in either name saying which is which; and a licence URL ending in a language, `.../by-sa/4.0/deed.en`, so the link rendered as "the licence" is one translation of it. Also observed and not served: a language code naming no wiki does not fail but redirects, in HTML, to the Incubator where such wikis are drafted -- and that body ends with an HTML comment naming the PHP file that produced it. Stated and not served: one page, whose latest revision moves whenever anyone edits Earth. Nothing on either registry is a client of this API, so it ships unmapped -- see the row below |
| Codecov | Assess — coverage reports and the commit they attach to |
| ~~SonarCloud~~ | Shipped. The quality gate has three outcomes and one of them means nobody set a gate: "The different statuses returned are: OK, WARN, ERROR, NONE. The NONE status is returned when there is no quality gate associated with the analysis." So `if (status !== "OK") fail()` fails the build for an ungated project and `if (status === "ERROR") fail()` passes everything on one. And **every number inside the gate is a string** -- the published example sends `"errorThreshold": "85"` beside `"actualValue": "82.50562381034781"`, with the direction of the comparison in a third field as `"comparator": "LT"` -- so evaluating a condition means parsing two strings and reading an operator, and comparing them as text is right on "14" against "0" and wrong on "9" against "10". Also pinned: **there is no verb but GET and POST** (87 of the 156 actions are POST, 69 are GET, and the description cannot express another -- each action carries a boolean called `post`, so deleting a comment is `POST api/issues/delete_comment` and the path is the verb); the analysis and the gate being different objects, since a Compute Engine task that FAILED produced no gate result at all; a failed task handing back a **Java stack trace inside a JSON string** beside a boolean saying whether to read it; and an issue carrying two names for its status (`issueStatus`/`status`) and two for its effort (`effort`/`debt`), agreeing today. Written against the description SonarCloud serves of itself at `api/webservices/list` with response examples at `api/webservices/response_example` -- and `project_status` carries `"deprecatedSince": "16 September, 2025"` there while its own response says nothing, so the only way to learn the endpoint is deprecated is to call a different one that describes it. Detection found an exclusion kind this collection had not recorded: packages that write a file for another program to upload, led by `vitest-sonar-reporter` at 841k downloads and no network request at all |
| Railway | Assess. Here for the same lapsed reason as New Relic: GraphQL-only was an exclusion before anything here spoke GraphQL, and three Recipes do now |

## Secrets and configuration

| Provider | Why |
|---|---|
| ~~HashiCorp Vault~~ | Shipped, and this row named the headline before it was written: the same read gives you the secret, or a box containing the secret. On KV v1 `response.data` is your secret; on KV v2 `response.data.data` is, and `response.data` is a box holding it beside a metadata block. Nothing in the path says which, because the version is a property of the mount. `secret = response.data` does not throw against v2 -- it yields an object, truthily, and surfaces later as a password that is an object. Written against openapi.json in `hashicorp/vault-client-go`, the document HashiCorp generates its own Go client from, which turned out to omit three things every real request and response has: the `/v1/` URL prefix (stated in `info.description` and on none of its 715 paths), the response envelope (the document's `KvV2ReadResponse` is the *inner* half; the wrapper carrying `request_id`, `lease_id`, `renewable` and `lease_duration` is hand-written outside the generated code as `Response[T any]`), and authentication (`securitySchemes` is `{}`, empty, for an API whose purpose is holding secrets). The generator emits the middle and a person supplies both ends. Also pinned: a write answers with metadata and not with the secret -- the request schema says so in passing, "will be stored and returned on read" -- while the v1 write answers 204 with no body at all; not-deleted is `deletion_time: ""` beside `destroyed: false`, an empty string rather than a null; and a version is an integer in `current_version` and a string when the same number is a key in `versions`. Stated and not served: `/sys/health`, whose five declared responses are 200 active, **429 "unsealed and standby"**, **472** (not an HTTP status code) for a DR replication secondary, **501 "not initialized"** and 503 sealed -- codes chosen so a load balancer's "2xx is healthy" rule finds only the active node, which means every generic monitor reports a working standby as rate-limited. Those are states of the server rather than properties of a request. Also unserved: listings, which answer with an array of names rather than records. Detection found the most contested word in the collection -- four vendors ship a Vault and the biggest is Azure's |
| ~~Doppler~~ | Shipped. The credential is checked before routing for every method and path, so it never tells a caller their path was wrong, only that they were |
| ~~1Password~~ | Shipped as Connect, and the row was right that they are separate APIs. The finding is in 1Password's own Go client, which contradicts itself in one file: one method treats a listing as complete items and another treats the identical listing as summaries and re-fetches each, because the list endpoint omits the fields holding the secret values |
| Infisical | Assess — secrets, environments, folder scoping |
| Terraform Cloud | Assess — runs, plans, applies. A run is a state machine with a confirmation step in the middle |
| Pulumi Cloud | Assess — stacks, updates, previews |

## Webhooks, queues and workflow

| Provider | Why |
|---|---|
| ~~Svix~~ | Shipped. Delivery attempts, backoff, and the endpoint disabled after repeated failure |
| ~~Knock~~ | Shipped. A trigger answers with a run id and nothing else, one trigger becomes several messages that disagree, a message that was never sent is still a message, and archived is not read |
| ~~Courier~~ | Shipped. Courier picks the channel at delivery time, so which one a notification went out on is a separate lookup rather than anything in the request or the response; routing that finds no channel is not an error; and SIMULATED is a real status |
| ~~Novu~~ | Shipped. A subscriber exists implicitly on first trigger, so a typo in an id creates one with nowhere to send and everything reports success; acknowledged means received rather than sent; and Novu's identifier and yours are different fields |
| ~~Hookdeck~~ | Shipped. "Did my webhook work" has three answers here and they disagree on purpose, because one arriving webhook becomes three objects: a Request is what arrived, an Event is what gets delivered, an Attempt is one try. **A request can arrive, verify, be acknowledged with a 200, and go to nobody** -- `rejection_cause` has nine values and one is `NO_CONNECTION`, with `events_count` at nought, and nothing anywhere failed. **A failed attempt leaves the event SCHEDULED**: `EventStatus` is SCHEDULED, QUEUED, HOLD, SUCCESSFUL, FAILED, CANCELLED while `AttemptStatus` is only FAILED and SUCCESSFUL, so FAILED on an event means the retries are exhausted and an alert on `status === "FAILED"` is silent through every failure until then. Also pinned: `webhook_id` is not the id of a webhook -- its own description is "ID of the associated connection (webhook)", the product having renamed webhooks to connections while the field kept the old word; **32 of the 33 attempt error codes mean the request never reached you** (DNS, TLS, resets, redirects) with `BAD_RESPONSE` the only one meaning your server answered, so those attempts carry `response_status: null` and counting failures by `response_status >= 500` counts none of them; `verified` is a *nullable* boolean, so never-checked and failed-verification are different states that `if (!verified)` merges; and every listing in the API holds its records under `models`, the same key whatever it contains. The error envelope declares its own HTTP status as `"type": "number", "format": "float"`. Detection found a shape this file had not recorded: **the vendor's other product under the vendor's own scope** -- `@hookdeck/outpost-sdk` is Hookdeck's self-hosted Outpost, a different API, which is why an npm scope cannot be trusted the way `@googleapis/` can |
| ~~Inngest~~ | Shipped. **Its failures declare `text/plain` and send JSON** -- a client trusting the content type treats valid JSON as a string, and one calling `.json()` succeeds despite being told not to. Two hosts take two credential systems and neither says when you have confused them: an event key goes in the path at `inn.gs`, a signing key in a Bearer header at `api.inngest.com`, and the same wrong key asked for two ways gives a different status and a different capitalisation. The caller's own deduplication id is never the key the event is filed under, so a client that remembers what it sent cannot find what it stored |
| Trigger.dev | Assess — runs, tasks, the resumed run |

## Scheduling and real-time

| Provider | Why |
|---|---|
| ~~Cal.com~~ | Shipped. A slot listing is a view rather than a reservation and has no identifier to hold on to, so booking one that was free is a 400 that reads like your bug; cancelling does not delete; and pending holds the slot without anybody agreeing |
| Acuity Scheduling | Assess — appointments, types, intake forms |
| ~~Daily.co~~ | Shipped as `daily`. An expired room is not deleted, it becomes what the reference calls a zombie |
| ~~LiveKit~~ | Shipped. A bad credential outranks routing and an absent one does not |
| ~~Agora~~ | Shipped. Five broken credential shapes collapse into the rejected sentence and only a truly absent header gets its own |
| Stream | Assess — chat channels, members, the message that is soft-deleted and still returned |
| SendBird | Assess — channels and messages |

## Lifecycle and marketing messaging

| Provider | Why |
|---|---|
| ~~Customer.io~~ | Shipped. **Its API answers its own marketing site's 404** -- a 3817-byte HTML page, byte-identical across both API hosts and both kinds of mistake, so a typo in an API path returns a web page and a client parsing it as JSON throws on the doctype. Two hosts take two credential schemes, Basic for tracking and Bearer for the app API, and sending the app's token to the tracking host answers exactly what sending nothing answers: missing, wrong and wrong-scheme are three mistakes with one response, separated only by a `WWW-Authenticate` header one host sends and the other does not. One authentication scheme is declared per Recipe with no per-route override, so this declares Basic and says plainly that the app route's enforcement is unfaithful while the body it serves there is byte-correct |
| ~~Braze~~ | Shipped. "message" is both the success channel and the error channel, a listing entry is a smaller object than the detail, and the export answers with a prefix instead of data |
| Iterable | Assess — users, campaigns, catalogues |
| ActiveCampaign | Assess — contacts, deals, automations |
| ~~Brevo~~ | Shipped. It **distinguishes a missing key from a wrong one and only one can be served**: both are 401 with the same `code` of `unauthorized`, differing only in prose, so the field a machine would switch on is identical while the sentence is not. **Its published description is itself behind the credential** -- fetching `swagger_definition.yml` without a key answers 401 -- so the document explaining how to authenticate cannot be read without authenticating, and no `upstream.spec` is recorded. Documented and not called: the transactional send returns a `messageId` shaped like a mail Message-ID rather than a bare token, and the matching webhook field is spelled `message-id` with a hyphen, unreachable with a dot in most languages |
| ~~Kit (ConvertKit)~~ | Shipped as `recipes/kit`. Listed twice under two spellings of one company, which the rename made easy to do; see the other row for what it found |
| Attentive | Assess — SMS subscribers and consent state, where consent is legally load-bearing |
| Beehiiv | Assess — publications, posts, subscribers |
| Loops | Assess — contacts and transactional sends |

## Social and content platforms

| Provider | Why |
|---|---|
| ~~Reddit~~ | Shipped, and small, because the reachable surface is small. Every endpoint that would answer a question about a post was blocked live on all three hosts |
| ~~Twitch~~ | Shipped. **It needs two credentials and checks one of them first** -- a correct `Client-Id` with no token answers byte-identically to sending nothing at all, so getting half of it right earns no acknowledgement. The last page of a listing carries `pagination` as an empty object rather than omitting it or nulling the cursor, so a client testing for the key loops forever and one testing the cursor does not |
| ~~Bluesky~~ | Shipped, and it did belong here after all. The same profile by handle and by DID returns byte-identical documents |
| ~~Mastodon~~ | Shipped. The reference and the wire disagree about what a missing status says, on the same server, on the same day |
| Telegram Bot API | Assess — every method is both GET and POST, errors come back with HTTP 200 in some client libraries, and updates arrive by long polling or webhook but never both |
| Buffer | Assess — profiles, updates, scheduling |
| ~~Spotify~~ | Shipped. **The body cannot tell two failures apart and the header can** -- a missing and a bogus token give identical JSON naming three possible causes at once, while `WWW-Authenticate` carries `missing_token` against `invalid_token`, so the diagnosis exists one layer above where anyone reading JSON looks. An unknown path answers **410 Gone**, a status meaning the resource existed and was removed, for one that never existed. Its description marks two endpoints deprecated while the wire sends no `Deprecation` or `Sunset` header at all. The token-refresh path this row asked about needs an account and is not modelled |
| ~~Strava~~ | Shipped. Four credential failures and no 405 anywhere, so the whole API has two sentences for everything a caller can get wrong about reaching it. Neither rate-limit window's headers appeared on any refusal collected |

## Identity verification and risk

| Provider | Why |
|---|---|
| ~~Persona~~ | Shipped. The decision arrives by webhook minutes after the inquiry is created |
| ~~Onfido~~ | Shipped. A check is complete and its report can still be consider rather than clear |
| ~~Veriff~~ | Shipped. The signature is never examined until the client identifier passes |
| Middesk | Assess — business verification and its partial matches |
| ~~Alloy~~ | Shipped. Three unauthenticated failures across two statuses |
| Sift | Assess — scores, decisions, the workflow that runs server side |

## Travel

| Provider | Why |
|---|---|
| ~~Duffel~~ | Shipped, and this row's premise is confirmed and unmodelled. **Offers do expire**, typically within thirty minutes, after which they answer 422 `offer_expired` -- and there is no way to expire an identifier on a clock here, so one fixture offer always answers expired and the real timer is described rather than simulated. What was verified instead: Duffel **refuses you for not naming a version before it checks who you are**, with distinct codes separating not saying from saying something retired |
| ~~Amadeus~~ | Shipped, with no live case. Both hosts are in DNS and carry no address record |
| ~~Hotelbeds~~ | Shipped. The stale-rate rejection publishes the tolerance the price was allowed to move within |

## Health

| Provider | Why |
|---|---|
| Redox | Assess — the data model is HL7 and FHIR shaped rather than REST shaped, which may make it a poor fit for this format |
| Metriport | Assess — medical records and device data behind one API |
| Health Gorilla | Assess — patient queries and document retrieval |

## Web data

| Provider | Why |
|---|---|
| ~~Firecrawl~~ | Shipped. Partial results are readable before the crawl finishes |
| ~~Apify~~ | Shipped. Assess — actor runs, datasets, the run that succeeds with zero items |
| ~~ScrapingBee~~ | Shipped. **The target's status survives in a header until you ask for it not to.** By default a scrape answers 200 whatever the target sent, and the real status lives only in `Spb-initial-status-code` -- so `response.status` says nothing about the page requested. Then `transparent_status_code` makes the proxy's own status *become* the target's, which destroys the one signal separating "ScrapingBee could not route you" from "the page does not exist". Its own failures are two shapes: a missing key is 400 with the field name nested three deep and repeated, a wrong key is 401 echoing the credential back. A path it never routed answers Werkzeug's HTML 404 with or without a key -- the one failure not gated behind the credential |
| ~~Browserless~~ | Shipped, and it is the sharper half of the proxy question ScrapingBee's Recipe asks. **It has no field for the target's status at all** -- not a header, not a body field -- so a page that 404d and a page that rendered are the same successful response, with nothing anywhere to separate them. Its own failures split by layer rather than status: no token is 401 as openresty's HTML from the edge, a wrong token is 401 as plain text from the application, and Browserless's own documentation describes them as one case. Sending the credential as an `Authorization` header, the convention nearly every other provider here uses, **crashes the edge with a 500**, reproduced three times |

## More commerce and billing

| Provider | Why |
|---|---|
| ~~Toast~~ | Shipped. businessDate is the day the money belongs to and openedDate is when the order happened, and they differ every night after midnight |
| ~~Medusa~~ | Shipped, Storefront API, and written against the published OpenAPI description at 2.19.0 rather than from memory. A field name without a + replaces the entire default set, so ?fields=id returns an order with one field and no error |
| ~~Wix~~ | Shipped, written against the API definition Wix ships inside its own reference pages -- a megabyte of OpenAPI with 632 component schemas, per-method triggered events and per-method error tables. The shipping address is not where it ships to: of shippingInfo.logistics.shippingDestination, "For pickup orders, this contains the pickup location address, not the recipient's address", so a label printed from the field with shipping in its name twice goes to the depot the customer was going to walk to. The customer is in recipientInfo, which says nothing about shipping. Also: importing an order triggers order_imported and never order_created, so a migration is silent to every handler subscribed to the obvious event; an order with nothing to ship is FULFILLED the moment it exists, because "Orders without shipping info are fulfilled automatically"; archived hides an order from the dashboard and not from the API, which is the inverse of Ecwid's and Squarespace's shape; taxSummary says it will be removed on September 30, 2024 and is still in a document updated in 2026; and two enumerations in the one document are called PaymentStatus and do not share a vocabulary. Detection found the busiest exclusion list here -- WiX is also the Windows Installer XML toolset, wixel/gump is a PHP validator with a million downloads, and the most installed package with the word in its name is Wix's own LLM test runner |
| ~~Squarespace~~ | Shipped, written against the OpenAPI 3.1.1 document Squarespace publishes and serves its own reference from. One payment state destroys the number that gives it meaning: "a partial refund of any amount sets the order to REFUNDED", so a two-pound goodwill refund on a four-hundred-and-eighty-pound order sits in the state a fully reversed order sits in, and the documentation's own remedy is arithmetic the response does not do -- "compare the order's refundedTotal with grandTotal". Also: an id that could exist and does not is a 404 while one that could not exist at all is a 400, which is a distinction no Recipe here could make until this one; test orders come back in the same listing as the money with no parameter to exclude them; and MonetaryAmount.currency is declared as an object with six properties and sent as the string "USD" by every example in the same document. Two rules are documented and not served, and the file says so: the cursor may not travel with any other parameter, and the order.create webhook for a Payment Plan fires only when the last instalment lands |
| ~~Ecwid~~ | Shipped, written against the documentation Ecwid publishes as source in its own GitHub organisation. The order list leaves out a whole category and does not say so: "If no filters are set in the URL, API will return all orders except for unfinished orders." Every number in the response agrees with every other and the set is short. Also: acceptMarketing is a boolean where null means yes, discount is everything except the coupon discount, and paymentMessage is cleared when the payment lands |
| ~~Royal Mail~~ | Shipped, written against the Swagger document Royal Mail serves from its own API host. Deleting an order can cost money, and the description says so: cancelled label information goes to Revenue Protection, and a cancelled label found on a parcel is charged with a handling fee. The path parameter is a list of up to a hundred identifiers separated by semicolons, where references must be quoted and percent-encoded -- and the document's own example is a reference full of semicolons |
| ~~Apideck~~ | Shipped, written against the OpenAPI document Apideck publishes, including its x-apideck-gotchas extension. The first unified API here -- one shape in front of Shopify, Walmart, TikTok, Wix and the rest -- and the seams are the subject: a 200 can carry meta.warnings holding somebody else's 429, present only when a step degraded. Also: a valid request can be unsupported by the connector, every response names which shop answered, and a unified field can be reported by one provider and inferred by another |
| ~~Voucherify~~ | Shipped, written against the OpenAPI document Voucherify publishes and generates its own SDKs from. Promotions are the first of their kind here, and the headline is what the domain does to the word failed: a failed redemption is a stored object, listed and counted like any other, with the reason in a field rather than a status code. Also: result and status are two verdicts on one object and status has a value result cannot express, a null quantity means unlimited, and the envelope names the key that holds its own data |
| ~~Akeneo~~ | Shipped, written against the Swagger description Akeneo publishes in its own docs repository, with fixtures built from the worked example beside the product listing. A product has no name: every attribute's value is an array of {locale, scope, data} entries, so one description is four rows in a two-locale two-channel catalogue and reading it means matching both. Also HAL at _embedded.items, a next link the document says never to construct by hand, and prices as unscaled decimal strings |
| ~~Allegro~~ | Shipped, written against the OpenAPI description Allegro publishes. An order exists before it is paid and its data can still change afterwards: the endpoint is /order/checkout-forms, and BOUGHT means the buyer has not filled the form in while FILLED_IN means, in the document's words, that payment is not complete "so data could still change". Also: two status machines on one object, the API version in the Accept header with a 406 for anything else, and two error envelopes |
| ~~Saleor~~ | Shipped, written against Saleor's own repository -- the schema, plus the Python beside it for the endpoint and the credential. Not a second Shopify GraphQL Recipe: the connection shape and the 200-on-refusal are stated as shared. What differs is the money, which arrives four ways in one object -- amount as a Float and fractionalAmount as an Int, with fractionDigits saying where the point goes -- so the field with the obvious name is the one that cannot be added up |
| ~~VTEX~~ | Shipped, written against the OpenAPI description VTEX publishes for its Orders API. One purchase is not one order: when several sellers fulfil a basket each gets its own order ID and they share an order group ID, so code that stores "the order id" stores one seller's share. Also: the listing stopped returning line items in 2018 and the schema still lists them, the same total is totalValue in a row and value on the order, an order can carry a named roundingError, and the API answers two different error envelopes |
| ~~commercetools~~ | Shipped, written against the RAML reference commercetools publishes and its own error-code table. The headline is the contract a test suite cannot test: every resource carries a version, every update is a POST of {version, actions}, and a write quoting the wrong one is refused with the right one in the reply. Code that ignores it passes everything and silently overwrites in production. Also: a product holds two complete copies of itself, and centAmount is the smallest indivisible unit, which for JPY is the yen |
| ~~Shopware~~ | Shipped, written against Shopware's own source rather than a description of it. One provider disagreeing with itself: /store-api/product and /store-api/product-listing/{categoryId} are both listings of products in one API and share no paging rule at all -- page against p, 100 against 24, a 400 against a silent trim, the page length against the real count. The default total is the size of the page you are holding, because counting costs a query the framework will not run uninvited |
| ~~Printify~~ | Shipped, written against Printify's own OpenAPI description. Not the wholesale-retail split, which is Printful's and is said to be repeated: what differs is that an order has no cost field at all, so the margin is one field minus the sum of a nested array. The page size is capped at ten and answered quietly |
| ~~Printful~~ | Shipped. Print-on-demand was not on this list at all, and neither is Printify, Wix, Squarespace, PrestaShop, Faire, Ecwid or Shopware -- a category gap rather than a provider gap. One order carries costs and retail_costs, wholesale and retail, with identical field names inside both and sometimes two currencies |
| ~~Lightspeed~~ | Shipped, R-Series. The row said two unrelated products and there are four: R-Series is retail, K-Series is restaurants, X-Series is the former Vend and eCom is the former SEOshop, each with its own API and credential. The Recipe says which one it is in its first line, and detection maps the R-Series client alone |
| ~~Clover~~ | Shipped. An order comes back with no items in it unless the request asked to expand them, and forgetting is answered rather than refused |
| ~~Recharge~~ | Shipped. Two sources of truth for one order: a charge and the Shopify order it produced carry different ids |
| ~~Lemon Squeezy~~ | Shipped, and see the row below — the merchant-of-record half is FastSpring's and is repeated rather than new. What is new: licence keys, store-scoped identifiers, and totals in two currencies at once |
| ~~Gumroad~~ | Shipped, written against Gumroad's own source, which is open. Not the sales half: verifying a licence spends one. POST /v2/licenses/verify defaults increment_uses_count to true, so an app that checks its licence on every launch burns a use every launch -- and the API ships PUT /decrement_uses_count to undo it. Also: not-found is a 200 and the code's own TODO says it should not be, there are three failure shapes in one v2 API, and a missing product_id answers 500 on purpose |
| ~~Polar~~ | Shipped, benefits rather than subscriptions -- the merchant-of-record half is FastSpring's and is not repeated. Paying for a thing and receiving it are two different records, and the second can fail on its own: granting a benefit means calling Discord or GitHub for the customer, so a BenefitGrant carries two booleans and an error "if the benefit grant failed with an unrecoverable error". Four states, not two, and the listing mixes them because is_granted is an opt-in filter |
| ~~Orb~~ | Shipped. An invoice is not final until the period closes |
| ~~Metronome~~ | Shipped, written against the OpenAPI document Metronome publishes at docs.metronome.com/openapi.json. The row below asked what would earn it a place and this is the answer, on the row's own terms and one step further: POST /v1/ingest carries every billable event in the system and its whole published response is "200: Success" -- no schema, no body -- while its own description says "Duplicate events are automatically detected and ignored (34-day deduplication window)", so an event that was counted and an event that was thrown away answer identically. The endpoint that could tell you is documented as not being for that: "Do not use this endpoint to check every event in your system" |
| ~~Lago~~ | Shipped, written against the OpenAPI document Lago publishes. An event whose code names no active billable metric is accepted and, in the document's own word, ignored: the request succeeds, the event is stored and returned, and no fee is ever produced. The self-hosted/cloud difference this row predicted is in the schema -- lago_id is "not guaranteed to be a UUID; on organizations using the ClickHouse events store it is a composite string". Also: the schema and its own description disagree about whether the timestamp is ISO or Unix seconds, and a field named cents holds '1234.56' |
| Zuora | Assess — very large, and the object model predates REST conventions |
| Maxio | Assess — the former Chargify, beside Chargebee and Recurly |

## Paging parameter names, closed

Offset and page numbering are implemented, and a route declaring either must
now name the provider's own parameters. The validator refuses one that does
not, so the gap cannot reopen.

146 routes declared a style with no names. That was harmless while nothing
read the style and became a claim the moment something did, because without
names the runtime reads `limit` and the style's own word — right for some
providers, wrong for plenty, and wrong in the way that passes: the page size
is ignored, one full page comes back, and the caller's paging loop runs once.

Rather than fill in sixty-one providers' names from memory, which is the
guessing this project refuses to do anywhere else, the names were read from
the providers' own OpenAPI descriptions with `cauldron check --paging`.
Verified and declared:

| Provider | Page size | Position |
|---|---|---|
| DigitalOcean | `per_page` | `page` |
| Twilio | `PageSize` | `Page` |
| Box | `limit` | `offset` |
| PagerDuty | `limit` | `offset` |
| GitLab | `per_page` | `page` |
| RingCentral | `perPage` | `page` |
| Gmail, Google Calendar | `maxResults` | `pageToken` |
| Jira | `maxResults` | `nextPageToken` |

The other 137 routes had their style withdrawn. They now page the way they
did before any of this existed, which is honest: a word nothing read was
never verified, and a claim nobody checked is worse than silence. Restoring
one is a matter of reading that provider's description and declaring its
names, and the validator will insist.

### A description of requests is not a description

Two providers were assessed this round and both were declined for the same
reason, which is worth naming so the next attempt does not spend the hour
again: **they publish a machine-readable description of their requests and
none of their responses.**

Katana publishes a real OpenAPI 3.0 document at
`api.katanamrp.com/v1/openapi.json` -- a hundred and twenty-two paths, request
bodies fully typed. Every response in it reads

```json
"200": {"description": "Return value of InventoryController.getInventories"}
```

with no schema at all. It is a NestJS gateway spec, and the generator wrote
down what the controllers accept and nothing about what they return.

USPS publishes a Postman collection in its own GitHub organisation, at
`USPS/api-examples`, covering ninety-odd requests across twenty-five folders --
addresses, prices, labels, tracking, the lot. Every single one has zero saved
example responses. The rest of the repository is a PDF, a `.docx` and an
`.xlsx` of error codes.

A Recipe built on either would be a Recipe whose every field name was invented,
which is the guessing this project exists to refuse. Both are worth revisiting
if a response-bearing description appears; neither is worth writing from what
is published today.

The distinction matters more than it sounds. A request schema tells you what a
client may send, and a client already knows that -- it is the thing writing the
request. What a client cannot know, and what an emulator exists to supply, is
what comes back.

### What is left, and why it is stuck

`verify` reports the remainder on every run: **15 routes across 7 Recipes page
by a parameter nobody named**. The rule is a declared page size with neither a
style nor a parameter name beside it, which is what `GuessedPagination`
counts.

They are:

| Recipe | Routes |
|---|---|
| AssemblyAI | `GET /v2/transcript` |
| Bill.com | `POST /api/v2/List/Vendor.json`, `POST /api/v2/List/Invoice.json` |
| Deel | `GET /rest/v2/contracts`, `GET /rest/v2/invoice-adjustments` |
| Mercury | `GET /api/v1/accounts`, `GET /api/v1/account/{accountId}/transactions` |
| QuickBooks | `GET /v3/company/{realm}/query`, `/invoices`, `/accounts` |
| Segment | `GET /v1beta/workspaces/{workspace}/sources`, `/destinations`, `/tracking-plans` |
| Vonage | `GET /account/numbers`, `GET /search/message` |

The blocker is obtaining the descriptions rather than declaring the names. The
method above wants a provider's own OpenAPI document to read the parameters
out of, and none of these seven publishes one at any obvious location -- the
repository paths and well-known spec URLs were tried and answered 404. The
nearest thing found was an AssemblyAI spec on a personal account with no
confirmable relationship to AssemblyAI, which is not evidence: reading
parameter names out of a stranger's fork and declaring them verified would be
the guessing this section exists to refuse, wearing the costume of the method
that replaced it.

So the number stays at 15 on purpose. Anybody holding a real description for
one of these seven can run `cauldron check --paging <recipe> <spec>` and fill
it in; the work is small and the evidence is the hard part.

QuickBooks may never move. Its listing endpoint takes a query language rather
than parameters, and `STARTPOSITION` and `MAXRESULTS` are clauses inside the
statement rather than things in the query string, so there may be no HTTP
parameter to name.

Two real bugs fell out of the verification. DigitalOcean was ignoring
`per_page` and its own conformance case had been asserting the emulator's
parameter name rather than the provider's. Twilio capitalises its parameters
and both spellings were being ignored.

### A note that recorded a failed attempt was worth retrying

Grepping for the rest of these found four more, and three are principled
refusals rather than tasks: observing Docker Hub's rate limit means
exhausting a pull budget shared with everyone on this address, and observing
GitHub's or Docker Hub's rejected-credential cases means sending a
credential-shaped header at somebody else's authentication endpoint. Those
stay unverified on purpose and say so.

The fourth was worth another go. Docker Hub's rate-limit note recorded an
attempt from 2026-08-22: an anonymous request answered with
`docker-ratelimit-source` and no count at all, so the note concluded that the
`100;w=21600` shape still rested on Docker's documentation.

It did, and the documentation was stale. A HEAD against
`registry-1.docker.io/v2/library/hello-world/manifests/latest` -- which does
not count as a pull -- answered on 2026-08-23 with all four headers reading
**`100;w=3600`**. An hour, not six.

Two things came out of it:

- The window is corrected. A client backing off for the window it read was
  waiting six times too long, which is the direction that looks like the
  emulator being careful rather than wrong.
- Docker sends `ratelimit-limit` **and** `x-ratelimit-limit`, both spellings
  of both headers on the same response. Only the un-prefixed pair was
  modelled, so a client reading the `x-` form found nothing.

The 429 itself is still unobserved and the case still carries no date. What
is unobserved now is narrower: a `remaining` of zero is what a limited
address reports, and this address is not limited.

### Two notes that said what would settle them, settled

Three WordPress cases carried no `verified:` date and, instead of one, a note
saying exactly why. Two of those notes named what was missing rather than
merely recording an absence, and both turned out to be findable:

- "Settling this needs a public site on some other timezone." wptavern.com is
  one. Post 185079 carries `date 2025-01-08T22:48:02` and `date_gmt
  2025-01-09T03:48:02` -- five hours apart and on different days, so reading
  the first as UTC files the post under the wrong date. ma.tt and
  techcrunch.com diverge the same way, so it is ordinary rather than a quirk.
- "A site with no child pages cannot answer this." ma.tt has five. Page 2545
  carries `parent: 2536` and page 2536 carries `parent: 0`, which settles
  both halves: the id is an integer pointing at another page, and zero is
  what a top-level page says rather than the field being absent.

The third is still open and now for a better-established reason. Its sticky
half needs a pinned post, and wordpress.org/news, wptavern.com, ma.tt and
techcrunch.com were each asked for `?sticky=true`: every one answered with
nothing. A pinned post is a setting most publications do not use.

Live-verified cases 56 to 58.

The general point is worth keeping. A note saying "this was not checked"
records an absence; a note saying "checking this needs a site on another
timezone" is a task somebody can pick up, and two of these were done in the
time it took to read them.

### An assertion that cannot fail is not an assertion

Mapping every field of the format to whether any test names it -- by its YAML
tag or its Go field -- found twenty-seven that Recipes use and nothing tests.
Most are shapes, and a wrong shape shows up as a failing case somewhere.

Six are not shapes. They are the ways a case *asserts* something:
`header_matches`, `absent_headers`, `body_matches`, `matches_header`,
`absent_events`, `signature_header`. A shape that is quietly broken fails
loudly; an assertion that is quietly broken passes, and every case built on
it passes with it. Nineteen cases use `signature_header` and eight use
`header_matches`, and the only evidence any of them did anything was that a
Recipe somewhere used them and stayed green -- which is what a silently
ignored field looks like from outside.

There is a test now that gives each kind a response it should refuse and one
it should accept: status, headers, header_matches, absent_headers, body,
matches, absent, body_matches, no_body. Neutering any of the three untested
ones makes it fail.

All nine work. The one that failed first was the test rather than the
mechanism: `body_matches` applies to the raw body, quotes and all, which
npm's own case gets right by writing the opening quote into its pattern and
which this test had wrong.

### The behaviour under five wrong cases, stated once

Five conformance cases had been found asking for a page with `limit` against
providers that do not take it -- Help Scout, Mailchimp, SES, Zoom and Docker
Hub -- and each was corrected in its own Recipe. What none of them said is
why the mistake was so easy: until a Recipe named its parameter, the runtime
read `limit` from everybody, so the wrong word worked and the right one was
never tried.

That behaviour has three branches and no test named for it. It has one now,
in the runtime rather than in five providers that happen to exercise it:

- a named parameter is the only one read, and every other spelling is inert
- `"-"` refuses all of them, including the default
- naming nothing reads `limit`, which is a guess for the providers that call
  it something else

The two mutations that matter both fail it. Making the named branch fall back
to `limit` fails the first; making `"-"` read `limit` fails the second.

A per-Recipe case for each of the twenty-odd Recipes that name a parameter
would have said the same thing twenty times. Rollbar and Pipedrive keep
theirs, because those two are also claims about what the provider does rather
than only about what the runtime does.

### The rest of the verified set holds up

Fifty-six cases carry a `verified:` date and three of them turned out to be
claiming more than they check: Docker Hub's asserted the right things about
the wrong request, GitLab's asserted them under the wrong name, and
Discourse's left half its name unasserted. Reading the other fifty-three
found nothing further.

That is worth recording rather than leaving as an absence. The three that
were wrong shared a shape -- the assertions held under both the claim and its
opposite -- and the rest do not: GitHub's page-two case names the issue it
expects rather than merely counting one, WordPress's names the slug, and the
two Link-header cases assert the header rather than inferring it from the
body.

An earlier pass over this set looked emptier than it was, because the dump
behind it did not print header assertions. The two GitHub cases that appeared
to assert nothing assert `Link`.

### One name whose second half nothing checked

Pipedrive's "the position is start, not offset or page" sends `start` and
asserts the second deal comes back, which shows `start` works and says
nothing about the other two words. That half is the one a client gets wrong,
because every neighbouring API spells it `offset` or `page`.

There is a case for it now: `offset=1` answers the first deal, unmoved.
Declaring `cursor_param` is what makes the other spellings inert, and
pointing it at `offset` fails both cases.

A scan for the same shape -- a name drawing a contrast with a single
assertion behind it -- returned seventy-eight, and almost all of them are
sound: "a bad api key is a 403, not a 401" is a contrast whose whole content
is the status, and the status is asserted. The heuristic counted body
assertions and not the status line, which is the sort of thing that makes a
scan look alarming and mean little.

### The same scan, widened past the page size

A response field whose value the request decides, declared as a fixed one, is
a class rather than an incident. Widening the scan past `limit` and `page` to
positions, counts and page totals found four more:

| Recipe | field | said | is now |
|---|---|---|---|
| Box | `offset` | 0 | the offset asked for |
| Contentful | `skip` | 0 | the skip asked for |
| Typeform | `page_count` | 1 | computed from the total |
| Zoom | `page_count` | 1 | computed from the total |

The two `page_count` constants are the worse pair. A client looping while
`page <= page_count` stopped after the first page however many there were,
and the pages it never asked for looked like an empty collection rather than
a mistake.

Both offset echoes needed `first_page: 0` beside them, which is the sort of
thing that would have been found in production rather than here: a page
number starts at one and an offset does not, so the echo would have reported
`offset: 1` before anything had been skipped.

Every one is asserted by a case that asks for something other than the
default -- an offset of one rather than none, a page size that makes more
than one page -- because an echo tested at the default is indistinguishable
from the constant it replaced. Restoring the four fails four cases.

**Miro is left, and is a different gap.** Its list envelope declares
`size: 0`, and Miro's `size` is how many items came back in this page --
neither the total, which `count_field` covers, nor the page size that was
asked for, which `limit_field` covers. There is no mechanism for the length
of the page actually served, so the constant stays and the Recipe cannot yet
say what it means. Its `offset: 0` is a constant too, on a cursor-style
listing where no offset is meaningful.

### Auditing the verified set found the bug in the audit's own case

Two changes running had turned up a case claiming more than it checked, so
the rest of the fifty-six were read the same way. The first thing the reading
found was one written in the change before: Discourse's "a topic listing
takes per_page and reports what it gave" asserted only that the second topic
was absent. The half about reporting was in the name and nowhere else.

It was reporting thirty. `topic_list.per_page` was a route constant, so
asking for one topic returned one and said thirty -- from the field whose
entire purpose is to say how big the page was.

`responses.list.limit_field` exists for exactly this, and its own comment
describes the same bug found in Algolia: "A constant cannot do this job, and
putting one there is worse than leaving the field out... a client that asked
for page 3 was told it was looking at page 0, by a field whose entire purpose
is to say where you are."

Scanning for the pattern found **five more**, every one a constant where an
echo belongs:

| Recipe | field | said | should echo |
|---|---|---|---|
| Box | `limit` | 100 | the size asked for |
| Contentful | `limit` | 100 | the size asked for |
| Help Scout | `page.size`, `page.number` | 25, 1 | both halves |
| Typesense | `page` | 1 | the page asked for |
| Zoom | `page_size` | 30 | the size asked for |

All six are echoes now, each asserted by a case that asks for a page smaller
than the default. Restoring the constants fails five cases across five
Recipes.

Zoom turned up a second thing on the way. Its listings declared a cursor
style and no `limit_param`, so they read `limit`, which Zoom does not take --
and its paging case, the one named for pinning `next_page_token`, was paging
by that word. That is the fifth case found doing this, after Help Scout,
Mailchimp, SES and Docker Hub, and the constant echo is why it stayed
invisible: a wrong page size and a page size that never changes look
identical.

Three stranded placeholder comments in Contentful were removed as well, left
behind when its paging was declared two changes ago.

### A case name can claim more than the case checks

GitLab had one called "the token header is PRIVATE-TOKEN, not Authorization".
It sends a bad token in the Authorization header and asserts a 401 whose body
is `{"message":"401 Unauthorized"}`. Every assertion is right and none of them
is about the name: the case never tests that a *valid* token in that header
is refused, which is the exclusivity the name claims.

GitLab's own authentication page contradicts it. A personal access token may
be sent as `PRIVATE-TOKEN` or as an OAuth2-style bearer, so "not
Authorization" is very likely false, and an anonymous probe cannot settle it
either way -- a bad token answers 401 in both headers, which is exactly why.

It is renamed to what it checks, and marked verified for that: the refusal
shape, observed against gitlab.com on 2026-08-23. The status appears twice in
that body and neither copy is a number, so code switching on a numeric code
finds none.

That makes two of these in consecutive changes -- Docker Hub's verified case
asserting the right things about the wrong request, and this one asserting
the right things under the wrong name. A `verified:` date says somebody
watched the provider. It does not say the case was asking the question its
title claims.

**GitLab's visibility case is deliberately not marked verified.** gitlab.com's
project listing reads without an account and everything in it is public, so
an anonymous request confirms that `visibility` is a real field carrying
`"public"`. The case turns on `internal` -- the value with no GitHub
equivalent, and the one worth checking. Seeing the field is not seeing the
claim.

### OpenRouter's links.next is declared

The follow-up from the last change. `cursor_field: links.next` with
`cursor_url: path`, so the catalogue's next page is a path a client follows
rather than an opaque token it could not use. The validator required a case
asserting that name where the value exists before it would accept it, which
is the rule working as intended.

Its query parameters come out in a different order from OpenRouter's --
`?limit=2&offset=2` against `?offset=2&limit=2` -- and the case does not
assert the order, because that is nothing to a client following a URL rather
than parsing one.

Stated gap: OpenRouter sends `links` on every response with `next` null on
the last page, and this omits the object entirely, because the runtime writes
a cursor field only when there is a next page. Both are falsy, so a loop
terminating on either is right and code distinguishing them is not.

### Four listings that said nothing about paging, asked

The providers this collection can reach without a credential are worth more
than their number suggests, because a question about them can be settled
rather than reasoned about. Four of their listings declared no paging at all,
and all four were wrong to:

| listing | observed on 2026-08-23 | what it declared |
|---|---|---|
| GitHub labels | `per_page=2` answers two with a `rel="next"` Link header | nothing -- paged at ten reading `limit`, which GitHub does not take |
| OpenRouter models | all 422 with `links.next` null; `?limit=2` answers two with `links.next` `/api/v1/models?offset=2&limit=2` | nothing |
| Discourse `/latest.json` | thirty topics, `topic_list.per_page` thirty; `?per_page=5` answers five and reports five | nothing |
| Discourse `/categories.json` | twelve categories, and `category_list` carries no paging field at all | nothing |

Three of the four are now exercised by a case that fails without the
declaration. **OpenRouter's is not**, and the reason is worth keeping: its
size parameter really is called `limit`, so the fallback happened to read the
right word and only the style was wrong. Nothing distinguishes offset paging
from the default here because the Recipe does not model `links` -- the field
that would show the difference is `links.next`, and it is not declared.

That is the next thing to do for OpenRouter, and it is a response-shape gap
rather than a paging one.

GitHub's labels fixture had one label, which cannot distinguish "pages by
per_page" from its opposite. It has two now, for the same reason Typeform's
webhooks fixture gained a second.

### A live-verified case can still carry a wrong claim

Docker Hub, and it is the fourth of these -- but the first with a `verified`
date on it, which is what makes it worth writing down.

Its paging case sent `limit=2` and had been checked against hub.docker.com.
Both of its assertions were true: four tags matched whatever the page size,
and a next link comes back regardless. So the live check confirmed the two
things the case asserted and said nothing about the request it made.

Docker Hub does not accept `limit`. Checked again on 2026-08-23 against
`library/registry`, which has seventy tags:

| request | returned |
|---|---|
| `?page_size=2` | 2, next reading `?page=2&page_size=2` |
| `?limit=2` | **10** -- the default -- next reading `?limit=2&page=2&page_size=10` |

The second row is the part worth having. An unknown parameter is not dropped,
it is copied into the next URL, so a client sending `limit=2` gets ten
results and a next link with `limit=2` in it. The parameter that was ignored
is right there in the response, reading exactly as though it had been
honoured.

Six conformance cases across the collection send a paging parameter to a
listing that declares none, relying on the runtime's fallback of ten read as
`limit`. Docker Hub was the one that could be settled without credentials,
and it was wrong. The other five -- Column, Increase, Orb, PostHog and Svix
-- are providers whose listings do take `limit`, so the default is probably
right for them, and "probably" is the word that matters: nothing is watching.

### A third case paging by a word its provider ignores

SES. And this one is the sharpest of the three, because the case exists
precisely to pin the name of the paging field -- "a name a Recipe chooses is
only a claim if something asserts it where the value exists" -- and it was
itself asking for a page with `limit`, which SES does not take. It sends
`PageSize` now.

That makes three: Help Scout, Mailchimp, SES. Every one passed, and every one
was receiving a whole collection where it believed it had received a page.
The pattern is not that these cases were written carelessly; it is that
`limit` worked for everybody until somebody said otherwise, so writing the
wrong word had no consequence a case could show.

Typeform needed a per-route answer rather than a per-recipe one, which is why
it was left last time: its forms listing counts pages, its responses listing
walks a token, and its webhooks listing is not paged at all. Pusher's channel
listing is not paged either.

Both of those "not paged" claims needed something to be wrong about. Pusher's
fixture had three channels and could carry a case straight away; Typeform's
had one webhook, where every page size returns the same thing and an
assertion holds whatever the Recipe says. It has two now.

### Five where the answer is that there is no parameter

Salesforce, Sanity, DynamoDB, ClickUp and Rollbar. Twelve routes, and in
every one the interesting part is what the provider does *not* take:

- **Salesforce** and **Sanity** put paging inside the query language.
  Salesforce's LIMIT and OFFSET are part of the SOQL a caller writes;
  Sanity's is a GROQ slice, `[0...10]`. Neither accepts a parameter beside
  the query, so a client sending one is sending nothing.
- **ClickUp** and **Rollbar** page by number and accept no size at all. The
  page is theirs, not the caller's, and a client asking for ten gets a
  hundred or twenty.
- **DynamoDB** takes `Limit` and `ExclusiveStartKey`, capitalised, **in the
  body**. Read from the query string they are not there at all.

`limit_param: "-"` says this, and it is a claim worth a case rather than a
comment: Rollbar's asks for one project, receives both, and fails the moment
the parameter is made live again. An emulator that honoured `limit` would
hand back exactly the page that was asked for, with a next link a real
Rollbar would never send.

Salesforce took one correction on the way in. Setting `cursor_param: "-"`
alongside refused the next page as well as the size, and Salesforce does
send a position -- `nextRecordsUrl`, a whole URL to follow rather than a
parameter to compose. Only the size is absent.

### Eight more, and two cases that were paging by a word nobody accepts

Auth0, Contentful, Freshdesk, Ghost, Help Scout, Mailchimp, Mux and
Statuspage are named now. Twenty-three routes, and none needed a description
fetched: `per_page`, `skip`, `count` and `size` are the sort of thing a
provider's paging is known for.

Naming them broke two conformance cases, which is the whole argument for
naming them. Help Scout's paging case sent `limit: 2` and Mailchimp's sent
`limit: 1`, and both passed -- because nothing named the parameter, the
runtime read `limit` from everybody, and neither provider accepts it. Each
case was asking for a page, receiving the entire collection, and asserting
against it happily. They send `size` and `count` now.

That is exactly the failure `limit_param` was written to describe, found in
the suite rather than in the field: an emulator that only understands
`limit` answers with its own default, the response has no next page in it,
and the paging loop a client carefully wrote runs once and passes.

Auth0 also gained `first_page: 0`, because the comment written for it said
its pages count from nought and the declaration did not. A comment claiming
something the Recipe does not do is worth less than no comment.

### Restoring them, one description at a time

**20 routes across 10 Recipes** still page by a parameter nobody named.

That number went up by a hundred and seven when the counter learned to see. It
had treated a declared style as though it named the parameters, on the reasoning
that neither is a word anybody writes down by accident -- and a style says how a
provider pages, not what it calls things. Front is the counterexample: three of
its listings declare cursor paging and name neither parameter, Front's own
description says they take `limit` and `page_token`, so the runtime read
`cursor`, Front ignored it, and every request came back on page one. Ninety-four
routes across forty Recipes were in exactly that state and neither counter could
see one of them -- this one required an empty style, the one below required an
empty everything. Naming one of the two is not naming them either, which is
thirteen more.

**Stripe was one of them.** Its own description names limit, "between 1 and 100,
and the default is 10", and `starting_after` for the position, with
`ending_before` for the other direction. It has no parameter called `cursor`,
and this Recipe -- the one most of the collection's vocabulary was copied from,
which is exactly why nobody looked -- declared the style, named neither, and
sent one.

Thirteen more went with it, all from descriptions those Recipes already record.
Asana's `offset` is an opaque continuation token rather than an index, so its
five listings are cursor paging under a parameter named like a page number.
Notion takes `page_size` and `start_cursor`; Lithic `page_size` and
`starting_after`; Cloudflare Stream `limit` and `after`, which is a timestamp.
Render had named its cursor and not its limit, and was right by luck -- which is
not the same as having been checked.

**A gap recorded in prose does not notice when it is closed.** Three Recipes --
Ably, Buildkite and GitHub -- carried a sentence saying "Cauldron does not send
Link headers", and all three declare `link_header` in the very same file. That
was true when each was written, Link was modelled afterwards, the envelopes were
updated and the sentences were not. A reader was being told a limitation that
had been fixed, in a file whose whole value is that its prose can be trusted.
Five sentences across three Recipes, corrected. Ably also names the parameter
inside its link now -- the word there is this fake's own, because a client
follows a Link URL as given rather than building one, and saying so is what
stops the runtime supplying the same word silently.

Dropbox writes both names down now, and records the one thing this format
cannot do about it: Dropbox does not take a cursor on `list_folder` at all. It
takes one at `/2/files/list_folder/continue`, a separate path whose only
argument is the cursor, and sending one to `list_folder` is an error there. The
Recipe accepts it on the listing because a listing is addressed by one route
here, and a fake that cannot serve a second page would be the worse of the two
inaccuracies -- but the file says so rather than leaving a reader to find out.

The US Census catalogue is the sharpest of the settled ones, because it does not
merely lack a page size -- it refuses one. Struck live on 2026-09-05: the
variables listing answers thirteen and a half megabytes to a bare request, and
`?limit=2` is a **400**, "error: the 'get' argument must be a comma separated
list of variable names", which is what any unrecognised parameter gets. Most
providers ignore a parameter they do not have; this one turns the request away,
so a client trying to page the catalogue fails loudly rather than quietly. That
listing is `style: none` now, which also removes the ceiling the Recipe had been
carrying -- a number about this fake rather than about the Census.

QuickBooks is a different shape of unanswerable and worth naming separately.
Its listing endpoint is `/query` and its only parameter is `query`, whose value
is a statement rather than a filter -- so whatever paging it offers is written
inside another parameter's value, which is not a place this format can address.
Both names now say the provider accepts none, which is true of the query string
and is the honest half of the picture; before, the runtime was reading `limit`
and `cursor` off a request that carries neither.

Snyk's three were already read -- its v1 API "answers these with everything it
has", so the size is fixed and no name is accepted -- and only the position was
left blank, which had the runtime reading `cursor` off a request that carries no
query string. Both names now say so.

**And the ones that could not be reached now say what was tried.** Brex
publishes a description per API behind its documentation site, and every address
that returns anything returns the site's own HTML; the JSON forms 404 and so
does the pagination guide the site links to. Recurly's v3 reference is a
single-page application that serves its navigation and not its content to
anything that is not a browser, and the description its client libraries are
generated from is not at any address those repositories suggest. Both Recipes
now carry the date, the addresses tried and the plain statement that their page
size is a number with no name attached -- a guess, labelled as one, rather than
a declaration nobody can check.

Duffel had named its cursor and not its size, and the size happened to be the
word the runtime supplies -- right by coincidence, which is not the same as
checked. Its reference gives both numbers as well: "default 50, range 1-200".

**Segment is the one that could not be settled, and why is the finding.** Its
current Public API documents `pagination.count` and `pagination.cursor`, with a
response carrying `next`, `current`, `previous` and `totalEntries` -- and two of
those four response names are what this Recipe already declares. Its paths are
not that API's. They are `/v1beta/workspaces/...`, the Config API, which
Segment's current documentation does not describe. The envelope in that file
came from one version and the routes from another, so declaring the current
API's parameter names would deepen the mix rather than settle it. Recorded in
the Recipe, left unstated.

Paddle and Cloudinary are the same again, and Paddle is the sharper of the two.
Its page size is `per_page` and its position is `after`, and its own reference
says to "use the `next` URL directly" because that URL already carries the
filters and sort the request was made with -- so `meta.pagination.next` is a
whole URL, and this Recipe served the bare cursor there. Two of its cases sent
`?limit=`, which Paddle does not have. Cloudinary takes `max_results`, ten by
default, and `next_cursor`, the same word the response carries it back under.

Webflow, Lob and Ramp finish the ones whose references could be read. Webflow's
two site-level listings declare no query parameters at all, so both names are
`"-"`; its items listing takes `limit` and `offset`, "default 100, maximum 100".
Lob takes `limit`, "1-100, default 10", and `after`, with `before` for the other
direction and only one allowed in a query at a time. Ramp takes `page_size`,
"between 2 and 100 ... the default value 20", and `start`, "the ID of the last
entity of the previous page" -- its description is published behind its
documentation site and is recorded and fingerprinted now.

Lob's next page is a URL too, under `next_url` with `previous_url` beside it,
and this Recipe serves neither. That one is stated rather than fixed: the field
would take the URL, and what it would not take is the pair, which is how a
client knows which direction it is going.

Shopify and Zendesk both advertise their next page somewhere this Recipe was
not looking, and both take a page size called something else. Shopify takes
`limit` (fifty by default, two hundred and fifty at most) and `page_info`, and
its link relations travel in **link response headers** with no cursor in the
JSON at all. Zendesk pages two ways: the offset one this Recipe serves takes
`per_page` and `page` and answers a `next_page` **URL**, where the fake was
serving a bare number; the newer cursor one takes `page[size]` and `page[after]`
and answers under `links.next`, which is recorded rather than modelled since a
route takes one set of names here. Zendesk's own reference is the reason to
prefer the newer one: offset requests "are limited to the first 100 pages and
10,000 resources".

A pattern worth naming, since it accounts for several of these: a Recipe that
has read its provider's document, found no parameter, and written that in prose
beside a declaration that still leaves the field blank. VTEX's order-group
route, Pub/Sub's pull, DynamoDB's Query and Scan, SQS's ReceiveMessage and
Daily's presence listing were all in that state -- the comment said "there is no
cursor here" and the runtime went on reading `cursor` from a query string the
request does not have. `"-"` is how that gets said once rather than twice, and
the counter only sees the declaration.

Calendly calls its page size `count` -- "the number of rows to return", twenty
by default and a hundred at most -- and its position `page_token`. Neither is a
word the runtime would supply, and its description was reachable once the index
at `developer.calendly.com/openapi.json` was followed to the document it names.
Recorded and fingerprinted. WorkOS takes `limit` and `after`, the `after` being
the `list_metadata.after` its Recipe already served with no way for a caller to
send it back.

Airtable, HubSpot and Okta came from their own reference pages rather than a
machine-readable document. Airtable takes `pageSize` and `offset` -- and its
`offset` is a continuation token rather than a number, despite the name, which
is why it is cursor paging. HubSpot takes `limit` and `after`, the `after` being
the `paging.next.after` its Recipe already served. Okta takes `limit` and
`after` too, and advertises the next page **in the Link header**: its own guide
says "Pagination links are included in the Link header of responses", and this
Recipe served no pointer of any kind, so a client following Okta's documented
mechanism found nothing to follow. Airtable's case was the fifth today sending a
page size its provider does not have.

DynamoDB was carrying Query's parameter name on ListTables. AWS's own service
model gives ListTables exactly two inputs -- `ExclusiveStartTableName` and
`Limit` -- while `ExclusiveStartKey` belongs to Query and Scan, which resume
from a key map rather than a name. A caller resuming a table listing was sending
a field ListTables does not read. Query and Scan themselves had a careful note
saying no cursor can be declared for them, and now say it in the declaration:
`cursor_param: "-"`, rather than leaving the runtime to read a query parameter
off a request that carries no query string. SQS's ReceiveMessage is the same
shape and got the same treatment -- AWS's model gives it no continuation input
of any kind, unlike the ListQueues beside it.

Sentry names its page size two ways on one API: `limit` on the issues listing
and `per_page` on the organisation's projects and releases, with the same
`cursor` on all three. The runtime read `limit` everywhere, which was right once
and wrong twice. Its description is recorded and fingerprinted.

Google Pub/Sub takes `pageSize` and `pageToken`, from the discovery document
Google publishes -- cited in the Recipe rather than recorded, because that
document is not OpenAPI and `cauldron drift` could not read it. Its pull route
was already carrying a careful note saying there is no cursor there at all, and
now says so in the declaration too: `cursor_param: "-"`, rather than leaving the
runtime to read a query parameter on a request that carries no query string.

Klaviyo, CircleCI and Miro were the same hunt and the sharpest three. Klaviyo
takes `page[size]` and `page[cursor]`, square brackets included -- JSON:API's
spelling, and not a parameter anybody arrives at by guessing. CircleCI takes
`page-token` and **no page size at all**, so `limit_param` there is `"-"`; its
`next_page_token` is typed "string, nullable", meaning a client testing whether
the key is present finds it on every page and only the value says whether there
is another. Miro pages its boards and members by `limit` and `offset` and its
items by `limit` and `cursor` -- two styles under one Recipe, from one document.

Three more providers publish a description in their own GitHub organisation
and nobody had recorded it: Square, Slack and Discord, all three now recorded
and fingerprinted. Square and Slack take `limit` and `cursor`, which are the two
words the runtime supplies -- so those six listings were right and had never
been checked, which is a different thing. Discord takes `limit` and `after`, a
snowflake rather than a token, and its guild-channels listing declares no query
parameters at all: however many channels a guild has, they come back in one
response.

Six more from the same source. Dub's is the one to remember: it declares
`startingAfter` and `endingBefore` as its cursors and `pageSize` as the size,
and it declares `page` too -- whose entire description is "DEPRECATED. Use
`startingAfter` instead." The parameter that reads like the obvious one is the
retired one, and a Recipe filled in from the names alone would have picked it.
Intercom pages its contacts and conversations by `starting_after` and its
companies by `page`, on the same API, with nothing in the response saying which.
Daily had named its cursor and not its size, and the size happened to be the
word the runtime supplies -- right by coincidence rather than by checking -- and
its presence listing takes no position parameter at all, which is `"-"`.

Notion's own case is the whole failure in one place. It sent `?limit=1`, a
parameter Notion does not have, against a Recipe that was reading `limit` too --
so the two wrongs agreed, a full page came back, and the case asserted a page it
never asked for. It sends `page_size` now. Intercom's case had the same
shape and now sends `per_page`; Klaviyo's sent `limit` where Klaviyo wants
`page[size]`; and CircleCI's sent `limit` to an API that has no page size at
all, which is four.

Vercel's three are settled the same way and are the sharper example: its
deployments and domains listings take `until`, its projects listing takes
`from`, and this Recipe was sending `cursor` to all three. Two names for the
position on one API, neither of them the word the runtime supplies. Its
`pagination.next` is a timestamp number as well, where the engine has only the
last record's identifier to offer -- so that one is stated in the file rather
than modelled, and the case that pins it now says whose value it is pinning.

Front's own three are settled by the description that found them, and its
`_pagination.next` turns out to be a whole URL as well -- the field's own
example in that document is
`https://yourCompany.api.frontapp.com/message_templates?page_token=...`. Its
case had pinned that field as `^cnv_[0-9]+$`, the bare cursor id this Recipe was
serving, so the case described the fake rather than the provider and passed.

Twenty-five providers have been settled. Eight were read from that provider's
own description or reference, none of them guessable:

| Provider | What the description said |
|---|---|
| SendGrid | Two listings paging two different ways. Templates declares `generations`, `page_size` and `page_token` and no `limit` at all; bounces declares `limit` and `offset` |
| Chargebee | `limit` and `offset`, and `offset` is typed **string** with the description saying to set it to the `next_offset` the last response returned. A cursor wearing an offset's name |
| OneSignal | `limit` and `offset` on both listings |
| Webflow | Eleven query parameters on the items listing including `limit` and `offset`, and none whatsoever on the two site-level listings, which are now `limit_param: "-"` |
| Pipedrive | `start` and `limit` in v1. Not `cursor` and `limit`, which is v2 -- and v2's `/deals`, `/persons` and `/organizations` live under `/api/v2` while v1's description no longer declares them at all. A Recipe modelling v1 paths pages the v1 way |
| AssemblyAI | `limit`, defaulting to ten and capped at two hundred, and a position that is another transcript's id: `before_id` for older, `after_id` for newer. `before_id` is the continuation, because the listing is newest first -- and AssemblyAI's own `page_details` calls that one `prev_url` while `next_url` points at transcripts newer than the ones you have. A loop following `next_url` walks towards the present and stops on its first page |
| Mercury | Two listings on one API paging two different ways, both from the OpenAPI document published with the reference. Accounts is cursor-based on `start_after`, with `end_before` for the other direction and the two mutually exclusive; account transactions is `limit` and `offset`. Nothing in either response says which one you are on. And both default to **a thousand**, where this Recipe had guessed twenty-five and fifty |
| Vonage | `size` and `index` on the Numbers API, `index` being a page number starting at one rather than an offset. `size` defaults to ten and caps at a hundred -- and the hundred was what this Recipe had written down as the default, so every page it served was ten times the real one. The SMS Search listing beside it is still unnamed: its reference page answered 504 twice and Vonage publishes the retired API's parameters nowhere else |

Chargebee is the one to remember. Filling that in from the parameter's name --
which is exactly what filling it in from memory would do -- declares a numeric
offset against an opaque token, and it fails in the quiet direction: page two
starts from the beginning again.

**The blocker is reaching the description, not the work.** Of the 29 Recipes
that needed names, 2 had a description already fetched. Guessing repository
URLs from the provider's name failed 6 times out of 6. Searching GitHub for an
org's own OpenAPI repository found 3 of 12 tried. Third-party copies of a
provider's description are on GitHub for several of the rest, and they are not
the provider's own description -- using one is a guess with more steps, and the
Chargebee row above is what that costs.

So the remaining 60 are not waiting on effort. Each one is waiting on somebody
finding where that provider publishes its description, or on watching the
provider page a real collection.

### And the count was the smaller half of itself

**41 more listings across 21 Recipes declare no paging at all**, and the
runtime pages them anyway: a route with no page size is given ten and reads
`limit`, exactly as a route declaring a size with no name is. The report could
not see them, because the count starts from a declared page size. So the
figure read as the whole answer and was a third of it -- 60 of 168.

Nothing is truncated by it today: no fixture behind one of these holds more
than ten records, which is why it stayed invisible. That is not a reason it is
fine. The claim is about the provider and the fixture is not the provider. A
listing a Recipe describes as unpaged answers at most ten and offers a cursor,
and the first collection large enough to notice will not be one of ours.

They are counted apart rather than folded in, because they are not the same
omission: one Recipe looked at paging and did not finish, the other has not
looked.

### A hundred and four of these can simply be asked

The blocker recorded above is reaching the provider's description. For a large
part of what is left, there is a better source than a description: the provider
itself. **104 of the unstated routes sit on 64 Recipes that already carry a
live-verified 200** -- this project has reached those hosts anonymously, which
means their paging can be observed rather than read about.

That is a better answer than a document, and the ones done that way have already
paid. Chuck Norris's search **ignores a page size entirely**: `?query=kick`
answers all 821 matches and `?query=kick&limit=2` answers the same 821, while
the fake was paging it at ten -- so a client sending `limit=2` got two here and
eight hundred and twenty-one from the provider. EPSS echoes `"offset": 0,
"limit": 100, "total": 368543` straight back in the envelope this Recipe already
serves, and `?limit=2&offset=2` moves the window. Datamuse has a size and no
position at all: `?max=5` answers five, `?max=1000` answers all 467 there are,
and `?offset=5` answers the same hundred beginning with the same word, which is
how an unrecognised parameter behaves.

**And every route that names only half its paging says why too.** The twenty that
page by a parameter nobody named are in the same position as the forty-one: each
carries a dated record of the attempt. The guard test covers both categories, so
neither kind of silence can arrive without an explanation again.

100ms is the last of them and the clearest statement of the rule this sweep keeps
proving. Its Rooms and Sessions listings are settled -- `limit` and `start`, "Min
: 10, Max : 100" -- and every address for its *room-codes* listing 404s, including
the one that page's own anchor link points at. The neighbouring pair is very
likely the same pair, because 100ms uses one shape across its listings. "Very
likely" is exactly why nothing is named there: copying a sibling's parameters
onto a route is the mistake this sweep keeps finding, and doing it deliberately
would be worse than doing it by accident.

**Every listing that still declares no paging now says what was tried.** That is
the line this sweep was heading for. 41 listings across 21 Recipes remain
unstated, and there are no silent ones left in the collection: each carries a
dated record naming the addresses that 404'd, the references that render from
JavaScript, the sites that answer 403 to anything that is not a browser, the
endpoints that answered 429, and in three cases a finding larger than the paging
-- Basiq's cited spec not containing its own route, Midtrans's listing not
appearing in its provider's documentation at all, DHL confirming it paginates
without saying how.

The distinction that matters is between *unstated* and *unexamined*. Silence used
to mean both. It now means only the first, and the second has been eliminated:
every one of the original 344 listings has either a declaration or a paragraph
explaining why there is not one. Heroku is the only one where the obstacle is
this format rather than the provider -- it pages by `Range` and `Next-Range`
headers, and `pagination.in` accepts `query` and `body`.

Xendit and New Relic close the un-noted list. Xendit's `after_id` is the fifth
spelling of "the last record you saw" in this collection, after Stripe's
`starting_after`, Persona's and BambooHR's `page[after]`, and TalkJS's
`startingAfter`.

Kustomer says only one thing about its page numbering and it is enough: `page` is
typed with a **minimum of 1**. A type constraint settles what no sentence on the
page states. Its Recipe already declared the response half as `links.next`, a
whole URL, so this route now carries three spellings of the same idea -- `page`
out, a URL back, and a `pageSize` that is neither `limit` nor `per_page`.

Basiq's dated record turned up something larger than its paging. Its reference
could not be reached at any address -- api.basiq.io/reference/<operation> 404s for
every name, docs.basiq.io does not have the page, api.docs.basiq.io does not
resolve -- and the connect.json that Recipe records as its own `upstream.spec`
describes `/users`, `/jobs` and `/connections` **without containing
`/institutions` at all**. The spec a Recipe cites not containing the route the
Recipe models is a bigger question than the page size, and it is written at the
route rather than papered over.

Avalara pages by `$top` and `$skip` -- OData's spelling, dollar signs and all:
"to view page 3 with a page size of 100, you'd use $top=100&$skip=200". A leading
dollar sign is the sort of thing a client library escapes by accident, and neither
name resembles `limit` or `cursor`. The response reports `@recordsetCount` and
`@nextLink`, at-signs this time, which this Recipe does not model.

Toast states its page numbering by arithmetic rather than by declaration: "set the
pageSize parameter to 10, and you set page to 2, the API returns a set of objects
that starts with the eleventh object." Eleventh after ten means page one held the
first ten, so `first_page: 1` is declared on that sum rather than on a sentence --
worth doing carefully, because two providers in this survey number from zero.

TalkJS carries two positions on one route -- `startingAfter`, which takes the
previous page's last conversation id, and an `offsetTs` that offsets by sort
order. The id one is declared because it walks a collection without depending on
the ordering the caller chose. It is Stripe's `starting_after` with the
underscore moved, which is the fourth spelling of that idea here after Column,
Persona and BambooHR.

Nutritionix fixes its page size in prose and gives nobody a parameter for it:
fifteen query parameters on the search endpoint and not one controls the count,
because "the results from the instant endpoint are separated into 2 arrays each
with a maximum 20 food items". A ceiling with no way to ask for the next one is
not a page. Its natural-language endpoint answers one record per food parsed out
of a sentence -- "two eggs and a slice of toast" is two records -- so half that
answer is not a page, it is a meal missing an egg.

Midtrans raises a bigger question than its paging. Its Core API documents
single-transaction lookups, `GET /v2/{order_id}/status` with no query parameters
at all, and no listing at `/v2/transactions` appears anywhere its documentation
index reaches -- so the route this Recipe models may not be an endpoint Midtrans
publishes. Declaring `style: none` there would assert that an endpoint of
unconfirmed existence serves whole collections: two claims for the price of none.

Electricity Maps answers its zone list as a **keyed object**, not an array -- 350
zones under 350 codes, struck live, with `?limit=2` ignored. A "page" of a keyed
object is not a smaller object, it is a different one with keys missing, and a
client looking up "DE" in a truncated one concludes Germany is not covered. That
is a sharper failure than a short list, and the fake was serving ten of them.

Sendcloud's spec corrected its own Recipe. The header called next and previous
"cursor **URLs**"; the v2 description types both as `{"type": "integer",
"nullable": true}`, and the parameter that consumes them is called `cursor` and
described as a "token". A token that is an integer, called a token, is the sort
of thing a client stores in the wrong column once and then forever.

SerpApi documents a position and no size, and gives away the size in its own
worked example: `start` "skips the given number of results ... 0 (default) is the
first page of results, 10 is the 2nd page of results, 20 is the 3rd page". The
arithmetic steps by ten, which is the underlying engine's page and not something
a caller can set -- so the size is `"-"` and the ten beside it is inferred from
the provider's example rather than from this engine's habit, which is stated at
the declaration.

Marqeta will not settle either, and one clause survived its JavaScript: "use the
/transactions endpoint to retrieve smaller datasets (**up to one page**)", with a
pointer to a different product for anything larger. One page and no more is a
stronger claim than "it pages" and a different one from "it does not", and it is
not declarable from a single clause in a paragraph about something else.

DHL is the sharpest of the ones that will not settle, because it is the only
provider so far that *confirms* it pages without saying how. Its Unified Tracking
reference documents six parameters, every one of them narrowing which shipment,
and its own changelog records "Added pagination" at version 1.0.10 while no page
names the parameters. The "Open API Specification" the reference points at is not
served at any address the page gives. Here `style: none` would contradict DHL's
own changelog and a guessed name would be worse than silence, so the route keeps
a note recording exactly that.

Gemini is `pageSize` and `pageToken` -- fifty per page, at most a thousand, and
the walk ends on a **missing** field rather than a null one. Worth noting against
Google Calendar's `maxResults` and `pageToken`, which is a different pair on
another Google API: guessing from the vendor does not work either. Groq's model
catalogue takes no parameters at all, the fourth model listing here to settle that
way after OpenAI's, xAI's and Perplexity's -- half a model catalogue is a client
that thinks a model does not exist. Honeycomb says "all" twice in one sentence
about its dataset listing and gives it nothing to ask for fewer with.

Zilliz pages by `pageSize` and `currentPage`, which is a name this survey has met
as a *response* field -- the Guardian reports one -- and here it is the request
parameter. No default and no maximum is stated for either, so neither is declared:
the 10 in Zilliz's example is an example, not a default, and writing it down as
one would be inventing a fact.

Moodle needed its own source rather than its documentation, and the two functions
this Recipe models split cleanly. `core_course_get_courses` takes exactly one
thing -- an options structure holding `ids`, "list of course id. **If empty return
all courses** except front page course" -- so it does not page.
`core_enrol_get_enrolled_users` does, and this format cannot spell it: the size is
`limitnumber` (not `limitnum`) and the position is `limitfrom`, but `options` is a
list of `{name, value}` pairs, so on the wire a page size is
`options[0][name]=limitnumber&options[0][value]=5` -- the parameter's name is a
*value* of another parameter. Both default to 0, which Moodle passes to
`get_recordset_sql`, where zero means no limit, so the default page is the whole
set. A note at the route says all of that instead of declaring something false.

Postscript's subscriber listing takes twenty-odd filters, a sort, and exactly one
paging parameter -- `page`, "page number of results to start from" -- with no
per-page, no limit and no cursor anywhere. What its description never says is what
number the first page has. `first_page` is left at the engine's default of one,
which is what that phrasing *implies* and not what it *states*, and if Postscript
numbers from zero the way Braze and Navitia do then the walk starts one page in.
That doubt is written at the declaration rather than resolved by assertion.

Brave Search is the third provider to document `may_undershoot` -- "the actual
number of results returned may be less than count" -- after Onfleet and Sumo
Logic, and the three of them were written by different companies for different
purposes. It also bounds the walk rather than the page, like TomTom: `offset` has
a **maximum of 9**, so the entire reachable collection is ten pages of twenty.

Alpaca carries two position mechanisms on one route -- `after`/`until`
timestamps with a `direction` defaulting to desc, and a separate
`before_order_id`/`after_order_id` cursor pair. The timestamp one is declared
because it is the documented default shape, and the existence of the other is
stated rather than hidden. Its positions listing takes no parameters at all: a
client that receives ten of forty open positions is looking at the wrong
portfolio, not a smaller one.

Fastly and Drip are both one-provider-two-listings splits where the *other*
listing is the tempting one. Fastly documents `page` and `per_page` on
`/service` and nothing at all on a service's versions. Drip's subscriber listing
takes `per_page` -- "defaults to 100. Maximum 1000" -- and its accounts listing
documents its arguments as "None." Copying the sibling's paging across is exactly
the mistake an emulator makes when a Recipe says nothing, and with nothing
declared both routes looked identical to this engine.

EasyPost's position is a mutually exclusive pair: `before_id` is "only records
created before the given ID will be included. **May not be used with after_id**",
and `after_id` is its mirror. `before_id` is the one declared, because EasyPost
lists newest-first and walking backwards in time is walking forwards through the
collection -- a client reaching for `after_id` would be paging away from the
records it has not seen. This format has one position per route and no way to say
two parameters exclude each other, so the other half is in prose.

Ashby's public job board answers all 72 postings and ignores `?limit=2`, struck
live. A job board is meant to be embedded whole in somebody's careers page, and
half of one is a careers page missing jobs. That Recipe's `syncToken` cursor
belongs to Ashby's *authenticated* API: one provider, two surfaces, and only the
private one pages.

Sumo Logic is the argument for having built `may_undershoot` this morning. Its
search-messages endpoint documents the limit parameter as: "limit the number of
messages returned in the response. **The number of messages returned may be less
than the limit.**" That is Onfleet's behaviour, in a description written years
apart from Onfleet's, found by accident four hours after the key existed. Both
its parameters are also *required* -- there is no request that omits them, so
there is no default page at all, which the format has no way to say and the
declaration says in prose.

FRED holds the largest default in the collection by three orders of magnitude:
`limit` defaults to **100000**. That is right for what FRED serves -- an economic
series is meant to arrive whole -- and the fake was answering ten observations. A
chart drawn from ten observations is not a short chart, it is a wrong one.

TomTom bounds the *walk* rather than the page: `ofs` caps at 1900 because "the
total number of results can be no more than 2000". HERE is the third geocoder
here with the same shape as Nominatim and Mapbox -- a size, no position, a ranked
answer rather than a walkable collection -- defaulting to 20. And Canvas's one
unstated listing turned out to be a JWKS document, which is the whole key set by
definition: a verifier that receives half of it cannot validate tokens signed by
the other half. Canvas pages plenty of things, at ten with `per_page` and Link
headers, and a Recipe copying that onto every listing would have put it here too.

The Guardian named its own paging parameters inside an error this Recipe has been
pinning all along. `?page=99999` answers "Content API does not support paging this
far. **Please change page or page-size** or consider filtering using a date
range" -- struck live, quoted verbatim in that file's errors block. The numbers
were already there too, in live responses the Recipe asserts: `response.pageSize`
is 10 and `response.currentPage` is 1. And its `/sections` route has asserted
since before it could declare it that `response.pageSize` is **absent** there,
with the comment "sections carries no paging metadata at all, unlike search" --
which is better evidence for a listing that does not page than any document.

Heroku is the one that stays undeclared on purpose, and the reason is a real gap
in this format. It pages by HTTP *headers*: `Range: id <start>..; max=200` out,
`206` and `Next-Range` back, and no query parameter of any kind. `pagination.in`
accepts `query` and `body`; a third value would need more than a name, because
Heroku's Range value is compound -- position and page size inside one string -- so
`limit_param` and `cursor_param` would have nothing separate to point at, and the
reply side needs a response header this format has no key for. `style: none` would
be the comfortable lie: it would clear the count and tell every reader Heroku
serves whole collections, which it does not. A note at each of the four routes
records why nothing is declared instead.

BambooHR's directory is the bluntest statement of the problem in any provider's
own words: "this endpoint returns the whole directory in one response and accepts
no name, department, or field filters, so narrow the results on the client side."
A company with four thousand employees receives four thousand employees. Its
newer `/v1/employees` listing is BambooHR's own answer to that, paging by
`page[limit]` and `page[after]` -- and the cursor is base64 of a JSON object that
decodes to `{"nextEmployee":125}`, so it is opaque by convention rather than by
construction, and somebody will eventually depend on the id inside it.

Backblaze puts a price on its page size, which is a thing no key here can hold.
maxFileCount defaults to 100 and caps at 10,000 -- and B2 counts a call returning
more than a thousand files as several transactions for billing. A client that
raises the page size to walk a bucket faster is choosing to be billed ten times
per call, and nothing in the response says so. max_limit can hold the ceiling and
has nowhere to put the cost, so the cost is in prose beside it. Its bucket
listing, meanwhile, takes filters and no paging at all: one provider, two
listings, and the one that pages is the one that can grow without bound.

Navitia is the second zero-based page numbering in this survey after Braze --
start_page defaults to 0, count to 25 -- with a ceiling Navitia states in bold:
"The number of objects returned for a request can not be superior than 200."

Better Stack retires a reason rather than a claim. Its header said route-level
paging was left undeclared because "what query parameter actually requests a
given page was never confirmed against a real response ... because no credential
exists to call one with". Still true -- but its own pagination page names them
plainly: `per_page` ("the default value is 50 and maximum value is 250") and
`page` ("the page number you want to return. Starts with 1"). A documented name
read from the provider beats the engine reading `limit` and `cursor`, neither of
which Better Stack takes, and the provenance is stated where the declaration is.

balldontlie is `per_page` (25, capped at 100) and `cursor` -- with the response
field named `next_cursor` and the request parameter named `cursor`, which is the
ordinary shape and worth saying because that Recipe already declared the response
half and had never declared the request half.

Unleash is the one where paging would itself be the bug. Its client and frontend
endpoints take filters and no size or position at all, and that is the only thing
a feature-flag API can do: an SDK holding half the flags evaluates the other half
*wrong*. It would not error on a short page -- it would quietly report the missing
thirty flags as off. The fake was serving ten of them.

Recharge is the opposite lesson. Both its parameter names are exactly the ones
the runtime guesses, `limit` and `cursor`, so nothing about the request shape
changes -- and the default is **50**, rising to 250, against the ten being served.
A guess that happens to use the right name is still a guess until it is written
down, and the number behind it was wrong the whole time.

Persona spells its paging `page[size]` and `page[after]`, and documents the cursor
as an "object ID" -- not opaque, the last record's own identifier, which is
Stripe's `starting_after` under a different name. Airbrake's `page` and `limit`
default to 1 and 20; the size was right by luck and the position was not, so a
client following the emulator's next page sent a cursor Airbrake ignores and read
page one again.

Google Books is `maxResults` and `startIndex` -- "the default is 10, and the
maximum allowable value is 40", "the index of the first item is 0" -- so the size
was right by coincidence under a name Google does not accept, and the ceiling of
forty is now enforced rather than described. That one was read from Google's own
guide rather than struck live, because the API answered **429**, "Quota exceeded
for quota metric", to an anonymous request. A rate limit is not something to
retry past, and the guide was the better source anyway.

Dwolla's `limit` (default 25, maximum 200) and `offset` are declared on all three
of its listings, with the provenance stated: the customers reference is the one
that states them, and the funding-source and micro-deposit listings take them as
Dwolla's API-wide convention rather than as separately confirmed fact.
MercadoLibre's country list answers all 46 in a bare array and ignores both
`limit` and `offset`, struck live.

Twilio Verify takes three paging parameters and one of them does nothing. Beside
`PageSize` ("The default is 50, and the maximum is 1000") and `PageToken` sits
`Page`, documented in full as: "The page index. **This value is simply for client
state.**" The server does not read it. A client that pages by incrementing `Page`
-- the obvious thing to do with a parameter called `Page`, on an API whose
response also carries a page index -- sends a number nobody looks at and receives
the first page forever. That is now the third distinct way a parameter's name has
lied about its job in this sweep, after Ory's non-sequential `page` and Polygon's
`limit` that is not a page size.

Bandwidth adds a constraint no key here can hold, so it is written down: "The sum
of limit and after cannot be more than 10000" -- a ceiling on the *walk* rather
than on the page. It also has `limitTotalCount`, which caps the reported
`totalCount` at 10,000 when true and gives "an accurate totalCount" when false.
The total is either true or fast and the caller picks; MongoDB Atlas's
`includeCount` is the same idea spelled differently. Amazon's Selling Partner API
just spells everything out -- `MaxResultsPerPage` ("Value must be 1 - 100.
Default 100") and `NextToken` -- which is a hundred where the fake served ten.

Polygon's aggregates endpoint takes a parameter called `limit` that is **not a
page size**, and it is the worst thing this survey has found. Polygon's own
words: "Limits the number of base aggregates queried to create the aggregate
results. Max 50000 and Default 5000." It bounds the underlying bars Polygon reads
in order to *compute* the answer. Asking for `limit=2` on a monthly aggregation
does not return two months -- it computes months out of two base aggregates and
returns whatever that produces.

Every other mistake in this document makes a client see too few records. This one
makes it see **wrong numbers**, in a response that looks complete, with no short
page and no error to notice. A client that read the name and paged on it would
have been quietly recomputing its own data. The route declares `style: none` --
there is no position parameter at all, the window is the `from` and `to` in the
path -- and the validator's refusal to let `limit_param` sit beside `none` is
exactly the right outcome here.

MongoDB Atlas pages by `itemsPerPage` (default 100, maximum 500) and `pageNum`,
so both halves were wrong and the page was a tenth of what Atlas serves; it also
takes `includeCount`, defaulting to true, which means the total a client loops on
is something a caller can switch off. GoCardless's bank-account-data API is
another one-provider-two-listings case: requisitions take `limit` and `offset`
(defaults 100 and 0), and the transaction listing takes `date_from` and `date_to`
and nothing else -- narrowing a date range is not paging, and a client swamped by
transactions has no parameter to ask for fewer.

Ory Kratos carries **four** paging parameters on one endpoint, and they disagree
with each other. The modern pair is `page_size` (default 250, maximum 500) and
`page_token`. Beside them sit `per_page` (default 250, maximum **1000**) and
`page`, both marked "DEPRECATED: Please use `page_token` instead" -- two page
sizes on the same route that do not agree on their own ceiling.

And the deprecated `page` is not a page number. Ory says so outright: "This value
is currently an integer, but it is not sequential. The value is not the page
number, but a reference. The next page can be any number and some numbers might
return an empty list. For example, page 2 might not follow after page 1. And even
if page 3 and 5 exist, but page 4 might not exist." A client that read the name,
saw an integer, and incremented it walks the collection in an order nobody
defined and lands on empty pages that mean nothing. This is the sharpest example
yet of why a parameter's *name* is not its *semantics*, and why guessing from the
shape of a query string is not a substitute for reading.

The page sizes keep being nothing like ten: Kratos serves **250**, InfluxDB 20
with a ceiling of 100, Qdrant's points query 10 (right by luck, in the body of a
POST, with `offset` being read as a cursor). Grafana's `/search` documents the
part that usually has to be inferred -- "Numbering starts at 1. limit param acts
as page size" -- and caps at 5000. Four more settle as `none` because the request
is the collection: Qdrant's collection list, Cohere's embed, Weaviate's batch
create, FusionAuth's user lookup. Cohere's rerank is a page size called `top_n`
in a body, with nothing to resume from.

Braze is the first route in this collection to need `first_page: 0`. Its own
description: "The page of campaigns to return, defaults to 0 (returns the first
set of up to 100)" -- and there is no per-page parameter at all, so a hundred is
what Braze serves and the caller cannot ask for fewer. A client that starts at
page 1, which is what every other paged provider here would want, silently skips
the first hundred records and is never told. That key existed for exactly this
and had never been used.

Hetzner repeats Aha!'s lesson. Its `/servers` operation enumerates `name`,
`label_selector`, `sort` and `status` in its own parameter list, and `page` and
`per_page` appear nowhere on it -- they are in the API-wide section, along with
the `meta.pagination` object every listing returns. Twice now, a generator that
reads only per-operation parameters would produce a client that cannot page at
all.

Eight GraphQL listings move too, and the honest framing matters. These Recipes
each say "nothing is filtered, ordered or paginated -- `first`, `after` ...
accepted and honoured by nothing", because no GraphQL is parsed here. That is
still true of a query that inlines its arguments in the document. It is no longer
true of the other spelling: `in: body` with `variables.first` and
`variables.after` reads the variables a Relay-style client actually sends, so
`{"query": "...", "variables": {"first": 2}}` gets two records. The alternative
being replaced is the engine reading `limit` out of the query string of a POST --
a name no GraphQL client has ever sent. Linear's own header supplies the one real
default among them: "first defaults to 50, which is enough that a small workspace
never notices and a real one silently truncates."

Ten more, and the defaults keep getting further from ten. Papertrail's log search
serves **1000** by default and caps at 10,000 -- a hundred times what the fake was
answering, on the endpoint where a short answer is least likely to be noticed,
because a log search that returns fewer lines than you expected looks like a
quiet system. Missive is the opposite shape: two listings on one provider sharing
one parameter name, where conversations default to 25 with a ceiling of 50 and
messages default to 10 with a ceiling of **10**, so a client that walks
conversations fifty at a time and reuses the number on messages is silently
answered with ten.

Onfleet supplies the inverse of everything `max_limit` was built for. It has no
page size parameter at all, and its own reference says the endpoint "will return
up to 64 tasks but may return fewer". A short page is therefore *not* evidence of
the end -- the end is `lastId` going absent -- so the loop-termination heuristic
that breaks against a provider which trims also breaks against one that simply
undershoots, in the same direction, for the opposite reason. This format could say
a provider trims (`max_limit`) and could say it refuses
(`over_limit`), and had no way at all to say either of these. It has one now.
`may_overshoot` and `may_undershoot` are declared on the two routes that needed
them, and the emulator acts on both: a page that is not the last serves one extra
record where the provider overshoots, and one fewer where it undershoots -- never
zero, because an empty page ends a walk for a reason that is not true.

Modern Treasury then proved half of that first design wrong within the day. The
validator refused `may_overshoot` and `may_undershoot` together as opposite
claims about one page. Modern Treasury makes both in one sentence about one
endpoint: "the actual number of records returned may be less than, equal to, or
more than the requested amount." They are not opposites, they are one provider.
The rule is gone, and when both are declared they take turns by position -- the
first page overshoots and every page after it undershoots -- because doing both
at once would cancel out and demonstrate neither. Modern Treasury draws the
conclusion the keys exist for: keep paginating until the cursor header is absent,
because the page size tells you nothing about whether you are finished.

Both are deliberate and deterministic rather than random. A fake that misbehaves
only sometimes is a fake nobody can write a test against, and the entire point is
to break the caller's loop *here*:

```
while len(page) == limit:
    page = fetch(next)
```

Overshooting breaks it on the very first page, because 26 != 25. Undershooting
breaks it in the middle, because 24 != 25. Neither errors, neither looks short --
the overshooting one hands back *more* data than was asked for -- and both report
a partial collection as the whole thing. The adjusted number is also the one the
offset arithmetic uses, so a page that served 24 advances the position by 24; that
is checked by a test that walks two pages and asserts the second does not begin
on the record the first ended with.

Two more positions are not really cursors and are declared as cursors anyway,
with the mismatch written down rather than hidden. Papertrail's `max_id` is a
window boundary a client sets to the `min_id` of the page it just read -- and the
Recipe had already noticed that `min_id`/`max_id` are JSON *strings* in the
envelope while `id`, the identical magnitude, is a bare number on each event
inside it, so the value you page with has a different type from the value you
lifted it off. Missive's `until` is a timestamp the caller computes from the last
record's own `last_activity_at`. None of cursor, offset or page describes "read it
off the last record yourself"; cursor is the closest, and saying so beats saying
nothing.

Thirteen listings settled without reading a single new document, because their
Recipes had already done the reading and had nowhere to put the answer. Shortcut:
"this endpoint takes no page size, no offset and no cursor, and answers with a
bare array of everything". Trello: "neither the board's lists nor its cards takes
a limit, an offset or a cursor". Xero: "the Accounts endpoint takes neither page
nor pageSize -- unlike Invoices and Contacts beside it". DocuSign: "the
recipients endpoint takes neither count nor start_position". Stytch: "this
endpoint takes user_id and nothing else". LiveKit: "ListRoomsRequest has no
limit, page size or cursor of any kind in the .proto -- names is its only field".
Google Maps: "the Geocoding API is not documented as paging at all -- it answers
whatever matched, once". Seven providers, thirteen routes, every one of them
checked against the provider's own description by somebody who then had to write
the finding in a comment and watch the engine page it at ten anyway.

Two of those comments said so out loud -- LiveKit's "Cauldron's own runtime pages
it at ten anyway, reading a `limit` parameter that does not exist on the real
request message", Google Maps' "Cauldron's own default (page at ten, reading a
`limit` parameter Google does not accept)". Those sentences are now false, which
makes them exactly the failure mode `TestNoRecipeAssertsARetiredLimitation`
exists to catch, so it grew three more phrases and was red-greened against a
restored copy of Mapbox's old wording.

Mapbox itself is the most uncomfortable of the batch. Its header recorded that
"Mapbox's own `limit` parameter ... is documented but not modelled", and because
the runtime's fallback name is `limit`, the route had been honouring the caller's
page size all along -- right name, no declaration, wrong numbers. Mapbox's real
default is **5** and its ceiling is **10**; the fake served ten and would have
served fifty to anyone who asked. A guess that uses the right name is the hardest
kind to notice, because everything works until the numbers matter.

SavvyCal is the clearest case yet of a Recipe holding its own answer. Its header
already recorded, from reading three of SavvyCal's schema pages side by side,
that two of its list endpoints page "cursor-based, nested under a sibling
`metadata` object -- `{"entries": [...], "metadata": {"after", "before",
"limit"}}` -- confirmed by both schemas' own worked examples". Those three names
*are* the parameters. The declaration underneath that paragraph was empty, so the
runtime read `limit` (right by luck) and `cursor` (which SavvyCal does not have),
while the paragraph two hundred lines up said `after`.

Four more will not settle, and now say what was tried instead of saying nothing.
Deel's reference has moved twice: `developer.deel.com/reference/<operation>` now
answers "You have been redirected to the new developer documentation ... Use the
search bar above this page", the replacement renders one endpoint at a time from
JavaScript, and its own `llms.txt` -- published for exactly this purpose --
indexes marketing pages and a create-contract example without reaching a list
endpoint. Mailgun's and Bill.com's references answer 404 with a quarter-megabyte
of application shell to anything that is not a browser. AppSignal has no API
reference at the address its own documentation index points at. Each keeps its
guessed page size, labelled as a guess, next to a dated record of the attempt --
which is worth more than a blank, because the next person to look starts from
what already failed.

Eight of the routes that declared paging with a parameter left unnamed are
settled, and FlightAware is the one worth reading twice. AeroAPI's only two
paging parameters are `cursor` and `max_pages`, and `max_pages` is not a page
size -- "Maximum number of **pages** to fetch", so a client that read it as one
would ask for 2 and be billed for two whole pages. There is no size parameter at
all. The page is **fifteen**, and FlightAware says so on its *pricing* page
rather than in its API description: "Pricing is based on result sets, with one
set equaling 15 records", beside "one page being equivalent to one result set".
The page size is a billing unit. The Recipe had been declaring fifty, which was
nobody's number.

That correction broke a conformance case, which is the useful part. The case sent
`?limit=1` and asserted a `links.next`, and it had always passed -- because the
Recipe read `limit` too, so the emulator produced a short page in response to a
parameter AeroAPI has never had, and the next link was the fake answering its own
invention. With the size named `"-"` the parameter went inert and the link
vanished. The fix is the one a real caller has: a fixture with sixteen airports
against a page of fifteen. That is the ninth conformance case this sweep has
found asserting behaviour that existed only because the Recipe and the case
shared a mistake.

ClinicalTrials.gov retires a claim rather than adding one. Its Recipe said the
position was deliberately undeclared because "nothing here ever reached a second
page, and a name nobody has seen the API answer to is a guess whatever its
provenance". It was walked on 2026-09-05: `pageSize=2`, follow `pageToken`, and
records three and four come back -- not eleven and twelve. Worth having checked,
because the tokens for `pageSize=2` and `pageSize=10` share their first twenty
characters and only diverge after that, so a client comparing cursors by prefix
sees two different positions as one page.

100ms has a page size **floor**, which is rarer than a cap: "Default : 10,
Allowed values : Min : 10, Max : 100". A client asking for two rooms is answered
with ten and not told. There is a key for the ceiling and none for the floor, so
that one is written in prose beside the declaration.

Sixteen more came out of the providers' own descriptions rather than out of a
response, and one of them argues against trusting a machine-readable spec on its
own. Aha!'s `openapi.json` declares **no query parameters at all** on
`/custom_field_definitions`, and its prose documentation, two clicks away, says
`page` is 1-indexed and `per_page` defaults to 30 "up to a maximum of 200". A
generator reading the spec produces a client that cannot page; a Recipe checked
against the spec would have concluded `none` and been confidently wrong. The
prose was the better source, which is not the direction that usually runs.

Transitland supplies the other warning. Its cursor is `after`, documented as "an
opaque value created by the server" -- with the immediate caveat that "for
historical reasons, this is based on the integer record ID values", which is an
opaque token that is not opaque and will be depended on by accident. And the same
route takes a query parameter literally named **`next`** that has nothing to do
with the next page: it selects "departures leaving within the next specified
number of seconds". Guessing from the parameter list would have produced a client
asking for the next ten minutes of departures and calling it page two.

Replicate names its cursor nowhere except inside the URL it hands back --
`"next": "https://api.replicate.com/v1/models/.../examples?cursor=..."` -- and
takes no size parameter at all. Tavily's `max_results` defaults to **5** and
lives in the POST body, next to a `chunks_per_source` that defaults to 3 and is
not a page size; Exa's is `numResults`, also in the body, next to a
`startPublishedDate` that looks like a position and is a date filter. Seven more
settle as `none` because the request already is the collection: Plaid's accounts,
Metronome's events by transaction id, Budibase's table search, Exa's contents by
id, Finch's pay statements by request, Adyen's payment methods for one
transaction, and Nexus's repository list, which is configuration rather than
data. Nexus's *components* listing keeps its `continuationToken` and gains a
`limit_param: "-"`: the spec lists two query parameters and neither is a size.

Sixteen more settled in one sweep of the anonymously reachable, and the pattern
holds: every provider probed was doing something the fake was not. Three of them
answered a question nobody had asked. DataCite's `totalPages` is not the total
divided by the size -- 400 at size 25, 5000 at size 2, both exactly ten thousand
records against a real total of 134,263,030. It is a **deep-paging cap wearing
the name of a total**, and a client that walks to the last page it advertises has
seen 0.0074% of the collection believing it reached the end. NOAA's search
service does the same thing in two fields at once: `"count": 10000` printed
beside `"totalCount": 86045`. And NCBI returns `count`, `retmax` and `retstart`
all as **strings**, so a client adding position to size to find the next offset
concatenates `"2"` and `"2"` into `"22"` in every language that overloads `+`.

Three had defaults far from the fake's ten. openFDA serves **one**. OSV serves
**three thousand** and takes no size parameter at all, its `page_token` riding in
the POST body rather than the query string -- the first `in: body` this sweep has
needed. The Food Standards Agency serves **five thousand**, which means the
paging loop a client writes against it runs exactly once until the first local
authority with more establishments than that.

Reaching the FSA at all took following its own refusal: an unfiltered
`/Establishments` answers **403** with "This is a CPU intensive query: please use
one of the documented filters" -- a 403 about query cost rather than about
credentials, recorded in that Recipe because a 403 that means "narrow your
question" is not one a client should retry with better keys.

The USGS earthquake catalogue is another whose answer was already in its own
header: "ask with one and count is gone, replaced by limit and offset". Both
names sat in a response the Recipe asserts while the declaration beside them
said nothing. The header also records the part the names cannot carry -- that
offset is **one-based**, so "a request that sent no offset at all comes back
saying offset: 1" -- and the engine's offset positions count from zero, which is
now stated at the declaration rather than only at the top of the file. The World
Bank echoes `per_page` and `page` back in the metadata object it puts in front
of its data, with `per_page` as the string "2" beside a `page` that is a number.

Spaceflight News prints `limit` and `offset` in its own next link, and its
default of ten happens to match the runtime's -- so the size was right by
coincidence and the position never was. Nominatim reads a `limit` and has no
position at all: a geocoder answers the best matches rather than a collection
somebody walks. TfL's modes and line statuses come back whole, because London
has the lines it has.

**openFDA's default page size is one.** A bare request answers a single drug
label with `meta.results {"skip": 0, "limit": 1, "total": 262663}` -- the API
echoing both parameter names back at the caller -- and this Recipe was serving
ten. A client that never sets a limit saw ten here and one from openFDA, which
is the difference between a loop that runs and a loop that thinks it has
finished. Zenodo takes `size` and `page`, both visible in its own next link,
against a collection of 7,266,591 records. UK Police's two listings are
reference lists -- fifteen crime categories, thirty-six months of dates -- and
come back whole.

**Open Trivia DB's answer was in this file.** The row for it above says "the
API has no pagination to declare -- `amount` caps at fifty and there is no
offset or cursor, which is why its five listings sit in the count above". The
count it refers to is this one, and the Recipe went on being paged at ten
reading `limit` while the backlog explained why it should not be. Struck live
again to confirm the size, and declared.

Five more do not page, all struck live: PoetryDB answers an author's 362 poems
whole and ignores a `limit`, TheSportsDB the same for its leagues, Nager.Date
returns a country's year of public holidays entire, TheMealDB answers the one
meal asked for, and Open Notify answers everyone currently in space -- twelve
people, with the count beside them, on an endpoint that takes no parameters and
a collection kept small by physics.

The Bank of Canada's Valet API names its size `recent`, and the name is the
point: a bare request answers the whole 2,415-observation series and `?recent=2`
answers the last two, counting backwards from today rather than forwards from a
start -- which is why there is no offset to go with it. Football-Data answers
all 189 competitions with `count: 189` beside them and ignores a `limit`.
deps.dev answers a whole dependency graph, 71 nodes and 128 edges, and the
runtime serving ten nodes of that was not a short answer but a broken one.

Deck of Cards is the one that changed a case. It draws **one** card when
nothing asks for more -- struck live, `remaining` 51 -- and its Recipe was
relying on the emulator's own default of ten, a number that API has never used.
Two of its cases asserted a third card they had not asked for and passed because
the fake was handing out ten. They send `count` now. Its size is `count` and
there is no position: a second page of a draw is another draw, and what is drawn
leaves the deck.

Six more do not page at all, each struck live: Advice Slip's search answers
everything that matched and ignores a `limit`, Dog CEO's sub-breed list comes
back whole, Kraken's ticker and Open Library's books answer exactly the keys the
caller named, PDBe answers one entry, Open Meteo answers a forecast.

iTunes cannot be paged past its first two hundred results, and nothing said so.
`?limit=200` answers two hundred, `?limit=201` answers two hundred as well, and
`?limit=2&offset=2` answers the same two records `?limit=2` does, beginning with
the same track -- so `offset` is not read and there is no other name for a
position. CoinGecko's `/coins/list` answers all 19,616 coins to a bare request
and the same 19,616 to `?per_page=2` and `?limit=2`; its price endpoint answers
exactly the ids the caller named. Carbon Intensity answers one reading in an
array of one. None of the four pages, and all four were being paged at ten.

Deezer and Gutendex went the same way. Deezer's own next link carries both
names -- a bare search answers twenty-five with `?q=eminem&index=25`, and
`?limit=2` answers two with `&limit=2&index=2`, so `index` is an offset in
records rather than a page number. Gutendex has a page and no size: thirty-two
books, `?page_size=2` answers the same thirty-two, and the position's name was
already sitting in a string this Recipe asserts -- its own live finding is that
an unparseable filter is "carried faithfully into the next link,
`?ids=abc&page=2`".

### The Recipe often knew already

Paystack is the clearest of them and was live-verified when it was written: /bank
"does when asked for with no parameters at all: the whole collection, unpaged --
which is also genuinely what the live endpoint does until use_cursor is named",
and ?perPage=abc "returns the entire list exactly as no perPage does". The
runtime was paging it at ten reading `limit`, a parameter Paystack has never had
on that route under any flag.

Europe PMC needs no reference at all: it echoes the request back inside the
response this Recipe already pins, as `{"cursorMark": "*", "pageSize": 25,
"synonym": false}` -- both names and the default, in a field the file was already
asserting. ClinicalTrials.gov names its own in a refusal: `?pageSize=notanumber`
answers "Value provided in parameter `pageSize` cannot be converted to 32-bit
integer", struck live and pinned in a case. A provider that names a parameter
while refusing it has named it.

And where a provider is public, the position can just be asked. FBI Wanted and
Open Brewery DB both take a page number beside their size, and rather than
inferring the partner of `pageSize` and `per_page` from convention, both were
struck live on 2026-09-05: `?pageSize=2&page=1` and `page=2` answer two
different pairs out of 1,240 wanted persons, and the same for breweries. Open
Brewery DB's size was already pinned here by its sharpest finding -- `?per_page=201`
does not clamp and does not refuse, it answers **302 to the landing page**, so a
client following redirects gets a 200 and the wrong body. A parameter a provider
redirects you for is one it reads.

ClinicalTrials.gov is also why the unnamed count went **up** by one while the
unstated count fell by four. Its size is established and its position is not --
nothing here ever reached a second page -- so the route names one and leaves the
other, and the counter is right to say so. A half-answer that says which half is
better than silence and better than a guess.


Six of these were settled without reading anything new, because the Recipe had
already read it and could only write it down. Toggl's header says "since,
before, start_date, end_date, meta and include_sharing are the whole query
parameter list" and then that "Cauldron's own runtime still applies its
undeclared-route default (page at ten, reading `limit`)". Agora's says every
project "comes back in one answer" and then the same thing. Amadeus's names the
`max` cap the endpoint really takes and then notes the default applying instead.
Hotelbeds' names `from`, `to` and `total` and then the default. OpenAIRE did not
even need prose: its own live-verified failure says "Size and page arguments
have exceeded the maximum number of returned results", and the response header
it pins carries `size` and `page`.

Every one of those was somebody doing the work and having nowhere to put the
answer -- the engine's fallback showing through a Recipe that had established
there was nothing to fall back to. That is what `style: none` and the parameter
names are for, and it is worth checking the prose before going looking for a
document.

### Silence could not say "this one does not page"

That was half the difficulty, and it was the format's rather than anybody's
diligence. A listing with an empty paging block and a listing whose provider
serves the whole collection looked identical in a Recipe, and both were served
ten records with a cursor. `pagination.style: none` says the second one now: no
page size, no position, no pointer to a next page, and none of the three read
from the request. Declaring anything else beside it is refused, because `none`
is the whole statement.

**Forty-three routes across twenty-seven Recipes were settled with it, from the
descriptions those Recipes already record.** Eighteen listings have a
description declaring no query parameters at all: OpenAI's `/v1/models`, xAI's,
Perplexity's, Supabase's projects and its secrets, Upstash's Redis databases,
Turso's organizations, Redis Cloud's subscriptions, Stability's engines,
Checkout.com's payment actions, Livepeer's assets, NocoDB's bases, Checkly's
runtimes, ClickHouse's organizations, SurrealDB's instances, Inngest's event
runs and two of Scout APM's. Fourteen more declare a query string with filters
in it and no paging parameter of any kind, which says more than silence does:
Fly's machines, GitBook's pages, Mistral's models, Together's, ElevenLabs'
voices, Flagsmith's flags, Spotify's tracks, Tink's providers, Basiq's
accounts, ClickHouse's services, Turso's databases, Watchmode's search, Podcast
Index's dead feeds and Scout APM's error groups.

Eleven more had paging parameters written down and unread, and are now
declared: Apify's dataset items (`limit` and `offset`, with "By default there is
no limit" stated outright), Belvo's links and accounts (`page` and `page_size`,
a hundred by default and a thousand at most), Checkly's check results (`limit`
and a `nextId` cursor, beside a checks listing on the same API that pages by
number), Chroma's collections, OpenAQ's parameters and locations (`limit` and
`page`), Scout APM's endpoints, and all three of Svix's listings (`limit` and
`iterator`, the same way every time).

Nine more were settled by finding the description first. `cauldron discover`
proposes a document only where it declares a path the Recipe already models, and
it found four the collection did not have: Orb's, Increase's, PostHog's and
Column's, now recorded and fingerprinted. Orb pages by `cursor` and `limit` on
all three of its listings, twenty by default and a hundred at most. Increase's
own description collapses two numbers into one sentence -- "The default (and
maximum) is 100 objects" -- so a caller cannot ask for more and never needs to
ask for the most. PostHog is Django REST Framework's `limit` and `offset`, which
is also why the next pointer beside it is a whole URL. Column takes `limit`,
`starting_after` and `ending_before`, a cursor that is a record's own id in both
directions.

All four were found the same way: those Recipes already served a `next_cursor`,
a `has_more` or a `next` in the response and said nothing about how a caller
sends one back. That contradiction is a good place to look -- twenty-four
listings across sixteen Recipes were in it, and it has already paid for itself
twice more. Typesense's search takes `page` and `per_page`, from the
description Typesense publishes in its own GitHub organisation, and documents
`limit` and `offset` as alternatives for the same job that this Recipe serves
one half of. CometChat's groups listing takes `page` and `perPage`, a hundred by
default and a thousand at most, out of the `chat-apis.json` this Recipe already
recorded -- `cauldron check` cannot match it automatically there, because the
description's server carries a `/v3` prefix the Recipe's routes do not.

Eight more were settled straight out of the twenty-nine descriptions
`cauldron discover` found. Netlify's sites and deploys take `page` and
`per_page`; Finch's employer directory takes `limit` and `offset`, with both
numbers in the parameter's own words -- "defaults to 100, maximum 10000";
OpenAlex takes `per_page` and `page`, and names its own ceiling in `page`'s
description -- "Use cursor for deep pagination beyond 10,000 results" -- so the
cursor half is recorded rather than modelled. Four do not page at all: Pinecone's
indexes and Finch's providers declare no query parameters and answer the whole
list; monday.com's single GraphQL endpoint carries paging inside the query
rather than beside it; and iNaturalist's single-taxon fetch answers inside a
paging envelope with nothing to page, which is the finding that Recipe's header
opens with.

Docker Hub's repositories listing had no paging block while the tags listing
beside it had a careful one, and both take the same `page` and `page_size` --
struck live with no credential, `?page_size=2` answering two of a hundred and
eighty-one. The same probe found something the Recipe had wrong for both:
`next` comes back as a whole URL, scheme and host and query included, and this
Recipe served a bare page number under a field whose value a client is meant to
request as it stands. Its own case asserted `next` matched `.`, which is true of
a page number and of a URL alike, so a live-verified case had been carrying the
wrong claim. It now pins the URL.

**And it found a listing that is not one.** AWS SES's `GET /v2/email/account`
was modelled as a collection, and AWS's own service model describes
`GetAccountResponse` with its members at the top level: no collection, no
identifier, and `Max24HourSend` and `SentLast24Hours` nested under `SendQuota`
where this Recipe had them flat. Code written against the fake read
`Accounts[0].SendingEnabled` and would have found nothing against SES; code
written against SES read `SendQuota.Max24HourSend` and found nothing here. The
route now collapses its single record and drops the identifier it was minting,
and `GetAccount` takes no input at all -- it is the one operation in that
service with no `NextToken` and no `PageSize`.

Scout APM's endpoints listing is the one to remember. Its description says
"Omit for the full listing (no default limit). Supplying this switches the
response to the paginated shape" -- so `results` is an array until you send a
`limit`, an `offset` or a `sort_by`, and then it is an object carrying
`endpoints`, `count`, `total_count` and `has_more`. A client reading
`results.length` works until somebody adds a sort. That switch is stated in the
Recipe and not modelled.

## What a delete answers with, mostly closed

A delete used to fabricate Stripe's receipt for every provider, using keys no
Recipe declares. It now answers 204 with no body unless the Recipe says
otherwise, which is what most providers do and is the safe direction: a client
calling `.json()` fails here exactly as it would in production.

Verified against the providers' own descriptions with `cauldron check`:

| Provider | Answers | Declared |
|---|---|---|
| Stripe | 200 with a receipt | `deleted_body: receipt` |
| Clerk | 200 with a DeletedObject | `deleted_body: receipt` |
| Airtable | 200 with id and a flag | `deleted_body: receipt` |
| Box, Twilio, PagerDuty, DigitalOcean, Adyen | 204, nothing | the default |
| Zoom | 204, nothing | the default, with a case |

Clerk was a real regression from the change that introduced the default, found
by the checker rather than by anybody noticing.

The delete routes on recipes with no published description are still on the
default. That is a default rather than a verified claim, and it is at least
wrong in the direction that shows up locally.

## Second sweep

Sixty-one, chosen for behaviour rather than for surface area. The same
convention applies: an entry whose reason starts with **Assess** names a
surface, and somebody still has to read the documentation and decide whether
there is a lie in there worth reproducing.

### Money that moves

The standing rule applies to every one of these: model the reads, leave the
transfers alone, and say so in the header.

| Provider | Why |
|---|---|
| ~~Increase~~ | Shipped. A return is a new object days later, the transfer stays unmarked, and the reason decides whether retrying is legal |
| ~~Modern Treasury~~ | Shipped. Every amount appears twice and cancels, one direction means opposite things on two accounts, three balances disagree on purpose |
| ~~Marqeta~~ | Shipped, without the JIT funding path. It puts a webhook in the authorisation with a response deadline on it, so a slow test is a declined payment, and Cauldron delivers webhooks without waiting for an answer. The Recipe states the gap rather than pretending to close it |
| ~~Dwolla~~ | Shipped. Three attempts and then the funding source is finished, not rate limited and not retryable tomorrow |
| ~~Checkout.com~~ | Shipped. The payment's action log is a separate resource, which is what this row was circling |
| ~~Klarna~~ | Shipped anyway, from live probing rather than from the portal. The documentation is still unreadable and both published locations for a machine-readable description answer 404, so the Recipe records none. What it found instead is better than the premise this row had: Klarna's own reference tells integrators not to rely on the status field it publishes |
| ~~Airwallex~~ | Shipped. The settlement currency is not the charge currency and the rate was fixed at a moment nobody picked, there is a separate balance per currency so a payout can fail on a funded account, and a partial capture leaves two amounts different forever |
| ~~Column~~ | Shipped. A notification of change is not a failure, and R01 and R07 are two characters apart with opposite obligations |
| ~~Moov~~ | Shipped. **Every wrong path is a 403** -- unknown path, wrong method and the bare root, checked four ways, never a 404 or 405 -- so the API declines to say that anything does not exist and a client cannot tell a typo from a permission it lacks. A missing credential and a well-formed fake token are indistinguishable in turn: identical 401, zero bytes, no content type. The payment-method polymorphism this row was queued for is real and models cleanly, because it is not a discriminated union at all -- every shape is declared as a sibling and one is populated. What is **not** expressible is the agreement between them: nothing here, and nothing in Moov's own schema, requires the populated sibling to match what `paymentMethodType` claims |
| ~~Ory Network~~ | Shipped as `recipes/ory`, and the row was right that most of the surface is Kratos and Hydra under another name -- those two already ship and the new Recipe deliberately does not copy their session and identity shapes. What it adds is what is left when you take those away: **the hosted playground answers one 404 to everything**, the same 155 bytes across four distinct paths, with a `ruleId` that never varies. And the row underestimated the obstacle. Every route that matters is behind a Cloudflare managed challenge -- 403 with `Cf-Mitigated: challenge` -- so Ory's own software is never reached at all. That is reported and not modelled, which is the only honest option Ory Network is the hosted form of software this collection already describes: its OAuth2 surface is Hydra, which ships, and its self-service and identity paths are Kratos's, which also ships -- the same routes, the same nested `error.code/message/reason/status` envelope. Probing confirmed it: every self-service, session and OAuth2 path on the public playground answers Cloudflare's bot challenge rather than Ory, and the one endpoint that did return JSON reproduced the envelope the Kratos Recipe already pins. The one surface that might genuinely differ is Ory Permissions, whose relationship-tuple model fits neither `auth` nor `identity` and has no reachable public instance. A third Recipe here would duplicate two |
| Unit | Assessed and left. That is the Marqeta Recipe's headline, already shipped: an authorization is a hold, the clearing is a second object pointing back at it, and the amounts legitimately differ. A Unit Recipe would need to earn its place on the banking-as-a-service parts -- the application lifecycle, the counterparty model -- rather than on the card shape |

### Billing that accrues

| Provider | Why |
|---|---|
| ~~Orb~~ | Shipped. A draft is a running total, a closed period can still be amended, and money is a decimal string |
| ~~Metronome~~ | Shipped. Still true that deduplication inside a window is Orb's fourth trap, and the Recipe says so rather than claiming it. What earned the place is what this row predicted: deduplication by idempotency key inside a window is the Orb Recipe's fourth trap, in those words. Checked before writing rather than after, which is what the Lemon Squeezy row taught. A Metronome Recipe would have to earn its place on what differs -- ingest is fire-and-forget, so a bad event is accepted with no body saying what was taken and no error to catch, and finding out means asking a second endpoint or waiting for the usage not to appear. The second endpoint exists and Metronome tells you not to lean on it -- "heavily rate limited and designed for sampling workflows only. Do not use this endpoint to check every event in your system" -- so the only observability over a billing pipeline is a sample. Metronome names the failure itself: "Silent revenue loss occurs when events are dropped, delayed, or fail to match billable metrics." Also: an event can be accepted and match no customer, or match a customer and no billable metric, and both are 200s; the whole error schema is {message} with nothing to branch on |
| ~~Lago~~ | Shipped. The difference this row asked about turned out to be modellable and is in the schema rather than in the deployment notes: lago_id is "not guaranteed to be a UUID; on organizations using the ClickHouse events store it is a composite string", so the identifier's format depends on which events store the organisation runs |
| ~~Lemon Squeezy~~ | Shipped on those terms, and the objection was right: four of the traps a first draft led with — merchant-of-record tax, amounts twice as number and display string, cancelled-but-still-running, test and live sharing an API — are all FastSpring's and are now stated as repeated rather than claimed as the reason. What the Recipe exists for is what this row asked for: licence keys, where the activation count runs out while the subscription stays healthy, and store-scoped identifiers. Two things beyond it: totals in two currencies on one record, and a partial refund whose refunded boolean is false. One request went unmet — licence activation is public and authenticated by the key itself, and Cauldron authenticates a Recipe rather than a route, so only the store-key half is modelled |
| ~~RevenueCat~~ | Shipped. There is no `is_active` field and RevenueCat's own guidance is to read one -- it is an SDK property, so the moment the question moves to a server somebody writes the comparison by hand and the advice stops applying. Four active entitlements in the fixture, active for four different reasons: cancelled, lifetime (`expires_date: null`, which every naive comparison reads as expired), failing to pay inside a grace period, and somebody else's family purchase on a trial. Entitlements are keyed by your names and subscriptions by the stores', and the endpoint is a GET that creates: 200 found, 201 invented |
| ~~Recharge~~ | Shipped, as the narrower thing this row described. The cross-system drift still needs both systems and a Recipe still holds one, so what shipped is what was suggested here: a subscription that is ACTIVE while its charge has reached MAX_RETRIES_REACHED, and external ids that are Shopify's numbers held as strings. A charge also carries the Shopify order id beside its own, which shows the two numbers for one payment without needing Shopify running -- it shows they differ, not that they drift |
| Maxio | Assess — the former Chargify, beside Chargebee and Recurly |

### Payroll, people and hiring

| Provider | Why |
|---|---|
| ~~Rippling~~ | Shipped, and the assessment stopped at the front door in a way worth recording. **The surface Rippling renamed did not move**: `GET /platform/api/workers` -- "worker" being the term its own current documentation uses -- 404s live, while `GET /platform/api/employees` still answers 401. An integration written from today's documentation reaches nothing; one written years ago still reaches the gate. Every credential shape collapses to one Django REST Framework sentence, `{"detail":"Authentication credentials were not provided."}`, including three shapes where a credential *was* provided. No employee field is modelled: developer.rippling.com will not show a schema without an account |
| ~~BambooHR~~ | Shipped, and it answers the group's question about terminated employees **in its own prose rather than by inference**: the directory excludes inactive and former employees, and BambooHR states the consequence outright -- absence from that response means the employee is not in the published directory, not that no such employee exists. So an integration reconciling a roster cannot tell "terminated last week" from "never sent". Confirmed against its field architecture too: the directory fieldset structurally cannot carry a status field and the single-employee endpoint can, which is two endpoints for one record disagreeing about whether the record's most important attribute exists. No case carries a date -- live probing never got past tenant-subdomain resolution, and the header says so |
| ~~Ashby~~ | Shipped, and the assessment was answered sideways. Stages are configurable, as the row guessed, but the finding that matters is that **Ashby's own auth documentation is wrong about Ashby**: it states a missing key is 401 and a wrong key is 403, and five credential shapes struck live against two endpoints all answered a byte-identical 401 in **plain text**, on an API whose documented success bodies are JSON. On the person-versus-application question the group asked, Ashby separates Candidate from Application and links them from the person outwards with `applicationIds` -- the opposite direction from Lever's `contact` |
| ~~Workable~~ | Shipped, and it gives the most decisive answer of the three to the person-versus-application question -- from its help centre rather than its schema. **A Candidate is scoped to a job.** The same person applying to two roles gets two entirely separate profiles, correlated only by matching the raw email field by eye in the interface, with no merge. Also verified: routing precedes authentication for every method, and Workable has **no method-not-allowed concept on the wire at all** -- PATCH and DELETE on a real path answer the identical 404 an unrouted path does, across five probes, with no `Allow` header anywhere |

### Messaging

| Provider | Why |
|---|---|
| ~~Plivo~~ | Shipped. **Its failure is one its own SDK cannot parse**: a 401 arrives as `text/html` on a JSON API, and Plivo's published Python client is written to read `{"error": ...}` -- reading that client's source shows parsing throw, the exception swallowed, and a generic sentence substituted, so the vendor's own SDK reports something Plivo never said. A path matching no route answers JSON keyed `message` where documented failures use `error`. The send model completes a four-way comparison across this collection: Twilio has no multi-recipient primitive, Bandwidth mints one shared id, Sinch mints one with opt-in detail, and **Plivo mints N bare ids with nothing tying them together** -- lose one and that message can never be asked about again |
| ~~Bandwidth~~ | Shipped, messaging only. The same message comes back with entirely different property names depending on which endpoint returned it, and a send is 202 Accepted rather than sent. Number provisioning is genuinely asynchronous and would be worth having, but that API answers in XML |
| ~~Twilio Verify~~ | Shipped. A wrong code is a 200 with the verdict in a word inside the body, so `if (!res.ok)` lets the wrong person through the second factor. Three different 429s, two of which are terminal and one of which is not. A verification is deleted on approval, so checking twice and never having started are the same 404 |
| ~~Customer.io~~ | Shipped, and this row had it right: the Track and App APIs are separate hosts with separate credential schemes. What the Recipe adds is that **nothing tells you when you have confused them** -- sending the App token to the Track host answers exactly what sending nothing answers, and what a wrong secret answers, so three mistakes share one response, separated only by a `WWW-Authenticate` header one host sends and the other does not. And **the API answers its own marketing site's 404**, a 3817-byte HTML page identical across both hosts and both kinds of mistake |
| ~~Braze~~ | Shipped. The export answers 201 with a prefix and no users; the file lands in cloud storage minutes later, so a test reading users off that response reads nothing forever |
| ~~Brevo~~ | Shipped, and this row's own correction still stands: the premise about limits was wrong and is kept here so it is not re-derived. What the Recipe found instead is that **Brevo's published description is behind the credential** -- fetching `swagger_definition.yml` without a key answers 401 -- so the document explaining how to authenticate cannot be read without authenticating, and no `upstream.spec` is recorded. It also distinguishes a missing key from a wrong one in prose while giving both the identical `code` of `unauthorized`, so the field a machine switches on cannot tell them apart |
| ~~Kit~~ | Shipped. The rename is not a migration: **api.convertkit.com and api.kit.com both serve both /v3 and /v4 identically**, one backend with two doorbells, and the two versions take incompatible credentials -- v3 reads `api_secret` from the query string and ignores an Authorization header entirely, v4 reads `X-Kit-Api-Key` or an OAuth bearer, and each scheme is blind to the other's. Kit's own upgrade guide says outright that V4 keys are not compatible with V3. The live 401 still carries `WWW-Authenticate: Bearer realm="ConvertKit API"`, the retired brand, in a header a client is meant to parse. On unsubscribing: one-way. The unsubscribe route is the only one that writes `state: cancelled`, the upsert says it cannot, and the update schema has no `state` property |
| Beehiiv | Assess — publications, posts, subscribers |

### Identity and risk

| Provider | Why |
|---|---|
| ~~AWS Cognito~~ | Shipped as `recipes/cognito`. Sharpest finding: **a pool id with the wrong region baked in is indistinguishable from one that never existed** -- the same `ResourceNotFoundException` sentence word for word, differing only in the id echoed back. A Cognito pool id carries its region in the first half, so the one thing the caller needs to know is the one thing the response declines to use. Three layers answer on that host, not two: a load balancer refuses non-POST verbs before Cognito's dispatcher runs, the dispatcher answers `UnknownOperationException` with no credential examined, and only a recognised operation reaches the credential check, which has three sentences of its own |
| ~~Ory Kratos~~ | Shipped, and the flow object was the unusual part. The server sends you the form: a flow carries the URL to post to, the method, and the list of fields to render, so a login page is a renderer for somebody else's JSON and a hardcoded {email, password} POST works against one deployment and not the next. An expired flow is 410 and 410 is not a retry -- Ory's words are that a new flow has to be initiated. The CSRF token is a node in the same array as the visible fields, a node's `type` and its attribute's `type` are different words at different depths, and `identity.state` is documented as having no effect while still carrying an enum |
| ~~Kinde~~ | Shipped, and the assessment it was queued for came out sharper than expected. Kinde **tells three credential failures apart under one status** -- a missing header, a bearer that is not JWT-shaped, and a JWT whose kid was never issued -- which is the first Recipe here to use all three of the credential verdicts the runtime learned this week. Auth resolves before routing, confirmed on `kinde.kinde.com`, Kinde's own tenant. The finding that was not anticipated is the contrast inside one hostname: the management API answers in rich JSON and `POST /oauth2/token`, the surface that exists to handle credentials, collapses every malformed request into ten bytes of plain text -- `not_found`. On the per-customer-hostname question this was grouped to answer, Kinde avoids the failure the Make Recipe records |
| ~~Persona~~ | Shipped. completed is not approved, needs_review is neither, and nothing is at the top level because it is JSON:API |
| ~~Onfido~~ | Shipped. complete is not clear, consider is neither a pass nor a failure, and the reason lives on the report rather than the check |
| Sift | Blocked on documentation, not on interest. The premise holds and the shape is good, and the public docs do not publish a score response body: the pages describe the 0-100 scale in prose while the API is widely reported to send decimals, and that number is the one every integration branches on. Building it would mean guessing the scale, which is the one thing a Recipe must not do. Worth revisiting with an account, where a single real response settles it |

### Storage and media

| Provider | Why |
|---|---|
| ~~Backblaze B2~~ | Shipped, and the row was right about where to look. Deleting a file brings back an older one -- Backblaze's words are that the most recent older version becomes the current version -- so delete-then-read is a 200 and stale bytes rather than a 404. There is no overwrite either, so the S3 habit of writing a key to replace it accumulates billed copies. A hide is a version with its own id, the action enum is documented open, the base URL is data and comes back from authorize, and part sizes are quoted strings beside a contentLength that is a number |
| ~~Uploadcare~~ | Shipped. A file exists before it is stored and unstored files are deleted after twenty-four hours, so the same code works on one project and loses files on another; a removed file still answers 200 with everything intact |
| ~~ImageKit~~ | Shipped, and the assessment came out as an absence. **There is no state field on a File because there is no state**: transformation happens at request time on the CDN edge, so nothing is ever pending and there is nothing to poll -- which is a real answer beside Transloadit's queue and Bunny's encoder rather than a missing one. It does have exactly one asynchronous job, and it is modelled for the contrast: a purge returns `{"requestId"}` with no status, and a second endpoint takes that same id and answers `Pending` or `Completed`. Its credential failures split by **status** rather than by wording -- 401 for absent and for a wrong scheme, 403 for a well-formed Basic credential naming a key nobody issued. No description exists to fingerprint: the documentation is a client-rendered application curl cannot reach, which is the dead end this backlog already recorded, so the source is the generated SDK's TypeScript interfaces |
| ~~Livepeer~~ | Shipped, and it gives the sharpest version of this collection's create-returns-what question. An upload request hands back **two statuses that both look authoritative and are not about the same thing**: the created Asset at phase `uploading`, and a separate `task.id` answering to a different vocabulary entirely. They share two words and disagree about the important one -- a caller polling the task and reading `completed` has learned nothing about the asset, whose terminal word is `ready` and is never `completed`. Also live: an unrouted path with no Authorization header 404s, and the identical path with a well-formed wrong bearer 401s naming the token back, so which layer answers depends on the **shape** of what was presented rather than on whether anything was |

### Data and search

| Provider | Why |
|---|---|
| ~~Qdrant~~ | Shipped, and **not** for the reason this row gave. The filter behaviour needs the search itself and Cauldron does not do vector arithmetic, so the header says outright that filters are not applied. What shipped instead is the envelope: `status` is the string `ok` on success and an object `{error}` on failure, so a typed client fails to parse one of the two paths and an untyped one reads `status.error` as undefined, which is falsy, which reads like no error. Plus: `result` nests twice, a collection listing hands back names only, a write answers `acknowledged`, `points_count` and `indexed_vectors_count` disagree on purpose, and `version` is on a query result but not on a point fetched by id |
| ~~Weaviate~~ | Shipped, and the thing worth having was the batch. A bulk import that half failed is a **200**: Weaviate's own words for the endpoint are that the request was processed successfully and individual object statuses are in the body, so each element carries `result.status` of SUCCESS or FAILED and the failures sit in the array beside the successes, identical in every other way. Plus: an error is an *array* of messages so `error.message` is undefined everywhere, the same object has two URLs of which one is deprecated, and the API says so in a `deprecations` field on the response body rather than in a header |
| ~~Typesense~~ | Shipped, and the thing worth having was the relevance score. `text_match` is an int64 up near 578730123365711993, which `JSON.parse` rounds to 578730123365712000 -- and the hit ranked below it, ...994, rounds to the same number, so two differently-ranked results compare equal. This is the worse half of a problem the format already knew: Discord sends snowflakes as strings so it cannot happen, and Typesense sends a number, so no client can avoid it. Plus `search_cutoff` is a documented boolean for "your results are incomplete" at 200, and there are three counts of which one counts the array |
| ~~Convex~~ | Shipped, and the assessment resolves the other way from the row's guess: the function-call routes are REST-shaped enough to model, and what does not fit is elsewhere. **Nothing tells you the deployment is wrong.** `*.convex.cloud` is a Cloudflare wildcard, so DNS never fails, and a gateway answering to `convex-usher` takes every request whatever the name -- validating the shape of the body before ever asking whether the deployment exists, so a malformed body is refused identically on a real deployment and an imaginary one, while a well-formed call to a fake deployment gets a bare empty 404 indistinguishable from an unmatched path. Also established: the function-call surface has **no platform credential gate at all**, because authorisation is the developer's own code. The control plane, where missing and wrong genuinely differ, is described in the header rather than served |

### Hosting and deployment

| Provider | Why |
|---|---|
| ~~Netlify~~ | Shipped. A deploy has an id before it has a site, ready is not published, and a missing URL means two different things |
| ~~Render~~ | Shipped, and both suspicions this row raised were right. The spec was not blocked, only mislaid: `api-docs.render.com` does 404 on the URLs guessed at here, but its own "OpenAPI Spec" reference page names the real one, `api-docs.render.com/v1.0/openapi/render-public-api-1.json`, which resolves. And `suspended` really is a string enum, `"suspended"`/`"not_suspended"`, not a boolean -- beside two more of the same shape, `autoDeploy` as `"yes"`/`"no"` and `notifyOnFail` as a three-way switch. List endpoints do wrap each element beside a cursor, `{cursor, service}`, which this format can reproduce for the wrapping and not for the per-item cursor sitting beside it |
| ~~Fly.io~~ | Shipped, and **not** for the reason this row gave: the older platform API is GraphQL at a different host, and this format speaks REST, so the two-shapes half is stated in the header and not modelled. What shipped is better anyway. `state: started` does not mean the application is up, and three fields on the same object independently say so -- `host_status` can be `unreachable`, `cordoned` can be true, and `checks[0].status` can be `critical`. Four answers to one question, disagreeing by design. Plus: `instance_id` is unique per *version*, so anything keyed on it loses its history at every deploy, and `nonce` is returned once, at creation, and only if a lease duration was asked for |
| ~~Heroku~~ | Shipped, and the header was the smaller half. A successful list is **206**: Heroku pages with the `Range` header and answers `206 Partial Content` while there is more, with the resume point in `Next-Range` rather than in the body -- so comparing against 200 rejects every page but the last, and testing `ok` accepts them and never looks for the rest. The `Accept` version header is a 406 when missing, errors are keyed by `id` rather than `code`, `url` is on an error only sometimes, and a formation with `quantity: 0` is a process type that exists and is not running |
| ~~Hetzner Cloud~~ | Shipped, and the row was exactly right. Powering off a server answers 201 with a job: status running, progress 0, finished null, and the machine still on. Every mutation in the API is that shape, so nothing that changes anything answers with the thing it changed. An action can fail long after its 201, its reason is an object rather than a sentence, progress reaches 100 while still running, a server has nine statuses of which eight are not running, and `locked` is a separate question from `status` whose refusal is a 423 |

### Observability and flags

| Provider | Why |
|---|---|
| ~~Bugsnag~~ | Shipped. An error carries counts and no stack trace, fixed is not terminal, and severity is a different question from unhandled |
| Grafana Cloud | Assess — the stack-management API and the Prometheus-shaped query API that is not REST. The Grafana HTTP API itself ships; see the row in the observability table above for what it pins |
| ~~Better Stack~~ | Shipped. Its 404 hardcodes the word GET whatever method was sent, and it is the only one of its group with both a last-checked time and paused as its own status word |
| ~~PostHog~~ | Shipped. A flag definition is not what a user gets, nought per cent is not inactive, and capture says the same thing whatever you send |
| Statsig | Assess — gates, experiments, exposure logging |
| ~~Flagsmith~~ | Shipped, and the point this row made is confirmed: a flag definition is not what a user gets, and the Recipe models the evaluated form. The finding it did not anticipate is that Flagsmith **answers three failures in three formats with no field in common** -- a printed Python tuple served as JSON, a zero-byte 404 for a wrong key, and a bare array holding one object for an unrouted path. Missing and wrong are both modelled, which most Recipes here cannot do, because they run through two different mechanisms rather than one credential gate |

### Webhook infrastructure

| Provider | Why |
|---|---|
| ~~Svix~~ | Shipped. Accepted is not delivered, the outcome lives on the attempts, a failing endpoint gets disabled and nothing announces it |
| ~~Inngest~~ | Shipped. Step state is documented rather than served, for the reason this collection keeps meeting: a run's `ended_at` and `output` stay unset for exactly as long as its status is Running and fill in together when it is not, and nothing here advances a fixture between requests. Verified instead: **its failures declare `text/plain` and send JSON**, so a client trusting the content type treats valid JSON as a string. Also pinned: the caller's own deduplication id is never the key the event is filed under |

### Web data

| Provider | Why |
|---|---|
| ~~Firecrawl~~ | Shipped. A running crawl already has results, total is an estimate completed never reaches, and a 200 can carry a failed fetch |
| ~~Apify~~ | Shipped. A run that SUCCEEDED tells you nothing about whether it produced anything: the status describes the process, and the results are in a dataset under a different id at a different endpoint. That endpoint is also the one thing this API does not wrap in `data`, so the code reading a run and the code reading its results unwrap differently against the same provider. Plus: three of the eight statuses are the `-ING` half of a pair, a TIMED-OUT run has partial data worth keeping, `isStatusMessageTerminal` says whether the prose will change, and the duration is there twice in two units |
| ~~ScrapingBee~~ | Shipped, and this row had it right: a failed fetch is a successful API call. The Recipe adds where the real status went -- it survives in `Spb-initial-status-code` and nowhere else, so `response.status` says nothing about the page requested. And the parameter that looks like the fix makes it worse: `transparent_status_code` makes ScrapingBee's own status **become** the target's, so a 404 is now either the proxy failing to route you or the page not existing, and the header invented to separate them agrees with the status line |

### Registries

| Provider | Why |
|---|---|
| ~~Docker Hub~~ | Shipped. The rate limit is a count and a window glued together, latest can be older, two tags can be one image |
| ~~npm registry~~ | Shipped. A tombstone in time that is absent from versions, latest is not the highest, deprecation is the presence of a string |

## Third sweep

Fifty-two more, none of them already above. The bar is the same as the first two
sweeps: a provider earns a row by having something specific an emulator can
pin and a mock cannot, not by being popular. A row that says "the usual CRUD"
is a row that should not be here.

### Brokerage and market data

An order is not a fill, and the gap between them is where the bugs live.

| Provider | Why |
|---|---|
| ~~Alpaca~~ | Shipped. A partially filled order is a real position that is still open, the listing hides everything that finished, and every number is a string |
| ~~Tradier~~ | Shipped, and the row was the smaller half. The shape of the answer depends on how many answers there are: Tradier's own words are that a single order is returned as a JSON obj/dict whereas multiple orders are returned as an array. Its own OpenAPI says `order` is always an array, contradicting the prose, so a generated client is written for the branch that is not always true. The account really is a path segment and nowhere in the order, an order id is an integer, and `quantity` is what was asked for while `exec_quantity` is what happened |
| ~~Polygon.io~~ | Shipped, and the row understated it. `t` is documented as the **start** of the window on the range endpoint and the **end** of it on the grouped one -- same one-letter field, same API, Polygon's own words on each. A missing bucket is also not a quiet one: Polygon does not populate an aggregate unless OHLC changed or *eligible* trades occurred, so absence means nothing eligible happened. Plus `T` is the ticker and `t` the timestamp on one object, `otc` is left off when false, `limit` limits base aggregates rather than results, and three counts mean three things |
| ~~Finnhub~~ | Shipped, though not on the row's claim, which is an account fact rather than a documented one. What is documented and worse: `c` is "List of close prices" on a candle and "Current price" on a quote -- Finnhub's own words -- so five one-letter fields are arrays on one endpoint and scalars on another. A candle response is seven parallel columns with no object for a bar, `no_data` is a 200 whose arrays are *absent* rather than empty, and the response echoes neither the symbol nor the range, which is exactly what makes a shortened range invisible |

### Open banking and financial aggregation

| Provider | Why |
|---|---|
| ~~TrueLayer~~ | Shipped. Consent expires ninety days from when it was granted rather than from last use, expired and revoked are both 403 with different fixes, pending and settled are two lists with different ids for the same money, and a balance never reconciles to either |
| ~~GoCardless Bank Account Data~~ | Shipped. Booked and pending are two arrays in one response, the same purchase moves from one to the other and changes its transactionId, the amount is a signed decimal string nested under transactionAmount, and a requisition status is a two-letter code that explains nothing |
| Salt Edge | A refresh is asynchronous and the connection reports success before the new transactions exist |
| Codat | Sync status and data are separate reads, so data is queryable while the sync that would change it is still running |

### Unified APIs

The normalisation is the product, which means the emulation has to reproduce
what normalisation does not fix.

| Provider | Why |
|---|---|
| ~~Merge.dev~~ | Shipped, HRIS only. A field is null because the provider lacks it or because it is empty and the response is identical; remote_data is null unless the account was configured for it; a linked account exists before it works, and can be authenticated with nothing synced |
| ~~Finch~~ | Shipped. /providers is the only thing distinguishing an unsupported field from an empty one, and almost nobody calls it; an assisted connection is a human logging in on a schedule; income is an amount and a unit and the unit is the meaning; a pay statement is a pay period, not a person |
| Nango | The proxy passes the provider's own errors through unchanged, so error handling has to know both vocabularies |
| Argyle | An account is connected before it is scanned, and the payroll history arrives over minutes |

### Payments outside the usual five

| Provider | Why |
|---|---|
| ~~Razorpay~~ | Shipped. Capture is a separate call and an uncaptured authorisation is auto-refunded, so doing nothing is a decision; amounts are in paise; the fee and tax do not exist until capture; and BAD_REQUEST_ERROR is the code on almost every failure, so reason is the only thing worth branching on |
| ~~Paystack~~ | Shipped, and written entirely from live responses: no credential, no account, every case dated. **The row's own premise is the part that could not be checked** -- every amount-bearing endpoint needs a secret key, so whether amounts really are in kobo is documented and unverified, and the Recipe models no amount field rather than claiming one. What was verified is a **negative** finding: every response carries `{status, message, data}` with a boolean beside the HTTP status, the shape that lets a body disagree with its own status line, and on the whole key-free surface the two always agree -- `true` with 200, `false` with 401, no 200 ever carrying `false`. Also pinned: a 401 whose machine-readable code is `invalid_Key`, with a capital K inside a snake_case value; an unknown country answering 200 with an empty array, indistinguishable from a real country it has no banks for and with the success flag saying `true` for both; malformed paging values (`perPage=abc`, `-1`, twenty digits) all accepted in silence; and cursor paging that exists only behind an opt-in `use_cursor=true`. Stated and not served: the nested `data.status` of `"failed"` its documentation describes on a verified transaction, which would be a third statement of one outcome and the likeliest to disagree |
| ~~Flutterwave~~ | Shipped, with four cases checked against the live API and four drafted from its SDK documentation, and meant to be read against Paystack's. **Both wrap everything in `{status, message, data}` and disagree about the first field**: Paystack sends the boolean `false`, Flutterwave the string `"error"`, so `if (!body.status)` treats every Flutterwave failure as a success. Flutterwave has no unauthenticated surface at all, where Paystack serves banks and countries to anyone. Documented and not called: the response shape changes by payment method rather than by endpoint -- a bank transfer has no `data` key at all, a Nigerian direct debit nests `meta.authorization` inside `data`, and a European one puts the same value at the top level as a sibling |
| Mercado Pago | A payment can sit in in_process for days, and the status detail rather than the status says why |
| dLocal | The settlement currency is not the charge currency and the rate is fixed at a moment you did not choose |
| ~~Rapyd~~ | Shipped. The id-namespace question could not be reached -- every endpoint that would answer it needs real HMAC credentials -- and what the probing found instead is sharper. Rapyd answers three things and **only two of them are about the credential**: a routing or method miss answers 401 `UNAUTHORIZED_API_CALL` with an empty message, identical whether credentials were absent, garbage, or well-formed, which is what proves routing resolves independently of authentication. And the signature's shape is never examined at all: `!!!not valid at all???` reaches the same tier as plausible hex, so a wrong signature and a credential that was never valid are one event. `auth.malformed_error` is deliberately unset, because nothing here distinguishes malformed from absent |
| ~~Worldpay~~ | Shipped, and the row's premise is confirmed from an angle it did not expect. The modern API dispatches on a **versioned vendor content type** -- `application/vnd.worldpay.payments-v7+json` against v6 -- so the media type decides which API you reached rather than the path. That mechanism is real, documented, and **proven unreachable without a credential**: Basic auth resolves before content negotiation for every credential state, verified across seven Content-Type and Accept combinations, all answering byte-identical 401s. A mechanism nobody can observe before authenticating is a mechanism nobody debugs before authenticating. The credential split itself is a genuine three-way one and all three are served |
| ~~FastSpring~~ | Shipped. One purchase is an order, a subscription and an account with three different ids and which one a webhook hands you depends on the event; FastSpring is the merchant of record so three amounts sit on one order and none is the price you set; a cancelled subscription is still active until its period ends |

### Notification and messaging infrastructure

| Provider | Why |
|---|---|
| ~~Xendit~~ | Shipped. A virtual account number exists in the 201 and does not work at the bank for minutes, so the customer is told it does not exist by their own bank; a closed account accepts one exact amount and a short payment bounces days later; amounts are integer rupiah with no subunit at all |
| ~~Midtrans~~ | Shipped. transaction_status and fraud_status are two fields and a payment is only safe when both agree; capture plus challenge means the card was charged and the funds are held; a bank transfer never passes through capture at all |
| PayU | The same merchant has different endpoints per country and the response fields differ between them |
| ~~Gorgias~~ | Shipped. A ticket and its messages are two paginated endpoints read at different moments, so the count on one disagrees with the array on the other; from_agent is true for automated replies; and a reopened ticket keeps its closing time with nothing marking the reopening |
| ~~Kustomer~~ | Shipped. A conversation is a customer timeline carrying every channel they ever used, assignment is on the conversation rather than any message, status and queue are unrelated, and everything is JSON:API so nothing a client wants is at the top level |
| ~~Jotform~~ | Shipped, and the premise was wrong in the useful direction: an answer embeds the question's own text and type inline, so a Jotform submission is the one in this group that survives its form being edited. Its not-found echoes the path parameter's name rather than the value |

### AI and inference

| Provider | Why |
|---|---|
| ~~Cohere~~ | Shipped. billed_units and tokens do not agree and pricing from the wrong one is wrong on every request; a finish reason is not an error and MAX_TOKENS is a 200 with a truncated body; an embedding carries nothing but its position |
| ~~Mistral~~ | Shipped, with four cases checked against the live API and five drafted from Mistral's own description and carrying no verified date -- the split is in the file. **It is offered as an OpenAI-compatible surface, and this Recipe marks where that stops**: the error envelope is flat `detail` where OpenAI nests under `error{}`, so `err.error.message` is undefined; completion ids are prefixed `cmpl-` rather than `chatcmpl-`; `AssistantMessage` has no `refusal` field at all, absent from the schema rather than null, so the sibling-field pattern OpenAI's own Recipe documents does not exist here; and `finish_reason` carries `model_length` and `error`, two values OpenAI has no equivalent for. Its unauthenticated failure does not distinguish a missing credential from a wrong one, answering byte-identically to both, and a malformed path answers a gateway's words rather than the API's -- `{"message": "no Route matched with those values"}`, with no `detail` key anywhere. `cauldron check` against its own description reports nothing contradicted |
| ~~Together AI~~ | Shipped, and the finding is one the row did not reach for. Together's own compatibility page tells you to point an OpenAI client at it and change two lines; the first call such a client makes is `GET /v1/models`, and **Together answers it with a bare JSON array** where OpenAI answers `{"object": "list", "data": [...]}`. That is Together's own published description, not an inference: `ModelInfoList` is `type: array` with no wrapper. So `response.data` is undefined on the first call, before any inference happens. Also pinned: absent and wrong credentials answer in **different media types** -- fifteen bytes of `text/plain` against JSON -- and a header with no `Bearer ` prefix is reported as missing, which is the exact opposite of how Fireworks classifies the byte-identical request |
| ~~Groq~~ | Shipped. The rate-limit headers this row asked about could not be reached without a key and are not modelled. What was verified: **the wrong verb and the wrong path are the same failure** -- a POST to a real GET-only path answers the identical 404 a path nobody defined answers, with the method folded into the sentence and no 405 anywhere. Routes also resolve before the credential, the opposite order from Scaleway's. No `upstream.spec`: Groq's description is embedded in a documentation page's client-side payload rather than served at a URL, so a fingerprint would be one drift could never recompute |
| ~~Meilisearch~~ | Shipped. A write answers 202 with a number and the word enqueued, the document is in neither the document listing nor the index until the task runs, and the task can fail after the 202 that accepted it |
| Hugging Face Inference | A cold model answers 503 with an estimated_time, and the correct behaviour is to wait rather than retry |
| ~~Langfuse~~ | Shipped. Traces are ingested asynchronously and are not readable immediately after being written |

### Data movement and warehousing

| Provider | Why |
|---|---|
| ~~Fivetran~~ | Shipped. Triggering a sync that is already running answers success and starts nothing, succeeded_at and failed_at both persist so which is later is the only health signal, and paused and sync_frequency are different things that each look fine alone |
| ~~Airbyte~~ | Shipped, and it answers the question this row was queued for: **a partial failure is not called a success**. `incomplete` is a status distinct from both succeeded and failed. Two things qualify it -- the name appears in the enum and nowhere in Airbyte's own prose about jobs, and once a job is incomplete the response gives only whole-job totals with no per-stream breakdown, so a caller knows something did not finish and cannot find out what. Census, written alongside, does the opposite: its own worked example shows a run with `records_failed: 1` carrying `status: "completed"`. Airbyte also checks the credential before routing without exception, which is exactly what this runtime does, so it needed no workaround -- the first time in that batch that was true |
| ~~Hightouch~~ | Shipped. A run with every row rejected finishes as success, warning is neither an error nor a success, the reasons live on an endpoint the run does not link to, and disabled and paused are different things that both stop a sync |
| dbt Cloud | A run has steps, and the run status is not the step status |
| ~~Snowflake SQL API~~ | Shipped. The same endpoint answers 200 with results or 202 with a handle depending on how fast the query was, every value is a string, a row is positional with the names in the metadata, and NULL is a real null inside the array of strings |
| ~~ClickHouse Cloud~~ | Shipped as `recipes/clickhouse`, and the row's question about paging belongs to the data plane, which this Recipe does not model. What the control plane gave instead: **its own root serves its own description** -- `GET /v1` answers 200 unauthenticated with 1.2MB of OpenAPI, to anybody -- so this Recipe's base URL and its `upstream.spec` are the same string. And a create nests the record **one level deeper than every other operation**, `result.service.*` against `result.*`, silently: both are 200, both are JSON, and the field a client wants is one hop further down. Every routing failure is Express's bare `Cannot {METHOD} {path}` in HTML, with no 405 anywhere |

### Incident response and observability

| Provider | Why |
|---|---|
| ~~incident.io~~ | Shipped. Status and severity are configured per workspace so neither is an enum you can hard-code, the category is the only fixed thing, the rank is the orderable thing rather than the name, and the two move independently |
| ~~Opsgenie~~ | Shipped. An alert and an incident have separate ids, lifecycles and close endpoints, so closing one leaves the other open; a create answers 202 with a request id that is not the alert; and a flapping monitor is one alert with a count rather than many alerts |
| ~~Checkly~~ | Shipped, with four cases checked against the live API and six drafted from its published description, which is fingerprinted. **One run disagrees with itself by location**: a check names several locations, every run produces one result each sharing a single `checkRunId`, and they differ independently on timing and on whether anything failed -- so "did the check pass" has no single answer and a client reading the first result gets a verdict true of one continent. Triggering is asynchronous, 201 with `status: RUNNING` and an empty results array, then polled. It does **not** share the property the Healthchecks Recipe documents -- Checkly originates every request itself -- except for its HEARTBEAT check type, which is exactly Healthchecks' shape and is named and not modelled |
| Rev.ai | A job is asynchronous and the transcript is a separate fetch with its own content types, so the job being done is not the transcript being readable |
| ~~Perplexity~~ | Shipped. The citations this row asked about are recorded from the description and not called: they arrive as a top-level sibling of `choices`, with no home in OpenAI's schema, so a client typed against it drops them silently. What was verified: **one route is versioned and the other is not** -- `POST /chat/completions` is a 401 and `POST /v1/chat/completions` is a 404 with zero bytes -- and `/v1/models` describes an unrelated product, so a client listing models to choose one for chat is reading the wrong catalogue |

### Auth, one more time

| Provider | Why |
|---|---|
| ~~PropelAuth~~ | Shipped, and the row's own claim is confirmed from PropelAuth's published worked examples: `org_id_to_org_info` really does arrive as `{}` for a user in no organisation, so an empty object is a valid user rather than a broken read. What the row did not anticipate is worse than what it did. **The status code tracks whether an Authorization header was sent, not what was in it** -- no header is 404, any header at all is 401, on a hostname invented for the probe and on PropelAuth's own quickstart host alike, with byte-identical empty bodies. So a typo in the subdomain and an expired key are the same event. That is the failure the Make Recipe records, repeated and widened |
| Helicone | It is a proxy, so a single call can fail as the upstream provider or as Helicone, and the two error shapes are unrelated |
| ~~Fireworks AI~~ | Shipped, and the row's first clause is confirmed in a stronger form than it guessed: a model's identifier on the control plane is a **resource name**, `accounts/my-account/models/my-model`, so it carries three slashes and the word "models" inside it. What the row missed is that there are two API designs on one hostname -- the OpenAI-compatible inference surface, and a control plane shaped by Google's API Improvement Proposals (`models`/`nextPageToken`/`totalSize`, sharing no key with the other half). Written as the pair to Together AI, and the two **split the middle credential case in opposite directions**: a header present with no scheme is "missing" to Together and "invalid" to Fireworks. Cold-start latency was not reachable without a key and is not claimed |
| ~~SurveyMonkey~~ | Shipped, and the sharper problem is upstream of that one: a response carries choice, row and column identifiers and never a label, so an answer cannot be read at all without joining to the survey's current questions. Its own published example gets its own wire wrong, quoting a numeric error id as a string |

### Storage and media

| Provider | Why |
|---|---|
| ~~Bunny.net~~ | Shipped as `recipes/bunny`. The purge question was overtaken by a stranger one: **`POST /videolibrary` with a zero-byte body and no credential answers 500**, with Bunny's own typos in the message, verified twice -- and attaching any body at all, still with no credential, reverts to an ordinary 401. So the unauthenticated request that fails hardest is the one that sends least. Also pinned: one vendor, two hosts, two entirely different backends and two incompatible error envelopes, each collapsing absent, malformed and wrong credentials into one message of its own -- so the split is by host rather than by mistake. Both check routing before credentials |
| ~~Transloadit~~ | Shipped, and the trap is sharper than the row guessed. It is not that partial results are readable -- it is that **the byte counts say finished while the status says executing**. Transloadit's own two published examples of one assembly show `bytes_received` equal to `bytes_expected`, 1687 of 1687, on a response still marked `ASSEMBLY_EXECUTING`; what actually moves between them is `execution_duration` and `results`. A client watching the bytes concludes it is done at the moment the work begins. Found while probing, and worth more than the row's own question: `GET /assemblies/{id}` answers **with no credential at all**, which Transloadit's documentation confirms is the design -- "This request can be issued by anyone who knows the assembly_ssl_url". The URL is the credential, so anything that logs one has published it |
| ~~Elastic Cloud~~ | Shipped, and the row's premise could not be reached -- no account exists and Elastic publishes no worked example of a create anywhere, so the create-returns question is answered in prose against Koyeb's and Northflank's rather than served. What the probing found instead is better: **attaching any credential hides the routing.** With no Authorization header, an unrouted path and a wrong method resolve first and answer properly -- `root.resource_not_found` and `root.method_not_allowed`, the latter with a real `Allow` header. Attach any credential at all, even garbage, and both collapse into a different, undocumented shape with no `errors` array and the method-versus-path distinction gone. Reproduced three times each way, so the more you look like a customer the less the API tells you about your own mistake |
| ~~Rootly~~ | Shipped. It is the first JSON:API provider here, and the outer envelope fits the way Lemon Squeezy's already does. What has no precedent is **a second complete JSON:API document nested as an attribute's value** -- `incident.attributes.severity` carries its own `data`/`id`/`type`/`attributes`, rather than a reference in `relationships` with the body in `included`, which is what the specification is for. Expressing it took every nesting primitive the format has, chained three deep. Its routing order is also **conditional**: the credential is checked first for every path Rootly has, and the router answers first for paths it does not, so the ordering depends on whether the route exists -- which a caller cannot know in advance, and which one gate before routing cannot reproduce |

## Listings that narrow themselves

A listing that quietly returns a subset is the worst-shaped failure there is:
nothing errors, and the list is correct. An emulator that returns everything
is being helpful in exactly the direction that hides it.

Route filters exist now, so this is a list of where they still need applying.
GitHub and Basecamp are done. The rest are known and open.

| Provider | The default, and why it is not applied yet |
|---|---|
| ~~GitHub~~ | Issues list open only. Applied |
| ~~Basecamp~~ | Projects list active only, and there is no value that asks for everything. Applied |
| Bitbucket | Pull requests list OPEN only. Applying it needs an open pull request in the fixture, which the fixture does not have, and it would rewrite the case that asserts three terminal states in one listing |
| ClickUp | Tasks exclude closed ones unless include_closed=true. The parameter is a boolean that inverts rather than a value that matches, which the filter cannot express yet |
| PagerDuty | Incidents default to triggered and acknowledged, excluding resolved. Two values in a repeated array parameter. The filter can carry a set of values now, so what is left is the repeated parameter itself |
| Sentry | Issues default to is:unresolved through a query DSL rather than a parameter, so the filter would have to parse a query language |
| Front | Conversations exclude archived in some views and not others, which needs checking before anything is declared |

Two shapes the filter still cannot describe are worth having: a boolean that
inverts, and a default expressed inside a query language. Neither should be
guessed at. The third, a parameter value covering several field values, exists
now: Alpaca needed it, because status=open means new or partially_filled and
matching the word literally would have hidden every partially filled order.

## Detection coverage

A Recipe that ships and that no dependency maps to still works: `cauldron up
<name>` and a config entry both reach it. What it does not do is turn up on its
own when somebody runs `cauldron detect` in a repository that uses it, which is
the thing the front of the README promises.

The table went from 12 Recipes to 91 in one pass, and from 91 to 147 in
another. These twenty-one are left, and every one of them has now been looked for
rather than remembered -- which is the whole rule: a package name written from
memory is exactly the guess detection forbids.

**Three names collided with something else entirely**, which is the reason the
rule exists and is worth recording so nobody maps them later on the strength
of the name:

| Name on npm | What it actually is |
|---|---|
| `basecamp` | A collection of Astro components for somebody's scout group |
| `mercury` | A modular frontend framework, not the bank |
| `persona-api` | A personalised RAG and LLM chatbot, not the identity provider |

Two more resolve and cannot be verified: `marqeta` carries no description and
no repository, and `gorgias` is a grunt plugin published by the company rather
than a client for its API.

| Recipe | Note |
|---|---|
| Attio | The npm package is a CLI for building apps on the platform rather than a client that calls the API |
| Basecamp | No official SDK, and the name is taken by something unrelated |
| Bill.com | No widely used client |
| Column | No official SDK |
| Deel | No official SDK found under any obvious name |
| FastSpring | No official SDK |
| Fivetran | No official Node or Go client found |
| Gorgias | The package of that name is a grunt plugin, not an API client |
| Greenhouse | Community clients only |
| Hightouch | No official client found |
| incident.io | No official Node or Go client found |
| Kustomer | Community clients |
| Marqeta | A package of that name exists with no description and no repository, so it cannot be verified |
| Mercury | No official SDK, and the name is taken by a frontend framework |
| Persona | The name is taken by a chatbot library |
| Ramp | A package of that name exists with no repository and no description beyond its own name |
| RevenueCat | The published packages are the mobile and hybrid SDKs rather than a server client |
| ShipHero | No official client; the API is GraphQL over raw HTTP |
| Where the ISS at | Nothing on either registry calls it -- npm and Packagist both answer "wheretheiss" with nothing at all. What the neighbourhood does have is everything that *reads* a two-line element set without fetching one from here: `tle.js`, `satellite.js`, `ootk-core` and `satellite-tracker` propagate orbits from the format this API serves, and the two that name a source fetch it from `celestrak.com` and `space-track.org` instead. `solarpatrol/yii2-tle` on Packagist is TLE storage for space-track too. The rest match on the word: `@tuwaio/satellite-react` and `@tuwaio/satellite-solana` are Solana wallet adapters, "satellite" being a layer name in somebody's architecture, and `@pipeworx/mcp-tle` is an MCP server carrying no host at all |
| FBI Wanted | Nothing on either registry calls it, and both words lead somewhere else. "Wanted" is a term of art in package managers -- npm's *wanted* version is the one a semver range prefers -- so the search returns `@pnpm/parse-wanted-dependency`, `npm-install-checks` and `path-scurry`. "FBI" on Packagist is Firebird: `satwareag/php-firebird-stubs`, `satag/doctrine-firebird-driver` and `satwareag/pdo-fbird`, whose PHP extension is spelled `fbird`. Beside them, `smarthub/fbinsights` and `mcstutterfish/fbia-rss` are Facebook Instant Articles, `fbi/apk-parser` is an Android package reader under that vendor prefix, and the bare npm name `fbi` is "a workflow tool in the command line". The one adjacent client is for a different API of the same agency: `@cyanheads/fbi-crime-mcp-server` exposes the Crime Data Explorer, not this |
| Spaceflight News | Nothing on either registry calls it. Packagist has no result for "spaceflight" at all, and npm's one candidate, `spaceflightnewsapi-node` -- "lets you easily consume content from the Spaceflight News API" -- carries no host in its published output. Two MCP servers carry the name and are not mapped. The rest of what npm returns is the popularity fallback, the habit Packagist showed for "datamuse": `nan`, `bullmq`, `draco3d`, `cytoscape` and `@plotly/d3`, none of which mentions space |
| arXiv | Nothing on either registry calls this API under a name that identifies it. The npm packages naming arXiv scrape the website's HTML or parse LaTeX sources, and Packagist returns bibliography formatters that read files already on disk. The clients that do call the API are Python, which neither registry indexes. Offering this emulator to a LaTeX parser would be wrong about the protocol entirely, so it ships unmapped and a test guards the decision |
| Healthchecks | One npm client exists, `healthchecks-io-client`, and it targets the **management API** at healthchecks.io -- creating and listing checks, all of it authenticated. This Recipe models that surface only as its unauthenticated 401, because no account was created to observe anything behind it. What the Recipe does serve is the ping host, which that client barely touches. Mapping the two would offer an emulator that answers 401 to every call the client makes, so it ships unmapped |
| NASA | Nothing on either registry is a client of api.nasa.gov. The npm packages bearing the name wrap the Astronomy Picture of the Day for hobby projects and are depended on by nothing; Packagist returns image-gallery libraries. This Recipe also covers a gateway fronting several unrelated applications, so a name match would not even identify which of them a caller meant, and it ships unmapped |
| Steam | The packages naming Steam on npm are overwhelmingly clients of **Steamworks**, the C++ game SDK, or of the community trading and market pages -- different hosts, different protocols, and in the SDK's case not an HTTP API at all. Packagist returns OpenID login helpers for signing in with Steam, which touch the auth flow and never the Web API. Offering this emulator to any of them would be wrong about the host, so it ships unmapped |
| NOAA | Nothing on either registry reaches the NCEI access service. The clients that exist are Python and R, which neither registry indexes, and the npm results for the name reach **api.weather.gov** -- the National Weather Service, a different NOAA API that this collection already describes under its own Recipe. Mapping this one to those packages would hand somebody the wrong NOAA, so it ships unmapped |
| Guardian | The npm packages naming the Guardian are almost all unrelated: route guards for web frameworks, which share the word and nothing else. The few genuine Open Platform clients are unmaintained and depended on by nothing. Packagist returns authorisation libraries for the same reason. A name match here is a homograph rather than a miss, so it ships unmapped |
| New Relic | The `newrelic` package on npm is the **APM agent**, not a client of this API. It instruments an application and reports telemetry to New Relic's ingest endpoints; it never calls the v2 REST API or NerdGraph that this Recipe describes, and `@newrelic/telemetry-sdk` is the same story with a smaller surface. A project holding either is being monitored by New Relic rather than querying it, so offering this emulator would intercept nothing it actually sends. It ships unmapped and a test guards the decision |
| NewsAPI | Nothing on either registry reaches it under a name that identifies it. The npm packages calling themselves NewsAPI clients are thin unmaintained wrappers with nothing depending on them, and Packagist returns unrelated content-management libraries. Offering this emulator on a name match would be wrong about the host and the key at once, so it ships unmapped and a test guards the decision |
| GNews | The same, and worse for the name: `gnews` on npm is a Google News **scraper** that parses the public RSS feed and never calls gnews.io at all -- a different service reached a different way. Mapping the two would hand somebody an emulator of an API their code does not use, so it ships unmapped |
| Dynatrace | `dynatrace` on npm says so itself: "This package is no longer supported and has been deprecated." Nothing has replaced it under any scope, which fits this Recipe's own finding -- the API lives behind a per-tenant hostname that answers an AWS gateway default to anyone without an account, so there is little for a public client to be written against |
| Scout APM | Nothing on either registry. Scout publishes a real, fetchable OpenAPI document and no client generated from it, which is the same pairing the Kit row records |
| Nexus | A homograph: `nexus` on npm is a **GraphQL schema-development framework**, unrelated to Sonatype's repository manager. Nothing else under any spelling |
| Snipcart | The second job-queue homograph in this collection, after Workable: `snipcart` on npm is a "Node.js job queue as a stream", version 0.0.1. Fitting, in a way -- this Recipe's own finding is that Snipcart's API creates nothing, so the name on the registry does more writing than the product does |
| Formstack | Nothing under that name, and the name is the problem the Recipe already documents: Formstack Forms and Formstack Documents are different products on different hosts, and a package called `formstack` would be ambiguous between them even if one existed |
| Sendcloud | A homograph, and a confusing one: `sendcloud` on npm is "send email quickly and simply" -- SendCloud the email service, not Sendcloud the shipping company. Two real products with one name, and the one on the registry is not this one |
| ShipBob | Nothing under any spelling, which is no longer surprising: this Recipe's own finding is that **the v1.0 host has been switched off**, answering 410 Gone to every credentialed request since 29 August 2026. A client library for it would have nothing left to call |
| ClickHouse | `@clickhouse/client` and `clickhouse` are both **data-plane** clients: they execute SQL against a database. This Recipe describes api.clickhouse.cloud, the control plane -- creating and sizing services, which is how you learn what host to point a SQL client at. The Milvus situation exactly: mapping either would offer an emulator of the address book to code that already has the address |
| Timescale | `@timescale/client` does not exist, and the company is Tiger Data now -- its API moved to console.cloud.tigerdata.com while console.timescale.com went dead and api.timescale.com became a load balancer that refuses TLS. Nothing on either registry has followed it |
| InfluxDB | `@influxdata/influxdb-client` is the **2.x data-plane** client -- write and query against a bucket. This Recipe describes the one endpoint on InfluxDB Cloud that is genuinely a credential-checked API, and the Recipe's own finding is that the data plane is not served on that host at all: `/api/v2/buckets` redirects to a login page and the documented write example 404s. A client of a surface that is not there |
| NetSuite | `netsuite` on npm wraps the **SuiteTalk SOAP API**. This Recipe describes the REST record service, which is a different protocol against a different endpoint -- the same wrong-surface miss the Milvus row records, one layer down. Nothing on either registry reaches the REST half, which follows from the Recipe's own finding that there is no address to reach without a provisioned account |
| Sage Intacct | Nothing under any spelling tried. Unsurprising for an API that is **XML over HTTP to one endpoint**: there are no resources to give a client library a shape, and integrations are written against the XML schema directly |
| Porter | `porter` on npm is a binary-streaming service, unrelated. Nothing exists under the vendor's own scope either, which follows from what the Recipe found: **there is no documented HTTP API to write a client for.** `api.porter.run` is dead, the documentation publishes no REST reference at all, and the only live surface is the dashboard's own backend. A client library would have to be a client of somebody's front end |
| Rippling | A homograph with a repository to prove it. `rippling` on npm is `e7h4n/rippling`, described as "Rippling Core" and unrelated to the HR company -- and no official client exists under any scope. Consistent with this Recipe's own finding, that the vendor will not show a schema without an account |
| Workday | Nothing on either registry, under any spelling. Workday integrations are built inside Workday, or through middleware, and the REST gateway this Recipe describes is reached by code that holds a tenant's OAuth credentials rather than by a published library |
| Milvus | The official clients speak **gRPC** to a Milvus cluster's data plane. `@zilliz/milvus2-sdk-node` and `milvus-sdk-go` both do, and this Recipe describes Zilliz Cloud's control-plane REST API -- listing and describing clusters, which is how you find out what host to point a gRPC client at in the first place. Mapping either would offer an emulator of the thing that hands out addresses to code that already has one |
| Baserow | `baserow` on npm is a **security holding package**, version 0.0.1-security: a reserved name with no code in it, the same miss the Fireworks row records. Nothing under any scope, and nothing on Packagist either -- which is consistent with Baserow being self-hosted as often as not, where the caller is a script rather than a dependency |
| Retool | Nothing under any spelling tried, on either registry. Consistent with this Recipe's own finding: the public API is administrative -- users, groups, permissions, folders -- and the code that calls it is somebody's onboarding script rather than a library anybody publishes |
| Payoneer | Nothing on either registry, which follows from what this Recipe found: developer.payoneer.com answers a bot challenge and its reference mirror renders client-side. The code that calls this API is written from a GitHub sample rather than installed |
| Melio | Nothing, and only a partner-issued token authenticates anything -- so there is no population of callers large enough to publish a client for |
| Treasury Prime | Nothing. Banking-as-a-service integrations are written against a specific bank programme rather than as a general-purpose library |
| Tipalti | `tipalti` on npm is a **URL helper** -- it composes iframe links and never calls the API this Recipe describes. The right shape of miss for a product whose front end is a hosted iframe |
| Anrok | Nothing under any spelling. Anrok publishes a real fetchable OpenAPI document and nothing generated from it, the pairing the Kit and Scout APM rows already record |
| Vertex | Nothing, and there would be nowhere to point it: this Recipe's own finding is that **every hostname in Vertex's own spec is NXDOMAIN** |
| TriNet | Nothing, and no public documentation exists either -- `apidocs.trinet.com`, the address TriNet's own site links to, does not resolve |
| Remote.com | Nothing under `remote`, `remotecom` or any scoped spelling. The word is too common to have been claimed for this |
| Kaleyra | Nothing. Kaleyra was acquired and its documentation moved to a ReadMe portal; no client followed |
| Textline | Nothing, which fits a product whose API is an afterthought to a web application -- this Recipe found it has no delivery-status vocabulary at all |
| Constant Contact | `constantcontact` on npm wraps **v2**. This Recipe describes v3, which is a different API with a different credential flow -- the Worldpay miss again, where a package resolves and targets a retired version |
| Omnisend | Nothing under any spelling, on either registry |
| Sardine | Nothing, which is consistent with the Recipe's own finding: **Sardine's field-level API reference is sales-gated**, so there is no public population of integrators to write a client for |
| MinIO | `minio` on npm is the **S3 client**. This Recipe models the admin API at `/minio/admin/*`, deliberately -- see its header on why the S3 surface could not be probed on the demo host. Two protocols, one vendor, and the package speaks the other one |
| Tigris | `@tigrisdata/core` is the client for Tigris's **former** product, a serverless document database. This Recipe describes the S3-compatible object storage the company pivoted to, whose client is an AWS SDK pointed elsewhere -- the same reason Backblaze, Wasabi and Filebase are all unreachable |
| Elastic Cloud | `@elastic/elasticsearch` is the **search** client: it queries an index. This Recipe describes the control plane that creates the deployment the index lives in, which is the ClickHouse and Milvus situation for the third time |
| Optimizely | `@optimizely/optimizely-sdk` evaluates experiments **on the client**, which is this Recipe's own central finding: there is no decide endpoint, assignment is computed locally from a datafile. The SDK never calls the account API this Recipe models |
| VWO | `vwo-node-sdk` is the same shape of miss for the same reason -- it buckets locally with MurmurHash3 and never calls the management API |
| Flexport | `flexport` on npm is version 1.0.0 with an empty description and no repository. Not enough to call it a client of anything |
| Bringg | Nothing, and this Recipe explains why nobody could write one: every real operation lives at a path of per-merchant UUIDs that are undiscoverable without an account |
| Radar | `radar-sdk-js` is the **browser** SDK and holds a publishable key. This Recipe models the server API, whose secret-key routes that SDK is specifically not allowed to reach -- the AppSignal and FullStory miss, in geolocation |
| TalkJS | `talkjs` on npm is the **embed** SDK: it renders a chat UI in a page. This Recipe models the REST API a server calls to provision the conversations that UI shows |
| CometChat | `@cometchat/chat-sdk-javascript` is likewise the client-side chat SDK. The REST API this Recipe describes is called by a customer's backend, and nothing on either registry represents that |
| Whereby | `@whereby.com/browser-sdk` embeds a video room in a page. This Recipe models the API that creates the room -- and its own finding is that the room URL is unsigned because the room itself expires, which is a server-side concern the embed never sees |
| Clio | A homograph, and a complete one: `clio` on npm is the **Clio programming language** |
| Aha! | Another: `aha` on npm is an ANSI-to-HTML adapter |
| Missive | And a third: `missive` on npm is a binary encoding library. Three unrelated homographs in one batch |
| JustCall | Nothing under any spelling. Its own OpenAPI contradicts itself about whether a list returns an array, which is not a document anyone would generate a client from |
| CloudTalk | Nothing, which fits a vendor whose two hosts disagree with each other about credential granularity and routing order |
| Adapty | Nothing on npm. Its published SDKs are for iOS, Android and React Native -- the platforms where a paywall lives -- and the server API this Recipe models is called from a backend that installs none of them |
| Appsmith | `appsmith` on npm is a **security holding package**, version 0.0.1-security -- a reserved name with no code, the third in this collection after Fireworks and Baserow. Fitting for a product whose public API this Recipe found to be administrative only |
| Budibase | `budibase` on npm is the **CLI**: it installs and runs Budibase. This Recipe describes the public API a caller reaches from outside an instance, which is the same distinction the n8n and NocoDB rows draw |
| SavvyCal | Nothing under any spelling, on either registry. Its API is small and recent, and the integrations that use it are written directly |
| Ayrshare | Nothing, which is unsurprising for a product whose whole value is one call to several networks -- there is little for a client library to abstract |
| Postscript | Nothing. Its callers are Shopify apps, which reach it through Shopify's own tooling rather than through a published client |
| Papaya Global | `papaya` on npm is a **dependency-injection container**, unrelated. Nothing exists under the company's own name either, which fits this Recipe's finding that the only reachable API is a payment-disbursement product rather than the employment one the name suggests |
| Zapier | `zapier-platform-core` is the SDK for **building** an integration that Zapier runs, not for calling Zapier. Same shape as the Pipedream row: a project holding it is a thing Zapier invokes rather than a thing that invokes Zapier, so the emulator would intercept nothing |
| n8n | Installing `n8n` means **running n8n**, not calling it, and `n8n-workflow` is its own base code. The public API this Recipe describes is called by the operator of an n8n instance, from outside it, and nothing on either registry represents that |
| Worldpay | `worldpay` on npm was last published in **2014**, at version 0.0.4, against the old worldpay.com endpoints. This Recipe describes Worldpay Access (`access.worldpay.com`) with its versioned vendor content types, which did not exist then. A package that resolves and targets a retired API is worse than none: it would map a modern integration to an emulator of something else |
| Rapyd | `rapyd` on npm is a single keyword -- "connector" -- with no repository, no description beyond that word, and nothing published since **2019**. Not enough to call it a client of anything, and the Recipe's own finding is that every reachable Rapyd endpoint needs real HMAC credentials, so there is little for a client to do unauthenticated anyway |
| Kit | Three misses in one word. `kit` on npm is a TypeScript utility library with Russian pluralization helpers; `convertkit`, the old brand, is version 0.0.0 with a one-word description and nothing in it; and no official SDK exists under any name on either registry. Kit publishes a real OpenAPI document and no client generated from it |
| Lever | `lever` on npm is a Chinese frontend continuous-integration tool. No official client exists on either registry, which is consistent with the Recipe's own finding that Lever publishes prose documentation and no machine-readable description |
| Ashby | Nothing at all, under any spelling tried. Ashby's API is POST-only and documented in prose; the integrations that call it are internal to customers rather than published |
| Workable | The best homograph in this collection: `workable` on npm is a **MongoDB-based job queue**. A job queue and a job board, one word. No official client exists either |
| Wasabi | A homograph, and the npm one is not close. `wasabi` is a realtime game-networking and general-purpose replication library; this Recipe is object storage. Nothing on either registry calls the S3 endpoint or the account-control API, which is unsurprising -- code reaching Wasabi's S3 gateway holds an AWS SDK and a different endpoint URL, and a dependency list cannot see an endpoint |
| Filebase | The same miss twice over. `filebase` on npm is git-backed file storage with i18n support, unrelated to the S3-on-IPFS service this Recipe describes, and here too a real client would be an AWS SDK pointed at another host. The three storage Recipes written together share this: **S3 compatibility is what makes them unreachable**, because the compatible client is somebody else's |
| Airbrake | A client exists and it is the wrong half. `airbrake-js` and `@airbrake/node` are **notifiers** -- they report exceptions to the ingestion endpoint, and this Recipe models the read API where you go looking for them afterwards. A project holding one is being watched by Airbrake rather than querying it, which is the same shape the AppSignal and New Relic rows record, now three times in one collection |
| Sumo Logic | Also split by half, also on the wrong side. The `sumologic` package syncs application logs to the collector endpoint. This Recipe describes the search-job API -- create a job, poll it, read the messages -- which is the other end of the same pipe, and nothing on either registry reaches it |
| Fireworks | `fireworks-ai` on npm is a **security holding package**, version 0.0.1-security: a name reserved so nobody else can take it, with no code in it. The official client is Python, which neither registry indexes. So the name resolves, the package installs, and mapping it would point the detector at an empty box |
| AppSignal | A client exists and it is the wrong one, the same shape of miss the New Relic row records. `@appsignal/nodejs` is the **instrumentation library**: it reports telemetry to `appsignal-endpoint.net`, the surface this Recipe deliberately does not model, having modelled the personal-token REST API at `appsignal.com/api` instead. A project holding it is being observed by AppSignal rather than querying it, so the emulator would intercept nothing it sends |
| SingleStore | The npm packages under the vendor's own scope speak **SQL**, not this API. `@singlestore/http-client` executes statements against a database; `@singlestore/client` wraps it. This Recipe describes the management plane -- creating and sizing workspaces and clusters -- which neither package touches. Mapping them would offer an emulator of the control plane to code that only ever runs queries |
| Pipedream | `pipedream` on npm is the **components monorepo**, the library of integration steps that run inside a Pipedream workflow. It is not a client of the REST API this Recipe describes, and a project holding it is being run by Pipedream rather than calling it. Nothing on Packagist is relevant either |
| Airbyte | Nothing on either registry calls this API. `faros-airbyte-cdk` is a **connector development kit** -- for building sources and destinations that Airbyte runs, not for driving Airbyte itself -- and the embedded widget is a UI component. The clients that do call the API are Python and Java, which neither registry indexes |
| Census | Nothing on either registry reaches it, and the name is unhelpful: the npm results for `census` are United States demographic-data wrappers, which speak to the Census Bureau API this collection already describes under `uscensus`. Two unrelated things sharing a word, and mapping either to the other would be wrong twice |
| Browserless | A homograph, and a costly one. The `browserless` package on npm is **microlinkhq's local Puppeteer wrapper** -- it drives a browser inside your own process and never calls browserless.io at all. A project holding it has no relationship with this service whatsoever, so mapping the two would offer an emulator of an API that project does not use. Packagist returns nothing relevant either. It ships unmapped and a test guards the decision |
| Tink | A homograph rather than a gap. `tink` on npm is a next-generation runtime and package manager that shares the name and nothing else; Packagist returns unrelated packages for the same reason. The clients that genuinely reach this API are Java and Python, which neither registry indexes. Offering an open-banking emulator to a project that installed a package manager would be wrong about everything at once, so it ships unmapped and a test guards the decision |
| Semantic Scholar | Nothing on either registry reaches the Academic Graph API. The clients that exist for it are Python, which neither npm nor Packagist indexes, and the npm results for the name are citation-formatting libraries that parse BibTeX already on disk rather than calling the service. Offering this emulator to a bibliography formatter would be wrong about the host and the protocol at once, so it ships unmapped and a test guards the decision |
| Mercado Libre | The npm and Packagist packages that name Mercado Libre are almost all clients of **Mercado Pago**, the payments product on a different host with a different credential, which this Recipe does not describe. The few that do target the marketplace API are unmaintained wrappers with nothing depending on them. Offering this emulator to a Mercado Pago integration would be wrong about the host, the credential and the vocabulary at once, so it ships unmapped |
| US Census | Nothing on either registry calls the Census data API. The packages that name it are R and Python, which neither registry indexes, and the npm results are unrelated geographic-boundary datasets shipped as static files rather than clients of anything. This Recipe also models only the metadata endpoints, since the data API is behind a key gate, so a name match would be doubly wrong -- and it ships unmapped |
| NCBI Entrez | Nothing on either registry reaches E-utilities. The packages that do exist for it are Python and R, which neither npm nor Packagist indexes, and the npm results for the name are unrelated bioinformatics file parsers that read FASTA and GenBank files already on disk rather than calling the service. Offering this emulator to a project parsing sequence files would be wrong about the host and the protocol at once, so it ships unmapped and a test guards the decision |
| Alpha Vantage | Nothing on either registry reaches this API under a name that identifies it. The npm packages calling themselves Alpha Vantage clients are thin wrappers republished many times over, none of them depended on by anything, and Packagist returns unrelated finance libraries that speak to other vendors entirely. Offering this emulator on the strength of a name match would be wrong about the host and the key at once, so it ships unmapped and a test guards the decision |
| SEC EDGAR | Nothing on either registry is a client of this API. The packages that name EDGAR parse **the filings**, not the API: they take a document that has already been fetched and pull XBRL facts out of it, so they make no request to data.sec.gov at all. What clients do exist for the API are Python, which neither registry indexes. Offering this emulator to a project that parses filing text would be wrong about the host, the protocol and the credential at once -- there is no credential -- so it ships unmapped and a test guards the decision |
| RCSB | Every client on npm speaks something else, and the vendor's own is the clearest case. `@rcsb/rcsb-api-tools` calls itself the "RCSB PDB API request library" and reaches `data.rcsb.org/graphql` and nothing else -- there is no `/rest/v1` anywhere in its published output, so the organisation's own client for its own host does not speak the REST Data API this Recipe describes. `@rcsb/rcsb-search-tools` reaches `search.rcsb.org`, **the vendor's other API**. `@rcsb/rcsb-molstar` reaches `data.rcsb.org/graphql` and `files.rcsb.org`, so it is on the right host by the other protocol. `pdbmine`, described as "a Node.js API for the RCSB Protein Data Bank", reaches `files.rcsb.org/view` and `www.rcsb.org/pdb` -- the file server and the website. `parse-pdb`, `pdb-parser-js` and `@i-vresse/pdbtbx-ts` parse **the file format** and make no request at all. `@pipeworx/mcp-pdbe` is **the same data from another organisation**, PDBe at the EBI. Three more are MCP servers, and the six `@rcsb/rcsb-saguaro*`, `rcsb-charts`, `rcsb-statistics` and `rcsb-pecos-app` packages are the vendor's viewers. Packagist returns SMS gateways |
| OpenAIRE | Nothing on either registry calls it, and Packagist's answer is the edit-distance near miss at a scale nothing else here has reached: **seven of its eight results are `openfire` or `openair`**. Five are clients for the Openfire XMPP server -- `gidkom/php-openfire-restapi`, `gnello/php-openfire-restapi`, `tuemmlerkon/php-openfire-restapi`, `tailonperin/openfire-restapi` and `pcabreus/openfire-userservice-plugin` -- which differs from "openaire" by swapping one letter, and two are `openair-xml-php`, for the Openair professional-services product, which differs by dropping one. The eighth, `acdh-oeaw/arche-openaire`, is an "OpenAIRE usage tracking plugin", which is **the vendor's other API**: the usage-statistics service rather than the search one. On npm the only candidate is `@pipeworx/mcp-openaire`, an MCP server, beside `doiget-cli` and `lit-search`, neither of which mentions api.openaire.eu anywhere in its published output, and a Zotero plugin for Volto |
| Nobel Prize | The only client on npm reaches the version this Recipe's own records point at, which is not the version that served them. `nobelprizeparser` -- "Nobel Prize data parser" -- reaches `api.nobelprize.org/2.0/laureates` and `api.nobelprize.org/2/laureate/102`, and this Recipe describes 2.1, so it is left out for the reason the MBTA's and ClinicalTrials.gov's retired-interface clients were. That the ecosystem speaks 2 is the same finding the Recipe pins from the other side: a 2.1 record's `links[0].href` points at 2. The rest of npm is **the data shipped as a file** twice over -- `nobel-prize-data`, "a JSON dataset about Nobel Prizes and Laureates", and `nobeldb-data`, which is a different Nobel database entirely -- beside two MCP servers, **a different sense of "prize"** in `@everymatrix/prize-engine-suite` and PoolTogether's lottery utilities, **a package about somebody who won one** in `hangang`, and a joke: `buttt`, "strings 'butts' together. Nobel prize here I come." Packagist's answer is **an author called Nobel** -- `nobelatunje/jwt`, `nobelatunje/wallet`, `nobelatunje/laravel-validator` and `nobelzsushank/nepali-date-converter` -- beside `noble-commerce`, `noeldemartin` and `nomelodic` |
| Wikimedia | Nothing on either registry is a client of the Core REST API, and the organisation's own namespaces are entirely internal engineering. Packagist's eight `wikimedia/*` results are a Composer plugin, a PHP port of LESS, a CSS sanitizer, a warning suppressor, an Aho-Corasick matcher, Unicode normalisation, a session serialiser and Parsoid -- the foundation's own libraries, published under the name a client would search for, and not one of them a client of anything. The `@wikimedia` npm scope is the same story: the Codex design system, its tokens and icons, three linting configurations, and -- pointedly -- `@wikimedia/wikimedia-spectral-ruleset`, a **ruleset for linting OpenAPI documents**, published by the organisation whose API this describes. The nearest real candidate is `wtf-plugin-api`, "a wikimedia api plugin for wtf_wikipedia", which reaches `/w/api.php` -- **the MediaWiki Action API**, a third interface again, neither this one nor the `rest_v1` this one supersedes |
| Toast | No client for the POS API on either registry, and the name means a notification popup: tall-toasts alone has over four hundred thousand installs. A mapping on the word would offer a restaurant point-of-sale emulator to every Laravel project that shows a toast when a form saves, so a test guards the decision |
| Royal Mail | Royal Mail runs three separate APIs on three hosts, and every client package on Packagist targets one of the other two: elliotjreed/royal-mail-tracking is the tracking service, and turtledesign/royalmail-php, zvps/royal-mail-shipping-rest-api-client and mobi-market/royalmail-shipping-v3 are all the Shipping API. Click & Drop has no client of its own, and offering its emulator to a project using a different Royal Mail API would be wrong about the host, the credential and the vocabulary at once, so a test guards the decision |
| Homebrew | Nothing calls it either. npm's results for "homebrew formula" are formula *parsers* -- `nth-check` at 260 million and `@sideway/formula` at 55 million, matched on the word -- beside two plugins that edit a formula file in a tap repository rather than calling the API. Packagist's is the hobby: `georgeh/php-beerxml` parses beer recipes, and `munkireport/homebrew` inventories Homebrew on managed Macs by running `brew` locally |
| Repology | Nothing on either registry calls it. npm has two MCP servers and Packagist's top result for the word is `pragmarx/countries` at 3.5 million downloads, a countries and currencies library -- so the empty result and the wrong result look identical, which is the reason to write down that it is empty |
| Tradier | A third-party SDK exists and is one person's; worth a look rather than a mapping on sight |
| Twilio Verify | Part of the main Twilio SDK, which maps to the twilio Recipe rather than this one |
| UPS | The dominant PHP client, `gabrielbull/ups-api` at two and a half million downloads, and both npm candidates (`ups-nodejs-sdk`, `node-ups`) all take an AccessKey, a UserId and a Password -- the legacy XML API this Recipe's own header says it deliberately does not model. The one package built against the OAuth2 REST API, `rahul-godiyal/php-ups-api-wrapper`, is a single person's project with seven stars |
| WeatherAPI | The bare name on npm, `weatherapi`, is a decade-old wrapper for a different vendor entirely, api.worldweatheronline.com. WeatherAPI.com publishes no SDK of its own, and the one Laravel package that names it correctly, `grigorygerasimov/laravel-weather`, has zero stars and twenty-six hundred downloads total |
| Tomorrow.io | Nothing dominant on either registry. The one composer package that genuinely targets it, `php-weather/tomorrow`, a provider inside the php-weather framework, has thirteen downloads total; npm's results are an n8n community node and an MCP wrapper, neither a library a project calls directly |
| Visual Crossing | Nothing on Packagist calls it -- the word "visual" returns image-optimisation SaaS clients instead. npm's one candidate, `@essamonline/weather-visualcrossing`, is a single version published by one person with eight downloads a month |
| Geoapify | The official npm package, `@geoapify/geocoder-autocomplete`, is a browser address-autocomplete input, not a server client -- and this Recipe models the IP Geolocation API, a different product on the same platform. Nothing on either registry reaches that endpoint specifically |
| Clearbit | `clearbit` on npm is Clearbit's own client and says so itself: "This package is no longer being maintained." The same finding as Dynatrace's row -- nothing has replaced it under any scope, which fits a company HubSpot acquired in 2023 and has been folding into its own product since |
| Apollo.io | The name is claimed many times over by unrelated things -- Apollo GraphQL, Ctrip's "阿波罗" configuration-management client, an Aliyun OSS adapter -- and every genuine npm result for Apollo.io itself is an MCP server wrapping the API for an agent, not a library a backend calls directly. The one composer package that targets the real REST API, `tommyoneill/apollo-api-client-php`, has six downloads total |
| Daily | The bare name on npm, `daily`, is an unrelated LevelDB-based logging system. Daily's own `@daily-co/daily-js` is the front-end call-object SDK that renders a video call in a page, not a client of the REST API this Recipe models -- the Whereby and TalkJS miss again. Three composer packages target the REST API directly (`steadfastcollective/laravel-dailyco`, `axenso/laravel-dailyco`, `abdelhamiderrahmouni/laravel-dailyco`), none of them with more than five hundred downloads or five stars |
| Agora | Every real candidate is the wrong half of the product. `agora-access-token` mints the dynamic keys client SDKs use to join a call; `agora-rtm-sdk` is the real-time-messaging client itself; the composer packages generate the same join tokens for Laravel. The one REST wrapper found, `agora-rest-client` on npm, covers only Cloud Recording, not the channel-management surface -- `/dev/v1/projects`, `/dev/v1/channel/{appid}` -- this Recipe models |
| CockroachDB | Nothing under the vendor's own GitHub org, and nothing found elsewhere, reaches the Cloud API this Recipe models. CockroachDB publishes the OpenAPI document (cockroachlabs.cloud/assets/docs/api/latest/openapi.json) and no client generated from it, the same pairing the Kit and Anrok rows already record |
| Redis Cloud | What exists under RedisLabs' own org are sample apps from the old Heroku add-on days -- `rediscloud-node-sample`, `rediscloud-php-sample` -- demonstrating a data-plane connection to a provisioned database, not a client of the account-management REST API (x-api-key, api.redislabs.com/v1) this Recipe models. `redis` and `ioredis` are the data-plane client the samples point at, the SingleStore miss again |
| SurrealDB | The official `surrealdb` package on npm is the database driver -- it runs queries against a SurrealDB instance, the ClickHouse and Milvus situation again. This Recipe models the Cloud control plane at api.surrealdb.com, and the organisation's own GitHub publishes only the OpenAPI document for it, `surrealdb/openapi`, with nothing generated from it |
| Alloy | The bare name is claimed by somebody else before it is claimed by this. `alloy` on npm is the TiDev Titanium MVC framework, and `eslint-config-alloy` is AlloyTeam's own lint rules -- two more unrelated projects carrying the word. Grafana ships an open-source project under the identical name that is not published to either registry as a package at all, which would still mislead on reputation alone. `alloy-frontend` and `alloy-node` are Alloy Automation, a workflow-integration company, not the identity-risk one this Recipe describes. The one candidate that genuinely belongs to Alloy Technologies, `@alloyidentity/web-sdk`, is a document-capture SDK a browser runs, not a caller of the Evaluations API this Recipe models. Packagist has nothing under the name at all. It ships unmapped |
| Taddy | Two unrelated products hold the name and neither is this one. `taddy` on npm is `lttb/taddy`, a compile-time atomic CSS-in-JS library, and `@taddy/core` and `@taddy/babel-plugin` are the rest of that build toolchain. The sharper collision is a second company: `taddy-sdk-web` on npm and `taddy/taddy-sdk` on Packagist are published by taddy.pro, a Telegram bot advertising and monetisation platform, and they are initialised with a Telegram bot token. A vendor prefix, the word SDK and a matching brand, and nothing in either repository mentions `api.taddy.org`, GraphQL, `X-API-KEY` or `X-USER-ID`. The only package that reaches the real podcast API is `n8n-nodes-taddy`, a community node that runs inside n8n rather than a client a project calls, which is where the Tomorrow.io row already drew the line. It ships unmapped |
| WattTime | Every client that exists is built against the previous API, on a host this Recipe does not model. `IuryAlves/watttime-go`, `rekuberate-io/carbon`, `siderolabs/kube-scheduler` and the WattTime half of `thegreenwebfoundation/grid-intensity-go` all default to `api2.watttime.org` with a `/v2` prefix and call `/index`, `/forecast` and `/register` -- the v2 endpoint set. This Recipe models `api.watttime.org` and the `/v3/*` routes behind its JWT gateway, so a v2 client points at a different hostname before it reaches a route that differs too. The Basic auth against `/login` is the one thing they share, which makes these look closer than they are. Neither npm nor Packagist has anything at all: both answer the name with `wakatime` wrappers. `grid-intensity-go` is mapped under Electricity Maps, whose provider inside the same module does target the current API. It ships unmapped |
| EIA | Three letters, so both registries answer with everything else first. Packagist returns Tencent's and Alibaba's `eiam` identity products, one letter longer, plus `wakatime` wrappers matched on a fuzzy edit distance. npm's bare `eia` is `arafathusayn/node-eia`, version 0.0.1: its entire published source is `const eia = () => "eia"`, a function that returns its own name and calls nothing. `@cyanheads/eia-energy-mcp-server` does default its `baseUrl` to `https://api.eia.gov/v2` and send `api_key` as a query parameter, read out of its published tarball, but it is the same publisher whose MCP servers this file already leaves unmapped for OpenFDA and Europe PMC, on the same reasoning: a project holding it wants an agent rather than an emulator of this host. It ships unmapped |
| OpenStates | The obvious name is twelve years stale. npm `openstates` wraps the Sunlight Foundation's original service (`sunlightlabs.github.io/openstates-api`), last published in 2013, before the project moved to Plural and became the v3 API this Recipe models. `walkergriggs/openstate` is an unrelated finite-state-machine library and `osflab/osf` is a PHP framework whose initials happen to spell OpenStates Framework. The one package that reaches the current `v3.openstates.org` host with an `X-API-KEY` header, confirmed by reading its published `dist/`, is `@cyanheads/openstates-mcp-server` -- the same MCP-server publisher already left unmapped above for EIA, OpenFDA and Europe PMC. Nothing changes about that reasoning just because this is the provider where its output happens to match. It ships unmapped |
| Aviationstack | Every candidate is either the wrong shape or too thin to trust. `apilayer/aviationstack`, the vendor's only GitHub presence for this name, is a documentation-only repository -- one README, archived -- not a client. The five npm packages that do call the right host and query parameter (`aviationstack-serversdk`, `apilayer-aviationstack`, `aviation-stack-api`, `aviation-stack-sdk`, `apimatic-avaitionstack-sdk`) are all APIMatic-generated and published from the identical apimatic.io maintainer account (`developer-sdksio`), fragmented across five near-duplicate names at six to ten downloads a month each -- generated inventory rather than a dependency a project chose. `dommer1/laravel-aviationstack`, the one Composer wrapper, is marked abandoned by its own author. It ships unmapped |
| Hardcover | The name loses to `hardcode`/`hardcoded` on both registries before a real candidate turns up -- security scanners, SAML bundles hardcoded for one IdP, none of it this API. Hardcover's own GitHub organisation (`hardcoverapp`) publishes documentation, static assets and a wiki, and nothing else; there is no official client anywhere. The two Go hits that do speak this host's GraphQL API, `blampe/rreading-glasses` and `freeeve/libcat`, are both subpackages buried inside somebody else's unrelated tool -- a Readarr metadata proxy and a personal library cataloguer -- implementation detail rather than a dependency a project holds on its own. It ships unmapped |
| Akamai | Akamai publishes an official EdgeGrid client in every ecosystem searched -- `akamai-open/edgegrid-client` and `akamai-open/edgegrid-auth` on Packagist (2.7 and 3.6 million downloads), `akamai-edgegrid` on npm (463,577 a month), `AkamaiOPEN-edgegrid-golang` in Go -- and every one of them is the same thing: signing middleware for the EdgeGrid HMAC scheme this Recipe's header checks, with no product baked in. The PHP client's own README shows the shape plainly: credentials from `~/.edgerc`, a bare Guzzle client or middleware, and the path -- `/papi/v1/...`, `/config-dns/v2/...`, `/cps/v2/...`, whichever of Akamai's hundred-plus REST APIs the calling code chooses -- supplied entirely by the caller. A project depending on it could be calling Property Manager, this Recipe's own subject, or Edge DNS, or Bot Manager, or Image and Video Manager, and the dependency name gives no way to tell which. The Go module even ships a `papi` subpackage naming Property Manager specifically, but `go.mod` records the module, not the import path inside it, so that specificity is invisible at the point this file has to decide. It ships unmapped |
| PatentsView | Confirmed retired rather than merely hard to find. This Recipe's own header already establishes that patentsview.org and api.patentsview.org both 301 to USPTO's transition guide, unconditionally, credentialled or not, checked live -- and neither registry turns up anything that still calls the old host either. Every current hit for the name is an MCP server wrapping the replacement, USPTO's Open Data Portal, not this API: `@pipeworx/mcp-patents` says so in its own description, and the rest (`mcp-uspto`, `patents-mcp-server`) are broader USPTO/EPO tool servers. Packagist has nothing at all. There is no package left to map, because there is no host left for one to call. It ships unmapped |
| Bitwarden | Every official package is the wrong surface, on purpose per this Recipe's own header, which models the Public API and identity host rather than the CLI or the vault client. `@bitwarden/cli` is the CLI itself; `@bitwarden/sdk-napi` and Packagist's `bitwarden/sdk` (published under a third party, `MaliRobot/bitwarden-sdk-php`) are bindings for the Secrets Manager vault SDK; `@bitwarden/mcp-server` is agent tooling, the exclusion this file already applies to the openFDA and Europe PMC MCP servers, and `@pikku/addon-bitwarden` is the same shape a layer further down -- generated scaffolding for an agent framework, not a client a backend calls. The community clients found wrap the same excluded surfaces from the other side: `jalismrs/bitwarden-service.php` shells out to a `bw` binary its own test fixtures ship a stub of; `fbegyn/bitwarden-api` in Go says so in its own description, "depends on the `bw` CLI tool"; and `mutablelogic/go-client/pkg/bitwarden` implements the master-password vault-sync protocol (KDF and crypto key derivation in `crypto/key.go`) rather than the organization-management REST API. Nothing on any registry reaches api.bitwarden.com or identity.bitwarden.com as a general-purpose library. It ships unmapped |
| 511 | The name invites exactly the collisions this batch's brief warned about. Packagist's only near-hit, `kirkaracha/511-data-exchange`, says what it is in its own description -- "Laravel implementation of Open 511 Data Exchange API" -- Open511, a separate traffic-incident standard with no relationship to 511.org's transit service. The one Go client that does say 511.org, `apokalyptik/511`, reaches `services.my511.org/traffic/getpathlist.aspx`, checked by reading `drive/main.go`: 511.org's own, older Driving Times API, a different subdomain and a different product from the SIRI-based transit service (`api.511.org`) this Recipe models. Nothing on npm reaches either 511.org surface at all -- the name loses immediately to Salesforce's `sf` abbreviation and the data-serialisation format called Transit. It ships unmapped |

**OpenAI is mapped on purpose without a Recipe.** A dependency Cauldron
recognises and cannot emulate is reported by name, so a developer is told
which calls will still reach the real network, which is more useful than the
"looks like an API client" heuristic that would otherwise catch it. Removing
the mapping made the warning vaguer, not safer, and two tests in the CLI
package said so.

What that arrangement cannot do on its own is tell a deliberate
forward-mapping from a typo: a mapping naming "postmarkk" would quietly become
a warning about a provider nobody has heard of. Every mapped name that does not
ship is now declared as intentional, and a test refuses the ones that are not.

An OpenAI Recipe is worth writing on its own merits.

## Identifiers that are numbers on the wire

Cauldron mints and stores every identifier as a string, because that is the
only form every style shares and the only form a path parameter arrives in.
Until `id.type` existed it also *sent* every identifier as a string, and a
great many providers send a number.

It is not cosmetic. `id === 1` fails against `"1"`, `typeof id === "number"`
fails, and a schema declaring `"type": "integer"` rejects the response
outright. That is the exact class of bug Cauldron exists to catch, committed
by Cauldron.

Ninety-seven Recipes send at least one identifier as a number now, and each
carries a case asserting an unquoted one, so removing the declaration fails
something. Three of them already had cases asserting the quoted form, which is
to say three cases were pinning the bug in place.

The rest of the `numeric` resources are below. Which of the two a provider
sends has to be read from its documentation rather than assumed: sending a
number where the provider sends a string is as wrong as the thing this fixes.

| Recipe | Expected | Note |
|---|---|---|
| RingCentral | needs checking | Message and extension ids may be numbers or numeric strings |
| Postmark | needs checking | Bounce ID is a number; the casing of the field needs confirming too |
| SendGrid | needs checking | Suppression ids are numbers on some endpoints |
| ~~Documenso~~ | **number** | Settled against its description, which declares the document id and the documentId a create answers with as numbers. Both were strings and are not now |
| ~~Intercom~~ | **string** | Settled against its description, which declares `conversation.id` as a string and gives `"1295"` as the example -- a quoted number, which is the case this table exists to keep people from "fixing" |
| HubSpot | **string** | Contact, deal and company ids are quoted, which is why the default stays string |
| Jira | **string** | Issue id is a quoted number; the key is the readable identifier |
| QuickBooks | **string** | Id is quoted everywhere in the JSON API |
| DocuSign | **string** | recipientId is quoted |

The ones marked **string** are correct as they stand and are listed so nobody
"fixes" them.

Three are still unchecked: RingCentral, Postmark and SendGrid. None of the
three publishes a description this repository has a copy of, so settling them
means reading the provider's documentation rather than running check.

A second question turned out to sit behind this one. A provider can send two
identifiers and accept only one of them in a path, which no amount of getting
the type right will fix: GitHub numbers an issue per repository and gives it a
global id as well, GitLab does the same with `iid`, and Buildkite addresses a
pipeline by slug and a build by number while sending a UUID for both. All
three are modelled on what the URL takes now, with the other identifier beside
it. Code that stored the wrong one and asked for it back got a 404 from the
provider and a record from the fake.

## Rules considered and rejected

The validator has grown a lot of rules, and most of them came from noticing
the same mistake in a third or fourth Recipe and deciding it should not be
possible to make again. That has worked well enough that the reflex now fires
on things it should not, so the ones that were tried and abandoned are worth
writing down. Each of these looked obviously right and was measurably wrong.

### A default no fixture exercises

The reasoning: four Recipes in a row had a field whose declared default every
fixture overrode, so mutating the default proved nothing. Razorpay's
`captured: false`, Midtrans's `fraud_status: accept`, Alpaca's `primaryKey`,
Meilisearch's. Each time the fix was to let the declaration carry the common
case so a mutation could reach it.

Measured across the portfolio: **301 fields**. Almost all of them legitimate.
A default is also documentation of the provider's initial value, and a safety
net for a fixture written later that omits the field. Forcing three hundred
fixtures to stop spelling out a status would make them harder to read, not
easier to trust: a fixture that says `status: settled` is clearer than one
relying on a default declared eighty lines above.

### A default every fixture contradicts

The narrower version, on the reasoning that a default nothing agrees with is
probably wrong. Measured: **66 fields**, and reading them settles it. A
counter defaulting to zero is correct even when every fixture holds a real
count. A status defaulting to the initial state is correct even when every
fixture holds a later one. That is what an initial value is.

### What the four real findings actually were

Judgement about which particular default carried the interesting claim, not a
property of defaults in general. `fraud_status: accept` matters because the
whole Recipe is about the one transaction where it is `challenge`. `count: 0`
on a counter matters to nobody. No mechanical rule separates those, and the
one that tried would have produced three hundred false alarms.

The check cost two queries and is the reason this is a note rather than a
rule. Worth doing before writing the rule, every time.

## What a create echoes

Twilio Verify's check takes the verification code in the request body. Cauldron
echoed it back, because a create answers with the record it built and a record
absorbs whatever was posted at it. So the sandbox was handing back the second
factor in the reply to the attempt to use it, and a test asserting on that
would have passed.

`returns` fixes it, and the fix is per-route: name what the provider sends and
the rest is dropped. Both of that Recipe's creates now do.

The general question is unanswered. There are 129 create routes across the
portfolio and 115 of them declare no `returns`, which is right for most
providers, because most creates really do echo what you sent. The ones that
matter are the creates whose *request* accepts something the *response* does
not return, and where that something is a credential: a password, a token, a
code, a key. No shipped suite demonstrates one -- a scan for conformance cases
posting an undeclared field found zero -- but that is a statement about the
cases, not about the mechanism. Whatever a user's own code posts is echoed
whether a case does or not.

Answering it properly means reading each provider's request schema beside its
response schema, one provider at a time, which is the work. Nothing here can
be inferred from the Recipes, because the field that would prove it is exactly
the field nobody wrote down.

## What `check` cannot disambiguate

Polygon's description declares four paths that differ only in the name of a
path parameter: `{stocksTicker}`, `{cryptoTicker}`, `{forexTicker}` and
`{optionsTicker}`, all at
`/v2/aggs/ticker/{...}/range/{multiplier}/{timespan}/{from}/{to}`.

`templatePattern` deliberately ignores parameter names, because a Recipe names
its parameters for itself and a description names them for its own reasons.
That is right almost everywhere and wrong here: the first match wins, and the
first of those four in sorted order is crypto.

The visible symptom is a finding that is not true. Only the stocks variant
carries `otc`, so a stocks Recipe checked against this description is told
`otc` is undeclared -- by the crypto schema, which the reader has no way to
know was consulted.

Nothing is broken enough to warrant a change yet. Fixing it means either
preferring the template whose parameter name matches the Recipe's, or
reporting every candidate rather than the first, and both want a second
provider showing the same shape before the right answer is obvious. Written
down so the next person to see a false finding here knows where it came from.

## A partition in the path, repeated in the body

A route that scopes by a path segment needs that segment as a field on the
resource, because that is how the record is partitioned. Nothing then stopped
it being emitted, and most providers do not repeat a partition they already put
in the URL: Fly does not send `app_name` on a machine, Tradier does not say
which account an order is in.

The only way to say so was a route's `returns` naming every other field --
twenty-three names to hide one on Fly, repeated per route, and silent about the
resource itself. `in: "-"` says it once, on the field, where the fact belongs.
It reuses the idiom `id.field: "-"` already had: a dash means the provider does
not send this.

**The audit, and what it does not license.** 119 scope fields across 39 Recipes
reach the wire with no conformance case mentioning them. That is not 119 bugs.
Some providers really do echo the partition -- Discord sends `channel_id` and
`guild_id` on a message, Mailchimp sends `list_id` on a member -- so the
difference is a fact about each provider and has to be read before it is
changed. Hetzner was checked during this work and left alone: its `resource_id`
really is on the wire, nested inside `resources[0]`.

Three are done, because all three were verified while writing their Recipes:
Fly, Tradier, and Hetzner-by-exclusion. The remaining list, largest first:

```
airtable 10   statuspage 7   dockerhub 6   gitlab 6      posthog 6
svix 6        webflow 6      cloudflare 5  googlecalendar 5
basecamp 4    bitbucket 4    miro 4        twilio 4      bugsnag 3
```

The work is one provider's response schema at a time, and the evidence is
exactly the field nobody wrote down. Anyone picking this up should change a
Recipe only after reading what the provider sends, and add a case asserting
the absence, so the next audit counts it as decided rather than as unexamined.

## A 200 that is not the thing you asked for

`https://api.tradier.com/v1/openapi.json` answers 200 with 115 kilobytes of
HTML: a ReadMe documentation page, not a specification. Nothing in the status,
the size or the URL says so.

That is the same shape half the Recipes here describe, arriving in the middle
of writing one, and worth recording for two reasons. It is the third
documentation dead end in this stretch -- after Sift, Klarna and Render -- and
it is the only one that looked like a success. A script fetching specs on a
schedule would store this and never notice.

The lesson for this project is small and specific: when a spec URL is used as
evidence, parse it before believing it. Content-type is not enough either;
ReadMe serves this as text/html and a fetcher that checked would have been
right, but the useful check is whether the bytes are the document you wanted.

## The verifier could not read the number it was checking

Typesense's score made this visible and it was never about Typesense.
`encoding/json` decodes every number into a float64 unless told otherwise, so
`cauldron verify` could not see any integer above nine quadrillion. It read
578730123365711993 as 578730123365712000, and read the adjacent score as the
same thing, which means no conformance case could have told those two apart.

The emulator was right the whole time -- the wire carried the exact digits, and
a curl showed it. The tool that exists to catch an emulator sending the wrong
number was the thing that could not see the right one.

`UseNumber` fixes it, and all 1524 cases still pass afterwards, which is the
useful part of the result: the change made the comparator stricter and nothing
had been relying on the looseness.

Two things worth keeping from how it was found. The first mutation pass on the
fix reported two survivors, and both were correct: the test decoded the body
itself, so `UseNumber` in the production path was never exercised. A test that
reimplements the thing it is testing is not a test. The second is that one of
the three branches added really was inert -- `json.Marshal` already renders a
`json.Number` as its literal digits -- so it was removed rather than kept, on
the same rule this project applies to Recipes.

## What `check` cannot read

`cauldron check` reads OpenAPI 3 and refuses Swagger 2.0 outright, with a
sentence saying so rather than a guess. That is the right behaviour and it has
a cost worth writing down: two providers shipped in this stretch publish
Swagger 2.0 -- Finnhub and Weaviate -- so neither Recipe could be
cross-checked against its provider's own description, and both rest on reading
alone plus a mutation pass.

That is a smaller claim than the ones that say "nothing in this Recipe is
contradicted by the description", and the difference is not visible from the
Recipe. Converting a 2.0 document to 3.0 before checking would close it;
whether that belongs in this tool or in a note telling the author to run a
converter is undecided.

## Picking providers by whether their documentation exists

Four dead ends in a row -- Sift, Klarna, Render, then Statsig and ImageKit in
one sitting -- came from choosing a provider for its shape and finding out
afterwards that nothing citable describes it. Tradier's was the worst of them:
`api.tradier.com/v1/openapi.json` answers 200 with a documentation page.

The order was wrong. Reachability is a cheap check and interest is an
expensive one, so reachability goes first now: probe several candidate spec
URLs in one command, look at the first thirty bytes of each, and pick from
what actually resolved.

That is how Ory Hydra was chosen. It also surfaced Keycloak, which resolves
and is unclaimed.

Providers whose API is open source or whose spec lives in a public repository
are the reliable seam. The ones checked so far:

```
resolve:      ory/hydra, ory/kratos, typesense, weaviate, qdrant, apify,
              twilio, hetzner, fly, polygon, finnhub (swagger 2.0)
do not:       statsig, imagekit, render, moov, checkout.com, inngest,
              gitea (moved), hashicorp/vault (moved)
lies:         tradier (200, HTML)
```

## GraphQL, now reachable

Seven providers were blocked on one missing feature and recorded as seven
separate judgements: Linear, Monday, Attio, New Relic, Railway, ShipHero, and
half of Fly.io. It was costed here one cycle and built the next.

What it took was smaller than the note guessed. A GraphQL API is one path and
one method, so the path cannot say which route answers -- but the query can.
`selects:` names the root field a route serves, and the route matches only
when the body's query mentions it. Nothing parses GraphQL: it looks for the
word, which is enough to pick a fixture and is the same bargain every Recipe
here already makes.

Two things came with it. `fields:` on a route, because a one-path API keys its
envelope metadata by the query that was asked -- ShipHero puts complexity and
request_id at `data.orders.*` and `data.products.*`, so a Recipe-wide constant
would stamp one connection's numbers onto the other's response. And an
unmodelled query is a 404 rather than a 405: counting selecting routes as
"methods this path allows" told a client to change the method it already had
right.

The Relay shape needed nothing new. `collection: data.orders.data.edges` with
every field declared `in: node` produces edges holding a node and a cursor
side by side, which is what Relay is.

ShipHero and Linear are shipped on it, and the two differ in every way that
matters: ShipHero nests its connection under a second `data` and answers a bad
query with 200, Linear puts `edges` and `nodes` side by side and answers 400.
A client written for one is written for the other until it is not, which is
the argument for having both.

Linear also tightened the mechanism. `selects: me` matched
`query { viewer { name email } }`, because "name" contains "me", so the match
is now a whole-word one. Short root fields are common -- me, user, node, team
-- and an accidental match sends a request to the wrong fixture, which is a
bug this project exists to catch rather than commit.

The other five are ordinary provider work: find the queries, find the
responses, write the Recipe.

## Blocked on documentation, with what would unblock each

The pattern is common enough now to be worth one list rather than a sentence
per row.

| Provider | What is reachable | What is not | What would unblock it |
|---|---|---|---|
| Logiwa | Marketing pages, a webhooks blog post | Every endpoint reference. The developer centre needs a purchased API user seat | An account, or a customer's copy of the reference |
| Intelcom / Dragonfly | The introduction page at developers.intelcomexpress.com, which names six APIs and their required call order | Every reference page under it: `/reference/tracking-api` and friends 404, and the site's own `llms.txt` answers 200 with the ReadMe shell | A partner login, or one captured response per endpoint |
| Sift | Prose describing a 0-100 score | The score response body. The scale is decimals on the wire and prose in the docs, and that number is what every integration branches on | One real response |
| Klarna | Navigation | Every schema. The portal renders through JavaScript | An account or a published spec |

Render was on this row too, on the same complaint -- `api-docs.render.com`
404s on the reference pages this was tried against. It was not blocked: its
own "OpenAPI Spec" page names the real file, `render-public-api-1.json`,
three redirects past the well-known `render.com/openapi.json`, and that one
resolves. Shipped; see the hosting and deployment table above.

The useful observation is that three of these four are ReadMe or JavaScript
portals, and the fourth is a login. A provider that publishes a spec file --
in its own repository, or at a URL that answers with the spec rather than
with a page about the spec -- is a provider that can be modelled honestly.
That is now the first thing checked, not the last.

## One provider, two shapes for the same status

The npm registry answers a 404 two ways, checked against registry.npmjs.org on
2026-08-22:

| Request | Answer |
|---|---|
| `GET /cauldron-no-such-package-xyz` | `{"error":"Not found"}` |
| `GET /left-pad/99.99.99` | `"version not found: 99.99.99"` |

Same registry, same status, one object and one bare JSON string. Code reading
`body.error` off the second finds `undefined` rather than failing, so a client
reporting `body.error` as the reason reports "undefined" -- quieter and worse
than a thrown error.

Closed. An error may carry its own `style:`, `string` is a style, and a route
may name the error its absences raise with `not_found:`. The npm Recipe
declares both shapes and a case asserts the raw bytes of the second, because
the quotes are what separates a JSON string from plain text.

One decision worth recording: a route naming its own not-found error gets the
identifier alone as `{detail}`, rather than the default's `resource: id`. The
composite belongs to the default message's phrasing. npm says "version not
found: 99.99.99", and "version not found: release: 99.99.99" is not a sentence
anybody sends.

## What check says once it stops needing to be told the prefix

`check` compared a Recipe's whole path against a description's relative one and
reported the difference as a missing path, so a provider whose description
declares `/files/{id}` beside a server of `https://api.box.com/2.0` had all
fifty-eight of its routes reported as paths the description does not have.
None of it was a disagreement. The prefix is written down in the description
and is now read from there, with `--base` still overriding and the inferred
value printed so a wrong guess is visible rather than silent.

Swept across the thirty-four descriptions already fetched: **twenty-one are
contradicted by nothing at all**, up from thirteen, and no description in that
set now fails to read.

Four more were fetched on 2026-08-23 -- Documenso, Apify, Fly and LaunchDarkly
all publish one at a stable URL -- and every one of the four contradicts
nothing. Thirty-eight descriptions, twenty-five clean.

Documenso's is the only one that said anything at all: its `rate_limit` error
answers 429 and the description declares no 429 anywhere, which is the
declined-to-say category rather than a disagreement.

Guessing at URLs for the providers whose paging is still unnamed failed
again, as it did before: AssemblyAI, Vonage, Deel and Segment were each tried
at the address a description would plausibly live at and none of the four
answered with one. That approach has now failed ten times out of ten and is
not worth a eleventh -- the descriptions that exist are found by reading a
provider's documentation, not by pattern.

Two of that rise came from a status nobody declares. A description that names
a status nowhere has not contradicted a Recipe answering with it -- it has
declined to say. Nineteen of the thirty-one descriptions here declare no 429
anywhere, and Stripe, Twilio, Slack and Square are among them, so treating
that silence as a contradiction was reporting the fact that rate limits get
documented in prose. Those findings are still printed, under their own heading
and in their own words.

Three of the thirty-four were not disagreements at all: the description and
the Recipe were about different APIs. CircleCI's cached description declares
`/api/v1` while the Recipe models v2, and Pipedrive's and Snyk's are the same
mistake. Every route in those Recipes came back as a path the description does
not have, which reads exactly like a Recipe full of invented paths and is one
error made by whoever fetched the file. `check` now says so before the list,
because the list is what somebody would otherwise start working through.

That leaves **eight**, and all eight have now been read. Every one is either
fixed or carries a note in its own Recipe saying why it is not, so the next
person to run `check` does not repeat the work.

What the sweep produced, counted honestly:

| | |
|---|---|
| Recipe defects found and fixed | 4 |
| `check` bugs found and fixed | 3 |
| descriptions wrong or partial, Recipe correct | 5 |
| behaviour nobody had defended, now pinned by a case | 3 |
| the wrong document entirely | 3 |

The four defects were worth the whole exercise. Increase modelled a collection
that does not exist and put a transfer's return on the wrong object. Adyen's
refund used the payment's name for its reference. OneSignal's create answered
with delivery counts on a delivery that had not happened. Xero explained the
Microsoft-date trap in a comment and then fell into it on the next field.

The five where the description was wrong are worth as much in a quieter way,
and each now says so in its Recipe -- because the obvious way to silence a
report is to change the Recipe, and that would have been wrong five times.

The mixed bag, for anybody reading a future report:

- **The Recipe is right and the description is stale.** Meilisearch's document
  add answers 202 because the work is enqueued, and its description says 200.
  That one already carries a comment saying so, because the obvious way to
  silence `check` is to make the Recipe wrong.
- **A field a create should not echo.** OneSignal's notification carries
  `completed_at`, `converted`, `failed` and `queued_at`, none of them declared
  on the create, and none of them things that exist before a send.
- **A tool bug rather than either.** Chargebee wraps each entry in its listing
  under the resource's own name, which `check` did not descend through, so all
  eleven of its subscription fields read as fields the description does not
  declare. Fixed; two findings remain and they are the next item.
- **Two catalogs behind one API version.** Chargebee's v2 has two incompatible
  subscription shapes and ships two descriptions for them. This Recipe models
  Product Catalog 1.0 -- `plan_id`, `plan_quantity` -- and PC 2.0 has neither,
  carrying `subscription_items[].item_price_id` instead. Now written down in
  the Recipe, because it was a decision nobody had recorded. A 2.0 Recipe
  would be a second Recipe.
- **A field a create should not have echoed at all**, which is what OneSignal's
  turned out to be: it answered with delivery counts on a notification that had
  delivered nothing, and its own description and this Recipe's own comment both
  said so already.
- **A path the description does not have**, which may be the description's gap
  or the Recipe's invention, and needs reading either way. Braze's export
  endpoint and PostHog's capture are both of these, and PostHog's is on a
  different host from the API the description covers.

Five descriptions in the cache cannot be read at all and are not in those
counts: Bitbucket's is Swagger 2.0, Linear's is not a description, and
Cal.com, GitLab and Gorgias answered with HTML pages rather than a file.

## One provider, two hosts

The format gives a Recipe one upstream and hangs every route off it. Three
providers in the queue do not work that way, and the third is already shipped:

- **PostHog.** `/capture/` is ingestion and lives on a different origin from
  the management API. Its own OpenAPI describes the management API and does
  not mention capture at all, which is why `check` reports that route as a
  path the description does not have. The route is right; the description is
  about the other host.
- **Customer.io.** Already in the queue with the same note: the Track and App
  APIs are separate hosts with separate credentials, "which is the sort of
  thing that only fails in one environment".
- **Braze.** Its REST endpoint is per-cluster -- `rest.iad-01.braze.com` and
  the rest -- which is a configurable host rather than two of them, so it is
  a different problem wearing similar clothes.

What a single-host emulator gets wrong is narrow and real: an application
integrating PostHog configures an ingestion host and, separately, an API host.
Code written against a fake that serves both from one origin runs green with
only one of them configured, and the missing one is not discovered until
something is deployed.

This is not built and is not obviously worth building. It is written down
because it has now been noticed three times from three directions, and a gap
found three times and recorded none of them is a gap that will be noticed a
fourth.

## A create echoes whatever you send it

Found while trying to mutation-test a field rename, and it is the reason that
rename had gone unnoticed.

A create stores the decoded request body as the record, without dropping the
fields the resource does not declare. So this:

    POST /v71/payments/{id}/refunds
    {"merchantAccount": "...", "totallyMadeUpField": "xyzzy", "amount": {...}}

answers with `totallyMadeUpField` in the body. No provider does that.

**The consequence worth naming: a conformance case asserting a value it sent
on a create cannot fail.** The validator already refuses a case whose *every*
claim is an echo, but a case with one real claim and several echoes still
contains claims that hold whatever the Recipe says. Adyen's refund carried
`merchantReference` where Adyen sends `reference`, and a case asserting the
name would have passed either way.

Filtering to declared fields was tried and the blast radius is almost nothing:
all 1697 conformance cases still pass. One runtime test does not, and it is
the reason this is not simply a bug.

`TestCreateAcceptsFormEncoding` sends `metadata[order_id]=42` to Stripe and
expects it back. Stripe really does accept arbitrary metadata and really does
echo it, and `metadata` is declared nowhere in that Recipe -- the echo is how
it works today.

So the echo was doing two jobs: supporting a genuine free-form container, and
silently accepting nonsense.

**Closed.** A field may now be declared `type: map`, which keeps whatever keys
it was sent, and creates and updates drop everything the resource does not
declare. Stripe's customer declares `metadata: {type: map}` and the form
encoding that fed this finding still works; nothing else in 167 Recipes needed
a change, because nothing else was relying on the echo.

The cost of leaving it open was not the stray key. It was that a conformance
case asserting a value it sent on a create could not fail, which is how
Adyen's refund kept the payment's name for its reference.

### What the sweep after it found

Both of Stripe's created resources now declare `metadata`, because declaring
it on the customer and not the payment intent would have been the worst of the
three options -- and the payment intent is where it matters most, since the
webhook that arrives later carries the intent and nothing else.

**No other Recipe models free-form data at all.** Not one of the twenty that
mention metadata declares a field for it, and no conformance case anywhere
sends metadata in a request. So the echo was the only thing making those
Recipes appear to accept it, and appearing to accept it is what a fake should
not do: a provider whose Recipe never modelled metadata now drops it, which is
the truthful answer to "does this Recipe describe that behaviour".

Declaring `type: map` across the providers that really do support metadata --
Adyen, Chargebee, Orb, Paddle, Shippo and the rest -- would mean modelling
from memory, which is the guessing this project refuses everywhere else. It
wants the same treatment every other field gets: read the provider, then
declare it.

**Nothing tells a Recipe author they need it.** A provider with free-form
metadata and a Recipe that does not declare it produces no warning, because
absence of a field is indistinguishable from a provider that has none. That is
the gap this closing leaves open.

## Declared events nothing could fire

The runtime emits `resource.action` -- `customer.created`, `invoice.updated`
-- and nothing else. **438 of the 482 events declared across these Recipes are
not named that way**, and every one of those names is the provider's own:
Freshdesk's `ticket_create`, Bitbucket's `repo:push`, Zoom's `meeting.started`,
Recurly's `new_subscription_notification`, Webflow's `collection_item_created`.

So for 97 Recipes a create fired nothing. Not a wrong event, not an error --
silence, which is indistinguishable from a provider that does not send one,
and a handler waiting for it waits for ever.

A route may now name what it fires with `emits:`, and the validator refuses a
name the Recipe does not declare. Freshdesk's ticket create and update are
wired as the demonstration.

Twenty-one more providers are wired now, forty-two routes in total, and the
count of events nothing can fire is **396** rather than 438. Every one was
reviewed rather than pattern-matched, and three suggestions were rejected on
reading:

- **Miro.** A substring match proposed `update:board` fires `board_item.updated`.
  It does not: `board_item.*` belongs to the item resource, and updating a
  board is not updating an item on it. Only `create:item` was wired, and
  `board.created` already fired by convention.
- **Salesforce.** `AccountChangeEvent` fires for any change to an account
  rather than for an update specifically, so mapping it to the update route
  would be a guess about Salesforce's change-data-capture semantics.
- **Google Calendar.** `calendar#event.created` carries Google's `kind`
  prefix, and whether that string is the webhook's event name or a marker on
  the resource is exactly the sort of thing to read rather than assume.

### And the number was overstated

Forty-seven routes are wired now and 391 declared events still cannot be
fired, which reads like a large backlog and mostly is not. Classified:

| | |
|---|---|
| not a lifecycle event at all | **263** |
| lifecycle-shaped, for a resource this Recipe cannot mutate | **88** |
| lifecycle events for a mutable resource -- actually wirable | **40** |

The 263 are things like `crawl.completed`, `video.asset.ready`,
`user.session.start` and `payment.failed`. No create, update or delete
produces any of them, and none ever should: they are what `cauldron emit`
exists for, and a Recipe declaring them is describing the provider correctly.

The 88 name a resource with no create, update or delete route. Those are
Recipes modelling less of a provider than the provider emits about, which is a
different question from wiring.

**So the remaining work is forty routes, not 391 events.** The earlier note
here said "the other 95 are the work this leaves", and that was wrong in the
direction that flatters the finding.

### And forty was overstated too

The same measurement, made properly, gives eleven. Forty counted an event as
wirable when the Recipe mutates that resource *somewhere*; the route the event
would actually hang off has to exist, and for most of them it does not. Jira
declares `jira:issue_updated` and has no update route at all. Dropbox declares
`file_added` and creates no files.

Of the eleven that survive, none is a plain wire:

- Six are substring accidents. `miro:board` matched `board_item.updated`
  because "board" is inside "board_item"; `okta:user` matched
  `group.user_membership.add`; `recurly:subscription/create` matched
  `renewed_subscription_notification`, which is a renewal and not a creation.
- The rest are a different shape entirely, and are the reason `emits_when`
  now exists: `ticket_status_change`, `ticket_priority_change` and
  `taskStatusUpdated` are not what an update does, they are what one
  particular change to one particular field does.

Freshdesk is wired that way now. The general lesson is the one from the
Stytch finding, arrived at from the other side: an event left unfired has two
possible fixes, and "wire it to the nearest route" is only the right one when
that route is genuinely what produces it. Hanging `ticket_status_change` off
every update would have made the emulator fire it constantly and production
fire it rarely, which is worse than the silence it replaced.

**The unconditional wiring is done.** What is left needs a provider read, and
now has somewhere to go when the answer is "only when this field moves".

ClickUp is wired that way too now, and wiring it turned up the thing worth
having done this for: `taskStatusUpdated` did not fire because the status had
not changed, and the status had not changed because a ClickUp task's status
could not be written at all. `flatten` dropped every declared wrapper name
rather than the ones it unpacked, and ClickUp's status field is nested under a
wrapper of its own name, so a write in the shape ClickUp documents was
silently a no-op that answered 200. Nine fields collide that way and one sits
on a writable resource, so the exposure was exactly that one -- but it was
found by asserting a webhook, not by any of the ten cases that already
exercised the route.

### The envelope is the larger webhook gap

Of the 99 Recipes that emit events, 41 declare a payload envelope and 58 fall
back to the default, which is Stripe's `{id, type, created, data.object}`.
The README has always said this is a default rather than a claim about the
provider, and now says so with the numbers in it and a test holding them.

Three are declared now -- Shopify, Slack and ClickUp -- and doing them turned
up two things the template cannot say and one it was saying by coin toss.

~~**A top-level array.**~~ **Built.** `webhooks.payload` is `any` now, so a
payload that is not an object can be declared. HubSpot's is: an array of
subscription objects, batched, carrying no record at all -- a delivery names
the object that changed and leaves you to fetch it, so an application reading
a contact's email off the webhook is reading something HubSpot has never sent.

Two of the three providers that motivated it are still not done, for
different reasons:

- **SendGrid** can now be described and still cannot be checked. None of its
  nine events is reachable from any route it has, so declaring the array
  would be writing a claim no case could exercise -- the Calendly rule.
- ~~**QuickBooks**~~ is unblocked. `name` and `operation` are the two halves
  of an event declared as `Customer.Create`, which looked like a gap and was
  half of one: `{resource}` already gave the first half and `{action}` now
  gives the second. Xero still needs the case change on top -- INVOICE and
  UPDATE -- so that one is a third of a mechanism away rather than a whole
  one.

  The judgement to build `{action}` at all was close, and turned on a third
  provider rather than a second: Asana splits an event the same way, sending
  `{resource: {resource_type}, action}`. One provider wanting a mechanism is
  a special case; two doing it independently is a shape.

  It was defined wrongly the first time, as the part after the last dot,
  which is only the action for providers that put the resource first.
  Pipedrive does not -- its events are `added.deal` -- so the suffix was the
  object and the prefix the action. It is the event with the resource removed
  now, which needs no knowledge of the order a provider chose and is the
  better definition anyway: the action is the part that is not the thing.

### Every signature was Stripe's shape

The same finding as the envelope and a worse one, because nobody parses a
signature by hand. An application hands the header to the provider's SDK, so
a value in the wrong shape fails there rather than anywhere a Recipe could
see -- and the comment above the signing code said exactly that, above code
that gave all seventy-four Recipes Stripe's `t=<unix>,v1=<hex>` over
`<unix>.<body>`. Sixty-eight distinct header names, one shape.

`signing.format` chooses now: `stripe` (the default, so nothing changes under
a Recipe nobody has looked at), `prefixed-hex` for GitHub-style `sha256=`,
`base64` for Shopify-style, and `v0-hex` for `v0=` over `v0:<ts>:<body>`.
Shopify, Bitbucket, Slack, Zoom and Jira are correct; **69 are still on the
default**.

The `v0-hex` format was called `slack` for one change and should not have
been. Zoom's signature is Slack's, deliberately, and naming a shared shape
after the first provider to want it is the same mistake as giving everyone
Stripe's, in miniature. `stripe` keeps its name only because no second
provider here wants that shape, and would lose it if one did.

Thirteen are shaped now and **61 are on the default**, of which two are
right: Mux sends `t=...,v1=...` over `<ts>.<body>` and Stripe invented it. The
default being a convention does not make it wrong everywhere, only wrong
by assumption.

### Some of these are not a format problem at all

Three Recipes declare `hmac-sha256` for a provider that computes no HMAC, and
in all three the header name says so plainly enough that nothing had to be
looked up:

- **Discord** -- `X-Signature-Ed25519`. A public-key signature verified with
  Discord's public key, not a shared secret. The scheme is wrong rather than
  unshaped and no format fixes it.
- **Okta** -- `x-okta-verification-challenge`. Not a per-delivery signature:
  it carries the one-time value Okta sends when an event hook is registered,
  which an endpoint echoes to prove it owns the URL. Okta authenticates
  deliveries with an Authorization header the subscriber chooses.
- **Documenso** -- `X-Documenso-Secret`. The configured secret, sent
  verbatim, compared against the endpoint's own copy. Not a digest of
  anything.

Each says so in its own Recipe now. The wider question -- what a Recipe
should declare when a provider signs with something other than HMAC-SHA256 --
is a scheme, not a format, and is not answered here.

Others were looked at and left because the base string is unreachable rather
than unknown: **Box** signs body plus timestamp, **Zendesk** timestamp plus
body, **Webflow** `<ts>:<body>`, **Lob** `<ts>.<body>`, **HubSpot v3** method
and URI and body and timestamp, **Trello** body plus the callback URL. All of
them need a second header or a value a delivery does not have, the same limit
Slack and Svix have.

**Square is deliberately left on the default.** Its value is base64, which
the format could send, but the digest is taken over the subscriber's
notification URL followed by the body -- and a delivery does not know the URL
it is going to. Declaring `base64` would give a value of exactly the right
shape over the wrong string, which is the one signature bug a conformance
case cannot see. Wrong and obvious beats wrong and plausible.

Two halves can be wrong independently and only one is assertable from a
Recipe. A conformance case can pin the wrapper with a pattern, and cannot see
what string was signed -- every digest of the right length looks alike.
Signing the body where Slack signs `v0:<ts>:<body>` produces a well-formed
`v0=` value that Slack's verifier rejects, and the suite passes nine of nine.
That is why the base strings are checked by a unit test computing the HMAC
the long way rather than by a case.

The same gap appeared again with the timestamp header, in a different place.
A conformance case reads the recorded delivery; a subscriber reads the HTTP
request; the header is set in two places to reach both. Removing it from the
request alone leaves the suite green at nine of nine and every real handler
without the value its verifier needs, so that has a unit test posting to a
real sink. Twice now the thing a case structurally cannot see has been the
thing worth checking.

~~Also owed: a delivery carries the signing header and no others.~~ **Built.**
`signing.timestamp_header` names where the signed timestamp travels, and a
delivery carries it. Slack and Zoom sign `v0:<ts>:<body>` and now send the
`<ts>` a verifier needs to rebuild that string; without it the delivery
carried a signature nothing could check, which is the failure the whole
signing surface exists to avoid, reached from the other side.

The validator refuses the header on a format that does not sign a timestamp,
and on `stripe`, whose signature carries the timestamp inside the value --
sending it twice would imply a verifier needs the second copy.

Four more were within reach and turned out to want one more thing each.

The named formats became templates first, because those four are what showed
the enum was the wrong shape: nine providers vary along three axes -- the
separator between timestamp and body, hex or base64, and what prefix the
value carries -- and naming each combination gives a list where most entries
have exactly one user. `over`, `encoding` and `value` say all nine and would
say the tenth. Empty still means Stripe's, so the Recipes nobody has shaped
keep what they had.

Two rules came with it that the enum could not express, and both catch a
typo that is otherwise invisible until a verifier somewhere else says no: a
signature whose `over` omits `{body}` covers nothing it is meant to
authenticate, and a `value` without `{digest}` carries no signature at all.

### Owed: a timestamp that is not Unix seconds

The reason Zendesk, Webflow and Lob are still not done. All three sign a
timestamp beside the body, which the templates now express, and none of them
sends Unix seconds: Zendesk's is ISO 8601, and Webflow's and Lob's are
milliseconds. `{timestamp}` is seconds everywhere, in the signed string and
in the header, and those have to agree or the signature is unverifiable in a
way that looks exactly like a wrong secret.

A `timestamp_format` beside `timestamp_header` -- `unix`, `unix_ms`, `iso` --
would do it, and would be one field rather than three tokens, because the two
places must never disagree.

Not built here because the confidence that Webflow's and Lob's are
milliseconds is moderate rather than settled, and a signature is the wrong
place to be approximately right: it fails inside somebody else's verifier
with nothing to point at. Worth ten minutes with the documentation open.

**Svix** (Clerk and two others) is the fourth, and is closer than it was --
`{id}` is a token now, so `v1,<base64>` over `<id>.<ts>.<body>` is
expressible. It still wants the id in a header of its own beside the
timestamp.

### The field a handler reaches for first is often a constant

Enough envelopes have been written now for one shape to have turned up
repeatedly, and it is worth naming because it is the same bug every time:

| provider | the field a handler tries | what it always holds | where the answer is |
|---|---|---|---|
| Slack | `type` | `event_callback` | `event.type` |
| Box | `type` | `webhook_event` | `trigger` |
| Intercom | `type` | `notification_event` | `topic` |
| Okta | `eventType` | `com.okta.event_hook` | `data.events[0].eventType` |

Four providers, four constants, and in every one the real discriminator has a
different name somewhere else. Code switching on the obvious field takes one
branch forever and never learns it was reading the envelope's description of
itself rather than the event.

Two more are the same idea wearing a different hat. Lob's `event_type` is an
object rather than a string, so `event.event_type === "letter.created"`
compares an object to a string and is false forever. Clerk's `data.object` is
the string `"user"` rather than the record, so `data.object.id` is undefined
one level further down than a reader expects.

The default envelope hid all six. Under it every one of these Recipes had a
`type` holding the event name, which is the one thing none of them does.

### The seven left, and why they are left

Recurly is one and is a decision, below. The other six are not blocked by the
format and are not done, which is a different thing from being hard:

- **Xero** wants `eventCategory: INVOICE` beside `eventType: CREATE`.
  `{action}` supplies the second half; the first is the resource shouted, and
  there is no substitution that changes case. One provider, so it stays a
  gap rather than a mechanism.
- **Cloudflare** and **Klaviyo** declare event names this note is not
  confident are real. Cloudflare's notifications are alert-shaped rather than
  `zone.created`-shaped, and Klaviyo's outbound payloads are configured per
  flow rather than fixed. Declaring an envelope on top of an event name that
  may itself be wrong would put two guesses in one place, so these want the
  provider read before anything is written.
- **LaunchDarkly**, **Miro** and **Stytch** are ordinary remaining work: a
  documented shape nobody has transcribed yet.

The distinction is worth keeping. Six of these could be filled in an
afternoon by someone with the documentation open, and two of the six would be
wrong to fill without checking the events first.

### Owed: a webhook body that is not JSON

Every delivery is `json.Marshal`ed. Nothing in the format chooses otherwise,
which was invisible while the only question was what shape the JSON took.

Recurly is where it stops being invisible. Its declared events are push
notifications and a push notification is XML, with the event name as the root
element -- `new_subscription_notification` is literally
`<new_subscription_notification>`. So the choice there was an XML shape the
runtime cannot send or a JSON one Recurly does not, and the Recipe says so
rather than picking. It is the only Recipe whose missing payload is a
decision rather than a gap not yet filled.

Worth having before claiming the webhook surface is done, and it is more than
a template change: the content type, the signature computed over the encoded
body, and what a conformance case asserts against all follow from it.

### A great many deliveries carry no record

Also worth stating together, because a fake supplying the record is teaching
the easy version of an integration that is not easy: Dropbox, ClickUp,
Mollie, HubSpot, QuickBooks, Okta, Notion, Asana and Airtable all send
references or pings and leave the fetching to the application. Airtable's is
the starkest -- a base id, a webhook id and a timestamp, the same three keys
whatever happened.

### Owed: a before state in a payload template

Pipedrive sends `previous` beside `current`, holding the record as it was
before the write. A payload template has no access to one -- `emits_when` is
the only thing in the format that compares a before and an after -- so
Pipedrive's `previous` is null always, which is right for an add and wrong
for an update.

Declaring nothing would be worse rather than safer: Pipedrive always sends
the key, so its absence would be a claim that it does not. The Recipe says
which half is modelled.

**A Recipe with no writes cannot show any of this.** Calendly declares five
events and has no create, update or delete route, so nothing it declares can
fire and no case can assert an envelope for it. That
is not an envelope problem and declaring one there would be writing an
unassertable claim, which is the mistake this collection keeps finding in its
own past.

**A value the template has to compute.** ClickUp's delivery carries a
`history_items` array whose contents depend on what changed -- a status move
carries before and after status objects, an assignment carries users. A fixed
literal would be inventing a shape rather than modelling one, so it is left
out and the Recipe says why. The same will be true of any provider whose
payload describes the change rather than the record.

**The default is not Stripe's shape either.** Stripe's Recipe declares its own
envelope now, and it is larger than the default that was modelled on it: the
default carries `id`, `type`, `created` and `data.object`, and Stripe also
sends `object: "event"`, `livemode`, `api_version`, `pending_webhooks` and a
`request` object present with null members. The first two matter --
`object: "event"` is what the libraries check before treating a body as an
event, and `livemode` is what an application branches on to decide whether a
charge was real. Under the default, `event.livemode` was undefined, which is
falsy, which is accidentally right and the worst way for a check like that to
pass. The README no longer calls the fallback Stripe's shape, because it is
not one.

**A template could not name one field.** Splicing the whole record was the
only way to get data into an envelope, which left a provider that renames or
prefixes what it sends undescribable. Freshdesk wraps its payload in
`freshdesk_webhook` and prefixes every key with `ticket_`, so the record's own
names appear nowhere in it -- a template that could only merge had to send the
wrong key names or send Stripe's envelope instead. `{record.field}` reads one
field at whatever type it already has, so an integer status stays an integer.

Freshdesk is also the last Recipe that was pinning the default envelope in a
case, and the case was one written in this session: it asserted
`data.object.status` with a comment saying "Freshdesk wraps the record", which
is true of Cauldron and not of Freshdesk. The validator refuses that now --
a webhook case may not assert a path inside the default envelope unless the
Recipe declares one -- so the Square mistake cannot be made again quietly.

Left open on purpose: Freshdesk's placeholders render status and priority as
labels rather than codes, so the real payload may carry "Open" where this
carries 2. The codes are what the Recipe knows to be true of the ticket and
the label mapping is not modelled anywhere in it, so inventing one in the
envelope would put a string in the payload nothing else could account for.

**Greenhouse was not the same, and this note said it was.** It was named here
beside Calendly as a Recipe that could not fire anything. It has a create
route carrying `emits: candidate_has_been_created`, so it could fire all
along; what it has no update route for is `candidate_stage_change`, which is
a narrower claim than the one written down. Its envelope is declared now --
Greenhouse names the event `action` and keys the record on the resource, so a
candidate arrives at `payload.candidate` and none of the default's four
fields are there at all.

### How many of the rest can be checked

Of the Recipes still on the default, **7 can fire at least one declared
event and 51 cannot** -- 41 declared, 7 and 51, which is the 99 that emit
events at all. The 51 are not an envelope problem: every event they
declare is unreachable from any route they have, so no case could assert a
payload shape for them whatever the envelope said. Calendly is the clearest,
declaring five events with no create, update or delete route anywhere.

The first count of this was wrong and said 36. It asked whether a route's
resource had any declared `resource.created`-shaped event, without checking
that the verb matched the route's operation -- so Sentry counted because it
updates and deletes issues while declaring `issue.created`, which an update
does not fire.

The correction was wrong too, in a smaller way: 34 and 51 do not add up to 99
with the 13 declared at the time, and nobody noticed because the two numbers
were reported without the total beside them. Every figure in this section is
measured together now and stated with the sum, which is the only way the
arithmetic can be checked at a glance.

**A payload key needed to vary.** Square nests the record under its own type
name -- a payment at `data.object.payment`, a customer at
`data.object.customer` -- so a Recipe-wide envelope could not describe it
while template keys were fixed strings. Keys are substituted now, and
`{resource}` names the thing an event is about, taken from the route that
declares the event and from the convention otherwise so `cauldron emit`
resolves it the same way a write does.

Square is the case that shows why this is worth having: it already had a
webhook case, and that case asserted `data.object.amount_money.amount` and
passed. It was pinning a path Square has never sent, under the default
envelope, while its own name claimed it was checking that the webhook matches
the response. A Go test had the same assumption hard-coded in a helper whose
comment said it dug the record out of "whatever envelope the Recipe
declares".

**A payload needed milliseconds.** Zoom's `event_ts` is milliseconds and every
other timestamp the format emits is seconds. The Recipe format already treats
that as material -- it has `timestamp` and `timestamp_ms` field types for
exactly this reason -- so the envelope templates now have `{created_ms}`
beside `{created}`.

**And the merge order was random.** A template naming a key the record also
carries got one payload or the other depending on Go's map iteration. Slack's
envelope has a literal `type` beside the merged event, which is exactly that
collision. The record is merged first now and explicit keys win, every time.

Two of the three also had their event *names* wrong, which is worth more than
the envelope: Shopify's topics are `orders/create`, not `order.created`, and
Slack's are `channel_created`, not `channel.created`. A handler keyed on the
real name matched nothing here, and one written against this Recipe would
match nothing there. Event names are worth checking wherever the envelope is.

It is the biggest remaining distance between a webhook from this emulator and
one from the provider. ClickUp's real payload is `{event, task_id,
history_items}` and shares nothing with the shape it currently sends. The
mechanism to fix it exists -- `webhooks.payload` is a template and five
Recipes use it -- so this is 94 provider reads rather than a design question,
and the same "read each provider" work the conditional events need.

Pattern-matching candidate mappings is also spent. Widening the net produced
`ghost: create:post -> member.added`, `slack: create:message ->
channel.created` and `lob: create:address -> letter.created` -- a hit rate low
enough that reading each provider is now cheaper than reviewing the
suggestions. Each needs the
provider read closely enough to know which change produces which event, which
is not the same question as what the events are called -- Freshdesk has
`ticket_status_change` and `ticket_priority_change` beside `ticket_update`,
and deciding which of the three an update fires is a matter of reading
Freshdesk rather than pattern-matching on names.

### The signing header nothing could assert

Seventy-four Recipes name the header a webhook signature travels in, and until
now no case could assert one. Not because nobody had written the assertion --
because the format had no way to express it, and the name was only ever
applied at send time, which a conformance case never reaches: it has no
subscriber.

A case may now claim it with `signature_header:`. Square's is
`x-square-hmacsha256-signature`, lower-case and hyphenated with the algorithm
in the middle, which is not a shape anybody guesses. It is also the first
thing a handler reads, so a Recipe naming the wrong one produces a handler
that looks for a header which never arrives, and verification that fails only
against the real provider.

Square is pinned. The other seventy-three are not, and they want a webhook
case each rather than a header assertion bolted to nothing.

### And 97 of 99 assert no webhook at all

Separately from the naming: 99 Recipes declare events, and 97 have no
conformance case asserting that one arrives. 73 of those declare a signing
scheme nothing exercises either.

Mass-adding cases would produce 97 weak ones, which is the opposite of what a
case is for. They want writing where the payload shape is worth pinning, a few
at a time.

## Assessed and deliberately not done

| Provider | Why not |
|---|---|
| ~~Linear~~ | Shipped, on the GraphQL support ShipHero brought. Priority counts down and zero is not the top -- Linear's own words are 0 = No priority, 1 = Urgent, 4 = Low -- so sorting ascending puts untriaged issues above the ones on fire and descending puts Low above Urgent. Plus: a state's `name` belongs to the team and its `type` does not, three of the seven types close an issue and `duplicate` is the forgotten one, a connection carries the same list twice as `edges` and `nodes`, and `number` is team-scoped so two issues are both 123 |
| ~~Attio~~ | Shipped. It was recorded here as GraphQL-only, which was simply wrong -- it is REST and publishes OpenAPI 3.1 at `https://api.attio.com/openapi/api`. Its records are queried by POST with the paging in the body, which is what found that bug in four other Recipes |
| Temporal Cloud | gRPC rather than HTTP. The format describes HTTP surfaces and nothing here would be a Temporal client |

~~New Relic and Railway were here too~~ -- both for being GraphQL-only, which
was a reason until ShipHero brought GraphQL support and Linear and Monday
shipped on it. Neither has been assessed on anything else, so they are back
in the queues above rather than sitting under a reason that has been
withdrawn. This table held five rows and four of them were wrong: two
providers that shipped, and two more kept out by a reason that lapsed
without anybody going back to the list that rested on it. Only Temporal
Cloud is still out for the reason written beside it.

~~Monday.com belongs here too~~ -- shipped, and it was the one that found the
next two bugs.

## What Monday found

Two of them, and both were asymmetries rather than mistakes: something the
format could say in one place and not honour in another.

**A conformance case could assert an array position that no Recipe could
emit.** The comparator has understood `data.boards[0].name` from the
beginning, and `nestedObject` has built arrays for a field's `in` for almost
as long, but `setPath` -- which puts a route's declared constants into the
body -- treated the whole segment as a key. So Monday's board went out under a
key literally spelled `boards[0]`, and a case asserting the array form
reported the field missing. A key like that has now been written four times by
four different mechanisms, and every one was silent: no provider sends it, so
nothing errors, the field simply is not the one anyone asked for.

**A declared identifier shape was never on the wire.** Monday mints ten-digit
item ids. Changing the Recipe to say six broke nothing at all, because every
conformance case reads a seeded record and the generator only runs when
something is created. The declaration could have disagreed with the provider
from the first commit and no case would have said so. A fixture now has to
seed an identifier the generator could have minted, and that check found four
more Recipes where it could not have:

| Recipe | What was seeded | What it should have been |
|---|---|---|
| Discord | 18-digit snowflakes under a 19-digit declaration | 19, which is what Discord mints now |
| Mailchimp | `c1a2m3p4g5`, and a member id ending `grc2`, under a hex declaration | Hex. A member id is the MD5 of the lower-cased email and no MD5 produces a g |
| Trello | `...l01` under a hex declaration | Hex. Board and card ids happened to be valid hex already; only list ids were not |
| Statuspage | `comp0000000001` under a hex declaration | The declaration was the wrong one. Statuspage ids are lower-case alphanumeric -- `kh2n5rn2rgsk` -- so hex would mint from a narrower alphabet than the provider uses, and a client validating one against `[a-f0-9]+` would pass here and fail there |

Cohere is the exception that shaped the rule: it keys embeddings `e1` and `e2`
and declares `field: "-"`, so no client ever sees them. An identifier that is
never emitted cannot disagree with anything, and its shape carries no claim.

## Two limits of the GraphQL support, stated

Neither is a bug, and both would read as one to somebody writing a Recipe.

**A route answers with everything it declares, not with your selection set.**
Nothing here parses GraphQL; `selects:` looks for a word. So a query naming
one field of one object is handed the whole modelled shape. Monday's
complexity case originally asserted that a complexity-only query carried no
board, which is true of Monday and false of Cauldron, and asserting it would
have been asserting a pruning step that does not exist.

**Two selecting routes on one path are ambiguous when a query names both.**
Monday's complexity started as a route of its own, and a realistic query --
asking for the cost alongside the data it is the cost of -- matched both
routes equally, so which answered was arbitrary. It is now a route field on
the query it belongs to, the way ShipHero's already was. The general shape of
the rule: envelope metadata rides with the query it describes rather than
becoming a route.

## Smaller, noted rather than done

`opaque` mints from an alphabet with upper case in it, and several providers
whose ids are opaque use lower case only -- Statuspage among them. A client
lower-casing an id before comparing would not notice; one comparing directly
would, and only against the emulator. A `lowercase` flag or a distinct style
would close it.

## Paging that travels in the request body

Attio asked one question -- it pages by `limit` and `offset` inside a POST
body -- and the answer turned out to be a bug in four shipped Recipes.

Every pagination parameter was read from the query string. A listing reached
by POST usually carries its paging in the JSON body, so for those routes there
was nothing to read. What that produced is worse than an error: the limit fell
back to the route's default, the default is larger than any fixture, so the
first response held the whole collection and reported no next page. A paging
loop written against that runs exactly once, takes neither branch, and passes.
The first collection large enough to page is in production.

Dropbox is the clearest case, because its own conformance suite had been
written around the gap:

```yaml
      json:
        path: /documents
      query:
        limit: "1"      # a parameter Dropbox does not read
```

That case passed. It was a record of what Cauldron did, not of what Dropbox
does. Sending what a real client sends -- `{"path": "/documents", "limit": 1}`
-- returned the whole folder with `has_more: false`.

`pagination.in: body` fixes it, and a dotted parameter name nests, because a
provider that puts paging in the body often puts it inside something.

| Recipe | Was | Is |
|---|---|---|
| Dropbox | `?limit=1`, a parameter Dropbox does not accept | `limit` and `cursor` in the body |
| Plaid | `?limit=2` beside a top-level `count`, neither of which Plaid reads | `options.count` and `options.offset`, and offset style rather than cursor |
| AWS Secrets Manager | `?limit=1` | `MaxResults` and `NextToken` in the body |
| DynamoDB | `?limit=1` | `Limit` in the body |

Two of those had a note in them already saying the style was unverified. The
note was right to be there and this is the check it asked for.

Three further corrections came out of it. Plaid's two account listings
declared cursor pagination and Plaid does not page them at all, so the
declaration claimed a round trip that does not exist. DynamoDB resumes from
`ExclusiveStartKey`, which is the whole primary key in attribute-value form
rather than a string, so no cursor name is declared for it -- naming one would
claim a round trip that does not work. And both Dropbox and Plaid now have a
second-page case, which is the branch a paging loop actually depends on and
the branch no case here had ever taken.

### Still to audit

Five more POST listings declare a `limit` that nothing reads. None of them
asserts anything false today, because no conformance case exercises the limit
-- which is also why none of them was caught. Each needs the same question
asked of its provider:

| Recipe | Path | Suspected |
|---|---|---|
| Algolia | `/1/indexes/{index}/query` | `hitsPerPage` in the body, and a `params` string that is URL-encoded inside the JSON |
| Google Pub/Sub | `...:pull` | `maxMessages` in the body |
| Bill.com | `/api/v2/List/*.json` | `start` and `max`, inside a JSON document inside a form field |
| Adyen | `/v71/paymentMethods` | Not paged at all, in which case the declaration should go |
| AWS SQS | `/` | Depends which protocol the Recipe models |

## What Attio found

An identifier that is an object. A record's `id` is three UUIDs -- workspace,
object and record -- so `record.id === other.id` is never true, `record.id` as
a map key is `[object Object]`, and anything that logs or compares one gets
nonsense rather than an error. Two of the three are constant, and declaring
them at `id.workspace_id` produced a key literally spelled that way: resource
constants were the fifth mechanism in the runtime to treat a dotted name as a
key rather than a path, after route fields, a renamed identifier, a field's
`in`, and the comparator.

The reason it keeps happening is now stated as a rule rather than fixed a
sixth time: a dotted name is a path only where somebody remembered to make it
one. `recipe.IsPath` decides it in one place, and the exception is what makes
it worth having -- Dropbox names a field `.tag`, where the leading dot is part
of the name. A path is at least two segments and every one of them is a name,
so `.tag` is a key and `id.workspace_id` is a path.

The rest of Attio is ordinary but unusually rich in things worth catching:
every attribute value is an array even when single-valued, values are
versioned rather than replaced so the current one is the entry whose
`active_until` is null rather than `[0]`, timestamps carry nine fractional
digits that JavaScript's `Date` silently truncates to three, there is no GET
listing anywhere in the API, and failures are flat with the HTTP status
repeated in the body as a number.

## The unverified pagination note, and what checking it costs

135 pagination declarations across 56 Recipes carry a note saying the style
was never verified and asking somebody to check what the provider actually
calls its parameters. The note is honest and it has been sitting there; what
it does not say is what an unchecked declaration *does*.

With no style, a route pages as though by cursor and reads `limit` from the
query string. For a provider that calls it something else -- and most do --
nothing is read, the limit falls back to the route's default, the default is
larger than any fixture, and the first response holds the whole collection
with no next page. The paging loop runs once and passes. It is the same shape
as the body-pagination bug, and there are 135 of them.

Two down. Each cost more than a parameter rename, because pulling the thread
found things beside it.

**Pub/Sub**, against Google's own discovery document: `pageSize` and
`pageToken` in the query on both listings, `maxMessages` in the body for pull,
which is not a paged listing at all. Beside that: no list response carried
`nextPageToken`, which is the thing client code loops on, so the loop ended
after one page; and every subscription and topic went out with a bare `id`
that Google does not send, which reads far more like the identifier than the
resource path beside it does.

**Algolia**, against its own OpenAPI description: `hitsPerPage` and `page` in
the body for search, the same two in the query for the index listing, and
`cursor` for browse. Beside that, three more:

- There is no `GET /1/indexes/{index}`. Algolia serves only DELETE and POST
  there, and the Recipe answered a GET with a listing. No conformance case
  ever used it, which is how a route that would 404 against the real API sat
  there unnoticed. It is now `POST .../browse`, which is the real one.
- Algolia counts pages **from nought**. Read as one-based, a client asking for
  page 1 is handed page 0 again -- the same records twice, no error, and a
  loop that never reaches the end. `positionOf` carried a comment warning
  about losing or duplicating a record at every page break; being hard-coded
  one-based, it was doing exactly that for every provider that counts from
  nought. `pagination.first_page` says which.
- The response's `page` and `hitsPerPage` were **constants**, 0 and 20. A
  client that asked for page 3 was told it was on page 0, by the field whose
  entire purpose is to say where you are. `page_field` and `limit_field` echo
  what was actually served, falling back to the provider's defaults.

Stated gap rather than hidden: browse honours an incoming cursor but does not
send the next one, because `cursor_field` is Recipe-wide and search shares it
-- and search has no cursor at all, so declaring one would put a field on
every search response that Algolia never sends.

**Adyen**: no pagination at all, checked against its own OpenAPI 3.1. The
request for `/paymentMethods` takes no limit, offset, page or cursor of any
kind and the response has nothing to resume from, so the declared page size
claimed something the provider does not accept. Removed; 95 of the 345 list
routes here already declare no pagination, which is the idiom for a listing
that does not page.

**AWS SQS**: `MaxNumberOfMessages` in the body for a receive -- a batch size
rather than a page size, because there is no cursor and asking again is how
you get more -- and `MaxResults` with `NextToken` for `ListQueues`.

## The AWS Recipes model paths AWS does not serve

Found while checking SQS's pagination, and larger than the thing it was found
under.

The AWS JSON protocol is RPC over `POST /`, with the operation named in the
`X-Amz-Target` header. Cauldron routes on method and path, so all three of the
AWS JSON-protocol Recipes encode the operation in the path instead:

| Recipe | Paths it serves | What AWS serves |
|---|---|---|
| Secrets Manager | `/ListSecrets`, `/GetSecretValue`, `/DescribeSecret`, `/CreateSecret`, `/PutSecretValue` | `POST /` with `X-Amz-Target: secretsmanager.ListSecrets` and friends |
| DynamoDB | `POST /` for query, plus `/items/{id}` and `/tables` | `POST /` with `X-Amz-Target: DynamoDB_20120810.Query` and friends |
| SQS | `POST /` for receive, plus `/queues` and `/messages/{id}` | `POST /` with `X-Amz-Target: AmazonSQS.ReceiveMessage` and friends |

SQS's own header comment already says the operation belongs in `X-Amz-Target`,
so this is a known simplification rather than a mistake -- but it is not
stated where it matters, and its consequence is that an SDK or a hand-written
client that works against these Recipes is addressing URLs AWS does not have.
The convenience routes are worse than the pagination note, because a client
can be written entirely against them and be entirely wrong.

The fix is the shape `selects:` already took for GraphQL: one path, several
routes, disambiguated by something other than the path -- there by a word in
the query body, here by a header value. `dispatch_on: X-Amz-Target` with a
per-route value would let these three Recipes describe the real protocol, and
the three of them are enough to justify it.

**GitHub**: `per_page` and `page` in the query, counting from one, checked
against GitHub's own OpenAPI description. Stated gap: GitHub also advertises
the next page in a `Link` header and most client libraries follow that rather
than counting pages, and Cauldron does not send one -- so a client written
against the header sees no next page and stops after the first.

**Cloudflare**: `page` and `per_page` in the query, from one, with
`result_info` reporting where you are beside the total. Both of those are now
`page_field` and `limit_field` rather than being absent.

**Telnyx**: the parameter is a `deepObject`, so what goes on the wire is
`page[size]` and `page[number]` -- brackets and all -- and the case that
covered it sent `?limit=1`, which Telnyx does not accept. Worse underneath:
**`meta.next_page` does not exist.** Telnyx's meta is four integers --
`total_pages`, `total_results`, `page_number`, `page_size` -- and this Recipe
invented a cursor beside them that its own case then asserted. A client
reading `meta.next_page` against the real API gets undefined, and a loop that
reads undefined as "no more pages" stops on the first one.

That is the third invented field found by pulling this thread, after Pub/Sub's
`id` and Algolia's `GET /1/indexes/{index}`. The pattern is worth naming: a
Recipe that declares a paging mechanism nobody checked tends to have invented
the response half of it too, because both were written from the same guess.

**Datadog** is the argument against reading one endpoint and assuming the
rest. Three listings page three different ways: monitors by `page_size` and
`page`, dashboards by `count` and `start`, events by `page` alone with the size
fixed at a thousand. `start` is an offset on the dashboard endpoint -- "the
specific offset to use as the beginning of the returned response" -- and a
POSIX timestamp on the events endpoint. Declaring it as the cursor for both
would page one of them by reading a clock.

That endpoint also needed a third thing a parameter name can be. Empty falls
back to reading `limit`, a spelling names the provider's word, and `"-"` now
says the provider accepts no name at all -- Datadog fixes the event page at a
thousand, and honouring an invented size would let code that sends one work
here and be ignored in production. The idiom already existed for `field`,
`message_field` and `code_field`.

**Shippo** had three inventions in one Recipe:

- `GET /rates` does not exist. Shippo serves `/rates/{RateId}` for one and
  `/shipments/{ShipmentId}/rates` for the list, and two conformance cases used
  the bare path.
- There is no `count`. The envelope is `next`, `previous` and `results` and
  nothing else, and a case asserted the count that was declared beside them.
- `?limit=1`, again, where Shippo calls it `results`.

Stated gap rather than hidden: `next` and `previous` are full URLs in this API
-- the spec's own example is `baseurl?page=3&results=10` -- and Cauldron puts
a position there instead. A client that follows `next` as a URL, which is the
entire point of that shape, cannot follow this one. A `cursor_field` that can
render a URL would close it, and this is not the only provider shaped that way.

### On what keeps being found

Every Recipe checked so far has had something invented beside the paging: a
field the provider does not send, or a path it does not serve. Five now --
Pub/Sub's `id`, Algolia's GET listing, Telnyx's `meta.next_page`, Shippo's
`count` and its `GET /rates`. The paging note was the thread; the guess it
records was never confined to the paging.

The second lesson is about the cases. Four mutations survived their first
conformance case this round because the fixture could not tell the outcomes
apart -- one dashboard makes `count=1` and the default identical, two make an
offset of one indistinguishable from a page size of one. A case only checks a
parameter when the data would differ if the parameter were ignored.

**Stytch** had two routes that are not listings.

`GET /v1/sessions` is not paged at all: it takes `user_id` and nothing else,
and answers with every session that user has. A page size was declared that
Stytch does not accept.

The member listing is a search, and it is a POST.
`/v1/b2b/organizations/{organization_id}/members` exists in Stytch for POST
only -- that path *creates* a member -- so the GET listing declared here was a
route Stytch does not serve, and a conformance case used it. The real one is
`POST /v1/b2b/organizations/members/search`, paging by `cursor` and `limit` in
the body, with the next page at `results_metadata.next_cursor`, which nothing
was emitting.

Stated gap: the real search narrows by an `organization_ids` array in the
body, and Cauldron scopes a listing from the path, so this one answers with
every member rather than one organisation's. The scope declared before came
from a path parameter that does not exist.

That is the sixth phantom route or invented field, out of ten Recipes checked.

**LaunchDarkly**: `limit` and `offset` in the query. `limit` happened to be
the name the runtime already guessed, so the page size looked like it worked
-- but `offset` was not read at all while no style was declared, so every
request answered with the first page and a loop that asks for the next offset
never moved.

**Bitbucket**: `pagelen` and `page`, from one, and the envelope's own
definition names them. Its `page` and `pagelen` were declared as the constants
1 and 10, the same shape as Algolia's, so the response said page one whatever
was asked for. Also: the one case covering paging existed **twice, byte for
byte**, and both copies sent `?limit=1` and asserted only that `next` was
non-empty. A duplicate adds no coverage and hides how little there is.

**Xero**: `page` and `pageSize` for Invoices and Contacts, and **neither for
Accounts** -- three routes in one Recipe where the third takes no paging at
all. That distinction is only visible in the spec; a Recipe written from the
shape of its siblings would have declared a page size Xero does not accept.

### Where the specs run out

Reachability first still holds, and two providers failed it this round.
Pipedrive publishes an OpenAPI document that does not contain its own top-level
collections -- `/deals`, `/persons` and `/organizations` are absent, and only
nested forms like `/pipelines/{id}/deals` are described. Snyk's
`api.snyk.io/rest/openapi` answers 405 to a HEAD. Neither is modelled from a
guess; both stay on the list.

**Resend**: `limit`, `after` and `before` in the query. The cursor is a
record's own identifier rather than an opaque token, and it points both ways.

**Shortcut**: no pagination on any of the three. `/epics`, `/iterations` and
`/epics/{id}/stories` take no page size, no offset and no cursor, and answer
with bare arrays. Shortcut does page -- but only on its search endpoints,
which take `page_size` and `next` and are a different set of URLs entirely. A
page size declared on the listings was one they do not accept.

**DocuSign**: `count` and `start_position` for the envelope listing, and
**neither for recipients** -- two routes, two answers, the same split as
Xero's. `start_position` counts records, so reading it as a page number would
be wrong by a whole page at a time.

DocuSign also carries a stated gap rather than a fix. Its envelope reports the
position it served and the Recipe declares that as the constant `'0'`, which
is the same lie `page_field` and `limit_field` were added to correct for
Algolia and Bitbucket. It cannot use them yet: every numeric field in that
envelope travels as a string (`count_as_string`) and the echo fields emit
numbers, so switching would trade a stale value for a wrong type. Closing it
means teaching the echo fields the same string rendering the count already has.

**Ably** is the first row where the checked answer is that no style belongs
here at all. `limit` is the page size, but the next page is not a parameter:
Ably returns RFC 5988 `Link` headers, one per relative link, and a client
follows the one with `rel="next"`. Cauldron does not send them, so the page
size is honoured and the next page is unreachable.

That makes two providers whose real paging mechanism is a response header --
GitHub is the other -- and between them they are the argument for modelling
`Link`. It is a small feature with a clear shape: render the next page's URL
into a header rather than a body field, from the same position the cursor
already comes from.

**Documenso** pages its documents by `page` and `perPage`, from one. That part
is done. The rest of the Recipe needs a rework rather than a paging fix, and
tidying its pagination would have polished the surface of something that
should not exist:

- `GET /api/v1/documents/{id}/recipients` is **POST only** in Documenso -- it
  creates a recipient. Three conformance cases read it as a listing.
- `GET /api/v1/recipients/{id}` does not exist either. Documenso addresses a
  recipient at `/documents/{id}/recipients/{recipientId}`, and only for DELETE
  and PATCH. Two more cases use it.
- Recipients come back **embedded in the document**, which is where those five
  cases should be reading them.
- `count_field: totalPages` is the right name carrying the wrong quantity.
  Documenso's list envelope is `{documents, totalPages}` and nothing else, and
  `count_field` emits the number of records rather than the number of pages.
  Three documents at ten per page is one page, and this reports three. That is
  worse than an invented field, because the name is real and the value looks
  plausible. Closing it means a `pages_field` that divides rather than counts.

**Done.** All four, plus one the rework surfaced: the listing and the
single-document fetch answer with different shapes. Documenso's list item
carries the document's own fields and nothing nested, so code that lists
documents and reads `doc.recipients` gets undefined for every one of them and
concludes nobody has signed. That is a wrong answer with no error attached,
and it is now a case.

The five cases that had been written against the two routes Documenso does not
serve now read the document, which is where the provider puts recipients. The
recipient resource is gone: it is a list on the document, and the fixture seeds
it inline.

`totalPages` is back, with the quantity it claims. `pages_field` divides the
total by the page size and rounds up, so three documents at ten per page is one
page and at two per page is two. It had been removed rather than left wrong,
which was the right call for one tick and not a fix.

An empty set is nought pages. Providers differ about whether it is nought or
one, and nought is the reading that stops a loop rather than sending it after
a page with nothing in it.

**WordPress**: `page` and `per_page`, from one, read out of a live WordPress
REST route index -- the API describes its own arguments, which is as
authoritative as a spec gets. Stated gap: WordPress reports `X-WP-Total` and
`X-WP-TotalPages` in response headers and sends `Link` headers for the next
and previous pages, and Cauldron sends none of them. That is the third
provider whose paging metadata lives in headers, after GitHub and Ably.

**Clerk**: `limit` and `offset` in the query. `limit` was already the name the
runtime guessed, so the page size looked like it worked while `offset` did
nothing.

Clerk also had the envelope wrong, and it is the eighth thing found beside the
paging rather than in it. **Its listings do not share a shape.** `/v1/users`
and `/v1/sessions` answer with bare arrays; only `/v1/organizations` is
`{data, total_count}`. The Recipe wrapped all three, and two conformance cases
asserted that shape -- one of them named for it -- so code written here reads
`response.data.map(...)` and receives an array from the real API, where
`.data` is undefined.

Two of the three are bare, the envelope is Recipe-wide, so the majority shape
wins and the organisation listing is now the wrong one. That is stated in the
Recipe rather than hidden, and it is the same missing feature the AWS Recipes
need from the other direction: **a route that carries its own response shape.**
Three Recipes want it now.

**Trello**: no pagination on either. Neither a board's lists nor its cards
takes a limit, an offset or a cursor -- both answer with everything on the
board. Trello does page elsewhere, on its actions and search endpoints, which
take `limit` with `before` and `since`. A page size declared on these two was
one they do not accept.

**Greenhouse**: `per_page` and `page`, read out of its own Harvest
documentation, which states the same pair on every listing -- `per_page` an
integer between 1 and 500 defaulting to 100, and `page`, "the n-th chunk of
`per_page` objects". Stated gap: Greenhouse also sends a `link` response header
carrying `next` and `last`, and `skip_count` removes `last` from it.

That is the **fourth** provider whose next page lives in a header, after
GitHub, Ably and WordPress. Four is enough to stop calling it a curiosity.

### Two features three audits keep arriving at

Worth stating plainly, because they are no longer speculative:

~~**A route that carries its own response shape.**~~ **Built, for the envelope
half.** A route may now declare `list:` and override the Recipe-wide envelope:
empty inherits, `"-"` clears, and a boolean can only be turned on -- an unset
boolean and a false one are the same value in YAML, and guessing which was
meant is how a Recipe ends up asserting something nobody wrote.

Two Recipes stopped lying because of it. Clerk's organisation listing is
`{data, total_count}` again while its users and sessions stay bare, so the
majority-shape compromise is gone. Algolia's browse carries the cursor it
really has, without putting one on every search response, which is what the
Recipe-wide field would have done.

**The AWS half is built too.** `matches_header` on a route picks it by a
request header's value, which is the same bargain `selects` already makes for
GraphQL: one path, several routes, told apart by something that is not the
path. A route declaring it beats an equally-scoring route that declares
nothing, and an unmatched target is a 404 rather than a 405 -- a 405 would tell
a client to change the method, which was never the problem.

**Secrets Manager is converted.** Its five operations are `POST /` with
`X-Amz-Target`, and all eleven of its conformance cases moved with them. The
paths it used to serve -- `/ListSecrets`, `/GetSecretValue` and the rest --
now 404, which is what AWS does with them.

The comparison is exact, and that needed a case of its own. A mutation making
it a substring comparison passed every other case in the Recipe: nothing sent
a target that *contained* a modelled one. `secretsmanager.ListSecretsAndMore`
is not an AWS operation and is in the suite as a near miss, because a loose
comparison would route it to `ListSecrets` and answer it with somebody else's
secrets.

**DynamoDB is converted too.** Query, GetItem, ListTables and DescribeTable
are `POST /` with `X-Amz-Target`, and all ten of its cases moved. GetItem
needed one thing the format did not have: its key is
`{"Key": {"id": {"S": "ORDER#1001"}}}`, three levels down and in
attribute-value form, so `id_from: body:` now accepts a dotted name.

It also turned up an invented shape that nothing was asserting. **ListTables
answers with `TableNames` as an array of strings** -- just the names -- and
this Recipe emits an array of table objects under that name. A client doing
`TableNames.forEach(name => describe(name))` gets objects and calls
`describe([object Object])`. Cauldron has no way to say "a collection of
scalars": a listing is a list of records, and returning only the identifier
still emits `{"TableName": "orders"}` rather than `"orders"`. Stated in the
Recipe, and costed here.

**SQS is converted, and it took the Documenso treatment.** All five operations
are `POST /` with `X-Amz-Target` and all twelve cases moved. The route serving
`GET /messages/{id}` is gone: there is no fetch-one-message-by-id operation in
SQS, and the four cases that used it read a receive now, which is the only way
SQS hands a message over.

Three things it could not fix, stated in the Recipe's own header where a reader
meets them:

- `DeleteMessage` is keyed by the receipt handle from a receive, and Cauldron
  finds a record by its identifier and nothing else. Every delete answers
  NonExistentQueue. The case covering it asserts exactly that, which is true
  of an expired handle and not of a live one.
- `GetQueueAttributes` answers with a flat map of strings under `Attributes`
  and this answers with the queue object. Nesting the fields would fix the
  shape and break `ListQueues`, which shares the resource.
- `ListQueues` answers with `QueueUrls` as an array of strings -- the same
  shape `TableNames` needs and Cauldron cannot express.

### ~~A route that looks a record up by something other than its id~~ Built

`lookup_by` names the field the value from `id_from` is matched against.
`id_from` says where the value comes from; this says what it is compared with.

SQS's DeleteMessage is the case that wanted it. A receipt handle is
deliberately not a message id -- it is issued per receive, two consumers
holding two handles for the same message is normal, and a handle from an
earlier receive is stale. Looked up as though it were an id it found nothing,
so **every delete failed**, and the only thing that Recipe could show about
DeleteMessage was the way it fails.

A value matching nothing is left as it is, so a stale handle falls through to
the ordinary not-found, which is what SQS does with it. Both halves are cases
now: a live handle deletes the message it names, and a stale one deletes
nothing.

### ~~A collection of scalars~~ Built

`entry_field` on a route's `list` makes each entry that one field's value
rather than the whole record. It belongs on the route rather than the Recipe
because the listing and the fetch disagree by design: DynamoDB's ListTables
answers with the names and DescribeTable answers with the table, and SQS
splits ListQueues and GetQueueAttributes the same way.

Both were emitting objects where the provider sends strings, and neither had a
case saying so -- the only cases on those two routes were about authentication
failures. A client doing `TableNames.forEach(name => describe(name))` received
objects and called `describe([object Object])`: a request built from a string
that is not a name, with nothing to show it went wrong.


`pages_field` is also still owed, for Documenso's `totalPages`.

~~**Paging carried in a response header.**~~ **Built.** GitHub, Ably,
WordPress, Greenhouse and Buildkite all advertise the next page in `Link`, and
for Ably it is the only mechanism there is. `responses.list.link_header: true`
emits it: the request that was made, with the position parameter moved on, as
an absolute URL a client can request as it stands.

Two things the building turned up.

The URL has to be built from what the caller asked for, not from what the
sandbox sees. The multi-provider server mounts each Recipe under its own name
and rewrites `URL.Path`, so the first working version advertised
`/repos/octocat/hello-world/issues?page=2` when the server serves
`/github/repos/...` -- a next page that 404s, which is worse than no next page
at all. `RequestURI` survives the rewrite and is what it uses now.

And the absence of the header needed to become assertable. No `Link` on the
last page is the claim that ends the loop, and there was no way to write it
down, so `absent_headers` exists now. A first attempt at the case passed
against an emulator that sent `Link` on every page; the mutation that forced
that is what showed the case was checking nothing.

Only `next` is emitted. Providers also advertise `prev`, `first` and `last`,
and `last` needs a total this does not have.

### Sources that answer 200 without answering the question

The list is worth keeping because probing them again costs time:

| Provider | What happens |
|---|---|
| Segment | `api.segmentapis.com/docs/openapi.yaml` returns an HTML page |
| Tradier | `api.tradier.com/v1/openapi.json` returns 115KB of HTML |
| Pipedrive | Publishes OpenAPI that omits its own top-level collections |
| Snyk | `api.snyk.io/rest/openapi` answers 405 to a HEAD |
| Ably | Publishes the Control API, not the REST data API its Recipe models |
| ClickUp | `jsapi.apiary.io/apis/clickup20.source` answers 200 with Apiary's **Polls starter template** -- "a simple API allowing consumers to view polls and vote in them" -- 2KB describing a different API entirely |
| Snyk | `snyk.docs.apiary.io/api-description-document` answers 200 with a 435-byte sunset notice |
| Mailchimp | Documentation renders through JavaScript: stripped of markup, "count" and "offset" appear **zero** times |

The ClickUp entry is the one worth remembering. It is not an HTML shell and it
is not a 404 -- it is a valid API Blueprint document that parses cleanly and
describes somebody else's API. **A spec that parses is not a spec that
describes your provider**, and a tool that reads one without looking at it
would produce a confident Recipe for a polls service.

**WooCommerce**: `page` and `per_page`, from one, defaulting to one and ten.

**Buildkite**: the same two, defaulting to one and thirty with a hundred as the
cap -- but its documentation leads with the header, not the parameters: "the
pagination information can be found in the `Link` HTTP response header
containing zero or more of `next`, `prev`, `first` and `last`". The query
parameters work here and the header does not exist, which makes Buildkite the
**fifth**.

**Mailchimp cannot be sourced, and that is now a finding rather than a
deferral.** Its published JSON schemas describe responses only. Its
documentation page renders through JavaScript: fetched and stripped of markup,
the words "count" and "offset" appear in it **zero times**, because the
parameter table is not in the served HTML at all. It joins Klarna on the list
of portals that cannot be read without a browser. Recalling its parameters
would be exactly the guess this audit exists to remove.

**Postmark**: `count` and `offset` in the query, from its own published
Swagger at `postmarkapp.com/swagger/server.yml` -- count is "number of messages
to return per request, max 500" and offset is "number of messages to skip".
Stated gap: **both are required**. Postmark refuses the listing without them
and this answers with everything, so code that forgets them works here and
fails in production -- the same shape of lie as reading the wrong parameter
name, approached from the other side.

**Basecamp** was the best-documented of the night and needed the most of what
was built for it. `page` in the query and **no page-size parameter at all** --
the size is Basecamp's to choose -- so `limit_param: "-"` says so. Its own
words on the rest: "follow this convention to retrieve the next page of data,
please don't build the pagination URLs yourself", and "if the Link header is
blank, that's the last page".

Stated gap: Basecamp uses **geared pagination** -- fifteen results on page one,
thirty on page two, fifty on three, a hundred on four and above. The page size
changes as you go, which is exactly why building the URLs yourself is the thing
their documentation warns against: a client computing offsets from page numbers
is wrong from page two onwards. Cauldron serves one size, so it cannot make
that mistake and cannot teach it either. Page one's fifteen is the honest number
to declare.

Its fixture needed sixteen to-dos before a second page existed. Without one,
neither the page parameter nor the Link header was checked by anything -- the
mutation renaming `cursor_param` passed cleanly until the fixture was big
enough to have a page two.

### The fixture that cannot tell two parameters apart

Four times now a mutation has survived its first case because the fixture was
too small to distinguish the page size from the position: with two records,
skipping one returns a single record whether or not the limit was read.
Datadog, Resend, Documenso and Postmark each needed a third record before the
case checked what it claimed to.

It is worth stating as a rule. **A paging case needs at least three records:
one to skip, one to return, and one to prove the page size stopped it.**

### The count is in the tool now

`verify` reports it beside the evidence:

```
1636 of 1636 cases passed across 163 recipe(s).
5 checked against the real API, 1631 from documentation only
73 route(s) across 30 recipe(s) page by a parameter nobody named.
```

A route counts when it declares a page size and names neither a style nor a
parameter, because that is the state where the runtime reads `limit` -- right
for some providers, wrong for plenty, and invisible either way. A Recipe can be
entirely green and still be guessing: a case that never sends a page size
cannot notice the wrong name being read.

The number it produces is exactly the number this file had been tracking by
hand, which is the useful part -- it is the same debt, counted by the tool
rather than by a person, and it falls on its own as Recipes are checked.

### Remaining, now counted by the tool as well

That debt was tracked by hand at 32 Recipes and 78 declarations. `verify`
counts it now, the same way it counts the other two, and the real figure is far
larger because today's sweep added roughly two hundred declarations of its own:

```
47 paging parameter name(s) across 24 recipe(s) are declared and sent by no
case, so renaming them would break nothing.
```

A name counts when it is declared and no conformance case sends it to that
route -- in the query string, in the path's own query, or in the body for the
routes that page there. `"-"` is excluded, because it is a claim about absence
and a case cannot send a parameter that does not exist.

This is the request-side twin of a rule the format already has. Validation
refuses a *response* field name that no case asserts, on the grounds that it
could be renamed to anything unnoticed. The request half had no such rule and is
the half a client gets wrong: asserting the response to a listing says nothing
about which parameter produced it.

The method that works: find the provider's own machine-readable description,
read the parameter names out of it, then write a case that *sends* them.
Pub/Sub's `cursor_param` could be renamed to `cursor` with every case still
passing, until a case sent a token back and a second page came out of it.

Then 71 more, generated and then made to earn it. A route qualifies when it
declares both names, wraps its list, and already has a green case to copy the
request from -- path, headers and fixture -- so the new case differs from the old
one by exactly the two parameters. It sends a page size of one and a position
that lands on the second record, and asserts that record came back **on its own**.

The generation is not the interesting part; the acceptance is. Every candidate
had to survive three runs: pass as written, fail with `limit_param` renamed, and
fail with `cursor_param` renamed. 121 candidates went in and 71 came out. Thirty-
five did not pass at all -- Zendesk, Xero, PostHog, Postmark, Typeform and thirty
more, each for its own reason, and each now known to need a hand-written case
rather than a copied one. Three passed and then survived a rename: Airtable, Box
and Confluence declare a size that this fixture cannot prove is being read, which
is the fixture-too-small trap in a fourth, fifth and sixth form.

A second pass took 91 more, by widening the generator twice. Bare lists index
from the top -- `"[0].symbol"` rather than `"data[0].symbol"` -- which is 57
routes it had been skipping for no reason. And a fixture holding exactly two
records now grows a third: a copy of the second under a fresh identifier, with a
comment saying that is all it is. Two records had already established the
collection holds more than one, so a third claims nothing new about the provider,
which is the line that keeps this honest. A fixture holding *one* is not grown --
it may be holding one because the provider only ever has one, and inventing a
second is a claim about cardinality nobody checked.

That widening exposed a hole in the acceptance itself. Growing a fixture adds a
record every other case in the file can see, and a case asserting "these are the
two dashboards" is now wrong about a collection of three -- so the first run of
the second pass left thirty Recipes red while every new case passed. Checking
that the new case passes is not the same as checking the Recipe is still green,
and only the second one is the question worth asking. With that fixed, 30 more
candidates were dropped, and the ones that fell are the Recipes whose fixtures
carry a count somebody asserted on purpose.

### The listing nothing has ever answered

Chasing the paging debt found something underneath it. Of the routes whose
parameter names no case sends, the largest single reason is not a small fixture
or an awkward envelope: **218 of them are on listings no case answers at all.**
Nothing can be added to a request that does not exist.

`verify` counts that too now:

```
14 listing(s) across 11 recipe(s) have no case that answers them
successfully at all.
```

A listing counts as shown when some case asks for it, by that method and a
matching path, and expects 200. Every case touching the other 168 is checking a
failure -- a missing credential, a wrong method, an unknown path -- so the
collection, the envelope, the page size and the cursor field beside them are a
description with no evidence under it. The Recipe is green and has never seen the
endpoint work.

This was found once before, sideways. The rule that a response field name needs a
case asserting it turned up "two of them ... with no successful list case at all:
every case touching the collection was checking a failure". Two is what that
angle could see. Counted head-on it was 168.

46 of them now have one, built the same way as the paging cases: the credential
headers and the fixture come from a green case elsewhere in the same file, the
path comes from the route, and the assertion is the two records the fixture holds
in the order it holds them. Ten more went red and were withdrawn -- Chargebee,
DigitalOcean, GitLab, Gorgias, Greenhouse, Heroku, Hetzner, Recharge, WooCommerce
and WordPress -- each because the listing answers something the fixture alone does
not predict, which is the finding for the next person rather than a failure of
the sweep.

One bug is worth keeping. The first version asserted `[0].id` and Axiom answered
`missing`, because Axiom's dataset declares `id.field: name`: the store's
identifier is printed under another key. A generator that assumes `id` is called
`id` writes cases that fail on every Recipe that renames it, and there are
plenty.

### Two cases beat a third record

The three-record rule was right about one case and wrong about the problem. One
case cannot prove both halves at once -- with two records, asking for a page of
one *after* the first returns a single record whether or not the size was read.
Two cases can, and they need nothing the fixture does not already hold:

1. `limit=1` with no position. The first record comes back and the second does
   not, which can only happen if the size was read -- ignore the name and the
   whole collection arrives, because the whole collection is two.
2. `limit=1` at the position after the first record. The second record comes
   back, which can only happen if the position was read.

That is strictly better than growing the fixture, and it is the answer to the 23
clones that had to be withdrawn because a slug or an airport code cannot be
copied honestly. 152 cases across 76 Recipes went in on this pattern, each still
required to break under the rename it is named for -- the size case under
`limit_param`, the position case under `cursor_param`. The debt fell from 648 to
531.

A fixture of **one** can still show half of it, and it is the half nothing else
proves. The size cannot be shown -- asking for one out of one answers the same
whether the name was read or not -- but the position can: past the only record the
collection is empty, and a name the runtime ignored would have answered with the
record again. 95 more cases went in on that, each saying in its own comment which
half it shows and which it cannot.

A size with no position is provable on its own, which the generator had also
been treating as unbuildable. Amadeus takes `max` and no cursor at all: ask for
one out of several, receive one, and the name was read. There is no second case
because there is nowhere to page to, and the comment says so rather than leaving
a reader to wonder where the other half went.

And a template with no request body is not a reason to give up on a route whose
paging travels in one. Cohere's rerank is declared `in: body` on a GET and its
green case sends nothing, so the paging fields *are* the whole body -- which is
exactly what a client reading that declaration would send.

### Two more gaps that looked like provider behaviour

A diagnostic said 91 routes had a "position search that found nothing that
moved". It had found it. The generated case name was built from the *parameter*
-- "page moves the page on" -- so a Recipe with three listings that page the same
way produced the same name three times, and the second and third were dropped as
duplicates. Naming the case after the route settled 25 of them immediately.

And a route may name one half and refuse the other. Basecamp's messages take a
page number and no page size at all: `limit_param` is `"-"`. The generator
required both names to be real and skipped the route whole, so the one name it
does have stayed unsent. The size there is genuinely unprovable; the position is
not, and it is now a single case saying exactly that. 53 cases across 35 Recipes
came out of the two fixes together, and the debt fell from 189 to 122.

Both were the same mistake in different clothes: treating "the generator cannot
do this" as "there is nothing here to do".

### A correction, one commit later

The commit that added the one-record half-proof also claimed that the routes
still unsettled were "listings whose Recipe and fixture disagree about what comes
back". That was wrong, and it was wrong in the direction that blames the Recipes.

Probing them showed the disagreement was in the tooling. A route may override the
whole list *shape*, not only its key: Jira's field listing is `bare` where the
Recipe is `wrapped`, SES's account listing is bare and collapsed to a single
object, Jotform's questions are a `map`. Reading only the Recipe-wide style
asserted paths those responses never had. Nine of those routes settled as soon as
the override was read.

Worth stating plainly because the mistake is easy to repeat: when a generated
check disagrees with a file, the file is not automatically the one that is wrong.

### Stop guessing at the response; ask the sandbox

Every version of this generator guessed at what a listing would answer, and every
guess was wrong somewhere. `id` is not always called `id` on the wire -- Axiom
prints `name`. The fixture's type is not the response's type -- Shopify stores
`"632910392"` and emits `632910392`. And finally the fixture's *order* is not the
response's order, which is what Alpaca and Etsy were failing on: "want 902, got
904".

The sandbox knows all three. `cauldron serve --headless` mounts the Recipes,
`POST /_cauldron/<recipe>/seed?fixture=...` loads the fixture, and the listing
answers. Making the two requests and asserting exactly what came back removes
every one of those guesses at once. What is still being *claimed* is only that
the parameter names belong to the provider, and that is what the rename check
tests -- so the generator now claims one thing and observes everything else.

72 cases across 36 Recipes, and the debt fell from 315 to 245.

The same treatment on the unshown listings took them from 47 to 37, and it
exposed one more assumption that was about this tooling rather than about any
Recipe: a green case was only accepted as a template if it sent *headers*. Several
providers put the credential in the query string -- FRED, 511, the Guardian -- so
eight listings sat unshown because the script insisted on a header the provider
does not use. Carrying whatever query the green case sends fixes it, and is what
should have been copied all along. The ones that
still do not build are the honest remainder: a listing with fewer than two
records served, or a response whose records carry no scalar field to point at.

Then three widenings, and one guard that had to come with them.

Paging in the **request body** is now sent in the body -- Plaid's `options.count`
and `options.offset` nest inside the JSON, Algolia's `hitsPerPage` sits at its
top level, and both are POSTs whose whole request is the body, so reading a query
string there settled nothing. The **default and `envelope` shapes** are the same
`{"object": "list", "data": [...]}` the runtime emits from its default arm, and
treating an undeclared style as unaddressable had left 33 names unprovable on the
commonest shape in the collection. And the two-case pattern works for a fixture of
**any** size, not just exactly two: ask for one and get the first without the
second, then ask for one from after the first.

The guard is the interesting part. Run without it, that widening produced **234
cases and settled 18 names** -- 216 assertions repeating what the file already
proved, because nothing checked whether a parameter was already sent before
writing a case to send it. A suite that grows faster than its coverage is a suite
nobody will read. With the check in place the same run produced 29 cases and
settled 25 names.

That is the whole remaining shape of this debt, incidentally. Half of it was never
about missing evidence; it was about asking for evidence a two-record fixture
cannot give and then declining to ask for the half it can.

The lesson generalises past paging: when an assertion cannot distinguish two
causes, the fix is usually a second assertion rather than a bigger fixture.

A fourth relaxation, and the sharpest of them: **the fixture is not the answer
about a value's type.** Shopify writes its product id as the string
`"632910392"` in the fixture and declares the resource as `{style: numeric,
type: number}`, so the response carries a *number*. A case asserting the string
fails with "want the string 632910392, got the number 632910392" -- a message
that reads like the value is wrong when only its quoting is. Six of eight sampled
failures were exactly this, across Shopify, Zendesk, Recharge, PostHog, Rollbar
and Chargebee. The generator now renders each value the way the Recipe says the
response will carry it, reading `id.type`, `id.style` and the field's own
declared type rather than the shape it happens to have in YAML.

Three more relaxations took the unshown count from 122 to 67, and each was a
mistake in the generator rather than a limit of the Recipes.

**One record is enough.** The first version required two, because it asserted an
order and an order needs something to order. But the question a listing case
answers is "what does this endpoint return", and one record answers it -- 62
listings were being skipped for wanting a second row they did not need. The case
now asserts what the fixture holds and, with `absent`, that there is no more of
it than that.

**A scoped listing needs the scope its own record carries.** Akamai's activations
are scoped by `propertyId`, and filling the placeholder with the first property in
the fixture asks a different property for its activations, gets an empty list, and
fails an assertion for a reason that looks like an emulator bug and is a bad
request.

**A route may override the collection key.** Akamai's activation listing wraps
twice, at `activations.items`, while the Recipe-wide key is `items`. Reading the
Recipe-wide one asserts a path the response does not have.

And one that would have been much louder: appending a case to the end of the file
works only when `conformance` is the last top-level key. For Akamai it is not, and
a case block landing after `webhooks:` parses as a sequence where a mapping
belongs -- the Recipe stops loading at all. Cases are inserted at the end of the
conformance list now.

The two counts feed each other, which is the useful part. Showing 46 listings
gave the paging generator 46 new green requests to copy, and re-running it took
another 84 cases -- so the debt fell from 763 to 695 without reading a single new
document. The order matters: a paging case cannot be built on a request nobody
has ever made work.

One failure that only a Go test could catch. Twilio's fixture grew a third
message, and `TestTwilioScopesMessagesByAccount` asserted two. The generator's
acceptance runs `verify`, which runs conformance cases and not the Go suite, so a
fixture change is invisible to it. Worth noticing while fixing the number: every
message in that fixture belongs to one account, so the test shows the scoped
listing answers them and not that another account's are kept out -- it has
nothing to exclude.

### The clone that shared an email address

Growing a fixture by copying its last record leaves every field duplicated but
the identifier, and most of that is fine -- two things can share a name, a
status, a currency. A handful of fields identify the record itself, and there
two rows sharing one is a state the provider does not have: two HubSpot contacts
with one email, two Duffel airports with one IATA code.

The first attempt at fixing it inferred which fields those were, on the rule that
a field differing across every existing record is being treated as unique. With
two records that is nearly every descriptive field, and the rule withdrew 29 of
33 clones. The second attempt kept the inference and freshened the values, which
was worse: it produced `LAX-2`, `KLAX-2` and `Los Angeles International
Airport-2`. An invented value in the wrong shape is worse than the duplicate it
replaces, because a duplicate is at least a thing the provider could print.

What works is naming the fields rather than inferring them. An email is
freshened, because a local part takes a suffix without changing shape. A slug, a
URL, an airport code is not, and those nine clones were withdrawn whole along
with the cases that needed them -- Duffel, Discourse, Svix, Twitch, Unsplash,
CourtListener, incident.io, Papertrail and SavvyCal, each of which now waits for
a hand-written third record instead of a generated one.

One more bug worth recording, because it is the same shape as the go:embed one.
The freshener searched the file for `        email: ...` and edited the *last*
match, which was a conformance assertion indented by ten spaces -- an eight-space
search string is a substring of a ten-space line. It left the duplicate in the
fixture and rewrote the case to describe a record that no longer existed, and the
suite went red on exactly one Recipe. Scoping the search to the fixtures block
and anchoring it to a whole line fixed it.

The first run of that generator found nothing at all, and the reason is worth
recording: Recipes are `go:embed`ed, so `verify` was reading the compiled-in copy
and every edit was invisible. It reported "clean" for cases that were never
loaded and "the mutation survived" for mutations that were never applied. The fix
was to batch the work into three phases with one `go build` each -- write all,
rename all, rename all again -- which is also forty times faster than rebuilding
per case.

Stripe is the first paid down. Its `starting_after` was declared on two routes
and sent by nothing, so it could have been renamed to `cursor` -- the name the
runtime reads when none is given -- and all thirteen cases would still have
passed. Mutating it now fails two. That needed a third customer and two more
payment intents first, because of the rule three sections up: **a paging case
needs at least three records: one to skip, one to return, and one to prove the
page size stopped it.**

## Neon, and a create that answers with unfinished work

Neon is serverless Postgres with branching, and almost nothing is finished when
the response arrives. Creating a branch answers **201 with the branch in state
`init`**, which Neon's own description defines as "being created but is not
available for querying". Code that creates a branch and connects to it on the
next line connects to nothing, and the status code gave it no reason to wait.

The response carries the operations it started, and they are the only thing in
that body that says when the branch becomes usable -- sitting beside the branch
a caller came for. That needed one change: route `fields` applied to listings
only, so a create could answer with the record and nothing else. A response
carrying just the record would be the helpful kind of wrong, because it looks
finished.

Four more things worth catching, all from Neon's published description:

- `current_state` and `pending_state` answer different questions and disagree
  exactly while something is happening. Neither says "is it ready" alone; the
  settled branch is the one with no `pending_state` at all.
- An operation's status has **eight** values and `finished` is the only
  success. `failed`, `error`, `cancelled` and `skipped` are all "stopped, and
  not because it worked", and `scheduling` has not started -- so a poll written
  as `while status === 'running'` exits immediately and calls it done.
- `default` and `primary` are both on a branch, `primary` is deprecated, and
  they can point at different branches. Code reading the deprecated one gets a
  branch rather than an error.
- The create response carries the connection string **with the password in
  it**, in a body that gets logged.

### Owed: a request envelope

Neon takes a create body as `{"branch": {...}}` and Cauldron reads it flat, so
a client sending Neon's real shape is told the name is missing. The Recipe's
cases send the flat form, which is the one thing in it that does not match the
provider. Plenty of APIs wrap a create body this way.

## Supabase, and two things it could not say

Supabase's management API has four traps worth reproducing, and the Recipe
reproduces three of them.

**Status has fifteen values and none of them is `ACTIVE`.** A working project
says `ACTIVE_HEALTHY`, so `status === 'ACTIVE'` is never true and the branch
guarding on it never runs. Worse in the other direction:
`status.startsWith('ACTIVE')` is true for `ACTIVE_UNHEALTHY`, which is a
project that does not work -- a check that passes about a thing that is broken.

**A created project has no database.** The listing carries a database object
with the host in it; the create response omits it entirely, because there is
nothing to connect to yet. Code that creates a project and reads
`project.database.host` reads a property of undefined. Both halves are cases,
because the absence alone would hold against an emulator that never sends a
database at all.

**Secrets come back with their values**, in a body that gets logged.

### Owed: a second identifier

A project has two, and only one works. `ref` is what every path takes; `id` is
marked "Deprecated: Use `ref` instead" in Supabase's own description and is
still sent beside it, so code that stores `id` and interpolates it addresses
nothing.

Only `ref` is modelled. A resource keys on one identifier and emits it under
one name, so a second cannot sit beside it -- and of the two, the one that
works is the one worth having. The trap is described in the Recipe's header
rather than reproduced, which is the weaker half of what it could do. Neon has
the same shape with `default` and `primary`, and could model both only because
neither of them is the identifier.

### ~~Owed: lookup_by should not fall through to an id~~ Fixed

Found while modelling the above. `lookup_by` returns the value unchanged when
it matches nothing, so the ordinary not-found path can report it -- which is
right for SQS, where a stale receipt handle matches no identifier either. It is
wrong when the value *is* a valid identifier: a route that says "look this up
by ref" then answers 200 for a project addressed by its deprecated `id`, which
is the exact failure the route was declared to prevent.

Made strict, and SQS's behaviour is unchanged because a handle is never an id.
What did change is the direction nobody had tried: **a message id in the
ReceiptHandle field used to answer 200 and delete the message**, because a
value matching no handle fell back to an identifier lookup and a message id is
one.

So the emulator was teaching precisely what SQS's own Recipe header warns
against -- that a handle and an id are interchangeable. Measured before the
fix: four messages, delete by id, 200, three messages left.

The value still comes back from the lookup so the failure names what the
caller asked about; only its usability changes.

## Discourse

Its shapes are the whole reason to model it.

**The topics are two levels down.** `/latest.json` answers with
`topic_list.topics`, and `per_page` is inside that object rather than beside
it, so a client reading `response.topics` or `response.per_page` finds nothing
and no error.

**The users arrive beside the topics, not inside them.** A topic carries
`last_poster_username`, a string; the object with the name and the avatar is in
a separate array in the same body. Rendering one row means joining two arrays
that arrived together.

**Authentication is two headers.** `Api-Key` and `Api-Username` are both
marked required in Discourse's own description, and a key alone is refused with
a 403 rather than a 401. Every other API here takes one credential, which is
why this is the failure a Discourse integration hits first.

**`title` and `fancy_title` are different strings** -- the same words, rendered,
so comparing `fancy_title` to anything a human typed does not match. And
`bumped_at` moves when anything happens while `last_posted_at` moves only for a
post, so ordering by one is not ordering by the other.

Not claimed: `avatar_template` is a template with a pixel size substituted into
it rather than a URL, which is a real trap and one the published description
does not evidence -- it carries `type: string` and no example. The fixture
holds a realistic value and no case asserts the substitution, because there is
nothing to cite for it.

## Langfuse, and a deprecation you can only see in the body

Langfuse announces a dying endpoint in the response, not the status code. A
deprecated endpoint answers **200 with everything you asked for**, and a
`_deprecation` object beside it: a message, the replacement endpoint, a link to
the migration guide, and **`sunsetAt`** -- the date after which it may stop
working. Nothing about the status, the shape or the data says anything is
wrong. A client reads the payload, works perfectly, and stops working on a date
nobody read.

Three of its endpoints are already there: `GET /traces`, `GET /observations`
and `GET /v2/scores`, replaced by `/v2/observations` and `/v3/scores`. Both
halves are cases, because "the deprecated one carries a signal" holds just as
well against an emulator that stamps one on everything.

**Migrating is not a change of URL.** The deprecated observation listing pages
by `page` and `limit` and reports `meta.limit`, `meta.page`, `meta.totalItems`
and `meta.totalPages`. The replacement pages by a **cursor** and its `meta` is
`{cursor}` and nothing else. Code that showed "page 3 of 12" has nothing left
to show it with, and code that keeps sending `page` gets the first page every
time.

That is the first Recipe to need a per-route envelope for a reason other than
inconsistency: the two endpoints disagree **on purpose**, because one replaces
the other.

Also modelled: a trace's `observations` is an array of ids rather than objects,
and only a `GENERATION` observation carries a model or a cost -- a `SPAN`
carries null for both, so reading `model` off every observation gets null for
two kinds out of three.

## Squarespace, and an absence that was two absences

Every Recipe here answered a request for a record that is not there with a 404,
because an identifier had a shape the emulator **minted** with and no shape it
**checked** against. Squarespace documents both answers on one route: `GET
/1.0/commerce/orders/{id}` is 404 "The requested Order was not found" for an
id that could exist and does not, and 400 "The id is not in the expected
format" for one that could not exist at all.

The distinction is not decorative, and it does not fail in a direction anybody
notices. A 404 is a fact about the account -- the order was deleted, or belongs
to somebody else -- and retrying will not help. A 400 is a fact about the
caller: an id from the wrong provider, a truncated string, an empty variable
interpolated into the path. An emulator that collapses them teaches an
application to log its own bugs as missing data. Worse, the id a test reaches
for when it wants a miss is `"nonexistent"` -- which is exactly the id that
does **not** behave the way the test assumes.

`id.pattern` closes it. A resource declaring one is refused before the store is
consulted, which is the order the providers who do this run it in: the id never
reached anything that could have looked for it. Three rules keep the
declaration from being decorative -- the pattern has to be anchored at both
ends, every fixture id has to match it, and some conformance case has to
address the resource with an identifier it rejects. Without the third, a
pattern nothing exercises looks exactly like a provider that does not check,
because both answer 404 to everything a case is likely to ask for.

Squarespace is the first user. Stripe, Intercom and everything else built on
ObjectIds behave the same way and are not yet declared, so this is a gap that
is now expressible rather than a gap that is closed.

## Two Squarespace rules that are documented and not served

Both are refusals keyed on a **combination** of parameters rather than on any
one value, and route selection here works on the values a request carries.

The cursor may not travel with anything else -- "Cannot be used with other
parameters", so page two of a filtered query is not expressible and the
ordinary pager that keeps its query string and appends `&cursor=...` gets a 400
on its second call and none on its first. And the date range is all or nothing:
`modifiedAfter` is "Required when `modifiedBefore` is passed" and the reverse,
so half a range is a refusal rather than an open interval.

Declaring either against particular fixture values would serve one pair of
strings while claiming a rule. `matches_query` compares values, not
combinations, and a route that refuses *whenever two parameters appear
together* is the shape neither it nor `selects` nor `matches_header` can say.

Also documented and not served: the `order.create` webhook for a Payment Plan
order "fires only when paymentState transitions to PAID (all installments
collected), not when the order is initially placed", and "No webhook fires for
the deposit or for intermediate installment captures". Nothing in the public
API moves an order's payment state, so there is no request an emulator could
answer with that silence. The other half of the same trap **is** served: the
default `paymentStates` filter leaves out `PARTIALLY_PAID`, so the order is
invisible to the listing for as long as the plan is collecting.

## Four commerce providers checked and not written: three unreadable, one unservable

The commerce and billing section is nearly worked through. What is left was
checked against the bar this project has been holding -- only write a Recipe
when the vendor publishes a **response-bearing** description -- and none of the
four clears it. The reasons differ and are worth keeping, because three of them
could change and one cannot.

**Faire.** `faire.github.io/external-api-v2-docs`, the address this collection
recorded, now serves a meta-refresh to `developers.faire.com/docs`, and that
page is a client-rendered application: the only scripts are
`cdn.faire.com/static/js/external-developers.*.js` and a cookie banner. No
OpenAPI at any conventional path, and the GitHub organisation has no docs
repository -- its twenty public repos are Kotlin and JavaScript tooling
(`mjml-react`, `yawn`, `faire-detekt-rules`), which is the shape Squarespace
and Wix both have: a provider publishing a great deal, none of it a client.

**Maxio**, the former Chargify. `developers.maxio.com` answers 121 bytes to
every path including `/openapi.json`. `maxio-com/ab-api-docs` is gone. The
generated `ab-java-sdk` carries a `.codegenignore`, which says a spec exists
somewhere, and does not carry the spec. Reconstructing response shapes from
generated Java model classes is reconstruction rather than reading, which is
the distinction the Katana row drew.

**Zuora.** The documented spec address, `developer.zuora.com/api-references
/api/spec/openapi.yaml`, is a 404 that answers with 229 kilobytes of
application shell. Same failure as Maxio, and the row's original objection --
"very large, and the object model predates REST conventions" -- still stands
behind it.

**PrestaShop is different, and its reason will not change.** The evidence is
not the problem: the webservice lives in the open-source tree and can be read.
The problem is that the interesting thing about it is not servable here. The
API answers **XML** unless a request asks for JSON with `output_format=JSON`,
so the headline -- a REST API whose default representation is not the one
every client assumes -- is a claim about a content type this emulator does not
produce. A Recipe could serve PrestaShop's JSON half faithfully and would have
to leave out the only reason to write it.

Three of these unblock the moment a spec becomes reachable. The fourth would
need Cauldron to serve XML, which is a larger decision than one Recipe.

## Gemini, and a response shape chosen by what you asked for

Written up rather than written, because the finding is real and the format
cannot serve it.

Of `promptFeedback.blockReason`, in the discovery document Google publishes at
`generativelanguage.googleapis.com/$discovery/rest?version=v1beta`:

> "If set, the prompt was blocked and no candidates are returned. Rephrase the
> prompt."

So a blocked prompt is **a 200 with the candidates taken away**. The response
still carries a `responseId`, a `modelVersion` and `usageMetadata`; the
`candidates` array, which is the only part anybody reads, is simply not there.
Code written the obvious way does `candidates[0].content.parts[0].text` and
throws on the first index, and the sentence explaining why sits in a sibling
object most integrations have never opened.

That is worth having because it differs from the relative already recorded.
The OpenAI Recipe's refusal arrives as a normal message with `content: null`
and the reason in `message.refusal`, so the read **succeeds and yields
nothing**. Gemini removes the element, so the read **throws**. One provider
hands you a null and the other an exception, for the same event, and neither
says anything in the status code.

**Why it is not a Recipe.** Demonstrating it needs one path to answer two
shapes depending on what the request asked for, and the only mechanism here
that reads a request body is `selects`, which is fed `graphQLQuery(r)` — the
GraphQL `query` field and nothing else. A create route echoes the request; a
route constant is fixed per route; `matches_query` and `matches_header` read
places Gemini does not use. So the choice cannot be made.

**This has since been done and the Recipe ships.** `selects_body` was added as
a new field beside `selects` rather than as a change to it: seven GraphQL
Recipes depend on the current behaviour, and a marker word that matches only
inside a query today would have started matching variables and arguments too.
What is above is left as it was written, because the reasoning that identified
the missing mechanism is the useful part of it.

Three more from the same document, for whoever writes it:

- **An empty finish reason means it has not finished.** Of `finishReason`:
  "If empty, the model has not stopped generating tokens" — a field whose
  absence means work in progress on the streaming endpoint and nothing at all
  on the other one, from one schema serving both.
- **The total counts thinking you never see.** `totalTokenCount` is "prompt +
  thoughts + response candidates", with `thoughtsTokenCount` beside it. That
  is OpenAI's reasoning-token bargain; what is Gemini's own is that the
  arithmetic is spelled out in the field's own description.
- **Being cached does not make a prompt smaller.** Of `promptTokenCount`:
  "When `cached_content` is set, this is still the total effective prompt
  size", so the number a cost model reads does not fall when caching works.

## Mixpanel, and an API that only speaks in batches

Assessed against the OpenAPI document Mixpanel publishes and its own
documentation site renders from, at `mixpanel/docs`,
`openapi/ingestion.openapi.yaml` -- 68 kilobytes, and the source of every
quotation below. Written up rather than written, for the reason at the end.

**The endpoint everybody uses answers with the number 1, and Mixpanel says
plainly that it does not mean what you think.** Of the `/track` response, which
is `text/plain` carrying an integer:

> "`1` - One or more objects provided are valid. **This does not signify a
> valid project token or secret.**"
>
> "`0` - No data objects in the body are valid."

Three separate things are in there.

The first is that the answer to a whole batch is **one bit**, and the bit is an
OR. Send fifty events, have forty-nine rejected, and the response is `1`.
Nothing in it says forty-nine, or which, or why.

The second is the sentence Mixpanel took the trouble to write: a `1` **does not
signify a valid project token**. `/track` carries `security: - {}` in the spec
-- no HTTP authentication at all -- because the credential is a field inside
each event, `properties.token`. So a request with a wrong token, an expired
token or no token is not a 401. It is a `1`, and the events are gone.

The third is that you have to opt in to being told anything. Of `verbose`:

> "If present and equal to 1, Mixpanel will respond with a JSON Object
> describing the success or failure of the tracking call."

The default is the bit. The diagnosis is a query parameter, and its own
description recommends it only "for debugging during implementation".

**A 400 from `/import` can have imported 999 records.** That is Mixpanel's own
example, verbatim, in the `StrictInvalid` response:

```
code: 400
num_records_imported: 999
status: Bad Request
failed_records:
- index: 0
  insert_id: 13c0b661-f48b-51cd-ba54-97c5999169c0
  field: properties.time
  message: "'properties.time' is invalid: must be specified as seconds since epoch"
```

The batch is not atomic and the status code says nothing about how much of it
landed. Retrying the request on a 400 -- which is what a client does when it
reads a 4xx as "this did not happen" -- re-sends 999 records that did.

And the failures name records by `index`: their **position in the array you
sent**. Not by id, because an event has no id until Mixpanel assigns one. So
interpreting the response requires still holding the request body, which a
client that streamed or freed its batch no longer does.

**Twenty paths, six URLs.** The document's path keys are `/engage#profile
-set`, `/engage#profile-set-once`, `/engage#profile-numerical-add`,
`/engage#profile-union`, `/engage#profile-list-append`,
`/engage#profile-list-remove`, `/engage#profile-unset`,
`/engage#profile-batch-update`, `/engage#profile-delete`, and the same trick
again seven times for `/groups`.

There is one `/engage` endpoint. It does nine different things depending on
which `$`-prefixed key is present in the body -- `$set`, `$set_once`, `$add`,
`$union`, `$append`, `$remove`, `$unset`, `$delete` -- and OpenAPI requires
path keys to be unique, so the fragments are there to make nine descriptions of
one URL fit in a document that has no way to say "same path, different body".

So sixteen of the document's twenty path keys describe two endpoints, and the
whole file describes six: `/track`, `/import`, `/engage`, `/groups`,
`/lookup-tables` and `/lookup-tables/{id}`.

A URL fragment is never sent over the wire, so a client generated from this
document requests `/engage` and works by accident. Every tool that reads the
document for inventory -- a gateway, a catalogue, a coverage report -- counts
twenty endpoints where there are six.

**`status` is three different things in one document.** In `verbose` mode it is
"the value 1 on success and 0 on failure", a number. In `StrictReceived` it is
the string `"OK"`, and in `StrictUnauthorized` the string `"Unauthorized"` --
HTTP reason phrases. In the shared `ErrorResponse` schema, which `/track`'s own
401 and 403 use, it is an enum of exactly one value, `"error"`. So the two 401s
in this document disagree with each other about what `status` says, and the
success field of the same name is a different type again.

### Why it is not written

**No endpoint in this API accepts a single JSON object.** `/track`, `/import`,
all nine `/engage` operations and all seven `/groups` ones -- eighteen of the
twenty path keys -- declare `requestBody: schema: type: array, minItems: 1`.
The two that do not are the lookup tables, which take `text/csv`. There is no
single-record form of anything.

This format models one record per request, and a conformance case's request
body is a mapping. Serving Mixpanel would mean accepting a single event object
at `/track` -- a request shape Mixpanel does not accept -- which is the
helpful kind of wrong this collection refuses: a suite that passes locally
against a body the provider would reject.

PostHog is the near neighbour that shows the line. Its Recipe ships and its
`/capture/` case pins the same `{"status": 1}` shape, because PostHog has a
single-event endpoint and a separate `/batch/` one. Mixpanel has only the
batch.

The second gap is smaller and real: a `text/plain` body containing `1` is not
an object, and every response this format emits is. Even ignoring the content
type -- `1` is valid JSON -- there is no way to say "this route answers with a
constant rather than a record".

**What would unblock it is a category, not a Recipe.** Array request bodies
would reach Mixpanel, Amplitude, Heap, Segment's batch endpoint and Customer.io,
all of which are queued above and all of which are batch-shaped for the same
reason: an analytics SDK buffers on the client and flushes. That is a decision
about what this format models, which is larger than one provider and should be
made on its own evidence rather than as a side effect of wanting this one.

**Resolved, 2026-09-03, on that evidence.** A conformance case's `json` is now
any, so it accepts an array. Mixpanel, Amplitude and Heap all ship. Serving the
array rather than describing it corrected the finding above: an array is not
accepted for being an array. Drop the token, the event name or the properties
and a well-formed array still answers zero, so the gate is the presence of three
fields and never their validity -- a token belonging to nobody passes, and a
missing one does not. That distinction was not visible while the behaviour lived
in prose.

The second gap named above stands. A body that is the single character `1` is
still not something a route can answer with, and Mixpanel's is served through
`raw:` for that reason.

## Five providers want a credential failure per route, and one wants it per scheme

`auth` describes one credential and one way of refusing it. That is right for
most of this collection, and five Recipes have now had to work around it the
same way -- `scheme: none`, plus a `matches_header` pair on every route, so the
credential is checked once per declared route rather than once for the Recipe.

The workaround is honest and it is expensive. Every route gains a duplicate,
declaration order stops being obvious, and a reader has to reconstruct the
credential rule from five places instead of reading it in one.

What each of them actually needs:

| Recipe | What it wants that `auth` cannot say |
|---|---|
| Flagsmith | A **different status and body per path** for the same wrong key. Three failures in three formats -- a printed Python tuple served as JSON, a zero-byte 404, and a bare array holding one object -- and which one you get depends on the endpoint, not on the credential |
| Anvil | A **different answer per scheme**. A wrong Basic credential is a GraphQL `{"errors":[...]}` envelope; a wrong Bearer one is a singular `{"error":...}` with a `WWW-Authenticate` header the Basic failure never carries. Same host, same request, two shapes chosen by the word in front of the credential |
| Baserow | Both at once, and a **third ordering** besides: URL pattern, then credential, then the specific id and the method. Its two live prefixes, `Token` and `JWT`, answer different error codes |
| Backblaze | **Two protocols on one vendor**, a JSON native API and an S3 gateway, whose credential failures share no vocabulary. `scheme: none` with per-route `matches_header` is the only way to gate both |
| Make | Four answers, three verdicts. Fixed for three of them by naming an error per verdict; the fourth -- `Bearer` with a uuid, "Invalid bearer token." -- is a **scheme Make recognises and refuses**, which is a fourth thing to know about a credential rather than a fourth way of failing |

Two shapes would cover all five: a credential failure declared on a route,
overriding the Recipe-wide one, and a verdict that can depend on which scheme
was presented. Neither is written yet. What is written is that each of these
files says plainly which of its provider's answers it is serving and which it is
not, which is the thing that must not be lost if this is ever built.

### One green case is not the credential, it is *a* credential

Seven of the listings nothing had ever answered were not blocked by the
provider or by the Recipe. They were blocked by the generator taking the first
green case in the file and treating it as the file's credential.

It usually is not. The first green case in a Recipe is very often a check
against the empty fixture -- the listing answers, the collection is empty, the
case is about the envelope -- so the generator seeded `empty`, read back `[]`,
found no record to assert and gave up. And a Recipe with more than one accepted
credential presents each of them on different routes, so the header it copied
came back 401 from the route it was trying to show.

Trying every green case in turn, best fixture for that route's resource first,
answered five of them. The other two were a different mistake: the record inside
the collection is sometimes a wrapper around the record. Chargebee's customers
listing answers a list of one-key objects, and looking for a scalar at the top
of one found nothing to assert. One step further in finds the id.

Both are the same lesson as the sandbox rewrite. The generator was not wrong
about the provider; it was wrong about the file, and the file was right there
to read.

### The first field in a record is not always a distinguishing one

Every Shippo transaction in the fixture carries the same `object_created`, and
the generator asserted it. The case passed. It also passed with `cursor_param`
renamed to something no client sends -- the page had moved, and the assertion
could not tell -- so the mutation check dropped the pair.

That was the right outcome for a case proving nothing, but the case did not have
to be that weak. Some other field on the same record does tell the two pages
apart, and asking for one that differs from the record on the first page rather
than the first scalar that appears settles eleven more names, Shippo's three
among them.

A generated case that survives its own mutation is not evidence. Choosing what
to assert is part of writing the evidence, not a detail after it.

### The record inside the record

Nineteen more names settled, and none of them needed anything from a provider.
Chargebee's list is a list of one-key objects, Cloudtalk's calls arrive under
`Cdr`, Contentful's assets keep everything under `fields`, Crowdin's languages
under `data`. The generator looked for a scalar at the top of the served record,
found an object, and passed the whole route over as unassertable.

One step further in finds the id. The same fix had already been made in the
generator that shows unanswered listings, a commit earlier, and not carrying it
across cost five Recipes their paging evidence for a week.

A cursor is also not always called `id`. Cloudflare Stream's videos carry a
`uid`, so scanning the fixture for ids returned nothing to try and the route
looked unpageable when it pages perfectly well.

### "The sandbox serves nothing" was the wrong fixture, again

Fourteen routes reported that the sandbox served no records for them. Several
had a fixture full of records; the generator was seeding the wrong one.

A route usually has more than one green case, and the first in the file is
frequently the one that checks the empty answer -- the listing responds, the
collection is empty, the case is about the envelope. Copying that case's request
copies its fixture too. Trying every green case for the route, fullest fixture
for that route's resource first, settles five more names.

This is the third time the same mistake has been found in a different generator.
The lesson is not about fixtures. It is that "the first thing in the file that
matches" is a guess wearing the clothes of a lookup.

### A record is not always an object, and a route is not always a path

DynamoDB's ListTables answers a list of table names. SQS's ListQueues answers a
list of URLs. The generator looked for a field on a string, found none, and
wrote the routes off -- when the assertion it wanted was simply on the element:
`QueueUrls[0]` is the queue.

Both Recipes then produced cases named "MaxResults is read on , and a page of
one is not the whole collection". AWS mounts every operation on `/` and tells
them apart by an `X-Amz-Target` header, so the path contributed nothing to the
name, two cases in a file collided, and the whole batch was dropped without a
word -- the case names never appeared in `verify`'s output, so the check that
looks for them found nothing to look at. Naming those cases after the operation
in the header fixed six of them at once.

Two more small things fell out on the way: a position is best read off the page
the sandbox actually served rather than off the fixture (Workable's jobs carry a
`shortcode` where the generator wanted an `id`), and a generator that drops
work silently is worse than one that fails loudly. The silent drop hid a
one-line naming bug for three sweeps.

### The request was right; the fixture it came with was not

Thirteen routes reported that the sandbox served no records. For six of them the
Recipe had a fixture full of records and the generator never seeded it, because
the only green case for the route was written against `empty` -- the listing
answers, the collection is empty, the case is about the envelope.

The request in that case is still the right request: a credential the Recipe
accepts, a path that resolves. Only the fixture is wrong. Seeding a fuller one
and sending the same request settles six more names.

That leaves seven routes whose Recipe genuinely holds no record of that
resource in any fixture -- Clio's matters, Frontegg's users, Kinde's users,
Make's scenarios, Hardcover's editions, HERE's geocode results. Those need a
fixture written, not a generator taught, and writing provider-shaped fixture
data from documentation is exactly the kind of guess this whole line of work has
been removing. They stay counted.

### The credential was in the Recipe all along

The generator that shows an unanswered listing copies its request from a green
case elsewhere in the file, because a path and a credential are two things it
cannot invent. That is true of the path. It was never true of the credential.

A Recipe's `auth` block names the scheme, the header, the prefix and the keys it
accepts. Building a credential from that is reading the file, not inventing a
secret -- and it works where copying does not, because several Recipes
authenticate every route while their only green cases are the two or three that
need no credential at all: Toggl's status endpoint, Navitia's root.

Two other small things came out of the same pass. Heroku answers **206 Partial
Content** to a listing, which is what its Range paging does and what the Recipe
models; insisting on 200 had hidden four listings behind a status code that was
never wrong. And a route whose path ends in a slash is not the route without
one -- Sentry writes `/organizations/{organization}/projects/` and answers 404
to the trimmed version, which is a real difference and worth reproducing.

### The counter was asking for evidence six Recipes already had

`UnshownListing` treated a listing as shown only by a case expecting 200.
Heroku answers **206 Partial Content** to a listing, because its paging is a
Range header and a partial answer is what that means. Agora answers **201** to
several of its reads. Both are shown, plainly, by cases already in their files.

Six listings were counted as unanswered on that basis, and the generator sent
after them did what it was told: it asked the sandbox, got an answer, and wrote
a case that already existed under a name one word longer. Two Recipes briefly
carried the same evidence twice.

The fix is four characters wide -- any 2xx, not 200 -- and it is the same
mistake as every other one in this section, at the other end of the pipe. A
measurement that is slightly wrong does not produce slightly wrong work. It
produces confident work on a problem nobody has.

Ten of the twenty-one that remain are in Recipes whose fixtures hold no record
of that resource at all: Dub's links, Melio's bills, payments and vendors,
Rippling's employees, FullStory's users and sessions, Transitland's departures.
Every one of those has exactly one fixture, `empty`, and nothing to serve. They
need fixture data written from the provider's documentation, which is the guess
this line of work has spent a fortnight removing, so they stay counted.

### One name was credited by a request the provider refuses

`UnsentPagingParam` asked whether any case sends the parameter. It did not ask
whether that case succeeds. Gumroad's `page_key` was sent by a case expecting a
400, which is no evidence at all: the refusal happens before anything reads the
paging, so the case stays green under any parameter name you like -- which is
exactly the condition this counter exists to find.

One name across the whole set, and both halves of the fix are two lines. It is
recorded because the interesting part is not the size. Every counter in this
file has now been wrong once in the same direction: too generous about what
counts as evidence. `UnstatedPagination` was right and got the reasoning wrong
in the docs; `UnshownListing` refused a 206; this one accepted a 400. A counter
is a claim about the corpus, and it deserves the same red-green treatment as
the code it measures.

### A collection of one is sometimes not a collection

Tradier answers an array of orders when there are several and a single object
when there is one -- which is precisely what asking for a page of one produces.
The generator asked for one, received the object it had asked for, could not
read it as a list and moved on. `collapse_single` is declared right there in the
Recipe; nothing had to be guessed, only read.

The assertion for a collapsed page has no index in it: `orders.order.id` rather
than `orders.order[0].id`, with `orders.order[1]` still asserted absent. Both
halves are needed. The absent alone would pass on the unpaged answer too.

Two smaller ones alongside. The paging generator was insisting on a 200 where
the listing generator had already been taught not to -- Agora answers 201 -- and
a position is any identifier the served record carries, not one of the six names
the generator knew: Workable's jobs are told apart by a `shortcode`.

### Two requests that were nearly right

SerpAPI's template case writes its whole query into the path, in the older
spelling: `/search?q=cauldron+search+recipe&api_key=...`. The generator read
that into a query dictionary without decoding it, so it asked for a phrase with
plus signs in it and the Recipe answered 404 -- correctly, since that is not a
search anyone ran. A path's query string is percent-encoded and a query
dictionary holds the decoded value; the two are not the same text.

PlanetScale's only green case for its databases listing asks for somebody
else's organisation, and is answered with an empty list, which is the point of
that case. The generator copied the path along with everything else and found
nothing to page. The listing generator already knows how to fill a path
parameter from the fixture; borrowing that and trying the other candidates when
the copied path serves nothing settles both names.

Neither was a provider being difficult. Both were a request that was nearly
right, which is the harder kind to see, because a nearly-right request comes
back with a real answer.

### An empty answer is still an answer, and it is not the same answer

Seven of the listings nothing had ever answered are in Recipes with no fixture
holding a record of that resource. Dub's links, Melio's bills, FullStory's
sessions and four more: one fixture, `empty`, nothing to serve. The generator
could not show a record because there is no record, and it had been treating
that as nothing to do.

There was something to do. The route answers. It answers with a status, an
envelope and an empty collection, and *every case in those files until now was
checking a refusal*, so none of that had ever been seen. A case that asserts the
collection is empty rather than absent is worth writing: it is the distinction a
client gets wrong, and it is true.

But it is not the same evidence as a listing that shows a record, and a single
counter would have called them equal the moment those cases landed. So the
report grew a second line. `HollowListing` counts listings answered only empty
-- shown, but with nothing asserted inside the collection, which leaves the
collection key, the field names and the identifier style beside them still
undescribed.

It found more than the seven just written. **Thirty-two listings across
twenty-six Recipes are answered only empty**, twenty-five of which have been
that way all along, counted by nothing. That is the third time a new counter has
turned up debt older than the counter. The work is not finding new problems; it
is finding out that the old measurement had a shape, and the shape had a shadow.

### Twelve of the thirty-two had a record all along

The generator that shows a listing had been aimed at the wrong question. It
looked for listings with no successful case at all, which is `UnshownListing`'s
question, and a listing answered only empty already has one. So the twenty-five
hollow listings the new counter surfaced were invisible to the tool built to
fill exactly that gap.

Aiming it at `HollowListing`'s question instead -- does any successful case
assert something *inside* the collection -- found twelve of them answerable
straight away. Stripe's customers and payment intents among them: three records
in the fixture, a listing every case in the file answers, and nothing anywhere
asserting a single field on a customer.

Two shapes needed teaching. A bare listing has no envelope, so anything asserted
on it is asserted on a record -- the counter was requiring an index that Alpaca's
account and SES's, both `collapse_single`, never have. And the generator needed
the same collapse handling the paging one had already grown: with one record the
response *is* the record.

Of the twenty-one still hollow, two are meant to be. iNaturalist's
`/v1/observations/999999999999` and OpenAIRE's `publications-no-match` are
routes whose whole purpose is the empty answer, and there is no record to
describe because a match would defeat the point. The rest have no fixture
holding a record of the resource.
