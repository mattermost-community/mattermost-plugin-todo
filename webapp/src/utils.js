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
