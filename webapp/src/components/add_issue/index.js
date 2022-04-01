import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { add, autocompleteUsers, openAssigneeModal, removeAssignee } from 'actions';
import { getMessage, getPostID, getAssignee, getCurrentTeamRoute, isAddCardVisible } from 'selectors';

import AddIssue from './add_issue';

const mapStateToProps = (state) => ({
    visible: isAddCardVisible(state),
    message: getMessage(state) + (getPostID(state) ? '\n[Permalink](' + getCurrentTeamRoute(state) + 'pl/' + getPostID(state) + ')' : ''),
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
