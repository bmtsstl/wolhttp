package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/bmtsstl/wolhttp/wol"
)

type Config struct {
	HTTPAddr string             `json:"http_addr"`
	Target   map[string]*Target `json:"target"`
}

func LoadConfig(filename string) (cfg *Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := f.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	cfg = &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	return cfg, err
}

type Target struct {
	Network      string `json:"network"`
	LocalAddr    string `json:"local_addr"`
	RemoteAddr   string `json:"remote_addr"`
	HardwareAddr string `json:"hardware_addr"`
}

func (e *Target) Request() (*wol.Request, error) {
	req := &wol.Request{
		Network:   "udp",
		LocalAddr: nil,
		RemoteAddr: &net.UDPAddr{
			IP:   net.IPv4bcast,
			Port: 9,
		},
	}
	var err error

	if e.Network != "" {
		req.Network = e.Network
	}
	if e.LocalAddr != "" {
		req.LocalAddr, err = net.ResolveUDPAddr(req.Network, e.LocalAddr)
		if err != nil {
			return nil, err
		}
	}
	if e.RemoteAddr != "" {
		req.RemoteAddr, err = net.ResolveUDPAddr(req.Network, e.RemoteAddr)
		if err != nil {
			return nil, err
		}
	}
	req.HardwareAddr, err = net.ParseMAC(e.HardwareAddr)
	if err != nil {
		return nil, err
	}
	return req, nil
}

type Handler struct {
	Target   map[string]*Target
	ErrorLog *log.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("target")
	if key == "" {
		http.Error(w, "target not specified", http.StatusBadRequest)
		return
	}
	tgt, ok := h.Target[key]
	if !ok {
		http.Error(w, "undefined target", http.StatusNotFound)
		return
	}
	wolreq, err := tgt.Request()
	if err != nil {
		h.logf("wolhttp: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	err = wol.WakeOnLANRequest(wolreq)
	if err != nil {
		h.logf("wolhttp: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

func (h *Handler) logf(format string, args ...interface{}) {
	if h.ErrorLog != nil {
		h.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func main() {
	var cfgpath string
	flag.StringVar(&cfgpath, "c", "config.json", "configuration file")
	flag.Parse()

	cfg, err := LoadConfig(cfgpath)
	if err != nil {
		log.Fatalf("wolhttp: %v", err)
		return
	}

	h := &Handler{Target: cfg.Target}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		h.ServeHTTP(w, req)
	})
	srv := http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}
	err = srv.ListenAndServe()
	log.Fatalf("wolhttp: %v", err)
}
