"use strict";

let productsData = [];
let currentPage = 1;
const itemsPerPage = 25;

const getData = async () => {
    try {
        const response = await fetch('https://inevitable-sukey-erwin-9f629ae2.koyeb.app/data/mitutoyo', {
        //const response = await fetch('http://127.0.0.1:8000/data/mitutoyo', {
            method: 'GET'
        });
        if (response.ok) {//dateMitutoyoData
            const jsonData = await response.json();
            productsData = jsonData[0].Data;
            let productDate = jsonData[0].UploadedAt
            formatISODateToHumanReadable(productDate)
            displayPage(currentPage);
            updatePaginationControls();

            console.log(productsData)
        } else {
            alert('Error al obtener los datos');
        }
    } catch (error) {
        console.error('Error:', error);
    }
};

const insertDataIntoTable = (dataSubset) => {
    const tableBody = document.querySelector('tbody');
    tableBody.innerHTML = '';
    dataSubset.forEach(dataRow => {
        const newRow = document.createElement('tr');
        newRow.classList.add('border-b', 'transition-colors');
        const columns = [
            'url', 'articlenumber', 'name', 'description', 'shortdescription', 
            'image', 'technicalimage', 'variants', 'leafletlinks', 
            'instructionpdflinks', 'accesories', 'imagelinks', 
            'youtubelinks', 'softwarelinks'
        ];
        columns.forEach(column => {
            const cell = document.createElement('td');
            cell.classList.add('p-4', 'align-middle');
            const cellData = dataRow.find(item => item.Key === column)?.Value || '';
            if (Array.isArray(cellData)) {
                cell.textContent = cellData.join(', ');
            } else {
                cell.textContent = cellData;
            }
            newRow.appendChild(cell);
        });
        tableBody.appendChild(newRow);
    });
};

const displayPage = (pageNumber) => {
    const startIndex = (pageNumber - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const pageData = productsData.slice(startIndex, endIndex);
    insertDataIntoTable(pageData);
    document.getElementById('pageNumber').textContent = pageNumber;
};
const updatePaginationControls = () => {
    const totalPages = Math.ceil(productsData.length / itemsPerPage);
    document.getElementById('prevPage').disabled = currentPage === 1;
    document.getElementById('nextPage').disabled = currentPage === totalPages;
};
document.getElementById('prevPage').addEventListener('click', () => {
    if (currentPage > 1) {
        currentPage--;
        displayPage(currentPage);
        updatePaginationControls();
    }
});
document.getElementById('nextPage').addEventListener('click', () => {
    const totalPages = Math.ceil(productsData.length / itemsPerPage);
    if (currentPage < totalPages) {
        currentPage++;
        displayPage(currentPage);
        updatePaginationControls();
    }
});
getData();

const searcherInput = document.getElementById('searcher');

searcherInput.addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
        const searchTerm = searcherInput.value.trim().toLowerCase();
        const filteredData = productsData.filter(item => {
            return item.some(field => {
                if (field.Key !== undefined && field.Value !== null) {
                    return field.Value.toString().toLowerCase().includes(searchTerm);
                }
                return false;
            });
        });
        insertDataIntoTable(filteredData);
    }
});

function resetTable() {
    insertDataIntoTable(productsData);
}

searcherInput.addEventListener('input', function() {
    if (searcherInput.value.trim() === '') {
        resetTable();
    }
});

function formatISODateToHumanReadable(isoDate) {
    const date = new Date(isoDate);
    const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: 'numeric',
        minute: 'numeric',
    };
    let dataDate = date.toLocaleString('es-ES', options);

    let mitutoyoDate = document.querySelector("#dateMitutoyoData");

    mitutoyoDate.innerHTML = dataDate + " datos de Mitutoyo"
}

// Function to export data to CSV
function exportToCSV(data, filename = 'data.csv') {
    const csv = Papa.unparse(data);
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

// Function to export data to Excel
function exportToExcel(data, filename = 'data.xlsx') {
    const ws = XLSX.utils.json_to_sheet(data);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
    XLSX.writeFile(wb, filename);
}

// Transform product data to a flat object format suitable for CSV/Excel export
function transformProductData(data) {
    return data.map(product => {
        const transformed = {};
        product.forEach(item => {
            transformed[item.Key] = Array.isArray(item.Value) ? item.Value.join(', ') : item.Value;
        });
        return transformed;
    });
}

// Get the export buttons
const csvExportBtn = document.getElementById('csv_export');
const excelExportBtn = document.getElementById('excel_export');

// Attach event listeners to export buttons
csvExportBtn.addEventListener('click', () => {
    const transformedData = transformProductData(productsData);
    exportToCSV(transformedData);
});

excelExportBtn.addEventListener('click', () => {
    const transformedData = transformProductData(productsData);
    exportToExcel(transformedData);
});
