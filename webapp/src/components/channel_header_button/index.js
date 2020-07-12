// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {showRHSPlugin} from '../../actions';
import {isPluginRhsOpen} from '../../selectors';

import ChannelHeaderButton from './channel_header_buttons.jsx';

function mapStateToProps(state) {
    return {
        issues: state['plugins-com.mattermost.plugin-todo'].issues,
        inIssues: state['plugins-com.mattermost.plugin-todo'].inIssues,
        outIssues: state['plugins-com.mattermost.plugin-todo'].outIssues,
        isTodoPluginRhsOpen: isPluginRhsOpen(state, 'plugins-com.mattermost.plugin-todo'),
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            showRHSPlugin,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(ChannelHeaderButton);