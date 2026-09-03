// Shaping: turning a stored record into the body a provider would send.

package runtime

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// nestedValue walks a dotted path to a leaf in a decoded body.
func nestedValue(body store.Record, path string) (any, bool) {
	var current any = map[string]any(body)

	for _, segment := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = object[segment]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

// trim reduces a record to the fields a route answers with. The trim happens
// before shaping, so a kept field still nests where the Recipe says it does.
//
// Jira's create hands back an id and a key and nothing else, so anything
// reading created.fields.summary gets undefined. That is the shape this was
// written for, and for a long time it was the only one it served: the trim
// lived inside the create path, so a Recipe declaring returns on a get, an
// update or a list was declaring nothing. The validator accepted it, the
// emulator ignored it, and the Recipe read as though a smaller body had been
// described when the full record was still going out.
//
// A listing is where that matters most. Braze's /campaigns/list hands back
// five properties and /campaigns/details hands back a different, larger set,
// so code reading channels off a list entry gets undefined against the real
// API and a populated array against an emulator that ignored the trim. GitHub
// does the same thing between its issue list and its issue detail, and so do
// enough providers that a listing being a summary should be the assumption
// rather than the surprise.
func trim(spec recipe.Route, record store.Record) store.Record {
	if len(spec.Returns) == 0 {
		return record
	}

	kept := make(store.Record, len(spec.Returns))

	for _, field := range spec.Returns {
		if value, ok := record[field]; ok {
			kept[field] = value
		}
	}

	return kept
}

// present renames the identifier to the property the provider actually uses.
// The store keeps every record keyed by "id" so fixtures and internal lookups
// stay uniform; only the wire shape changes, which is where it matters.
// idAs renames the identifier on one route alone, for the providers that call
// it one thing there and another everywhere else. Empty means the resource's
// own name, which is the ordinary case.
func (s *Sandbox) present(resource string, record store.Record, idAs string) store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return record
	}

	if idAs != "" {
		spec.ID.Field = idAs
	}

	// "-" is not a rename but a suppression: the identifier is how Cauldron
	// finds the record, not something the provider puts on the wire.
	hidden := spec.ID.Field == "-"
	renamed := spec.ID.Field != "" && spec.ID.Field != "id" && !hidden

	// A provider whose identifier is a JSON number. The store keeps it as a
	// string, because that is the only form every style shares and the only
	// form a path parameter arrives in, so the conversion belongs here at the
	// edge rather than anywhere a lookup happens.
	retyped := spec.ID.Type == "number"

	if !renamed && !hidden && !retyped && !s.nests(spec) {
		return record
	}

	out := make(store.Record, len(record))

	for key, value := range record {
		if key == "id" && hidden {
			continue
		}

		if key == "id" && retyped {
			setPath(out, identifierName(spec), numeric(value))

			continue
		}

		if key == "id" && renamed {
			// setPath, so a dotted name nests. Contentful keeps the
			// identifier at sys.id rather than at the top level.
			setPath(out, spec.ID.Field, value)
			continue
		}

		// A field declared with "in" moves under that sub-object on the wire.
		// HubSpot reads contact.properties.email, and a client written against
		// it finds nothing at the top level.
		//
		// "-" moves it nowhere: the record holds it and the wire never sees
		// it. That is what a route's scope needs, because a partition living
		// in the path is not repeated in the body by most providers, and the
		// only way to say so before this was to name every other field in the
		// route's returns.
		if field, declared := spec.Fields[key]; declared && field.In == "-" {
			continue
		} else if declared && field.In != "" {
			// A dotted name nests twice. Brex puts a card's limit at
			// spend_controls.limit.amount, and treating the name as one
			// literal key produced a flat "spend_controls.limit" key that no
			// provider sends and nothing was checking.
			nested := nestedObject(out, field.In)

			// The wire name, which is the field's own unless it says
			// otherwise. Without that, a resource needing both title.rendered
			// and content.rendered had to name one of them something else, and
			// the name it was given leaked out as title.title_rendered.
			nested[field.WireName(key)] = value

			continue
		}

		// A declared constant at a dotted name nests, the same way a renamed
		// identifier and a field's "in" already do. Attio keeps two of a
		// record's three identifier UUIDs beside the third -- the whole id is
		// an object -- and writing id.workspace_id as one literal key
		// produced a shape no provider sends. That is the fifth mechanism in
		// this file to have made that mistake, and the reason it keeps
		// happening is that a dotted name is only a path in the places
		// somebody remembered to make it one.
		if _, declared := spec.Fields[key]; !declared && recipe.IsPath(key) {
			setPath(out, key, value)

			continue
		}

		// A top-level field can name itself something else on the wire, for
		// the one thing it cannot otherwise say: that the name it wants is
		// already the record's own. GitHub gives an issue a number and an id,
		// a path addresses it by the number, so the record keys on the number
		// and "id" is left for the field carrying GitHub's id.
		if field, declared := spec.Fields[key]; declared && field.As != "" {
			out[field.As] = value

			continue
		}

		out[key] = value
	}

	return out
}

