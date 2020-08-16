import React from 'react';

import {id as pluginId} from './manifest';

import Root from './components/root';
import SidebarRight from './components/sidebar_right';

import {openRootModal, list, setShowRHSAction, updateConfig, setHideTeamSidebar} from './actions';
import reducer from './reducer';
import PostTypeTodo from './components/post_type_todo';
import TeamSidebar from './components/team_sidebar';
import ChannelHeaderButton from './components/channel_header_button';

let activityFunc;
let lastActivityTime = Number.MAX_SAFE_INTEGER;
const activityTimeout = 60 * 60 * 1000; // 1 hour

export default class Plugin {
    initialize(registry, store) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(Root);

        registry.registerBottomTeamSidebarComponent(TeamSidebar);

        registry.registerPostDropdownMenuAction(
            'Add Todo',
            (postID) => store.dispatch(openRootModal(postID)),
        );

        const {toggleRHSPlugin, showRHSPlugin} = registry.registerRightHandSidebarComponent(SidebarRight, 'Todo List');
        store.dispatch(setShowRHSAction(() => store.dispatch(showRHSPlugin)));
        registry.registerChannelHeaderButtonAction(<ChannelHeaderButton/>, () => store.dispatch(toggleRHSPlugin), 'Todo', 'Open your list of Todo issues.');

        const refresh = () => {
            store.dispatch(list(false, 'my'));
            store.dispatch(list(false, 'in'));
            store.dispatch(list(false, 'out'));
        };

        registry.registerWebSocketEventHandler(`custom_${pluginId}_refresh`, refresh);
        registry.registerReconnectHandler(refresh);

        store.dispatch(list(true));
        store.dispatch(list(false, 'in'));
        store.dispatch(list(false, 'out'));

        // register websocket event to track config changes
        const configUpdate = ({data}) => {
            store.dispatch(setHideTeamSidebar(data.hide_team_sidebar));
        };

        registry.registerWebSocketEventHandler(`custom_${pluginId}_config_update`, configUpdate);

        store.dispatch(updateConfig());

        activityFunc = () => {
            const now = new Date().getTime();
            if (now - lastActivityTime > activityTimeout) {
                store.dispatch(list(true));
            }
            lastActivityTime = now;
        };

        document.addEventListener('click', activityFunc);

        registry.registerPostTypeComponent('custom_todo', PostTypeTodo);
    }

    deinitialize() {
        document.removeEventListener('click', activityFunc);
    }
}

window.registerPlugin(pluginId, new Plugin());
