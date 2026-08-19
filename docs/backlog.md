# Provider backlog

The providers queued for Recipes, roughly in the order they are worth doing.
This is a working list, not a promise, and it exists so the queue survives
between sessions rather than living in somebody's head.

Two rules apply to everything here:

1. **A Recipe has to earn its place.** The test is whether the provider has
   behaviour worth reproducing — a third state nobody branches on, a shape that
   breaks a client, an asynchronous outcome, a failure that only exists in
   production. Linear and Attio were assessed and left out rather than
   approximated, and that stays the standard. A thin Recipe that raises the
   count is worse than an honest gap.
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
| Google Cloud Storage | Objects, signed URLs, permissions |
| ~~Google Pub/Sub~~ | Shipped. Base64 bodies, ack deadlines, delivery attempts |
| Azure Blob Storage | Containers, blobs, SAS and auth behaviour |
| Cloudflare R2 | Sits beside the existing Cloudflare Recipe |
| Vultr | Instances, block storage, the async provisioning lifecycle |

The queueing ones are the most valuable and the least like anything shipped so
far: at-least-once delivery, a message that comes back after a visibility
timeout, and a dead-letter queue are behaviours no fake reproduces by accident.

**XML is the blocker for several of these.** S3, SNS and Azure Blob answer in
XML, not JSON, and so do UPS, FedEx, USPS and Avalara further down. Cauldron
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
| Shopify GraphQL Admin API | Separate from the REST Shopify Recipe, not a replacement |
| BigCommerce | Orders, products, customers, webhooks |
| Magento / Adobe Commerce | Carts, orders, inventory, customers |
| Etsy | Listings, orders, receipts, inventory |
| eBay | Listings, fulfilment, orders |
| Amazon Selling Partner API | Orders, inventory, reports, throttling. Very large |

## Shipping and logistics

| Provider | Why |
|---|---|
| ShipStation | Shipments, labels, carriers, tracking |
| EasyPost | Rates, shipments, labels, trackers |
| ShipEngine | Rates, labels, tracking, carrier errors |
| AfterShip | Tracking states, checkpoints, carrier behaviour |
| Easyship | Rates, shipments, labels |
| USPS | Addresses, rates, tracking |
| UPS | Rating, labels, tracking |
| FedEx | Rating, shipping, tracking |
| DHL | Shipment lifecycle and international customs states |

## Productivity and identity

| Provider | Why |
|---|---|
| Microsoft Graph | Mail, calendars, users, files. Enormous surface area |
| ~~Google Calendar~~ | Shipped. Listings return the series not the occurrence, all-day events have no time, cancelled instances are almost empty |
| ~~Gmail~~ | Shipped. A listing carries no message, headers are an array, no read flag, trash is not delete |
| Google Drive | Files, permissions, revisions, shared drives |
| Microsoft OneDrive | Files, sharing, async operations |
| Microsoft Teams | Channels, messages, members |
| GitHub Actions | Could extend the GitHub Recipe, probably its own |
| ~~WorkOS~~ | Shipped. Inactive-but-present users, per-IdP raw attributes, draft connections |
| ~~Stytch~~ | Shipped. Per-factor verification, session factors, invited members |
| Descope | Authentication flows and identities |
| FusionAuth | Users, tenants, applications |
| Keycloak | Realms, users, clients, tokens |

## Analytics and flags

| Provider | Why |
|---|---|
| ~~LaunchDarkly~~ | Shipped. Per-environment state, variations as indices |
| ~~PostHog~~ | Shipped. A property is nested under properties and a flag is a string, a boolean or false, all in the same field |
| Mixpanel | Events, profiles, exports |
| Amplitude | Events, cohorts, user properties |

## Models and inference

