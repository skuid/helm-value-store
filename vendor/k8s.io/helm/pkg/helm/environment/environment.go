/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

/*Package environment describes the operating environment for Tiller.

Tiller's environment encapsulates all of the service dependencies Tiller has.
These dependencies are expressed as interfaces so that alternate implementations
(mocks, etc.) can be easily generated.
*/
package environment

import (
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/helm/helmpath"
)

const (
	// HomeEnvVar is the HELM_HOME environment variable key.
	HomeEnvVar = "HELM_HOME"
	// PluginEnvVar is the HELM_PLUGIN environment variable key.
	PluginEnvVar = "HELM_PLUGIN"
	// PluginDisableEnvVar is the HELM_NO_PLUGINS environment variable key.
	PluginDisableEnvVar = "HELM_NO_PLUGINS"
	// HostEnvVar is the HELM_HOST environment variable key.
	HostEnvVar = "HELM_HOST"
)

// DefaultHelmHome gets the configured HELM_HOME, or returns the default.
func DefaultHelmHome() string {
	if home := os.Getenv(HomeEnvVar); home != "" {
		return home
	}
	return filepath.Join(os.Getenv("HOME"), ".helm")
}

// DefaultHelmHost returns the configured HELM_HOST or an empty string.
func DefaultHelmHost() string {
	return os.Getenv(HostEnvVar)
}

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	// TillerHost is the host and port of Tiller.
	TillerHost string
	// TillerNamespace is the namespace in which Tiller runs.
	TillerNamespace string
	// Home is the local path to the Helm home directory.
	Home helmpath.Home
	// PluginDirs is the path to the plugin directories.
	PlugDirs string
	// Debug indicates whether or not Helm is running in Debug mode.
	Debug bool
}
