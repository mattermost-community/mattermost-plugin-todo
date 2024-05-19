// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import { test } from "@e2e-support/test_fixture";
import commands from "./todo_plugin.spec";

import "../support/init_test";

// Test if plugin shows the correct suggestions
// test.describe("testing todo command", commands.todo);

// Test if adding todo works correctly
// test.describe("testing add todo command", commands.addTodo);

// Test if listing todo works correctly
test.describe("testing list todo command", commands.listTodo);

// Test if plugin actions work correctly
// test.describe("testing help command", commands.help);
