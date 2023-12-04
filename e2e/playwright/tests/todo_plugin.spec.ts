// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// ***************************************************************

import { expect, test } from "@e2e-support/test_fixture";
import SlashCommandSuggestions from "support/components/slash_commands";
import {
  fillMessage,
  getLastPost,
  getTodoBotDMPageURL,
  postMessage,
} from "support/utils";

test.beforeEach(async ({ page, pw }) => {
  const { adminClient, adminUser } = await pw.getAdminClient();
  if (adminUser === null) {
    throw new Error("can not get adminUser");
  }
  const dmURL = await getTodoBotDMPageURL(adminClient, "", adminUser.id);
  await page.goto(dmURL, { waitUntil: "load" });
});

export default {
  setup: () => {
    test("checking available commands", async ({ pages, page, pw }) => {
      const slash = new SlashCommandSuggestions(
        page.locator("#suggestionList")
      );

      // # Run command to trigger todo
      await fillMessage("/todo", page);

      // * Assert suggestions are visible
      await expect(slash.container).toBeVisible();

      // * Assert todo [command] is visible
      await expect(slash.getItemTitleNth(0)).toHaveText("todo [command]");

      await expect(slash.getItemDescNth(0)).toHaveText(
        "Available commands: list, add, pop, send, settings, help"
      );
    });
  },
  actions: () => {
    test("help action", async ({ pages, page, pw }) => {
      const c = new pages.ChannelsPage(page);

      // # Run command to trigger help
      postMessage("/todo help", page);

      // # Grab the last post
      const post = await getLastPost(page);
      const postBody = post.locator(".post-message__text-container");

      // * Assert /todo add [message] command is visible
      await expect(postBody).toContainText(`add [message]`);

      // * Assert /todo list command is visible
      await expect(postBody).toContainText("list");

      // * Assert /todo list [listName] command is visible
      await expect(postBody).toContainText("list [listName]");

      // * Assert /todo pop command is visible
      await expect(postBody).toContainText("pop");

      // * Assert /todo send [user] [message] command is visible
      await expect(postBody).toContainText("send [user] [message]");

      // * Assert /todo settings summary [on, off] command is visible
      await expect(postBody).toContainText("settings summary [on, off]");

      // * Assert /todo settings allow_incoming_task_requests [on, off] command is visible
      await expect(postBody).toContainText(
        "settings allow_incoming_task_requests [on, off]"
      );

      // * Assert /todo help command is visible
      await expect(postBody).toContainText("help");
    });
  },
};
