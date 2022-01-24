// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import TextOption from './textOption'

import './menu.scss'

type Props = {
    children: React.ReactNode
    position?: 'top'|'bottom'|'left'|'right'
}

export default class Menu extends React.PureComponent<Props> {
    static Text = TextOption

    public render(): JSX.Element {
        const {position, children} = this.props

        return (
            <div className={'Menu noselect ' + (position || 'bottom')}>
                <div className='menu-contents'>
                    <div className='menu-options'>
                        {children}
                    </div>

                    <div className='menu-spacer hideOnWidescreen'/>

                    <div className='menu-options hideOnWidescreen'>
                        <Menu.Text
                            id='menu-cancel'
                            name={'Cancel'}
                            className='menu-cancel'
                            onClick={this.onCancel}
                        />
                    </div>
                </div>
            </div>
        )
    }

    private onCancel = () => {
        // No need to do anything, as click bubbled up to MenuWrapper, which closes
    }
}
