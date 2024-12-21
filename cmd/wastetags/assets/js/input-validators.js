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

export const casNumberValid = () => {
    const casValue = `${document.getElementById("cas1").value}${document.getElementById("cas2").value
        }`
        .split("")
        .reverse()
        .reduce(
            (accumulator, currentValue, index) =>
                accumulator + currentValue * (index + 1),
            0
        );

    const checkDigit = parseInt(document.getElementById("cas3").value, 10);

    if (casValue % 10 == checkDigit) {
        return true;
    } else {
        alert("Invalid CAS Number. Please fix it to submit this form.");
        return false;
    }
};