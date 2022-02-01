import React from 'react';
import PropTypes from 'prop-types';

import {
    makeStyleFromTheme,
    changeOpacity,
} from 'mattermost-redux/utils/theme_utils';

import FullScreenModal from '../modals/full_screen_modal.jsx';
import Button from '../../widget/buttons/button';
import Chip from '../../widget/chip/chip';
import AutocompleteSelector from '../user_selector/autocomplete_selector.tsx';
import './add_issue.scss';
import CompassIcon from '../icons/compassIcons';
import { getProfilePicture } from '../../utils';

const PostUtils = window.PostUtils;

export default class AddIssue extends React.Component {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        message: PropTypes.string.isRequired,
        postID: PropTypes.string.isRequired,
        assignee: PropTypes.object.isRequired,
        closeAddBox: PropTypes.func.isRequired,
        submit: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
        autocompleteUsers: PropTypes.func.isRequired,
        openAssigneeModal: PropTypes.func.isRequired,
        removeAssignee: PropTypes.func.isRequired,
    };

    shouldComponentUpdate() {
        return true;
    }

    constructor(props) {
        super(props);

        this.state = {
            message: null,
            sendTo: null,
            attachToThread: false,
            previewMarkdown: false,
            assigneeModal: false,
        };
    }

    static getDerivedStateFromProps(props, state) {
        if (props.visible && state.message == null) {
            return { message: props.message };
        }
        if (!props.visible && (state.message != null || state.sendTo != null)) {
            return {
                message: null,
                sendTo: null,
                attachToThread: false,
                previewMarkdown: false,
            };
        }
        return null;
    }

    handleAttachChange = (e) => {
        const value = e.target.checked;
        if (value !== this.state.attachToThread) {
            this.setState({
                attachToThread: value,
            });
        }
    };

    close = () => {
        const { closeAddBox, removeAssignee } = this.props;
        removeAssignee();
        closeAddBox();
    }

    submit = () => {
        const { assignee, submit, postID, closeAddBox, removeAssignee } = this.props;
        const { message, attachToThread } = this.state;
        if (attachToThread) {
            submit(message, assignee.username, postID);
        } else {
            submit(message, assignee.username);
        }

        removeAssignee();
        closeAddBox();
    };

    toggleAssigneeModal = (value) => {
        this.setState({ assigneeModal: value });
    }

    render() {
        const { assignee, visible, theme } = this.props;

        if (!visible) {
            return null;
        }

        const { message } = this.state;

        const style = getStyle(theme);

        return (
            <div className='AddIssueBox'>
                <div className='AddIssueBox__body'>
                    <div className='AddIssueBox__check'/>
                    <div className='AddIssueBox__content'>
                        <div className='todoplugin-issue'>
                            {this.state.previewMarkdown ? (
                                <div
                                    className='todoplugin-input'
                                    style={style.markdown}
                                >
                                    {PostUtils.messageHtmlToComponent(
                                        PostUtils.formatText(this.state.message),
                                    )}
                                </div>
                            ) : (
                                <input
                                    className='todoplugin-input'
                                    placeholder='Enter a title'
                                    autoFocus={true}
                                    value={message}
                                    onChange={(e) =>
                                        this.setState({
                                            message: e.target.value,
                                        })
                                    }
                                />
                            )}
                        </div>
                        {this.props.postID && (
                            <div className='todoplugin-add-to-thread'>
                                <input
                                    type='checkbox'
                                    checked={this.state.attachToThread}
                                    onChange={this.handleAttachChange}
                                />
                                <b>{' Add to thread'}</b>
                                <div className='help-text'>
                                    {
                                        'Select to have the Todo Bot respond to the thread when the attached todo is added, modified or completed.'
                                    }
                                </div>
                            </div>
                        )}

                        <div style={style.chipsContainer}>
                            {!assignee && (
                                <Chip
                                    icon={<CompassIcon icon='account-outline'/>}
                                    onClick={() => this.props.openAssigneeModal('')}
                                >
                                    {'Assign to…'}
                                </Chip>
                            )}
                            {assignee && (
                                <button
                                    style={style.assigneeContainer}
                                    onClick={() => this.props.openAssigneeModal('')}
                                >
                                    <img
                                        style={style.assigneeImage}
                                        src={getProfilePicture(assignee.id)}
                                        alt={assignee.username}
                                    />
                                    <span>{assignee.username}</span>
                                </button>
                            )}
                        </div>

                        <FullScreenModal
                            show={this.state.assigneeModal}
                            onClose={() => this.toggleAssigneeModal(false)}
                        >
                            <AutocompleteSelector
                                id='send_to_user'
                                loadOptions={this.props.autocompleteUsers}
                                onSelected={(selected) =>
                                    this.setState({
                                        sendTo: selected?.username,
                                    })
                                }
                                label={'Send to user'}
                                helpText={
                                    'Select a user if you want to send this todo.'
                                }
                                placeholder={''}
                                theme={theme}
                            />
                        </FullScreenModal>
                    </div>
                </div>
                <div
                    className='todoplugin-button-container'
                    style={style.buttons}
                >
                    <Button
                        emphasis='tertiary'
                        size='small'
                        onClick={this.close}
                    >
                        {'Cancel'}
                    </Button>
                    <Button
                        emphasis='primary'
                        size='small'
                        onClick={this.submit}
                        disabled={!message}
                    >
                        {'Save'}
                    </Button>
                </div>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        modal: {
            color: changeOpacity(theme.centerChannelColor, 0.88),
        },
        textarea: {
            backgroundColor: theme.centerChannelBg,
        },
        helpText: {
            color: changeOpacity(theme.centerChannelColor, 0.64),
        },
        button: {
            color: theme.buttonColor,
            backgroundColor: theme.buttonBg,
            marginRight: 4,
            fontSize: 11,
            height: 24,
            padding: '0 10px',
        },
        inactiveButton: {
            color: changeOpacity(theme.buttonColor, 0.88),
            backgroundColor: changeOpacity(theme.buttonBg, 0.32),
        },
        markdown: {
            minHeight: '149px',
            fontSize: '16px',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'end',
        },
        assigneeImage: {
            width: 12,
            height: 12,
            marginRight: 6,
        },
        assigneeContainer: {
            borderRadius: 50,
            backgroundColor: changeOpacity(theme.centerChannelColor, 0.08),
            height: 24,
            padding: '4px 10px',
            fontWeight: 600,
            alignItems: 'center',
            justifyContent: 'center',
            display: 'inline-flex',
            border: 0,
        },
        buttons: {
            marginTop: 16,
        },
        chipsContainer: {
            marginTop: 8,
        },
    };
});
