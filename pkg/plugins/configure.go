/*
Copyright 2021 The NitroCI Authors.

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
package plugins

import (
	"os"
	"strings"

	"github.com/nitroci/nitroci-core/pkg/core/configs"
	"github.com/nitroci/nitroci-core/pkg/core/contexts"
	"github.com/nitroci/nitroci-core/pkg/core/terminal"
	"github.com/nitroci/nitroci-core/pkg/core/net/http"
)

func OnConfigure(context *contexts.RuntimeContext, args []string, fields map[string]interface{}) {
	var domain string
	if fields["jfrog-domain"] != nil {
		domain = fields["jfrog-domain"].(string)
	}
	if len(domain) == 0 {
		domain, _ = terminal.PromptGlobalConfigKey(context.Cli.Profile, "Domain", false)
	}
	var username string
	if fields["jfrog-user"] != nil {
		username = fields["jfrog-user"].(string)
	}
	if len(username) == 0 {
		username, _ = terminal.PromptGlobalConfigKey(context.Cli.Profile, "Username", false)
	}
	var password string
	if fields["jfrog-pass"] != nil {
		password = fields["jfrog-pass"].(string)
	}
	if len(password) == 0 {
		password, _ = terminal.PromptGlobalConfigKey(context.Cli.Profile, "Password", true)
	}
	httpResult, err := http.HttpGetWithAuth("https://"+domain+".jfrog.io/"+domain+"/api/npm/auth", username, password)
	if err != nil || httpResult.StatusCode != 200 {
		errMessage := "Operation cannot be completed. Please verify the inputs."
		terminal.Print(&terminal.TerminalOutput{
			Messages:    []string{errMessage},
			MessageType: terminal.Error,
		})
		os.Exit(1)
	}
	for _, line := range strings.Split(strings.TrimSuffix(httpResult.ToString(), "\n"), "\n") {
		s := strings.Split(line, " = ")
		if s[0] == "_auth" {
			configs.SetGlobalConfigString(context.Cli.Profile, "jfrog_secret", s[1])
		} else if s[0] == "email" {
			configs.SetGlobalConfigString(context.Cli.Profile, "jfrog_username", s[1])
		}
	}
}
