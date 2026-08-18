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
| PostHog | Events, persons, feature flags |
| Mixpanel | Events, profiles, exports |
| Amplitude | Events, cohorts, user properties |

## Models and inference

| Provider | Why |
|---|---|
| OpenAI | Responses, streaming, tool calls, rate limits, structured output |
| Anthropic | Messages, content blocks, tool use, streaming |
| Google Gemini | Generation, multimodal requests, safety responses |
| Pinecone | Indexes, namespaces, vectors, metadata filtering |
| Replicate | Predictions, async jobs, failures |
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
| Increase | Assess — the state machine is the product. A transfer is pending, then submitted, then it settles or returns days later, and a return can arrive after everything downstream has treated it as done |
| Modern Treasury | Assess — payment orders, ledgers, reconciliation. Double-entry means an amount appears twice with opposite signs and summing naively gives zero |
| Column | Assess — ACH returns and their reason codes |
| Dwolla | Assess — transfers, funding sources, micro-deposit verification |
| Unit | Assess — accounts, cards, authorisations. An authorisation is not a transaction and the two have separate ids |
| Marqeta | Assess — card issuing, JIT funding, authorisation webhooks with a response deadline |
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
| Bugsnag | Assess — errors, events, the difference between an error and its occurrences |
| Honeybadger | Assess — faults and notices |
| Grafana Cloud | Assess — dashboards, alert rules, the Prometheus-shaped query API |
| Honeycomb | Assess — datasets, triggers, query results |
| Better Stack | Assess — logs, uptime monitors, incidents |
| Heap | Assess — events and user properties |
| LogRocket | Assess — sessions and issues |

## Hosting, deployment and package registries

| Provider | Why |
|---|---|
| Netlify | Assess — deploys are asynchronous and a deploy id exists before the site is live |
| Render | Assess — services, deploys, the build-then-live gap |
| Fly.io | Assess — machines, apps, the Machines API against the older platform API |
| Heroku | Assess — the API still uses `Accept: application/vnd.heroku+json; version=3`, so a missing header is a different response rather than an error |
| Linode | Assess — instances and the async provisioning lifecycle, beside Vultr and DigitalOcean |
| Hetzner Cloud | Assess — servers, actions. Every mutation returns an action object you have to poll, rather than the thing you changed |
| Scaleway | Assess — instances and object storage |
| Docker Hub | Assess — repositories, tags, the rate limit that is counted per IP and not per token |
| npm registry | Assess — packages, versions, dist-tags. Unpublished versions leave a tombstone |
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
| Svix | Assess — the delivery-attempt model, exponential backoff, the endpoint that gets disabled after repeated failure |
| Knock | Assess — workflows, preferences, batching |
| Courier | Assess — routing across channels and the fallback that fires silently |
| Novu | Assess — subscribers, workflows, digest |
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
| Braze | Assess — users, campaigns, the export API being asynchronous |
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
| Persona | Assess — inquiries, verifications, the decision that arrives by webhook minutes later |
| Onfido | Assess — applicants, checks, reports. A check is complete and its report can still be `consider` rather than `clear` |
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
| Firecrawl | Assess — crawl jobs are asynchronous and partial results are readable before the job finishes |
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
| Orb | Assess — usage-based billing where an invoice is not final until the period closes |
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