// lookupNested walks a dotted path through an incoming request body.
func lookupNested(record store.Record, path string) (map[string]any, bool) {
	var current any = map[string]any(record)

	for _, segment := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = object[segment]
		if !ok {
			return nil, false
		}
	}

	object, ok := current.(map[string]any)

	return object, ok
}

// nestedObject walks a dotted path, creating the objects it passes through,
// and returns the one the field belongs in.
//
// Most providers nest one level, which is a single key. Brex nests twice for a
// card's spending limit, and a dotted name is how a Recipe says so.
func nestedObject(out map[string]any, path string) map[string]any {
	current := out

	for _, segment := range strings.Split(path, ".") {
		name, indexes := splitIndex(segment)

		if len(indexes) == 0 {
			next, _ := current[segment].(map[string]any)
			if next == nil {
				next = map[string]any{}
				current[segment] = next
			}

			current = next

			continue
		}

		// An array with objects in it. RingCentral's to is a list of one
		// recipient, and without this the whole segment became a literal key
		// spelled "to[0]" — a shape no provider sends, produced silently, and
		// invisible to a conformance suite that has no case naming it. That is
		// the third time a key like this has shipped, so the validator now
		// refuses an index the runtime cannot honour rather than trusting the
		// next Recipe author to notice.
		current = descendIndexed(current, name, indexes)
	}

	return current
}

// descendIndexed walks into an array-valued key, growing the array as needed,
// and returns the object at the innermost index.
func descendIndexed(current map[string]any, name string, indexes []int) map[string]any {
	held, _ := current[name].([]any)
	if held == nil {
		held = []any{}
	}

	for depth, index := range indexes {
		for len(held) <= index {
			held = append(held, map[string]any{})
		}

		if depth == len(indexes)-1 {
			object, _ := held[index].(map[string]any)
			if object == nil {
				object = map[string]any{}
				held[index] = object
			}

			current[name] = held

			return object
		}

		inner, _ := held[index].([]any)
		if inner == nil {
			inner = []any{}
		}

		current[name] = held
		held = inner
	}

	return map[string]any{}
}

// splitIndex separates a path segment from any [n] suffixes on it.
func splitIndex(segment string) (string, []int) {
	open := strings.Index(segment, "[")
	if open < 0 {
		return segment, nil
	}

	name := segment[:open]
	rest := segment[open:]

	var indexes []int

	for rest != "" {
		if !strings.HasPrefix(rest, "[") {
			return segment, nil
		}

		close := strings.Index(rest, "]")
		if close < 0 {
			return segment, nil
		}

		n, err := strconv.Atoi(rest[1:close])
		if err != nil || n < 0 {
			return segment, nil
		}

		indexes = append(indexes, n)
		rest = rest[close+1:]
	}

	return name, indexes
}

