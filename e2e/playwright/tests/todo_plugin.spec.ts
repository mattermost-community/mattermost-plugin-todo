// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// ***************************************************************

import {expect, test} from "@e2e-support/test_fixture";
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
  autocomplete: () => {
    test('/todo', async ({ page }) => {
      const slash = new SlashCommandSuggestions(
        page.locator("#suggestionList")
      );

      // # Type command to show suggestions
      await fillMessage('/todo', page);

      // * Assert suggestions are visible
      await expect(slash.container).toBeVisible();

      // * Assert todo [command] is visible
      await expect(slash.getItemTitleNth(0)).toHaveText("todo [command]");

      await expect(slash.getItemDescNth(0)).toContainText(
        "Available commands:"
      );
    });
  },

  help: () => {
    test('/todo help', async ({ page }) => {
      // # Run command to trigger help
    await postMessage('/todo help', page);

    // # Get ephemeral post response from todo help command
    const lastPost = await getLastPost(page);

    // * Assert "help" is in the post body
    await expect(lastPost).toContainText("help");

    // * Assert if length of content shown is greater than 10 lines
    const postBody = await lastPost.textContent();
    const postBodyLines = postBody ? postBody.split('\n') : [];
    expect(postBodyLines.length).toBeGreaterThanOrEqual(10);
    });
  },

  add: () => {
    test("/todo add <message>", async ({ page }) => {
      const todoMessage = "Don't forget to be awesome";

      // # Run command to add todo
      await postMessage(`/todo add ${todoMessage}`, page);

      // # Get ephemeral post response from todo add command
      const post = await getLastPost(page);

      await expect(post).toBeVisible();

      await expect(post).toContainText("Added Todo. Todo List:");
      await expect(post).toContainText(todoMessage);

      // * Assert added todo is visible
      await expect(post).toContainText(todoMessage);
    });
  },

  list: () => {
    test("/todo list", async ({ page }) => {
      const todoMessage = "Don't forget to be awesome";

      // # Run command to add todo
      await postMessage(`/todo add ${todoMessage}`, page);

      // # Type command to list todo
      await postMessage('/todo list', page);

      // # Get ephemeral post response from todo list command
      const post = await getLastPost(page);

      await expect(post).toBeVisible();

      await expect(post).toContainText("Todo List:");
      await expect(post).toContainText(todoMessage);

      // * Assert added todo is visible
      await expect(post).toContainText(todoMessage);
    });
  },
};
