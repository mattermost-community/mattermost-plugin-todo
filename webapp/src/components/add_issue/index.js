import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {add, autocompleteUsers, openAssigneeModal, removeAssignee} from 'actions';
import {getMessage, getPostID, getAssignee, getCurrentTeamRoute, isAddCardVisible} from 'selectors';

import AddIssue from './add_issue';
import { openTodoToast } from 'src/actions';

function mapStateToProps(state) {
    const postID = getPostID(state);

    let permalink = '';
    if (postID) {
        permalink = `\n[Permalink](${getCurrentTeamRoute(state)}pl/${postID})`;
    }

    const message = getMessage(state) + permalink;

    return {
        visible: isAddCardVisible(state),
        message,
        postID: getPostID(state),
        assignee: getAssignee(state),
    };
}

const mapDispatchToProps = (dispatch) => bindActionCreators({
    removeAssignee,
    submit: add,
    autocompleteUsers,
    openAssigneeModal,
    openTodoToast,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddIssue);
