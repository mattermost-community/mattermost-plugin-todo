// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {
    makeStyleFromTheme,
    changeOpacity,
} from 'mattermost-redux/utils/theme_utils';

import TodoItem from '../todo_item';
import Tada from '../../illustrations/tada';

function ToDoIssues(props) {
    const style = getStyle(props.theme);
    const {theme, siteURL, accept, complete, list, remove, bump, addVisible, issues} = props;

    let emptyState = (
        <div style={style.completed.container}>
            <Tada/>
            <h3 style={style.completed.title}>{'All tasks completed'}</h3>
            <p style={style.completed.subtitle}>
                {'Nicely done, you’ve finished all of your tasks! Why not reward yourself with a little break.'}
            </p>
        </div>
    );

    if (addVisible) {
        emptyState = null;
    }

    if (!issues.length) {
        return emptyState;
    }
        issues.map((issue) =>
            (
                <TodoItem
                    issue={issue}
                    theme={theme}
                    siteURL={siteURL}
                    accept={accept}
                    complete={complete}
                    list={list}
                    remove={remove}
                    bump={bump}
                    key={issue.id}
                />
            ),
        )
    ) : (
        emptyState
    );
}

ToDoIssues.propTypes = {
    addVisible: PropTypes.bool.isRequired,
    remove: PropTypes.func.isRequired,
    issues: PropTypes.arrayOf(PropTypes.object),
    theme: PropTypes.object.isRequired,
    siteURL: PropTypes.string.isRequired,
    complete: PropTypes.func.isRequired,
    accept: PropTypes.func.isRequired,
    bump: PropTypes.func.isRequired,
    list: PropTypes.string.isRequired,
};

const getStyle = makeStyleFromTheme((theme) => {
    return {
        container: {
            padding: '8px 20px',
            display: 'flex',
            alignItems: 'flex-start',
        },
        completed: {
            container: {
                textAlign: 'center',
                padding: '116px 40px',
            },
            title: {
                fontSize: 20,
                fontWeight: 600,
            },
            subtitle: {
                fontSize: 14,
                color: changeOpacity(theme.centerChannelColor, 0.72),
            },
        },
        itemContent: {
            padding: '0 0 0 16px',
        },
        issueTitle: {
            color: theme.centerChannelColor,
            lineHeight: 1.7,
            fontWeight: 'bold',
        },
        subtitle: {
            fontSize: '13px',
        },
        message: {
            width: '100%',
            overflowWrap: 'break-word',
            whiteSpace: 'pre-wrap',
        },
    };
});

export default ToDoIssues;
