// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// ***************************************************************

import {expect, test} from '@e2e-support/test_fixture';
import SlashCommandSuggestions from 'support/components/slash_commands';
import {fillMessage, getLastPost, getTodoBotDMPageURL, postMessage} from 'support/utils';

test.beforeEach(async ({ page, pw }) => {
  const {adminClient, adminUser} = await pw.getAdminClient();
  if (adminUser === null) {
    throw new Error('can not get adminUser');
  }
  const dmURL = await getTodoBotDMPageURL(adminClient, '', adminUser.id);
  await page.goto(dmURL, {waitUntil: 'load'});
});

export default {
  setup: () => {
    test('checking available commands', async ({ page }) => {
      const slash = new SlashCommandSuggestions(page.locator('#suggestionList'));

      // # Run command to trigger todo
      await fillMessage('/todo', page);

      // * Assert suggestions are visible
      await expect(slash.container).toBeVisible();

      // * Assert todo [command] is visible
      await expect(slash.getItemTitleNth(0)).toHaveText('todo [command]');

      await expect(slash.getItemDescNth(0)).toHaveText('Available commands: list, add, pop, send, settings, help');
    });
  },
  commands: () => {
    test("list action", async ({ pages, page, pw }) => {
      const c = new pages.ChannelsPage(page);
      const todoMessage = "Don't forget to be awesome";

      // # Run command to add todo
      postMessage(`/todo add ${todoMessage}`, page);

      // # Type command to list todo
      await fillMessage("/todo list ", page);
      const slash = new SlashCommandSuggestions(
        page.locator("#suggestionList")
      );
      // * Assert suggestions are visible
      await expect(slash.container).toBeVisible();
      await expect(slash.getItemTitleNth(1)).toHaveText("in (optional)");
      await expect(slash.getItemDescNth(1)).toHaveText("Received Todos");
      await expect(slash.getItemTitleNth(2)).toHaveText("out (optional)");
      await expect(slash.getItemDescNth(2)).toHaveText("Sent Todos");

      // # Run command to list todo
      await postMessage('/todo list', page);

      // # Grab the last post
      const post = await getLastPost(page);
      const postBody = post.locator(".post-message__text-container");

      // * Assert post body has correct title
      await expect(postBody).toContainText("Todo List:");

      // * Assert added todo is visible
      await expect(postBody).toContainText(todoMessage);
    });
  },
};

