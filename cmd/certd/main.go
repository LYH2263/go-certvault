package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/certvault"
	"example.com/certvault/internal/api"
)

func main() {
	addr := flag.String("addr", ":8091", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	persist := flag.String("persist", "", "证书快照 JSON 路径（可选）")
	flag.Parse()

	opts := []certvault.Option{}
	if *persist != "" {
		opts = append(opts, certvault.WithPersistPath(*persist))
	}
	v := certvault.New(opts...)
	defer v.Close()
	if *persist != "" {
		_ = v.LoadPersist()
	}

	srv := api.New(v, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("certd 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
