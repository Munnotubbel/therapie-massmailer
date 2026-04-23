#!/bin/bash
cd "$(dirname "$0")"/..
./bin/massmailer_mac_apple_silicon || ./bin/massmailer_mac_intel || ./bin/massmailer
