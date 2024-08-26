// // Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// // See LICENSE.txt for license information.

// // ***************************************************************
// // - [#] indicates a test step (e.g. # Go to a page)
// // - [*] indicates an assertion (e.g. * Check the title)
// // ***************************************************************

import { test, expect, Locator } from "@playwright/test";
import {
  MattermostContainer,
  MattermostPlugin,
  login,
  logout,
} from "mattermost-plugin-e2e-test-utils";

class SlashCommandSuggestions {
  constructor(readonly container: Locator) {
    this.container = container;
  }

  getItems() {
    return this.container.getByRole("button");
  }

  getItemNth(n: number) {
    return this.container.getByRole("button").nth(n);
  }
  getItemTitleNth(n: number) {
    return this.getItemNth(n).locator(".slash-command__title");
  }
  getItemDescNth(n: number) {
    return this.getItemNth(n).locator(".slash-command__desc");
  }

  // The text must be exact and complete, otherwise won't match the item
  getItemByText(text: string) {
    return this.container.getByRole("button", { name: text });
  }
}

type PluginConfig = {
  clientId: string;
};

let mattermost: MattermostContainer;
let pluginInstance: MattermostPlugin<PluginConfig>;

test.beforeAll(async () => {
  pluginInstance = new MattermostPlugin<PluginConfig>({
    pluginId: "com.mattermost.demo-plugin",
    pluginConfig: {
      clientId: "client-id",
    },
  }).withLocalBinary("./dist");

  mattermost = await new MattermostContainer()
    .withPlugin(pluginInstance)
    .withEnv("MM_FILESETTINGS_ENABLEPUBLICLINK", "true")
    .startWithUserSetup();
});

test.afterAll(async () => {
  await mattermost.stop();
});

test.describe("available commands", () => {
  test("with just the main command", async ({ page }) => {
    const url = mattermost.url();
    await login(page, url, "regularuser", "regularuser");

    const adminClient = await mattermost.getAdminClient();

    const teams = await adminClient.getTeamsForUser(adminClient.userId);
    const team = teams[0].name;

    await page.goto(`${team}/messages/@todo`, { waitUntil: "load" });

    const slash = new SlashCommandSuggestions(page.locator("#suggestionList"));

    // # Run incomplete command to trigger help
    await page.getByTestId("post_textbox").fill("/todo");

    await expect(slash.container).toBeVisible();

    // * Assert help is visible
    await expect(slash.getItemTitleNth(0)).toHaveText("todo [command]");

    await expect(slash.getItemDescNth(0)).toHaveText(
      "Available commands: list, add, pop, send, settings, help"
    );

    await logout(page);
  });
});
