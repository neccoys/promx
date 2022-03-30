package promx

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Proms struct {
	Name string
	Path string
	Port string
}

func NewServe(name, path, port string) *Proms {
	return &Proms{
		Name: name,
		Path: path,
		Port: port,
	}
}

func New() *Proms {
	return &Proms{
		Name: "prometheus_metrics",
		Path: "/metrics",
		Port: "9091",
	}
}

func (p *Proms) SetName(name string) *Proms {
	p.Name = name
	return p
}

func (p *Proms) SetPort(port string) *Proms {
	p.Port = port
	return p
}

func (p *Proms) SetPath(path string) *Proms {
	p.Path = path
	return p
}

func (p *Proms) Start() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: p.Name,
	})

	prometheus.MustRegister(counter)

	var (
		done = make(chan bool, 1)
		quit = make(chan os.Signal, 1)
	)

	router := http.NewServeMux()
	router.Handle(p.Path, promhttp.Handler())
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", p.Port),
		Handler: router,
	}

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go gracefulShutdown(srv, quit, done)

	log.Println("[HTTP] Start...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panicln("[HTTP] Error:", err)
	}

	<-done
	fmt.Println("[HTTP] Shutdown!")

}

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("[HTTP] Shutdown Error:", err)
	}
	close(done)
}
