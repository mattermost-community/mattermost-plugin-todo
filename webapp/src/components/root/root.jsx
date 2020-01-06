import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import FullScreenModal from '../modals/full_screen_modal.jsx';

import './root.scss';

export default class Root extends React.Component {
    static propTypes = {
        visible: PropTypes.bool.isRequired,
        message: PropTypes.string.isRequired,
        postID: PropTypes.string.isRequired,
        close: PropTypes.func.isRequired,
        submit: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
    }
    constructor(props) {
        super(props);

        this.state = {
            message: null,
        };
    }

    static getDerivedStateFromProps(props, state) {
        if (props.visible && state.message == null) {
            return {message: props.message};
        }
        if (!props.visible && state.message != null) {
            return {message: null};
        }
        return null;
    }

    submit = () => {
        const {submit, close} = this.props;
        const {message} = this.state;
        submit(message);
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
                    <div className='todoplugin-item'>
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
                            {'Adding a to do will add an item to your to do list. You will get daily reminders about your to do items until you mark them as complete.'}
                        </div>
                        <div className='todoplugin-question'>
                            {'How is this different from flagging a post?'}
                        </div>
                        <div className='todoplugin-answer'>
                            {'To do items are disconnected from posts. You can generate to do items from posts but they have no other assoication to the posts. This allows for a cleaner to do list that does not rely on post history or someone else not deleting or editing the post.'}
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
            color: changeOpacity(theme.buttonColor, 0.5),
            backgroundColor: changeOpacity(theme.buttonBg, 0.1),
        },
    };
});