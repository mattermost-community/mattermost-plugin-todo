import React, { useState } from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../icons/compassIcons';

const CompleteButton = (props) => {
    const [active, setActive] = useState(false);

    const markAsDone = () => {
        setActive(true);
        props.markAsDone();
        setTimeout(() => {
            props.completeToast();
        }, 1000);
    };

    return (
        <button
            className={`todo-item__checkbox ${active ? 'todo-item__checkbox--active' : ''}`}
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
    issueId: PropTypes.string.isRequired,
    completeToast: PropTypes.func.isRequired,
    markAsDone: PropTypes.func.isRequired,
};

export default CompleteButton;