| Provider | Why |
|---|---|
| OpenAI | Responses, streaming, tool calls, rate limits, structured output |
| Anthropic | Messages, content blocks, tool use, streaming |
| Google Gemini | Generation, multimodal requests, safety responses |
| Pinecone | Indexes, namespaces, vectors, metadata filtering |
| ~~Replicate~~ | Shipped. A created prediction has no output property at all rather than a null one, succeeded is not the same as produced something, a cold start is a minute with no signal but a boot_time, and the output is a link to a file deleted after an hour |
| Hugging Face | Inference endpoints and model responses |
| ElevenLabs | Speech generation, voices, async jobs |
| ~~Deepgram~~ | Shipped. Four-level nesting, seconds against milliseconds, err_code |
| ~~AssemblyAI~~ | Shipped. Failure at HTTP 200, present-and-null fields, ms against s |

Streaming is the open question here. Cauldron serves whole responses, and
server-sent events are a different shape. Whether that is modelled, or stated
as a gap, needs deciding before the first of these ships rather than after.

## Data and messaging platforms

| Provider | Why |
|---|---|
| Supabase | Auth, storage, database REST and realtime |
| Redis Cloud | Keys, TTLs, streams, pub/sub |
| Upstash | Redis REST, queues, rate limits |
| MongoDB Atlas | Clusters, users, projects, the data API |
| Neon | Branches, databases, endpoints |
| PlanetScale | Databases, branches, deploy requests |
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
| Dropbox Sign | Signature requests and callbacks |
| Adobe Acrobat Sign | Agreements and the signing lifecycle |
| Ironclad | Contracts and approval workflows |
| DocSpring | PDF generation and templates |

## Tax, accounting and payroll

| Provider | Why |
|---|---|
| ~~TaxJar~~ | Shipped. Decimal rates, present-and-zero jurisdictions, nexus |
| Avalara | Sales tax calculations and filings |
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
| Vimeo | Uploads, videos, privacy |
| YouTube Data API | Videos, playlists, channels |

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
| Jira Service Management | Assess — requests, queues, SLAs. An SLA clock pauses and the paused time is not the elapsed time |
| Confluence Cloud | Bodies are XHTML storage format rather than Markdown or HTML, and the `representation` field decides how the same string parses. Sending the wrong one stores markup as literal text |
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
| RevenueCat | Assess — mobile subscriptions and entitlements, where the entitlement is the thing an app reads and the subscription is the thing that renews, and they can disagree for a whole billing period |
| Firebase Auth | Assess — identity, and the emulator Google already ships for it is the question: a Recipe has to be better than the official one to earn its place |
| Coinbase Commerce | Assess — charges that expire, underpayment and overpayment as distinct outcomes |
| Circle | Assess — USDC transfers and their settlement states |

## Feature flags and experimentation

| Provider | Why |
|---|---|
| Statsig | Assess — gates, experiments, exposure logging |
| GrowthBook | Assess — features and their environment overrides |
| Flagsmith | Assess — flags per environment, identity overrides |
| Split | Assess — treatments and targeting rules |
| ConfigCat | Assess — flags, segments, config JSON delivery |
| Unleash | Assess — toggles, strategies, the client and admin APIs being different shapes |

## Observability and product analytics

| Provider | Why |
|---|---|
| ~~Bugsnag~~ | Shipped. An error is not its occurrences, and the counts are on the error |
| Honeybadger | Assess — faults and notices |
| Grafana Cloud | Assess — dashboards, alert rules, the Prometheus-shaped query API |
| Honeycomb | Assess — datasets, triggers, query results |
| Better Stack | Assess — logs, uptime monitors, incidents |
| Heap | Assess — events and user properties |
| LogRocket | Assess — sessions and issues |

## Hosting, deployment and package registries

