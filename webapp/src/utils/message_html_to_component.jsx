// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {Parser, ProcessNodeDefinitions} from 'html-to-react';

/*
 * Converts HTML to React components using html-to-react.
 */
export function messageHtmlToComponent(html, isRHS, options = {}) {
    if (!html) {
        return null;
    }

    const parser = new Parser();
    const processNodeDefinitions = new ProcessNodeDefinitions(React);

    function isValidNode() {
        return true;
    }

    const processingInstructions = [

        // Workaround to fix MM-14931
        {
            replaceChildren: false,
            shouldProcessNode: (node) => node.type === 'tag' && node.name === 'input' && node.attribs.type === 'checkbox',
            processNode: (node) => {
                const attribs = node.attribs || {};
                node.attribs.checked = Boolean(attribs.checked);

                return React.createElement('input', {...node.attribs});
            },
        },
    ];

    processingInstructions.push({
        shouldProcessNode: () => true,
        processNode: processNodeDefinitions.processDefaultNode,
    });

    return parser.parseWithInstructions(html, isValidNode, processingInstructions);
}

export default messageHtmlToComponent;
