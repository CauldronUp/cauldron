// What a Recipe may say about the provider's own description of itself, and
// about the versions of it that came before.

package recipe

import "regexp"

// dateOnly is the form every dated field in this project uses, and the reason
// it is checked is that "last Tuesday" is a date to a person and nothing at
// all to a scan comparing two of them.
var dateOnly = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// validateUpstream checks the parts of upstream that only make sense together.
//
// The rules are all of one kind: a Recipe may say nothing about a subject, but
// it may not say half of it. Half a declaration reads as a complete one to
// anybody skimming, and reads as nothing to the code that acts on it, which is
// the worst combination available.
func (r *Recipe) validateUpstream(add func(string, ...any)) {
	up := r.Upstream

	if up.SpecHash != "" && up.Spec == "" {
		add("upstream.spec_hash records a fingerprint of a description upstream.spec does not name, so nothing can ever recompute it")
	}

	if up.Spec != "" && up.SpecHash != "" && up.SpecSeen == "" {
		add("upstream.spec_seen is required beside a fingerprint: undated evidence cannot tell a Recipe checked yesterday from one checked three years ago")
	}

	if up.SpecSeen != "" && !dateOnly.MatchString(up.SpecSeen) {
		add("upstream.spec_seen %q must be a date like 2026-08-30, because a scan comparing two of them cannot read prose", up.SpecSeen)
	}

	if len(up.Supersedes) == 0 {
		return
	}

	if up.Provider == "" {
		// Without it the two versions are unrelated, which is the state the
		// format was already in and the reason supersedes exists.
		add("upstream.provider is required beside upstream.supersedes, or the version this replaced belongs to nobody in particular")
	}

	for i, was := range up.Supersedes {
		if was.Version == "" {
			add("upstream.supersedes[%d] names no version", i)
		}

		if was.Host == "" {
			// The host is the whole operational value. A client is told to be
			// a v2 client by what it talks to, and detection has nothing else
			// to go on.
			add("upstream.supersedes[%d] names no host, and the host is the only thing that tells one version's clients from another's", i)
		}

		if was.Version != "" && was.Version == up.API {
			add("upstream.supersedes[%d] is %q, which is the version this Recipe describes", i, was.Version)
		}
	}
}
