import React from 'react';
import PropTypes from 'prop-types';

const EnqueueButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.enqueue(props.itemId)}
        >{'Add to my list'}</button>
    );
};

EnqueueButton.propTypes = {
    itemId: PropTypes.string.isRequired,
    enqueue: PropTypes.func.isRequired,
};

export default EnqueueButton;
