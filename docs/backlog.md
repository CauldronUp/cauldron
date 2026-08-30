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
| ~~AWS SES~~ | Shipped. Accepted-not-delivered, the invisible suppression list |
| ~~AWS DynamoDB~~ | Shipped. Typed attributes, omitted Items, table states |
| ~~AWS Secrets Manager~~ | Shipped. Undeducible ARNs, rotation stages, scheduled deletion |
| ~~Google Cloud Storage~~ | Shipped, written against the discovery document Google publishes. The folders come back in a different array from the files: there are no directories in a bucket, only names with slashes, and of delimiter -- "Objects whose names, aside from the prefix, contain delimiter will have their name, truncated after the delimiter, returned in prefixes. Duplicate prefixes are omitted." So a bucket of ten thousand files answers with two items and a handful of prefixes, and code reading items has a complete, accurate array that is not the answer to its question. And the paging counts both arrays at once -- maxResults is "items plus prefixes" and "fewer total results may be returned than requested" -- so the loop everybody writes stops early on the first page where two objects shared a folder. Also: the id includes the generation, so overwriting a file changes it and name is the stable one; size, generation and metageneration are digits in strings; and being deleted is two states, timeDeleted and softDeleteTime. Fetching one object by name is stated and not served -- an object name goes into the path URL-encoded and Go decodes it before the router sees it |
| ~~Google Pub/Sub~~ | Shipped. Base64 bodies, ack deadlines, delivery attempts |
| Azure Blob Storage | Containers, blobs, SAS and auth behaviour |
| Cloudflare R2 | Sits beside the existing Cloudflare Recipe |
| Vultr | Instances, block storage, the async provisioning lifecycle |

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
| USPS | Addresses, rates, tracking. Blocked on evidence rather than on interest -- see "A description of requests is not a description" below |
| UPS | Rating, labels, tracking |
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
| Descope | Authentication flows and identities |
| ~~FusionAuth~~ | Shipped, written against the OpenAPI document FusionAuth publishes. A header that is optional until somebody else changes something: X-FusionAuth-TenantId is declared required: false on nearly every operation, and its own description says "Only required when there is more than one tenant and the API key is not tenant-scoped". Both of those are administrative acts in a different system, so code written against a single-tenant install works, ships, and breaks the day a second tenant appears -- and a developer's sandbox is single-tenant, which is precisely the state in which this cannot be found. This Recipe seeds two. Also: an account's password can go from fine to breached without the account being touched, insertInstant and lastUpdateInstant identical while breachedPasswordStatus changes; active is a boolean and expiry is a date and neither mentions the other; and timestamps are called insertInstant and lastUpdateInstant, in milliseconds, so mapping by convention finds neither |
| ~~Keycloak~~ | Shipped, the Admin REST API, written against the OpenAPI document Keycloak publishes. Creating a user tells you nothing about the user you created: POST /admin/realms/{realm}/users is documented as "201 Created" and nothing else -- no content, no schema, no headers -- so the document describes no response body for the endpoint that makes the central object of the API, and does not say the identifier went into a Location header either. Also: enabled does not mean able to log in, because a requiredActions entry stops a login dead while enabled stays true; search is prefix-based by default with *foo* for infix and quotes for exact, so the punctuation is the API; and first is the pagination offset while the exact parameter's own description names "first" as a field to match on. The search punctuation is stated and not served |

## Analytics and flags

| Provider | Why |
|---|---|
| ~~LaunchDarkly~~ | Shipped. Per-environment state, variations as indices |
| ~~PostHog~~ | Shipped. A property is nested under properties and a flag is a string, a boolean or false, all in the same field |
| Mixpanel | Events, profiles, exports. **Assessed and written up rather than written** -- see "Mixpanel, and an API that only speaks in batches" below. The findings are recorded there; the blocker is that every ingestion endpoint takes an array and this format models one record per request |
| Amplitude | Events, cohorts, user properties. Batch-shaped for the same reason Mixpanel is, so read that assessment first |

## Models and inference

