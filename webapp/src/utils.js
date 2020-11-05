import {blendColors, changeOpacity, makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import Constants from './constants';

export function canRemove(myList, foreignList) {
    return myList === 'my' || myList === 'in' || foreignList === 'in';
}

export function canComplete(myList) {
    return myList === 'my' || myList === 'in';
}

export function canAccept(myList) {
    return myList === 'in';
}

export function canBump(myList, foreignList) {
    return myList === 'out' && foreignList === 'in';
}

export function setSelectionRange(input, selectionStart, selectionEnd) {
    if (input.setSelectionRange) {
        input.focus();
        input.setSelectionRange(selectionStart, selectionEnd);
    } else if (input.createTextRange) {
        var range = input.createTextRange();
        range.collapse(true);
        range.moveEnd('character', selectionEnd);
        range.moveStart('character', selectionStart);
        range.select();
    }
}

export function setCaretPosition(input, pos) {
    setSelectionRange(input, pos, pos);
}

export function isKeyPressed(event, key) {
    // There are two types of keyboards
    // 1. English with different layouts(Ex: Dvorak)
    // 2. Different language keyboards(Ex: Russian)

    if (event.keyCode === Constants.KeyCodes.COMPOSING[1]) {
        return false;
    }

    // checks for event.key for older browsers and also for the case of different English layout keyboards.
    if (typeof event.key !== 'undefined' && event.key !== 'Unidentified' && event.key !== 'Dead') {
        const isPressedByCode = event.key === key[0] || event.key === key[0].toUpperCase();
        if (isPressedByCode) {
            return true;
        }
    }

    // used for different language keyboards to detect the position of keys
    return event.keyCode === key[1];
}

export function getFullName(user) {
    if (user.first_name && user.last_name) {
        return user.first_name + ' ' + user.last_name;
    } else if (user.first_name) {
        return user.first_name;
    } else if (user.last_name) {
        return user.last_name;
    }

    return '';
}
export function handleFormattedTextClick(e) {
    const linkAttribute = e.target.getAttributeNode('data-link');

    if (linkAttribute) {
        const MIDDLE_MOUSE_BUTTON = 1;

        if (!(e.button === MIDDLE_MOUSE_BUTTON || e.altKey || e.ctrlKey || e.metaKey || e.shiftKey)) {
            e.preventDefault();

            window.WebappUtils.browserHistory.push(linkAttribute.value);
        }
    }
}

// return an object that contains the theme for react-select
export const getColorStyles = makeStyleFromTheme((theme) => {
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
