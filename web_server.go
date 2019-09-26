package main

import (
	"encoding/json"
	"fmt"
	"github.com/amaxlab/go-lib/log"
	"github.com/go-chi/chi"
	"net/http"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type PostReleStateResponse struct {
	State string `json:"state"`
}

type GetReleStateResponse struct {
	State string `json:"state"`
}

type ClientStateResponse struct {
	Connected bool `json:"connected"`
}

type RouteHandler struct {
	relayClient KC868Client
}

type WebServer struct {
	port         int
	routeHandler RouteHandler
}

func NewWebServer(port int, relayClient *KC868Client) *WebServer {
	return &WebServer{port: port, routeHandler: RouteHandler{relayClient: *relayClient}}
}

func (s *WebServer) start() error {
	router := chi.NewRouter()

	router.Get("/", s.routeHandler.HomePage)
	router.Get("/relay", s.routeHandler.GetRelays)
	router.Get("/relay/{id}", s.routeHandler.GetRelayById)
	router.Get("/relay/{id}/state", s.routeHandler.GetRelayState)
	router.Post("/relay/{id}/state", s.routeHandler.SetRelayState)
	router.Get("/healthCheck", s.routeHandler.HealthCheck)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)
}

func (r *RouteHandler) HealthCheck(w http.ResponseWriter, req *http.Request) {
	r.JsonResponse(w, HealthCheckResponse{Status: "OK"}, 200)
}

func (r *RouteHandler) HomePage(w http.ResponseWriter, req *http.Request) {
	r.JsonResponse(w, ClientStateResponse{Connected: r.relayClient.Connected}, 200)
}

func (r *RouteHandler) GetRelays(w http.ResponseWriter, req *http.Request) {
	r.JsonResponse(w, r.relayClient.Relays, 200)
}

func (r *RouteHandler) GetRelayById(w http.ResponseWriter, req *http.Request) {
	relay := r.relayClient.Relays[chi.URLParam(req, "id")]
	if relay == nil {
		http.Error(w, fmt.Sprintf("Relay with id: %s not found", chi.URLParam(req, "id")), http.StatusNotFound)
		return
	}

	r.JsonResponse(w, relay, 200)
}

func (r *RouteHandler) SetRelayState(w http.ResponseWriter, req *http.Request) {
	relay := r.relayClient.Relays[chi.URLParam(req, "id")]
	action := make([]byte, req.ContentLength)
	_, _ = req.Body.Read(action)
	if relay == nil {
		http.Error(w, fmt.Sprintf("Relay with id: %s not found", chi.URLParam(req, "id")), http.StatusNotFound)
		return
	}

	newState := false
	if string(action) == "on" {
		newState = true
	}

	r.relayClient.ChangeRelayState(chi.URLParam(req, "id"), newState)
	r.JsonResponse(w, PostReleStateResponse{State: string(action)}, 201)
}

func (r *RouteHandler) GetRelayState(w http.ResponseWriter, req *http.Request) {
	relay := r.relayClient.Relays[chi.URLParam(req, "id")]
	state := "off"
	if relay.Value  {
		state = "on"
	}
	r.JsonResponse(w, GetReleStateResponse{State: state}, 200)
}

func (r *RouteHandler) JsonResponse(w http.ResponseWriter, data interface{}, c int) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Error.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/j")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", j)
}
