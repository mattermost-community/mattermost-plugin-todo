import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {add, autocompleteUsers, openAssigneeModal, removeAssignee} from 'actions';
import {getMessage, getPostID, getAssignee, getCurrentTeamRoute, isAddCardVisible} from 'selectors';

import AddIssue from './add_issue';

function mapStateToProps(state) {
    const postID = getPostID(state);

    let postPermalink = '';
    if (postID) {
        postPermalink = `${getCurrentTeamRoute(state)}pl/${postID}`;
    }

    const message = getMessage(state);

    return {
        visible: isAddCardVisible(state),
        message,
        postPermalink,
        postID: getPostID(state),
        assignee: getAssignee(state),
    };
}

const mapDispatchToProps = (dispatch) => bindActionCreators({
    removeAssignee,
    submit: add,
    autocompleteUsers,
    openAssigneeModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddIssue);
