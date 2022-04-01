// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { openAssigneeModal, openTodoToast } from '../../actions';

import TodoItem from './todo_item';

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            openAssigneeModal,
            openTodoToast,
        }, dispatch),
    };
}

export default connect(null, mapDispatchToProps)(TodoItem);
