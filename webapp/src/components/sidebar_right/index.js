// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getItems, getInItems, getOutItems} from '../../selectors';
import {remove, list, openRootModal, complete, bump, enqueue} from '../../actions';

import SidebarRight from './sidebar_right.jsx';

function mapStateToProps(state) {
    return {
        todos: getItems(state),
        inTodos: getInItems(state),
        outTodos: getOutItems(state),
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            remove,
            complete,
            enqueue,
            bump,
            list,
            openRootModal,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);
