import {Client4} from 'mattermost-redux/client';
import * as TeamSelector from 'mattermost-redux/selectors/entities/teams';
import * as UserActions from 'mattermost-redux/actions/users';

import {
    OPEN_ASSIGNEE_MODAL,
    CLOSE_ASSIGNEE_MODAL,
    OPEN_TODO_TOAST,
    CLOSE_TODO_TOAST,
    RECEIVED_SHOW_RHS_ACTION,
    GET_ISSUES,
    GET_IN_ISSUES,
    GET_OUT_ISSUES,
    UPDATE_RHS_STATE,
    SET_RHS_VISIBLE,
    SET_HIDE_TEAM_SIDEBAR_BUTTONS,
    GET_ASSIGNEE,
    REMOVE_ASSIGNEE,
    OPEN_ADD_CARD,
    CLOSE_ADD_CARD,
    SET_EDITING_TODO,
    REMOVE_EDITING_TODO,
    GET_ALL_ISSUES,
} from './action_types';

import {getPluginServerRoute} from './selectors';

export const openAddCard = (postID) => (dispatch) => {
    dispatch({
        type: OPEN_ADD_CARD,
        postID,
    });
};

export const closeAddCard = () => (dispatch) => {
    dispatch({
        type: CLOSE_ADD_CARD,
    });
};

export const openTodoToast = (message) => (dispatch) => {
    dispatch({
        type: OPEN_TODO_TOAST,
        message,
    });
};

export const closeTodoToast = () => (dispatch) => {
    dispatch({
        type: CLOSE_TODO_TOAST,
    });
};

export const getAssignee = (assignee) => (dispatch) => {
    dispatch({
        type: GET_ASSIGNEE,
        assignee,
    });
};

export const removeAssignee = () => (dispatch) => {
    dispatch({
        type: REMOVE_ASSIGNEE,
    });
};

export const setEditingTodo = (issueID) => (dispatch) => {
    dispatch({
        type: SET_EDITING_TODO,
        issueID,
    });
};

export const removeEditingTodo = () => (dispatch) => {
    dispatch({
        type: REMOVE_EDITING_TODO,
    });
};

export const openAssigneeModal = () => (dispatch) => {
    dispatch({
        type: OPEN_ASSIGNEE_MODAL,
    });
};

export const closeAssigneeModal = () => (dispatch) => {
    dispatch({
        type: CLOSE_ASSIGNEE_MODAL,
    });
};

/**
 * Stores`showRHSPlugin` action returned by
 * registerRightHandSidebarComponent in plugin initialization.
 */
export function setShowRHSAction(showRHSPluginAction) {
    return {
        type: RECEIVED_SHOW_RHS_ACTION,
        showRHSPluginAction,
    };
}

export function setRhsVisible(payload) {
    return {
        type: SET_RHS_VISIBLE,
        payload,
    };
}

export function updateRhsState(rhsState) {
    return {
        type: UPDATE_RHS_STATE,
        state: rhsState,
    };
}

export const telemetry = (event, properties) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/telemetry', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({event, properties}),
    }));
};

export const add = (message, description, sendTo, postID) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/add', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({send_to: sendTo, message, description, post_id: postID}),
    }));
};

export const editIssue = (postID, message, description) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/edit', Client4.getOptions({
        method: 'put',
        body: JSON.stringify({id: postID, message, description}),
    }));
};

export const changeAssignee = (id, assignee) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/change_assignment', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({id, send_to: assignee}),
    }));
};

export const list = (reminder = false, listName = 'my') => async (dispatch, getState) => {
    let resp;
    let data;
    try {
        resp = await fetch(getPluginServerRoute(getState()) + '/list?reminder=' + reminder + '&list=' + listName, Client4.getOptions({
            method: 'get',
        }));
        data = await resp.json();
    } catch (error) {
        return {error};
    }

    let actionType = GET_ISSUES;
    switch (listName) {
    case 'my':
        actionType = GET_ISSUES;
        break;
    case 'in':
        actionType = GET_IN_ISSUES;
        break;
    case 'out':
        actionType = GET_OUT_ISSUES;
        break;
    }

    dispatch({
        type: actionType,
        data,
    });

    return {data};
};

export const fetchAllIssue = () => async (dispatch, getState) => {
    let data;
    try {
        const resp = await fetch(getPluginServerRoute(getState()) + '/lists', Client4.getOptions({
            method: 'get',
        }));
        data = await resp.json();
    } catch (error) {
        return {error};
    }

    dispatch({
        type: GET_ALL_ISSUES,
        data,
    });

    return {data};
};

export const remove = (id) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/remove', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({id}),
    }));
};

export const complete = (id) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/complete', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({id}),
    }));
};

export const accept = (id) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/accept', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({id}),
    }));
};

export const bump = (id) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/bump', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({id}),
    }));
};

export function autocompleteUsers(username) {
    return async (doDispatch, getState) => {
        const team = TeamSelector.getCurrentTeam(getState());
        const {data} = await doDispatch(UserActions.autocompleteUsers(username, team.id));
        return data.users.filter((user) => user.delete_at === 0);
    };
}

export function setHideTeamSidebar(payload) {
    return {
        type: SET_HIDE_TEAM_SIDEBAR_BUTTONS,
        payload,
    };
}

export const updateConfig = () => async (dispatch, getState) => {
    let data;
    try {
        const resp = await fetch(getPluginServerRoute(getState()) + '/config', Client4.getOptions({
            method: 'get',
        }));
        data = await resp.json();
    } catch (error) {
        return {error};
    }

    dispatch(setHideTeamSidebar(data.hide_team_sidebar));

    return {data};
};
