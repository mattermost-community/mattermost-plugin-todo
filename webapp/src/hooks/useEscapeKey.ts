import React, {useEffect} from 'react';

export const useEscapeKey = (close: () => void) => {
    useEffect(() => {
        function handleKeypress(e: any) {
            if (e.key === 'Escape') {
                close();
            }
        }

        document.addEventListener('keyup', handleKeypress);

        return () => {
            document.removeEventListener('keyup', handleKeypress);
        };
    }, [close]);
};
