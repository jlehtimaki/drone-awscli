package main

import (
  "os"

  "github.com/sirupsen/logrus"
  "github.com/urfave/cli"
)

var revision string // build number set at compile-time

func main() {
  app := cli.NewApp()
  app.Name = "awscli plugin"
  app.Usage = "awscli plugin"
  app.Action = run
  app.Version = revision
  app.Flags = []cli.Flag{

    //
    // plugin args
    //

    cli.StringFlag{
      Name:   "assume_role",
      Usage:  "A role to assume before running the awscli commands",
      EnvVar: "PLUGIN_ASSUME_ROLE",
    },
    cli.StringFlag{
      Name:   "awscli_version",
      Usage:  "AWSCli version number",
      EnvVar: "PLUGIN_AWSCLI_VERSION",
    },
    cli.StringSliceFlag{
      Name:    "awscli_commands",
      Usage:   "AWSCli commands to be run",
      EnvVar:  "PLUGIN_COMMANDS",
    },
    cli.StringFlag{
      Name:     "shell",
      Usage:    "Run awscli in shell",
      EnvVar:   "PLUGIN_SHELL",
      Value:    "false",
    },
  }

  if err := app.Run(os.Args); err != nil {
    logrus.Fatal(err)
  }
}

func run(c *cli.Context) error {
  plugin := Plugin{
    Config: Config{
      RoleARN:          c.String("assume_role"),
      Shell:            c.Bool("shell"),
    },
    AWSCli: AWSCli{
      Version:          c.String("awscli_version"),
      Commands:          c.StringSlice("awscli_commands"),
    },
  }

  return plugin.Exec()
}
