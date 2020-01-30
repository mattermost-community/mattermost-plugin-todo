// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import Scrollbars from 'react-custom-scrollbars';

import ToDoItems from './todo_items';

import './sidebar_right.scss';

export function renderView(props) {
    return (
        <div
            {...props}
            className='scrollbar--view'
        />);
}

export function renderThumbHorizontal(props) {
    return (
        <div
            {...props}
            className='scrollbar--horizontal'
        />);
}

export function renderThumbVertical(props) {
    return (
        <div
            {...props}
            className='scrollbar--vertical'
        />);
}

const OwnListName = 'own';
const SentListName = 'sent';
const InboxListName = 'inbox';

export default class SidebarRight extends React.PureComponent {
    static propTypes = {
        todos: PropTypes.arrayOf(PropTypes.object),
        inboxTodos: PropTypes.arrayOf(PropTypes.object),
        sentTodos: PropTypes.arrayOf(PropTypes.object),
        theme: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            remove: PropTypes.func.isRequired,
            complete: PropTypes.func.isRequired,
            enqueue: PropTypes.func.isRequired,
            list: PropTypes.func.isRequired,
            openRootModal: PropTypes.func.isRequired,
        }).isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            list: OwnListName,
        };
    }

    openList(listName) {
        if (this.state.list !== listName) {
            this.setState({list: listName});
        }
    }

    componentDidMount() {
        this.props.actions.list(false, 'own');
        this.props.actions.list(false, 'inbox');
        this.props.actions.list(false, 'sent');
    }

    getInboxImportantItems() {
        return this.props.todos.length;
    }

    getSentImportantItems() {
        return this.props.sentTodos.length;
    }

    getOwnImportantItems() {
        return this.props.inboxTodos.length;
    }

    render() {
        let todos = [];
        switch (this.state.list) {
        case OwnListName:
            todos = this.props.todos || [];
            break;
        case SentListName:
            todos = this.props.sentTodos || [];
            break;
        case InboxListName:
            todos = this.props.inboxTodos || [];
            break;
        }

        return (
            <React.Fragment>
                <Scrollbars
                    autoHide={true}
                    autoHideTimeout={500}
                    autoHideDuration={500}
                    renderThumbHorizontal={renderThumbHorizontal}
                    renderThumbVertical={renderThumbVertical}
                    renderView={renderView}
                    className='SidebarRight'
                >
                    <div className='header-menu'>
                        <div
                            className={this.state.list === InboxListName ? 'selected' : ''}
                            onClick={() => this.openList(InboxListName)}
                        >
                            {'Inbox'} {this.getInboxImportantItems() > 0 ? ' (' + this.getInboxImportantItems() + ')' : ''}
                        </div>
                        <div
                            className={this.state.list === SentListName ? 'selected' : ''}
                            onClick={() => this.openList(SentListName)}
                        >
                            {'Sent'} {this.getSentImportantItems() > 0 ? ' (' + this.getSentImportantItems() + ')' : ''}
                        </div>
                        <div
                            className={this.state.list === OwnListName ? 'selected' : ''}
                            onClick={() => this.openList(OwnListName)}
                        >
                            {'Own'} {this.getOwnImportantItems() > 0 ? ' (' + this.getOwnImportantItems() + ')' : ''}
                        </div>
                    </div>
                    <div
                        className='section-header'
                        onClick={() => this.props.actions.openRootModal('')}
                    >
                        {'Add new item '}
                        <i className='icon fa fa-plus-circle'/>
                    </div>
                    <div>
                        <ToDoItems
                            items={todos}
                            remove={this.props.actions.remove}
                            theme={this.props.theme}
                        />
                    </div>
                </Scrollbars>
            </React.Fragment>
        );
    }
}