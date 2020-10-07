import {combineReducers} from 'redux';

import {OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL, GET_ISSUES, GET_IN_ISSUES, GET_OUT_ISSUES, RECEIVED_SHOW_RHS_ACTION, UPDATE_RHS_STATE, SET_RHS_VISIBLE, SET_HIDE_TEAM_SIDEBAR_BUTTONS} from './action_types';

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

const postID = (state = '', action) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        return action.postID;
    case CLOSE_ROOT_MODAL:
        return '';
    default:
        return state;
    }
};

const selectedPost = (state = null, action) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        if (typeof action.selectedPost !== 'undefined') {
            return action.selectedPost;
        }
        return null;
    case CLOSE_ROOT_MODAL:
        return null;
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
    rootModalVisible,
    postID,
    selectedPost,
    issues,
    inIssues,
    outIssues,
    rhsState,
    rhsPluginAction,
    isRhsVisible,
    isTeamSidebarHidden,
});
