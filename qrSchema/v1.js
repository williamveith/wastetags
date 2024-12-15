const config = {
    topFields: [
        { jsonKey: "wasteTags", targetId: "CS_WASTE__REQUESTITEM__ITEMID", action: "pop" },
        { jsonKey: "chemName", targetId: "CS_WASTE__REQUESTITEM__NAME1" },
        { jsonKey: "location", targetId: "CS_WASTE__REQUESTITEM__LNAME" },
        { jsonKey: "contCount", targetId: "CS_WASTE__REQUESTITEM__CONTCOUNT" },
        { jsonKey: "contSize", targetId: "CS_WASTE__REQUESTITEM__CCONTSIZE" },
        { jsonKey: "sizeUnit", targetId: "CS_WASTE__REQUESTITEM__SIZEUNIT", type: "dropdown" },
        { jsonKey: "contType", targetId: "CS_WASTE__REQUESTITEM__CONTTYPE", type: "dropdown" },
        { jsonKey: "quantity", targetId: "CS_WASTE__REQUESTITEM__QUANTITY" },
        { jsonKey: "unit", targetId: "CS_WASTE__REQUESTITEM__UNIT", type: "dropdown" },
        { jsonKey: "physState", targetId: "CS_WASTE__REQUESTITEM__LSG" }
    ],
    componentFields: [
        { jsonKey: "component_name", targetId: "CS_WASTE__ADDCOMPOSITION__NAME1" },
        { jsonKey: "cas", targetId: "CS_WASTE__ADDCOMPOSITION__CAS" },
        { jsonKey: "percent", targetId: "CS_WASTE__ADDCOMPOSITION__PERCENT1" },
        { jsonKey: "unit", targetId: "CS_WASTE__ADDCOMPOSITION__UNIT", type: "dropdown" }
    ],
    modal: {
        id: "divWASTE__ADDCOMPOSITIONModal",
        addLinkId: "REQUESTITEM__SPREADITEM_AddLink",
        saveActionId: "SaveRecord"
    },
    finalSaveButtonId: "CMD_WASTE__REQUESTITEM__SaveRecord"
};

function setFieldValue(targetId, value, type = "input") {
    const element = document.getElementById(targetId);
    if (!element) return;

    if (type === "dropdown") {
        for (let option of element.options) {
            if (option.value === value) {
                element.value = option.value;
                break;
            }
        }
    } else {
        element.value = value;
    }
}

function fillInTop(data, config) {
    config.topFields.forEach(field => {
        const value = field.action === "pop" ? data[field.jsonKey].pop() : data.values[field.jsonKey];
        setFieldValue(field.targetId, value, field.type);
    });
}

let currentComponentIndex = 0;

function processComponent(component, config) {
    config.componentFields.forEach(field => {
        const value = component[field.jsonKey];
        setFieldValue(field.targetId, value, field.type);
    });

    document.querySelector(`[actionid="${config.modal.saveActionId}"]`).click();

    waitForModalToClose(() => {
        currentComponentIndex++;
        if (currentComponentIndex < components.length) {
            openModalAndProcessNextComponent(config);
        } else {
            document.getElementById(config.finalSaveButtonId).click();
        }
    });
}

function waitForModalToClose(callback) {
    const checkModalClose = setInterval(() => {
        const modal = document.getElementById(config.modal.id);
        if (!modal || modal.style.display === "none") {
            clearInterval(checkModalClose);
            callback();
        }
    }, 100);
}

function openModalAndProcessNextComponent(config) {
    const addLink = document.getElementById(config.modal.addLinkId);
    const observer = new MutationObserver((mutations, obs) => {
        const modal = document.getElementById(config.modal.id);
        if (modal && modal.style.display !== "none") {
            obs.disconnect();
            processComponent(components[currentComponentIndex], config);
        }
    });

    observer.observe(document.body, { childList: true, subtree: true });
    addLink.click();
}

function runInjector(data, config) {
    fillInTop(data, config);
    openModalAndProcessNextComponent(config);
    currentComponentIndex = 0;
}

const data = JSON.parse(prompt("Hazwaste String"));
runInjector(data, config);
