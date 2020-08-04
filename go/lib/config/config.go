// Copyright 2019 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config provides an unified pattern for configuration structs.
//
// Usage
//
// Every configuration struct should implement the Config interface. There
// are four parts to a configuration: Initialization, validation,
// sample generation and configuration from a template.
//
// Initialization
//
// A config struct is initialized by calling InitDefaults. This recursively
// initializes all uninitialized fields. Fields that should not be
// initialized to default must be set before calling InitDefaults.
//
// Validation
//
// A config struct is validated by calling Validate. This recursively
// validates all fields.
//
// Sample Generation
//
// A config struct can be used to generate a commented sample toml config
// by calling Sample. Unit tests guarantee the consistency between
// implementation and the generated sample. To this end, each config struct
// has to provide a composable unit test to check that the sample is
// parsable and consistent with the default values. See lib/envtest for an
// example.
//
// Warning: The method Sample is allowed to panic if an error occurs during
// sample generation.
//
// Configuration
// The sample configuration is interactively edited and validated.
//
package config

import (
	"fmt"
	"github.com/scionproto/scion/go/lib/common"
	"io"
	"os"
)

const ID = "id"

// Config is the interface that config structs should implement to allow for
// streamlined initialization, configuration, validation and sample generation.
type Config interface {
	Sampler
	Validator
	Defaulter
	// Configure creates a customized config and writes it to dst. Ctx provides
	// additional information.
	Configure(dst io.Writer, path Path, ctx CtxMap)
}

// Validator defines the validation part of Config.
type Validator interface {
	// Validate recursively checks that all fields contain valid values.
	Validate() error
}

// Defaulter defines the initialization part of Config.
type Defaulter interface {
	// InitDefaults recursively initializes the default values of all
	// uninitialized fields.
	InitDefaults()
}

// Sampler defines the sample generation part of Config.
type Sampler interface {
	// Sample creates a sample config and writes it to dst. Ctx provides
	// additional information. Sample is allowed to panic if an error
	// occurs.
	Sample(dst io.Writer, path Path, ctx CtxMap)
	// ConfigName returns the name of the config block. This forces
	// consistency between samples for different services for the same
	// config block.
	ConfigName() string
}

// Path is the header of a config block possibly consisting of multiple parts.
type Path []string

// Extend creates a copy of the path with string s appended.
func (p Path) Extend(s string) Path {
	c := append(Path(nil), p...)
	return append(c, s)
}

// NoValidator implements a Validator that never fails to validate. It can
// be embedded in config structs that do not need to validate.
type NoValidator struct{}

// Validate always returns nil.
func (NoValidator) Validate() error {
	return nil
}

// NoDefaulter implements a Defaulter that does a no-op on InitDefaults.
// It can be embedded in config structs that do not have any defaults.
type NoDefaulter struct{}

// InitDefaults is a no-op.
func (NoDefaulter) InitDefaults() {}

// StringSampler implements a Sampler that writes string Text and provides
// Name as ConfigName.
type StringSampler struct {
	// Text the sample string.
	Text string
	// Name the config name.
	Name string
}

// Sample writes the text to dst.
func (s StringSampler) Sample(dst io.Writer, _ Path, _ CtxMap) {
	WriteString(dst, s.Text)
}

// ConfigName returns the name.
func (s StringSampler) ConfigName() string {
	return s.Name
}

// NoConfigurator implements a Configurator that only applies the default values.
// It can be embedded in config structs that do not need any non-default values.
type NoConfigurator struct{
	Defaulter
	Validator
}

// Configure with defaults.
func (c NoConfigurator) Configure(dst io.Writer, _ Path, _ CtxMap) {
	c.InitDefaults()
	if err := c.Validate(); err != nil {
		return
	}
	WriteConfiguration(dst, nil, nil,  c)
}

// ValidateAll validates all validators. The first error encountered is returned.
func ValidateAll(validators ...Validator) error {
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			return common.NewBasicError("Unable to validate", err, "type", fmt.Sprintf("%T", v))
		}
	}
	return nil
}

// InitAll initializes all defaulters.
func InitAll(defaulters ...Defaulter) {
	for _, d := range defaulters {
		d.InitDefaults()
	}
}

// ConfigureAll configures and validates all configurators. The first error encountered is returned.
func ConfigureAll(configurators ...Config) error {
	for _, c := range configurators {
		if err := c.Validate(); err != nil {
			return common.NewBasicError("Unable to validate", err, "type", fmt.Sprintf("%T", c))
		}
		c.Sample(os.Stdout, nil, nil)
	}
	return nil
}
