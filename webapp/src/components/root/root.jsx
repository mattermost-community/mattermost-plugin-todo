import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import FullScreenModal from '../modals/full_screen_modal.jsx';

import './root.scss';
import AutocompleteSelector from '../user_selector/autocomplete_selector.jsx';
import GenericUserProvider from '../user_selector/generic_user_provider.jsx';

export default class Root extends React.Component {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        message: PropTypes.string.isRequired,
        postID: PropTypes.string.isRequired,
        close: PropTypes.func.isRequired,
        submit: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
        autocompleteUsers: PropTypes.func.isRequired,
    }
    constructor(props) {
        super(props);

        this.state = {
            message: null,
            sendTo: null,
        };
    }

    static getDerivedStateFromProps(props, state) {
        if (props.visible && state.message == null) {
            return {message: props.message};
        }
        if (!props.visible && (state.message != null || state.sendTo != null)) {
            return {message: null, sendTo: null};
        }
        return null;
    }

    submit = () => {
        const {submit, close} = this.props;
        const {message, sendTo} = this.state;
        submit(message, sendTo);
        close();
    }

    render() {
        const {visible, theme, close} = this.props;

        if (!visible) {
            return null;
        }

        const {message} = this.state;

        const style = getStyle(theme);

        return (
            <FullScreenModal
                show={visible}
                onClose={close}
            >
                <div
                    style={style.modal}
                    className='ToDoPluginRootModal'
                >
                    <h1>{'Add a To Do'}</h1>
                    <div className='todoplugin-issue'>
                        <h2>
                            {'To Do Message'}
                        </h2>
                        <textarea
                            className='todoplugin-input'
                            style={style.textarea}
                            value={message}
                            onChange={(e) => this.setState({message: e.target.value})}
                        />
                    </div>
                    <div>
                        <AutocompleteSelector
                            id='send_to_user'
                            providers={[new GenericUserProvider(this.props.autocompleteUsers)]}
                            onSelected={(selected) => this.setState({sendTo: selected.username})}
                            label={'Send to user'}
                            helpText={'Select a user if you want to send this todo.'}
                            placeholder={''}
                            value={this.state.sendTo}
                        />
                    </div>
                    <div className='todoplugin-button-container'>
                        <button
                            className={'btn btn-primary'}
                            style={message ? style.button : style.inactiveButton}
                            onClick={this.submit}
                            disabled={!message}
                        >
                            {'Add To Do'}
                        </button>
                    </div>
                    <div className='todoplugin-divider'/>
                    <div className='todoplugin-clarification'>
                        <div className='todoplugin-question'>
                            {'What does this do?'}
                        </div>
                        <div className='todoplugin-answer'>
                            {'Adding a to do will add an issue to your to do list. You will get daily reminders about your to do issues until you mark them as complete.'}
                        </div>
                        <div className='todoplugin-question'>
                            {'How is this different from flagging a post?'}
                        </div>
                        <div className='todoplugin-answer'>
                            {'To do issues are disconnected from posts. You can generate to do issues from posts but they have no other assoication to the posts. This allows for a cleaner to do list that does not rely on post history or someone else not deleting or editing the post.'}
                        </div>
                    </div>
                </div>
            </FullScreenModal>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        modal: {
            color: changeOpacity(theme.centerChannelColor, 0.9),
        },
        textarea: {
            backgroundColor: theme.centerChannelBg,
        },
        helpText: {
            color: changeOpacity(theme.centerChannelColor, 0.6),
        },
        button: {
            color: theme.buttonColor,
            backgroundColor: theme.buttonBg,
        },
        inactiveButton: {
            color: '#000000',
            backgroundColor: changeOpacity(theme.buttonBg, 0.1),
        },
    };
});