| Provider | Why |
|---|---|
| ~~Netlify~~ | Shipped. A deploy id exists long before the site is live |
| Render | Assess — services, deploys, the build-then-live gap |
| Fly.io | Assess — machines, apps, the Machines API against the older platform API |
| Heroku | Assess — the API still uses `Accept: application/vnd.heroku+json; version=3`, so a missing header is a different response rather than an error |
| Linode | Assess — instances and the async provisioning lifecycle, beside Vultr and DigitalOcean |
| Hetzner Cloud | Assess — servers, actions. Every mutation returns an action object you have to poll, rather than the thing you changed |
| Scaleway | Assess — instances and object storage |
| ~~Docker Hub~~ | Shipped. The rate limit is counted per IP, not per token |
| ~~npm registry~~ | Shipped. Unpublished versions leave a tombstone, and dist-tags are a flat map |
| PyPI | Assess — the JSON API is read-only and upload is a separate protocol entirely |
| Codecov | Assess — coverage reports and the commit they attach to |
| SonarCloud | Assess — issues, quality gates, the gate status arriving after analysis finishes |

## Secrets and configuration

| Provider | Why |
|---|---|
| HashiCorp Vault | Assess — KV v1 and KV v2 have different paths and different response shapes for the same secret, and v2 nests the value under `data.data` |
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
| Hookdeck | Assess — connections, retries, the paused destination |
| Inngest | Assess — events, functions, step state |
| Trigger.dev | Assess — runs, tasks, the resumed run |

## Scheduling and real-time

| Provider | Why |
|---|---|
| Cal.com | Assess — bookings, availability, the slot that is free when you read it and taken when you book it |
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
| Apify | Assess — actor runs, datasets, the run that succeeds with zero items |
| ScrapingBee | Assess — a failed fetch is a successful API call |
| Browserless | Assess — sessions and timeouts |

## More commerce and billing

| Provider | Why |
|---|---|
| Toast | Assess — orders, checks, the restaurant's business day not matching the calendar day |
| Lightspeed | Assess — the Retail and Restaurant APIs are unrelated products sharing a brand |
| Clover | Assess — merchants, orders, payments |
| Recharge | Assess — subscriptions on top of Shopify, so two sources of truth for one order |
| Lemon Squeezy | Assess — orders, subscriptions, the merchant-of-record tax handling |
| Gumroad | Assess — products and sales |
| Polar | Assess — subscriptions and benefits |
| ~~Orb~~ | Shipped. An invoice is not final until the period closes |
| Metronome | Assess — usage events and their deduplication window |
| Lago | Assess — open source usage billing, self-hosted and cloud behaving differently |
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

Two real bugs fell out of the verification. DigitalOcean was ignoring
`per_page` and its own conformance case had been asserting the emulator's
parameter name rather than the provider's. Twilio capitalises its parameters
and both spellings were being ignored.

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
| Klarna | The session expires between authorise and capture, so an order can be approved and unpayable |
| ~~Airwallex~~ | Shipped. The settlement currency is not the charge currency and the rate was fixed at a moment nobody picked, there is a separate balance per currency so a payout can fail on a funded account, and a partial capture leaves two amounts different forever |
| ~~Column~~ | Shipped. A notification of change is not a failure, and R01 and R07 are two characters apart with opposite obligations |
| Moov | Transfers across rails, where the rail decides the settlement time and nothing else does |
| Unit | An authorisation is not a transaction and the two carry different identifiers for the same purchase |

### Billing that accrues

| Provider | Why |
|---|---|
| ~~Orb~~ | Shipped. A draft is a running total, a closed period can still be amended, and money is a decimal string |
| Metronome | Usage events are deduplicated inside a window, so sending the same event twice is sometimes one event and sometimes two |
| Lago | Open source, and self-hosted behaves differently from cloud. Assess whether that difference is modellable or a reason not to |
| Lemon Squeezy | Merchant of record, so the tax is theirs, the payout is net, and the order total is not what arrives |
| RevenueCat | The entitlement an app reads and the subscription that renews are different objects and can disagree for a whole billing period |
| Recharge | Subscriptions on top of Shopify, so one order has two sources of truth and they drift |
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
| Twilio Verify | A code expires and a check consumes an attempt, so verifying twice with the right code fails the second time |
| Customer.io | The Track and App APIs are separate hosts with separate credentials, which is the sort of thing that works in one environment and not the other |
| ~~Braze~~ | Shipped. The export answers 201 with a prefix and no users; the file lands in cloud storage minutes later, so a test reading users off that response reads nothing forever |
| Brevo | Contacts and transactional email share one quota, so sending mail exhausts the budget for reading contacts |
| Kit | Assess — subscribers, sequences, tags |
| Beehiiv | Assess — publications, posts, subscribers |

