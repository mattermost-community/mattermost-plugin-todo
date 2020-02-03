// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {remove, complete, enqueue} from '../../actions';

import PostTypeTodo from './post_type_todo';

function mapStateToProps(state, ownProps) {
    return {
        ...ownProps,
        pendingAnswer: state['plugins-com.mattermost.plugin-todo'].inItems.some((item) => item.id === ownProps.post.props.itemId),
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            remove,
            complete,
            enqueue,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(PostTypeTodo);
