// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {CSSProperties} from 'react';

import AsyncSelect from 'react-select/async';
import {OptionsType, ValueType, Theme as ComponentTheme} from 'react-select/src/types';
import {Props as ComponentProps, StylesConfig} from 'react-select/src/styles';
import {ThemeConfig} from 'react-select/src/theme';

import {Theme} from 'mattermost-redux/types/preferences';
import {UserProfile} from 'mattermost-redux/types/users';

import {FormatOptionLabelContext} from 'react-select/src/Select';

import {getColorStyles, getDescription, getProfilePicture} from '../../utils';

import './autocomplete_selector.scss';

type Props = {
    loadOptions: (inputValue: string, callback: ((options: OptionsType<UserProfile>) => void)) => Promise<unknown> | void,
    autoFocus?: boolean,
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
    const mmColors = getColorStyles(mattermostTheme);

    const styles: StylesConfig = {
        option: (provided: CSSProperties, state: ComponentProps) => ({
            ...provided,
            color: state.isDisabled ? mmColors.neutral30 : mmColors.neutral90,
        }),
    };

    const compTheme: ThemeConfig = (componentTheme: ComponentTheme) => ({
        ...componentTheme,
        colors: {
            ...componentTheme.colors,
            ...mmColors,
        },
    });

    return [styles, compTheme];
};

const renderOption = (option: UserProfile, {context} : {context: FormatOptionLabelContext}) => {
    const {username} = option;
    const name = `@${username}`;
    const description = getDescription(option);

    if (context === 'menu') {
        return (
            <div>
                <img
                    className={'option-image'}
                    src={getProfilePicture(option.id)}
                    alt={option.username}
                />
                <span className={'option-username'}>{name}</span>
                {description !== '' && (
                    <span className={'option-nickname'}>{description}</span>
                )}
            </div>
        );
    }

    return <div>{name}</div>;
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

    const [styles, componentTheme] = useTheme(theme);

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

    return (
        <div
            data-testid='autoCompleteSelector'
            className='form-group todo-select'
        >
            {labelContent}
            <div className={inputClassName}>
                <AsyncSelect
                    autoFocus={props.autoFocus || false}
                    cacheOptions={true}
                    loadOptions={loadOptions}
                    defaultOptions={true}
                    isClearable={true}
                    disabled={disabled}
                    placeholder={placeholder}
                    getOptionLabel={(option: UserProfile) => option.username}
                    getOptionValue={(option: UserProfile) => option.id}
                    formatOptionLabel={renderOption}
                    onChange={handleSelected}
                    styles={styles}
                    theme={componentTheme}
                />
                {helpTextContent}
            </div>
        </div>
    );
}
