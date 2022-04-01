import React from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../icons/compassIcons';

const CompleteButton = (props) => {
    const markAsDone = () => {
        props.markAsDone();
        setTimeout(() => {
            props.completeToast();
        }, 1000);
    };

    return (
        <button
            className={`todo-item__checkbox ${props.active ? 'todo-item__checkbox--active' : ''}`}
            onClick={() => markAsDone()}
        >
            <CompassIcon
                icon='check'
                className='CheckIcon'
            />
        </button>
    );
};

CompleteButton.propTypes = {
    active: PropTypes.bool,
    issueId: PropTypes.string.isRequired,
    completeToast: PropTypes.func.isRequired,
    markAsDone: PropTypes.func.isRequired,
};

export default CompleteButton;
