package api

import (
	"encoding/json"
	"net/http"

	"securecoda/internal/coda"
	"securecoda/internal/detector"
	"securecoda/internal/store"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func RegisterRoutes(r *mux.Router, st *store.Store, det *detector.Detector, c *coda.Client, logger interface{}) {
	r.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		alerts, err := st.ListAlerts()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(alerts)
	}).Methods("GET")

	r.HandleFunc("/api/remediate/unpublish/{doc}", func(w http.ResponseWriter, r *http.Request) {
		doc := mux.Vars(r)["doc"]
		if err := c.UnpublishDoc(doc); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	}).Methods("POST")

	r.HandleFunc("/api/remediate/delete-doc/{doc}", func(w http.ResponseWriter, r *http.Request) {
		doc := mux.Vars(r)["doc"]
		// Placeholder: Coda deletion may require different API
		if err := c.UnpublishDoc(doc); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	}).Methods("POST")

	r.HandleFunc("/api/trigger-scan", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			if err := det.RunCycle(); err != nil {
				log.Error().Err(err).Msg("manual scan failed")
			}
		}()
		w.WriteHeader(202)
	}).Methods("POST")
}