### Identity and risk

| Provider | Why |
|---|---|
| AWS Cognito | The SDK is the API, and the token lifecycle is the whole integration |
| Ory | Assess — identities, sessions, flows. The flow object is the unusual part |
| Kinde | Assess — users, organizations, feature flags in one product |
| ~~Persona~~ | Shipped. completed is not approved, needs_review is neither, and nothing is at the top level because it is JSON:API |
| ~~Onfido~~ | Shipped. complete is not clear, consider is neither a pass nor a failure, and the reason lives on the report rather than the check |
| Sift | The score is computed server side and changes without any request from you |

### Storage and media

| Provider | Why |
|---|---|
| Backblaze B2 | S3-compatible and not identical, which is the interesting part: the differences are where a working integration breaks |
| ~~Uploadcare~~ | Shipped. A file exists before it is stored and unstored files are deleted after twenty-four hours, so the same code works on one project and loses files on another; a removed file still answers 200 with everything intact |
| ImageKit | Assess — transformations and delivery, beside Cloudinary and Imgix |
| Livepeer | Assess — streams and recordings, asynchronous throughout |

### Data and search

| Provider | Why |
|---|---|
| Qdrant | Vector search, where a filter changes which vectors are even considered, so the same query with a filter is not a subset of the same query without one |
| Weaviate | Assess — objects, classes, hybrid search |
| Typesense | Assess — collections, documents, search parameters |
| Convex | Assess — the query and mutation model is not REST-shaped, so this may belong in the not-done table |

### Hosting and deployment

| Provider | Why |
|---|---|
| ~~Netlify~~ | Shipped. A deploy has an id before it has a site, ready is not published, and a missing URL means two different things |
| Render | Services, deploys, and the gap between build finished and live |
| Fly.io | The Machines API and the older platform API are different shapes for the same account |
| Heroku | The API still wants `Accept: application/vnd.heroku+json; version=3`, so a missing header is a different response rather than an error |
| Hetzner Cloud | Every mutation answers with an action to poll rather than with the thing you changed |

### Observability and flags

| Provider | Why |
|---|---|
| ~~Bugsnag~~ | Shipped. An error carries counts and no stack trace, fixed is not terminal, and severity is a different question from unhandled |
| Grafana Cloud | Assess — dashboards, alert rules, and a Prometheus-shaped query API that is not REST |
| Better Stack | Assess — logs, uptime monitors, incidents |
| ~~PostHog~~ | Shipped. A flag definition is not what a user gets, nought per cent is not inactive, and capture says the same thing whatever you send |
| Statsig | Assess — gates, experiments, exposure logging |
| Flagsmith | Assess — flags per environment with identity overrides, so the same flag has many answers |

### Webhook infrastructure

| Provider | Why |
|---|---|
| ~~Svix~~ | Shipped. Accepted is not delivered, the outcome lives on the attempts, a failing endpoint gets disabled and nothing announces it |
| Inngest | Assess — events, functions, step state, where a step is durable and a function is not |

### Web data

| Provider | Why |
|---|---|
| ~~Firecrawl~~ | Shipped. A running crawl already has results, total is an estimate completed never reaches, and a 200 can carry a failed fetch |
| Apify | A run can succeed with zero items, which is not the same as failing and is handled as if it were |
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
| Tradier | The same order in two accounts is two objects, and the account is a path segment rather than a field |
| Polygon.io | Aggregates are bucketed and a bucket with no trades is absent rather than zero, so a chart drawn from the response has invisible gaps |
| Finnhub | The free tier truncates history rather than refusing, so a backtest runs on less data than it asked for and says nothing |

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
| FastSpring | One purchase is an order, a subscription and a sale, with three different ids, and which one a webhook names depends on the event |

### Notification and messaging infrastructure