// nests reports whether any of a resource's fields live under a sub-object.
//
// A field marked "-" is not sent, and flatten below looks for it under a key
// literally spelled that, finds nothing, and moves on. Excluding it here
// changed no behaviour, so it is not excluded: a guard that guards nothing is
// the thing this project keeps telling Recipes not to write.
func (s *Sandbox) nests(spec recipe.Resource) bool {
	for _, field := range spec.Fields {
		// A field that only renames itself still needs present to run, or the
		// rename is declared and never applied -- which is the quiet kind of
		// wrong this whole file keeps finding.
		if field.In != "" || field.As != "" {
			return true
		}
	}

	// A dotted constant nests too, and a resource can have one without any
	// field declaring an "in" at all -- which is exactly Attio's records,
	// where the nesting is entirely in the identifier and the constants
	// beside it.
	for name := range spec.Constants {
		if recipe.IsPath(name) {
			return true
		}
	}

	return false
}

// flatten is the inverse of present for an incoming request: a client sends
// {"properties": {"email": "..."}} and the store keeps fields flat.
func (s *Sandbox) flatten(resource string, record store.Record) store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok || !s.nests(spec) {
		return record
	}

	// Extracted separately, and applied after the wrappers are dropped,
	// because a field can be nested under a wrapper of its own name. ClickUp's
	// status is: the response carries {"status": {"status": "review", ...}},
	// so the field is named status and sits in status. Writing the value
	// before deleting the wrapper deleted the value.
	extracted := store.Record{}
	unpacked := map[string]bool{}

	for name, field := range spec.Fields {
		if field.In == "" {
			continue
		}

		nested, isObject := lookupNested(record, field.In)
		if !isObject {
			continue
		}

		unpacked[strings.Split(field.In, ".")[0]] = true

		if value, present := nested[field.WireName(name)]; present {
			extracted[name] = value
		}
	}

	// The sub-objects have been unpacked, so drop them rather than storing the
	// same values twice under two shapes.
	//
	// Only the ones actually unpacked. Dropping every declared wrapper name
	// discarded a plain value that happened to share one, which is how a
	// ClickUp task's status became unsettable: the real API takes a flat
	// {"status": "review"} on a write and answers with the object, and a
	// request in the shape the provider documents was silently a no-op. The
	// write returned 200 with the old status, so nothing looked wrong.
	for wrapper := range unpacked {
		delete(record, wrapper)
	}

	for name, value := range extracted {
		record[name] = value
	}

	return record
}

// presentAll renames identifiers across a page.
func (s *Sandbox) presentAll(resource string, records []store.Record) []store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok || ((spec.ID.Field == "" || spec.ID.Field == "id") && spec.ID.Type != "number" && !s.nests(spec)) {
		return records
	}

	out := make([]store.Record, 0, len(records))

	for _, record := range records {
		out = append(out, s.present(resource, record, ""))
	}

	return out
}

// countValue renders a total as the provider sends it. Docusign sends its
// counts as strings, and emitting a number would quietly fix a bug the caller
// has to handle for themselves.
func (s *Sandbox) countValue(spec recipe.ListResponse, total int) any {
	if spec.CountAsString {
		return strconv.Itoa(total)
	}

	return total
}

