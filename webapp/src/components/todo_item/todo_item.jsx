import React, { useState } from 'react';
import PropTypes from 'prop-types';

import { makeStyleFromTheme } from 'mattermost-redux/utils/theme_utils';

import CompleteButton from '../buttons/complete';
import AcceptButton from '../buttons/accept';
import BumpButton from '../buttons/bump';
import {
    canComplete,
    canRemove,
    canAccept,
    canBump,
    handleFormattedTextClick,
} from '../../utils';
import CompassIcon from '../icons/compassIcons';
import Menu from '../../widget/menu';
import MenuItem from '../../widget/menuItem';
import MenuWrapper from '../../widget/menuWrapper';

const PostUtils = window.PostUtils; // import the post utilities

function TodoItem(props) {
    const { issue, theme, siteURL, accept, complete, list, remove, bump, actions } = props;

    const [editTodo, setEditTodo] = useState(false);

    const style = getStyle(theme);
    const handleClick = (e) => handleFormattedTextClick(e);

    const htmlFormattedText = PostUtils.formatText(issue.message, {
        siteURL,
    });
    const issueComponent = PostUtils.messageHtmlToComponent(htmlFormattedText);

    let createdMessage = 'Created ';
    let listPositionMessage = '';
    if (issue.user) {
        if (issue.list === '') {
            createdMessage = 'Sent to ' + issue.user;
            listPositionMessage =
                'Accepted. On position ' + (issue.position + 1) + '.';
        } else if (issue.list === 'in') {
            createdMessage = 'Sent to ' + issue.user;
            listPositionMessage =
                'In Inbox on position ' + (issue.position + 1) + '.';
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

    const acceptButton = (
        <AcceptButton
            issueId={issue.id}
            accept={accept}
        />
    );

    const completeButton = (
        <CompleteButton
            theme={theme}
            issueId={issue.id}
            complete={complete}
        />
    );

    const bumpButton = (
        <BumpButton
            issueId={issue.id}
            bump={bump}
        />
    );

    const actionButtons = (
        <div className='action-buttons'>
            {canAccept(list) && acceptButton}
            {canBump(list, issue.list) && bumpButton}
        </div>
    );

    return (
        <div
            key={issue.id}
            className='todo-item'
        >
            <div className='todo-item__content'>
                {canComplete(list) && completeButton}
                <div style={style.itemContent}>
                    {editTodo ? (
                        <textarea value={issueComponent}/>
                    ) : (
                        <div
                            className='todo-text'
                            onClick={handleClick}
                        >
                            {issueComponent}
                        </div>
                    )}

                    {(canRemove(list, issue.list) ||
                        canComplete(list) ||
                        canAccept(list) ||
                        canBump(list, issue.list)) &&
                        actionButtons}
                    {listPositionMessage && listDiv}
                </div>
            </div>
            <MenuWrapper>
                <button className='todo-item__dots'>
                    <CompassIcon icon='dots-vertical'/>
                </button>
                <Menu position='left'>
                    {canAccept(list) && (
                        <MenuItem
                            action={() => accept(issue.id)}
                            text='Accept todo'
                            icon='check'
                        />
                    )}
                    <MenuItem
                        text='Edit todo'
                        icon='pencil-outline'
                        action={() => setEditTodo(true)}
                    />
                    <MenuItem
                        text='Assign to…'
                        icon='account-plus-outline'
                        action={() => actions.openAssigneeModal('')}
                    />
                    {canRemove(list, issue.list) && (
                        <MenuItem
                            action={() => remove(issue.id)}
                            text='Delete todo'
                            icon='trash-can-outline'
                        />
                    )}
                </Menu>
            </MenuWrapper>
        </div>
    );
}

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

TodoItem.propTypes = {
    remove: PropTypes.func.isRequired,
    issue: PropTypes.arrayOf(PropTypes.object),
    theme: PropTypes.object.isRequired,
    siteURL: PropTypes.string.isRequired,
    complete: PropTypes.func.isRequired,
    accept: PropTypes.func.isRequired,
    bump: PropTypes.func.isRequired,
    list: PropTypes.func.isRequired,
    actions: PropTypes.shape({
        openAssigneeModal: PropTypes.func.isRequired,
    }).isRequired,
};

export default TodoItem;