| Provider | Why |
|---|---|
| Xendit | A virtual account exists before it can receive money, and the callback is the only signal that it can |
| ~~Midtrans~~ | Shipped. transaction_status and fraud_status are two fields and a payment is only safe when both agree; capture plus challenge means the card was charged and the funds are held; a bank transfer never passes through capture at all |
| PayU | The same merchant has different endpoints per country and the response fields differ between them |
| Gorgias | A ticket and its messages page separately, so a message can arrive between the two reads and belong to neither |
| Kustomer | A conversation is a customer timeline rather than a thread, so assignment is a state of the conversation and not of any message |
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
| Langfuse | Traces are ingested asynchronously and are not readable immediately after being written |

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

The table went from 12 Recipes to 91 in one pass. These twenty-six are left,
and each needs its published client libraries checked rather than remembered,
because a package name written from memory is exactly the guess the detection
rule forbids.

| Recipe | Note |
|---|---|
| Alpaca | alpaca-trade-api and alpaca-py exist; the Node and Go names need checking |
| Bandwidth | Official SDKs exist for several languages under names that have changed at least once |
| Basecamp | No official SDK. Community clients only, so a mapping may not be justifiable |
| Bill.com | No widely used client |
| Braze | Server-side use is usually raw HTTP; the SDKs are mobile |
| Buildkite | go-buildkite is the obvious one |
| Column | No official SDK |
| Deel | No official SDK |
| Docker Hub | No SDK; usually the docker CLI or raw HTTP |
| Documenso | Has a TypeScript SDK |
| Greenhouse | Community clients only |
| Gusto | Has an official Node SDK |
| Help Scout | Community clients |
| Increase | Official SDKs in several languages |
| Marqeta | Official SDKs exist and are generated |
| Mercury | No official SDK |
| Modern Treasury | Official SDKs in several languages |
| npm registry | Usually npm itself or raw HTTP |
| Onfido | Official SDKs in several languages |
| Orb | Official SDKs in several languages |
| Persona | Community clients |
| Ramp | No official SDK |
| Shortcut | Community clients |
| Svix | Official SDKs in several languages |
| WooCommerce | Official REST clients in PHP and JS |
| WordPress | Usually raw HTTP or wp-cli |

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

Sixteen Recipes are fixed and each carries a case asserting an unquoted
identifier, so removing the declaration fails something. Three of them already
had cases asserting the quoted form, which is to say three cases were pinning
the bug in place.

The rest of the `numeric` resources are below. Which of the two a provider
sends has to be read from its documentation rather than assumed: sending a
number where the provider sends a string is as wrong as the thing this fixes.

| Recipe | Expected | Note |
|---|---|---|
| RingCentral | needs checking | Message and extension ids may be numbers or numeric strings |
| Postmark | needs checking | Bounce ID is a number; the casing of the field needs confirming too |
| SendGrid | needs checking | Suppression ids are numbers on some endpoints |
| Documenso | needs checking | Document and recipient ids are numbers |
| Intercom | needs checking | Conversation ids are quoted, but the older API differed |
| HubSpot | **string** | Contact, deal and company ids are quoted, which is why the default stays string |
| Jira | **string** | Issue id is a quoted number; the key is the readable identifier |
| QuickBooks | **string** | Id is quoted everywhere in the JSON API |
| DocuSign | **string** | recipientId is quoted |

The five marked **string** are correct as they stand and are listed so nobody
"fixes" them.

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

## Assessed and deliberately not done

| Provider | Why not |
|---|---|
| Linear | GraphQL-only. A REST-shaped Recipe would be an approximation of a different API, and saying so is more useful than shipping one |
| Attio | Same reason |
| New Relic | NerdGraph is GraphQL-only. Same reason again |
| Railway | GraphQL-only. Same reason |
| Temporal Cloud | gRPC rather than HTTP. The format describes HTTP surfaces and nothing here would be a Temporal client |

Monday.com belongs here too, on the same grounds, unless the format grows a way
to describe a GraphQL surface honestly.
