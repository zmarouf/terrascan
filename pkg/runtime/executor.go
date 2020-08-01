/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package runtime

import (
	"go.uber.org/zap"

	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
)

// Executor object
type Executor struct {
	filePath    string
	dirPath     string
	cloudType   string
	iacType     string
	iacVersion  string
	iacProvider iacProvider.IacProvider
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion, cloudType, filePath, dirPath string) (e *Executor, err error) {
	e = &Executor{
		filePath:   filePath,
		dirPath:    dirPath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
	}

	// initialized executor
	if err = e.Init(); err != nil {
		return e, err
	}

	return e, nil
}

// Init validates input and initializes iac and cloud providers
func (e *Executor) Init() error {

	// validate inputs
	err := e.ValidateInputs()
	if err != nil {
		return err
	}

	// create new IacProvider
	e.iacProvider, err = iacProvider.NewIacProvider(e.iacType, e.iacVersion)
	if err != nil {
		zap.S().Errorf("failed to create a new IacProvider for iacType '%s'. error: '%s'", e.iacType, err)
		return err
	}

	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute() (normalized interface{}, err error) {

	if e.dirPath != "" {
		normalized, err = e.iacProvider.LoadIacDir(e.dirPath)
	} else {
		// create config from IaC
		normalized, err = e.iacProvider.LoadIacFile(e.filePath)
	}
	if err != nil {
		return normalized, err
	}

	// write output

	// successful
	return normalized, nil
}