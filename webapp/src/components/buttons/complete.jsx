import React, { useState } from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../icons/compassIcons';

const CompleteButton = (props) => {
    const [active, setActive] = useState(false);

    const markAsDone = () => {
        setActive(true);
        setTimeout(() => {
            props.complete(props.issueId);
        }, 500);
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
    complete: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
};

export default CompleteButton;
