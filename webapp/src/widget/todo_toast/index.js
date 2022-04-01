// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { getTodoToast } from '../../selectors';
import { closeTodoToast } from '../../actions';

import TodoToast from './todo_toast';

function mapStateToProps(state) {
    return {
        content: getTodoToast(state),
    };
}

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeTodoToast,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TodoToast);
