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
| Klarna | Blocked on documentation, like Sift. The premise is good -- an authorization token expires while the order it approved has not been created, so a customer can be approved and unpayable -- and docs.klarna.com renders its API reference through a portal that serves navigation rather than schemas, so the response bodies cannot be read or cited. Revisit with an account or a published spec |
| ~~Airwallex~~ | Shipped. The settlement currency is not the charge currency and the rate was fixed at a moment nobody picked, there is a separate balance per currency so a payout can fail on a funded account, and a partial capture leaves two amounts different forever |
| ~~Column~~ | Shipped. A notification of change is not a failure, and R01 and R07 are two characters apart with opposite obligations |
| Moov | Transfers across rails, where the rail decides the settlement time and nothing else does |
| Unit | Assessed and left. That is the Marqeta Recipe's headline, already shipped: an authorization is a hold, the clearing is a second object pointing back at it, and the amounts legitimately differ. A Unit Recipe would need to earn its place on the banking-as-a-service parts -- the application lifecycle, the counterparty model -- rather than on the card shape |

### Billing that accrues

| Provider | Why |
|---|---|
| ~~Orb~~ | Shipped. A draft is a running total, a closed period can still be amended, and money is a decimal string |
| Metronome | Usage events are deduplicated inside a window, so sending the same event twice is sometimes one event and sometimes two |
| Lago | Open source, and self-hosted behaves differently from cloud. Assess whether that difference is modellable or a reason not to |
| Lemon Squeezy | Merchant of record, so the tax is theirs and the order total is not what arrives. That is FastSpring's headline, shipped this cycle, and a second Recipe saying it would add a name rather than a shape. Worth doing only for what differs: the licence-key API and the store-scoped identifiers |
| ~~RevenueCat~~ | Shipped. There is no `is_active` field and RevenueCat's own guidance is to read one -- it is an SDK property, so the moment the question moves to a server somebody writes the comparison by hand and the advice stops applying. Four active entitlements in the fixture, active for four different reasons: cancelled, lifetime (`expires_date: null`, which every naive comparison reads as expired), failing to pay inside a grace period, and somebody else's family purchase on a trial. Entitlements are keyed by your names and subscriptions by the stores', and the endpoint is a GET that creates: 200 found, 201 invented |
| Recharge | Still open, and deprioritised for a reason worth writing down: the drift is between Recharge and Shopify, so telling the story needs both systems and a Recipe can only hold one. What a single Recipe could carry is narrower -- a subscription that is ACTIVE while its charges have reached MAX_RETRIES_REACHED, and external ids that are Shopify's numbers held as strings |
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
| Grafana Cloud | Assess — dashboards, alert rules, and a Prometheus-shaped query API that is not REST |
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

## Assessed and deliberately not done

| Provider | Why not |
|---|---|
| ~~Linear~~ | Shipped, on the GraphQL support ShipHero brought. Priority counts down and zero is not the top -- Linear's own words are 0 = No priority, 1 = Urgent, 4 = Low -- so sorting ascending puts untriaged issues above the ones on fire and descending puts Low above Urgent. Plus: a state's `name` belongs to the team and its `type` does not, three of the seven types close an issue and `duplicate` is the forgotten one, a connection carries the same list twice as `edges` and `nodes`, and `number` is team-scoped so two issues are both 123 |
| ~~Attio~~ | Shipped. It was recorded here as GraphQL-only, which was simply wrong -- it is REST and publishes OpenAPI 3.1 at `https://api.attio.com/openapi/api`. Its records are queried by POST with the paging in the body, which is what found that bug in four other Recipes |
| New Relic | NerdGraph is GraphQL-only. Same reason again |
| Railway | GraphQL-only. Same reason |
| Temporal Cloud | gRPC rather than HTTP. The format describes HTTP surfaces and nothing here would be a Temporal client |

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

### Remaining

54 Recipes, 129 declarations. The method that works: find the provider's own
machine-readable description, read the parameter names out of it, then write a
case that *sends* them. Asserting only the response is not enough -- Pub/Sub's
`cursor_param` could be renamed to `cursor` with every case still passing,
because nothing sent a token back until a second-page case did.
