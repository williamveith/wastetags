document.addEventListener("DOMContentLoaded", () => {
    Array.from(document.getElementsByClassName("info-icon")).forEach((element) => {
        element.addEventListener("click", () => {
            alert(element.getElementsByClassName("help-text")[0].textContent)
        });
    });
});