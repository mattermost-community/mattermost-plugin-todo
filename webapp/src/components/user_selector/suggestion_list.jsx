// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import PropTypes from 'prop-types';
import React from 'react';
import ReactDOM from 'react-dom';

import Constants from '../../constants';

export default class SuggestionList extends React.PureComponent {
    static propTypes = {
        ariaLiveRef: PropTypes.object,
        open: PropTypes.bool.isRequired,
        location: PropTypes.string,
        renderNoResults: PropTypes.bool,
        onCompleteWord: PropTypes.func.isRequired,
        preventClose: PropTypes.func,
        onItemHover: PropTypes.func.isRequired,
        pretext: PropTypes.string.isRequired,
        cleared: PropTypes.bool.isRequired,
        matchedPretext: PropTypes.array.isRequired,
        items: PropTypes.array.isRequired,
        terms: PropTypes.array.isRequired,
        selection: PropTypes.string.isRequired,
        components: PropTypes.array.isRequired,
        wrapperHeight: PropTypes.number,
    };

    static defaultProps = {
        renderNoResults: false,
    };

    constructor(props) {
        super(props);

        this.contentRef = React.createRef();
        this.itemRefs = new Map();
        this.suggestionReadOut = React.createRef();
        this.currentLabel = '';
        this.currentItem = {};
    }

    componentDidUpdate(prevProps) {
        if (this.props.selection !== prevProps.selection && this.props.selection) {
            this.scrollToItem(this.props.selection);
        }
    }

    getContent = () => {
        return $(this.contentRef.current);
    }

    scrollToItem = (term) => {
        const content = this.getContent();
        if (!content || content.length === 0) {
            return;
        }

        const visibleContentHeight = content[0].clientHeight;
        const actualContentHeight = content[0].scrollHeight;

        if (visibleContentHeight < actualContentHeight) {
            const contentTop = content.scrollTop();
            const contentTopPadding = parseInt(content.css('padding-top'), 10);
            const contentBottomPadding = parseInt(content.css('padding-top'), 10);

            const item = $(ReactDOM.findDOMNode(this.itemRefs.get(term)));
            if (item.length === 0) {
                return;
            }

            const itemTop = item[0].offsetTop - parseInt(item.css('margin-top'), 10);
            const itemBottomMargin = parseInt(item.css('margin-bottom'), 10) + parseInt(item.css('padding-bottom'), 10);
            const itemBottom = item[0].offsetTop + item.height() + itemBottomMargin;

            if (itemTop - contentTopPadding < contentTop) {
                // the item is off the top of the visible space
                content.scrollTop(itemTop - contentTopPadding);
            } else if (itemBottom + contentTopPadding + contentBottomPadding > contentTop + visibleContentHeight) {
                // the item has gone off the bottom of the visible space
                content.scrollTop((itemBottom - visibleContentHeight) + contentTopPadding + contentBottomPadding);
            }
        }
    }

    renderNoResults() {
        return (
            <div
                key='list-no-results'
                className='suggestion-list__no-results'
                ref={this.contentRef}
            >
                <span id='suggestion_list.no_matches'>
                    {'No items match ' + (this.props.pretext || '""')}
                </span>
            </div>
        );
    }

    render() {
        if (!this.props.open || this.props.cleared) {
            return null;
        }

        const items = [];
        if (this.props.items.length === 0) {
            if (!this.props.renderNoResults) {
                return null;
            }
            items.push(this.renderNoResults());
        }

        for (let i = 0; i < this.props.items.length; i++) {
            const item = this.props.items[i];
            const term = this.props.terms[i];
            const isSelection = term === this.props.selection;

            // ReactComponent names need to be upper case when used in JSX
            const Component = this.props.components[i];

            if (item.loading) {
                items.push(
                    <i
                        className='fa fa-spinner fa-fw fa-pulse spinner'
                        key={item.type}
                    />
                );
                continue;
            }

            if (isSelection) {
                this.currentItem = item;
            }

            items.push(
                <Component
                    key={term}
                    ref={(ref) => this.itemRefs.set(term, ref)}
                    item={this.props.items[i]}
                    term={term}
                    matchedPretext={this.props.matchedPretext[i]}
                    isSelection={isSelection}
                    onClick={this.props.onCompleteWord}
                    onMouseMove={this.props.onItemHover}
                />
            );
        }
        const mainClass = 'suggestion-list suggestion-list--' + this.props.location;
        const contentClass = 'suggestion-list__content suggestion-list__content--' + this.props.location;
        let maxHeight = Constants.SUGGESTION_LIST_MAXHEIGHT;
        if (this.props.wrapperHeight) {
            maxHeight = Math.min(
                window.innerHeight - (this.props.wrapperHeight + Constants.SUGGESTION_LIST_MAXHEIGHT),
                Constants.SUGGESTION_LIST_MAXHEIGHT
            );
        }

        const contentStyle = {maxHeight};

        return (
            <div className={mainClass}>
                <div
                    id='suggestionList'
                    ref={this.contentRef}
                    style={{...contentStyle}}
                    className={contentClass}
                    onMouseDown={this.props.preventClose}
                >
                    {items}
                </div>
            </div>
        );
    }
}
