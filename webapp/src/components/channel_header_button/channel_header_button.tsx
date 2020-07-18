// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

type Props = {
    isVisible: boolean,
};

export default function ChannelHeaderButton(props: Props) {
    return (
        <span className={props.isVisible ? 'channel-header__icon--active' : ''} >
            <i className='icon fa fa-list '/>
        </span>
    );
}
