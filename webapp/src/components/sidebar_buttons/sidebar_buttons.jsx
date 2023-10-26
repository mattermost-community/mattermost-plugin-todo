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
        countIssues: PropTypes.object,
        actions: PropTypes.shape({
            count: PropTypes.func.isRequired,
            updateRhsState: PropTypes.func.isRequired,
            telemetry: PropTypes.func.isRequired,
        }).isRequired,
    };

    constructor(props) {
        super(props);

        this.state = {
            refreshing: false,
        };
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

        const countIssues = this.props.countIssues;

        return (
            <div style={container}>
                <OverlayTrigger
                    key='myTodosLink'
                    placement={placement}
                    overlay={<Tooltip id='myTodosTooltip'>{'Your Todos'}</Tooltip>}
                >
                    <a
                        style={button}
                        onClick={() => {
                            this.props.actions.telemetry('click_lhs_my');
                            this.openRHS(RHSStates.InListName);
                        }}
                    >
                        <i className='icon icon-check'/>
                        {' ' + countIssues.my}
                    </a>
                </OverlayTrigger>
                <OverlayTrigger
                    key='incomingTodosLink'
                    placement={placement}
                    overlay={<Tooltip id='incomingTodosTooltip'>{'Incoming Todos'}</Tooltip>}
                >
                    <a
                        onClick={() => {
                            this.props.actions.telemetry('click_lhs_in');
                            this.openRHS(RHSStates.InListName);
                        }}
                        style={button}
                    >
                        <i className='icon icon-arrow-down'/>
                        {' ' + countIssues.in}
                    </a>
                </OverlayTrigger>
                <OverlayTrigger
                    key='outgoingTodosLink'
                    placement={placement}
                    overlay={<Tooltip id='outgoingTodosTooltip'>{'Outgoing Todos'}</Tooltip>}
                >
                    <a
                        onClick={() => {
                            this.props.actions.telemetry('click_lhs_out');
                            this.openRHS(RHSStates.OutListName);
                        }}
                        style={button}
                    >
                        <i className='icon icon-arrow-up'/>
                        {' ' + countIssues.out}
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
