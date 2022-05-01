/* SPDX-License-Identifier: GPL-2.0-or-later */

package main

// var commands = map[string][]string{
// 	"version": {
// 		`uom get "relay/model"`,    // Model of the power controller
// 		`uom get "config/version"`, // Display the Application (left-side) firmware version
// 		`uom get "relay/version"`,  // Display the AVR/maintenance (right-side) firmware version
// 	},
// 	// Relays (0 indexed)
// 	"state": {
// 		`uom get "relay/outlets/%s/state"`,
// 	},
// 	"status": {
// 		`uom get "relay/outlets/%s/physical_state"`,
// 	},
// 	"on": {
// 		`uom set "relay/outlets/%s/transient_state" "true"`,
// 		`uom get "relay/outlets/%s/physical_state"`,
// 	},
// 	"off": {
// 		`uom set "relay/outlets/%s/transient_state" "false"`,
// 		`uom get "relay/outlets/%s/physical_state"`,
// 	},
// 	"cycle": {
// 		`uom invoke "relay/outlets/%s/cycle"`,
// 	},
// 	"name": {
// 		`uom get "relay/outlets/%s/name"`,
// 	},
// }
