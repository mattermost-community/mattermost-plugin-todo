import {combineReducers} from 'redux';

import {
    GET_ASSIGNEE,
    REMOVE_ASSIGNEE,
    OPEN_ASSIGNEE_MODAL,
    CLOSE_ASSIGNEE_MODAL,
    OPEN_TODO_TOAST,
    CLOSE_TODO_TOAST,
    OPEN_ADD_CARD,
    CLOSE_ADD_CARD,
    SET_EDITING_TODO,
    REMOVE_EDITING_TODO,
    GET_ALL_ISSUES,
    RECEIVED_SHOW_RHS_ACTION,
    UPDATE_RHS_STATE,
    SET_RHS_VISIBLE,
    SET_HIDE_SIDEBAR_BUTTONS,
} from './action_types';

const addCardVisible = (state = false, action) => {
    switch (action.type) {
    case OPEN_ADD_CARD:
        return true;
    case CLOSE_ADD_CARD:
        return false;
    default:
        return state;
    }
};

const assigneeModalVisible = (state = false, action) => {
    switch (action.type) {
    case OPEN_ASSIGNEE_MODAL:
        return true;
    case CLOSE_ASSIGNEE_MODAL:
        return false;
    default:
        return state;
    }
};

const todoToast = (state = null, action) => {
    switch (action.type) {
    case OPEN_TODO_TOAST:
        return action.message;
    case CLOSE_TODO_TOAST:
        return null;
    default:
        return state;
    }
};

const editingTodo = (state = null, action) => {
    switch (action.type) {
    case SET_EDITING_TODO:
        return action.issueID;
    case REMOVE_EDITING_TODO:
        return null;
    default:
        return state;
    }
};

const currentAssignee = (state = null, action) => {
    switch (action.type) {
    case GET_ASSIGNEE:
        return action.assignee;
    case REMOVE_ASSIGNEE:
        return null;
    default:
        return state;
    }
};

const postID = (state = '', action) => {
    switch (action.type) {
    case OPEN_ADD_CARD:
        return action.postID;
    case CLOSE_ADD_CARD:
        return '';
    default:
        return state;
    }
};

const allIssues = (state = {my: [], in: [], out: []}, action) => {
    switch (action.type) {
    case GET_ALL_ISSUES:
        return action.data ?? state;
    default:
        return state;
    }
};

function rhsPluginAction(state = null, action) {
    switch (action.type) {
    case RECEIVED_SHOW_RHS_ACTION:
        return action.showRHSPluginAction;
    default:
        return state;
    }
}

function rhsState(state = '', action) {
    switch (action.type) {
    case UPDATE_RHS_STATE:
        return action.state;
    default:
        return state;
    }
}

function isRhsVisible(state = false, action) {
    switch (action.type) {
    case SET_RHS_VISIBLE:
        return action.payload;
    default:
        return state;
    }
}

function isButtonSidebarHidden(state = false, action) {
    switch (action.type) {
    case SET_HIDE_SIDEBAR_BUTTONS:
        return action.payload;
    default:
        return state;
    }
}

export default combineReducers({
    currentAssignee,
    addCardVisible,
    assigneeModalVisible,
    todoToast,
    editingTodo,
    postID,
    allIssues,
    rhsState,
    rhsPluginAction,
    isRhsVisible,
    isButtonSidebarHidden,
});
