// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import DeleteButton from '../buttons/delete'
import CompleteButton from '../buttons/complete'
import EnqueueButton from '../buttons/enqueue'
import {canComplete, canDelete, canEnqueue} from '../../utils'

const PostUtils = window.PostUtils; // import the post utilities

const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

function ToDoItems(props) {
    const style = getStyle(props.theme);

    return props.items.length > 0 ? props.items.map((item) => {
        const date = new Date(item.create_at);
        const year = date.getFullYear();
        const month = MONTHS[date.getMonth()];
        const day = date.getDate();
        const hours = date.getHours();
        const minutes = '0' + date.getMinutes();
        const seconds = '0' + date.getSeconds();
        const formattedTime = hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
        const formattedDate = month + ' ' + day + ', ' + year;

        const htmlFormattedText = PostUtils.formatText(item.message);
        const itemComponent = PostUtils.messageHtmlToComponent(htmlFormattedText);

        let createdMessage = 'Created ';
        let orderMessage = '';
        if (item.user) {
            if (item.list === '') {
                createdMessage = 'Sent to ' + item.user;
                orderMessage = 'ENQUEUED on position ' + item.position + '.';
            } else if (item.list === 'in') {
                createdMessage = 'Sent to ' + item.user;
                orderMessage = 'In Inbox on position ' + item.position + '.';
            } else if (item.list === 'out') {
                createdMessage = 'Received from ' + item.user;
                orderMessage = '';
            }
        }

        const orderDiv = (
            <div
                className='light'
                style={style.subtitle}
            >
                {orderMessage}
            </div>
        );

        const deleteButton = (
            <DeleteButton
                itemId={item.id}
                remove={props.remove}
                list={props.list}
            />
        );

        const enqueueButton = (
            <EnqueueButton
                itemId={item.id}
                enqueue={props.enqueue}
            />
        );

        const completeButton = (
            <CompleteButton
                itemId={item.id}
                complete={props.complete}
            />
        );

        const actionButtons = (<div className='action-buttons'>
            {canDelete(props.list, item.list) && deleteButton}
            {canEnqueue(props.list) && enqueueButton}
            {canComplete(props.list) && completeButton}
        </div>);

        return (
            <div
                key={item.id}
                style={style.container}
            >
                <div style={style.message}>
                    {itemComponent}
                </div>
                {(canDelete(props.list, item.list) || canComplete(props.list) || canEnqueue(props.list)) && actionButtons}
                <div
                    className='light'
                    style={style.subtitle}
                >
                    {createdMessage + ' on ' + formattedDate + ' at ' + formattedTime}
                </div>
                {orderMessage && orderDiv}
            </div>
        );
    }) : <div style={style.container}>{'You have no to do items'}</div>;
}

ToDoItems.propTypes = {
    items: PropTypes.array.isRequired,
    theme: PropTypes.object.isRequired,
    list: PropTypes.string.isRequired,
    remove: PropTypes.func.isRequired,
    complete: PropTypes.func.isRequired,
    enqueue: PropTypes.func.isRequired,
};

const getStyle = makeStyleFromTheme((theme) => {
    return {
        container: {
            padding: '15px',
            borderTop: `1px solid ${changeOpacity(theme.centerChannelColor, 0.2)}`,
        },
        itemTitle: {
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

export default ToDoItems;
