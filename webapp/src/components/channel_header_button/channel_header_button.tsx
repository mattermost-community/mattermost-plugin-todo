// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

type Props = {
    shouldHighlight: boolean
};

export default function ChannelHeaderButton(props: Props) {
    return (
        <span className={props.shouldHighlight ? 'channel-header__icon--active' : ''}>
            <i className='icon fa fa-list '/>
        </span>
    );
}
