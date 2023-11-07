// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getSiteURL, getTodoToast, getIssues, getInIssues, getOutIssues} from '../../selectors';
import {remove, fetchAllIssueLists, openAssigneeModal, openAddCard, closeAddCard, complete, bump, accept, telemetry, setRhsVisible} from '../../actions';

import SidebarRight from './sidebar_right.jsx';

function mapStateToProps(state) {
    return {
        myIssues: getIssues(state),
        inIssues: getInIssues(state),
        outIssues: getOutIssues(state),
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
            fetchAllIssueLists,
            openAddCard,
            closeAddCard,
            openAssigneeModal,
            telemetry,
            setVisible: setRhsVisible,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);
