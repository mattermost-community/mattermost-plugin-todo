// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import AsyncSelect from 'react-select/async';
import {blendColors, changeOpacity, makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import {Theme} from 'mattermost-redux/types/preferences';
import {OptionsType, ValueType} from 'react-select/src/types';
import {UserProfile} from 'mattermost-redux/types/users';

type Props = {
    loadOptions: (inputValue: string, callback: ((options: OptionsType<UserProfile>) => void)) => Promise<unknown> | void,
    label?: string,
    labelClassName?: string,
    helpText?: string,
    inputClassName?: string,
    placeholder?: string,
    disabled?: boolean,
    onSelected?: (value: ValueType<UserProfile>) => void,
    theme: Theme,
}

// override react select theme based on mattermost theme
const getColorStyles = makeStyleFromTheme((theme) => {
    return {
        primary: changeOpacity(theme.centerChannelColor, 0.32),
        primary75: changeOpacity(theme.centerChannelColor, 0.24),
        primary50: changeOpacity(theme.centerChannelColor, 0.16),
        primary25: changeOpacity(theme.centerChannelColor, 0.08),
        danger: theme.errorTextColor,
        dangerLight: changeOpacity(theme.errorTextColor, 0.48),
        neutral0: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0),
        neutral5: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.08),
        neutral10: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.16),
        neutral20: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.24),
        neutral30: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.32),
        neutral40: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.40),
        neutral50: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.56),
        neutral60: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.64),
        neutral70: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.72),
        neutral80: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.8),
        neutral90: blendColors(theme.centerChannelBg, theme.centerChannelColor, 0.88),
    };
});

export default function AutocompleteSelector(props: Props) {
    const {
        loadOptions,
        label,
        labelClassName,
        helpText,
        inputClassName,
        placeholder,
        disabled,
        onSelected,
        theme,
    } = props;

    const handleSelected = (selected: ValueType<UserProfile>) => {
        if (onSelected) {
            onSelected(selected);
        }
    };

    let labelContent;
    if (label) {
        labelContent = (
            <label
                className={'control-label ' + labelClassName}
            >
                {label}
            </label>
        );
    }

    let helpTextContent;
    if (helpText) {
        helpTextContent = (
            <div className='help-text'>
                {helpText}
            </div>
        );
    }

    const getOptionData = (option: UserProfile) => option.username;

    const colors = getColorStyles(theme);

    return (
        <div
            data-testid='autoCompleteSelector'
            className='form-group'
        >
            {labelContent}
            <div className={inputClassName}>
                <AsyncSelect
                    cacheOptions={true}
                    loadOptions={loadOptions}
                    defaultOptions={true}
                    isClearable={true}
                    disabled={disabled}
                    placeholder={placeholder}
                    getOptionLabel={getOptionData}
                    getOptionValue={getOptionData}
                    onChange={handleSelected}
                    styles={{
                        option: (provided, state) => ({
                            ...provided,
                            color: state.isDisabled ? colors.neutral30 : colors.neutral90,
                        }),
                    }}
                    theme={(componentTheme) => ({
                        ...componentTheme,
                        colors: {
                            ...componentTheme.colors,
                            ...colors,
                        },
                    })}
                />
                {helpTextContent}
            </div>
        </div>
    );
}
