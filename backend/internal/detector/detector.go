package detector

import (
	"fmt"
	"regexp"
	"time"

	"securecoda/internal/coda"
	"securecoda/internal/store"
)

var (
	reEmail    = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	reSSN      = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
	reCard     = regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`)
	rePassword = regexp.MustCompile(`(?i)password\s*[:=]\s*\S+`)
)

type Detector struct {
	coda  *coda.Client
	store *store.Store
}

func NewDetector(c *coda.Client, s *store.Store, logger interface{}) *Detector {
	return &Detector{coda: c, store: s}
}

type Alert struct {
	ID        string    `json:"id"`
	DocID     string    `json:"doc_id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"` // unused, public_share, sensitive_row, sensitive_page
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *Detector) RunCycle() error {
	docs, err := d.coda.ListDocs()
	if err != nil {
		return err
	}
	now := time.Now()
	for _, doc := range docs {
		id := doc["id"].(string)
		title := ""
		if t, ok := doc["title"].(string); ok {
			title = t
		}
		// Check last updated
		if updated, ok := doc["updatedAt"].(string); ok {
			t, _ := time.Parse(time.RFC3339, updated)
			if now.Sub(t) > 90*24*time.Hour {
				d.store.CreateAlert(&store.Alert{DocID: id, Title: title, Type: "unused", Detail: "Not modified for 90 days"})
			}
		}
		// Check published status
		if published, ok := doc["isPublished"].(bool); ok && published {
			d.store.CreateAlert(&store.Alert{DocID: id, Title: title, Type: "public_share", Detail: "Document is published (public)"})
		}
		// Check tables
		tables, _ := d.coda.ListTables(id)
		for _, t := range tables {
			tableID := t["id"].(string)
			rows, _ := d.coda.GetRows(id, tableID)
			for _, r := range rows {
				raw := ""
				if vals, ok := r["values"].(map[string]interface{}); ok {
					for _, v := range vals {
						raw += " " + toString(v)
					}
				} else {
					// fallback: stringify row
				}
				if d.detectSensitive(raw) {
					d.store.CreateAlert(&store.Alert{DocID: id, Title: title, Type: "sensitive_row", Detail: "Sensitive data detected in table"})
				}
			}
		}
		// Check pages: export and scan
		if pages, ok := doc["pages"].([]interface{}); ok {
			for _, p := range pages {
				if pageMap, ok := p.(map[string]interface{}); ok {
					if pid, ok := pageMap["id"].(string); ok {
						html, _ := d.coda.ExportPageHTML(id, pid)
						if d.detectSensitive(html) {
							d.store.CreateAlert(&store.Alert{DocID: id, Title: title, Type: "sensitive_page", Detail: "Sensitive content in page"})
						}
					}
				}
			}
		}
	}
	return nil
}

func (d *Detector) detectSensitive(s string) bool {
	if reEmail.MatchString(s) {
		return true
	}
	if reSSN.MatchString(s) {
		return true
	}
	if reCard.MatchString(s) {
		return true
	}
	if rePassword.MatchString(s) {
		return true
	}
	return false
}

func toString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return fmt.Sprintf("%v", x)
	default:
		return ""
	}
}
