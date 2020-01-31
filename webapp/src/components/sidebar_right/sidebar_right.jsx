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
            enqueue: PropTypes.func.isRequired,
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

    getInItems() {
        return this.props.inTodos.length;
    }

    getOutItems() {
        return this.props.outTodos.length;
    }

    getMyItems() {
        return this.props.todos.length;
    }

    render() {
        let todos = [];
        switch (this.state.list) {
        case MyListName:
            todos = this.props.todos || [];
            break;
        case OutListName:
            todos = this.props.outTodos || [];
            break;
        case InListName:
            todos = this.props.inTodos || [];
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
                            {'My'} {this.getMyItems() > 0 ? ' (' + this.getMyItems() + ')' : ''}
                        </div>
                        <div
                            className={this.state.list === InListName ? 'selected' : ''}
                            onClick={() => this.openList(InListName)}
                        >
                            {'In'} {this.getInItems() > 0 ? ' (' + this.getInItems() + ')' : ''}
                        </div>
                        <div
                            className={this.state.list === OutListName ? 'selected' : ''}
                            onClick={() => this.openList(OutListName)}
                        >
                            {'Out'} {this.getOutItems() > 0 ? ' (' + this.getOutItems() + ')' : ''}
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
                            theme={this.props.theme}
                            list={this.state.list}
                            remove={this.props.actions.remove}
                            complete={this.props.actions.complete}
                            enqueue={this.props.actions.enqueue}
                        />
                    </div>
                </Scrollbars>
            </React.Fragment>
        );
    }
}