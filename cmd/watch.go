// Copyright © 2018 Jonathan Pentecost <pentecostjonathan@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch all events sent from a chromecaset device",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := castApplication(cmd, args)
		if err != nil {
			fmt.Printf("unable to get cast application: %v\n", err)
			return
		}
		app.SetDebug(true)
		if err := app.Update(); err != nil {
			fmt.Printf("unable to update cast application: %v\n", err)
			return
		}
		castApplication, castMedia, castVolume := app.Status()
		if castApplication == nil {
			fmt.Printf("Idle, volume=%0.2f muted=%t\n", castVolume.Level, castVolume.Muted)
		} else if castApplication.IsIdleScreen {
			fmt.Printf("Idle (%s), volume=%0.2f muted=%t\n", castApplication.DisplayName, castVolume.Level, castVolume.Muted)
		} else if castMedia == nil {
			fmt.Printf("Idle (%s), volume=%0.2f muted=%t\n", castApplication.DisplayName, castVolume.Level, castVolume.Muted)
		} else {
			metadata := "unknown"
			if castMedia.Media.Metadata.Title != "" {
				md := castMedia.Media.Metadata
				metadata = fmt.Sprintf("title=%q, artist=%q", md.Title, md.Artist)
			}
			fmt.Printf("%s (%s), %s, time remaining=%.0fs/%.0fs, volume=%0.2f, muted=%t\n", castApplication.DisplayName, castMedia.PlayerState, metadata, castMedia.CurrentTime, castMedia.Media.Duration, castVolume.Level, castVolume.Muted)
		}
		// Wait forever
		c := make(chan bool, 1)
		<-c
		return
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}