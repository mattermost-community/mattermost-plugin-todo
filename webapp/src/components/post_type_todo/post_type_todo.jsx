import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import DeleteButton from '../buttons/delete';
import CompleteButton from '../buttons/complete';
import EnqueueButton from '../buttons/enqueue';

export default class PostTypeTodo extends React.PureComponent {
    static propTypes = {
        post: PropTypes.object.isRequired,
        pendingAnswer: PropTypes.bool.isRequired,
        theme: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            complete: PropTypes.func.isRequired,
            remove: PropTypes.func.isRequired,
            enqueue: PropTypes.func.isRequired,
        }).isRequired,
    };

    static defaultProps = {};

    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {
        const style = getStyle(this.props.theme);

        const preText = 'Automated message';
        const title = this.props.post.props.message;
        const subtitle = this.props.post.props.todo;

        const content = (
            <div style={style.body}>
                <DeleteButton
                    itemId={this.props.post.props.itemId}
                    remove={this.props.actions.remove}
                    list={'my'}
                />
                <EnqueueButton
                    itemId={this.props.post.props.itemId}
                    enqueue={this.props.actions.enqueue}
                />
                <CompleteButton
                    itemId={this.props.post.props.itemId}
                    complete={this.props.actions.complete}
                />
            </div>
        );

        return (
            <div>
                {preText}
                <div style={style.attachment}>
                    <div style={style.content}>
                        <div style={style.container}>
                            <h1 style={style.title}>
                                {title}
                            </h1>
                            {subtitle}
                            <div>
                                {this.props.pendingAnswer && content}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        attachment: {
            marginLeft: '-20px',
            position: 'relative',
        },
        content: {
            borderRadius: '4px',
            borderStyle: 'solid',
            borderWidth: '1px',
            borderColor: '#BDBDBF',
            margin: '5px 0 5px 20px',
            padding: '2px 5px',
        },
        container: {
            borderLeftStyle: 'solid',
            borderLeftWidth: '4px',
            padding: '10px',
            borderLeftColor: '#89AECB',
        },
        body: {
            overflowX: 'auto',
            overflowY: 'hidden',
            paddingRight: '5px',
            width: '100%',
        },
        title: {
            fontSize: '16px',
            fontWeight: '600',
            height: '22px',
            lineHeight: '18px',
            margin: '5px 0 1px 0',
            padding: '0',
        },
        button: {
            fontFamily: 'Open Sans',
            fontSize: '12px',
            fontWeight: 'bold',
            letterSpacing: '1px',
            lineHeight: '19px',
            marginTop: '12px',
            borderRadius: '4px',
            color: theme.buttonColor,
        },
        buttonIcon: {
            paddingRight: '8px',
            fill: theme.buttonColor,
        },
        summary: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            fontWeight: '600',
            lineHeight: '26px',
            margin: '0',
            padding: '14px 0 0 0',
        },
        summaryItem: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            lineHeight: '26px',
        },
    };
});
