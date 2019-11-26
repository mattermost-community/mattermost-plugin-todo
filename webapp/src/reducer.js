import {combineReducers} from 'redux';

import {OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL, GET_ITEMS, RECEIVED_SHOW_RHS_ACTION} from './action_types';

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

const items = (state = '', action) => {
    switch (action.type) {
    case GET_ITEMS:
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

export default combineReducers({
    rootModalVisible,
    postID,
    items,
    rhsPluginAction,
});

