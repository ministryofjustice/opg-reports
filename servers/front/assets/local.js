document.addEventListener("DOMContentLoaded", ready);

function ready(){
    tableToggler()
    tableChoiceFilter()
    tableSorter()
    generateFormFilterOptions()
    tableOptionFilters()
}

// ==========================
// TABLE ITEM VISIBILITY
// ==========================
// Show / hide elements within table rows when + icons clicked
function tableToggler() {
    [].forEach.call( document.querySelectorAll( ".js-table-toggler" ), function ( ele ) {

        ele.addEventListener('click', function(eve){
            var toggleDisplay = eve.target.dataset.toggle;

            [].forEach.call( document.querySelectorAll('.'+toggleDisplay), function(info) {
                if (info.style.display == 'none' || info.style.display == "") {
                    info.style.display = 'table-cell'
                } else {
                    info.style.display = 'none'
                }
                return false
            })

            return false
        }, false)
    } )
}


// ==========================
// TABLE FILTERS
// ==========================
// Filter table based on radio button selections
function tableChoiceFilter() {
    let choices = document.querySelectorAll(".js-table-choice-filter")
    choices.forEach(ele => {
        ele.addEventListener("change", function(eve){
            let allQ = eve.target.dataset.all;
            let showQ = eve.target.dataset.show;
            let selector = allQ + '['+showQ+']'

            if (showQ != undefined){
                document.querySelectorAll(allQ).forEach(r => {r.style.display = 'none'})
                document.querySelectorAll(selector).forEach(r => {r.style.display = 'table-row'})
            } else {
                document.querySelectorAll(allQ).forEach(r => {r.style.display = 'table-row'})
            }

        })
    } )
}

function tableFilterByAll(tableData, filters){

    tableData.forEach(row => {
        let show = true;
        let cols = row.querySelectorAll('td,th')
        // console.log(cols);
        filters.forEach(filter => {
            let colNum = filter.dataset.col;
            let sel = filter.querySelector('select');
            let val = sel.value;
            let comp = cols[colNum].textContent;
            if (val != "all" && comp != val) {
                show = false;
            }
        })
        if (show == true) {
            row.style.display = 'table-row';
        } else {
            row.style.display = 'none';
        }
    })
}

function tableOptionFilters() {
    let formSelector = `.js-table-filter-options`
    let filterSelector = `.js-table-filter-select`
    let forms = document.querySelectorAll(formSelector)
    forms.forEach(f => {
        let done = f.dataset.filterActive;
        let tableData = document.querySelectorAll(f.dataset.filterrows);
        if (done !== true){
            let filters = f.querySelectorAll(filterSelector)

            filters.forEach(filter => {
                let sel = filter.querySelector('select');
                sel.addEventListener("change", e => tableFilterByAll(tableData, filters))
            })
        }

    })

}
// generateFormFilterOptions adds all the unique values from a column into a select within a form
// allowing dynamic data sets pulled from the page content
function generateFormFilterOptions() {
    let formSelector = `.js-table-filter-options-generate`
    let filterSelect = `.js-table-generate-options`
    let generators = document.querySelectorAll(formSelector)

    generators.forEach(block => {
        let done = block.dataset.generated
        if (done !== true) {
            let dataSource = block.dataset.optionssource;
            let filterWrappers = block.querySelectorAll(filterSelect)
            filterWrappers.forEach(wrapper => {
                let col = wrapper.dataset.col;
                let sel = wrapper.querySelector('select');
                let selectedCols = document.querySelectorAll(dataSource + `[data-col="${col}"]`);
                let colValues = [];
                selectedCols.forEach(c => { colValues.push(c.textContent) })
                let unique = [...new Set(colValues)];
                unique = unique.sort()
                // add the data
                unique.forEach(val => {
                    sel.add(new Option(val, val))
                })
                wrapper.classList.add(`js-table-filter-select`)
            })
            block.dataset.generated = true

        }
    })

}


// ==========================
// TABLE SORTING
// ==========================
// Compare values - strip $ and ,
// - try sorting as a number via float parsing
function tableSortCompareValues(a, b) {
    aVal = a.replace("$", "").replace(",", "");
    bVal = b.replace("$", "").replace(",", "");
    comp = 0
    if (parseFloat(aVal)){
        aVal = parseFloat(aVal);
    }
    if (parseFloat(bVal)){
        bVal = parseFloat(bVal);
    }

    if (aVal < bVal) {
        comp = -1
    } else if (aVal > bVal) {
        comp = 1
    }
    return comp

}
// Handle the sort trigger event and reorder the table body
function tableSortHandler(thead, tbody, th, colNum) {
    let defaultDir = "asc";
    let rows = Array.from(tbody.querySelectorAll(`tr`));
    let sortDir = th.dataset.sortdir;
    if (sortDir == undefined){
        sortDir = defaultDir
    }
    if (sortDir == "asc") {
        sortDir = "desc"
    } else {
        sortDir = "asc"
    }
    // and then just... sort the rows:
    rows.sort( (r1,r2) => {
        // get each row's relevant column
        var t1 = r1.querySelectorAll('td,th')[colNum];
        var t2 = r2.querySelectorAll('td,th')[colNum];
        var c1 = t1.textContent
        var c2 = t2.textContent
        comp = tableSortCompareValues(c1, c2)
        return comp
    });
    if (sortDir == "desc") {
        rows = rows.reverse()
    } else {
        tableSortDesc = true
    }

    thead.querySelectorAll(`th`).forEach(t => {t.classList.remove('sorted')})
    th.classList.add('sorted');
    th.dataset.sortdir = sortDir;
    // and then the magic part that makes the sorting appear on-page:
    rows.forEach(row => tbody.appendChild(row));
}
// Main table sorter call
function tableSorter() {
    tables = document.querySelectorAll("table.js-table-sorter")
    tables.forEach( table => {
        thead = table.querySelector(`thead`);
        tbody = table.querySelector(`tbody`);

        thead.querySelectorAll(`th`).forEach((th, position) => {
            th.addEventListener(`click`, evt => tableSortHandler(thead, tbody, th, position));
        });
        // trigger a click on the last column
        var sel = ".govuk-table__head th:last-of-type"
        table.querySelector(sel).click()
    })

}
