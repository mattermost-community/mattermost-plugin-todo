// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import SidebarButtons from '../sidebar_buttons';

export default class SidebarHeader extends React.PureComponent {
    static propTypes = {
        show: PropTypes.bool.isRequired,
        theme: PropTypes.object.isRequired,
        visible: PropTypes.bool.isRequired,
    };

    render() {
        if (this.props.show || !this.props.visible) {
            return null;
        }

        return (
            <SidebarButtons
                theme={this.props.theme}
                isTeamSidebar={false}
            />
        );
    }
}
