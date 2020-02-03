import React from 'react';
import PropTypes from 'prop-types';

const DeleteButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.remove(props.itemId)}
        >{props.list === 'out' ? 'Cancel' : 'Won\'t do'}</button>
    );
};

DeleteButton.propTypes = {
    itemId: PropTypes.string.isRequired,
    remove: PropTypes.func.isRequired,
    list: PropTypes.string,
};

DeleteButton.defaultProps = {
    list: 'my',
};

export default DeleteButton;
