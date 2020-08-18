// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import {isTeamSidebarVisible} from 'selectors';

import TeamSidebar from './team_sidebar.jsx';

function mapStateToProps(state) {
    const members = state.entities.teams.myMembers || {};
    return {
        show: Object.keys(members).length > 1,
        visible: isTeamSidebarVisible(state),
    };
}

export default connect(mapStateToProps)(TeamSidebar);
