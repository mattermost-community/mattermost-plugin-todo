import React from 'react';
import PropTypes from 'prop-types';

import {
    makeStyleFromTheme,
    changeOpacity,
} from 'mattermost-redux/utils/theme_utils';

import TextareaAutosize from 'react-textarea-autosize';

import FullScreenModal from '../modals/modals.jsx';
import Button from '../../widget/buttons/button';
import Chip from '../../widget/chip/chip';
import AutocompleteSelector from '../user_selector/autocomplete_selector.tsx';
import './add_issue.scss';
import CompassIcon from '../icons/compassIcons';
import {getProfilePicture} from '../../utils';

const PostUtils = window.PostUtils;

export default class AddIssue extends React.PureComponent {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        message: PropTypes.string.isRequired,
        postPermalink: PropTypes.string,
        postID: PropTypes.string.isRequired,
        assignee: PropTypes.object,
        closeAddBox: PropTypes.func.isRequired,
        submit: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
        autocompleteUsers: PropTypes.func.isRequired,
        openAssigneeModal: PropTypes.func.isRequired,
        removeAssignee: PropTypes.func.isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            message: props.message || '',
            postPermalink: props.postPermalink || '',
            description: '',
            sendTo: null,
            attachToThread: false,
            previewMarkdown: false,
            assigneeModal: false,
            isTyping: false,
        };
    }

    static getDerivedStateFromProps(props, state) {
        if (state.isTyping) {
            return null;
        }
        if (props.visible && (props.message !== state.message || props.postPermalink !== state.postPermalink)) {
            return {
                message: props.message,
                postPermalink: props.postPermalink,
            };
        }
        if (!props.visible && (state.message || state.sendTo)) {
            return {
                message: '',
                postPermalink: '',
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
        const {closeAddBox, removeAssignee} = this.props;
        removeAssignee();
        closeAddBox();
    }

    submit = () => {
        const {submit, postID, assignee, closeAddBox, removeAssignee} = this.props;
        const {message, postPermalink, description, attachToThread, sendTo} = this.state;
        this.setState({
            message: '',
            description: '',
            postPermalink: '',
            isTyping: false,
        });

        if (attachToThread) {
            if (assignee) {
                submit(message, postPermalink, description, assignee.username, postID);
            } else {
                submit(message, postPermalink, description, sendTo, postID);
            }
        } else if (assignee) {
            submit(message, postPermalink, description, assignee.username);
        } else {
            submit(message, postPermalink, description);
        }

        removeAssignee();
        closeAddBox();
    };

    toggleAssigneeModal = (value) => {
        this.setState({assigneeModal: value});
    }

    onKeyDown = (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            this.submit();
        }

        if (e.key === 'Escape') {
            this.close();
        }
    }

    handleInputChange = (e, field) => {
        this.setState({
            [field]: e.target.value,
            isTyping: true,
        });
    }

    render() {
        const {assignee, visible, theme} = this.props;

        if (!visible) {
            return null;
        }

        const {message, description, postPermalink} = this.state;
        const style = getStyle(theme);
        const formattedMessage = message.includes(`[Permalink](${postPermalink})`) ? message : message + (postPermalink ? `\n[Permalink](${postPermalink})` : '');

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
                                        PostUtils.formatText(formattedMessage),
                                    )}
                                </div>
                            ) : (
                                <React.Fragment>
                                    <TextareaAutosize
                                        style={style.textareaResizeMessage}
                                        placeholder='Enter a title'
                                        autoFocus={true}
                                        onKeyDown={(e) => this.onKeyDown(e)}
                                        value={formattedMessage}
                                        onChange={(e) => this.handleInputChange(e, 'message')}
                                    />
                                    <TextareaAutosize
                                        style={style.textareaResizeDescription}
                                        placeholder='Enter a description'
                                        onKeyDown={(e) => this.onKeyDown(e)}
                                        value={description}
                                        onChange={(e) => this.handleInputChange(e, 'description')}
                                    />
                                </React.Fragment>
                            )}
                        </div>
                        {this.props.postID && (
                            <div className='todo-add-issue'>
                                <label>
                                    <input
                                        type='checkbox'
                                        checked={this.state.attachToThread}
                                        onChange={this.handleAttachChange}
                                    />
                                    <b>{' Add to thread'}</b>
                                </label>
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
            borderRadius: 12,
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
        textareaResizeMessage: {
            border: 0,
            padding: 0,
            fontSize: 14,
            width: '100%',
            backgroundColor: 'transparent',
            resize: 'none',
            boxShadow: 'none',
        },
        textareaResizeDescription: {
            fontSize: 12,
            color: changeOpacity(theme.centerChannelColor, 0.72),
            marginTop: 1,
            border: 0,
            padding: 0,
            width: '100%',
            backgroundColor: 'transparent',
            resize: 'none',
            boxShadow: 'none',
        },
    };
});
