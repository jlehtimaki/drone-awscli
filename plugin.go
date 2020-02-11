package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
)

type (
	// Config holds input parameters for the plugin
	Config struct {
		Actions          []string
		Sensitive        bool
		RoleARN          string
	}

	AWSCli struct {
		Version         string
		Command         string
	}

	// Plugin represents the plugin instance to be executed
	Plugin struct {
		Config      Config
		AWSCli      AWSCli
	}
)

// Exec executes the plugin
func (p Plugin) Exec() error {
	// Install specified version of awscli
	if p.AWSCli.Version == "" {
		err := installAWSCli()

		if err != nil {
			return err
		}
	}

	if p.Config.RoleARN != "" {
		assumeRole(p.Config.RoleARN)
	}

	// Initialize commands
	var commands []*exec.Cmd

	// Print AWSCli version
	commands = append(commands, exec.Command(awsCliExe, "--version"))

	// Set user defined command
	commands = append(commands, exec.Command(awsCliExe, strings.Split(p.AWSCli.Command," ")...))

	// Run commands
	for _, c := range commands {
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if !p.Config.Sensitive {
			trace(c)
		}

		err := c.Run()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Failed to execute a command")
		}
		logrus.Debug("Command completed successfully")
	}

	return nil
}

func assumeRole(roleArn string) {
	client := sts.New(session.New())
	duration := time.Hour * 1
	stsProvider := &stscreds.AssumeRoleProvider{
		Client:          client,
		Duration:        duration,
		RoleARN:         roleArn,
		RoleSessionName: "drone",
	}

	value, err := credentials.NewCredentials(stsProvider).Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error assuming role!")
	}
	os.Setenv("AWS_ACCESS_KEY_ID", value.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", value.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", value.SessionToken)
}



func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}