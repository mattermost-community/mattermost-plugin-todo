// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import XRegExp from 'xregexp';

import * as Markdown from './markdown';

const punctuation = XRegExp.cache('[^\\pL\\d]');

const htmlEmojiPattern = /^<p>\s*(?:<img class="emoticon"[^>]*>|<span data-emoticon[^>]*>[^<]*<\/span>\s*|<span class="emoticon emoticon--unicode">[^<]*<\/span>\s*)+<\/p>$/;

// Performs formatting of user posts including converting urls, hashtags,
// @mentions and ~channels to links by taking a user's message and returning a string of formatted html. Also takes
// a number of options as part of the second parameter:
// - singleline - Specifies whether or not to remove newlines. Defaults to false.
// - markdown - Enables markdown parsing. Defaults to true.
// - siteURL - The origin of this Mattermost instance. If provided, links to channels and posts will be replaced with internal
//     links that can be handled by a special click handler.
// - channelNamesMap - An object mapping channel display names to channels. If provided, ~channel mentions will be replaced with
//     links to the relevant channel.
// - team - The current team.
// - minimumHashtagLength - Minimum number of characters in a hashtag. Defaults to 3.
export function formatText(text, inputOptions) {
    if (!text || typeof text !== 'string') {
        return '';
    }

    let output = text;
    const options = Object.assign({}, inputOptions);

    if (!('markdown' in options) || options.markdown) {
        // the markdown renderer will call doFormatText as necessary
        output = Markdown.format(output, options);
        if (output.includes('class="markdown-inline-img"')) {
            /*
            ** remove p tag to allow other divs to be nested,
            ** which allows markdown images to open preview window
            */
            const replacer = (match) => {
                return match === '<p>' ? '<div className="markdown-inline-img__container">' : '</div>';
            };
            output = output.replace(/<p>|<\/p>/g, replacer);
        }
    } else {
        output = sanitizeHtml(output);
        output = doFormatText(output, options);
    }

    // replace newlines with spaces if necessary
    if (options.singleline) {
        output = replaceNewlines(output);
    }

    if (htmlEmojiPattern.test(output.trim())) {
        output = '<span class="all-emoji">' + output.trim() + '</span>';
    }

    return output;
}

// Performs most of the actual formatting work for formatText. Not intended to be called normally.
export function doFormatText(text, options) {
    let output = text;

    const tokens = new Map();

    if (options.channelNamesMap) {
        output = autolinkChannelMentions(output, tokens, options.channelNamesMap, options.team);
    }

    output = autolinkEmails(output, tokens);
    output = autolinkHashtags(output, tokens, options.minimumHashtagLength);

    // reinsert tokens with formatted versions of the important words and phrases
    output = replaceTokens(output, tokens);

    return output;
}

