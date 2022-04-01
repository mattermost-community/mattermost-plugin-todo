// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';

import './todo_toast.scss';
import { CSSTransition } from 'react-transition-group';

import { generateClassName } from '../../utils';
import CompassIcon from '../../components/icons/compassIcons';
import IconButton from '../iconButton/iconButton';

type Props = {
    close: () => void,
    removeLastTodo: () => void,
    children?: React.ReactNode,
    submit: PropTypes.func.isRequired,
    title?: string,
    content: {
        icon: string,
        message: string,
    },
    lastPost: {
        id: string,
        create_at: number,
        post_id: string,
        username: string,
        list: string,
        position: number,
        message: string,
        description: string,
    }
    className?: string,
}

function TodoToast(props: Props): JSX.Element {
    const { close, removeLastTodo, submit, title, content, className, lastPost } = props;
    const classNames: Record<string, boolean> = {
        TodoToast: true,
    };
    classNames[`${className}`] = Boolean(className);

    useEffect(() => {
        const timer = setTimeout(() => {
            closeToast();
        }, 5000);
        return () => clearTimeout(timer);
    }, []);

    const closeToast = () => {
        close();
        removeLastTodo();
    };

    const undoTodo = () => {
        submit(lastPost.message, lastPost.description, lastPost.username, lastPost.id);
        closeToast();
    };

    return (
        <CSSTransition
            in={Boolean(content)}
            classNames='slide'
            mountOnEnter={true}
            unmountOnExit={true}
            timeout={300}
            appear={true}
        >
            <div
                className={generateClassName(classNames)}
                title={title}
            >
                <div>
                    <CompassIcon icon={content.icon}/>
                    <span>{content.message}</span>
                    <button
                        onClick={undoTodo}
                        className='TodoToast__undo'
                    >{'Undo'}</button>
                </div>
                <IconButton
                    onClick={closeToast}
                    icon={<CompassIcon icon='close'/>}
                    size='small'
                    inverted={true}
                />
            </div>
        </CSSTransition>
    );
}

export default React.memo(TodoToast);
