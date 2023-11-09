// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import PropTypes from 'prop-types';
import Scrollbars from 'react-custom-scrollbars';
import {Tooltip, OverlayTrigger} from 'react-bootstrap';

import AddIssue from '../add_issue';
import Button from '../../widget/buttons/button';
import TodoToast from '../../widget/todo_toast';
import CompassIcon from '../icons/compassIcons';

import Menu from '../../widget/menu';
import MenuItem from '../../widget/menuItem';
import MenuWrapper from '../../widget/menuWrapper';

import ToDoIssues from '../todo_issues';
import {isKeyPressed} from '../../utils.js';
import Constants from '../../constants';

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
        myIssues: PropTypes.array.isRequired,
        inIssues: PropTypes.array.isRequired,
        outIssues: PropTypes.array.isRequired,
        todoToast: PropTypes.object,
        theme: PropTypes.object.isRequired,
        siteURL: PropTypes.string.isRequired,
        rhsState: PropTypes.string,
        actions: PropTypes.shape({
            remove: PropTypes.func.isRequired,
            complete: PropTypes.func.isRequired,
            accept: PropTypes.func.isRequired,
            bump: PropTypes.func.isRequired,
            fetchAllIssueLists: PropTypes.func.isRequired,
            openAddCard: PropTypes.func.isRequired,
            closeAddCard: PropTypes.func.isRequired,
            openAssigneeModal: PropTypes.func.isRequired,
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
    }

    openList(listName) {
        if (this.state.list !== listName) {
            this.setState({list: listName});
        }
    }

    toggleInbox() {
        this.props.actions.telemetry('toggle_inbox', {action: this.state.showInbox ? 'collapse' : 'expand'});
        this.setState({showInbox: !this.state.showInbox});
    }

    toggleMy() {
        this.props.actions.telemetry('toggle_my', {action: this.state.showMy ? 'collapse' : 'expand'});
        this.setState({showMy: !this.state.showMy});
    }

    componentDidMount() {
        document.addEventListener('keydown', this.handleKeypress);
        this.props.actions.fetchAllIssueLists();
        this.props.actions.setVisible(true);
    }

    componentWillUnmount() {
        document.removeEventListener('keydown', this.handleKeypress);
        this.props.actions.setVisible(false);
    }

    handleKeypress = (e) => {
        if (e.altKey && isKeyPressed(e, Constants.KeyCodes.A)) {
            e.preventDefault();
            this.props.actions.openAddCard('');
        }
    };

    componentDidUpdate(prevProps) {
        if (prevProps.rhsState !== this.props.rhsState) {
            this.openList(this.props.rhsState);
        }
    }

    addTodoItem() {
        this.props.actions.openAddCard('');
    }

    closeAddBox = () => {
        this.props.actions.closeAddCard();
    }

    render() {
        const style = getStyle();
        let todos = [];
        let listHeading = 'My Todos';
        let addButton = '';
        let inboxList = [];

        switch (this.state.list) {
        case MyListName:
            todos = this.props.myIssues;
            addButton = 'Add Todo';
            inboxList = this.props.inIssues;
            break;
        case OutListName:
            todos = this.props.outIssues;
            listHeading = 'Sent Todos';
            addButton = 'Request a Todo from someone';
            break;
        }

        let inbox;

        if (inboxList.length > 0) {
            const actionName = this.state.showInbox ? (
                <CompassIcon
                    style={style.todoHeaderIcon}
                    icon='chevron-down'
                />
            ) : (
                <CompassIcon
                    style={style.todoHeaderIcon}
                    icon='chevron-right'
                />
            );
            inbox = (
                <div>
                    <div
                        className='todo-separator'
                        onClick={() => this.toggleInbox()}
                    >
                        {actionName}
                        <div>{`Incoming Todos (${inboxList.length})`}</div>
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
            const actionName = this.state.showMy ? (
                <CompassIcon
                    style={style.todoHeaderIcon}
                    icon='chevron-down'
                />
            ) : (
                <CompassIcon
                    style={style.todoHeaderIcon}
                    icon='chevron-right'
                />
            );
            separator = (
                <div
                    className='todo-separator'
                    onClick={() => this.toggleMy()}
                >
                    {actionName}
                    {`My Todos (${todos.length})`}
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
                    <div className='todolist-header'>
                        <MenuWrapper>
                            <button
                                className='todolist-header__dropdown'
                            >
                                {listHeading}
                                <CompassIcon
                                    style={style.todoHeaderIcon}
                                    icon='chevron-down'
                                />
                            </button>
                            <Menu position='right'>
                                <MenuItem
                                    onClick={() => this.openList(MyListName)}
                                    action={() => this.openList(MyListName)}
                                    text={'My Todos'}
                                />
                                <MenuItem
                                    action={() => this.openList(OutListName)}
                                    text={'Sent Todos'}
                                />
                            </Menu>
                        </MenuWrapper>
                        {this.state.list === MyListName && (
                            <OverlayTrigger
                                id='addOverlay'
                                placement={'bottom'}
                                overlay={(
                                    <Tooltip
                                        id='addTooltip'
                                    >
                                        <div className='shortcut-line'>
                                            <mark className='shortcut-key shortcut-key--tooltip'>{'OPT'}</mark>
                                            <mark className='shortcut-key shortcut-key--tooltip'>{'A'}</mark>
                                        </div>
                                    </Tooltip>
                                )}
                            >
                                <div>
                                    <Button
                                        emphasis='primary'
                                        icon={<CompassIcon icon='plus'/>}
                                        size='small'
                                        onClick={() => {
                                            this.props.actions.telemetry('rhs_add', {
                                                list: this.state.list,
                                            });
                                            this.addTodoItem();
                                        }}
                                    >
                                        {addButton}
                                    </Button>
                                </div>
                            </OverlayTrigger>
                        )}
                    </div>
                    <div>
                        {inbox}
                        {separator}
                        <AddIssue
                            theme={this.props.theme}
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
                    {this.props.todoToast && (
                        <TodoToast/>
                    )}
                </Scrollbars>
            </React.Fragment>
        );
    }
}

const getStyle = () => {
    return {
        todoHeaderIcon: {
            fontSize: 18,
            marginLeft: 2,
        },
    };
};
