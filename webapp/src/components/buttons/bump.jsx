import React from 'react';
import PropTypes from 'prop-types';

const BumpButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.bump(props.issueId)}
        >{'Bump'}</button>
    );
};

BumpButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    bump: PropTypes.func.isRequired,
};

export default BumpButton;