| Provider | Why |
|---|---|
| ~~OpenAI~~ | Shipped, written against the OpenAPI document OpenAI publishes and generates its own SDKs from. The refusal is not in the answer: content and refusal are two sibling nullable strings on the message, only one is ever filled in, and a declined request is a 200 with finish_reason stop -- so the obvious read gets null and logs an outage while the model answered in the field beside it. Also: completion_tokens counts tokens that never arrive ("like reasoning tokens, these tokens are still counted in the total completion tokens for purposes of billing"), max_completion_tokens is spent on reasoning before a word is emitted, max_tokens is deprecated and "not compatible with o-series models", store defaults to false so the completion you just made cannot be fetched, and content_filter is a stopping reason rather than an error. Streaming and the Responses API are stated as not modelled |
| Anthropic | Messages, content blocks, tool use, streaming |
| ~~Google Gemini~~ | Shipped, and the section near the end of this file that recorded it as unservable is now the story of how it was unblocked. A blocked prompt is a 200 with the candidates taken away -- promptFeedback.blockReason is "If set, the prompt was blocked and no candidates are returned" -- so candidates[0] throws where OpenAI's refusal hands back a null. Serving the contrast needed one path to answer two shapes chosen by the request body, which selects could not do, so selects_body was added beside it as a separate field. Also: an empty finishReason means the model "has not stopped generating tokens", from one schema serving the streaming and non-streaming calls; totalTokenCount is documented as "prompt + thoughts + response candidates"; and promptTokenCount "is still the total effective prompt size" when content is cached |
| ~~Pinecone~~ | Shipped. The response to "make me an index" is an address somewhere else. Pinecone is two APIs with two base URLs: the control plane at `api.pinecone.io` creates and lists indexes, and the data plane -- query, upsert, fetch, which is every operation an application actually performs -- lives at the index's own host, which you cannot know until the control plane tells you. The data plane's own OpenAPI document says so by being unable to name its server: `url: https://{index_host}` with the variable's `default: unknown`, so a client generated from it and left unconfigured posts to `https://unknown/query`. `host` is required on every IndexModel because there is no other way to find the thing you just made. Also pinned: `ready` and `state` are two fields and the vendor's own example has them disagreeing -- `{"ready": true, "state": "ScalingUpPodSize"}`, ready and not in the state called Ready, out of nine states; a create answers 201 with `host` filled in and `ready` false, so the address exists before the index does; a delete is a **202** ("The request to delete the index has been accepted") and `Terminating` is one of the nine states, so the index you deleted answers describe for a while afterwards; a sparse index has no `dimension` at all rather than a null one; and the error vocabulary is gRPC's -- the seventeen canonical status codes with FORBIDDEN, UNPROCESSABLE_ENTITY and PAYMENT_REQUIRED bolted on, which is why **`OK` is a documented value of a field called "The error code"** -- with the HTTP status repeated inside the body. `X-Pinecone-Api-Version` is `required: true` on every operation in both documents. Stated and not served: the data plane itself, because it is a second base URL and this emulator is one address, and serving `/query` here would teach the exact mistake the Recipe describes |
| ~~Replicate~~ | Shipped. A created prediction has no output property at all rather than a null one, succeeded is not the same as produced something, a cold start is a minute with no signal but a boot_time, and the output is a link to a file deleted after an hour |
| Hugging Face | Inference endpoints and model responses |
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
| Redis Cloud | Keys, TTLs, streams, pub/sub |
| Upstash | Redis REST, queues, rate limits |
| ~~MongoDB Atlas~~ | Shipped, the Administration API, written against the OpenAPI document MongoDB publishes. The version of the API you are talking to is a date inside a content type: GET /groups/{groupId}/clusters documents one 200 with three content types under it -- vnd.atlas.2023-01-01+json, 2023-02-01+json and 2024-08-05+json -- resolving to three separately named schemas, all current. And the difference is not cosmetic: the legacy view carries mongoURI, srvAddress, diskSizeGB and providerSettings and the newest carries none of them, so the field an application connects with is present under one date and absent under the next, same URL, same credentials, same instant. Also: the schema names carry the date (ClusterDescription20240805), failure is better documented than success (more operations describe a 401 than a 200), and the credential is HTTP Digest -- challenge-response, on a cloud API written in the 2020s. Digest is stated and not served, and so is the Content-Type echo |
| ~~Neon~~ | Shipped. Branches, databases, endpoints |
| ~~PlanetScale~~ | Shipped, written against the Swagger 2.0 document PlanetScale serves from its own API host. The field called state is not the state of the deployment: a deploy request carries state -- "Whether the deploy request is open or closed" -- beside deployment_state, "The deployment state of the deploy request", so a request abandoned after review and one whose migration ran an hour ago both read closed. Also: the id is "The ID of the deploy request" and every path takes the number instead, which is per database, so the globally unique identifier addresses nothing; a request outlives the branch it came from, still naming it while branch_deleted says it is gone; and next_page is "null when this is the last page" where Confluence omits the field entirely. The ten deploy-lifecycle endpoints are stated as not modelled, because each advances a state machine this format cannot express |
| CockroachDB Cloud | Clusters, SQL users, operations |
| Turso | Databases, replicas, tokens |
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
| Sinch | SMS, verification, messaging |
| ~~Telnyx~~ | Shipped. Per-recipient status, cost after the fact, silent MMS |
| Infobip | SMS, WhatsApp, verification |
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
| Smartsheet | Assess — sheets, rows, columns. Column ids are numeric and per-sheet, so nothing is portable between sheets |
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
| Checkout.com | Assess — payments, 3DS, the difference between authorised and captured |
| Authorize.Net | Assess — the XML-shaped API and its own result codes, which are not HTTP statuses |
| Klarna | Assess — order lifecycle across authorise, capture and refund, with a session that expires |
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
| Flagsmith | Assess — flags per environment, identity overrides |
| Split | Assess — treatments and targeting rules |
| ~~ConfigCat~~ | Shipped, and the row's third item is now a separate piece of work -- see below. **The Public Management API** is what ships. Two paths that differ by one character answer with different documents: `/v1/settings/{settingKeyOrId}/value` gives `SettingValueModel` -- `value`, `rolloutRules`, `rolloutPercentageItems` -- and `/v2/...` gives `SettingFormulaModel` -- `defaultValue`, `targetingRules`, `percentageEvaluationAttribute`. The field holding the value is called `value` in one and `defaultValue` in the other, and it is wrapped in the second; the two rule arrays become one. A client that upgrades the URL and not the reader gets undefined from a 200. Also pinned: **two credentials in one request doing different jobs** -- Basic in `Authorization` says who you are, `X-CONFIGCAT-SDKKEY` says which config and environment you mean, and only 6 of the API's 76 operations take the second; **one path segment holding two identifier spaces**, since `settingKeyOrId` is "The key or id of the Setting" and a setting has both a human-chosen `key` and an int32 `settingId`, so a flag keyed "42" collides with setting number 42; "may I save this" having four answers on v2 (`readOnly`, `approveRequired`, `canBypassApproval`, `reasonRequired`) and one on v1, so the same flag looks editable on one path and gated on the other; and a `value` whose JSON type follows `settingType` (boolean, string, int, double), with `isJson` beside it for the strings that are documents. **The config-delivery surface is the remaining work and it is where the users are**: ConfigCat SDKs fetch a static configuration file from the CDN and evaluate locally, never touching `api.configcat.com`, and `configcat/configcat-client` has 1.29M downloads against about 5k for the clients that speak the Management API -- a ratio of roughly 260 to 1. Different host, different shape, different job |
| ~~Unleash~~ | Shipped, and the row was right that the APIs are different shapes -- it is worse than that. Three APIs expose the same flag, each schema has an `enabled` boolean, and in none of them does it mean "this flag is on". `frontendApiFeatureSchema.enabled` is documented as **"Always set to `true`."** -- the Frontend API returns only the flags that already evaluated on, so the field is a constant and the real answer is whether the object is in the array at all, which makes `toggles.find(t => t.name === x).enabled` a question whose answer is always yes when it does not throw. `clientFeatureSchema.enabled` is one input: "This is ANDed with the evaluation results of the strategies list", and the strategies are themselves "evaluated and ORed together", so the rule is an AND over an OR and it exists nowhere except as two clauses in two field descriptions. And the Admin API's `enabled` is a summary of an `environments` array whose members disagree with it -- on in development, off in production, one boolean at the top. Also pinned: the two SDK APIs do not agree on what to call the collection (`toggles` on one, `features` on the other, same flags, same server, because Unleash renamed the concept and one endpoint moved); the response format carries its own version number in the body, where version 2 "includes segments as a separate array" so an evaluation that reads strategies without joining segments applies a different rule while looking complete; and `description` and `variants` arrive as null rather than absent. Two authentication schemes share one header -- `apiKey` raw and `bearerToken` behind "Bearer " -- and only the raw form is served, which is pinned as a refusal so the gap is visible. Detection found the shape next door to YouTube's simile: matched on an imperative, where the biggest npm result for "unleash" is `precinct` at twenty million downloads, a dependency parser described as "Unleash the detectives" |

