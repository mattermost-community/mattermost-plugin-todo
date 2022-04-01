import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { add, autocompleteUsers, openAssigneeModal } from 'actions';
import { getMessage, getPostID, getCurrentTeamRoute } from 'selectors';

import AddIssue from './add_issue';

const mapStateToProps = (state) => ({
    message: getMessage(state) + (getPostID(state) ? '\n[Permalink](' + getCurrentTeamRoute(state) + 'pl/' + getPostID(state) + ')' : ''),
    postID: getPostID(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    submit: add,
    autocompleteUsers,
    openAssigneeModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AddIssue);
