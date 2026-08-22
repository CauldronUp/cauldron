// Emission: the webhook a provider sends after something changes.

package runtime

import (
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
	_, _ = s.webhooks.Emit(event, s.present(resource, record, ""))
}
