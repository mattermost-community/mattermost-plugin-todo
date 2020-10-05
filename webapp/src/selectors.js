import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;
export const getPostID = (state) => getPluginState(state).postID;
export const getSelectedPost = (state) => getPluginState(state).selectedPost;
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
export const getIssues = (state) => getPluginState(state).issues;
export const getInIssues = (state) => getPluginState(state).inIssues;
export const getOutIssues = (state) => getPluginState(state).outIssues;
export const getCurrentTeamRoute = (state) => {
    const basePath = getSiteURL();
    const teamName = state.entities.teams.teams[state.entities.teams.currentTeamId].name;

    return basePath + '/' + teamName + '/';
};

export const getSiteURL = () => {
    let siteURL = window.location.protocol + '//' + window.location.hostname + (window.location.port ? ':' + window.location.port : '');
    if (window.location.origin) {
        siteURL = window.location.origin;
    }

    if (siteURL[siteURL.length - 1] === '/') {
        siteURL = siteURL.substring(0, siteURL.length - 1);
    }

    if (window.basename) {
        siteURL += window.basename;
    }

    if (siteURL[siteURL.length - 1] === '/') {
        siteURL = siteURL.substring(0, siteURL.length - 1);
    }

    return siteURL;
};

export const isRhsVisible = (state) => getPluginState(state).isRhsVisible;
export const isTeamSidebarVisible = (state) => !getPluginState(state).isTeamSidebarHidden;
