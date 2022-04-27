package main

// POWER CONTROLLER DEFINITIONS
// Persistent state - The outlet state will revert to the persistent state after a power cycle or reboot.
// Set persistent state: uom set "relay/outlets/2/state" true
// Get persistent state: uom get relay/outlets/0/state
// *Geting the persistent state may not reflect the physical state, if the transient state has been set.

// Transient state - Temporary outlet state. The outlet may not return to the set state, but will revert to the persistent state after a power cycle or reboot.
// Set transient state: uom set relay/outlets/0/transient_state true
// Get transient state: uom get relay/outlets/0/transient_state
// *Getting the transient state is "usually" the pysical state, but if there is a delay,  e.g. cycle(), the physical state may not have been set yet.

// Physical state - Get current physical outlet state of the outlet/relay.
// Physical state: uom get relay/outlets/0/physical_state
// *The physical state cannot be set by the user. Users set the Transient state or Persistent state. The unit will change the physical state.

// Note: Setting the "state" or "transient state" to "on" may not happen immediatey, as the on sequence delays and cycle delays will be honored.

// SAMPLE UOM COMMANDS
// // Lines preceded by these slashes "//" are comments and not part of the code or command

// // VERSION INFO
// uom get "relay/model"              // Model of the power controller
// uom get "config/version            // Display the Application (left-side) firmware version
// uom get "relay/version"            // Display the AVR/maintenance (right-side) firmware version

// // RELAY CONTROL
// Relays (0 indexed)
// uom set "relay/outlets/2/transient_state" "true"   // Turn on relay 3. The state will revert to the persistent state after a power cycle or reboot.
// uom set "relay/outlets/2/state" "true"                 // Turn on relay 3
// uom set "relay/outlets/2/state" "false"                // Turn off relay 3
// uom get "relay/outlets/2/state"                           // This is the configured (persistent) state and may not be the physical state
// uom get "relay/outlets/2/physical_state"             // Get the physical state of the outlet/relay

// Change the state of several outlets as simultaneously as possible. (Requires firmware 1.10.11+)
// uom invoke /relay/set_outlet_transient_states "[[0,true],[3,false],[4,true],[5,false],[7,false]]"

// uom invoke "relay/outlets/0/cycle"          // Cycle relay 1

// uom get relay/outlets/0/name                // Get the outet/relay 1 name

// // RUNNING A SCRIPT
// // Some items must be quoted. Quotes within quotes must be escaped ( "\"an item\"" ) or encased in the single quote ( '"an item"' )
// // I will use both in some of the examples below
// uom invoke "script/start/" '{"user_function":"my_custom_script"}'              // start the script
// or
// uom invoke "script/start/" "{\"user_function\":\"lighting_schedule\"}"          // start the script

// // Start a script with argument(s)
// uom invoke "script/start" "{\"source\":\"cycle_relay(3)\",\"user_function\":\"cycle_relay\"}"
// or
// uom invoke "script/start" '{"source":"cycle_relay(3)","user_function":"cycle_relay"}'

// uom get "script/threads"                      // show running threads

// uom get "script/threads/4/label"            // Get the name of a running thread

// uom invoke "script/stop/" '"1"'               // stop the script
// uom invoke "script/stop/" "\"all\""          // stop all scripts

// // AUTHENTICATION
// uom get "auth/cookie_timeout"
// uom set "auth/cookie_timeout" "28800"   // 8 hours
// uom set "auth/cookie_timeout" "36000"   // 10 hours
// The uom library can be used in scripts. Scripts using the uom library can be launched from shell.

// # get outlet states
// #!/usr/bin/env lua
// local uom=require"uom"
// local result={}
// for i,outlet in uom.ipairs(uom.relay.outlets) do
//    result[i]=outlet.physical_state
// end
// print(uom.json.encode(result))

// #list outlet names
// #!/usr/bin/env lua
// local uom=require"uom"
// for i,outlet in uom.ipairs(uom.relay.outlets) do
//   print(uom.relay.outlets[i].name)
// end

// #Change the state of several outlets as simultaneously as possible. (Requires firmware 1.10.11+)
// #!/usr/bin/env lua

// local uom=require("uom")

// local null=uom.null

// uom.relay.set_outlet_transient_states({{2,false},{4,false},{5,true},{7,true}})
