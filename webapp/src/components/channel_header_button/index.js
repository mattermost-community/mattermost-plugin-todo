// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import {isPluginRhsOpen} from 'selectors';

import ChannelHeaderButton from './channel_header_buttons.tsx';

function mapStateToProps(state) {
    return {
        shouldHighlight: isPluginRhsOpen(state),
    };
}

export default connect(mapStateToProps)(ChannelHeaderButton);
