import React from 'react';
import PropTypes from 'prop-types';

import Button from 'src/widget/buttons/button';

const AcceptButton = (props) => {
    return (
        <Button
            emphasis={'secondary'}
            onClick={() => props.accept(props.issueId)}
        >
            {'Add to my list'}
        </Button>
    );
};

AcceptButton.propTypes = {
    issueId: PropTypes.string.isRequired,
    accept: PropTypes.func.isRequired,
};

export default AcceptButton;
