import { combineReducers } from 'redux';

import {
    OPEN_ROOT_MODAL,
    GET_ASSIGNEE,
    REMOVE_ASSIGNEE,
    CLOSE_ROOT_MODAL,
    OPEN_ASSIGNEE_MODAL,
    CLOSE_ASSIGNEE_MODAL,
    OPEN_TODO_TOAST,
    CLOSE_TODO_TOAST,
    OPEN_ADD_CARD,
    CLOSE_ADD_CARD,
    GET_ISSUES,
    GET_IN_ISSUES,
    GET_OUT_ISSUES,
    RECEIVED_SHOW_RHS_ACTION,
    UPDATE_RHS_STATE,
    SET_RHS_VISIBLE,
    SET_HIDE_TEAM_SIDEBAR_BUTTONS,
} from './action_types';

const rootModalVisible = (state = false, action) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        return true;
    case CLOSE_ROOT_MODAL:
        return false;
    default:
        return state;
    }
};

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

const issues = (state = [], action) => {
    switch (action.type) {
    case GET_ISSUES:
        return action.data;
    default:
        return state;
    }
};

const inIssues = (state = [], action) => {
    switch (action.type) {
    case GET_IN_ISSUES:
        return action.data;
    default:
        return state;
    }
};

const outIssues = (state = [], action) => {
    switch (action.type) {
    case GET_OUT_ISSUES:
        return action.data;
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

function isTeamSidebarHidden(state = false, action) {
    switch (action.type) {
    case SET_HIDE_TEAM_SIDEBAR_BUTTONS:
        return action.payload;
    default:
        return state;
    }
}

export default combineReducers({
    currentAssignee,
    addCardVisible,
    rootModalVisible,
    assigneeModalVisible,
    todoToast,
    postID,
    issues,
    inIssues,
    outIssues,
    rhsState,
    rhsPluginAction,
    isRhsVisible,
    isTeamSidebarHidden,
});
