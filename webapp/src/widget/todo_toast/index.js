// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { getTodoToast, getlastTodo } from '../../selectors';
import { closeTodoToast, add, removeLastTodo } from '../../actions';

import TodoToast from './todo_toast';

function mapStateToProps(state) {
    return {
        content: getTodoToast(state),
        lastPost: getlastTodo(state),
    };
}

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeTodoToast,
    submit: add,
    removeLastTodo,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TodoToast);
