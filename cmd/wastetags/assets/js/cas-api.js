export const CasNumberValid = async (cas1, cas2, cas3, returnCasNumber = false) => {
    const casValue = `${cas1}${cas2}`
        .split("")
        .reverse()
        .reduce(
            (accumulator, currentValue, index) =>
                accumulator + currentValue * (index + 1),
            0
        );

    const isValid = casValue % 10 === parseInt(cas3, 10);
    const casNumber = isValid && returnCasNumber ? `${cas1}-${cas2}-${cas3}` : null;

    return { isValid, casNumber };
};

export const CheckDatabaseForCas = async (casNumber, chemicalName = "") => {
    try {
        const response = await fetch("/api/get-cas", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "cas": casNumber,
                "chem_name": chemicalName
            }),
        });
        return { "ok": response.ok, ...await response.json() };
    } catch (error) {
        return {
            "ok": false,
            "message": `Error: ${error}`,
            "cas": casNumber,
            "chem_name": chemicalName
        }
    }
};  