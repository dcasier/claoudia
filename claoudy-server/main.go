package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dcasier/claoudy"
	"os"
	"log"
	"crypto/tls"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
}

func main() {

	finish := make(chan bool)

	/*
	PutActivity("sport")
	PutActivity("musique")
	GetActivities()
	s, _ := GetActivity("sport")
	e := new(metamodel.Event)
	e.Name = "match du 24"
	e.Desc = "Description du match"
	PutEvent("sport", e)
	m := new(metamodel.Event)
	m.Name ="concert du 25"
	m.Desc = "Description du match"
	PutEvent("musique", m)
	file, _ := os.Open("minio.exe")
	PutMedia("match du 24", "minio", file)
	AddGrantToEvent("match du 24", "user1", "READ")
	AddGrantToEvent("match du 24", "user2", "FULLCONTROL")
	AddGrantToEvent("match du 24", "user3", "READ")
	AddGrantToEvent("match du 24", "user4", "FULLCONTROL")
	AddGrantToEvent("match du 24", "user5", "READ")
	AddGrantToEvent("match du 24", "user6", "READ")
	DelGrantToEvent("match du 24", "user5")
	AddGrantToEvent("match du 24", "user7", "READ")
	DelGrantToEvent("match du 24", "user7")

	fmt.Println(s)
	*/
	r := mux.NewRouter()

	r.HandleFunc("/", Gz).Methods("GET")
	r.HandleFunc("/api/v0.1/", GetCacheHandler).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v0.1/", PutCacheHandler).Methods("PUT")
	r.HandleFunc("/api/v0.1/config", GetConfigHandler).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/v0.1/spheres", GetActivitiesHandler).Methods("OPTIONS", "GET")

	r.HandleFunc("/api/v0.1/register", RegisterHandler).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/v0.1/login",  LoginHandler).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/v0.1/logout", LogoutHandler).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/v0.1/register", RegisterHandler).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/v0.1/setpassword", SetPasswordHandler).Methods("OPTIONS", "POST")	

    // Routes consist of a path and a handler function.
    r.HandleFunc("/api/v0.1/{activity}", GetActivityHandler).Methods("OPTIONS", "GET")
    r.HandleFunc("/api/v0.1/{activity}", PutActivityHandler).Methods("PUT")
	r.HandleFunc("/api/v0.1/{activity}", PostActivityHandler).Methods("POST")
    r.HandleFunc("/api/v0.1/{activity}", YourHandler).Methods("HEAD")
    r.HandleFunc("/api/v0.1/{activity}", DeleteActivityHandler).Methods("DELETE")

	r.HandleFunc("/api/v0.1/{activity}/{event}", PostEventHandler).Methods("OPTIONS", "POST")
    r.HandleFunc("/api/v0.1/{activity}/{event}", ListMediaHandler).Methods("GET")
    r.HandleFunc("/api/v0.1/{activity}/{event}", PutEventHandler).Methods("PUT")
    r.HandleFunc("/api/v0.1/{activity}/{event}", YourHandler).Methods("HEAD")
    r.HandleFunc("/api/v0.1/{activity}/{event}", YourHandler).Methods("DELETE")	
    r.HandleFunc("/api/v0.1/{activity}/{event}/{media}", GetMediaHandler).Methods("OPTIONS", "GET")
    r.HandleFunc("/api/v0.1/{activity}/{event}/{media}", PutMediaHandler).Methods("PUT")
    r.HandleFunc("/api/v0.1/{activity}/{event}/{media}", PostMediaHandler).Methods("OPTIONS", "POST")
    r.HandleFunc("/api/v0.1/{activity}/{event}/{media}", YourHandler).Methods("HEAD")
    r.HandleFunc("/api/v0.1/{activity}/{event}/{media}", YourHandler).Methods("DELETE")	

    cfg := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
    }
	srv2 := &http.Server{
        Addr:         ":9443",
        Handler:      r,
        TLSConfig:    common.MustGetTlsConfiguration(),
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }
    go func() {
		log.Fatal(srv2.ListenAndServeTLS(os.Args[2], os.Args[1]))
	}()
	srv := &http.Server{
        Addr:         ":8443",
        Handler:      r,
        TLSConfig:    cfg,
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }
    go func() {
		log.Fatal(srv.ListenAndServeTLS(os.Args[2], os.Args[1]))
	}()
    // Bind to a port and pass our router in
    //log.Fatal(http.ListenAndServe(":8000", r))

	<- finish
}