## Observability and product analytics

| Provider | Why |
|---|---|
| ~~Bugsnag~~ | Shipped. An error is not its occurrences, and the counts are on the error |
| Honeybadger | Assess — faults and notices |
| Grafana Cloud | Assess — the stack-management API, which is a second surface on top of the one that now ships. **The Grafana HTTP API is shipped** as `grafana`, written against `api-merged.json` in `grafana/grafana`: two unique identifiers for one dashboard, and the one called `id` is the one you cannot use. A save's required fields are `["status", "title", "version", "id", "uid", "url"]` and the two identifiers carry the same sentence -- "The unique identifier (id)" and "The unique identifier (uid)" -- for an int64 and a string, while the only path that fetches a dashboard is `/api/dashboards/uid/{uid}`. So the integer is required in every save response and can address nothing; it is the instance's row number, and a deploy that stored it stored the identifier that will not survive the move. The document says the same about folders outright: `folderId` is "Deprecated: use FolderUID instead", beside `folderUid`. Also pinned: a field called `title` whose description is **"Slug The slug of the dashboard."** with the example "my-dashboard" -- a struct comment one field out of place, on the field a UI shows back to the user; the `version` that guards the next save living in `meta` rather than on the dashboard, which is free-form JSON, with a **412** declared on the save; five permission booleans (`canAdmin`, `canDelete`, `canEdit`, `canSave`, `canStar`) and a sixth field, `provisioned`, that none of them accounts for; and a search that answers with a bare array, no envelope and no count, whose hits can be dashboards already in the bin. Detection is the thinnest in the collection and the reason is the finding: eighteen npm results for "grafana" and not one calls this API -- the whole `@grafana/` scope is plugin tooling, and the biggest Packagist numbers are Loki, a different product from the same vendor |
| Honeycomb | Assess — datasets, triggers, query results |
| Better Stack | Assess — logs, uptime monitors, incidents |
| Heap | Assess — events and user properties |
| LogRocket | Assess — sessions and issues |
| New Relic | Assess. It was in "assessed and deliberately not done" for being GraphQL-only, which stopped being a reason the day ShipHero brought GraphQL support and Linear shipped on it. NerdGraph is one endpoint whose body decides the shape, which is what `selects` is for |

## Hosting, deployment and package registries

