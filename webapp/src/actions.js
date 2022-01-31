import {Client4} from 'mattermost-redux/client';
import * as UserActions from 'mattermost-redux/actions/users';

import {
    OPEN_ROOT_MODAL,
    CLOSE_ROOT_MODAL,
    RECEIVED_SHOW_RHS_ACTION,
    GET_ISSUES,
    GET_IN_ISSUES,
    GET_OUT_ISSUES,
    UPDATE_RHS_STATE,
    SET_RHS_VISIBLE,
    SET_HIDE_TEAM_SIDEBAR_BUTTONS,
} from './action_types';

import {getPluginServerRoute} from './selectors';

export const openRootModal = (postID) => (dispatch) => {
    dispatch({
        type: OPEN_ROOT_MODAL,
        postID,
    });
};

export const closeRootModal = () => (dispatch) => {
    dispatch({
        type: CLOSE_ROOT_MODAL,
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

export const add = (message, sendTo, postID) => async (dispatch, getState) => {
    await fetch(getPluginServerRoute(getState()) + '/add', Client4.getOptions({
        method: 'post',
        body: JSON.stringify({message, send_to: sendTo, post_id: postID}),
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
    return async (doDispatch) => {
        const {data} = await doDispatch(UserActions.autocompleteUsers(username));
        return data.users;
    };
}

export function setHideTeamSidebar(payload) {
    return {
        type: SET_HIDE_TEAM_SIDEBAR_BUTTONS,
        payload,
    };
}

export const updateConfig = () => async (dispatch, getState) => {
    let resp;
    let data;
    try {
        resp = await fetch(getPluginServerRoute(getState()) + '/config', Client4.getOptions({
            method: 'get',
        }));
        data = await resp.json();
    } catch (error) {
        return {error};
    }

    dispatch(setHideTeamSidebar(data.hide_team_sidebar));

    return {data};
};
