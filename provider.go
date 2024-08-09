package main

import (
	"context"
	"github.com/libdns/libdns"
)

type Entry struct {
	record libdns.Record
	zone   string
}

type Provider struct {
	entries []Entry
}

func (provider *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	for _, rec := range recs {
		// NOTE: Currently, only TXT records are handled, I'm not sure if we need other kind of records.
		if rec.Type == "TXT" {
			provider.entries = append(provider.entries, Entry{record: rec, zone: zone})
		}
	}

	return recs, nil
}

func (provider *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	var entries []Entry

	for _, entry := range provider.entries {
		// Keep the entry if the zone is different.
		if entry.zone != zone {
			entries = append(entries, entry)
			continue
		}

		for _, rec := range recs {
			// Keep the entry if none of the record to remove matches.
			if entry.record != rec {
				entries = append(entries, entry)
			}
		}
	}

	provider.entries = entries
	return recs, nil
}
