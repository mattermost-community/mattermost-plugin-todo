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
  getBotDMPageURL,
  getLastPost,
  getTeamName,
  postMessage,
} from "support/utils";

const botUserName = "todo";
let teamName = "";

test.beforeAll(async ({ pw }) => {
  const { adminClient, adminUser } = await pw.getAdminClient();
  if (adminUser === null) {
    throw new Error("can not get adminUser");
  }
  if (teamName === "") {
    teamName = await getTeamName(adminClient, adminUser.id);
  }
});

test.beforeEach(async ({ page }) => {
  const dmURL = await getBotDMPageURL(teamName, botUserName);
  await page.goto(dmURL, { waitUntil: "load" });
});

export default {
  todo: () => {
    const command = "/todo";

    test(`${command}`, async ({ page }) => {
      const slash = new SlashCommandSuggestions(
        page.locator("#suggestionList")
      );

      // # Type command to show suggestions
      await fillMessage(command, page);

      // * Assert suggestions are visible
      await expect(slash.container).toBeVisible();

      // * Assert todo [command] is visible
      await expect(slash.getItemTitleNth(0)).toHaveText("todo [command]");

      await expect(slash.getItemDescNth(0)).toHaveText(
        "Available commands: list, add, pop, send, settings, help"
      );
    });
  },

  help: () => {
    const command = "/todo help";

    test(`${command}`, async ({ pages, page, pw }) => {
      const c = new pages.ChannelsPage(page);

      // # Run command to trigger help
      postMessage(command, page);

      // # Grab the last post
      const lastPost = await getLastPost(page);

      // * Assert all commands are shown in the help text output
      await expect(lastPost).toContainText("add [message]");
      await expect(lastPost).toContainText("list");
      await expect(lastPost).toContainText("list [listName]");
      await expect(lastPost).toContainText("pop");
      await expect(lastPost).toContainText("send [user] [message]");
      await expect(lastPost).toContainText("settings summary [on, off]");
      await expect(lastPost).toContainText(
        "settings allow_incoming_task_requests [on, off]"
      );
      await expect(lastPost).toContainText("help");
    });
  },
};