// resourceBody shapes a single object according to the Recipe's declared style.
// Shopify nests it under the singular resource name; most providers return the
// object itself, and a client written for one shape breaks on the other.
// The envelope is the route's, which is the Recipe's unless the route says
// otherwise: Datadog wraps a created event and not a created monitor, and
// Vercel wraps a domain and neither a project nor a deployment.
func (s *Sandbox) resourceBody(spec recipe.ResourceResponse, idAs, resource string, record store.Record) any {
	record = s.present(resource, record, idAs)

	success := s.recipe.Responses.Success.Fields

	if spec.Style != "wrapped" {
		// A bare array of one: [{...}], no key and no envelope, and still a
		// list. PoetryDB answers /title/Ozymandias that way, so client code
		// reads body[0] exactly as it would read Xero's Invoices[0] -- with
		// nothing to name it by.
		//
		// Array was read only on the wrapped path below, so a Recipe declaring
		// it here got the object back and nothing said otherwise. That is the
		// declared-and-ignored shape again, and PoetryDB is the first provider
		// that needed the unwrapped form.
		//
		// The Recipe-wide success fields have nowhere to go once the body is a
		// list rather than an object, which is why validation refuses the two
		// together rather than dropping one of them here.
		if spec.Array {
			return []store.Record{record}
		}

		if len(success) == 0 {
			return record
		}

		// Nowhere to put the envelope fields without colliding with the
		// object's own, so wrap rather than silently overwrite a real field.
		out := make(map[string]any, len(record)+len(success))

		for key, value := range record {
			out[key] = value
		}

		return withFields(out, success)
	}

	// The default is the resource's own singular name, which is what Shopify,
	// Slack, Square and Zendesk all use. When the object is wrapped in a list
	// the default is the plural collection name instead, because a list of one
	// is still a collection: Xero answers with Invoices, not Invoice.
	key := spec.Key
	if key == "" {
		key = resource

		if spec.Array {
			key = s.collectionName(resource, "")
		}
	}

	var payload any = record

	// Xero answers a request for one invoice with a list of one, so client
	// code reads Invoices[0]. An emulator returning the object directly lets
	// code ship that breaks against the real API on the first call.
	if spec.Array {
		payload = []store.Record{record}
	}

	// A dotted name nests, the same way it does for a list envelope, for an
	// error envelope and for the constants in withFields below.
	//
	// This was the one of the four that did not, and it failed quietly: a
	// Recipe wrapping a single record under data.tracking got a body with one
	// key that had a dot in its name, which is a shape no provider sends. The
	// Recipe validated, the sandbox answered 200, and only a conformance case
	// asserting the nested path noticed. AfterShip is the first Recipe to
	// need it -- its listing is data.trackings and its single record is
	// data.tracking -- and no shipped Recipe had a dotted resource key, so
	// nothing that exists changes shape.
	wrapped := map[string]any{}
	setPath(wrapped, key, payload)

	return withFields(wrapped, success)
}

// withFields stamps a provider's constant envelope fields onto a body. A dotted
// name nests, so response_metadata.next_cursor needs no second mechanism.
func withFields(body map[string]any, fields map[string]any) map[string]any {
	for name, value := range fields {
		// Copied, because setPath ends in an assignment and a nested constant
		// would otherwise go into the response by reference -- and anything
		// that writes into that path afterwards writes into the parsed
		// Recipe. limit_field, page_field and pages_field are applied after
		// this runs, so a Recipe declaring both a `meta` constant and
		// `page_field: meta.page` had its own constant rewritten by serving a
		// request:
		//
		//   before: {"source": "cauldron"}
		//   after:  {"source": "cauldron", "page": 3, "limit": 1}
		//
		// The next request then carried the previous one's numbers, on a
		// different route, and Reset does not undo it: it rewinds the store,
		// the clock, the faults, the log and the webhooks, and never touches
		// the Recipe. A long-lived serve was poisoned by one request.
		//
		// This is what store.Record.Clone was made deep to prevent. The same
		// shape reached here and was missed.
		setPath(body, name, store.DeepCopy(value))
	}

	return body
}

