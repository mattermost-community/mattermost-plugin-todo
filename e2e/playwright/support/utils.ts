// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {Page} from '@playwright/test';

import Client4 from '@mattermost/client/client4';
import {UserProfile} from '@mattermost/types/users';


export const waitForNewMessages = async (page: Page) => {
    await page.waitForTimeout(1000);

    // This should be able to be waited based on locators instead of pure time-based
    // The following code work "almost" always. Commented for now to have green tests.
    // await page.locator('#postListContent').getByTestId('NotificationSeparator').getByText('New Messages').waitFor();
};

export const getTodoBotDMPageURL = async (client: Client4, teamName: string, userId: string) => {
    let team = teamName;
    if (team === '') {
        const teams = await client.getTeamsForUser(userId);
        team = teams[0].name;
    }
    return `${team}/messages/@todo`;
};

export const fillTextField = async (name: string, value: string, page: Page) => {
    await page.getByTestId(`${name}`).fill(value);
};

export const submitDialog = async (page: Page) => {
    await page.click('#interactiveDialogSubmit');
};

export const fillMessage = async (message: string, page: Page) => {
    await fillTextField('post_textbox', message, page )
};

export const postMessage = async (message: string, page: Page) => {
    await fillMessage(message, page)
    await page.getByTestId('SendMessageButton').click();
};

export const cleanUpBotDMs = async (client: Client4, userId: UserProfile['id'], botUsername: string) => {
    const bot = await client.getUserByUsername(botUsername);

    const userIds = [userId, bot.id];
    const channel = await client.createDirectChannel(userIds);
    const posts = await client.getPosts(channel.id);

    const deletePostPromises = Object.keys(posts.posts).map(client.deletePost);
    await Promise.all(deletePostPromises);
};

export const getSlackAttachmentLocatorId = (postId: string) => {
    return `#post_${postId} .attachment__body`;
};

export const getPostMessageLocatorId = (postId: string) => {
    return `#post_${postId} .post-message`;
};

export const getLastPost = async (page: Page) => {
  const lastPost = page.getByTestId("postView").last();
  await lastPost.waitFor();
  return lastPost;
};
