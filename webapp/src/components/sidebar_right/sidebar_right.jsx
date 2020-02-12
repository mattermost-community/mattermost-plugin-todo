// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import Scrollbars from 'react-custom-scrollbars';

import ToDoIssues from './todo_issues';

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

const MyListName = 'my';
const OutListName = 'out';
const InListName = 'in';

export default class SidebarRight extends React.PureComponent {
    static propTypes = {
        todos: PropTypes.arrayOf(PropTypes.object),
        inTodos: PropTypes.arrayOf(PropTypes.object),
        outTodos: PropTypes.arrayOf(PropTypes.object),
        theme: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            remove: PropTypes.func.isRequired,
            complete: PropTypes.func.isRequired,
            accept: PropTypes.func.isRequired,
            bump: PropTypes.func.isRequired,
            list: PropTypes.func.isRequired,
            openRootModal: PropTypes.func.isRequired,
        }).isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            list: MyListName,
        };
    }

    openList(listName) {
        if (this.state.list !== listName) {
            this.setState({list: listName});
        }
    }

    componentDidMount() {
        this.props.actions.list(false, 'my');
        this.props.actions.list(false, 'in');
        this.props.actions.list(false, 'out');
    }

    getInIssues() {
        return this.props.inTodos.length;
    }

    getOutIssues() {
        return this.props.outTodos.length;
    }

    getMyIssues() {
        return this.props.todos.length;
    }

    render() {
        let todos = [];
        let addButton = '';
        switch (this.state.list) {
        case MyListName:
            todos = this.props.todos || [];
            addButton = 'Add new To-do';
            break;
        case OutListName:
            todos = this.props.outTodos || [];
            addButton = 'Request a To-do from someone';
            break;
        case InListName:
            todos = this.props.inTodos || [];
            addButton = 'Add new To-do';
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
                            className={this.state.list === MyListName ? 'selected' : ''}
                            onClick={() => this.openList(MyListName)}
                        >
                            {'My Todos'} {this.getMyIssues() > 0 ? ' (' + this.getMyIssues() + ')' : ''}
                        </div>
                        <div
                            className={this.state.list === InListName ? 'selected' : ''}
                            onClick={() => this.openList(InListName)}
                        >
                            {'Inbox'} {this.getInIssues() > 0 ? ' (' + this.getInIssues() + ')' : ''}
                        </div>
                        <div
                            className={this.state.list === OutListName ? 'selected' : ''}
                            onClick={() => this.openList(OutListName)}
                        >
                            {'Sent'} {this.getOutIssues() > 0 ? ' (' + this.getOutIssues() + ')' : ''}
                        </div>
                    </div>
                    <div
                        className='section-header'
                        onClick={() => this.props.actions.openRootModal('')}
                    >
                        {addButton + ' '}
                        <i className='icon fa fa-plus-circle'/>
                    </div>
                    <div>
                        <ToDoIssues
                            issues={todos}
                            theme={this.props.theme}
                            list={this.state.list}
                            remove={this.props.actions.remove}
                            complete={this.props.actions.complete}
                            accept={this.props.actions.accept}
                            bump={this.props.actions.bump}
                        />
                    </div>
                </Scrollbars>
            </React.Fragment>
        );
    }
}
