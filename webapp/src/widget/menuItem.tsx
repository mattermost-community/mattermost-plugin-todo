import React, {useState} from 'react';
import PropTypes from 'prop-types';
import CompassIcon from '../components/icons/compassIcons'

const MenuItem = (props) => {
    return (
        <button
            className='menu-option'
            onClick={() => props.action()}
        >
            <CompassIcon
                icon={props.icon}
                className='MenuItemIcon'
            />
            <span>{props.text}</span>
        </button>
    );
};

MenuItem.propTypes = {
    icon: PropTypes.string.isRequired,
    text: PropTypes.string.isRequired,
};

export default MenuItem;
