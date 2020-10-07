import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {closeRootModal, add, update, autocompleteUsers} from 'actions';
import {isRootModalVisible, getMessage, getPostID, getSelectedIssue, getCurrentTeamRoute} from 'selectors';

import Root from './root';

const mapStateToProps = (state) => ({
    visible: isRootModalVisible(state),
    message: getMessage(state) + (getPostID(state) ? '\n[Permalink](' + getCurrentTeamRoute(state) + 'pl/' + getPostID(state) + ')' : ''),
    postID: getPostID(state),
    selectedIssue: getSelectedIssue(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeRootModal,
    submit: add,
    update,
    autocompleteUsers,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Root);
