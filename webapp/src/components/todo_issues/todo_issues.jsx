// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import { makeStyleFromTheme } from 'mattermost-redux/utils/theme_utils';

import TodoItem from '../todo_item';

function ToDoIssues(props) {
    const style = getStyle(props.theme);
    const { theme, siteURL, accept, complete, list, remove, bump } = props;

    return props.issues.length > 0 ? (
        props.issues.map((issue) =>
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
        <div style={style.container}>{'You have no Todo issues'}</div>
    );
}

ToDoIssues.propTypes = {
    remove: PropTypes.func.isRequired,
    issues: PropTypes.arrayOf(PropTypes.object),
    theme: PropTypes.object.isRequired,
    siteURL: PropTypes.string.isRequired,
    complete: PropTypes.func.isRequired,
    accept: PropTypes.func.isRequired,
    bump: PropTypes.func.isRequired,
    list: PropTypes.func.isRequired,
};

const getStyle = makeStyleFromTheme((theme) => {
    return {
        container: {
            padding: '8px 20px',
            display: 'flex',
            alignItems: 'flex-start',
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
