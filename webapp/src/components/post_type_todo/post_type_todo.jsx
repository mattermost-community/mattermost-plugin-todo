import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import RemoveButton from '../buttons/remove';
import CompleteButton from '../buttons/complete';
import AcceptButton from '../buttons/accept';

const PostUtils = window.PostUtils; // import the post utilities

export default class PostTypeTodo extends React.PureComponent {
    static propTypes = {
        post: PropTypes.object.isRequired,
        pendingAnswer: PropTypes.bool.isRequired,
        theme: PropTypes.object.isRequired,
        siteURL: PropTypes.string.isRequired,
        actions: PropTypes.shape({
            complete: PropTypes.func.isRequired,
            remove: PropTypes.func.isRequired,
            accept: PropTypes.func.isRequired,
            telemetry: PropTypes.func.isRequired,
        }).isRequired,
    };

    static defaultProps = {};

    constructor(props) {
        super(props);

        this.state = {done: false};
    }

    render() {
        const style = getStyle(this.props.theme);

        const preText = 'Automated message';

        const titleHTMLFormattedText = PostUtils.formatText(this.props.post.props.message, {siteURL: this.props.siteURL});
        const title = PostUtils.messageHtmlToComponent(titleHTMLFormattedText);

        const subtitleHTMLFormattedText = PostUtils.formatText(this.props.post.props.todo, {siteURL: this.props.siteURL});
        const subtitle = PostUtils.messageHtmlToComponent(subtitleHTMLFormattedText);

        const postPermalink = (
            <a
                className='theme markdown_link'
                href={this.props.post.props.postPermalink}
                rel='noreferrer'
                data-link={this.props.post.props.postPermalink}
            >
                <span data-link={this.props.post.props.postPermalink}>
                    {'Permalink'}
                </span>
            </a>
        );

        const content = (
            <div
                className={`todo-post d-flex flex-row-reverse align-items-center justify-content-end ${this.state.done ? 'todo-item--done' : ''}`}
                style={style.body}
            >
                <RemoveButton
                    issueId={this.props.post.props.issueId}
                    remove={(issueID) => {
                        this.props.actions.telemetry('custom_post_remove');
                        this.props.actions.remove(issueID);
                    }}
                />
                <AcceptButton
                    issueId={this.props.post.props.issueId}
                    accept={(issueID) => {
                        this.props.actions.telemetry('custom_post_accept');
                        this.props.actions.accept(issueID);
                    }}
                />
                <CompleteButton
                    active={this.state.done}
                    markAsDone={() => this.setState({done: true})}
                    issueId={this.props.post.props.issueId}
                    completeToast={() => {
                        this.props.actions.telemetry('custom_post_complete');
                        this.props.actions.complete(this.props.post.props.issueId);
                    }}
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
                            {postPermalink}
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
        summaryIssue: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            lineHeight: '26px',
        },
    };
});
