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
| PayPal | Payments, captures, refunds, disputes, webhooks |
| Braintree | Transactions, customers, payment methods |
| Mollie | Payment lifecycle and asynchronous states |
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
| Google Calendar | Events, attendees, recurrence, cancellations |
| Gmail | Threads, labels, drafts, messages |
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

## Assessed and deliberately not done

| Provider | Why not |
|---|---|
| Linear | GraphQL-only. A REST-shaped Recipe would be an approximation of a different API, and saying so is more useful than shipping one |
| Attio | Same reason |

Monday.com belongs here too, on the same grounds, unless the format grows a way
to describe a GraphQL surface honestly.
