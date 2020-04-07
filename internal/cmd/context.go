/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package cmd

import (
	"context"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"io"

	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

// Context is used as a context for the commands
type Context struct {
	Ctx context.Context

	// Output should be used when printing in commands, instead of directly writing to stdout/stderr, to enable unit testing.
	Output io.Writer

	Client smclient.Client

	Verbose bool

	Configuration configuration.Configuration

	Parameters query.Parameters
}
