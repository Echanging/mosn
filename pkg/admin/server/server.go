/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"fmt"
	"net/http"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	jsoniter "github.com/json-iterator/go"
	"sofastack.io/sofa-mosn/pkg/admin/store"
	"sofastack.io/sofa-mosn/pkg/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// apiHandleFuncStore stores the supported admin api
// can register more admin api
var apiHandleFuncStore map[string]func(http.ResponseWriter, *http.Request)

func RegisterAdminHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	apiHandleFuncStore[pattern] = handler
	log.StartLogger.Infof("[admin server] [register api] register a new api %s", pattern)
}

func init() {
	// default admin api
	apiHandleFuncStore = map[string]func(http.ResponseWriter, *http.Request){
		"/api/v1/config_dump":     configDump,
		"/api/v1/stats":           statsDump,
		"/api/v1/update_loglevel": updateLogLevel,
		"/api/v1/enable_log":      enableLogger,
		"/api/v1/disbale_log":     disableLogger,
		"/api/v1/states":          getState,
		"/":                       help,
	}
}

type Server struct {
	*http.Server
}

func (s *Server) Start(config Config) {
	var addr string
	if config != nil {
		// merge MOSNConfig into global context
		store.SetMOSNConfig(config)
		// get admin config
		adminConfig := config.GetAdmin()
		if adminConfig == nil {
			// no admin config, no admin start
			log.DefaultLogger.Warnf("no admin config, no admin api served")
			return
		}
		address := adminConfig.GetAddress()
		if xdsPort, ok := address.GetSocketAddress().GetPortSpecifier().(*core.SocketAddress_PortValue); ok {
			addr = fmt.Sprintf("%s:%d", address.GetSocketAddress().GetAddress(), xdsPort.PortValue)
		}
	}

	mux := http.NewServeMux()
	for pattern, handler := range apiHandleFuncStore {
		mux.HandleFunc(pattern, handler)
	}

	srv := &http.Server{Addr: addr, Handler: mux}
	store.AddService(srv, "Mosn Admin Server", nil, nil)
	s.Server = srv
}
