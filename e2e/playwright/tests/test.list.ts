// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {test} from "@e2e-support/test_fixture";
import commands from "./todo_plugin.spec";

import "../support/init_test";

// Test if plugin shows the correct suggestions for command autocomplete
test.describe("command autocomplete", commands.autocomplete);

// Test `/todo add` commands
// test.describe("commands/add", commands.add);

// Test `/todo list` commands
// test.describe("commands/list", commands.list);

// Test `/todo help` commands
// test.describe("commands/help", commands.help);
