package main

import "net/http"

func (apic *apiConfig) resetCounterHits(w http.ResponseWriter, r *http.Request) {
	apic.fileserverHits = 0
}
