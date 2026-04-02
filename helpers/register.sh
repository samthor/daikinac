#!/bin/bash

# This is a script which helps register a new UUID with a Daikin aircon device.
#
# Usage:
#   ./register.sh <IP> <KEY>
#
# where IP is the device's IP, and KEY is the auth key physically printed on the WiFi dongle.

set -eu

IP=$1
KEY=$2
UUID=$(uuid | sed 's/-//g')

curl --insecure -H "X-Daikin-uuid: $UUID" "https://$IP/common/register_terminal?key=$KEY"

echo "Registered UUID: $UUID"
