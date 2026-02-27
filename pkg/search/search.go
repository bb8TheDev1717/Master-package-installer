// Package search provides combined pacman + AUR search with fuzzy ranking.
package search

import (
	"masterr/pkg/aur"
	"masterr/pkg/pacman"
	"masterr/pkg/util"
	"sort"
	"strings"
	"sync"
)

// Result wraps a package from either pacman or AUR.
type Result struct {
	Package pacman.Package
}

// Search queries pacman and AUR in parallel, deduplicates results,
// and returns them ranked by fuzzy relevance to the query.
func Search(query string) ([]Result, error) {
	var repoPkgs, aurPkgs []pacman.Package
	var repoErr, aurErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		repoPkgs, repoErr = pacman.Search(query)
	}()
	go func() {
		defer wg.Done()
		aurPkgs, aurErr = aur.Search(query)
	}()
	wg.Wait()

	if repoErr != nil && aurErr != nil {
		return nil, repoErr
	}

	// Deduplicate: repo packages take priority over AUR
	seen := map[string]bool{}
	var results []Result
	for _, p := range repoPkgs {
		if !seen[p.Name] {
			seen[p.Name] = true
			results = append(results, Result{Package: p})
		}
	}
	for _, p := range aurPkgs {
		if !seen[p.Name] {
			seen[p.Name] = true
			results = append(results, Result{Package: p})
		}
	}

	return fuzzyRank(query, results), nil
}

// fuzzyRank scores and sorts results by relevance to the query.
// Scoring: exact match = 1000, prefix = 500, contains = 200,
// fallback to Levenshtein distance (subtracted from 200).
func fuzzyRank(query string, results []Result) []Result {
	type scored struct {
		r     Result
		score int
	}
	q := strings.ToLower(query)
	items := make([]scored, len(results))
	for i, r := range results {
		name := strings.ToLower(r.Package.Name)
		var s int
		switch {
		case name == q:
			s = 1000
		case strings.HasPrefix(name, q):
			s = 500
		case strings.Contains(name, q):
			s = 200
		default:
			s = 200 - util.Levenshtein(q, name)
		}
		items[i] = scored{r, s}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].score > items[j].score
	})

	out := make([]Result, len(items))
	for i, s := range items {
		out[i] = s.r
	}
	return out
}
