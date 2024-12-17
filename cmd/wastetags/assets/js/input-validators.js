export const onlyNumbersAllowed = (event) => {
    event.target.value = event.target.value.replace(/(^0|[^0-9])/g, "");
    event.target.setSelectionRange(
        event.target.value.length,
        event.target.value.length
    );
};

export const onlyNumberDigitsAllowed = (event) => {
    event.target.value = event.target.value.replace(/[^0-9]/g, "");
    event.target.setSelectionRange(
        event.target.value.length,
        event.target.value.length
    );
};

export const upToTwoDecimals = (event) => {
    event.target.value = event.target.value.replace(
        /[^0-9.]|(\.\d{2})\d+|(\..*)\./g,
        (match, group1) => group1 || ""
    );
    event.target.setSelectionRange(
        event.target.value.length,
        event.target.value.length
    );
};
