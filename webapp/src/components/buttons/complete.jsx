import React, {useCallback} from 'react';
import PropTypes from 'prop-types';

import CompassIcon from '../icons/compassIcons';

const CompleteButton = (props) => {
    const {markAsDone, completeToast} = props;

    const markTodoAsDone = useCallback(() => {
        markAsDone();
        setTimeout(() => {
            completeToast();
        }, 1000);
    }, [markAsDone, completeToast]);

    return (
        <button
            className={`todo-item__checkbox ${props.active ? 'todo-item__checkbox--active' : ''}`}
            onClick={markTodoAsDone}
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
