// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {test} from '@playwright/test';
import core from './todo_plugin.spec';

import '../mattermost-plugin-e2e-test-utils/support/init_test';

test.describe(core.connected);
