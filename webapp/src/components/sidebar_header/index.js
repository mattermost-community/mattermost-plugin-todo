// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import {isButtonSidebarVisible} from 'selectors';

import SidebarHeader from './sidebar_header.jsx';

function mapStateToProps(state) {
    const members = state.entities.teams.myMembers || {};
    return {
        show: Object.keys(members).length > 1,
        visible: isButtonSidebarVisible(state),
    };
}

export default connect(mapStateToProps)(SidebarHeader);