export function sanitizeHtml(text) {
    let output = text;

    // normal string.replace only does a single occurrance so use a regex instead
    output = output.replace(/&/g, '&amp;');
    output = output.replace(/</g, '&lt;');
    output = output.replace(/>/g, '&gt;');
    output = output.replace(/'/g, '&apos;');
    output = output.replace(/"/g, '&quot;');

    return output;
}

// Copied from our fork of commonmark.js
var emailAlphaNumericChars = '\\p{L}\\p{Nd}';
var emailSpecialCharacters = '!#$%&\'*+\\-\\/=?^_`{|}~';
var emailRestrictedSpecialCharacters = '\\s(),:;<>@\\[\\]';
var emailValidCharacters = emailAlphaNumericChars + emailSpecialCharacters;
var emailValidRestrictedCharacters = emailValidCharacters + emailRestrictedSpecialCharacters;
var emailStartPattern = '(?:[' + emailValidCharacters + '](?:[' + emailValidCharacters + ']|\\.(?!\\.|@))*|\\"[' + emailValidRestrictedCharacters + '.]+\\")@';
var reEmail = XRegExp.cache('(^|[^\\pL\\d])(' + emailStartPattern + '[\\pL\\d.\\-]+[.]\\pL{2,4}(?=$|[^\\p{L}]))', 'g');

// Convert emails into tokens
function autolinkEmails(text, tokens) {
    function replaceEmailWithToken(fullMatch, prefix, email) {
        const index = tokens.size;
        const alias = `$MM_EMAIL${index}$`;

        tokens.set(alias, {
            value: `<a class="theme" href="mailto:${email}" rel="noreferrer" target="_blank">${email}</a>`,
            originalText: email,
        });

        return prefix + alias;
    }

    let output = text;
    output = XRegExp.replace(text, reEmail, replaceEmailWithToken);

    return output;
}

function autolinkChannelMentions(text, tokens, channelNamesMap, team) {
    function channelMentionExists(c) {
        return Boolean(channelNamesMap[c]);
    }
    function addToken(channelName, mention, displayName) {
        const index = tokens.size;
        const alias = `$MM_CHANNELMENTION${index}$`;
        let href = '#';
        if (team) {
            href = (window.basename || '') + '/' + team.name + '/channels/' + channelName;
        }

        tokens.set(alias, {
            value: `<a class="mention-link" href="${href}" data-channel-mention="${channelName}">~${displayName}</a>`,
            originalText: mention,
        });
        return alias;
    }

    function replaceChannelMentionWithToken(fullMatch, mention, channelName) {
        let channelNameLower = channelName.toLowerCase();

        if (channelMentionExists(channelNameLower)) {
            // Exact match
            const alias = addToken(channelNameLower, mention, escapeHtml(channelNamesMap[channelNameLower].display_name));
            return alias;
        }

        // Not an exact match, attempt to truncate any punctuation to see if we can find a channel
        const originalChannelName = channelNameLower;

        for (let c = channelNameLower.length; c > 0; c--) {
            if (punctuation.test(channelNameLower[c - 1])) {
                channelNameLower = channelNameLower.substring(0, c - 1);

                if (channelMentionExists(channelNameLower)) {
                    const suffix = originalChannelName.substr(c - 1);
                    const alias = addToken(channelNameLower, '~' + channelNameLower,
                        escapeHtml(channelNamesMap[channelNameLower].display_name));
                    return alias + suffix;
                }
            } else {
                // If the last character is not punctuation, no point in going any further
                break;
            }
        }

        return fullMatch;
    }

    let output = text;
    output = output.replace(/\B(~([a-z0-9.\-_]*))/gi, replaceChannelMentionWithToken);

    return output;
}

export function escapeRegex(text) {
    if (text == null) {
        return '';
    }
    return text.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');
}

const htmlEntities = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
};

export function escapeHtml(text) {
    return text.replace(/[&<>"']/g, (match) => htmlEntities[match]);
}

export function convertEntityToCharacter(text) {
    return text.
        replace(/&lt;/g, '<').
        replace(/&gt;/g, '>').
        replace(/&#39;/g, '\'').
        replace(/&quot;/g, '"').
        replace(/&amp;/g, '&');
}

function autolinkHashtags(text, tokens, minimumHashtagLength = 3) {
    let output = text;

    var newTokens = new Map();
    for (const [alias, token] of tokens) {
        if (token.originalText.lastIndexOf('#', 0) === 0) {
            const index = tokens.size + newTokens.size;
            const newAlias = `$MM_HASHTAG${index}$`;

            newTokens.set(newAlias, {
                value: `<a class='mention-link' href='#' data-hashtag='${token.originalText}'>${token.originalText}</a>`,
                originalText: token.originalText,
                hashtag: token.originalText.substring(1),
            });

            output = output.replace(alias, newAlias);
        }
    }

    // the new tokens are stashed in a separate map since we can't add objects to a map during iteration
    for (const newToken of newTokens) {
        tokens.set(newToken[0], newToken[1]);
    }

    // look for hashtags in the text
    function replaceHashtagWithToken(fullMatch, prefix, originalText) {
        const index = tokens.size;
        const alias = `$MM_HASHTAG${index}$`;

        if (originalText.length < minimumHashtagLength + 1) {
            // too short to be a hashtag
            return fullMatch;
        }

        tokens.set(alias, {
            value: `<a class='mention-link' href='#' data-hashtag='${originalText}'>${originalText}</a>`,
            originalText,
            hashtag: originalText.substring(1),
        });

        return prefix + alias;
    }

    return output.replace(XRegExp.cache('(^|\\W)(#\\pL[\\pL\\d\\-_.]*[\\pL\\d])', 'g'), replaceHashtagWithToken);
}

export function replaceTokens(text, tokens) {
    let output = text;

    // iterate backwards through the map so that we do replacement in the opposite order that we added tokens
    const aliases = [...tokens.keys()];
    for (let i = aliases.length - 1; i >= 0; i--) {
        const alias = aliases[i];
        const token = tokens.get(alias);
        output = output.replace(alias, token.value);
    }

    return output;
}

function replaceNewlines(text) {
    return text.replace(/\n/g, ' ');
}