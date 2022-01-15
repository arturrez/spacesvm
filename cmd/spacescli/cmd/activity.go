// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ava-labs/spacesvm/client"
)

var activityCmd = &cobra.Command{
	Use:   "activity [options]",
	Short: "View recent activity on the network",
	RunE:  activityFunc,
}

// TODO: move all this to a separate client code
func activityFunc(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "expected exactly 0 arguments, got %d", len(args))
		os.Exit(128)
	}
	cli := client.New(uri, requestTimeout)
	activity, err := cli.RecentActivity()
	if err != nil {
		return err
	}
	if err := client.PPActivity(activity); err != nil {
		return err
	}
	return nil
}