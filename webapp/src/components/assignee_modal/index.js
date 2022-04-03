import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {autocompleteUsers, closeAssigneeModal, getAssignee, removeAssignee, removeEditingTodo, changeAssignee} from 'actions';
import {isAssigneeModalVisible, subMenu, getEditingTodo} from 'selectors';

import AssigneeModal from './assignee_modal';

const mapStateToProps = (state) => ({
    visible: isAssigneeModalVisible(state),
    subMenu: subMenu(state),
    editingTodo: getEditingTodo(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({
    changeAssignee,
    getEditingTodo,
    removeEditingTodo,
    getAssignee,
    removeAssignee,
    close: closeAssigneeModal,
    autocompleteUsers,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(AssigneeModal);
