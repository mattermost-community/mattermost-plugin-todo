import React from 'react';
import PropTypes from 'prop-types';

const AcceptButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.accept(props.issueId)}
        >{'Add to my list'}</button>
    );
};

AcceptButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    accept: PropTypes.func.isRequired,
};

export default AcceptButton;
