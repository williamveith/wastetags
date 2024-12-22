export const setDecimalPlaces = (element, numberOfDecimalPlaces = 2) => {
    const value = parseFloat(element.value);

    if (!isNaN(value)) {
        element.value = value.toFixed(numberOfDecimalPlaces);
    }
};