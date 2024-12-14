const v0 = {
    "components": [
        {
            "cas": "67-64-1",
            "chem_name": "CONTAMINATED DEBRIS W/ACETONE 33%,ISOPROPANOL 34%,METHANOL 33%",
            "component_name": "ACETONE",
            "component_order": 0,
            "id": 97,
            "percent": "33"
        },
        {
            "cas": "67-63-0",
            "chem_name": "CONTAMINATED DEBRIS W/ACETONE 33%,ISOPROPANOL 34%,METHANOL 33%",
            "component_name": "ISOPROPANOL",
            "component_order": 1,
            "id": 98,
            "percent": "33"
        },
        {
            "cas": "67-56-1",
            "chem_name": "CONTAMINATED DEBRIS W/ACETONE 33%,ISOPROPANOL 34%,METHANOL 33%",
            "component_name": "METHANOL",
            "component_order": 2,
            "id": 99,
            "percent": "34"
        }
    ],
    "values": {
        "chemName": "CONTAMINATED DEBRIS W/ACETONE 33%,ISOPROPANOL 34%,METHANOL 33%",
        "contCount": "1",
        "contSize": "55.00",
        "contType": "DF Drum Poly Open",
        "location": "MER 1.010",
        "physState": "SOLID",
        "quantity": "1.50",
        "sizeUnit": "G GALLONS",
        "unit": "P POUNDS"
    },
    "wasteTags": [
        "QWyo4HnSzSkrQH6nhzEh"
    ]
}

let { wasteTags, values, components } = JSON.parse(prompt("Hazwaste String"));

function fillInTop() {
    document.getElementById('CS_WASTE__REQUESTITEM__ITEMID').value = wasteTags.pop();

    document.getElementById('CS_WASTE__REQUESTITEM__NAME1').value = values.chemName;

    document.getElementById('CS_WASTE__REQUESTITEM__LNAME').value = values.location;

    // Set Quantity to 1
    document.getElementById('CS_WASTE__REQUESTITEM__CONTCOUNT').value = values.contCount;

    // Set size to 1.00
    document.getElementById('CS_WASTE__REQUESTITEM__CCONTSIZE').value = values.contSize;

    // Set units to G GALLONS
    const unitsDropdown = document.getElementById('CS_WASTE__REQUESTITEM__SIZEUNIT');
    for (let option of unitsDropdown.options) {
        if (option.value === values.sizeUnit) {
            unitsDropdown.value = option.value;
            break;
        }
    }

    // Set type to GB Glass Bottle
    const typeDropdown = document.getElementById('CS_WASTE__REQUESTITEM__CONTTYPE');
    for (let option of typeDropdown.options) {
        if (option.value === values.contType) {
            typeDropdown.value = option.value;
            break;
        }
    }

    // Set Quantity to 1
    document.getElementById('CS_WASTE__REQUESTITEM__QUANTITY').value = values.quantity;

    // Set Units to G GALLONS
    const unitsDropdown2 = document.getElementById('CS_WASTE__REQUESTITEM__UNIT');
    // Iterate over the options to find the one with value "G GALLONS" and select it
    for (let option of unitsDropdown2.options) {
        if (option.value === values.unit) {
            unitsDropdown2.value = option.value;
            break; // Stop the loop once the correct option is found and selected
        }
    }

    document.getElementById('CS_WASTE__REQUESTITEM__LSG').value = values.physState;
}
let currentComponentIndex = 0;

function processComponent(component) {
    document.getElementById('CS_WASTE__ADDCOMPOSITION__NAME1').value = component.component_name;
    document.getElementById('CS_WASTE__ADDCOMPOSITION__CAS').value = component.cas;
    document.getElementById('CS_WASTE__ADDCOMPOSITION__PERCENT1').value = component.percent;

    let unitSelect = document.getElementById('CS_WASTE__ADDCOMPOSITION__UNIT');
    Array.from(unitSelect.options).forEach(option => {
        if (option.value === component.unit) {
            unitSelect.value = option.value;
        }
    });

    document.querySelector('[actionid="SaveRecord"]').click();

    waitForModalToClose(() => {
        currentComponentIndex++;
        if (currentComponentIndex < components.length) {
            openModalAndProcessNextComponent();
        } else {
            document.getElementById("CMD_WASTE__REQUESTITEM__SaveRecord").click();
        }
    });
}

function waitForModalToClose(callback) {
    const checkModalClose = setInterval(() => {
        const modal = document.getElementById('divWASTE__ADDCOMPOSITIONModal');
        if (!modal || modal.style.display === 'none') {
            clearInterval(checkModalClose);
            callback();
        }
    }, 100);
}

function openModalAndProcessNextComponent() {
    const addLink = document.getElementById('REQUESTITEM__SPREADITEM_AddLink');
    const observer = new MutationObserver((mutations, obs) => {
        const modal = document.getElementById('divWASTE__ADDCOMPOSITIONModal');
        if (modal && modal.style.display !== 'none') {
            obs.disconnect();
            processComponent(components[currentComponentIndex]);
        }
    });

    observer.observe(document.body, { childList: true, subtree: true });
    addLink.click();
}

function run() {
    fillInTop();
    openModalAndProcessNextComponent();
    currentComponentIndex = 0;
}

run();