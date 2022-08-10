// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {openAssigneeModal, openTodoToast, setEditingTodo, editIssue} from '../../actions';

import TodoItem from './todo_item';

const mapDispatchToProps = (dispatch) => bindActionCreators({
    editIssue,
    openAssigneeModal,
    setEditingTodo,
    openTodoToast,
}, dispatch);

export default connect(null, mapDispatchToProps)(TodoItem);
