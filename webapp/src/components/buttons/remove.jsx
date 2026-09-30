import React from 'react';
import PropTypes from 'prop-types';

import Button from 'src/widget/buttons/button';

const RemoveButton = ({issueId, remove, list = 'my'}) => {
    return (
        <Button
            emphasis='tertiary'
            onClick={() => remove(issueId)}
        >
            {list === 'out' ? 'Cancel' : 'Won\'t do'}
        </Button>
    );
};

RemoveButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    remove: PropTypes.func.isRequired,
    list: PropTypes.string,
};
export default RemoveButton;
