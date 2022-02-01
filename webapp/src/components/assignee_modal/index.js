import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { autocompleteUsers, closeAssigneeModal, getAssignee, removeAssignee } from 'actions';
import { isAssigneeModalVisible, subMenu } from 'selectors';

import AssigneeModal from './assignee_modal';

const mapStateToProps = (state) => ({
    visible: isAssigneeModalVisible(state),
    subMenu: subMenu(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    getAssignee,
    removeAssignee,
    close: closeAssigneeModal,
    autocompleteUsers,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AssigneeModal);
