import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {closeRootModal, add, autocompleteUsers} from 'actions';
import {isRootModalVisible, getMessage, getPostID} from 'selectors';

import Root from './root';

const mapStateToProps = (state) => ({
    visible: isRootModalVisible(state),
    message: getMessage(state),
    postID: getPostID(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeRootModal,
    submit: add,
    autocompleteUsers,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(Root);
