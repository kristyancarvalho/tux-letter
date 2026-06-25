package article

func Dedupe(items []Article) []Article {
	seenURL := make(map[string]bool)
	seenHash := make(map[string]bool)
	var out []Article

	for _, a := range items {
		key := CanonicalKey(a.URL)
		hash := a.Hash
		if hash == "" {
			hash = Hash(a)
		}
		if key == "" && hash == "" {
			continue
		}
		if key != "" && seenURL[key] {
			continue
		}
		if hash != "" && seenHash[hash] {
			continue
		}
		if key != "" {
			seenURL[key] = true
		}
		if hash != "" {
			seenHash[hash] = true
		}
		out = append(out, a)
	}
	return out
}

func NormalizeAll(items []Article) []Article {
	out := make([]Article, 0, len(items))
	for _, a := range items {
		n := Normalize(a)
		if n.IsZero() {
			continue
		}
		out = append(out, n)
	}
	return out
}
