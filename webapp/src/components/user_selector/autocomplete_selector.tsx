// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {components, FormatOptionLabelMeta, InputProps, OptionsOrGroups, GroupBase, SingleValue, StylesConfig, ThemeConfig, Theme as ComponentTheme} from 'react-select';
import AsyncSelect from 'react-select/async';

import {Theme} from 'mattermost-redux/types/preferences';
import {UserProfile} from 'mattermost-redux/types/users';

import {getColorStyles, getDescription, getProfilePicture} from '../../utils';

import './autocomplete_selector.scss';

type Props = {
    loadOptions: (inputValue: string, callback: ((options: OptionsOrGroups<UserProfile, GroupBase<UserProfile>>) => void)) => Promise<OptionsOrGroups<UserProfile, GroupBase<UserProfile>>> | void,
    autoFocus?: boolean,
    label?: string,
    labelClassName?: string,
    helpText?: string,
    inputClassName?: string,
    placeholder?: string,
    disabled?: boolean,
    onSelected?: (value: SingleValue<UserProfile>) => void,
    theme: Theme,
}

const useTheme = (mattermostTheme: Theme): [StylesConfig<UserProfile, false>, ThemeConfig] => {
    const mmColors = getColorStyles(mattermostTheme);

    const styles: StylesConfig<UserProfile, false> = {
        option: (provided, state) => ({
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

const renderOption = (option: UserProfile, {context}: FormatOptionLabelMeta<UserProfile>) => {
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

    const handleSelected = (selected: SingleValue<UserProfile>) => {
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

    const Input = (inputProps: InputProps<UserProfile, false>) => (
        <components.Input
            {...inputProps}
            maxLength={22}
        />);

    return (
        <div
            data-testid='autoCompleteSelector'
            className='form-group todo-select'
        >
            {labelContent}
            <div className={inputClassName}>
                <AsyncSelect
                    components={{Input}}
                    autoFocus={Boolean(props.autoFocus)}
                    cacheOptions={true}
                    loadOptions={loadOptions}
                    defaultOptions={true}
                    isClearable={true}
                    isDisabled={disabled}
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
