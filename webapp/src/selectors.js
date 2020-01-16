import {id as pluginId} from './manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const isRootModalVisible = (state) => getPluginState(state).rootModalVisible;
export const getPostID = (state) => getPluginState(state).postID;
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
export const getItems = (state) => getPluginState(state).items;
export const getCurrentTeamRoute = (state) => {
    const basePath = getSiteURL();
    const teamName = state.entities.teams.teams[state.entities.teams.currentTeamId].name;

    return basePath + '/' + teamName + '/';
};

function getSiteURL() {
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
}