func setPath(body map[string]any, path string, value any) {
	head, rest, nested := strings.Cut(path, ".")

	if !nested {
		if name, indexes := splitIndex(head); len(indexes) > 0 {
			// asObject rather than a bare assertion, for the same reason the
			// descend below uses it: a store.Record is a named map[string]any
			// and does not match the unnamed one. A record placed at an
			// indexed leaf skipped this branch entirely and fell through to
			// the assignment underneath, so an envelope key of
			// output.completeTrackResults[0] produced a key literally spelled
			// "completeTrackResults[0]".
			//
			// That is the same shape of failure the comment further down
			// records, and this is the third place it has had to be fixed.
			// The pattern is the assertion, not the callers.
			if object, ok := asObject(value); ok {
				target := descendIndexed(body, name, indexes)
				for key, nestedValue := range object {
					setPath(target, key, nestedValue)
				}

				return
			}
		}

		// A declared constant must not destroy data already in the body.
		// Intercom's pages object carries both a declared type and a computed
		// next cursor, and whichever landed second used to erase the other.
		if object, ok := value.(map[string]any); ok {
			if existing, ok := asObject(body[head]); ok {
				for key, nestedValue := range object {
					setPath(existing, key, nestedValue)
				}

				return
			}
		}

		body[head] = value

		return
	}

	// An indexed segment, the same way nestedObject already handles one for a
	// field's "in". Without it a constant declared at data.boards[0].name
	// produced a key literally spelled "boards[0]" -- a shape no provider
	// sends, produced in silence, and the fourth time a key like that has
	// been written here.
	//
	// The asymmetry was the tell: a conformance case could already assert
	// data.boards[0].name and a Recipe could not emit it.
	if name, indexes := splitIndex(head); len(indexes) > 0 {
		setPath(descendIndexed(body, name, indexes), rest, value)

		return
	}

	child, ok := asObject(body[head])
	if !ok {
		child = map[string]any{}
		body[head] = child
	}

	setPath(child, rest, value)
}

// asObject reads a map out of a body whichever of its two spellings it is
// stored under.
//
// store.Record is a named map[string]any, and a type assertion to the unnamed
// one does not match it. setPath asserted the unnamed one alone, so descending
// into a value that happened to be a record failed -- and the failure was not
// a no-op: the branch below it replaced the record with a fresh empty map, so
// a constant declared at data.type erased the record sitting at data.
//
// Nothing shipped hit it, because a Recipe reaches this only by declaring a
// constant inside its own envelope key, and Lemon Squeezy is the first: JSON:API
// puts a type beside every record, at data.type, with the record itself at
// data. The runtime already knew about this exact trap one layer up -- the
// handler switches on both spellings before applying a route's constants --
// and this is the same trap one layer down.
func asObject(value any) (map[string]any, bool) {
	switch object := value.(type) {
	case map[string]any:
		return object, true
	case store.Record:
		return object, true
	}

	return nil, false
}

