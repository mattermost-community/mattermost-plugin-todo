import React, { useState } from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../components/icons/compassIcons';

type Props = {
    action: () => void,
    icon: string,
    text?: string,
}

const MenuItem = (props: Props) => {
    return (
        <button
            className='menu-option'
            onClick={() => props.action()}
        >
            {props.icon && (
                <CompassIcon
                    icon={props.icon}
                    className='MenuItemIcon'
                />
            )}
            <span>{props.text}</span>
        </button>
    );
};

MenuItem.propTypes = {
    icon: PropTypes.node,
    text: PropTypes.string.isRequired,
    action: PropTypes.func.isRequired,
};

export default MenuItem;
