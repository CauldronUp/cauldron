package runtime

import "testing"

// The action is what is left of an event name once the resource is taken out
// of it, which is the only definition that works across providers: Asana puts
// the resource first and Pipedrive puts it last, and both want the other half.
//
// The case-insensitive match is here because a provider's event names and a
// Recipe's resource names do not have to agree on it. Xero declares
// INVOICE.CREATE against a resource called invoice and Documenso declares
// DOCUMENT_CREATED against document. Matching exactly returned the empty
// string for both, which is the failure this substitution exists to avoid: a
// present field with nothing in it, which reads as a provider that sent an
// empty action rather than an emulator that could not work one out.
func TestTheActionIsTheEventWithoutTheResource(t *testing.T) {
	for _, c := range []struct {
		event, resource, want string
	}{
		// Resource first, resource last, and an underscore instead of a dot.
		{"task.added", "task", "added"},
		{"added.deal", "deal", "added"},
		{"folder_added", "folder", "added"},
		{"Customer.Create", "Customer", "Create"},

		// Case need not agree.
		{"INVOICE.CREATE", "invoice", "CREATE"},
		{"DOCUMENT_CREATED", "document", "CREATED"},

		// Nothing to say rather than a guess.
		{"crawl.completed", "", ""},
		{"page.created", "database", ""},
	} {
		if got := actionOf(c.event, c.resource); got != c.want {
			t.Errorf("actionOf(%q, %q) = %q, want %q", c.event, c.resource, got, c.want)
		}
	}
}
