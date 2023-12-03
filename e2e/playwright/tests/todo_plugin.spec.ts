// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// ***************************************************************

import { expect, test } from "@e2e-support/test_fixture";
import Client4 from "@mattermost/client/client4";
import { UserProfile } from "@mattermost/types/users";
import SlashCommandSuggestions from "support/components/slash_commands";
import { fillMessage, getTodoBotDMPageURL } from "support/utils";

let adminClient: Client4, adminUser: UserProfile | null;

test.beforeEach(async ({ page, pw }) => {
  const data = await pw.getAdminClient();
  adminClient = data.adminClient;
  adminUser = data.adminUser;
  if (adminUser === null) {
    throw new Error("can not get adminUser");
  }
  const dmURL = await getTodoBotDMPageURL(adminClient, "", adminUser.id);
  await page.goto(dmURL, { waitUntil: "load" });
});

export default {
  setup: () => {
    test("checking available commands", async ({ pages, page, pw }) => {
      if (adminUser) {
        const dmURL = await getTodoBotDMPageURL(adminClient, "", adminUser.id);
        await page.goto(dmURL, { waitUntil: "load" });
        const c = new pages.ChannelsPage(page);
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
      }
    });
  },
  actions: () => {
    test("help action", async ({ pages, page, pw }) => {
      if (adminUser) {
        const dmURL = await getTodoBotDMPageURL(adminClient, "", adminUser.id);
        await page.goto(dmURL, { waitUntil: "load" });
        const c = new pages.ChannelsPage(page);
        const slash = new SlashCommandSuggestions(
          page.locator("#suggestionList")
        );

        // # Run command to trigger help
        await c.postMessage("/todo help");

        // # Grab the last post
        const post = await c.getLastPost();
        const postBody = post.container.locator(
          ".post-message__text-container"
        );

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
      }
    });
  },
};
