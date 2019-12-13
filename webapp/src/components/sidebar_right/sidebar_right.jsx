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

export default class SidebarRight extends React.PureComponent {
    static propTypes = {
        todos: PropTypes.arrayOf(PropTypes.object),
        theme: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            remove: PropTypes.func.isRequired,
            list: PropTypes.func.isRequired,
            openRootModal: PropTypes.func.isRequired,
        }).isRequired,
    };

    componentDidMount() {
        this.props.actions.list();
    }

    render() {
        const todos = this.props.todos || [];

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