## Overview

Provides UDP cli tools for testing and integrating with Emanate PowerPath tags.  Each cli tool is pre-compiled for OSX, Linux, and Windows.

The 'udp_sender_tool' simulates Emanate PowerPath tags sending UDP datagrams in the Cisco-CCX format.

The 'udp_receiver_tool' receives Emanate PowerPath tag UDP datagrams, parses the binary messages, and dumps the parsed results to the console.

## Quick Start

### 1. Download the UDP cli tools from the Github 'Releases' screen
   * Mac OSX
     * [Download OSX UDP sender](https://github.com/EmanateWireless/emanate-udp-tools/releases/download/v1.0.2/emanate_udp_sender_osx)
     * [Download OSX UDP receiver](https://github.com/EmanateWireless/emanate-udp-tools/releases/download/v1.0.2/emanate_udp_receiver_osx)
   * Windows
     * [Download Windows 'UDP sender'](https://github.com/EmanateWireless/emanate-udp-tools/releases/download/v1.0.2/emanate_udp_sender.exe)
     * [Download Windows 'UDP receiver'](https://github.com/EmanateWireless/emanate-udp-tools/releases/download/v1.0.2/emanate_udp_receiver.exe)
   * Linux x86
     * [Download Linux 'UDP sender'](https://github.com/EmanateWireless/emanate-udp-tools/releases/download/v1.0.2/emanate_udp_sender_linux_x86)
     * [Download Linux 'UDP receiver'](https://github.com/EmanateWireless/emanate-udp-tools/releases/download/v1.0.2/emanate_udp_receiver_linux_x86)

### 2. Start the 'UDP Receiver'

```
$ ./emanate_udp_receiver_osx

Starting UDP receiver listening on port '9999'
```

### 3. Start the 'UDP Sender'

The '--all' option sends the kitchen-sink of every Emanate UDP option (for testing purposes).

```
$ ./emanate_udp_sender_osx --all

2016/07/05 19:12:55 Sending udp packet to '127.0.0.1:9999' (323 bytes)
2016/07/05 19:12:55 DONE!
```

The UDP receiver will dump the fully parsed UDP packet to the console. Here is an example of a parsed kitchen-sink UDP packet (--all)...

```
UDP PACKET RECEIVED
===================

  - Total Bytes = 323
  - Remote Addr = 127.0.0.1:49536
  - Sequence = 1
  - Header
    - Protocol Version = 0
    - Transmit Power = 17
    - Wifi Channel = 1
    - Burst Length = 1
  - System Group
    - ID = 0
    - Length = 2
    - Product Type = 0
  - Battery Group
    - ID = 2
    - Length = 7
    - Tolerance = 0 %
    - Charge = 80 %
    - Days Remaining = 100
    - Age = 10 days
  - Status Group
    - ID = 3
    - Length = 42
    - Type = 8
    - Status = 'UTIL_STATE=UNPLUGGED'
  - Temperature Group
    - ID = 3
    - Length = 5
    - Type = 1
    - Temperature = 12.34 C
  - Status Group
    - ID = 3
    - Length = 42
    - Type = 8
    - Status = 'DOOR_OPEN_PERCENT=22'
  - Status Group
    - ID = 3
    - Length = 54
    - Type = 8
    - Status = 'HIGH_POWER_MODE_PERCENT=33'
  - Status Group
    - ID = 3
    - Length = 30
    - Type = 8
    - Status = 'BUTTON=PRESSED'
  - Status Group
    - ID = 3
    - Length = 54
    - Type = 8
    - Status = 'TEMP_PROBE_ERROR=UNPLUGGED'
  - Status Group
    - ID = 3
    - Length = 62
    - Type = 8
    - Status = 'TEMP_PROBE_ERROR=INVALID_VALUE'
```

## Advanced Usage

### UDP Sender

The 'emanate_udp_sender' tool provides many options to allow any combination of CCX fields to be included in the UDP packet.

```
$ ./emanate_udp_sender_osx -h

NAME:
   emanate_udp_sender - Emanate PowerPath UDP CCX packet transmitter

USAGE:
   emanate_udp_sender --host <IP> --port <PORT> [options]

VERSION:
   v1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value                    udp target hostname or ip-address (default: "127.0.0.1")
   --port value                    udp target port number (default: 9999)
   --all                           sends all possible udp message options for testing
   --num-dups value                number of duplicate udp packets to send (default: 0)
   --dup-interval-ms value         delay interval between duplicate udp packets (default: 100)
   --seq value                     sequence number of the emanate udp packet (default: 1)
   --util-state value              utility state of 'unplugged', 'off', 'idle', or 'active' (default: "unplugged")
   --temp value                    temperature floating-point value (in celsius) (default: 12.34)
   --battery-charge value          battery charge percentage remaining (0-100) (default: 80)
   --battery-days-remaining value  number of days remaining for battery charge (default: 100)
   --battery-age value             battery age in days (default: 10)
   --battery-tolerance value       battery prediction tolerance percentage (0-100) (default: 0)
   --button-pressed                adds the button-pressed telemetry status
   --door-open-percent value       percentage of time the fridge door has been open (default: 22)
   --high-power-percent value      percentage of time the device ran in high-power-mode (default: 33)
   --probe-unplugged               adds the 'temp probe unplugged' alert telemetry status
   --probe-invalid-value           adds the 'temp probe invalid value' alert telemetry status
   --product-type value            product-type code set in the ccx 'system group' (default: 0)
   --help, -h                      show help
   --version, -v                   print the version
```

### UDP Receiver

The 'emanate_udp_receiver' tool current just allows the user to change to listening UDP port.

```
$ ./emanate_udp_receiver_osx -h

NAME:
   emanate_udp_receiver - Emanate PowerPath UDP CCX packet receiver

USAGE:
   emanate_udp_receiver --port <LISTENING-PORT>

VERSION:
   v1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value   local udp receiver port number (default: 9999)
   --help, -h     show help
   --version, -v  print the version
```
