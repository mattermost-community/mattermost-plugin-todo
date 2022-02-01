import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import AutocompleteSelector from '../user_selector/autocomplete_selector.tsx';
import Button from '../../widget/buttons/button';
import IconButton from '../../widget/iconButton/iconButton';

import CompassIcon from '../icons/compassIcons';

const AssigneeModal = ({ visible, close, autocompleteUsers, theme, getAssignee, removeAssignee }) => {
    const [assignee, setAssignee] = useState();

    useEffect(() => {
        function handleKeypress(e) {
            if (e.key === 'Escape' && visible) {
                close();
            }
        }

        document.addEventListener('keydown', handleKeypress);

        return () => {
            document.removeEventListener('keydown', handleKeypress);
        };
    }, []);

    if (!visible) {
        return null;
    }

    const submit = () => {
        console.log(assignee);
        if (assignee) {
            console.log('add');
            getAssignee(assignee);
        } else {
            console.log('remove');
            removeAssignee();
        }
        close();
    };

    const changeAssignee = (selected) => {
        console.log(selected);
        setAssignee(selected);
    };

    const style = getStyle(theme);

    return (
        <div
            style={style.backdrop}
        >
            <div style={style.modal}>
                <h1 style={style.heading}>{'Assign todo to…'}</h1>
                <IconButton
                    size='medium'
                    style={style.closeIcon}
                    onClick={() => close()}
                    icon={<CompassIcon icon='close'/>}
                />
                <AutocompleteSelector
                    id='send_to_user'
                    loadOptions={autocompleteUsers}
                    onSelected={(selected) => changeAssignee(selected)}
                    placeholder={''}
                    theme={theme}
                />
                <div
                    className='todoplugin-button-container'
                    style={style.buttons}
                >
                    <Button
                        emphasis='tertiary'
                        size='medium'
                        onClick={() => close()}
                    >
                        {'Cancel'}
                    </Button>
                    <Button
                        emphasis='primary'
                        size='medium'
                        onClick={submit}
                        disabled={!assignee}
                    >
                        {'Assign'}
                    </Button>
                </div>
            </div>
        </div>
    );
};

AssigneeModal.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
    theme: PropTypes.object.isRequired,
    autocompleteUsers: PropTypes.func.isRequired,
    getAssignee: PropTypes.func.isRequired,
    removeAssignee: PropTypes.func.isRequired,
};

const getStyle = (theme) => ({
    backdrop: {
        position: 'absolute',
        display: 'flex',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.50)',
        zIndex: 2000,
        alignItems: 'center',
        justifyContent: 'center',
    },
    modal: {
        position: 'relative',
        width: 600,
        padding: 24,
        borderRadius: 8,
        maxWidth: '100%',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
    buttons: {
        marginTop: 24,
    },
    heading: {
        fontSize: 20,
        fontWeight: 600,
        margin: '0 0 24px 0',
    },
    closeIcon: {
        position: 'absolute',
        top: 8,
        right: 8,
    },
});

export default AssigneeModal;
