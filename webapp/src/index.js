import React from 'react';

import {id as pluginId} from './manifest';

import Root from './components/root';
import SidebarRight from './components/sidebar_right';

import {openRootModal, list, setShowRHSAction} from './actions';
import reducer from './reducer';
import PostTypeTodo from './components/post_type_todo';
import TeamSidebar from './components/team_sidebar';

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

        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(SidebarRight, 'Todo List');
        store.dispatch(setShowRHSAction(() => store.dispatch(showRHSPlugin)));

        registry.registerChannelHeaderButtonAction(<i className='icon fa fa-list'/>, () => store.dispatch(showRHSPlugin), 'Todo', 'Open your list of Todo issues.');

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
