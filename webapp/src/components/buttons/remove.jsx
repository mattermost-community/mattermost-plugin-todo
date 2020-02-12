import React from 'react';
import PropTypes from 'prop-types';

const RemoveButton = (props) => {
    return (
        <button
            className='button'
            onClick={() => props.remove(props.issueId)}
        >{props.list === 'out' ? 'Cancel' : 'Won\'t do'}</button>
    );
};

RemoveButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    remove: PropTypes.func.isRequired,
    list: PropTypes.string,
};

RemoveButton.defaultProps = {
    list: 'my',
};

export default RemoveButton;
