import React from 'react';
import PropTypes from 'prop-types';

const BumpButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.bump(props.itemId)}
        >{'Bump'}</button>
    );
};

BumpButton.propTypes = {
    itemId: PropTypes.string.isRequired,
    bump: PropTypes.func.isRequired,
};

export default BumpButton;