| Provider | Why |
|---|---|
| ~~Netlify~~ | Shipped. A deploy id exists long before the site is live |
| Render | Assess — services, deploys, the build-then-live gap |
| ~~Fly.io~~ | Shipped, and the comparison this row asked for is in the file: the Machines API is served and the older platform API is recorded as a stated gap, because it is GraphQL at a different host and is a second Recipe rather than more of this one. The headline is that `started` does not mean your application is up -- a machine has a `state`, everybody reads it, and three other fields on the same object independently contradict it. `host_status` can be `unreachable` (started, on a host Fly cannot currently talk to -- not stopped, not failed, not anything `state` has a word for), `cordoned` can be true (started and taking no new traffic, on purpose, and nothing in `state` says so), and `checks` can be `critical` (the process is running and failing the health check that exists to describe exactly that). Four answers to "is this machine up" on one object, disagreeing by design |
| ~~Heroku~~ | Shipped. Assess — the API still uses `Accept: application/vnd.heroku+json; version=3`, so a missing header is a different response rather than an error |
| Linode | Assess — instances and the async provisioning lifecycle, beside Vultr and DigitalOcean |
| ~~Hetzner Cloud~~ | Shipped, and this row asked for the sentence the Recipe opens with. Powering off a server does not answer with the server: it answers with an action -- a job, with a status of `running`, a progress percentage and a `finished` timestamp that is null. The machine is still on. Whether it will go off is a second request, to a different endpoint, about a different object, and there is no version of the call that waits. Every mutation in this API is that shape -- attach a volume, change a type, rebuild, enable rescue, request a console, each answers 201 with a job -- so nothing that changes anything answers with the thing it changed, and a test that powers a server off and asserts it is off passes against a mock, passes against a fixture, and is a race against the real thing. Also pinned: an action that fails long after its successful 201, an `error` that is an object rather than a sentence, a `progress` that reaches 100 on a job still running, nine server statuses of which "not running" is eight, `locked` as a field separate from `status` refusing the next mutation with a 423, and a root password returned exactly once |
| Scaleway | Assess — instances and object storage |
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
| Codecov | Assess — coverage reports and the commit they attach to |
| ~~SonarCloud~~ | Shipped. The quality gate has three outcomes and one of them means nobody set a gate: "The different statuses returned are: OK, WARN, ERROR, NONE. The NONE status is returned when there is no quality gate associated with the analysis." So `if (status !== "OK") fail()` fails the build for an ungated project and `if (status === "ERROR") fail()` passes everything on one. And **every number inside the gate is a string** -- the published example sends `"errorThreshold": "85"` beside `"actualValue": "82.50562381034781"`, with the direction of the comparison in a third field as `"comparator": "LT"` -- so evaluating a condition means parsing two strings and reading an operator, and comparing them as text is right on "14" against "0" and wrong on "9" against "10". Also pinned: **there is no verb but GET and POST** (87 of the 156 actions are POST, 69 are GET, and the description cannot express another -- each action carries a boolean called `post`, so deleting a comment is `POST api/issues/delete_comment` and the path is the verb); the analysis and the gate being different objects, since a Compute Engine task that FAILED produced no gate result at all; a failed task handing back a **Java stack trace inside a JSON string** beside a boolean saying whether to read it; and an issue carrying two names for its status (`issueStatus`/`status`) and two for its effort (`effort`/`debt`), agreeing today. Written against the description SonarCloud serves of itself at `api/webservices/list` with response examples at `api/webservices/response_example` -- and `project_status` carries `"deprecatedSince": "16 September, 2025"` there while its own response says nothing, so the only way to learn the endpoint is deprecated is to call a different one that describes it. Detection found an exclusion kind this collection had not recorded: packages that write a file for another program to upload, led by `vitest-sonar-reporter` at 841k downloads and no network request at all |
| Railway | Assess. Here for the same lapsed reason as New Relic: GraphQL-only was an exclusion before anything here spoke GraphQL, and three Recipes do now |

## Secrets and configuration

