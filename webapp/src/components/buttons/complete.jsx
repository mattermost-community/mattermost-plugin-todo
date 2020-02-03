import React from 'react';
import PropTypes from 'prop-types';

const CompleteButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.complete(props.itemId)}
        >{'Done'}</button>
    );
};

CompleteButton.propTypes = {
    itemId: PropTypes.string.isRequired,
    complete: PropTypes.func.isRequired,
};

export default CompleteButton;