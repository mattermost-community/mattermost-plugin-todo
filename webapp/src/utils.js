export function canDelete(myList, foreignList) {
    return myList === 'my' || myList === 'in' || foreignList === 'in';
}

export function canComplete(myList) {
    return myList === 'my' || myList === 'in';
}

export function canEnqueue(myList) {
    return myList === 'in';
}