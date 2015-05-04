// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Marc Berhault (marc@cockroachlabs.com)

package base

// Base context defaults.
const (
	defaultCerts  = "certs"
	defaultPort   = 8080
	defaultRegion = "aws:us-east-1"
)

// Context is the base context object.
type Context struct {
	// Certificates directory.
	Certs string
	// Port for cockroach nodes to listen on.
	Port int64
	// Region to run in.
	Region string
}

// NewContext returns a context with initialized values.
func NewContext() *Context {
	ctx := &Context{}
	ctx.InitDefaults()
	return ctx
}

// InitDefaults sets up the default values for a context.
func (ctx *Context) InitDefaults() {
	ctx.Certs = defaultCerts
	ctx.Port = defaultPort
	ctx.Region = defaultRegion
}
