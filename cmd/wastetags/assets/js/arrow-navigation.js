const arrowNavigation = (event, elements) => {
    const elementIndex = elements.indexOf(document.activeElement);
    switch (event.key) {
        case "ArrowRight":
            if (elementIndex < elements.length - 1) {
                elements[elementIndex + 1].focus();
            } else {
                elements[0].focus();
            }
            break;
        case "ArrowLeft":
            if (elementIndex > 0) {
                elements[elementIndex - 1].focus();
            } else {
                elements[elements.length - 1].focus();
            }
            break;
        case "ArrowDown":
            elements[elements.length - 1].focus();
            break;
        case "ArrowUp":
            elements[0].focus();
            break;
        default:
            return;
    }
    event.preventDefault();
}

document.addEventListener("DOMContentLoaded", () => {
    const inputElements = Array.from(document.querySelectorAll("input"));
    inputElements.forEach((element) => {
        element.addEventListener("keydown", (event) => {
            arrowNavigation(event, inputElements);
        });
    });
});
