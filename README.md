<!--
  Copyright 2022 TrueLevel SA

  This Source Code Form is subject to the terms of the Mozilla Public
  License, v. 2.0. If a copy of the MPL was not distributed with this
  file, You can obtain one at https://mozilla.org/MPL/2.0/.

  SPDX-License-Identifier: MPL-2.0
-->

# go.swapbox.cash/nd300cli (ND-300KM Driver)

[![License](https://img.shields.io/badge/license-MPL--2.0-orange)](https://gitlab.com/TrueLevel/swapbox/nd300cli/-/blob/main/LICENSE)

Interactive CLI tool to interact with an [ND-300KM Note dispencer from ICT Corp.][product] using [`go.swapbox.cash/nd300`][lib].

> ⚠️ **This tool is in active development and is meant for testing & evaluation purposes only.** ⚠️
> 
> This software is distributed as is, without warranty.
> The authors are not liable for any claim, damage, or financial loss 
> related to the use of this tool.

## Installation

```shell
go install go.swapbox.cash/nd300cli@latest
```
## Usage

Just run `nd300cli` and select the serial interface of your note dispenser.

The serial port and machine number can also be specified via the command line:

```shell
Usage of nd300cli:
  -m int
    	The machine number (must be a number between 0 and 7 included, -1 to ignore.) (default -1)
  -p string
    	The machine serial port (optional)
```

### Tests

TODO 😞

[lib]: https://gitlab.com/TrueLevel/swapbox/nd300
[product]: http://www.ictgroup.tw/pro_cen.php?prod_id=70
