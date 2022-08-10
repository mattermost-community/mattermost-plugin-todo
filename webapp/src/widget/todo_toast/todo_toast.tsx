// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React, {useEffect, useCallback} from 'react';

import './todo_toast.scss';
import {CSSTransition} from 'react-transition-group';

import {generateClassName} from '../../utils';
import CompassIcon from '../../components/icons/compassIcons';
import IconButton from '../iconButton/iconButton';

type Props = {
    close: () => void,
    children?: React.ReactNode,
    submit: () => void,
    title?: string,
    content: {
        icon: string,
        message: string,
        undo: () => void,
    },
    className?: string,
}

function TodoToast(props: Props): JSX.Element {
    const {close, content} = props;

    const closeToast = useCallback(close, [close]);
    const undoTodo = useCallback(() => {
        content.undo();
        close();
    }, [content.undo, close]);

    const classNames: Record<string, boolean> = {
        TodoToast: true,
    };
    classNames[`${props.className}`] = Boolean(props.className);

    useEffect(() => {
        const timer = setTimeout(() => {
            props.close();
        }, 5000);

        return () => clearTimeout(timer);
    }, []);

    return (
        <CSSTransition
            in={Boolean(props.content)}
            classNames='slide'
            mountOnEnter={true}
            unmountOnExit={true}
            timeout={300}
            appear={true}
        >
            <div
                className={generateClassName(classNames)}
                title={props.title}
            >
                <div>
                    <CompassIcon icon={props.content.icon}/>
                    <span>{props.content.message}</span>
                    <button
                        onClick={undoTodo}
                        className='TodoToast__undo'
                    >
                        {'Undo'}
                    </button>
                </div>
                <IconButton
                    onClick={closeToast}
                    icon={<CompassIcon icon='close'/>}
                    size='small'
                    inverted={'true'}
                />
            </div>
        </CSSTransition>
    );
}

export default React.memo(TodoToast);
