/*
Copyright © 2021 Kaku Li <1154584512@qq.com>

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

package app

import (
	"context"
	"errors"
	"flag"

	"github.com/likakuli/generic-project-template/cmd/apiserver/app/options"
	"github.com/likakuli/generic-project-template/pkg/config"
	"github.com/likakuli/generic-project-template/pkg/server"
	"github.com/likakuli/generic-project-template/pkg/util/signals"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	// BuildDate date string of when build was performed filled in by -X compile flag
	BuildDate string

	// LatestCommit date string of when build was performed filled in by -X compile flag
	LatestCommit string

	// BuildNumber date string of when build was performed filled in by -X compile flag
	BuildNumber string

	// BuiltOnIP date string of when build was performed filled in by -X compile flag
	BuiltOnIP string

	// BuiltOnOs date string of when build was performed filled in by -X compile flag
	BuiltOnOs string

	// RuntimeVer date string of when build was performed filled in by -X compile flag
	RuntimeVer string
)

func NewServerCommand() *cobra.Command {
	opts := options.NewOptions()
	cmd := &cobra.Command{
		Use:                "generic-project-template",
		Long:               "A generic restful controllers and command line project template.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			defer glog.Flush()

			cmd.Flags().AddGoFlagSet(flag.CommandLine)
			if err := cmd.Flags().Parse(args); err != nil {
				return err
			}

			help, err := cmd.Flags().GetBool("help")
			if err != nil {
				return errors.New(`"help" flag declared as non-bool. Please correct your code`)
			}
			if help {
				return cmd.Help()
			}

			version, err := cmd.Flags().GetBool("version")
			if err != nil {
				return errors.New(`"version" flag declared as non-bool. Please correct your code`)
			}
			if version {
				cmd.Printf(`Application build information
  Build date      : %s
  Build number    : %s
  Git commit      : %s
  Runtime version : %s
  Built on OS     : %s
`, BuildDate, BuildNumber, LatestCommit, RuntimeVer, BuiltOnOs)
				return nil
			}

			options.PrintFlags(cmd.Flags())

			if err = opts.Validate(); err != nil {
				return err
			}

			cfg, err := opts.Complete()
			if err != nil {
				return err
			}

			ctx := signals.SetupSignalContext()
			if err = Run(ctx, cfg); err != nil {
				return err
			}

			return nil
		},
	}

	opts.AddFlags(cmd.Flags())
	return cmd
}

func Run(ctx context.Context, cfg *config.Config) error {
	glog.Infof("Starting server.")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server.NewServer(cfg).Run(ctx)

	return nil
}
