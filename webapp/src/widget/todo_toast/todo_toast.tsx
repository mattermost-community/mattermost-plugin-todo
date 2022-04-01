// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React, { useEffect } from 'react';

import './todo_toast.scss';
import { generateClassName } from '../../utils';
import CompassIcon from '../../components/icons/compassIcons';
import IconButton from '../iconButton/iconButton';

type Props = {
    close: () => void,
    children?: React.ReactNode,
    title?: string,
    content: {
        icon: string,
        message: string,
    },
    className?: string,
}

function TodoToast(props: Props): JSX.Element {
    const classNames: Record<string, boolean> = {
        TodoToast: true,
    };
    classNames[`${props.className}`] = Boolean(props.className);

    useEffect(() => {
        const timer = setTimeout(() => {
            props.close();
        }, 3000);
        return () => clearTimeout(timer);
    }, []);

    const closeToast = () => {
        props.close();
    };

    return (
        <div
            className={generateClassName(classNames)}
            title={props.title}
        >
            <div>
                <CompassIcon icon={props.content.icon}/>
                <span>{props.content.message}</span>
            </div>
            <IconButton
                onClick={closeToast}
                icon={<CompassIcon icon='close'/>}
                size='small'
                inverted={true}
            />
        </div>);
}

export default React.memo(TodoToast);
