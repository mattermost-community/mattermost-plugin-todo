// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import { makeStyleFromTheme, changeOpacity } from 'mattermost-redux/utils/theme_utils';

export default class ChannelHeaderButton extends React.PureComponent {
    static propTypes = {
        isTodoPluginRhsOpen: PropTypes.bool,
    };

    constructor(props) {
        super(props);

        this.state = {
            refreshing: false,
        };
    }

    render() {
        return (
            <span className={`channel-header__icon channel-header__icon--wide ${this.props.isTodoPluginRhsOpen ? 'channel-header__icon--active' : ''}`}>
                <i className='icon fa fa-list '/>
            </span>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        buttonTeam: {
            color: changeOpacity(theme.sidebarText, 0.6),
            display: 'block',
            marginBottom: '10px',
            width: '100%',
        },
        buttonHeader: {
            color: changeOpacity(theme.sidebarText, 0.6),
            textAlign: 'center',
            cursor: 'pointer',
        },
        containerHeader: {
            marginTop: '10px',
            marginBottom: '5px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-around',
            padding: '0 10px',
        },
        containerTeam: {
        },
    };
});