| Provider | Why |
|---|---|
| ~~HashiCorp Vault~~ | Shipped, and this row named the headline before it was written: the same read gives you the secret, or a box containing the secret. On KV v1 `response.data` is your secret; on KV v2 `response.data.data` is, and `response.data` is a box holding it beside a metadata block. Nothing in the path says which, because the version is a property of the mount. `secret = response.data` does not throw against v2 -- it yields an object, truthily, and surfaces later as a password that is an object. Written against openapi.json in `hashicorp/vault-client-go`, the document HashiCorp generates its own Go client from, which turned out to omit three things every real request and response has: the `/v1/` URL prefix (stated in `info.description` and on none of its 715 paths), the response envelope (the document's `KvV2ReadResponse` is the *inner* half; the wrapper carrying `request_id`, `lease_id`, `renewable` and `lease_duration` is hand-written outside the generated code as `Response[T any]`), and authentication (`securitySchemes` is `{}`, empty, for an API whose purpose is holding secrets). The generator emits the middle and a person supplies both ends. Also pinned: a write answers with metadata and not with the secret -- the request schema says so in passing, "will be stored and returned on read" -- while the v1 write answers 204 with no body at all; not-deleted is `deletion_time: ""` beside `destroyed: false`, an empty string rather than a null; and a version is an integer in `current_version` and a string when the same number is a key in `versions`. Stated and not served: `/sys/health`, whose five declared responses are 200 active, **429 "unsealed and standby"**, **472** (not an HTTP status code) for a DR replication secondary, **501 "not initialized"** and 503 sealed -- codes chosen so a load balancer's "2xx is healthy" rule finds only the active node, which means every generic monitor reports a working standby as rate-limited. Those are states of the server rather than properties of a request. Also unserved: listings, which answer with an array of names rather than records. Detection found the most contested word in the collection -- four vendors ship a Vault and the biggest is Azure's |
| Doppler | Assess — configs, secrets, the inherited value that looks local |
| 1Password | Assess — Connect and Service Accounts are separate APIs |
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
| Inngest | Assess — events, functions, step state |
| Trigger.dev | Assess — runs, tasks, the resumed run |

## Scheduling and real-time

| Provider | Why |
|---|---|
| ~~Cal.com~~ | Shipped. A slot listing is a view rather than a reservation and has no identifier to hold on to, so booking one that was free is a 400 that reads like your bug; cancelling does not delete; and pending holds the slot without anybody agreeing |
| Acuity Scheduling | Assess — appointments, types, intake forms |
| Daily.co | Assess — rooms, meeting tokens, recordings that finish after the call |
| LiveKit | Assess — rooms, participants, egress |
| Agora | Assess — channels, tokens, recording |
| Stream | Assess — chat channels, members, the message that is soft-deleted and still returned |
| SendBird | Assess — channels and messages |

## Lifecycle and marketing messaging

| Provider | Why |
|---|---|
| Customer.io | Assess — the Track and App APIs are separate hosts with separate credentials, which is the sort of thing that only fails in one environment |
| ~~Braze~~ | Shipped. "message" is both the success channel and the error channel, a listing entry is a smaller object than the detail, and the export answers with a prefix instead of data |
| Iterable | Assess — users, campaigns, catalogues |
| ActiveCampaign | Assess — contacts, deals, automations |
| Brevo | Assess — contacts, transactional email, the shared quota between them |
| Kit (ConvertKit) | Assess — subscribers, sequences, tags |
| Attentive | Assess — SMS subscribers and consent state, where consent is legally load-bearing |
| Beehiiv | Assess — publications, posts, subscribers |
| Loops | Assess — contacts and transactional sends |

## Social and content platforms

| Provider | Why |
|---|---|
| Reddit | Assess — listings page by fullname rather than offset, and a deleted comment is still in the tree with its body replaced by a marker |
| Twitch | Assess — Helix pagination cursors, token scopes, the EventSub subscription lifecycle |
| Bluesky | Assess — the AT Protocol is its own model with DIDs and records rather than REST resources, so it may belong in the not-done table |
| Mastodon | Assess — pagination is Link-header based and instance behaviour varies, which is itself the interesting part |
| Telegram Bot API | Assess — every method is both GET and POST, errors come back with HTTP 200 in some client libraries, and updates arrive by long polling or webhook but never both |
| Buffer | Assess — profiles, updates, scheduling |
| Spotify | Assess — the token expires in an hour and the refresh path is the whole integration |
| Strava | Assess — activities, webhooks, the rate limit counted in two windows at once |

## Identity verification and risk

| Provider | Why |
|---|---|
| ~~Persona~~ | Shipped. The decision arrives by webhook minutes after the inquiry is created |
| ~~Onfido~~ | Shipped. A check is complete and its report can still be consider rather than clear |
| Veriff | Assess — sessions and decisions |
| Middesk | Assess — business verification and its partial matches |
| Alloy | Assess — evaluations and their outcomes |
| Sift | Assess — scores, decisions, the workflow that runs server side |

## Travel

| Provider | Why |
|---|---|
| Duffel | Assess — offers expire, and an offer that was valid when you displayed it is gone when the customer presses book. That gap is the entire domain |
| Amadeus | Assess — self-service against enterprise APIs, and the test environment carrying different inventory from production |
| Hotelbeds | Assess — availability, rates, the rate that changes between check and book |

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
| ScrapingBee | Assess — a failed fetch is a successful API call |
| Browserless | Assess — sessions and timeouts |

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

**15 routes across 7 Recipes** still page by a parameter nobody named.
Twenty-two providers have been settled. Five were read from that provider's own
description, none of them guessable:

| Provider | What the description said |
|---|---|
| SendGrid | Two listings paging two different ways. Templates declares `generations`, `page_size` and `page_token` and no `limit` at all; bounces declares `limit` and `offset` |
| Chargebee | `limit` and `offset`, and `offset` is typed **string** with the description saying to set it to the `next_offset` the last response returned. A cursor wearing an offset's name |
| OneSignal | `limit` and `offset` on both listings |
| Webflow | Eleven query parameters on the items listing including `limit` and `offset`, and none whatsoever on the two site-level listings, which are now `limit_param: "-"` |
| Pipedrive | `start` and `limit` in v1. Not `cursor` and `limit`, which is v2 -- and v2's `/deals`, `/persons` and `/organizations` live under `/api/v2` while v1's description no longer declares them at all. A Recipe modelling v1 paths pages the v1 way |

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

**168 more listings across 102 Recipes declare no paging at all**, and the
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
| Checkout.com | Authorised is not captured, and the two have separate identifiers |
| Klarna | Blocked on documentation, like Sift. The premise is good -- an authorization token expires while the order it approved has not been created, so a customer can be approved and unpayable -- and docs.klarna.com renders its API reference through a portal that serves navigation rather than schemas, so the response bodies cannot be read or cited. Revisit with an account or a published spec |
| ~~Airwallex~~ | Shipped. The settlement currency is not the charge currency and the rate was fixed at a moment nobody picked, there is a separate balance per currency so a payout can fail on a funded account, and a partial capture leaves two amounts different forever |
| ~~Column~~ | Shipped. A notification of change is not a failure, and R01 and R07 are two characters apart with opposite obligations |
| Moov | Transfers across rails, where the rail decides the settlement time and nothing else does |
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
| Rippling | Assess — employees, groups, app provisioning. Same standing rule as Gusto and Deel |
| BambooHR | Assess — employees, time off, custom fields whose ids are per-account |
| Ashby | Assess — candidate stages, where the stage is configurable per pipeline so branching on its name breaks |
| Workable | Assess — jobs, candidates, stages |

### Messaging

| Provider | Why |
|---|---|
| Plivo | Assess — SMS and voice, beside Twilio and Telnyx |
| ~~Bandwidth~~ | Shipped, messaging only. The same message comes back with entirely different property names depending on which endpoint returned it, and a send is 202 Accepted rather than sent. Number provisioning is genuinely asynchronous and would be worth having, but that API answers in XML |
| ~~Twilio Verify~~ | Shipped. A wrong code is a 200 with the verdict in a word inside the body, so `if (!res.ok)` lets the wrong person through the second factor. Three different 429s, two of which are terminal and one of which is not. A verification is deleted on approval, so checking twice and never having started are the same 404 |
| Customer.io | The Track and App APIs are separate hosts with separate credentials, which is the sort of thing that works in one environment and not the other |
| ~~Braze~~ | Shipped. The export answers 201 with a prefix and no users; the file lands in cloud storage minutes later, so a test reading users off that response reads nothing forever |
| Brevo | The premise is wrong and the row is kept to stop it being re-derived. Brevo's limits are **per-endpoint**, not one pool: `POST /v3/smtp/email` gets 1,000 RPS while `GET /v3/contacts/{id}` gets 10. What is true is duller and still real -- two endpoints on one API differ by a hundredfold, so a throttle calibrated on the send is wrong on the read -- and Cauldron models rate limits as armable errors rather than as counters, so a Recipe could only restate it |
| Kit | Assess — subscribers, sequences, tags |
| Beehiiv | Assess — publications, posts, subscribers |

### Identity and risk

| Provider | Why |
|---|---|
| AWS Cognito | The SDK is the API, and the token lifecycle is the whole integration |
| ~~Ory Kratos~~ | Shipped, and the flow object was the unusual part. The server sends you the form: a flow carries the URL to post to, the method, and the list of fields to render, so a login page is a renderer for somebody else's JSON and a hardcoded {email, password} POST works against one deployment and not the next. An expired flow is 410 and 410 is not a retry -- Ory's words are that a new flow has to be initiated. The CSRF token is a node in the same array as the visible fields, a node's `type` and its attribute's `type` are different words at different depths, and `identity.state` is documented as having no effect while still carrying an enum |
| Kinde | Assess — users, organizations, feature flags in one product |
| ~~Persona~~ | Shipped. completed is not approved, needs_review is neither, and nothing is at the top level because it is JSON:API |
| ~~Onfido~~ | Shipped. complete is not clear, consider is neither a pass nor a failure, and the reason lives on the report rather than the check |
| Sift | Blocked on documentation, not on interest. The premise holds and the shape is good, and the public docs do not publish a score response body: the pages describe the 0-100 scale in prose while the API is widely reported to send decimals, and that number is the one every integration branches on. Building it would mean guessing the scale, which is the one thing a Recipe must not do. Worth revisiting with an account, where a single real response settles it |

### Storage and media

| Provider | Why |
|---|---|
| ~~Backblaze B2~~ | Shipped, and the row was right about where to look. Deleting a file brings back an older one -- Backblaze's words are that the most recent older version becomes the current version -- so delete-then-read is a 200 and stale bytes rather than a 404. There is no overwrite either, so the S3 habit of writing a key to replace it accumulates billed copies. A hide is a version with its own id, the action enum is documented open, the base URL is data and comes back from authorize, and part sizes are quoted strings beside a contentLength that is a number |
| ~~Uploadcare~~ | Shipped. A file exists before it is stored and unstored files are deleted after twenty-four hours, so the same code works on one project and loses files on another; a removed file still answers 200 with everything intact |
| ImageKit | Assess — transformations and delivery, beside Cloudinary and Imgix |
| Livepeer | Assess — streams and recordings, asynchronous throughout |

### Data and search

| Provider | Why |
|---|---|
| ~~Qdrant~~ | Shipped, and **not** for the reason this row gave. The filter behaviour needs the search itself and Cauldron does not do vector arithmetic, so the header says outright that filters are not applied. What shipped instead is the envelope: `status` is the string `ok` on success and an object `{error}` on failure, so a typed client fails to parse one of the two paths and an untyped one reads `status.error` as undefined, which is falsy, which reads like no error. Plus: `result` nests twice, a collection listing hands back names only, a write answers `acknowledged`, `points_count` and `indexed_vectors_count` disagree on purpose, and `version` is on a query result but not on a point fetched by id |
| ~~Weaviate~~ | Shipped, and the thing worth having was the batch. A bulk import that half failed is a **200**: Weaviate's own words for the endpoint are that the request was processed successfully and individual object statuses are in the body, so each element carries `result.status` of SUCCESS or FAILED and the failures sit in the array beside the successes, identical in every other way. Plus: an error is an *array* of messages so `error.message` is undefined everywhere, the same object has two URLs of which one is deprecated, and the API says so in a `deprecations` field on the response body rather than in a header |
| ~~Typesense~~ | Shipped, and the thing worth having was the relevance score. `text_match` is an int64 up near 578730123365711993, which `JSON.parse` rounds to 578730123365712000 -- and the hit ranked below it, ...994, rounds to the same number, so two differently-ranked results compare equal. This is the worse half of a problem the format already knew: Discord sends snowflakes as strings so it cannot happen, and Typesense sends a number, so no client can avoid it. Plus `search_cutoff` is a documented boolean for "your results are incomplete" at 200, and there are three counts of which one counts the array |
| Convex | Assess — the query and mutation model is not REST-shaped, so this may belong in the not-done table |

### Hosting and deployment

| Provider | Why |
|---|---|
| ~~Netlify~~ | Shipped. A deploy has an id before it has a site, ready is not published, and a missing URL means two different things |
| Render | Blocked on documentation, like Sift and Klarna. The shape looked good -- list endpoints appear to wrap each element beside a cursor, and `suspended` appears to be a string enum rather than a boolean -- and none of that could be confirmed: `api-docs.render.com` 404s on its reference pages and the OpenAPI spec its own docs point to is not at the URL given. Revisit with an account or a spec that resolves |
| ~~Fly.io~~ | Shipped, and **not** for the reason this row gave: the older platform API is GraphQL at a different host, and this format speaks REST, so the two-shapes half is stated in the header and not modelled. What shipped is better anyway. `state: started` does not mean the application is up, and three fields on the same object independently say so -- `host_status` can be `unreachable`, `cordoned` can be true, and `checks[0].status` can be `critical`. Four answers to one question, disagreeing by design. Plus: `instance_id` is unique per *version*, so anything keyed on it loses its history at every deploy, and `nonce` is returned once, at creation, and only if a lease duration was asked for |
| ~~Heroku~~ | Shipped, and the header was the smaller half. A successful list is **206**: Heroku pages with the `Range` header and answers `206 Partial Content` while there is more, with the resume point in `Next-Range` rather than in the body -- so comparing against 200 rejects every page but the last, and testing `ok` accepts them and never looks for the rest. The `Accept` version header is a 406 when missing, errors are keyed by `id` rather than `code`, `url` is on an error only sometimes, and a formation with `quantity: 0` is a process type that exists and is not running |
| ~~Hetzner Cloud~~ | Shipped, and the row was exactly right. Powering off a server answers 201 with a job: status running, progress 0, finished null, and the machine still on. Every mutation in the API is that shape, so nothing that changes anything answers with the thing it changed. An action can fail long after its 201, its reason is an object rather than a sentence, progress reaches 100 while still running, a server has nine statuses of which eight are not running, and `locked` is a separate question from `status` whose refusal is a 423 |

### Observability and flags

| Provider | Why |
|---|---|
| ~~Bugsnag~~ | Shipped. An error carries counts and no stack trace, fixed is not terminal, and severity is a different question from unhandled |
| Grafana Cloud | Assess — the stack-management API and the Prometheus-shaped query API that is not REST. The Grafana HTTP API itself ships; see the row in the observability table above for what it pins |
| Better Stack | Assess — logs, uptime monitors, incidents |
| ~~PostHog~~ | Shipped. A flag definition is not what a user gets, nought per cent is not inactive, and capture says the same thing whatever you send |
| Statsig | Assess — gates, experiments, exposure logging |
| Flagsmith | Assessed and left. "A flag definition is not what a user gets" is point one of the PostHog Recipe, already shipped, so a second Recipe would add a name rather than a shape. What differs is worth a row of its own if anyone wants it: Flagsmith answers per environment *key* rather than per SDK evaluation, so the same flag has a different answer per environment and per identity from the server rather than from the client |

### Webhook infrastructure

| Provider | Why |
|---|---|
| ~~Svix~~ | Shipped. Accepted is not delivered, the outcome lives on the attempts, a failing endpoint gets disabled and nothing announces it |
| Inngest | Assess — events, functions, step state, where a step is durable and a function is not |

### Web data

| Provider | Why |
|---|---|
| ~~Firecrawl~~ | Shipped. A running crawl already has results, total is an estimate completed never reaches, and a 200 can carry a failed fetch |
| ~~Apify~~ | Shipped. A run that SUCCEEDED tells you nothing about whether it produced anything: the status describes the process, and the results are in a dataset under a different id at a different endpoint. That endpoint is also the one thing this API does not wrap in `data`, so the code reading a run and the code reading its results unwrap differently against the same provider. Plus: three of the eight statuses are the `-ING` half of a pair, a TIMED-OUT run has partial data worth keeping, `isStatusMessageTerminal` says whether the prose will change, and the duration is there twice in two units |
| ScrapingBee | A failed fetch is a successful API call, so the status code says nothing about whether you got the page |

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
| Paystack | Amounts are in kobo, the smallest unit of a currency whose smallest unit most code has never met |
| Flutterwave | The verify endpoint is the source of truth and the redirect parameters are not, which is the whole of the fraud surface |
| Mercado Pago | A payment can sit in in_process for days, and the status detail rather than the status says why |
| dLocal | The settlement currency is not the charge currency and the rate is fixed at a moment you did not choose |
| Rapyd | Payouts and payments have separate id namespaces that look identical |
| Worldpay | Two APIs, one modern and one not, and the older one is what most accounts are still on |
| ~~FastSpring~~ | Shipped. One purchase is an order, a subscription and an account with three different ids and which one a webhook hands you depends on the event; FastSpring is the merchant of record so three amounts sit on one order and none is the price you set; a cancelled subscription is still active until its period ends |

### Notification and messaging infrastructure

| Provider | Why |
|---|---|
| ~~Xendit~~ | Shipped. A virtual account number exists in the 201 and does not work at the bank for minutes, so the customer is told it does not exist by their own bank; a closed account accepts one exact amount and a short payment bounces days later; amounts are integer rupiah with no subunit at all |
| ~~Midtrans~~ | Shipped. transaction_status and fraud_status are two fields and a payment is only safe when both agree; capture plus challenge means the card was charged and the funds are held; a bank transfer never passes through capture at all |
| PayU | The same merchant has different endpoints per country and the response fields differ between them |
| ~~Gorgias~~ | Shipped. A ticket and its messages are two paginated endpoints read at different moments, so the count on one disagrees with the array on the other; from_agent is true for automated replies; and a reopened ticket keeps its closing time with nothing marking the reopening |
| ~~Kustomer~~ | Shipped. A conversation is a customer timeline carrying every channel they ever used, assignment is on the conversation rather than any message, status and queue are unrelated, and everything is JSON:API so nothing a client wants is at the top level |
| Jotform | Form fields are keyed by numeric ids that change when the form is edited, so yesterday's submission maps to today's form incorrectly |

### AI and inference

| Provider | Why |
|---|---|
| ~~Cohere~~ | Shipped. billed_units and tokens do not agree and pricing from the wrong one is wrong on every request; a finish reason is not an error and MAX_TOKENS is a 200 with a truncated body; an embedding carries nothing but its position |
| Mistral | Streaming and non-streaming answer with different shapes for the same request |
| Together AI | Model names change and a retired model answers 404 rather than falling back |
| Groq | Rate limits are on tokens per minute as well as requests, and the headers report both |
| ~~Meilisearch~~ | Shipped. A write answers 202 with a number and the word enqueued, the document is in neither the document listing nor the index until the task runs, and the task can fail after the 202 that accepted it |
| Hugging Face Inference | A cold model answers 503 with an estimated_time, and the correct behaviour is to wait rather than retry |
| ~~Langfuse~~ | Shipped. Traces are ingested asynchronously and are not readable immediately after being written |

### Data movement and warehousing

| Provider | Why |
|---|---|
| ~~Fivetran~~ | Shipped. Triggering a sync that is already running answers success and starts nothing, succeeded_at and failed_at both persist so which is later is the only health signal, and paused and sync_frequency are different things that each look fine alone |
| Airbyte | A job has attempts, and an attempt failing is not the job failing |
| ~~Hightouch~~ | Shipped. A run with every row rejected finishes as success, warning is neither an error nor a success, the reasons live on an endpoint the run does not link to, and disabled and paused are different things that both stop a sync |
| dbt Cloud | A run has steps, and the run status is not the step status |
| ~~Snowflake SQL API~~ | Shipped. The same endpoint answers 200 with results or 202 with a handle depending on how fast the query was, every value is a string, a row is positional with the names in the metadata, and NULL is a real null inside the array of strings |
| ClickHouse Cloud | Results are paged by the client rather than the server, and the format is chosen by a parameter |

### Incident response and observability

| Provider | Why |
|---|---|
| ~~incident.io~~ | Shipped. Status and severity are configured per workspace so neither is an enum you can hard-code, the category is the only fixed thing, the rank is the orderable thing rather than the name, and the two move independently |
| ~~Opsgenie~~ | Shipped. An alert and an incident have separate ids, lifecycles and close endpoints, so closing one leaves the other open; a create answers 202 with a request id that is not the alert; and a flapping monitor is one alert with a count rather than many alerts |
| Checkly | A check run has assertions, and a passing run containing a failed assertion is possible |
| Rev.ai | A job is asynchronous and the transcript is a separate fetch with its own content types, so the job being done is not the transcript being readable |
| Perplexity | Citations are a separate array whose indices point into the text, so dropping either half makes the other meaningless |

### Auth, one more time

| Provider | Why |
|---|---|
| PropelAuth | Org membership and role are separate reads, and a user in no org is valid |
| Helicone | It is a proxy, so a single call can fail as the upstream provider or as Helicone, and the two error shapes are unrelated |
| Fireworks AI | Model names are paths rather than identifiers, and a cold deployment answers slowly rather than failing |
| SurveyMonkey | A response in progress is returned and then changes, so a page read twice is not the same page |

### Storage and media

| Provider | Why |
|---|---|
| Bunny.net | Purging is eventually consistent and the API confirms the request rather than the purge |
| Transloadit | An assembly is asynchronous and partial results are readable, which is a crawl's trap in a different shape |
| Elastic Cloud | A deployment is created before it is reachable, and the endpoint appears in a later read than the one that created it |
| Rootly | An incident timeline is append-only and events arrive out of order, so the last write is not the latest event |

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
| Toast | No client for the POS API on either registry, and the name means a notification popup: tall-toasts alone has over four hundred thousand installs. A mapping on the word would offer a restaurant point-of-sale emulator to every Laravel project that shows a toast when a form saves, so a test guards the decision |
| Royal Mail | Royal Mail runs three separate APIs on three hosts, and every client package on Packagist targets one of the other two: elliotjreed/royal-mail-tracking is the tracking service, and turtledesign/royalmail-php, zvps/royal-mail-shipping-rest-api-client and mobi-market/royalmail-shipping-v3 are all the Shipping API. Click & Drop has no client of its own, and offering its emulator to a project using a different Royal Mail API would be wrong about the host, the credential and the vocabulary at once, so a test guards the decision |
| Homebrew | Nothing calls it either. npm's results for "homebrew formula" are formula *parsers* -- `nth-check` at 260 million and `@sideway/formula` at 55 million, matched on the word -- beside two plugins that edit a formula file in a tap repository rather than calling the API. Packagist's is the hobby: `georgeh/php-beerxml` parses beer recipes, and `munkireport/homebrew` inventories Homebrew on managed Macs by running `brew` locally |
| Repology | Nothing on either registry calls it. npm has two MCP servers and Packagist's top result for the word is `pragmarx/countries` at 3.5 million downloads, a countries and currencies library -- so the empty result and the wrong result look identical, which is the reason to write down that it is empty |
| Tradier | A third-party SDK exists and is one person's; worth a look rather than a mapping on sight |
| Twilio Verify | Part of the main Twilio SDK, which maps to the twilio Recipe rather than this one |

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

Forty-six Recipes send at least one identifier as a number now, and each
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
| Render | Marketing | `api-docs.render.com` 404s on its reference pages and the OpenAPI its own docs name is not at that URL | A spec that resolves |

The useful observation is that four of these five are ReadMe or JavaScript
portals, and the fifth is a login. A provider that publishes a spec file --
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

### Remaining

32 Recipes, 78 declarations. The method that works: find the provider's own
machine-readable description, read the parameter names out of it, then write a
case that *sends* them. Asserting only the response is not enough -- Pub/Sub's
`cursor_param` could be renamed to `cursor` with every case still passing,
because nothing sent a token back until a second-page case did.

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
