// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {test} from "@playwright/test";
import core from "./todo_plugin.spec";

import "../support/init_test";

// Test if plugin is setup correctly
test.describe("setup", core.setup);

// Test various plugin actions
test.describe("actions", core.actions);
