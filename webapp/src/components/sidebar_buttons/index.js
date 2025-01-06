// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {fetchAllIssueLists, updateRhsState, telemetry} from '../../actions';

import {getMyIssues, getInIssues, getOutIssues} from '../../selectors';

import SidebarButtons from './sidebar_buttons.jsx';

function mapStateToProps(state) {
    return {
        myIssues: getMyIssues(state),
        inIssues: getInIssues(state),
        outIssues: getOutIssues(state),
        showRHSPlugin: state['plugins-com.mattermost.plugin-todo'].rhsPluginAction,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            fetchAllIssueLists,
            updateRhsState,
            telemetry,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarButtons);
