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
	"strconv"
	"time"

	"github.com/99nil/gopkg/signals"
	"golang.org/x/net/http2"
	"golang.org/x/sync/errgroup"
)

const DefaultPort = 9090

type Config struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type Server struct {
	*http.Server
}

func New(cfg *Config) *Server {
	if cfg == nil {
		cfg = new(Config)
	}

	port := DefaultPort
	if cfg.Port > 0 {
		port = cfg.Port
	}

	var addr string
	if len(cfg.Host) > 0 {
		addr += cfg.Host
	}
	addr += ":" + strconv.Itoa(port)

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
	return s.ListenAndServe()
}

func (s *Server) RunAndStop(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return signals.Exit(ctx) })
	eg.Go(func() error { return s.Run(ctx) })
	eg.Go(func() error { return WaitForShutdown(ctx, s.Server, 0) })
	return eg.Wait()
}

func (s *Server) RunTLS(ctx context.Context) error {
	return s.ListenAndServeTLS("", "")
}

func (s *Server) RunTLSAndStop(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return signals.Exit(ctx) })
	eg.Go(func() error { return s.RunTLS(ctx) })
	eg.Go(func() error { return WaitForShutdown(ctx, s.Server, 0) })
	return eg.Wait()
}

func WaitForShutdown(ctx context.Context, srv *http.Server, timeout time.Duration) error {
	<-ctx.Done()

	if timeout <= 0 {
		timeout = time.Second * 5
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}
	return nil
}
