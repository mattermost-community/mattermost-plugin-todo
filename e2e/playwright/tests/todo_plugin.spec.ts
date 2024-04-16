// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// ***************************************************************

import {expect, test} from '@e2e-support/test_fixture';
import SlashCommandSuggestions from 'support/components/slash_commands';
import {fillMessage, getTodoBotDMPageURL} from 'support/utils';

export default {
    connected: () => {
        test.describe('available commands', () => {
            test('with just the main command', async ({pages, page, pw}) => {
               
            const {adminClient, adminUser} = await pw.getAdminClient();
            if (adminUser === null) {
                throw new Error('can not get adminUser');
            }
            const dmURL = await getTodoBotDMPageURL(adminClient, '', adminUser.id);
            await page.goto(dmURL, {waitUntil: 'load'});

            const c = new pages.ChannelsPage(page);
            const slash = new SlashCommandSuggestions(page.locator('#suggestionList'));

            // # Run incomplete command to trigger help
            await fillMessage('/todo', page);

            // * Assert suggestions are visible
            await expect(slash.container).toBeVisible();

            // * Assert help is visible
            await expect(slash.getItemTitleNth(0)).toHaveText('todo [command]');

            await expect(slash.getItemDescNth(0)).toHaveText('Available commands: list, add, pop, send, settings, help');
            });
        });
    },
};

