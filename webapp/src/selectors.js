import {getConfig} from 'mattermost-redux/selectors/entities/general';

import manifest from './manifest';

const {id: pluginId} = manifest;
const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;
export const isAddCardVisible = (state) => getPluginState(state).addCardVisible;
export const isAssigneeModalVisible = (state) => getPluginState(state).assigneeModalVisible;
export const subMenu = (state) => getPluginState(state).subMenu;
export const getPostID = (state) => getPluginState(state).postID;
export const getAssignee = (state) => getPluginState(state).currentAssignee;
export const getEditingTodo = (state) => getPluginState(state).editingTodo;
export const getTodoToast = (state) => getPluginState(state).todoToast;
export const getShowRHSAction = (state) => getPluginState(state).rhsPluginAction;
export const getMessage = (state) => {
    const postID = getPluginState(state).postID;
    if (!postID) {
        return '';
    }
    const post = state.entities.posts.posts[postID];
    if (!post) {
        return '';
    }
    return post.message;
};
export const getMyIssues = (state) => getAllIssues(state).my;
export const getInIssues = (state) => getAllIssues(state).in;
export const getOutIssues = (state) => getAllIssues(state).out;
export const getAllIssues = (state) => getPluginState(state).allIssues;
export const getCurrentTeamRoute = (state) => {
    const basePath = getSiteURL(state);
    const teamName = state.entities.teams.teams[state.entities.teams.currentTeamId].name;

    return basePath + '/' + teamName + '/';
};

// TODO: Move this into mattermost-redux or mattermost-webapp.
export const getSiteURL = (state) => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath;
};

export const getPluginServerRoute = (state) => {
    const siteURL = getSiteURL(state);
    return siteURL + '/plugins/' + pluginId;
};

export const isRhsVisible = (state) => getPluginState(state).isRhsVisible;
export const isButtonSidebarVisible = (state) => !getPluginState(state).isButtonSidebarHidden;
