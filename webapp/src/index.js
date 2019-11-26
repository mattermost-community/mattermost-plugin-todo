import React from 'react';

import {id as pluginId} from './manifest';

import Root from './components/root';
import SidebarRight from './components/sidebar_right';

import {openRootModal, list, setShowRHSAction} from './actions';
import reducer from './reducer';

let activityFunc;
let lastActivityTime = Number.MAX_SAFE_INTEGER;
const activityTimeout = 60 * 60 * 1000; // 1 hour

export default class Plugin {
    initialize(registry, store) {
        registry.registerRootComponent(Root);
        registry.registerReducer(reducer);

        registry.registerPostDropdownMenuAction(
            'Add To Do',
            (postID) => store.dispatch(openRootModal(postID)),
        );

        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(SidebarRight, 'To Do List');
        store.dispatch(setShowRHSAction(() => store.dispatch(showRHSPlugin)));

        registry.registerChannelHeaderButtonAction(<i className='icon fa fa-list'/>, () => store.dispatch(showRHSPlugin), "To Do", "Open your list of to do items.")

        registry.registerWebSocketEventHandler(`custom_${pluginId}_refresh`, () => store.dispatch(list()));

        store.dispatch(list(true));

        activityFunc = () => {
            const now = new Date().getTime();
            if (now - lastActivityTime > activityTimeout) {
                store.dispatch(list(true));
            }
            lastActivityTime = now;
        };

        document.addEventListener('click', activityFunc);
    }

    deinitialize() {
        document.removeEventListener('click', activityFunc);
    }
}


window.registerPlugin(pluginId, new Plugin());
