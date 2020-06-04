// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getIssues, getInIssues, getOutIssues, getSiteURL} from '../../selectors';
import {remove, list, openRootModal, complete, bump, accept} from '../../actions';

import SidebarRight from './sidebar_right.jsx';

function mapStateToProps(state) {
    return {
        todos: getIssues(state),
        inTodos: getInIssues(state),
        outTodos: getOutIssues(state),
        siteURL: getSiteURL(),
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
            openRootModal,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);
