// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {fetchAllIssue, updateRhsState, telemetry} from '../../actions';

import SidebarButtons from './sidebar_buttons.jsx';

function mapStateToProps(state) {
    return {
        allIssues: state['plugins-com.mattermost.plugin-todo'].allIssues,
        showRHSPlugin: state['plugins-com.mattermost.plugin-todo'].rhsPluginAction,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            fetchAllIssue,
            updateRhsState,
            telemetry,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarButtons);
