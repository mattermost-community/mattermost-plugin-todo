// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export function isUrlSafe(url) {
    let unescaped;

    try {
        unescaped = decodeURIComponent(url);
    } catch (e) {
        unescaped = unescape(url);
    }

    unescaped = unescaped.replace(/[^\w:]/g, '').toLowerCase();

    return !unescaped.startsWith('javascript:') && // eslint-disable-line no-script-url
        !unescaped.startsWith('vbscript:') &&
        !unescaped.startsWith('data:');
}

export function getScheme(url) {
    const match = (/([a-z0-9+.-]+):/i).exec(url);

    return match && match[1];
}