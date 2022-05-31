// Copyright © 2021 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/net/http2"
)

const DefaultPort = 9090

type Config struct {
	Host   string `json:"host" yaml:"host"`
	Port   int    `json:"port" yaml:"port"`
	Secret string `json:"secret" yaml:"secret"`
}

type Server struct {
	*http.Server
}

func New(cfg *Config) *Server {
	port := DefaultPort
	if cfg.Port > 0 {
		port = cfg.Port
	}
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{
		Addr:           addr,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &Server{Server: server}
}

func NewHTTP2(cfg *Config, conf *http2.Server) (*Server, error) {
	srv := New(cfg)
	if err := http2.ConfigureServer(srv.Server, conf); err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *Server) Run(ctx context.Context) error {
	go s.ShutdownGraceful(ctx)
	return s.ListenAndServe()
}

func (s *Server) RunTLS(ctx context.Context) error {
	go s.ShutdownGraceful(ctx)
	return s.ListenAndServeTLS("", "")
}

func (s *Server) ShutdownGraceful(ctx context.Context) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case <-ctx.Done():
	case <-ch:
		fmt.Println("Server shutdown.")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			fmt.Println("Server shutdown failed: ", err)
		}
	}
}
