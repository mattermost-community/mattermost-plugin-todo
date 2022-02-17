import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../components/icons/compassIcons';

type Props = {
    action: () => void,
    icon: string,
    text?: string,
    shortcut?: string,
}

const MenuItem = (props: Props) => {
    const { icon, shortcut, action, text } = props;

    useEffect(() => {
        function handleKeypress(e: {key: string}) {
            if (e.key === shortcut) {
                action();
            }
        }

        document.addEventListener('keyup', handleKeypress);

        return () => {
            document.removeEventListener('keyup', handleKeypress);
        };
    }, []);

    return (
        <button
            className='menu-option'
            onClick={() => action()}
        >
            <div className='menu-option__left'>
                {icon && (
                    <CompassIcon
                        icon={icon}
                        className='MenuItemIcon'
                    />
                )}
                <span>{text}</span>
            </div>
            <div className='menu-option__shortcut'>{shortcut}</div>
        </button>
    );
};

MenuItem.propTypes = {
    icon: PropTypes.node,
    text: PropTypes.string.isRequired,
    action: PropTypes.func.isRequired,
};

export default MenuItem;
