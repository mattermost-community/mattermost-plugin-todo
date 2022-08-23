import React, {useEffect} from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../components/icons/compassIcons';
import Constants from 'src/constants';
import {isKeyPressed} from 'src/utils';

export type MenuOptionProps = {
    id: string,
    name: string,
    onClick: (id: string) => void,
}

type Props = {
    action: () => void,
    icon: string,
    text?: string,
    shortcut?: string,
}

const MenuItem = (props: Props) => {
    const {icon, shortcut, action, text} = props;

    useEffect(() => {
        function handleKeypress(e: KeyboardEvent) {
            if (e.key === shortcut) {
                e.preventDefault();
                e.target.dispatchEvent(new Event('menuItemClicked'));
                action();
            }

            if (!isKeyPressed(e, Constants.KeyCodes.TAB)) {
                e.preventDefault();
            }
        }

        document.addEventListener('keydown', handleKeypress);

        return () => {
            document.removeEventListener('keydown', handleKeypress);
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
