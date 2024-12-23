function updateQrCode(clonedLabel, newInfo) {
    clonedImage = clonedLabel.querySelector("#qrCodeDataUri");
    clonedImage.src = newInfo["dataURI"];
    clonedImage.setAttribute('alt', newInfo["jsonContent"]);
    clonedLabel.querySelector("#qrCodeValue").value =
        newInfo["jsonContent"];
    const detailsElement = clonedLabel.querySelector(".details");
    const idMatch = detailsElement.textContent.match(/ID:\s*([^\s]+)/);
    if (idMatch) {
        detailsElement.innerHTML = detailsElement.innerHTML.replace(
            idMatch[1],
            newInfo["wasteTag"]
        );
    }
}

export const numberOfCopies = async (event) => {
    const desiredNumberOfCopies = parseInt(event.target.value, 10);
    if (desiredNumberOfCopies >= 1) {
        const tagElements = Array.from(
            document.getElementsByClassName("label")
        );
        const changeInNumberOfCopies =
            desiredNumberOfCopies - tagElements.length;

        if (changeInNumberOfCopies < 0) {
            while (tagElements.length > desiredNumberOfCopies) {
                const tagElement = tagElements.pop();
                tagElement.remove();
            }
        }

        for (let i = 0; i < changeInNumberOfCopies; i++) {
            const clonedTag = tagElements[0].cloneNode(true);

            try {
                const response = await fetch("/api/generate-qr-code", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: tagElements[0].querySelector("#qrCodeValue").value
                });

                if (!response.ok) {
                    throw new Error("Failed to generate QR code");
                }

                const result = await response.json();
                updateQrCode(clonedTag, result);
                tagElements[0].after(clonedTag);
            } catch (error) {
                console.error("Error:", error);
                alert("Error generating QR code");
            }
        }
    }
};