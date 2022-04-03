// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react';

import './chip.scss';
import {generateClassName} from '../../utils';

type Props = {
    onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void
    onBlur?: (e: React.FocusEvent<HTMLButtonElement>) => void
    children?: React.ReactNode
    title?: string
    icon?: React.ReactNode
    size?: string
    className?: string
    rightIcon?: boolean
}

function Chip(props: Props): JSX.Element {
    const classNames: Record<string, boolean> = {
        Chip: true,
    };

    classNames[`size--${props.size}`] = Boolean(props.size);
    classNames[`${props.className}`] = Boolean(props.className);

    return (
        <button
            type={'button'}
            onClick={props.onClick}
            className={generateClassName(classNames)}
            title={props.title}
            onBlur={props.onBlur}
        >
            {!props.rightIcon && props.icon}
            <span>{props.children}</span>
            {props.rightIcon && props.icon}
        </button>);
}

export default React.memo(Chip);
