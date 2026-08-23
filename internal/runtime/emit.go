// Emission: the webhook a provider sends after something changes.

package runtime

import (
	"fmt"

	"github.com/CauldronUp/cauldron/internal/recipe"
	"github.com/CauldronUp/cauldron/internal/store"
)

// emitFor sends the conventional lifecycle event for a resource, when the
// Recipe declares it. Providers emit customer.created after a customer is
// created; a fake that stays silent teaches an application the wrong lesson.
func (s *Sandbox) emitFor(resource, action, named string, record store.Record) {
	event := named
	if event == "" {
		event = resource + "." + action
	}

	if !s.webhooks.Known(event) {
		return
	}

	// Shaped, exactly as an HTTP response is. Without this the payload carried
	// the store's own field names, so the same record from the same sandbox at
	// the same instant had two shapes depending on how you looked at it: the
	// response nested amount under amount_money as the Recipe declares and the
	// webhook did not. A handler written against the emulator read
	// event.data.object.amount and was entirely green here and entirely wrong
	// against the provider.
	//
	// present rather than resourceBody, because the response envelope belongs
	// to the response. A webhook has its own envelope and the record sits
	// inside it unwrapped.
	_, _, _ = s.webhooks.Emit(event, s.present(resource, record, ""))
}

// emitChanges sends the events a route owes only when the field each names
// actually moved.
//
// The comparison is before-and-after rather than "did the request mention the
// field", because a write that sets status to the value it already held is not
// a status change and the provider does not announce one.
func (s *Sandbox) emitChanges(spec recipe.Route, before, after store.Record) {
	for _, conditional := range spec.EmitsWhen {
		if same(before[conditional.Field], after[conditional.Field]) {
			continue
		}

		if !s.webhooks.Known(conditional.Event) {
			continue
		}

		_, _, _ = s.webhooks.Emit(conditional.Event, s.present(spec.Resource, after, ""))
	}
}

// same compares two decoded JSON values for equality.
//
// Formatting them is enough for the scalars these events key off -- a status,
// a priority, a stage -- and avoids depending on the numeric type JSON
// decoding happened to choose. It is not a general deep equality and does not
// claim to be: EmitsWhen rejects a field that is not scalar, so this is never
// asked to compare a map whose key order would make it unreliable.
func same(before, after any) bool {
	return fmt.Sprintf("%v", before) == fmt.Sprintf("%v", after)
}
