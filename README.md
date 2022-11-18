# Mattermost Todo Plugin

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-todo/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-todo)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-todo/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-todo)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-todo)](https://github.com/mattermost/mattermost-plugin-todo/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-todo/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-todo/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

**Maintainer:** [@larkox](https://github.com/larkox)
**Co-Maintainer:** [@jfrerich](https://github.com/jfrerich)

A plugin to track Todo issues in a list and send you daily reminders about your Todo list.

**[Help Wanted](https://github.com/mattermost/mattermost-plugin-todo/issues?utf8=%E2%9C%93&q=is%3Aopen+label%3A%22up+for+grabs%22+label%3A%22help+wanted%22+sort%3Aupdated-desc)**

## Install

1. Go the releases page and download the latest release.
2. In Mattermost, go to **System Console > Plugin Management** and upload it.
3. Start using it!

## Usage

![image](https://user-images.githubusercontent.com/83284294/194702321-995ea864-c11e-400b-9129-0e9e16e65851.png)


To add an issue to your Todo list, do one one of the following:

* Open the sidebar from the channel header and select **Add new issue**.
* Type `/todo add <your Todo message here>` into the textbox and send.
* Select the dropdown menu from a post and click "Add Todo".

To view your Todo list, do one of the following:

* Select the button in the channel header to open the Todo list in the right sidebar.
* Type `/todo list` into the textbox and send.

To remove an issue from your list:

* Open the sidebar from the channel header and select **Done** or **Won't Do** below the issue you want to remove.
* Type `/todo pop` into the text and send to remove the top issue in the list.

To send an issue to another user:

* Open the sidebar from the channel header and select **Add new issue**, the select the user you want to send the issue to.
* Type `/todo send <username> <your Todo message here>` into the textbox and send.

You'll get a daily reminder of the issues you need to complete from the `Todo` bot. The message is only sent if you have issues on your Todo list.
