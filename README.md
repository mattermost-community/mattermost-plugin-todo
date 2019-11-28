# Mattermost To Do Plugin

A plugin to track to do items in a list and send you daily reminders about your to do list.

## Install

1. Go the releases page and download the latest release.
2. On your Mattermost, go to System Console -> Plugin Management and upload it.
3. Start using it!

## Usage

To add an item to your to do list, do one one of the following:

* Open the sidebar from the channel header and click the "Add new item" button
* Type `/todo add <your to do message here>` into the textbox and send
* Click the on the dropdown menu from a post and click "Add To Do"

To view your to do list, do one of the following:

* Click on the button in the channel header to open the to do list in the right sidebar.
* Type `/todo list` into the textbox and send

To remove an item from your list:

* Open the sidebar from the channel header and click the "X" next to the item you want to remove
* Type `/todo pop` into the text and send to remove the top item in the list

Every day you will get a reminder of the items you need to complete from the `todo` bot. The message is only sent if you have items on your to do list.