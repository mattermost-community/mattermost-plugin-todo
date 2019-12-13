// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

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

        return (
            <div
                key={item.id}
                style={style.container}
            >
                <a
                    href='#'
                    className='light pull-right'
                    onClick={() => props.remove(item.id)}
                >
                    {'X'}
                </a>
                <div style={style.message}>
                    <strong>
                        {item.message}
                    </strong>
                </div>
                <div
                    className='light'
                    style={style.subtitle}
                >
                    {'Created on ' + formattedDate + ' at ' + formattedTime}
                </div>
            </div>
        );
    }) : <div style={style.container}>{'You have no to do items'}</div>;
}

ToDoItems.propTypes = {
    items: PropTypes.array.isRequired,
    theme: PropTypes.object.isRequired,
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
            width: '95%',
        },
    };
});

export default ToDoItems;
