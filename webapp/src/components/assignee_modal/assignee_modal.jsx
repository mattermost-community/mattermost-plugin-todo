import React from 'react';
import PropTypes from 'prop-types';

import AssigneeForm from './assignee_form';

const AssigneeModal = (
    {
        visible,
        close,
        autocompleteUsers,
        theme,
        getAssignee,
        removeAssignee,
        removeEditingTodo,
        changeAssignee,
        editingTodo,
    },
) => {
    if (!visible) {
        return null;
    }

    return (
        <AssigneeForm
            autocompleteUsers={autocompleteUsers}
            changeAssignee={changeAssignee}
            close={close}
            editingTodo={editingTodo}
            getAssignee={getAssignee}
            removeAssignee={removeAssignee}
            removeEditingTodo={removeEditingTodo}
            theme={theme}
            visible={visible}
        />
    );
};

AssigneeModal.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
    autocompleteUsers: PropTypes.func.isRequired,
    getAssignee: PropTypes.func.isRequired,
    editingTodo: PropTypes.string.isRequired,
    removeAssignee: PropTypes.func.isRequired,
    removeEditingTodo: PropTypes.func.isRequired,
    changeAssignee: PropTypes.func.isRequired,
};

export default AssigneeModal;
