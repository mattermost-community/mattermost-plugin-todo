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
export const getIssues = (state) => getPluginState(state).issues;
export const getInIssues = (state) => getPluginState(state).inIssues;
export const getOutIssues = (state) => getPluginState(state).outIssues;
