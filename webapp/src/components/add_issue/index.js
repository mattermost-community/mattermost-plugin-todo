import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {add, autocompleteUsers, openAssigneeModal, removeAssignee} from 'actions';
import {getMessage, getPostID, getAssignee, getCurrentTeamRoute, isAddCardVisible} from 'selectors';

import AddIssue from './add_issue';

const mapStateToProps = (state) => ({
    visible: isAddCardVisible(state),
    const postID = getPostID(state);

    let permalink = '';
    if (postID) {
        permalink = `[Permalink](${getCurrentTeamRoute(state)}pl/${postID})`;
    }

    const message = `${getMessage(state)}\n${permalink}`;
    postID: getPostID(state),
    assignee: getAssignee(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    removeAssignee,
    submit: add,
    autocompleteUsers,
    openAssigneeModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddIssue);
