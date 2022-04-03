// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getIssues, getInIssues, getOutIssues, getSiteURL, getTodoToast} from '../../selectors';
import {remove, list, openAssigneeModal, openAddCard, closeAddCard, complete, bump, accept, telemetry, setRhsVisible} from '../../actions';

import SidebarRight from './sidebar_right.jsx';

function mapStateToProps(state) {
    return {
        todos: getIssues(state),
        todoToast: getTodoToast(state),
        inTodos: getInIssues(state),
        outTodos: getOutIssues(state),
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
            list,
            openAddCard,
            closeAddCard,
            openAssigneeModal,
            telemetry,
            setVisible: setRhsVisible,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);
