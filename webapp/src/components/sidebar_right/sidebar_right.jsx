// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import Scrollbars from 'react-custom-scrollbars';

import AddIssue from '../add_issue';
import Button from '../../widget/buttons/button';
import CompassIcon from '../icons/compassIcons';

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
        siteURL: PropTypes.string.isRequired,
        rhsState: PropTypes.string,
        actions: PropTypes.shape({
            remove: PropTypes.func.isRequired,
            complete: PropTypes.func.isRequired,
            accept: PropTypes.func.isRequired,
            bump: PropTypes.func.isRequired,
            list: PropTypes.func.isRequired,
            openRootModal: PropTypes.func.isRequired,
            setVisible: PropTypes.func.isRequired,
            telemetry: PropTypes.func.isRequired,
        }).isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            list: props.rhsState || MyListName,
            showInbox: true,
            showMy: true,
            addTodo: false,
        };

        this.closeAddBox = this.closeAddBox.bind(this);
    }

    openList(listName) {
        if (this.state.list !== listName) {
            this.setState({ list: listName });
        }
    }

    toggleInbox() {
        this.props.actions.telemetry('toggle_inbox', { action: this.state.showInbox ? 'collapse' : 'expand' });
        this.setState({ showInbox: !this.state.showInbox });
    }

    toggleMy() {
        this.props.actions.telemetry('toggle_my', { action: this.state.showMy ? 'collapse' : 'expand' });
        this.setState({ showMy: !this.state.showMy });
    }

    componentDidMount() {
        this.props.actions.list(false, 'my');
        this.props.actions.list(false, 'in');
        this.props.actions.list(false, 'out');
        this.props.actions.setVisible(true);
    }

    componentWillUnmount() {
        this.props.actions.setVisible(false);
    }

    componentDidUpdate(prevProps) {
        if (prevProps.rhsState !== this.props.rhsState) {
            this.openList(this.props.rhsState);
        }
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

    addTodoItem() {
        this.setState({ addTodo: true });
    }

    closeAddBox() {
        this.setState({ addTodo: false });
    }

    render() {
        let todos = [];
        let listHeading = 'My Todos';
        let addButton = '';
        let inboxList = [];
        switch (this.state.list) {
        case MyListName:
            todos = this.props.todos || [];
            addButton = 'Add new Todo';
            inboxList = this.props.inTodos || [];
            break;
        case OutListName:
            todos = this.props.outTodos || [];
            listHeading = 'Sent Todos';
            addButton = 'Request a Todo from someone';
            break;
        }

        let inbox;
        if (inboxList.length > 0) {
            const actionName = this.state.showInbox ? 'collapse' : 'expand';
            inbox = (
                <div>
                    <div
                        className='todo-separator'
                        onClick={() => this.toggleInbox()}
                    >
                        {`Incoming Todos (${inboxList.length}) (${actionName})`}
                    </div>
                    {this.state.showInbox ?
                        <ToDoIssues
                            issues={inboxList}
                            theme={this.props.theme}
                            list={InListName}
                            remove={this.props.actions.remove}
                            complete={this.props.actions.complete}
                            accept={this.props.actions.accept}
                            bump={this.props.actions.bump}
                        /> : ''}
                </div>
            );
        }

        let separator;
        if ((inboxList.length > 0) && (todos.length > 0)) {
            const actionName = this.state.showMy ? 'collapse' : 'expand';
            separator = (
                <div
                    className='todo-separator'
                    onClick={() => this.toggleMy()}
                >
                    {`My Todos (${todos.length}) (${actionName})`}
                </div>
            );
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
                    <div className='todolist-tabs'>
                        <div
                            className={'todolist-tab' + (this.state.list === MyListName ? ' selected' : '')}
                            onClick={() => this.openList(MyListName)}
                        >
                            {'Todos'} {this.getMyIssues() > 0 ? ' (' + this.getMyIssues() + ')' : ''} {this.getInIssues() > 0 ? ' (' + this.getInIssues() + ' received)' : ''}
                        </div>
                        <div
                            className={'todolist-tab' + (this.state.list === OutListName ? ' selected' : '')}
                            onClick={() => this.openList(OutListName)}
                        >
                            {'Sent'} {this.getOutIssues() > 0 ? ' (' + this.getOutIssues() + ')' : ''}
                        </div>
                    </div>
                    <div
                        className='todolist-header'
                    >
                        <div className='todolist-header__heading'>{listHeading}</div>
                        <Button
                            emphasis='primary'
                            icon={<CompassIcon icon='plus'/>}
                            size='small'
                            onClick={() => {
                                this.props.actions.telemetry('rhs_add', { list: this.state.list });
                                this.addTodoItem();
                            }}
                        >
                            {addButton}
                        </Button>
                    </div>
                    <div>
                        {inbox}
                        {separator}
                        <AddIssue
                            theme={this.props.theme}
                            visible={this.state.addTodo}
                            closeAddBox={this.closeAddBox}
                        />
                        {(inboxList.length === 0) || (this.state.showMy && todos.length > 0) ?
                            <ToDoIssues
                                issues={todos}
                                theme={this.props.theme}
                                list={this.state.list}
                                remove={this.props.actions.remove}
                                complete={this.props.actions.complete}
                                accept={this.props.actions.accept}
                                bump={this.props.actions.bump}
                                siteURL={this.props.siteURL}
                            /> : ''}
                    </div>
                </Scrollbars>
            </React.Fragment>
        );
    }
}
