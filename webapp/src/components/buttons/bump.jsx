import React from 'react';
import PropTypes from 'prop-types';

const BumpButton = (props) => {
    return (
        <button
            className='btn btn-primary'
            onClick={() => props.bump(props.issueId)}
        >{'Bump'}</button>
    );
};

BumpButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    bump: PropTypes.func.isRequired,
};

export default BumpButton;
