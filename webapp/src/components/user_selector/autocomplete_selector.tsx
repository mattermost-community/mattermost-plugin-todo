// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState, CSSProperties} from 'react';

import AsyncSelect from 'react-select/async';
import {OptionsType, ValueType, Theme as ComponentTheme} from 'react-select/src/types';
import {Props as ComponentProps, StylesConfig} from 'react-select/src/styles';
import {ThemeConfig} from 'react-select/src/theme';

import {Theme} from 'mattermost-redux/types/preferences';
import {UserProfile} from 'mattermost-redux/types/users';

import {getColorStyles} from '../../utils';

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

const useTheme = (mattermostTheme: Theme): [StylesConfig, ThemeConfig] => {
    const colors = getColorStyles(mattermostTheme);

    const styles: StylesConfig = {
        option: (provided: CSSProperties, state: ComponentProps) => ({
            ...provided,
            color: state.isDisabled ? colors.neutral30 : colors.neutral90,
        }),
    };

    const compTheme: ThemeConfig = (componentTheme: ComponentTheme) => ({
        ...componentTheme,
        colors: {
            ...componentTheme.colors,
            ...colors,
        },
    });

    return [styles, compTheme];
};

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

    const [styles, componentTheme] = useTheme(theme);

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
                    styles={styles}
                    theme={componentTheme}
                />
                {helpTextContent}
            </div>
        </div>
    );
}