// listBody shapes a page according to the declared list style, which is the
// Recipe's unless the route overrode it.
func (s *Sandbox) listBody(spec recipe.ListResponse, page store.Page, limit int, resource, path, nextURL string) any {
	page.Records = s.presentAll(resource, page.Records)

	// Chargebee wraps every item under the resource's own name, so a client
	// reads list[0].subscription.id. Anyone indexing straight into the item
	// finds nothing at all.
	var items any = page.Records

	// A collection of identifiers rather than of records. DynamoDB's
	// ListTables sends TableNames as an array of strings and keeps the table
	// object for DescribeTable beside it; SQS's ListQueues does the same with
	// QueueUrls. Emitting objects there hands a client something it will
	// interpolate into a URL, and [object Object] is what it builds.
	if spec.EntryField != "" {
		values := make([]any, 0, len(page.Records))

		for _, record := range page.Records {
			if value, ok := nestedValue(record, spec.EntryField); ok {
				values = append(values, value)
			}
		}

		items = values
	}

	if spec.EntryStyle == "wrapped" {
		wrapped := make([]any, 0, len(page.Records))

		for _, record := range page.Records {
			wrapped = append(wrapped, map[string]any{resource: record})
		}

		items = wrapped
	}

	// A collection of one, sent as the object. Tradier's words: if you have a
	// single order, it will be returned as a JSON obj/dict whereas multiple
	// orders will be returned as an array.
	//
	// This runs after the entry wrapping above, so a Recipe that does both
	// collapses to the wrapped item rather than to the bare record, which is
	// what an XML-shaped API sending one child element would produce.
	if spec.CollapseSingle {
		if list, ok := items.([]any); ok && len(list) == 1 {
			items = list[0]
		}

		if list, ok := items.([]store.Record); ok && len(list) == 1 {
			items = list[0]
		}
	}

	switch spec.Style {
	case "map":
		// Keyed by identifier rather than ordered. Pusher answers with an
		// object of channel names, so looping over it as a list finds
		// nothing, and a channel nobody is on is absent from the object
		// entirely rather than present with a zero. The key is the
		// identifier, so it does not repeat inside the value.
		keyed := map[string]any{}

		// Under the name the identifier actually goes out as, which is not
		// always "id". presentAll has already renamed it by the time this
		// runs, so reading "id" literally found nothing for any resource
		// declaring id.field -- and finding nothing here meant skipping the
		// record, so the whole entry vanished from the response with no error
		// anywhere. A keyed list that silently drops its contents is the worst
		// shape of this bug: the caller gets a smaller object, not a failure.
		key := "id"
		if spec, ok := s.recipe.Resources[resource]; ok {
			key = identifierName(spec)
		}

		for _, record := range page.Records {
			name, _ := record[key].(string)
			if name == "" {
				continue
			}

			value := map[string]any{}

			for field, held := range record {
				if field == key {
					continue
				}

				value[field] = held
			}

			// PDBe answers {"4hhb": [ {...} ]}: the map is keyed by the
			// identifier and the value is an array holding one record. The
			// difference is not cosmetic -- res["4hhb"].title is undefined
			// and res["4hhb"][0].title is the answer -- and without this a
			// Recipe could only claim the object, which is the reading that
			// breaks against the real provider.
			if spec.EntryStyle == "list" {
				keyed[name] = []any{value}

				continue
			}

			keyed[name] = value
		}

		// "-" is no wrapper at all: the map is the whole response. CoinGecko
		// answers /simple/price with {"bitcoin": {...}, "ethereum": {...}} at
		// the top level, so there is no property to put it under and no name
		// the collection could take. The convention is the one id.field and
		// message_field already use, where "-" means the provider does not
		// send it.
		//
		// Recipe-wide success fields have nowhere to go here, the same way
		// they have nowhere to go in a bare array, and a validation rule says
		// so rather than letting the runtime drop one shape in silence.
		if spec.Key == "-" {
			return keyed
		}

		body := map[string]any{}
		setPath(body, s.collectionName(resource, spec.Key), keyed)

		return withFields(body, spec.Fields)
	case "bare":
		// GitHub and friends return the array itself, with paging in headers.
		// A caller doing json.Unmarshal into a slice must not receive an object.
		return items
	case "tuple":
		// Two elements: the paging object first, the records second. The World
		// Bank's v2 API answers [{page, pages, per_page, total}, [...]], so
		// body[1] is the collection and body[0] is everything about it.
		//
		// It is worth having as its own style because of what the failure
		// looks like. A bad country code answers [{message: [...]}] -- one
		// element, not two -- with the same HTTP 200, so body[1] is the
		// collection on success and undefined on failure, and the length of
		// the outer array is the only structural signal that anything went
		// wrong.
		meta := map[string]any{}

		if spec.CountField != "" {
			setPath(meta, spec.CountField, s.countValue(spec, countTotal(spec, page, limit)))
		}

		if spec.PageCountField != "" {
			setPath(meta, spec.PageCountField, s.countValue(spec, len(page.Records)))
		}

		meta = withFields(meta, s.recipe.Responses.Success.Fields)

		return []any{withFields(meta, spec.Fields), items}
	case "wrapped":
		// A dotted key nests, so a collection can sit two levels down.
		// Segment answers with data.sources rather than a top-level array.
		body := map[string]any{}

		// SQS leaves the key out entirely when there is nothing to send, and
		// that is the difference between a consumer that waits on an idle
		// queue and one that throws. Sending an empty array is the helpful
		// kind of wrong: every test passes and the first quiet minute in
		// production does not.
		if !spec.OmitWhenEmpty || len(page.Records) > 0 {
			setPath(body, s.collectionName(resource, spec.Key), items)
		}

		if spec.CursorField != "" && page.NextCursor != "" {
			setPath(body, spec.CursorField, cursorValue(spec, page.NextCursor, nextURL))
		} else if spec.CursorField != "" && spec.CursorNull {
			setPath(body, spec.CursorField, nil)
		}

		// Salesforce's done is has_more with the sense reversed, and false is
		// its interesting value: a query that matched more rows than it
		// returned says done: false and expects the caller to follow
		// nextRecordsUrl. Sending true would tell a client its partial result
		// set was the whole thing.
		if spec.CompleteField != "" {
			setPath(body, spec.CompleteField, !page.HasMore)
		}

		// Only on the last page, which is the entire point. Google Calendar
		// gives you a page token or a sync token and never both, so code that
		// grabs whichever one it finds first stores the wrong one.
		if spec.FinalField != "" && !page.HasMore {
			setPath(body, spec.FinalField, finalToken(resource, page))
		}

		if spec.HasMoreField != "" {
			setPath(body, spec.HasMoreField, page.HasMore)
		}

		if spec.CountField != "" {
			setPath(body, spec.CountField, s.countValue(spec, countTotal(spec, page, limit)))
		}

		// How many are on this page, for the providers that send that beside
		// the total rather than instead of it.
		if spec.PageCountField != "" {
			setPath(body, spec.PageCountField, s.countValue(spec, len(page.Records)))
		}

		body = withFields(body, s.recipe.Responses.Success.Fields)

		return withFields(body, spec.Fields)
	default:
		body := map[string]any{
			"object":   "list",
			"data":     items,
			"has_more": page.HasMore,
		}

		if spec.URL {
			body["url"] = path
		}

		// A cursor field is opt in because sending one the provider does not
		// send is the more dangerous mistake: code written against it works
		// locally and fails in production, which is the exact failure this
		// project exists to prevent.
		if spec.CursorField != "" && page.NextCursor != "" {
			setPath(body, spec.CursorField, cursorValue(spec, page.NextCursor, nextURL))
		} else if spec.CursorField != "" && spec.CursorNull {
			setPath(body, spec.CursorField, nil)
		}

		body = withFields(body, s.recipe.Responses.Success.Fields)

		return withFields(body, spec.Fields)
	}
}

