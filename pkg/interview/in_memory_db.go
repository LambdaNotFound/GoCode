package interview

import (
	"sort"
	"strings"
)

// InMemoryDBV1 stores records as a two-level hash map: key -> field -> value.
// Point operations (Set/Get/Delete) are O(1); scans sort on demand rather than
// maintaining sorted order, keeping writes cheap at the cost of O(f log f) reads.
type InMemoryDBV1 struct {
	records map[string]map[string]string // key -> field -> value

}

func NewInMemoryDBV1() *InMemoryDBV1 {
	return &InMemoryDBV1{
		records: make(map[string]map[string]string),
	}
}

// Set upserts a field-value pair under key, lazily creating the inner map
// on first write to a key. O(1).
func (db *InMemoryDBV1) Set(key, field, value string) {
	if _, ok := db.records[key]; !ok {
		db.records[key] = make(map[string]string)
	}
	db.records[key][field] = value
}

// Get looks up key then field, returning nil at either miss so the caller can
// distinguish "absent" from an empty-string value. O(1).
func (db *InMemoryDBV1) Get(key, field string) *string {
	rec, ok := db.records[key]
	if !ok {
		return nil
	}
	val, ok := rec[field]
	if !ok {
		return nil
	}
	return &val
}

// Delete removes a single field under key, reporting whether it existed.
// The inner map is checked before deleting so a miss returns false rather
// than silently no-op'ing; an emptied record is left in place. O(1).
func (db *InMemoryDBV1) Delete(key, field string) bool {
	rec, ok := db.records[key]
	if !ok {
		return false
	}
	if _, ok := rec[field]; !ok {
		return false
	}
	delete(rec, field)
	return true
}

// Scan returns all fields of a record as "field(value)" strings in
// lexicographic field order — a prefix scan with the empty prefix, which
// matches every field. O(f log f) for f fields.
func (db *InMemoryDBV1) Scan(key string) []string {
	rec, ok := db.records[key]
	if !ok {
		return []string{}
	}
	return formatAndSort(rec, "")
}

// ScanByPrefix filters a record's fields by prefix, formatted and sorted the
// same way as Scan. Fields are unordered in the map, so this is a full pass
// over the record plus a sort: O(f log f).
func (db *InMemoryDBV1) ScanByPrefix(key, prefix string) []string {
	rec, ok := db.records[key]
	if !ok {
		return []string{}
	}
	return formatAndSort(rec, prefix)
}

// formatAndSort is the shared scan core: collect fields matching prefix as
// "field(value)", then sort. Sorting the formatted strings still orders by
// field name because "(" sorts below all alphanumerics and fields are unique.
func formatAndSort(rec map[string]string, prefix string) []string {
	var result []string
	for field, value := range rec {
		if strings.HasPrefix(field, prefix) {
			result = append(result, field+"("+value+")")
		}
	}
	sort.Strings(result) // lexicographic on field name, and since format is "field(value)", this works
	return result
}
