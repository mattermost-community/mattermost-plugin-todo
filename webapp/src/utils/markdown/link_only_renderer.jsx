// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import RemoveMarkdown from './remove_markdown';
import {getScheme} from 'utils/url';

export default class LinkOnlyRenderer extends RemoveMarkdown {
    link(href, title, text) {
        let outHref = href;

        if (!getScheme(href)) {
            outHref = `http://${outHref}`;
        }

        let output = `<a class="theme markdown__link" href="${outHref}" target="_blank"`;

        if (title) {
            output += ' title="' + title + '"';
        }

        output += `>${text}</a>`;

        return output;
    }
}

function getScheme(url) {
    const match = (/([a-z0-9+.-]+):/i).exec(url);

    return match && match[1];
}
