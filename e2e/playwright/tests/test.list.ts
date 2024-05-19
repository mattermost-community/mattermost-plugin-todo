// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import { test } from "@playwright/test";
import commands from "./todo_plugin.spec";

import "../support/init_test";

// Test if plugin shows the correct suggestions
test.describe("testing todo command", commands.todo);

// Test if plugin actions work correctly
test.describe("testing help command", commands.help);