// finalToken derives the opaque token a provider hands back when a listing is
// exhausted.
//
// Stable for a given final page, so a suite can assert it rather than only
// matching a shape, and opaque so nothing is tempted to read structure out of
// it. The real ones carry no parseable structure either, and clients that
// guessed otherwise are the reason they are documented as opaque.
func finalToken(resource string, page store.Page) string {
	last := ""

	if n := len(page.Records); n > 0 {
		last, _ = page.Records[n-1]["id"].(string)
	}

	sum := sha256.Sum256([]byte(resource + ":" + last))

	return base64.RawURLEncoding.EncodeToString(sum[:15])
}

// collectionName resolves the key a wrapped list is nested under: the
// resource's declared collection, then a recipe-wide override, then the
// resource name itself.
func (s *Sandbox) collectionName(resource, override string) string {
	if spec, ok := s.recipe.Resources[resource]; ok && spec.Collection != "" {
		return spec.Collection
	}

	if override != "" {
		return override
	}

	return resource
}

// listName strips the array marker from a field name.
func listName(field string) string {
	return strings.TrimSuffix(field, "[]")
}

// numberOrString keeps an all-digit code a number, because Twilio's error codes
// are integers and a client comparing against 20404 must not be handed "20404".
func numberOrString(value string) any {
	if n, err := strconv.Atoi(value); err == nil {
		return n
	}

	return value
}

