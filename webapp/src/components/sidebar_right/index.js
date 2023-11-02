// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getSiteURL, getTodoToast} from '../../selectors';
import {remove, fetchAllIssue, openAssigneeModal, openAddCard, closeAddCard, complete, bump, accept, telemetry, setRhsVisible} from '../../actions';

import SidebarRight from './sidebar_right.jsx';

function mapStateToProps(state) {
    return {
        allIssues: state['plugins-com.mattermost.plugin-todo'].allIssues,
        todoToast: getTodoToast(state),
        siteURL: getSiteURL(state),
        rhsState: state['plugins-com.mattermost.plugin-todo'].rhsState,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            remove,
            complete,
            accept,
            bump,
            fetchAllIssue,
            openAddCard,
            closeAddCard,
            openAssigneeModal,
            telemetry,
            setVisible: setRhsVisible,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);
