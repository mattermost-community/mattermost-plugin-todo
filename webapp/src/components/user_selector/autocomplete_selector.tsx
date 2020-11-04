// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import AsyncSelect from 'react-select/async';
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
                />
                {helpTextContent}
            </div>
        </div>
    );
}
