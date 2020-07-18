// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import { bindActionCreators } from 'redux';

import {showRHSPlugin} from 'actions';

import ChannelHeaderButton from './channel_header_button';

function mapStateToProps(state) {
    return {
        isVisible: state['plugins-com.mattermost.plugin-todo'].isVisible,
    };
}

function mapDispatchToProp(dispatch) {
    return {
        actions: bindActionCreators({
            showRHSPlugin,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProp)(ChannelHeaderButton);
