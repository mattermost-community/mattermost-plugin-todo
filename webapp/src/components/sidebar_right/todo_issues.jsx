// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import RemoveButton from '../buttons/remove';
import CompleteButton from '../buttons/complete';
import AcceptButton from '../buttons/accept';
import BumpButton from '../buttons/bump';
import {canComplete, canRemove, canAccept, canBump} from '../../utils';

const PostUtils = window.PostUtils; // import the post utilities

const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

function ToDoIssues(props) {
    const style = getStyle(props.theme);

    return props.issues.length > 0 ? props.issues.map((issue) => {
        const date = new Date(issue.create_at);
        const year = date.getFullYear();
        const month = MONTHS[date.getMonth()];
        const day = date.getDate();
        const hours = date.getHours();
        const minutes = '0' + date.getMinutes();
        const seconds = '0' + date.getSeconds();
        const formattedTime = hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
        const formattedDate = month + ' ' + day + ', ' + year;

        const htmlFormattedText = PostUtils.formatText(issue.message, {siteURL: props.siteURL});
        const issueComponent = PostUtils.messageHtmlToComponent(htmlFormattedText);

        let createdMessage = 'Created ';
        let listPositionMessage = '';
        if (issue.user) {
            if (issue.list === '') {
                createdMessage = 'Sent to ' + issue.user;
                listPositionMessage = 'Accepted. On position ' + issue.position + '.';
            } else if (issue.list === 'in') {
                createdMessage = 'Sent to ' + issue.user;
                listPositionMessage = 'In Inbox on position ' + issue.position + '.';
            } else if (issue.list === 'out') {
                createdMessage = 'Received from ' + issue.user;
                listPositionMessage = '';
            }
        }

        const listDiv = (
            <div
                className='light'
                style={style.subtitle}
            >
                {listPositionMessage}
            </div>
        );

        const removeButton = (
            <RemoveButton
                issueId={issue.id}
                remove={props.remove}
                list={props.list}
            />
        );

        const acceptButton = (
            <AcceptButton
                issueId={issue.id}
                accept={props.accept}
            />
        );

        const completeButton = (
            <CompleteButton
                issueId={issue.id}
                complete={props.complete}
            />
        );

        const bumpButton = (
            <BumpButton
                issueId={issue.id}
                bump={props.bump}
            />
        );

        const actionButtons = (<div className='action-buttons'>
            {canRemove(props.list, issue.list) && removeButton}
            {canAccept(props.list) && acceptButton}
            {canComplete(props.list) && completeButton}
            {canBump(props.list, issue.list) && bumpButton}
        </div>);

        return (
            <div
                key={issue.id}
                style={style.container}
            >
                <div style={style.message}>
                    {issueComponent}
                </div>
                {(canRemove(props.list, issue.list) || canComplete(props.list) || canAccept(props.list) || canBump(props.list, issue.list)) && actionButtons}
                <div
                    className='light'
                    style={style.subtitle}
                >
                    {createdMessage + ' on ' + formattedDate + ' at ' + formattedTime}
                </div>
                {listPositionMessage && listDiv}
            </div>
        );
    }) : <div style={style.container}>{'You have no to do issues'}</div>;
}

ToDoIssues.propTypes = {
    issues: PropTypes.array.isRequired,
    theme: PropTypes.object.isRequired,
    list: PropTypes.string.isRequired,
    remove: PropTypes.func.isRequired,
    complete: PropTypes.func.isRequired,
    accept: PropTypes.func.isRequired,
    bump: PropTypes.func.isRequired,
    siteURL: PropTypes.string.isRequired,
};

const getStyle = makeStyleFromTheme((theme) => {
    return {
        container: {
            padding: '15px',
            borderTop: `1px solid ${changeOpacity(theme.centerChannelColor, 0.2)}`,
        },
        issueTitle: {
            color: theme.centerChannelColor,
            lineHeight: 1.7,
            fontWeight: 'bold',
        },
        subtitle: {
            margin: '5px 0 0 0',
            fontSize: '13px',
        },
        message: {
            width: '100%',
        },
    };
});

export default ToDoIssues;
