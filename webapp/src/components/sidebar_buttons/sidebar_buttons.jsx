// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {Tooltip, OverlayTrigger} from 'react-bootstrap';
import PropTypes from 'prop-types';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import {RHSStates} from '../../constants';

export default class SidebarButtons extends React.PureComponent {
    static propTypes = {
        theme: PropTypes.object.isRequired,
        isTeamSidebar: PropTypes.bool,
        showRHSPlugin: PropTypes.func.isRequired,
        issues: PropTypes.arrayOf(PropTypes.object),
        inIssues: PropTypes.arrayOf(PropTypes.object),
        outIssues: PropTypes.arrayOf(PropTypes.object),
        actions: PropTypes.shape({
            list: PropTypes.func.isRequired,
            updateRhsState: PropTypes.func.isRequired,
        }).isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            refreshing: false,
        };
    }

    componentDidUpdate() {
        this.getData();
    }

    getData = async (e) => {
        if (this.state.refreshing) {
            return;
        }

        if (e) {
            e.preventDefault();
        }

        this.setState({refreshing: true});
        await Promise.all([
            this.props.actions.list(false, 'my'),
            this.props.actions.list(false, 'in'),
            this.props.actions.list(false, 'out'),
        ]);
        this.setState({refreshing: false});
    }

    openRHS = (rhsState) => {
        this.props.actions.updateRhsState(rhsState);
        this.props.showRHSPlugin();
    }

    render() {
        const style = getStyle(this.props.theme);
        const isTeamSidebar = this.props.isTeamSidebar;

        let container = style.containerHeader;
        let button = style.buttonHeader;
        let placement = 'bottom';
        if (isTeamSidebar) {
            placement = 'right';
            button = style.buttonTeam;
            container = style.containerTeam;
        }

        const issues = this.props.issues || [];
        const inIssues = this.props.inIssues || [];
        const outIssues = this.props.outIssues || [];
        const refreshClass = this.state.refreshing ? ' fa-spin' : '';

        return (
            <div style={container}>
                <OverlayTrigger
                    key='todoLink'
                    placement={placement}
                    overlay={<Tooltip id='todoTooltip'>{'Your todos'}</Tooltip>}
                >
                    <a
                        style={button}
                        onClick={() => this.openRHS(RHSStates.ISSUES)}
                    >
                        <i className='fa fa-check'/>
                        {' ' + issues.length}
                    </a>
                </OverlayTrigger>
                <OverlayTrigger
                    key='todoReviewsLink'
                    placement={placement}
                    overlay={<Tooltip id='reviewTooltip'>{'Todos received'}</Tooltip>}
                >
                    <a
                        onClick={() => this.openRHS(RHSStates.IN_ISSUES)}
                        style={button}
                    >
                        <i className='fa fa-angle-double-down'/>
                        {' ' + inIssues.length}
                    </a>
                </OverlayTrigger>
                <OverlayTrigger
                    key='todoAssignmentsLink'
                    placement={placement}
                    overlay={<Tooltip id='reviewTooltip'>{'Todos sent'}</Tooltip>}
                >
                    <a
                        onClick={() => this.openRHS(RHSStates.OUT_ISSUES)}
                        style={button}
                    >
                        <i className='fa fa-angle-double-up'/>
                        {' ' + outIssues.length}
                    </a>
                </OverlayTrigger>
                <OverlayTrigger
                    key='todoRefreshButton'
                    placement={placement}
                    overlay={<Tooltip id='refreshTooltip'>{'Refresh'}</Tooltip>}
                >
                    <a
                        href='#'
                        style={button}
                        onClick={this.getData}
                    >
                        <i className={'fa fa-refresh' + refreshClass}/>
                    </a>
                </OverlayTrigger>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        buttonTeam: {
            color: changeOpacity(theme.sidebarText, 0.6),
            display: 'block',
            marginBottom: '10px',
            width: '100%',
        },
        buttonHeader: {
            color: changeOpacity(theme.sidebarText, 0.6),
            textAlign: 'center',
            cursor: 'pointer',
        },
        containerHeader: {
            marginTop: '10px',
            marginBottom: '5px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-around',
            padding: '0 10px',
        },
        containerTeam: {
        },
    };
});