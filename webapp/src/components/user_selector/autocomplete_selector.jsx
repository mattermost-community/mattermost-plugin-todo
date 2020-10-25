// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import PropTypes from 'prop-types';
import React from 'react';

import AsyncSelect from 'react-select/async';

export default function AutocompleteSelector(props) {
    const {
        loadOptions,
        footer,
        label,
        labelClassName,
        helpText,
        inputClassName,
        placeholder,
        disabled,
        onSelected,
    } = props;

    const handleSelected = (selected) => {
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

    const getOptionData = (option) => option.username;

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
                {footer}
            </div>
        </div>
    );
}

AutocompleteSelector.propTypes = {
    loadOptions: PropTypes.func.isRequired,
    onSelected: PropTypes.func,
    label: PropTypes.node,
    labelClassName: PropTypes.string,
    inputClassName: PropTypes.string,
    helpText: PropTypes.node,
    placeholder: PropTypes.string,
    footer: PropTypes.node,
    disabled: PropTypes.bool,
};

AutocompleteSelector.defaultProps = {
    value: '',
    id: '',
    labelClassName: '',
    inputClassName: '',
};
