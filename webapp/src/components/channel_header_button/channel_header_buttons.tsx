// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, { FunctionComponent } from 'react';

type HighlightProps = {
    shouldHighlight: boolean
};

export const Highlight: FunctionComponent<HighlightProps> = ({ shouldHighlight }) => <aside>
    <span className={shouldHighlight ? 'channel-header__icon--active' : ''}>
        <i className='icon fa fa-list '/>
    </span>
</aside>
