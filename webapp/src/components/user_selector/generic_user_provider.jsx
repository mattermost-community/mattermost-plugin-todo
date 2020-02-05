// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import * as Utils from 'utils';

import Provider from './provider.jsx';
import Suggestion from './suggestion.jsx';

class UserSuggestion extends Suggestion {
    render() {
        const {item, isSelection} = this.props;

        let className = 'suggestion-list__item mentions__name';
        if (isSelection) {
            className += ' suggestion--selected';
        }

        const username = item.username;
        let description = '';

        if ((item.first_name || item.last_name) && item.nickname) {
            description = `- ${Utils.getFullName(item)} (${item.nickname})`;
        } else if (item.nickname) {
            description = `- (${item.nickname})`;
        } else if (item.first_name || item.last_name) {
            description = `- ${Utils.getFullName(item)}`;
        }

        return (
            <div
                className={className}
                onClick={this.handleClick}
                onMouseMove={this.handleMouseMove}
                {...Suggestion.baseProps}
            >
                <span className='admin-setting-user--align'>
                    {'@' + username}
                </span>
                <span className='admin-setting-user__fullname'>
                    {' '}
                    {description}
                </span>
            </div>
        );
    }
}

export default class UserProvider extends Provider {
    constructor(searchUsersFunc) {
        super();
        this.autocompleteUsers = searchUsersFunc;
    }
    async handlePretextChanged(pretext, resultsCallback) {
        const normalizedPretext = pretext.toLowerCase();
        this.startNewRequest(normalizedPretext);

        const data = await this.autocompleteUsers(normalizedPretext);

        if (this.shouldCancelDispatch(normalizedPretext)) {
            return false;
        }

        const users = Object.assign([], data.users);

        resultsCallback({
            matchedPretext: normalizedPretext,
            terms: users.map((user) => user.username),
            items: users,
            component: UserSuggestion,
        });

        return true;
    }
}