// writeRaw answers with a body exactly as given, past the identifier renaming
// and the resource envelope that writeRecord applies.
//
// A delete receipt keyed on its own name is not the resource, so presenting it
// as one would rename the very key the route just named.
func (s *Sandbox) writeRaw(w http.ResponseWriter, matched route, body map[string]any) int {
	status := matched.spec.Status
	if status == 0 {
		status = http.StatusOK
	}

	// The identifier is retyped the same way present would have retyped it.
	// A provider whose ids are JSON numbers sends a number here too, and a
	// receipt that disagrees with every other route about the type is its
	// own small lie.
	if spec, ok := s.recipe.Resources[matched.spec.Resource]; ok && spec.ID.Type == "number" {
		for key, value := range body {
			body[key] = numeric(value)
		}
	}

	if len(matched.spec.Fields) > 0 {
		body = withFields(body, matched.spec.Fields)
	}

	writeJSON(w, status, body)

	return status
}

// declaredOnly drops request fields the resource does not declare.
//
// A create used to store the decoded body as the record, so anything sent
// came back: posting {"totallyMadeUpField": "xyzzy"} to a refund answered with
// totallyMadeUpField in it. No provider does that, and the cost is not the
// stray key. It is that a conformance case asserting a value it sent on a
// create cannot fail, because the value comes back whether or not the Recipe
// declares the field. Adyen's refund carried the payment's name for its
// reference for exactly that long.
//
// A field declared type: map is free-form and keeps whatever it was sent,
// because some providers really do accept arbitrary keys. That is a Recipe
// saying so rather than the runtime assuming it.
//
// The identifier is taken under whichever name the resource publishes it. It
// used to be kept under that name as an ordinary key, which meant a create
// naming its own identifier -- Pinecone's POST /indexes sends {"name": ...},
// and the index name is the identifier -- stored no id at all, minted one, and
// then answered with two writers for a single key: the stored value under the
// wire name, and the minted identifier rendered under the same wire name.
//
// Which won depended on map iteration order, so the response was a coin flip.
// That is the worst failure this emulator can have: a fake that is
// non-deterministic makes a suite flaky, and the first place anybody looks is
// their own code.
func (s *Sandbox) declaredOnly(resource string, record store.Record) store.Record {
	spec, ok := s.recipe.Resources[resource]
	if !ok {
		return record
	}

	kept := store.Record{}

	for name, value := range record {
		if name == "id" {
			kept[name] = value

			continue
		}

		// The wire name of the identifier, stored as the identifier. Not when
		// the resource also declares a field of that name: then the Recipe has
		// said the two are different things, and the field wins because the
		// identifier still has "id" to arrive under.
		if name == spec.ID.Field && spec.ID.Field != "" && spec.ID.Field != "-" {
			if _, declared := spec.Fields[name]; !declared {
				if _, already := kept["id"]; !already {
					kept["id"] = value
				}

				continue
			}
		}

		if _, declared := spec.Fields[name]; declared {
			kept[name] = value
		}
	}

	return kept
}

// countTotal is the number the count field carries, which is not always how
// many records matched.
//
// Three quantities travel under the same name and a fixture small enough to
// fit on one page makes all three agree, so this is the one place the
// difference is decided rather than a property of the data.
func countTotal(spec recipe.ListResponse, page store.Page, limit int) int {
	switch spec.CountMeans {
	case "page":
		// The length of the page in front of you, which is what Shopware
		// reports unless a request asks it to count. Not a mistake on its
		// part: counting is a second query and it does not run one uninvited.
		return len(page.Records)
	case "lookahead":
		// A bounded count. The provider reads a few pages past this one and
		// reports what it found, so the number is real up to the window and
		// stops there. The extra row is the sentinel that says the window is
		// not the end.
		bound := limit*spec.CountLookahead + 1
		if page.Total < bound {
			return page.Total
		}

		return bound
	default:
		return page.Total
	}
}
