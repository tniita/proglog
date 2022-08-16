package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpServ := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpServ.handleProcuce).Methods("POST")
	r.HandleFunc("/", httpServ.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProcuceRequest struct {
	Record Record `json:"record"`
}

type ProcuceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProcuce(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ProcuceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProcuceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
