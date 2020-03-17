import React from 'react';
import PropTypes from 'prop-types';

const CompleteButton = (props) => {
    return (
        <button
            className='btn btn-primary'
            onClick={() => props.complete(props.issueId)}
        >{'Done'}</button>
    );
};

CompleteButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    complete: PropTypes.func.isRequired,
};

export default CompleteButton;
