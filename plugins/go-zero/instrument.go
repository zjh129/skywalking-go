// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package gozero

import (
	"embed"
	"strings"

	"github.com/apache/skywalking-go/plugins/core/instrument"
)

//go:embed *
var fs embed.FS

type Instrument struct {
}

func NewInstrument() *Instrument {
	return &Instrument{}
}

func (i *Instrument) Name() string {
	return "gozero"
}

func (i *Instrument) BasePackage() string {
	return "github.com/zeromicro/go-zero"
}

func (i *Instrument) VersionChecker(version string) bool {
	return strings.HasPrefix(version, "v1.")
}

func (i *Instrument) Points() []*instrument.Point {
	return []*instrument.Point{
		{
			PackagePath: "core/stores/redis",
			At:          instrument.NewMethodEnhance("*Redis", "Do"),
			Interceptor: "RedisDoInterceptor",
		},
		{
			PackagePath: "core/stores/sqlx",
			At:          instrument.NewMethodEnhance("*SqlConn", "QueryRow"),
			Interceptor: "SqlQueryRowInterceptor",
		},
		{
			PackagePath: "core/stores/sqlx",
			At:          instrument.NewMethodEnhance("*SqlConn", "Exec"),
			Interceptor: "SqlExecInterceptor",
		},
	}
}

func (i *Instrument) FS() *embed.FS {
	return &fs
}
